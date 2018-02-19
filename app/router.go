package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/", home)
	r.With(statusCode, delayer).HandleFunc("/debug", debug)

	return r
}
