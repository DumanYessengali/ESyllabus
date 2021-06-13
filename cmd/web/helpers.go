package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	if app.IsTeacher(r) {
		td.IsTeacher = app.IsTeacher(r)
	} else if app.IsStudent(r) {
		td.IsStudent = app.IsStudent(r)
	} else if app.IsCoordinator(r) {
		td.IsCoordinator = app.IsCoordinator(r)
	} else if app.IsDean(r) {
		td.IsDean = app.IsDean(r)
	}
	return td
}
func (app *application) IsTeacher(r *http.Request) bool {
	return app.session.Exists(r, "sessionTeacherID")
}

func (app *application) IsStudent(r *http.Request) bool {
	return app.session.Exists(r, "sessionStudentsID")
}

func (app *application) IsCoordinator(r *http.Request) bool {
	return app.session.Exists(r, "sessionCoordinatorID")
}

func (app *application) IsDean(r *http.Request) bool {
	return app.session.Exists(r, "sessionDeanID")
}
