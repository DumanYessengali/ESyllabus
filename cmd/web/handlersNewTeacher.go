package main

import (
	"errors"
	"examFortune/pkg/forms"
	"examFortune/pkg/models"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) NewTeacherHomePage(w http.ResponseWriter, r *http.Request) {
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

	app.render(w, r, "newTeacher.page.tmpl", &templateData{
		Flash:    flash,
		Syllabus: syllabus,
	})
}

func (app *application) createSyllabusGetForNewTeacher(w http.ResponseWriter, r *http.Request) {
	discipline, _ := app.student.SelectAllDiscipline()
	app.render(w, r, "createSyllabusForNewTeacher.page.tmpl", &templateData{
		Form:       forms.New(nil),
		Discipline: discipline,
	})
}

func (app *application) createSyllabusForNewTeacher(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	//form.Required(
	//	"title",
	//	"subjectName",
	//	"credits_num",
	//	"course_goal",
	//	"skills",
	//	"objectives",
	//	"outcomes",
	//	"prerequisites",
	//	"post_requisites",
	//	"instructors",
	//	"week_num",
	//	"lecture",
	//	"lecture_h",
	//	"practice",
	//	"practice_h",
	//	"assignment",
	//	"week_num1",
	//	"lecture1",
	//	"lecture_h1",
	//	"practice1",
	//	"practice_h1",
	//	"assignment1",
	//	"week_num2",
	//	"lecture2",
	//	"lecture_h2",
	//	"practice2",
	//	"practice_h2",
	//	"assignment2",
	//	"week_num3",
	//	"lecture3",
	//	"lecture_h3",
	//	"practice3",
	//	"practice_h3",
	//	"assignment3",
	//	"week_num4",
	//	"lecture4",
	//	"lecture_h4",
	//	"practice4",
	//	"practice_h4",
	//	"assignment4",
	//	"week_num5",
	//	"lecture5",
	//	"lecture_h5",
	//	"practice5",
	//	"practice_h5",
	//	"assignment5",
	//	"week_num6",
	//	"lecture6",
	//	"lecture_h6",
	//	"practice6",
	//	"practice_h6",
	//	"assignment6",
	//	"week_num7",
	//	"lecture7",
	//	"lecture_h7",
	//	"practice7",
	//	"practice_h7",
	//	"assignment7",
	//	"week_num8",
	//	"lecture8",
	//	"lecture_h8",
	//	"practice8",
	//	"practice_h8",
	//	"assignment8",
	//	"week_num9",
	//	"lecture9",
	//	"lecture_h9",
	//	"practice9",
	//	"practice_h9",
	//	"assignment9",
	//	"week_nums1",
	//	"table2_topic1",
	//	"hours1",
	//	"literature1",
	//	"submission1",
	//	"week_nums2",
	//	"table2_topic2",
	//	"hours2",
	//	"literature2",
	//	"submission2",
	//	"week_nums3",
	//	"table2_topic3",
	//	"hours3",
	//	"literature3",
	//	"submission3",
	//	"week_nums4",
	//	"table2_topic4",
	//	"hours4",
	//	"literature4",
	//	"submission4",
	//	"week_nums5",
	//	"table2_topic5",
	//	"hours5",
	//	"literature5",
	//	"submission5",
	//	"week_nums6",
	//	"table2_topic6",
	//	"hours6",
	//	"literature6",
	//	"submission6",
	//	"week_nums7",
	//	"table2_topic7",
	//	"hours7",
	//	"literature7",
	//	"submission7",
	//	"week_nums8",
	//	"table2_topic8",
	//	"hours8",
	//	"literature8",
	//	"submission8",
	//	"week_nums9",
	//	"table2_topic9",
	//	"hours9",
	//	"literature9",
	//	"submission9",
	//	"week_nums10",
	//	"table2_topic10",
	//	"hours10",
	//	"literature10",
	//	"submission10",
	//)
	//form.MaxLength("username", 100)
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	assesment, _ := strconv.ParseInt(form.Get("assessment"), 10, 64)
	week1_1, _ := strconv.ParseInt(form.Get("week_num"), 10, 64)
	week1_2, _ := strconv.ParseInt(form.Get("week_num1"), 10, 64)
	week1_3, _ := strconv.ParseInt(form.Get("week_num2"), 10, 64)
	week1_4, _ := strconv.ParseInt(form.Get("week_num3"), 10, 64)
	week1_5, _ := strconv.ParseInt(form.Get("week_num4"), 10, 64)
	week1_6, _ := strconv.ParseInt(form.Get("week_num5"), 10, 64)
	week1_7, _ := strconv.ParseInt(form.Get("week_num6"), 10, 64)
	week1_8, _ := strconv.ParseInt(form.Get("week_num7"), 10, 64)
	week1_9, _ := strconv.ParseInt(form.Get("week_num8"), 10, 64)
	week1_10, _ := strconv.ParseInt(form.Get("week_num9"), 10, 64)

	lec1_1, _ := strconv.ParseInt(form.Get("lecture_h"), 10, 64)
	lec1_2, _ := strconv.ParseInt(form.Get("lecture_h1"), 10, 64)
	lec1_3, _ := strconv.ParseInt(form.Get("lecture_h2"), 10, 64)
	lec1_4, _ := strconv.ParseInt(form.Get("lecture_h3"), 10, 64)
	lec1_5, _ := strconv.ParseInt(form.Get("lecture_h4"), 10, 64)
	lec1_6, _ := strconv.ParseInt(form.Get("lecture_h5"), 10, 64)
	lec1_7, _ := strconv.ParseInt(form.Get("lecture_h6"), 10, 64)
	lec1_8, _ := strconv.ParseInt(form.Get("lecture_h7"), 10, 64)
	lec1_9, _ := strconv.ParseInt(form.Get("lecture_h8"), 10, 64)
	lec1_10, _ := strconv.ParseInt(form.Get("lecture_h9"), 10, 64)

	prac1_1, _ := strconv.ParseInt(form.Get("practice_h"), 10, 64)
	prac1_2, _ := strconv.ParseInt(form.Get("practice_h1"), 10, 64)
	prac1_3, _ := strconv.ParseInt(form.Get("practice_h2"), 10, 64)
	prac1_4, _ := strconv.ParseInt(form.Get("practice_h3"), 10, 64)
	prac1_5, _ := strconv.ParseInt(form.Get("practice_h4"), 10, 64)
	prac1_6, _ := strconv.ParseInt(form.Get("practice_h5"), 10, 64)
	prac1_7, _ := strconv.ParseInt(form.Get("practice_h6"), 10, 64)
	prac1_8, _ := strconv.ParseInt(form.Get("practice_h7"), 10, 64)
	prac1_9, _ := strconv.ParseInt(form.Get("practice_h8"), 10, 64)
	prac1_10, _ := strconv.ParseInt(form.Get("practice_h9"), 10, 64)

	week2_1, _ := strconv.ParseInt(form.Get("week_nums1"), 10, 64)
	week2_2, _ := strconv.ParseInt(form.Get("week_nums2"), 10, 64)
	week2_3, _ := strconv.ParseInt(form.Get("week_nums3"), 10, 64)
	week2_4, _ := strconv.ParseInt(form.Get("week_nums4"), 10, 64)
	week2_5, _ := strconv.ParseInt(form.Get("week_nums5"), 10, 64)
	week2_6, _ := strconv.ParseInt(form.Get("week_nums6"), 10, 64)
	week2_7, _ := strconv.ParseInt(form.Get("week_nums7"), 10, 64)
	week2_8, _ := strconv.ParseInt(form.Get("week_nums8"), 10, 64)
	week2_9, _ := strconv.ParseInt(form.Get("week_nums9"), 10, 64)
	week2_10, _ := strconv.ParseInt(form.Get("week_nums10"), 10, 64)

	lec2_1, _ := strconv.ParseInt(form.Get("hours1"), 10, 64)
	lec2_2, _ := strconv.ParseInt(form.Get("hours2"), 10, 64)
	lec2_3, _ := strconv.ParseInt(form.Get("hours3"), 10, 64)
	lec2_4, _ := strconv.ParseInt(form.Get("hours4"), 10, 64)
	lec2_5, _ := strconv.ParseInt(form.Get("hours5"), 10, 64)
	lec2_6, _ := strconv.ParseInt(form.Get("hours6"), 10, 64)
	lec2_7, _ := strconv.ParseInt(form.Get("hours7"), 10, 64)
	lec2_8, _ := strconv.ParseInt(form.Get("hours8"), 10, 64)
	lec2_9, _ := strconv.ParseInt(form.Get("hours9"), 10, 64)
	lec2_10, _ := strconv.ParseInt(form.Get("hours10"), 10, 64)
	//form.Get("username")
	var t1 = []*models.TopicWeek{
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_1),
			LectureTopic:  form.Get("lecture"),
			LectureHours:  int(lec1_1),
			PracticeTopic: form.Get("practice"),
			PracticeHours: int(prac1_1),
			Assignment:    form.Get("assignment"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_2),
			LectureTopic:  form.Get("lecture1"),
			LectureHours:  int(lec1_2),
			PracticeTopic: form.Get("practice1"),
			PracticeHours: int(prac1_2),
			Assignment:    form.Get("assignment1"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_3),
			LectureTopic:  form.Get("lecture2"),
			LectureHours:  int(lec1_3),
			PracticeTopic: form.Get("practice2"),
			PracticeHours: int(prac1_3),
			Assignment:    form.Get("assignment2"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_4),
			LectureTopic:  form.Get("lecture3"),
			LectureHours:  int(lec1_4),
			PracticeTopic: form.Get("practice3"),
			PracticeHours: int(prac1_4),
			Assignment:    form.Get("assignment3"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_5),
			LectureTopic:  form.Get("lecture4"),
			LectureHours:  int(lec1_5),
			PracticeTopic: form.Get("practice4"),
			PracticeHours: int(prac1_5),
			Assignment:    form.Get("assignment4"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_6),
			LectureTopic:  form.Get("lecture5"),
			LectureHours:  int(lec1_6),
			PracticeTopic: form.Get("practice5"),
			PracticeHours: int(prac1_6),
			Assignment:    form.Get("assignment5"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_7),
			LectureTopic:  form.Get("lecture6"),
			LectureHours:  int(lec1_7),
			PracticeTopic: form.Get("practice6"),
			PracticeHours: int(prac1_7),
			Assignment:    form.Get("assignment6"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_8),
			LectureTopic:  form.Get("lecture7"),
			LectureHours:  int(lec1_8),
			PracticeTopic: form.Get("practice7"),
			PracticeHours: int(prac1_8),
			Assignment:    form.Get("assignment7"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_9),
			LectureTopic:  form.Get("lecture8"),
			LectureHours:  int(lec1_9),
			PracticeTopic: form.Get("practice8"),
			PracticeHours: int(prac1_9),
			Assignment:    form.Get("assignment8"),
		},
		&models.TopicWeek{
			TopicWeekID:   0,
			WeekNumber:    int(week1_10),
			LectureTopic:  form.Get("lecture9"),
			LectureHours:  int(lec1_10),
			PracticeTopic: form.Get("practice9"),
			PracticeHours: int(prac1_10),
			Assignment:    form.Get("assignment9"),
		},
	}
	var t2 = []*models.StudentTopicWeek{
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_1),
			Topics:                form.Get("table2_topic1"),
			Hours:                 int(lec2_1),
			RecommendedLiterature: form.Get("literature1"),
			SubmissionForm:        form.Get("submission1"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_2),
			Topics:                form.Get("table2_topic2"),
			Hours:                 int(lec2_2),
			RecommendedLiterature: form.Get("literature2"),
			SubmissionForm:        form.Get("submission2"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_3),
			Topics:                form.Get("table2_topic3"),
			Hours:                 int(lec2_3),
			RecommendedLiterature: form.Get("literature3"),
			SubmissionForm:        form.Get("submission3"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_4),
			Topics:                form.Get("table2_topic4"),
			Hours:                 int(lec2_4),
			RecommendedLiterature: form.Get("literature4"),
			SubmissionForm:        form.Get("submission4"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_5),
			Topics:                form.Get("table2_topic5"),
			Hours:                 int(lec2_5),
			RecommendedLiterature: form.Get("literature5"),
			SubmissionForm:        form.Get("submission5"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_6),
			Topics:                form.Get("table2_topic6"),
			Hours:                 int(lec2_6),
			RecommendedLiterature: form.Get("literature6"),
			SubmissionForm:        form.Get("submission6"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_7),
			Topics:                form.Get("table2_topic7"),
			Hours:                 int(lec2_7),
			RecommendedLiterature: form.Get("literature7"),
			SubmissionForm:        form.Get("submission7"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_8),
			Topics:                form.Get("table2_topic8"),
			Hours:                 int(lec2_8),
			RecommendedLiterature: form.Get("literature8"),
			SubmissionForm:        form.Get("submission8"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_9),
			Topics:                form.Get("table2_topic9"),
			Hours:                 int(lec2_9),
			RecommendedLiterature: form.Get("literature9"),
			SubmissionForm:        form.Get("submission9"),
		},
		&models.StudentTopicWeek{
			StudentTopicWeekID:    0,
			WeekNumber:            int(week2_10),
			Topics:                form.Get("table2_topic10"),
			Hours:                 int(lec2_10),
			RecommendedLiterature: form.Get("literature10"),
			SubmissionForm:        form.Get("submission10"),
		},
	}

	syllabus := &models.Syllabus{
		ID:                0,
		Title:             form.Get("title"),
		Goals:             form.Get("course_goal"),
		SkillsCompetences: form.Get("skills"),
		Status:            "ready",
		Feedback:          " ",
		Objectives:        form.Get("objectives"),
		LearningOutcomes:  form.Get("outcomes"),
		Prerequisites:     form.Get("prerequisites"),
		Postrequisites:    form.Get("post_requisites"),
		Instructors:       form.Get("instructors"),
		Assessment:        int(assesment),
		SyllabusInfoID:    0,
		Table1:            t1,
		Table2:            t2,
	}

	//disciplineId, _ := strconv.ParseInt(form.Get("discipline"), 10, 64)
	teacherId, _ := app.student.GetTeacherId()
	syllabusId, sId, _ := app.student.InsertSyllabusForNewTeacher(syllabus, teacherId, form.Get("title"))
	dId, _ := strconv.ParseInt(r.PostFormValue("discipline"), 10, 64)
	disciplineId, _ := app.student.InsertDisciplineForNewTeacher(int(dId), sId)
	fmt.Println(syllabusId, disciplineId)
	fmt.Println(" disciplineId ", dId)
	app.session.Put(r, "flash", "Syllabus successfully created!")

	app.session.Put(r, "flash", "Syllabus successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/new_teacher"), http.StatusSeeOther)
}

func (app *application) getSyllabusByIdFOrNewTeacher(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	topic, independent, syllabus, teacher, assessment, err := app.student.GetSyllabusById(id)
	fmt.Println(syllabus)
	fmt.Println(topic)
	if err != nil {
		app.notFound(w)
		return
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "selectForNewTeacher.page.tmpl", &templateData{
		Flash:          flash,
		Syllabus:       syllabus,
		Topic:          topic,
		Independent:    independent,
		Teacher:        teacher,
		AssessmentType: assessment,
	})
}
