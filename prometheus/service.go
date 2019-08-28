package prometheus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	HealthCheck() error
}

type service struct {
	http      *http.Server
	runErr    error
	readiness bool
}

func New(port int) Service {
	return &service{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler(),
		},
	}
}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("prometheus service: begin run")

	go func() {
		defer wg.Done()
		log.Debug("prometheus service addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil {
			s.runErr = err
			log.Error("prometheus service end run:", err.Error())
			return
		}
		log.Info("prometheus service: end run")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("prometheus service shutdown error:", err.Error())
		}
	}()

	s.readiness = true
}

func handler() http.Handler {
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	return handler
}

func (s *service) HealthCheck() error {
	if !s.readiness {
		return errors.New("prometheus service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("run prometheus service issue")
	}
	return nil
}
