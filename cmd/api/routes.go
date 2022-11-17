package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.handleHealthcheck())

	mux.HandlerFunc(http.MethodPost, "/v1/positions", app.handleCreatePosition())

	mux.HandlerFunc(http.MethodPost, "/v1/clearances", app.handleCreateClearance())

	mux.HandlerFunc(http.MethodPost, "/v1/resources", app.handleCreateResource())
	mux.HandlerFunc(http.MethodGet, "/v1/resources/:id", app.handleShowResource())

	return app.recoverPanic(mux)
}
