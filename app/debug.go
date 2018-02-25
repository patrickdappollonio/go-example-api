package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func debug(w http.ResponseWriter, r *http.Request) {
	data := debugresponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: cleanup(r.Header),
		Time:    time.Now(),
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
