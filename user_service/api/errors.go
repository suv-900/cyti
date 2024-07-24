package api

import "net/http"

func (app *application) sendErrorResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, env, statusCode)
	if err != nil {
		//err while sending
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) invalidAuthenticationHeader(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication headers."
	app.sendErrorResponse(w, http.StatusExpectationFailed, message)
}

func (app *application) invalidToken(w http.ResponseWriter, r *http.Request) {
	message := "token invalid"
	app.sendErrorResponse(w, http.StatusUnauthorized, message)
}

func (app *application) invalidTokenDeletedUser(w http.ResponseWriter, r *http.Request) {
	message := "user is deleted[dead token] please create new account."
	app.sendErrorResponse(w, http.StatusNotFound, message)
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request) {
	message := "server error"
	app.sendErrorResponse(w, http.StatusInternalServerError, message)
}

func (app *application) authenticationRequired(w http.ResponseWriter, r *http.Request) {
	message := "authentication required."
	app.sendErrorResponse(w, http.StatusUnauthorized, message)
}

func (app *application) routeNotFound(w http.ResponseWriter, r *http.Request) {
	message := "route not available."
	app.sendErrorResponse(w, http.StatusNotFound, message)
}

func (app *application) wrongcredentials(w http.ResponseWriter, r *http.Request) {
	message := "wrong credentials."
	app.sendErrorResponse(w, http.StatusUnauthorized, message)
}

func (app *application) userNotFound(w http.ResponseWriter, r *http.Request) {
	message := "user not found"
	app.sendErrorResponse(w, http.StatusNotFound, message)
}

func (app *application) userExists(w http.ResponseWriter, r *http.Request) {
	message := "user exists."
	app.sendErrorResponse(w, http.StatusConflict, message)
}

func (app *application) fileSizeLimitExceded(w http.ResponseWriter, r *http.Request) {
	message := "file should be less than 32mb"
	app.sendErrorResponse(w, http.StatusBadRequest, message)
}
