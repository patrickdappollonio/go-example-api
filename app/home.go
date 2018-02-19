package app

import (
	"fmt"
	"net/http"
)

const project = "https://github.com/patrickdappollonio/go-example-api"

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Usage at:", project)
}
