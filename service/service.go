package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/akhripko/dummy/metrics"
	"github.com/gorilla/mux"
)

func New(port int, db DB, cache Cache) *Service {
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	var srv Service
	srv.setupHTTP(&httpSrv)
	srv.db = db
	srv.cache = cache

	return &srv
}

type Service struct {
	http      *http.Server
	runErr    error
	readiness bool
	db        DB
	cache     Cache
}

func (s *Service) setupHTTP(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

func (s *Service) buildHandler() http.Handler {
	r := mux.NewRouter()
	// path -> handlers

	// hello request
	hello := metrics.Counter(metrics.HelloRequestCounts, s.hello)
	hello = metrics.Timer(metrics.HelloRequestTiming, hello)
	r.HandleFunc("/hello", hello).Methods("GET")

	// ==============
	return r
}

func (s *Service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("Service: begin run")

	go func() {
		defer wg.Done()
		log.Debug("Service addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil {
			s.runErr = err
			log.Error("Service end run:", err)
			return
		}
		log.Info("Service: end run")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("Service shutdown error:", err)
		}
	}()

	s.readiness = true
}

func (s *Service) HealthCheck() error {
	if !s.readiness {
		return errors.New("Service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("run Service issue")
	}
	if s.db == nil || s.db.Ping() != nil {
		return errors.New("db issue")
	}
	if s.cache == nil || s.cache.Ping() != nil {
		return errors.New("cache issue")
	}
	return nil
}
