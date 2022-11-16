package main

import (
	"net/http"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/data"
	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

func (app *application) handleCreatePosition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title string `json:"title"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		position := &data.Position{
			Title: input.Title,
		}

		v := validator.New()

		if data.ValidatePosition(v, position.Title); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Positions.Insert(*position)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}
