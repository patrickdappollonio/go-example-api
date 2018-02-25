package app

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func cleanup(m interface{}) map[string][]string {
	h := make(map[string][]string)

	switch a := m.(type) {
	case http.Header:
		h = a
	case url.Values:
		h = a
	default:
		return nil
	}

	local := make(map[string][]string)

	for k, v := range h {
		if strings.HasPrefix(k, "X-Appengine") {
			continue
		}

		if strings.HasPrefix(k, "X-Goog-Cloud-Shell") {
			continue
		}

		if k == "X-Cloud-Trace-Context" {
			continue
		}

		key := strings.ToLower(k)
		values := make([]string, 0, len(v))

		for _, value := range v {
			if key == "via" && fmt.Sprintf("%v", value) == "1.1 google" {
				continue
			}

			values = append(values, value)
		}

		local[key] = values
	}

	return local
}
