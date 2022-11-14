package main

import "net/http"

func (app *application) handleHealthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env := envelope{
			"status": "available",
			"system_info": map[string]string{
				"environment": app.cfg.env,
				"version":     version,
			},
		}

		err := app.writeJSON(w, http.StatusOK, env, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
