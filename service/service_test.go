package service

import (
	"testing"

	mock "github.com/akhripko/dummy/service/mock"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv Service

	srv.readiness = false
	assert.Equal(t, "service is't ready yet", srv.HealthCheck().Error())
}

func TestService_StatusCheckRunErr(t *testing.T) {
	var s Service

	s.readiness = true

	c := gomock.NewController(t)
	storage := mock.NewMockStorage(c)
	storage.EXPECT().Check().DoAndReturn(func() error { return nil }).Times(1)
	s.storage = storage

	cache := mock.NewMockCache(c)
	cache.EXPECT().Check().DoAndReturn(func() error { return nil }).Times(1)
	s.cache = cache

	assert.Nil(t, s.HealthCheck())
}

func TestService_StatusCheckStorage(t *testing.T) {
	var srv Service

	srv.readiness = true

	c := gomock.NewController(t)

	storage := mock.NewMockStorage(c)
	storage.EXPECT().Check().DoAndReturn(func() error { return errors.New("some db error") }).Times(1)
	srv.storage = storage

	cache := mock.NewMockCache(c)
	cache.EXPECT().Check().DoAndReturn(func() error { return nil }).Times(1)
	srv.cache = cache

	assert.Equal(t, "service: storage issue", srv.HealthCheck().Error())
}

func TestService_StatusCheckCache(t *testing.T) {
	var srv Service

	srv.readiness = true

	c := gomock.NewController(t)

	storage := mock.NewMockStorage(c)
	storage.EXPECT().Check().DoAndReturn(func() error { return nil }).Times(1)
	srv.storage = storage

	cache := mock.NewMockCache(c)
	cache.EXPECT().Check().DoAndReturn(func() error { return errors.New("some cache error") }).Times(1)
	srv.cache = cache

	assert.Equal(t, "service: cache issue", srv.HealthCheck().Error())
}
