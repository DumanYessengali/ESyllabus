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

//signupUserForm
func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

//signupUser
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	username := form.Get("username")
	password := form.Get("password")
	password2 := form.Get("password2")
	fullname := form.Get("fullname")
	degree := form.Get("degree")
	rank := form.Get("rank")
	position := form.Get("position")
	contacts := form.Get("contacts")
	interests := form.Get("interests")
	form.Required(
		"username",
		"password",
		"password2",
		"fullname",
		"degree",
		"rank",
		"position",
		"contacts",
		"interests",
	)

	fmt.Println("password: ", password)
	fmt.Println("password2: ", password2)
	if password != password2 {
		form.Errors.Add("generic", "Password should be same")
		newI++

		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	user := &models.User{
		ID:       0,
		Username: username,
		Password: password,
		Role:     "teacher",
	}
	teacherInfo := &models.TeacherInfo{
		FullName:  fullname,
		Degree:    degree,
		Rank:      rank,
		Position:  position,
		Contacts:  contacts,
		Interests: interests,
	}

	id, err := app.student.InsertTeacherAuthTable(user)
	fmt.Println(id)
	_, err = app.student.InsertTeacherTable(teacherInfo, id)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}

	app.session.Put(r, "flash", "Teacher successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/signin"), http.StatusSeeOther)
}

//loginUserForm
func (app *application) signInForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

var wait string

//loginUser
var newI = 0

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
		newI++

		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		//} else {
		//	app.serverError(w, err)
		//}
		fmt.Println(newI)
		if newI >= 5 {
			http.Redirect(w, r, "/wait", http.StatusSeeOther)
			newI = 0
		}
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
	return
}

func (app *application) waitPage(w http.ResponseWriter, r *http.Request) {

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "wait.page.tmpl", &templateData{
		Flash: flash,
		Time:  true,
	})
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
