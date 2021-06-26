package main

import (
	"errors"
	"examFortune/pkg/forms"
	"examFortune/pkg/models"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"strconv"
)

func (app *application) getMainPageTeacherConfirmed(w http.ResponseWriter, r *http.Request) {

	app.student.GetTeacherId()

	syllabus, err := app.student.GetNameSyllabus("ready")

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

func (app *application) getMainPageCoordinator(w http.ResponseWriter, r *http.Request) {

	app.student.GetCoordinatorId()

	syllabus, err := app.student.GetNameSyllabusFromCoordinator("approvement")

	for i := 0; i < len(syllabus)-1; i++ {
		for j := i + 1; j < len(syllabus); j++ {
			if syllabus[i].SyllabusInfoID == syllabus[j].SyllabusInfoID {
				syllabus = removeSyllabusSlice(syllabus, i)
			}
		}
	}

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "coordinator.page.tmpl", &templateData{
		Flash:    flash,
		Syllabus: syllabus,
	})
}

func (app *application) showFeedback(w http.ResponseWriter, r *http.Request) {

	app.student.GetCoordinatorId()

	syllabus, err := app.student.GetNameSyllabusFromCoordinator("in_process")

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "coordinatorFeedback.page.tmpl", &templateData{
		Flash:    flash,
		Syllabus: syllabus,
	})
}

func (app *application) getMainPageTeacherInProcess(w http.ResponseWriter, r *http.Request) {

	app.student.GetTeacherId()

	syllabus, err := app.student.GetNameSyllabus("in_process")

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "inProcess.page.tmpl", &templateData{
		Flash:    flash,
		Syllabus: syllabus,
	})
}

func (app *application) getDeanFeedback(w http.ResponseWriter, r *http.Request) {
	syllabus, err := app.student.GetSyllabusForDean("confirmed")

	for i := 0; i < len(syllabus)-1; i++ {
		for j := i + 1; j < len(syllabus); j++ {
			if syllabus[i].SyllabusInfoID == syllabus[j].SyllabusInfoID {
				syllabus = removeSyllabusSlice(syllabus, i)
			}
		}
	}
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "mainDean.page.tmpl", &templateData{
		Flash:    flash,
		Syllabus: syllabus,
	})
}

func (app *application) getMainPageTeacherApprovement(w http.ResponseWriter, r *http.Request) {

	app.student.GetTeacherId()

	syllabus, err := app.student.GetNameSyllabus("approvement")

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "adjustment.page.tmpl", &templateData{
		Flash:    flash,
		Syllabus: syllabus,
	})
}

func (app *application) sendSyllabus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}
	//-----------------TABLE 1
	var table1 []*models.TopicWeek
	for ind := 1; ind <= 10; ind++ {
		topic := app.tempSyl.GetTempOneSessionTopic(id, ind, "session_topics")
		table1 = append(table1, topic)
	}

	//-----------------TABLE 2
	var table2 []*models.StudentTopicWeek
	for ind := 1; ind <= 10; ind++ {
		topic := app.tempSyl.GetTempOneStudentTopic(id, ind, "student_topics")
		table2 = append(table2, topic)
	}

	goals := app.tempSyl.GetTempFields(id, "goals")[0]
	objectives := app.tempSyl.GetTempFields(id, "objectives")[0]
	outcomes := app.tempSyl.GetTempFields(id, "outcomes")[0]

	_ = app.student.UpdateSyllabusInfoTemp(goals.Content, objectives.Content, outcomes.Content, id)

	_, err = app.student.InsertSessionPlan(table1, id)
	_, err = app.student.InsertIndependentStudyPlan(table2, id)

	err = app.student.SendSyllabus(id, "approvement")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}

	http.Redirect(w, r, "/inProcess", http.StatusSeeOther)
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

