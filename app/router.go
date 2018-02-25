package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func Router() http.Handler {
	setupTemplate()

	r := chi.NewRouter()
	r.Use(canonical("request.tools"))
	r.Use(cors.Default().Handler)

	r.Get("/", home)
	r.With(statusCode, delayer).HandleFunc("/debug", debug)

	r.With(statusCode, delayer).Post("/save", save)
	r.With(statusCode, delayer, validID).Get("/get/{id}", get)

	r.Get("/inspector", bincreate)
	r.With(validID, loadPrevious).Get("/inspector/{id}", binget)
	r.With(safeVerbs, validID, loadPrevious).HandleFunc("/r/{id}", binsave)
	r.With(validID).Get("/inspector/{id}/delete", bindelete)

	return r
}
