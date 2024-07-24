package api

import (
	"context"
	"net/http"

	"user_service/internal/data"
)

type contextKey string

var userContextKey = contextKey("user")

// updates request context with key and value
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("cannot typecast to *data.User probably unknown value")
	}
	return user
}
