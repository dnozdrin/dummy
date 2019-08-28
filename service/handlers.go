package service

import (
	"net/http"
)

func (s *service) hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello"))
}
