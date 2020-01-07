package srvgql

import (
	"testing"

	mock "github.com/akhripko/dummy/srvhttp/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv HTTPSrv

	srv.readiness = false
	assert.Equal(t, "gql service is't ready yet", srv.HealthCheck().Error())
}

func TestService_StatusCheckRunErr(t *testing.T) {
	var srv HTTPSrv

	srv.readiness = true

	c := gomock.NewController(t)
	service := mock.NewMockService(c)
	service.EXPECT().HealthCheck().DoAndReturn(func() error { return nil }).Times(1)
	srv.service = service
	assert.Nil(t, srv.HealthCheck())

	srv.runErr = errors.New("some run error")
	assert.Equal(t, "gql service: run issue", srv.HealthCheck().Error())
}

func TestService_StatusCheckDB(t *testing.T) {
	var srv HTTPSrv

	srv.readiness = true

	c := gomock.NewController(t)
	service := mock.NewMockService(c)
	service.EXPECT().HealthCheck().DoAndReturn(func() error { return errors.New("some service error") }).Times(1)
	srv.service = service

	assert.Equal(t, "gql service: storage issue", srv.HealthCheck().Error())
}
