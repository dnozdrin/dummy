package service

import (
	"testing"

	"github.com/akhripko/dummy/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv service

	srv.readiness = false
	assert.Equal(t, "service is't ready yet", srv.HealthCheck().Error())
}

func TestService_StatusCheckRunErr(t *testing.T) {
	var srv service

	srv.readiness = true

	c := gomock.NewController(t)
	db := mock_service.NewMockStorage(c)
	db.EXPECT().Ping().DoAndReturn(func() error { return nil }).Times(1)
	cache := mock_service.NewMockCache(c)
	cache.EXPECT().Ping().DoAndReturn(func() error { return nil }).Times(1)
	srv.db = db
	srv.cache = cache
	assert.Nil(t, srv.HealthCheck())

	srv.runErr = errors.New("some run error")
	assert.Equal(t, "run service issue", srv.HealthCheck().Error())
}

func TestService_StatusCheckDB(t *testing.T) {
	var srv service

	srv.readiness = true

	c := gomock.NewController(t)
	db := mock_service.NewMockStorage(c)
	db.EXPECT().Ping().DoAndReturn(func() error { return errors.New("some db error") }).Times(1)
	cache := mock_service.NewMockCache(c)
	cache.EXPECT().Ping().DoAndReturn(func() error { return nil }).Times(1)
	srv.db = db
	srv.cache = cache

	assert.Equal(t, "db issue", srv.HealthCheck().Error())
}

func TestService_StatusCheckCache(t *testing.T) {
	var srv service

	srv.readiness = true

	c := gomock.NewController(t)
	db := mock_service.NewMockStorage(c)
	db.EXPECT().Ping().DoAndReturn(func() error { return nil }).Times(1)
	cache := mock_service.NewMockCache(c)
	cache.EXPECT().Ping().DoAndReturn(func() error { return errors.New("some cache error") }).Times(1)
	srv.db = db
	srv.cache = cache

	assert.Equal(t, "cache issue", srv.HealthCheck().Error())
}
