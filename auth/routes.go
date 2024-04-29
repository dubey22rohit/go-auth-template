package main

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	mux.HandleFunc("POST /users/login/", app.loginUserHandler)
	mux.HandleFunc("POST /users/logout/", app.logoutUserHandler)
	mux.HandleFunc("POST /users/register/", app.registerUserHandler)
	mux.HandleFunc("POST /users/current-user/", app.currentUserHandler)
	// mux.HandleFunc("PUT /users/activate/:id/", app.activateUserHandler) : Use Later

	return app.recoverPanic(mux)
}
