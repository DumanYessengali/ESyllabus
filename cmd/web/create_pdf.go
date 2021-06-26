package main

import (
	"examFortune/pkg/models"
	"fmt"
	"github.com/rongfengliang/maroto/pkg/consts"
	"github.com/rongfengliang/maroto/pkg/pdf"
	"github.com/rongfengliang/maroto/pkg/props"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (app *application) getCreatePDF2(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}
	app.student.GetTeacherId()

	fmt.Println(id)
	//_, independent, syllabus, teacher, _, err := app.student.GetSyllabusById(id)
	topic, independent, syllabus, _, assessment, err := app.student.GetSyllabusById(id)

	teacher, err1 := app.student.GetFullInfoByTeacherId()

	if err1 != nil {

		fmt.Println(err1)
	}
	fmt.Println(teacher.FullName)
	teacherTable := teacher.FullName + ", " + teacher.Degree + ", " + teacher.Rank + ", " +
		teacher.Position + ", " + teacher.Contacts + ", " + teacher.Interests
	syllabusTable := [][]string{{"Syllabus title", syllabus[0].Title}, {"Discipline", syllabus[0].Discipline},
		{"Number of credits", strconv.Itoa(syllabus[0].Credits)}, {"Prerequisites", syllabus[0].Prerequisites},
		{"Postrequisites", syllabus[0].Postrequisites}, {"Lecturer(s)", teacherTable}}
	syllabusTable2 := [][]string{{"Course goal(s)", syllabus[0].Goals}, {"Course objectives", syllabus[0].Objectives},
		{"Skills & competences", syllabus[0].SkillsCompetences}, {"Course learning outcomes", syllabus[0].LearningOutcomes},
		{"Course instructor(s)", syllabus[0].Instructors}}
	topicTable := [][]string{{"Week Number", "Course Topic", "Lectures (H\\W)", "Practice session (H\\W)", "Lab. sessions (H\\W)", "SIS (H\\W)"}}
	independentTable := [][]string{{"Week Number", "Assignments (topics) for Independent study", "Hours", "Recommended literature and other sources (links)", "Submission Form"}}
	word := []string{"1st attestation", "2nd attestation", "Final exam", "Total"}
	fmt.Println(assessment.Assignment1)
	assessmentTable := [][]string{{"Period", "Assignments", "Number of points", "Total"},
		{word[0], strings.Join(assessment.Assignment1, "\t"), strings.Join(assessment.PointsNum1, " "), "100"},
		{word[1], strings.Join(assessment.Assignment2, "\t"), strings.Join(assessment.PointsNum2, " "), "100"},
		{word[2], word[2], "", "100"},
		{word[3], "0,3 * 1st Att + 0,3 * 2nd Att + 0,4*final", "", "100"},
	}
	fmt.Println("Hello")
	grade := [][]string{{"Grade", "Criteria to be satisfied"}, {"A", "Performs accurate calculations. Uses adequate mathematical operations without errors. Draws logical conclusions, supported by a graph. Provides detailed and correct explanations for the calculations performed."},
		{"B", "Performs well calculations. Uses adequate mathematical operations with few errors. Draws logical conclusions, supported by a graph. Explains the calculations done well."},
		{"C", "I tried to make calculations, but many of them are not accurate. Uses inappropriate mathematical operations, but no errors. Draws conclusions that are not supported by a graph. Provides a small explanation for the calculations performed."},
		{"D", "Does inaccurate calculations. Uses inappropriate mathematical operations. Doesn't draw any conclusions on the schedule. Does not offer an explanation for the calculations performed."},
		{"F", "No response. The student did not try to complete the assignment."},
	}
	gradeTable := [][]string{{"Letter Grade", "Numerical equivalent", "Percentage", "Grade according to the traditional system"},
		{"A", "4,0", "95-100", "Excellent"},
		{"A-", "3,67", "90-94", "Excellent"},
		{"B+", "3,33", "85-89", "Good"},
		{"B", "3,0", "80-84", "Good"},
		{"B-", "2,67", "75-79", "Good"},
		{"С+", "2,33", "70-74", "Satisfactory"},
		{"C", "2,0", "65-69", "Satisfactory"},
		{"C-", "1,67", "60-64", "Satisfactory"},
		{"D+", "1,33", "55-59", "Satisfactory"},
		{"D", "1,0", "50-54", "Satisfactory"},
		{"F", "0", "0-49", "Fail"},
	}
	m := pdf.NewMaroto(consts.Portrait, consts.Letter)
	//m.SetBorder(true)
	topicTableWeek := topicTableFor(topicTable, topic)
	independentTableWeek := topicTableFor1(independentTable, independent)
	m.Row(40, func() {
		m.Col(4, func() {
			_ = m.FileImage("ui/static/img/img.png", props.Rect{
				Center:  true,
				Percent: 60,
			})
		})
		m.Col(4, func() {
			m.Text("Syllabus", props.Text{
				Top:         12,
				Size:        20,
				Extrapolate: true,
			})
			m.Text("Academic Year 2020-2021", props.Text{
				Size: 12,
				Top:  22,
			})
		})
		m.ColSpace(4)
	})

	m.SetBorder(true)
	headersForTable(m, "1. General information")
	tableRow1(m, syllabusTable, 10)
	m.SetBorder(false)
	headersForTable(m, "2. Goals, objectives and learning outcomes of the course")
	m.SetBorder(true)

	tableRow1(m, syllabusTable2, 15)

	m.SetBorder(false)

	headersForTable(m, "3.1 Lecture, practical/seminar/laboratory session plans")
	m.SetBorder(true)
	tableRow2(m, topicTableWeek, 10)
	e := os.Remove("./ui/pdf/Syllabus.pdf")
	if e != nil {
		log.Fatal(e)
	}
	m.SetBorder(false)
	headersForTable(m, "3.2 List of assignments for Student Independent Study")

	m.SetBorder(true)
	tableRow3(m, independentTableWeek, 10)
	m.SetBorder(false)
	headersForTable(m, "4. Student performance evaluation system for the course")
	m.SetBorder(true)
	tableRow4(m, assessmentTable, 20)
	m.SetBorder(false)
	headersForTable(m, "Based on the specific grade for each assignment, and the final grade, following criteria must be satisfied:")
	m.SetBorder(true)
	tableRow1(m, grade, 15)
	m.SetBorder(false)
	headersForTable(m, "Achievement level as per course curriculum shall be assessed according to the evaluation chart adopted by the academic credit system")
	m.SetBorder(true)
	tableRow4(m, gradeTable, 10)
	err = m.OutputFileAndClose("./ui/pdf/Syllabus.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}
	flash := app.session.PopString(r, "flash")

	app.render(w, r, "pdf.page.tmpl", &templateData{
		Flash: flash,
	})

}
func headersForTable(m pdf.Maroto, str string) {
	m.Row(18, func() {
		m.Col(12, func() {
			m.Text(str, props.Text{
				Size:  12,
				Top:   6,
				Align: consts.Center,
			})
		})
	})
}

func tableRow1(m pdf.Maroto, arr [][]string, h float64) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < 1; j++ {
			m.Row(h, func() {
				m.Col(2, func() {
					m.Text(arr[i][j], props.Text{
						Size:  10,
						Align: consts.Center,
					})
				})
				m.Col(10, func() {
					m.Text(arr[i][j+1], props.Text{
						Size: 10,
					})
				})
			})
		}
	}
}

