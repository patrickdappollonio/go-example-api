package app

import (
	"fmt"
	"net/http"
)

const project = "https://github.com/patrickdappollonio/go-example-api"

func nocontent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Usage at:", project)
}
