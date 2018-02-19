package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type hijacker struct{ http.ResponseWriter }

func (c *hijacker) WriteHeader(statusCode int) {}

func statusCode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if st := r.URL.Query().Get("status"); st != "" {
			if status, _ := strconv.Atoi(st); status >= 100 && status < 600 {
				q := r.URL.Query()
				q.Del("status")
				r.URL.RawQuery = q.Encode()

				// need to copy the body before writing the
				// status code
				var buf bytes.Buffer
				if _, err := io.Copy(&buf, r.Body); err != nil {
					errorJSON(w, "Unable to clone request body: %s", err.Error())
				}

				// write status code
				w.WriteHeader(status)

				// then re-set request body, since it was closed
				// on the WriteHeader() call
				r.Body = ioutil.NopCloser(&buf)
				defer r.Body.Close()

				next.ServeHTTP(&hijacker{w}, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