func tableRow2(m pdf.Maroto, arr [][]string, h float64) {
	maxl := 0
	for i := 0; i < len(arr); i++ {
		if len(arr[i][1]) > maxl {
			maxl = len(arr[i][1])
		}
	}
	if maxl/60 >= 1 {
		h = h * math.Floor(float64(maxl/60)+1)
	} else if maxl/30 >= 1 {
		h = h * (math.Floor(float64(maxl/30)) + 1)
	}
	for i := 0; i < len(arr); i++ {
		m.Row(h, func() {
			m.Col(1, func() {
				m.Text(arr[i][0], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				m.Text(arr[i][1], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(1, func() {
				m.Text(arr[i][2], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				m.Text(arr[i][3], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(2, func() {
				m.Text(arr[i][4], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(2, func() {
				m.Text(arr[i][5], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
		})
	}
}

func tableRow3(m pdf.Maroto, arr [][]string, h float64) {
	maxl := 0
	for i := 0; i < len(arr); i++ {
		if len(arr[i][1]) > maxl {
			maxl = len(arr[i][1])
		}
	}
	if maxl/80 >= 1 {
		h = h * math.Floor(float64(maxl/80)+1)
	}
	for i := 0; i < len(arr); i++ {
		m.Row(h, func() {
			m.Col(1, func() {
				m.Text(arr[i][0], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(5, func() {
				m.Text(arr[i][1], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(1, func() {
				m.Text(arr[i][2], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				m.Text(arr[i][3], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(2, func() {
				m.Text(arr[i][4], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
		})
	}
}

func tableRow4(m pdf.Maroto, arr [][]string, h float64) {
	for i := 0; i < len(arr); i++ {
		m.Row(h, func() {
			m.Col(3, func() {
				m.Text(arr[i][0], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				m.Text(arr[i][1], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				m.Text(arr[i][2], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				m.Text(arr[i][3], props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
		})
	}
}

func topicTableFor(topicTable [][]string, topic []*models.TopicWeek) [][]string {
	for j := 0; j < 10; j++ {
		arr := []string{strconv.Itoa(topic[j].WeekNumber), string(topic[j].LectureTopic), strconv.Itoa(topic[j].LectureHours), string(topic[j].PracticeTopic), strconv.Itoa(topic[j].PracticeHours), string(topic[j].Assignment)}
		topicTable = append(topicTable, arr)
	}
	return topicTable
}
func topicTableFor1(independentTable [][]string, independent []*models.StudentTopicWeek) [][]string {
	for j := 0; j < 10; j++ {
		arr := []string{strconv.Itoa(independent[j].WeekNumber), string(independent[j].Topics), strconv.Itoa(independent[j].Hours), string(independent[j].RecommendedLiterature), independent[j].SubmissionForm}
		independentTable = append(independentTable, arr)
	}
	return independentTable
}
