package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// movies
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:read", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:write", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))

	// users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// authentication
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
