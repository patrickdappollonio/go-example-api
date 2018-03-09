package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func Router() http.Handler {
	setupTemplate()
	generateRandomData()

	r := chi.NewRouter()
	r.Use(canonical("request.tools"))
	r.Use(disableAccess)
	r.Use(cors.Default().Handler)

	r.Get("/", home)
	r.Get("/favicon.ico", favicon)
	r.With(statusCode, delayer).HandleFunc("/debug", debug)

	r.With(statusCode, delayer).Post("/save", save)
	r.With(statusCode, delayer, validID).Get("/get/{id}", get)

	r.Get("/inspector", bincreate)
	r.With(validID, loadPrevious).Get("/inspector/{id}", binget)
	r.Handle("/r", http.RedirectHandler("/inspector", http.StatusFound))
	r.With(safeVerbs, validID, loadPrevious).HandleFunc("/r/{id}", binsave)
	r.With(validID).Get("/inspector/{id}/delete", bindelete)

	r.Route("/f", func(r chi.Router) {
		r.Use(asJSON)
		r.Get("/", fakerest)
		r.Get("/users", fakeusers)
		r.With(validNumericID).Get("/users/{id}", fakesingleuser)
		r.Get("/posts", fakeposts)
		r.With(validNumericID).Get("/posts/{id}", fakesinglepost)
		r.Get("/domains", fakedomains)
		r.With(validNumericID).Get("/domains/{id}", fakesingledomain)
		r.Get("/products", fakeproducts)
		r.With(validNumericID).Get("/products/{id}", fakesingleproduct)
	})

	return r
}
