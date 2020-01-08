package srvhttp

import "github.com/akhripko/dummy/models"

type Service interface {
	Hello(name string) (*models.HelloMessage, error)
}
