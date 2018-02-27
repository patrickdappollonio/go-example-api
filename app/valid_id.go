package app

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

var (
	reIsValidParam = regexp.MustCompile(`^[\w|\-]+$`)
	reIsValidNumID = regexp.MustCompile(`^\d+$`)
)

const ID_KEY = "__routed_id"

func validID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			id = chi.URLParam(r, "id")
			ct = r.Header.Get("Content-Type")
		)

		if !reIsValidParam.MatchString(id) {
			w.WriteHeader(http.StatusNotFound)
			content := fmt.Sprintf("Invalid ID: %s", id)

			if strings.HasPrefix(ct, "application/json") {
				errorJSON(w, content)
				return
			}

			fmt.Fprint(w, content)
			return
		}

		ctx := context.WithValue(r.Context(), ID_KEY, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validNumericID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			id = chi.URLParam(r, "id")
			ct = r.Header.Get("Content-Type")
		)

		numid, _ := strconv.Atoi(id)

		if !reIsValidNumID.MatchString(id) || numid <= 0 {
			w.WriteHeader(http.StatusNotFound)
			content := fmt.Sprintf("Invalid ID: %s", id)

			if strings.HasPrefix(ct, "application/json") {
				errorJSON(w, content)
				return
			}

			fmt.Fprint(w, content)
			return
		}

		ctx := context.WithValue(r.Context(), ID_KEY, numid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
