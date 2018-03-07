package app

import (
	"fmt"
	"net/http"
	"strings"
)

var disabledLocations = map[string]struct{}{
	"ru":             struct{}{},
	"ru#kya#norilsk": struct{}{}, // this location has been hitting URLs like crazy
}

func disableAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		country := getHeaderLower(r, "X-Appengine-Country")

		if _, found := disabledLocations[country]; found {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		region := getHeaderLower(r, "X-Appengine-Region")
		city := getHeaderLower(r, "X-Appengine-City")

		hash := fmt.Sprintf("%s#%s#%s", country, region, city)

		if _, found := disabledLocations[hash]; found {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getHeaderLower(r *http.Request, key string) string {
	return strings.ToLower(r.Header.Get(key))
}