func (app *application) updateSyllabus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	syllabus, err := app.student.SelectSyllabusTableRow(id)
	if err != nil {
		app.notFound(w)
		return
	}

	goals := app.tempSyl.GetTempFields(id, "goals")
	objectives := app.tempSyl.GetTempFields(id, "objectives")
	outcomes := app.tempSyl.GetTempFields(id, "outcomes")

	if err != nil {
		app.notFound(w)
		return
	}

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}

	flash := app.session.PopString(r, "flash")
	app.render(w, r, "updateSyllabus.page.tmpl", &templateData{
		Flash:          flash,
		Syllabus:       syllabus,
		TempGoals:      goals,
		TempOutcomes:   outcomes,
		TempObjectives: objectives,
	})
}

func (app *application) updateSyllabuss(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	teacherId, _ := app.student.GetTeacherId()
	username, _ := app.student.GetTeacherUsername(teacherId)
	tempGoals := &models.TempFields{
		TeacherId:      teacherId,
		TeacherName:    username,
		SyllabusInfoId: id,
		Content:        form.Get("goals"),
	}
	tempObj := &models.TempFields{
		TeacherId:      teacherId,
		TeacherName:    username,
		SyllabusInfoId: id,
		Content:        form.Get("objectives"),
	}
	tempOut := &models.TempFields{
		TeacherId:      teacherId,
		TeacherName:    username,
		SyllabusInfoId: id,
		Content:        form.Get("outcomes"),
	}

	if form.Get("goals") == "" {
		tempGoals.Content = form.Get("Goals")
	}
	if form.Get("objectives") == "" {
		tempObj.Content = form.Get("Objectives")
	}
	if form.Get("outcomes") == "" {
		tempOut.Content = form.Get("Outcomes")
	}

	app.tempSyl.UpdateTempField(tempGoals, "goals")
	app.tempSyl.UpdateTempField(tempObj, "objectives")
	app.tempSyl.UpdateTempField(tempOut, "outcomes")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}
	url := "/admin/syllabusinfo?id="
	http.Redirect(w, r, url+strconv.Itoa(id), http.StatusSeeOther)
}

func (app *application) updateTopicOpen(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	weekNum, err := strconv.Atoi(r.URL.Query().Get("weekNum"))

	tempTopics := app.tempSyl.GetTempSessionTopicByWeek(id, weekNum, "session_topics")

	var lt []*models.TempFields
	var lh []*models.TempFields
	var pt []*models.TempFields
	var ph []*models.TempFields
	var a []*models.TempFields

	for _, t := range tempTopics {
		if t.LectureTopic != "" {
			lt = append(lt, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				Content:        t.LectureTopic,
			})
		}
		if t.LectureHours != 0 {
			lh = append(lh, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				ContentInt:     t.LectureHours,
			})
		}
		if t.PracticeTopic != "" {
			pt = append(pt, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				Content:        t.PracticeTopic,
			})
		}
		if t.PracticeHours != 0 {
			ph = append(ph, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				ContentInt:     t.PracticeHours,
			})
		}

		if t.Assignment != "" {
			a = append(a, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				Content:        t.Assignment,
			})
		}
	}

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}

	flash := app.session.PopString(r, "flash")
	app.render(w, r, "updateTopicOpen.page.tmpl", &templateData{
		Flash: flash,
		TopicOneRow: &models.SessionWeek{
			SyllabusInfoId: id,
			WeekNumber:     weekNum,
			LectureTopic:   lt,
			LectureHours:   lh,
			PracticeTopic:  pt,
			PracticeHours:  ph,
			Assignment:     a,
		},
	})
}

