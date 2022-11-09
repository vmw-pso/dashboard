package main

import "net/http"

func (app *application) handleHealthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("status: available"))
	}
}
