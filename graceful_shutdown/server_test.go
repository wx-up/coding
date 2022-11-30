package service

import (
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	srv1 := NewServer("web", "localhost:8081")
	srv1.Handle("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello web"))
	})

	srv2 := NewServer("admin", "localhost:8082")
	srv2.Handle("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello admin"))
	})

	app := NewApp([]*Server{srv1, srv2})
	app.StartAndServe()
}