func (app *application) updateTopic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	weekNumber, _ := strconv.ParseInt(form.Get("WeekNumber"), 10, 64)
	lectureHours, _ := strconv.ParseInt(form.Get("LectureHours"), 10, 64)
	practiceHours, _ := strconv.ParseInt(form.Get("PracticeHours"), 10, 64)
	lt := form.Get("LectureTopic")
	pt := form.Get("PracticeTopic")
	a := form.Get("Assignment")

	//---------------------------TUUUT
	teacherId, _ := app.student.GetTeacherId()
	teacherUsername, _ := app.student.GetTeacherUsername(teacherId)

	if lectureHours == 0 {
		lectureHours, _ = strconv.ParseInt(form.Get("lectureHours"), 10, 64)
	}
	if practiceHours == 0 {
		practiceHours, _ = strconv.ParseInt(form.Get("practiceHours"), 10, 64)
	}
	if pt == "" {
		pt = form.Get("practiceTopic")
	}
	if lt == "" {
		lt = form.Get("lectureTopic")
	}
	if a == "" {
		a = form.Get("assignment")
	}

	topic := &models.TopicWeek{
		WeekNumber:     int(weekNumber),
		SyllabusInfoId: id,
		TeacherId:      teacherId,
		TeacherName:    teacherUsername,
		LectureTopic:   lt,
		LectureHours:   int(lectureHours),
		PracticeTopic:  pt,
		PracticeHours:  int(practiceHours),
		Assignment:     a,
	}
	app.tempSyl.UpdateTempSessionTopic(topic, "session_topics")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}
	url := "/inProcess"
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *application) updateIndepTopicOpen(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	weekNum, err := strconv.Atoi(r.URL.Query().Get("weekNum"))

	tempTopics := app.tempSyl.GetTempStudentTopicByWeek(id, weekNum, "student_topics")

	var st []*models.TempFields
	var h []*models.TempFields
	var l []*models.TempFields
	var sf []*models.TempFields
	for _, t := range tempTopics {
		if t.Topics != "" {
			st = append(st, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				Content:        t.Topics,
				ContentInt:     0,
			})
		}
		if t.Hours != 0 {
			h = append(h, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				ContentInt:     t.Hours,
				Content:        "",
			})
		}
		if t.RecommendedLiterature != "" {
			l = append(l, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				Content:        t.RecommendedLiterature,
				ContentInt:     0,
			})
		}
		if t.SubmissionForm != "" {
			sf = append(sf, &models.TempFields{
				TeacherId:      t.TeacherId,
				TeacherName:    t.TeacherName,
				SyllabusInfoId: id,
				Content:        t.SubmissionForm,
				ContentInt:     0,
			})
		}
	}

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}

	flash := app.session.PopString(r, "flash")
	app.render(w, r, "updateIndepTopicOpen.page.tmpl", &templateData{
		Flash: flash,
		IndepTopicOneRow: &models.StudentWeek{
			SyllabusInfoId:        id,
			WeekNumber:            weekNum,
			Topics:                st,
			Hours:                 h,
			RecommendedLiterature: l,
			SubmissionForm:        sf,
		},
	})
}

func (app *application) updateIndepTopic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	weekNumber, _ := strconv.ParseInt(form.Get("WeekNumber"), 10, 64)
	hours, _ := strconv.ParseInt(form.Get("Hours"), 10, 64)
	st := form.Get("Topics")
	l := form.Get("RecommendedLiterature")
	sf := form.Get("SubmissionForm")

	teacherId, _ := app.student.GetTeacherId()
	teacherUsername, _ := app.student.GetTeacherUsername(teacherId)

	if hours == 0 {
		hours, _ = strconv.ParseInt(form.Get("hours"), 10, 64)
	}
	if st == "" {
		st = form.Get("topics")
	}
	if l == "" {
		l = form.Get("recommendedLiterature")
	}
	if sf == "" {
		sf = form.Get("submissionForm")
	}

	topic := &models.StudentTopicWeek{
		StudentTopicWeekID:    0,
		SyllabusInfoId:        id,
		WeekNumber:            int(weekNumber),
		TeacherId:             teacherId,
		TeacherName:           teacherUsername,
		Topics:                st,
		Hours:                 int(hours),
		RecommendedLiterature: l,
		SubmissionForm:        sf,
	}
	app.tempSyl.UpdateTempStudentTopic(topic, "student_topics")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}
	url := "/inProcess"
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *application) confirmSyllabus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		app.notFound(w)
		return
	}

	err = app.student.SendSyllabus(id, "confirmed")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}

	url := "/coordinator"
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *application) readySyllabus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		app.notFound(w)
		return
	}
	err = app.student.SendSyllabus(id, "ready")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		println(err.Error())
		return
	}
	url := "/dean"
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *application) rejectSyllabus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	err = app.student.DeleteForRejection(id)
	_, err = app.student.InsertFeedback(form.Get("feedback"), id)
	fmt.Println(err)
	err = app.student.SendSyllabus(id, "in_process")
	if err != nil {
		app.notFound(w)
		return
	}
	url := "/coordinator"
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *application) rejectSyllabusDean(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	_, err = app.student.InsertFeedback(form.Get("feedback"), id)

	err = app.student.SendSyllabus(id, "approvement")
	if err != nil {
		app.notFound(w)
		return
	}
	url := "/dean"
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *application) getSyllabusByIdForCoordinator(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	topic, independent, syllabus, teacher, assessment, err := app.student.GetSyllabusById(id)

	if err != nil {
		app.notFound(w)
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "selectCoordinator.page.tmpl", &templateData{
		Form:           forms.New(nil),
		Flash:          flash,
		Syllabus:       syllabus,
		Topic:          topic,
		Independent:    independent,
		Teacher:        teacher,
		AssessmentType: assessment,
	})
}

