package postgres

import (
	"context"
	"database/sql"
	"errors"
	"examFortune/pkg/models"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"math/rand"
	"time"
)

const (
	insertSql = "INSERT INTO student (username, password, group_name, subject_name, life_time, is_last)" +
		" VALUES ($1,$2,$3,$4,$5,$6)"
	getNameSyllabus   = "select syllabus_id,name,syllabus_info_id from syllabus where teacher_id=$1"
	getTeacherId      = "SELECT teacher_id from teacher where authorization_id=$1"
	getRoleByUsername = "SELECT authorization_id, role FROM auth WHERE username=$1"
	auth              = "SELECT authorization_id, password FROM auth WHERE username = $1"

	deleteTopicWithPlan                 = "delete from topic where plan_id=(select plan_id from session_plan where syllabus_info_id=$1)"
	deleteSessionPlanWithSyllabusInfoId = "DELETE FROM session_plan WHERE syllabus_info_id=$1"
	deleteIndependentStudyTopic         = "delete from independent_study_topic where independent_study_plan_id=(select independent_study_plan_id from independent_study_plan where syllabus_info_id=$1)"
	deleteIndependentStudyPlan          = "DELETE FROM independent_study_plan WHERE syllabus_info_id=$1"
	deleteStudentSyllabus               = "delete from student_syllabus where syllabus_id=(select syllabus_id from syllabus where syllabus_info_id=$1)"
	deleteSyllabusTableRow              = "DELETE FROM  syllabus WHERE syllabus_info_id=$1"
	deleteSyllabusInfo                  = "DELETE FROM syllabus_info WHERE syllabus_info_id=$1"

	selectTopicWithPlan = "select lecture,lecture_hours,practice,practice_hours,assignment,week_number " +
		"from topic where plan_id=(select plan_id from session_plan where syllabus_info_id=$1)"
	selectIndependentStudyTopic = "select week_numbers,topics,hours,recommended_literature,sudmission_form " +
		"from independent_study_topic where independent_study_plan_id=(select independent_study_plan_id from independent_study_plan where syllabus_info_id=$1)"
	selectSyllabusTableRow = "select name FROM  syllabus WHERE syllabus_info_id=$1"
	selectSyllabusInfo     = "select credits_num,goals,skills_competences,objectives,learning_outcomes,prerequisites,postrequisites,instructors " +
		"FROM syllabus_info WHERE syllabus_info_id=$1"
)

type PgModel struct {
	Pool *pgxpool.Pool
}

func (m *PgModel) InsertSyllabus(username, password, group_name, subject_name string) (int, error) {

	var id uint64
	row := m.Pool.QueryRow(context.Background(), insertSql,
		username, password, group_name, subject_name, 3, false)
	err := row.Scan(id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

var authID int
var teacherIDWithSyllabus int

func (m *PgModel) GetNameSyllabus() ([]*models.Syllabus, error) {
	var students []*models.Syllabus
	rows, err := m.Pool.Query(context.Background(), getNameSyllabus, teacherIDWithSyllabus)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := &models.Syllabus{}
		err = rows.Scan(&s.ID, &s.Title, &s.SyllabusInfoID)
		if err != nil {
			return nil, err
		}

		students = append(students, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (m *PgModel) GetTeacherId() (int, error) {
	var id int
	err := m.Pool.QueryRow(context.Background(), getTeacherId, authID).
		Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Teacher id", id)
	teacherIDWithSyllabus = id
	return id, nil
}

func (m *PgModel) DeleteStudentById(id int) error {
	_, err := m.Pool.Exec(context.Background(), deleteStudentSyllabus, id)
	_, err = m.Pool.Exec(context.Background(), deleteTopicWithPlan, id)
	_, err = m.Pool.Exec(context.Background(), deleteSessionPlanWithSyllabusInfoId, id)
	_, err = m.Pool.Exec(context.Background(), deleteIndependentStudyTopic, id)
	_, err = m.Pool.Exec(context.Background(), deleteIndependentStudyPlan, id)
	_, err = m.Pool.Exec(context.Background(), deleteSyllabusTableRow, id)
	_, err = m.Pool.Exec(context.Background(), deleteSyllabusInfo, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *PgModel) Authenticate(username, password string) (int, error) {
	var id int
	var pass string
	row := m.Pool.QueryRow(context.Background(), auth, username)
	err := row.Scan(&id, &pass)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	if pass != password {
		return 0, models.ErrInvalidCredentials
	}
	authID = id
	return id, nil
}

//func (m *PgModel) Get(id int) (*models.Student, error) {
//	return nil, nil
//}
//

func (m *PgModel) GetSyllabusById(id int) (*models.TopicWeek, *models.StudentTopicWeek, *models.Syllabus, error) {
	topic := &models.TopicWeek{}
	independent := &models.StudentTopicWeek{}
	syllabus := &models.Syllabus{}

	err := m.Pool.QueryRow(context.Background(), selectTopicWithPlan, id).
		Scan(&topic.LectureTopic, &topic.LectureHours, &topic.PracticeTopic, &topic.PracticeHours, &topic.Assignment, &topic.WeekNumber)

	err = m.Pool.QueryRow(context.Background(), selectIndependentStudyTopic, id).
		Scan(&independent.WeekNumber, &independent.Topics, &independent.Hours, &independent.RecommendedLiterature, &independent.SubmissionForm)

	err = m.Pool.QueryRow(context.Background(), selectSyllabusTableRow, id).
		Scan(&syllabus.Title)

	err = m.Pool.QueryRow(context.Background(), selectSyllabusInfo, id).
		Scan(&syllabus.Credits, &syllabus.Goals, &syllabus.SkillsCompetences, &syllabus.Objectives, &syllabus.LearningOutcomes, &syllabus.Prerequisites, &syllabus.Postrequisites, &syllabus.Instructors)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil, nil, models.ErrNoRecord
		} else {
			return nil, nil, nil, err
		}
	}
	return topic, independent, syllabus, nil
}

//

func (m *PgModel) GetRoleByUsername(username string) (*models.User, error) {
	s := &models.User{}
	err := m.Pool.QueryRow(context.Background(), getRoleByUsername, username).
		Scan(&s.ID, &s.Role)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

//
///*func (m *PgModel) DeleteStudentByUsername(username string) error {
//	_, err := m.Pool.Exec(context.Background(), deleteStudentByUsername, username)
//	if err != nil {
//		return err
//	}
//	return nil
//}*/
//
//func (m *PgModel) UpdateStudent(s *models.Student) error {
//	_, err := m.Pool.Exec(context.Background(), updateStudent, s.Username, s.Password, s.GroupName, s.SubjectName, s.LifeTime, s.IsLast, s.ID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
func init() {

	rand.Seed(time.Now().UnixNano())
}
