package srvgql

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/99designs/gqlgen/cmd" // nolint
	gqlhandler "github.com/99designs/gqlgen/handler"
	log "github.com/sirupsen/logrus"
)

func New(port int, service Service) (*HTTPSrv, error) {
	// build http server
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	// build HTTPSrv
	var srv HTTPSrv
	srv.service = service
	srv.setupHTTP(&httpSrv)

	return &srv, nil
}

type HTTPSrv struct {
	http      *http.Server
	runErr    error
	readiness bool
	service   Service
}

func (s *HTTPSrv) setupHTTP(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

func (s *HTTPSrv) buildHandler() http.Handler {
	handler := http.NewServeMux()
	// path -> handlers

	// build gql resolver
	resolver := Resolver{
		service: s.service,
	}

	graphQLHandler := gqlhandler.GraphQL(
		NewExecutableSchema(
			Config{Resolvers: &resolver},
		),
	)

	handler.Handle("/gql/query", graphQLHandler)

	// ==============
	return handler
}

func (s *HTTPSrv) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("gql service: begin run")

	go func() {
		defer wg.Done()
		log.Debug("gql service: addr=", s.http.Addr)
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info("gql service: end run > ", err)
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Info("gql service shutdown (", err, ")")
		}
	}()

	s.readiness = true
}

func (s *HTTPSrv) HealthCheck() error {
	if !s.readiness {
		return errors.New("gql service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("gql service: run issue")
	}
	if s.service == nil || s.service.HealthCheck() != nil {
		return errors.New("gql service: storage issue")
	}
	return nil
}