func (app *application) getSyllabusByIdForDean(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	topic, independent, syllabus, teacher, assessment, err := app.student.GetSyllabusById(id)

	if err != nil {
		app.notFound(w)
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "selectDean.page.tmpl", &templateData{
		Form:           forms.New(nil),
		Flash:          flash,
		Syllabus:       syllabus,
		Topic:          topic,
		Independent:    independent,
		Teacher:        teacher,
		AssessmentType: assessment,
	})
}

func (app *application) getCreatePDF(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}
	fmt.Println(id)
	topic, independent, syllabus, teacher, _, err := app.student.GetSyllabusById(id)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(190, 7, syllabus[0].Goals, "0", 5, "CM", false, 0, "")
	pdf.CellFormat(190, 7, topic[0].LectureTopic, "0", 0, "CM", false, 0, "")
	pdf.CellFormat(190, 7, independent[0].RecommendedLiterature, "0", 0, "CM", false, 0, "")
	pdf.CellFormat(190, 7, teacher[0].FullName, "0", 0, "CM", false, 0, "")
	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)

	err = pdf.OutputFileAndClose(syllabus[0].Title + ".pdf")
	if err != nil {
		panic(err)
	}

	if err != nil {
		app.notFound(w)
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "pdf.page.tmpl", &templateData{
		Flash: flash,
	})
}

