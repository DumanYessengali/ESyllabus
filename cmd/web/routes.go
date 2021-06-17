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

	mux.Get("/admin", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.getMainPageTeacherConfirmed))
	mux.Get("/inProcess", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.getMainPageTeacherInProcess))
	mux.Get("/approvement", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.getMainPageTeacherApprovement))
	mux.Get("/admin/create", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.createSyllabusGet))
	mux.Post("/admin/create", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.createSyllabus))
	mux.Get("/admin/delete", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.deleteStudent))
	mux.Get("/admin/send", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.sendSyllabus))
	mux.Get("/admin/updateSyllabus", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.updateSyllabus))
	mux.Post("/admin/updateSyllabuss", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.updateSyllabuss))
	mux.Get("/admin/updateTopicOpen", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.updateTopicOpen))
	mux.Post("/admin/updateTopic", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.updateTopic))
	mux.Get("/admin/updateIndepOpen", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.updateIndepTopicOpen))
	mux.Post("/admin/updateIndep", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.updateIndepTopic))
	mux.Get("/admin/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication, app.requireTeacher).ThenFunc(app.getSyllabusById))
	mux.Get("/createPDF", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getCreatePDF2))

	mux.Get("/coordinator/confirm/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication, app.requireCoordinator).ThenFunc(app.confirmSyllabus))
	mux.Post("/coordinator/reject/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication, app.requireCoordinator).ThenFunc(app.rejectSyllabus))
	mux.Get("/coordinator", dynamicMiddleware.Append(app.requireAuthentication, app.requireCoordinator).ThenFunc(app.getMainPageCoordinator))
	mux.Get("/coordinator/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication, app.requireCoordinator).ThenFunc(app.getSyllabusByIdForCoordinator))
	mux.Get("/coordinator/feedback", dynamicMiddleware.Append(app.requireAuthentication, app.requireCoordinator).ThenFunc(app.showFeedback))

	mux.Get("/dean", dynamicMiddleware.Append(app.requireAuthentication, app.requireDean).ThenFunc(app.getDeanFeedback))
	mux.Get("/dean/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication, app.requireDean).ThenFunc(app.getSyllabusByIdForDean))
	mux.Post("/dean/reject/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication, app.requireDean).ThenFunc(app.rejectSyllabusDean))
	mux.Get("/dean/ready/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication, app.requireDean).ThenFunc(app.readySyllabus))

	mux.Get("/", dynamicMiddleware.Append(app.requireAuthentication, app.requireStudent).ThenFunc(app.home))
	mux.Get("/syllabusinfo", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getSyllabusByIdForStudents))
	mux.Get("/signin", dynamicMiddleware.ThenFunc(app.signInForm))
	mux.Post("/signin", dynamicMiddleware.ThenFunc(app.signIn))
	mux.Post("/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
