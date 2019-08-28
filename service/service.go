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

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	HealthCheck() error
}

func New(port int, db DB, cache Cache) Service {
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	var srv service
	srv.setupHTTP(&httpSrv)
	srv.db = db
	srv.cache = cache

	return &srv
}

type service struct {
	http      *http.Server
	runErr    error
	readiness bool
	db        DB
	cache     Cache
}

func (s *service) setupHTTP(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

func (s *service) buildHandler() http.Handler {
	r := mux.NewRouter()
	// path -> handlers

	// hello request
	hello := metrics.Counter(metrics.HelloRequestCounts, s.hello)
	hello = metrics.Timer(metrics.HelloRequestTiming, hello)
	r.HandleFunc("/hello", hello).Methods("GET")

	// ==============
	return r
}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("service: begin run")

	go func() {
		defer wg.Done()
		log.Debug("service addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil {
			s.runErr = err
			log.Error("service end run:", err)
			return
		}
		log.Info("service: end run")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("service shutdown error:", err)
		}
	}()

	s.readiness = true
}

func (s *service) HealthCheck() error {
	if !s.readiness {
		return errors.New("service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("run service issue")
	}
	if s.db == nil || s.db.Ping() != nil {
		return errors.New("db issue")
	}
	if s.cache == nil || s.cache.Ping() != nil {
		return errors.New("cache issue")
	}
	return nil
}
