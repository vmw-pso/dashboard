package main

import (
	"net/http"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/data"
	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

func (app *application) handleCreateResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ID             int64    `json:"id"`
			FirstName      string   `json:"firstName"`
			LastName       string   `json:"lastName"`
			Position       string   `json:"position"`
			Clearance      string   `json:"clearance"`
			Specialties    []string `json:"specialties"`
			Certifications []string `json:"certifications"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		resource := data.Resource{
			ID:             input.ID,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			Position:       input.Position,
			Clearance:      input.Clearance,
			Specialties:    input.Specialties,
			Certifications: input.Certifications,
		}

		v := validator.New()

		if data.ValidateResource(v, &resource); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Resources.Insert(&resource)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusCreated, envelope{"resource": resource}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
