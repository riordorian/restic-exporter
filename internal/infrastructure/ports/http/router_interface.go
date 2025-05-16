package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type RouterInterface interface {
	Handle(path string, handler http.Handler) *mux.Route
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
