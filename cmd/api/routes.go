package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc(http.MethodGet, "/healthcheck", app.handleHealthcheck())

	return mux
}
