package app

import (
	"fmt"
	"net/http"
	"strings"
)

func canonical(hostname string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Host, ":") || r.Host == "127.0.0.1" || r.Host == "localhost" {
				next.ServeHTTP(w, r)
				return
			}

			if r.Host != hostname {
				prefix := "http"
				if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
					prefix += "s"
				}

				query := ""
				if q := r.URL.RawQuery; q != "" {
					query += "?" + q
				}

				http.Redirect(w, r, fmt.Sprintf("%s://%s%s%s", prefix, hostname, r.URL.Path, query), http.StatusMovedPermanently)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
