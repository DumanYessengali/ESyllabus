package main

import (
	"errors"
	"examFortune/pkg/forms"
	"examFortune/pkg/models"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) getMainPageTeacher(w http.ResponseWriter, r *http.Request) {

	app.student.GetTeacherId()
	//fmt.Print(app.syllabus.GetAllSyllabuses(1))
	syllabus, err := app.student.GetNameSyllabus()
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
		Syllabus: syllabus,
	})
}

func (app *application) createStudentFormAdmin(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

//func (app *application) createStudent(w http.ResponseWriter, r *http.Request) {
//	if err := r.ParseForm(); err != nil {
//		app.serverError(w, err)
//		return
//	}
//
//	form := forms.New(r.PostForm)
//	form.Required("username", "password", "groupName", "subjectName")
//	form.MaxLength("username", 100)
//
//	if !form.Valid() {
//		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
//		return
//	}
//
//	app.student.InsertSyllabus(form.Get("username"), form.Get("password"), form.Get("groupName"), form.Get("subjectName"))
//
//	app.session.Put(r, "flash", "Student successfully created!")
//
//	http.Redirect(w, r, fmt.Sprintf("/admin"), http.StatusSeeOther)
//}

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

func (app *application) getSyllabusById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	topic, independent, syllabus, teacher, err := app.student.GetSyllabusById(id)

	if err != nil {
		app.notFound(w)
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "select.page.tmpl", &templateData{
		Flash:       flash,
		Syllabus:    syllabus,
		Topic:       topic,
		Independent: independent,
		Teacher:     teacher,
	})
}

func (app *application) createSyllabus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required(
		"title",
		"fullname",
		"degree",
		"rank",
		"position",
		"contacts",
		"interests",
		"subjectName",
		"credits_num",
		"course_goal",
		"skills",
		"objectives",
		"outcomes",
		"prerequisites",
		"post_requisites",
		"instructors",
		"week_num",
		"lecture",
		"lecture_h",
		"practice",
		"practice_h",
		"assignment",
		"week_num2",
		"table2_topic",
		"hours",
		"literature",
		"submission",
	)
	//form.MaxLength("username", 100)

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
	//form.Get("username")
	syllabus := &models.Syllabus{
		ID:                0,
		Title:             form.Get("title"),
		Teacher:           nil,
		Credits:           0,
		Goals:             "",
		SkillsCompetences: "",
		Objectives:        "",
		LearningOutcomes:  "",
		Prerequisites:     "",
		Postrequisites:    "",
		Instructors:       "",
		SyllabusInfoID:    0,
		Table1:            nil,
		Table2:            nil,
	}

	teacherId, _ := app.student.GetTeacherId()
	syllabusId, _ := app.student.InsertSyllabus(syllabus, teacherId, 1, form.Get("title"))

	fmt.Println(syllabusId)

	app.session.Put(r, "flash", "Syllabus successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/admin"), http.StatusSeeOther)
}
