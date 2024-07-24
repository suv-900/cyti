package api

import (
	"errors"
	"net/http"
	"strings"

	"user_service/internal/data"
)

// adds user with the request
func (app *application) authenticator(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("Authorization")

		if len(authorizationHeader) == 0 {
			app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			//invalid request header send 401 response
			app.invalidAuthenticationHeader(w, r)
			return
		}

		//expired malformed token or empty string
		userid, err := app.verifyToken(headerParts[1])
		if err != nil {
			app.invalidToken(w, r)
			return
		}

		//not accounted errors yet
		user, err := app.models.Users.GetUser(userid)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidTokenDeletedUser(w, r)
			default:
				app.internalServerError(w, r)
			}
		}

		app.contextSetUser(r, &user)

		next.ServeHTTP(w, r)
	}

}

// ordinary function , HandlerFunc(w,r)->(ServeHTTP) , Handler(ServeHTTP)
func (app *application) requireAuthentication(next *http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := app.contextGetUser(r)

		if user.IsAnonymousUser() {
			app.authenticationRequired(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
}
