package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/", nocontent)
	r.With(statusCode, delayer).HandleFunc("/debug", debug)

	return r
}
