package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type debugresponse struct {
	Headers map[string]interface{} `json:"headers"`
	Body    interface{}            `json:"body"`
	Query   map[string]interface{} `json:"query,omitempty"`
}

func debug(w http.ResponseWriter, r *http.Request) {
	data := debugresponse{
		Headers: cleanup(r.Header),
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		errorJSON(w, "Unable to read request body: %s", err.Error())
		return
	}

	if buf.Len() > 0 {
		var parsed map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &parsed); err == nil {
			data.Body = parsed
		} else {
			data.Body = buf.String()
		}
	}

	if d := cleanup(r.URL.Query()); len(d) > 0 {
		data.Query = d
	}

	json.NewEncoder(w).Encode(data)
}

func cleanup(m interface{}) map[string]interface{} {
	h := make(map[string][]string)

	switch a := m.(type) {
	case http.Header:
		h = a
	case url.Values:
		h = a
	default:
		return nil
	}

	local := make(map[string]interface{})

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

		var (
			value interface{}
			key   = strings.ToLower(k)
		)

		if len(v) == 1 {
			value = v[0]
		} else {
			value = v
		}

		if key == "via" && fmt.Sprintf("%v", value) == "1.1 google" {
			continue
		}

		local[key] = value
	}

	return local
}
