package healthcheck

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type service struct {
	http *http.Server
}

func New(port int, healthChecks ...func() error) Service {
	return &service{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: buildHandler(healthChecks),
		},
	}
}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("healthcheck service: begin run")

	go func() {
		defer wg.Done()
		log.Debug("healthcheck service addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil {
			log.Error("healthcheck service end run:", err.Error())
			return
		}
		log.Info("healthcheck service: end run")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("healthcheck service shutdown error:", err.Error())
		}
	}()
}

func buildHandler(healthChecks []func() error) http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/version", serveVersion)
	var checks = func(w http.ResponseWriter, _ *http.Request) { serveCheck(w, healthChecks) }
	handler.HandleFunc("/", checks)
	handler.HandleFunc("/health", checks)
	handler.HandleFunc("/ready", checks)
	return handler
}

func writeFile(file string, response http.ResponseWriter) {
	if data, err := ioutil.ReadFile(file); err == nil { // nolint
		response.WriteHeader(http.StatusOK)
		response.Write(data) // nolint
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

func serveCheck(w http.ResponseWriter, checks []func() error) {
	writtenHeader := false
	for _, check := range checks {
		if err := check(); err != nil {
			if !writtenHeader {
				w.WriteHeader(http.StatusInternalServerError)
				writtenHeader = true
			}
			w.Write([]byte(err.Error())) // nolint
			w.Write([]byte("\n\n"))      // nolint
		}
	}

	if !writtenHeader {
		w.WriteHeader(http.StatusNoContent)
	}
}
