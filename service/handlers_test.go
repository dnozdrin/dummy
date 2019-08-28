package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hello(t *testing.T) {
	srv := service{}

	// check http
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	//req.Header.Add("X-Forwarded-Proto", "http")

	handler := srv.buildHandler()

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "hello", string(rr.Body.Bytes()))
}
