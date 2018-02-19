package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Default().Handler)

	r.Get("/", home)
	r.With(statusCode, delayer).HandleFunc("/debug", debug)

	r.With(statusCode, delayer).Post("/save", save)
	r.With(statusCode, delayer).Get("/get/{id}", get)

	return r
}