func (app *application) getSyllabusById(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	syllabus, teacher, assessment, err := app.student.GetSyllabusTeacherInProcess(id)

	//-----------------SYLLABUS INFO
	goals := app.tempSyl.GetTempFields(id, "goals")
	objectives := app.tempSyl.GetTempFields(id, "objectives")
	outcomes := app.tempSyl.GetTempFields(id, "outcomes")

	for i := 0; i < len(goals)-1; i++ {
		for j := i + 1; j < len(goals); j++ {
			if goals[i].Content == goals[j].Content {
				goals[i].TeacherName = ""
				goals[j].TeacherName = ""
				goals = remove(goals, i)
			}
		}
	}
	for i := 0; i < len(outcomes)-1; i++ {
		for j := i + 1; j < len(outcomes); j++ {
			if outcomes[i].Content == outcomes[j].Content {
				outcomes[i].TeacherName = ""
				outcomes[j].TeacherName = ""
				outcomes = remove(outcomes, i)
			}
		}
	}
	for i := 0; i < len(objectives)-1; i++ {
		for j := i + 1; j < len(objectives); j++ {
			if objectives[i].Content == objectives[j].Content {
				objectives[i].TeacherName = ""
				objectives[j].TeacherName = ""
				objectives = remove(objectives, i)
			}
		}
	}

	//-----------------TABLE 1
	var sessionTopics []*models.SessionWeek
	for ind := 1; ind <= 10; ind++ {
		t1 := app.tempSyl.GetTempSessionTopicByWeek(id, ind, "session_topics")
		var lt []*models.TempFields
		var lh []*models.TempFields
		var pt []*models.TempFields
		var ph []*models.TempFields
		var a []*models.TempFields
		for _, t := range t1 {
			if t.LectureTopic != "" {
				lt = append(lt, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					Content:        t.LectureTopic,
					ContentInt:     0,
				})
			}
			if t.LectureHours != 0 {
				lh = append(lh, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					ContentInt:     t.LectureHours,
					Content:        "",
				})
			}
			if t.PracticeTopic != "" {
				pt = append(pt, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					Content:        t.PracticeTopic,
					ContentInt:     0,
				})
			}
			if t.PracticeHours != 0 {
				ph = append(ph, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					ContentInt:     t.PracticeHours,
					Content:        "",
				})
			}

			if t.Assignment != "" {
				a = append(a, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					Content:        t.Assignment,
					ContentInt:     0,
				})
			}
		}
		sessionTopics = append(sessionTopics, &models.SessionWeek{
			SyllabusInfoId: id,
			WeekNumber:     ind,
			LectureTopic:   lt,
			LectureHours:   lh,
			PracticeTopic:  pt,
			PracticeHours:  ph,
			Assignment:     a,
		})
	}
	for _, v := range sessionTopics {
		for i := 0; i < len(v.LectureTopic)-1; i++ {
			for j := i + 1; j < len(v.LectureTopic); j++ {
				if v.LectureTopic[i].Content == v.LectureTopic[j].Content {
					v.LectureTopic[i].TeacherName = ""
					v.LectureTopic[j].TeacherName = ""
					v.LectureTopic = remove(v.LectureTopic, i)
				}
			}
		}
		for i := 0; i < len(v.PracticeTopic)-1; i++ {
			for j := i + 1; j < len(v.PracticeTopic); j++ {
				if v.PracticeTopic[i].Content == v.PracticeTopic[j].Content {
					v.PracticeTopic[i].TeacherName = ""
					v.PracticeTopic[j].TeacherName = ""
					v.PracticeTopic = remove(v.PracticeTopic, i)
				}
			}
		}
		for i := 0; i < len(v.LectureHours)-1; i++ {
			for j := i + 1; j < len(v.LectureHours); j++ {
				if v.LectureHours[i].ContentInt == v.LectureHours[j].ContentInt {
					v.LectureHours[i].TeacherName = ""
					v.LectureHours[j].TeacherName = ""
					v.LectureHours = remove(v.LectureHours, i)
				}
			}
		}
		for i := 0; i < len(v.PracticeHours)-1; i++ {
			for j := i + 1; j < len(v.PracticeHours); j++ {
				if v.PracticeHours[i].ContentInt == v.PracticeHours[j].ContentInt {
					v.PracticeHours[i].TeacherName = ""
					v.PracticeHours[j].TeacherName = ""
					v.PracticeHours = remove(v.PracticeHours, i)
				}
			}
		}
		for i := 0; i < len(v.Assignment)-1; i++ {
			for j := i + 1; j < len(v.Assignment); j++ {
				if v.Assignment[i].Content == v.Assignment[j].Content {
					v.Assignment[i].TeacherName = ""
					v.Assignment[j].TeacherName = ""
					v.Assignment = remove(v.Assignment, i)
				}
			}
		}
	}

	//-----------------TABLE 2
	var studentTopics []*models.StudentWeek
	for ind := 1; ind <= 10; ind++ {
		t2 := app.tempSyl.GetTempStudentTopicByWeek(id, ind, "student_topics")
		var st []*models.TempFields
		var h []*models.TempFields
		var l []*models.TempFields
		var sf []*models.TempFields
		for _, t := range t2 {
			if t.Topics != "" {
				st = append(st, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					Content:        t.Topics,
					ContentInt:     0,
				})
			}
			if t.Hours != 0 {
				h = append(h, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					ContentInt:     t.Hours,
					Content:        "",
				})
			}
			if t.RecommendedLiterature != "" {
				l = append(l, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					Content:        t.RecommendedLiterature,
					ContentInt:     0,
				})
			}
			if t.SubmissionForm != "" {
				sf = append(sf, &models.TempFields{
					TeacherId:      t.TeacherId,
					TeacherName:    t.TeacherName,
					SyllabusInfoId: id,
					Content:        t.SubmissionForm,
					ContentInt:     0,
				})
			}
		}
		studentTopics = append(studentTopics, &models.StudentWeek{
			SyllabusInfoId:        id,
			WeekNumber:            ind,
			Topics:                st,
			Hours:                 h,
			RecommendedLiterature: l,
			SubmissionForm:        sf,
		})
	}
	for _, v := range studentTopics {
		for i := 0; i < len(v.Topics)-1; i++ {
			for j := i + 1; j < len(v.Topics); j++ {
				if v.Topics[i].Content == v.Topics[j].Content {
					v.Topics[i].TeacherName = ""
					v.Topics[j].TeacherName = ""
					v.Topics = remove(v.Topics, i)
				}
			}
		}
		for i := 0; i < len(v.Hours)-1; i++ {
			for j := i + 1; j < len(v.Hours); j++ {
				if v.Hours[i].ContentInt == v.Hours[j].ContentInt {
					v.Hours[i].TeacherName = ""
					v.Hours[j].TeacherName = ""
					v.Hours = remove(v.Hours, i)
				}
			}
		}
		for i := 0; i < len(v.RecommendedLiterature)-1; i++ {
			for j := i + 1; j < len(v.RecommendedLiterature); j++ {
				if v.RecommendedLiterature[i].Content == v.RecommendedLiterature[j].Content {
					v.RecommendedLiterature[i].TeacherName = ""
					v.RecommendedLiterature[j].TeacherName = ""
					v.RecommendedLiterature = remove(v.RecommendedLiterature, i)
				}
			}
		}
		for i := 0; i < len(v.SubmissionForm)-1; i++ {
			for j := i + 1; j < len(v.SubmissionForm); j++ {
				if v.SubmissionForm[i].Content == v.SubmissionForm[j].Content {
					v.SubmissionForm[i].TeacherName = ""
					v.SubmissionForm[j].TeacherName = ""
					v.SubmissionForm = remove(v.SubmissionForm, i)
				}
			}
		}
	}

	if err != nil {
		app.notFound(w)
		return
	}
	flash := app.session.PopString(r, "flash")

	app.render(w, r, "select.page.tmpl", &templateData{
		Flash:          flash,
		Syllabus:       syllabus,
		TempGoals:      goals,
		TempObjectives: objectives,
		TempOutcomes:   outcomes,
		TempTable1:     sessionTopics,
		TempTable2:     studentTopics,
		Teacher:        teacher,
		AssessmentType: assessment,
	})
}

