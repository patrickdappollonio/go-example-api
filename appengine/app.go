package main

import (
	"net/http"

	"github.com/patrickdappollonio/go-example-api/app"
)

func init() {
	http.Handle("/", app.Router())
}
