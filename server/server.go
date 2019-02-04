package server

import "net/http"

type CustomHandler struct {
}

func (s *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Hello World"))
	}
}
