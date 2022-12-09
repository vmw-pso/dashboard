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
			Sex            string   `json:"sex"`
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
			Sex:            input.Sex,
		}

		v := validator.New()

		if data.ValidateResource(v, resource); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Resources.Insert(&resource)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// app.logger.PrintInfo("resource created", map[string]string{"id": fmt.Sprintf("%d", resource.ID)})

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

func (app *application) handleUpdateResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := app.readIDParam(r)
		if err != nil || id < 1 {
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

		var input struct {
			FirstName      *string  `json:"firstName"`
			LastName       *string  `json:"lastName"`
			Position       *string  `json:"position"`
			Clearance      *string  `json:"clearance"`
			Specialties    []string `json:"specialties"`
			Certifications []string `json:"certifications"`
			Active         *bool    `json:"active"`
			Sex            *string  `json:"sex"`
		}

		err = app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		if input.FirstName != nil {
			resource.FirstName = *input.FirstName
		}

		if input.LastName != nil {
			resource.LastName = *input.LastName
		}

		if input.Position != nil {
			resource.Position = *input.Position
		}

		if input.Clearance != nil {
			resource.Clearance = *input.Clearance
		}

		if input.Specialties != nil {
			resource.Specialties = input.Specialties
		}

		if input.Certifications != nil {
			resource.Certifications = input.Certifications
		}

		if input.Active != nil {
			resource.Active = *input.Active
		}

		if input.Sex != nil {
			resource.Sex = *input.Sex
		}

		v := validator.New()

		if data.ValidateResource(v, *resource); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Resources.Update(resource)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrEditConflict):
				app.editConflictResponse(w, r)
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

func (app *application) handleDeleteResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := app.readIDParam(r)
		if err != nil || id < 1 {
			app.notFoundResponse(w, r)
			return
		}

		err = app.models.Resources.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"message": "successfully deleted"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func (app *application) handleListResources() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Clearance      string
			Specialties    []string
			Certifications []string
			Active         bool
			data.Filters
		}

		v := validator.New()

		qs := r.URL.Query()

		input.Clearance = app.readString(qs, "clearance", "*")
		input.Specialties = app.readCSV(qs, "specialties", []string{})
		input.Certifications = app.readCSV(qs, "certifications", []string{})
		input.Active = app.readBool(qs, "active", true, v)
		input.Filters.Page = app.readInt(qs, "page", 1, v)
		input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
		input.Filters.Sort = app.readString(qs, "sort", "id")
		input.Filters.SortSafelist = []string{"id", "first_name", "last_name", "-id", "-first_name", "-last_name"}

		if data.ValidateFilters(v, input.Filters); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		// currently not filtering for clearance also. Need to fix
		resources, metadata, err := app.models.Resources.GetAll(input.Specialties, input.Certifications, input.Active, input.Filters)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"resources": resources, "metadata": metadata}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
