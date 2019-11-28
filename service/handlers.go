package service

import (
	"net/http"
)

func (s *Service) hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello"))
}
