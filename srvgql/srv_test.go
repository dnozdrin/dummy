package srvgql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv HTTPSrv

	srv.readiness = false
	assert.Equal(t, "gql service is't ready yet", srv.HealthCheck().Error())
}