func (app *application) getSyllabusByIdForStudents(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	topic, independent, syllabus, teacher, assessment, err := app.student.GetSyllabusById(id)

	if err != nil {
		app.notFound(w)
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "selectForStudents.page.tmpl", &templateData{
		Flash:          flash,
		Syllabus:       syllabus,
		Topic:          topic,
		Independent:    independent,
		Teacher:        teacher,
		AssessmentType: assessment,
	})
}

func (app *application) createSyllabusGet(w http.ResponseWriter, r *http.Request) {
	discipline, _ := app.student.SelectAllDiscipline()
	app.render(w, r, "create.page.tmpl", &templateData{
		Form:       forms.New(nil),
		Discipline: discipline,
	})
}

func (app *application) createSyllabus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	assesment, _ := strconv.ParseInt(form.Get("assessment"), 10, 64)

	syllabus := &models.Syllabus{
		ID:                0,
		Title:             form.Get("title"),
		Goals:             "",
		SkillsCompetences: form.Get("skills"),
		Objectives:        "",
		LearningOutcomes:  "",
		Prerequisites:     form.Get("prerequisites"),
		Postrequisites:    form.Get("post_requisites"),
		Instructors:       form.Get("instructors"),
		Assessment:        int(assesment),
		SyllabusInfoID:    0,
		Table1:            nil,
		Table2:            nil,
	}
	dId, _ := strconv.ParseInt(r.PostFormValue("discipline"), 10, 64)
	teacherId, _ := app.student.GetTeacherId()
	teacherIds, _ := app.student.GetTeacherIds(int(dId))
	_, sId, _ := app.student.InsertSyllabus(syllabus, teacherId, form.Get("title"))
	for _, v := range teacherIds {
		if v != teacherId {
			_, _ = app.student.InsertSyllabusForOtherTeachers(v, form.Get("title"), sId)
		}
	}
	_, _ = app.student.InsertDiscipline(int(dId), sId)

	//----------------MONGO INSERT
	teacherUsername, _ := app.student.GetTeacherUsername(teacherId)

	goals := &models.TempFields{
		TeacherId:      teacherId,
		TeacherName:    teacherUsername,
		SyllabusInfoId: sId,
		Content:        form.Get("course_goal"),
	}
	app.tempSyl.InsertTempField(goals, "goals")

	obj := &models.TempFields{
		TeacherId:      teacherId,
		TeacherName:    teacherUsername,
		SyllabusInfoId: sId,
		Content:        form.Get("objectives"),
	}
	app.tempSyl.InsertTempField(obj, "objectives")

	out := &models.TempFields{
		TeacherId:      teacherId,
		TeacherName:    teacherUsername,
		SyllabusInfoId: sId,
		Content:        form.Get("outcomes"),
	}
	app.tempSyl.InsertTempField(out, "outcomes")

	//-----------table 1 and 2 topics
	var t1 []*models.TopicWeek
	var t2 []*models.StudentTopicWeek

	for i := 1; i < 11; i++ {
		lh, _ := strconv.ParseInt(form.Get(fmt.Sprintf("lecture_h%d", i)), 10, 64)
		ph, _ := strconv.ParseInt(form.Get(fmt.Sprintf("practice_h%d", i)), 10, 64)

		t1 = append(t1, &models.TopicWeek{
			TopicWeekID:    0,
			SyllabusInfoId: sId,
			TeacherId:      teacherId,
			TeacherName:    teacherUsername,
			WeekNumber:     i,
			LectureTopic:   form.Get(fmt.Sprintf("lecture%d", i)),
			LectureHours:   int(lh),
			PracticeTopic:  form.Get(fmt.Sprintf("practice%d", i)),
			PracticeHours:  int(ph),
			Assignment:     form.Get(fmt.Sprintf("assignment%d", i)),
		})

		sh, _ := strconv.ParseInt(form.Get(fmt.Sprintf("hours%d", i)), 10, 64)

		t2 = append(t2, &models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			SyllabusInfoId:        sId,
			TeacherId:             teacherId,
			TeacherName:           teacherUsername,
			WeekNumber:            i,
			Topics:                form.Get(fmt.Sprintf("table2_topic%d", i)),
			Hours:                 int(sh),
			RecommendedLiterature: form.Get(fmt.Sprintf("literature%d", i)),
			SubmissionForm:        form.Get(fmt.Sprintf("submission%d", i)),
		})

	}

	for _, t := range t1 {
		app.tempSyl.InsertTempSessionTopic(t, "session_topics")
	}

	//-----------table2: student independent topics
	for _, t := range t2 {
		app.tempSyl.InsertTempStudentTopic(t, "student_topics")
	}
	//------------------------

	app.session.Put(r, "flash", "Syllabus successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/inProcess"), http.StatusSeeOther)
}

func remove(slice []*models.TempFields, s int) []*models.TempFields {
	return append(slice[:s], slice[s+1:]...)
}

func removeSyllabusSlice(slice []*models.Syllabus, s int) []*models.Syllabus {
	return append(slice[:s], slice[s+1:]...)
}
