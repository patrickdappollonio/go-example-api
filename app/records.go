package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
)

var reIsValidParam = regexp.MustCompile(`^[a-v0-9]+$`)

func save(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorJSON(w, "Unable to read request body in buffer: %s", err.Error())
		return
	}

	var m interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorJSON(w, "Unable to parse JSON body from request. Body doesn't seem to be JSON: %s", err.Error())
		return
	}

	key, err := saveContent(r, buf.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorJSON(w, "Unable to save content to App Engine datastore: %s", err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/get/%v", key))
	w.WriteHeader(http.StatusCreated)
}

func get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !reIsValidParam.MatchString(id) {
		w.WriteHeader(http.StatusNotFound)
		errorJSON(w, "Record with ID %q not found", id)
		return
	}

	data, err := getContent(r, id)
	if err != nil {
		if err == notfound {
			w.WriteHeader(http.StatusNotFound)
			errorJSON(w, "Record with ID %q not found", id)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		errorJSON(w, "Unable to process record with ID %q: %s", id, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, data)
}
