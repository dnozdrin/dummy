package srvgql

import "github.com/akhripko/dummy/models"

type Service interface {
	HealthCheck() error
	Hello(name string) (*models.HelloMessage, error)
}
