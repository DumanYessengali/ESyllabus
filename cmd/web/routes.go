package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	mux.Get("/admin", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getMainPage))
	mux.Get("/admin/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createStudentFormAdmin))
	mux.Post("/admin/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createStudent))
	mux.Get("/admin/delete", dynamicMiddleware.ThenFunc(app.deleteStudent))

	mux.Get("/", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.home))
	mux.Get("/signin", dynamicMiddleware.ThenFunc(app.signInForm))
	mux.Post("/signin", dynamicMiddleware.ThenFunc(app.signIn))
	mux.Post("/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
