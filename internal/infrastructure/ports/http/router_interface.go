package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

// TODO:  add mux.route interface

type RouterInterface interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
