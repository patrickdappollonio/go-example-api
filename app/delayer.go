package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func delayer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if str := r.URL.Query().Get("delay"); str != "" {
			if n, _ := strconv.Atoi(str); n > 0 && n <= 50 {
				q := r.URL.Query()
				q.Del("delay")
				r.URL.RawQuery = q.Encode()

				if sl, err := time.ParseDuration(fmt.Sprintf("%ds", n)); err == nil {
					time.Sleep(sl)
				}

			}
		}

		next.ServeHTTP(w, r)
	})
}
