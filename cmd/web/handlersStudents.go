package main

import (
	"errors"
	"examFortune/pkg/forms"
	"examFortune/pkg/models"
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	app.student.GetStudentId()
	syllabus, err := app.student.GetNameSyllabusWithStudent()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.session.PopString(r, "flash")
	app.render(w, r, "home.page.tmpl", &templateData{
		Flash:    flash,
		Syllabus: syllabus,
	})
}

//loginUserForm
func (app *application) signInForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

//loginUser
func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	username := form.Get("username")
	password := form.Get("password")
	id, err := app.student.Authenticate(username, password)
	if err != nil {
		//if errors.Is(err, models.ErrInvalidCredentials) {
		form.Errors.Add("generic", "Username or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		//} else {
		//	app.serverError(w, err)
		//}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	role, err := app.student.GetRoleByUsername(username)
	if err != nil {
		fmt.Print(err.Error())
	}
	if role.Role == "teacher" {
		app.session.Put(r, "sessionTeacherID", id)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	} else if role.Role == "student" {
		app.session.Put(r, "sessionStudentsID", id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if role.Role == "coordinator" {
		app.session.Put(r, "sessionCoordinatorID", id)
		http.Redirect(w, r, "/coordinator", http.StatusSeeOther)
		return
	} else if role.Role == "dean" {
		app.session.Put(r, "sessionDeanID", id)
		http.Redirect(w, r, "/dean", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	if app.IsTeacher(r) {
		app.session.Remove(r, "sessionTeacherID")
	} else if app.IsStudent(r) {
		app.session.Remove(r, "sessionStudentsID")
	} else if app.IsCoordinator(r) {
		app.session.Remove(r, "sessionCoordinatorID")
	} else if app.IsDean(r) {
		app.session.Remove(r, "sessionDeanID")
	}
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
