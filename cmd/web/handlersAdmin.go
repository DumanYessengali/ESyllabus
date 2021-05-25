package main

import (
	"errors"
	"examFortune/pkg/forms"
	"examFortune/pkg/models"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) getMainPage(w http.ResponseWriter, r *http.Request) {

	students, err := app.student.GetAllStudents()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "admin.page.tmpl", &templateData{
		Flash:    flash,
		Students: students,
	})
}

func (app *application) createStudentFormAdmin(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createStudent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "password", "groupName", "subjectName")
	form.MaxLength("username", 100)

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	app.student.InsertStudent(form.Get("username"), form.Get("password"), form.Get("groupName"), form.Get("subjectName"))

	app.session.Put(r, "flash", "Student successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/admin"), http.StatusSeeOther)
}

func (app *application) deleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.student.DeleteStudentById(id)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}

	app.render(w, r, "afterDelete.page.tmpl", &templateData{})
}
