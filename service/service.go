package service

import (
	"github.com/akhripko/dummy/models"
	"github.com/pkg/errors"
)

type Storage interface {
	Check() error
	Hello(name string) (*models.HelloMessage, error)
}

type Cache interface {
	Check() error
	Read(name string) (string, error)
	WriteTTL(name, msg string, ttl int) error
}

func New(storage Storage, cache Cache) (*Service, error) {
	// build service
	srv := Service{
		storage:   storage,
		cache:     cache,
		readiness: true,
	}

	return &srv, nil
}

type Service struct {
	storage   Storage
	cache     Cache
	readiness bool
}

func (s *Service) HealthCheck() error {
	if !s.readiness {
		return errors.New("service is't ready yet")
	}
	if s.storage == nil || s.storage.Check() != nil {
		return errors.New("service: storage issue")
	}
	if s.cache == nil || s.cache.Check() != nil {
		return errors.New("service: cache issue")
	}
	return nil
}
