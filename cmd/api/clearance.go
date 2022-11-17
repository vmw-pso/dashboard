package main

import (
	"net/http"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/data"
	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

func (app *application) handleCreateClearance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Description string `json:"description"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		clearance := &data.Clearance{
			Description: input.Description,
		}

		v := validator.New()

		if data.ValidateDescription(v, clearance.Description); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Clearances.Insert(clearance)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusCreated, envelope{"clearance": clearance}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
