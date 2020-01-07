package srvgrpc

import (
	"context"
	"testing"

	"github.com/akhripko/dummy/api"
	"google.golang.org/grpc"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	mock "github.com/akhripko/dummy/srvgrpc/mock"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv Srv

	srv.readiness = false
	assert.Equal(t, "grpc service is't ready yet", srv.HealthCheck().Error())
}

func TestService_StatusCheckRunErr(t *testing.T) {
	var srv Srv

	srv.readiness = true

	c := gomock.NewController(t)
	service := mock.NewMockService(c)
	service.EXPECT().HealthCheck().DoAndReturn(func() error { return nil }).Times(1)
	srv.service = service
	assert.Nil(t, srv.HealthCheck())

	srv.runErr = errors.New("some run error")
	assert.Equal(t, "grpc service: run issue", srv.HealthCheck().Error())
}

func TestService_StatusCheckDB(t *testing.T) {
	var srv Srv

	srv.readiness = true

	c := gomock.NewController(t)
	service := mock.NewMockService(c)
	service.EXPECT().HealthCheck().DoAndReturn(func() error { return errors.New("some error") }).Times(1)
	srv.service = service

	assert.Equal(t, "grpc srv: service issue", srv.HealthCheck().Error())
}

func newClient(ctx context.Context, target string) (api.DummyServiceClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return api.NewDummyServiceClient(conn), nil
}
