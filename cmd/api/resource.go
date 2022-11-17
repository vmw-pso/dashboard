package main

import (
	"errors"
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
			Active         bool     `json:"active"`
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
			Active:         input.Active,
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

func (app *application) handleShowResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := app.readIDParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		resource, err := app.models.Resources.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"resource": resource}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
