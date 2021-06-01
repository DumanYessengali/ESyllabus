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
	selectTeacherInfo  = "select fullname, degree, rank, position, contacts, interests from teacher where authorization_id = $1"
	insertSyllabus     = "insert into syllabus (teacher_id, syllabus_info_id, discipline_id, name) values ($1, $2, $3, $4) returning syllabus_id"
	insertSyllabusInfo = "insert into syllabus_info(credits_num, goals, skills_competences, objectives, learning_outcomes," +
		"prerequisites, postrequisites, instructors) values($1, $2, $3, $4, $5, $6, $7, $8) returning syllabus_info_id"
	insertSessionPlan  = "insert into session_plan(total_time, syllabus_info_id) values($1, $2) returning plan_id"
	insertSessionTopic = "insert into topic (lecture, lecture_hours, practice, practice_hours, assignment, week_number, plan_id)" +
		"values ($1, $2, $3, $4, $5, $6, $7) returning topic_id"
	insertIndependentStudyPlan      = "insert into independent_study_plan(total_time, syllabus_info_id) values($1, $2) returning independent_study_plan_id"
	insertIndependentStudyPlanTopic = "insert into independent_study_topic (week_numbers, topics, hours, recommended_literature, sudmission_form, independent_study_plan_id)" +
		"values ($1, $2, $3, $4, $5, $6) returning independent_study_topic_id"
	GetStudentId = "SELECT student_id from student where authorization_id=$1"
)

type PgModel struct {
	Pool *pgxpool.Pool
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
////
//
//func (m *PgModel) GetSyllabusById(id int) ([]*models.TopicWeek, []*models.StudentTopicWeek, []*models.Syllabus, []*models.TeacherInfo, error) {
//
//	var topic []*models.TopicWeek
//	var independent []*models.StudentTopicWeek
//	var syllabus []*models.Syllabus
//	var teacher []*models.TeacherInfo
//
//	rows1, err := m.Pool.Query(context.Background(), selectTopicWithPlan, id)
//	rows2, err := m.Pool.Query(context.Background(), selectIndependentStudyTopic, id)
//	rows3, err := m.Pool.Query(context.Background(), selectSyllabusTableRow, id)
//	rows4, err := m.Pool.Query(context.Background(), selectSyllabusInfo, id)
//	rows5, err := m.Pool.Query(context.Background(), selectTeacherInfo, authID)
//	if err != nil {
//		return nil, nil, nil, nil, err
//	}
//
//	for rows3.Next() && rows4.Next() && rows5.Next() {
//
//		s := &models.Syllabus{}
//		te := &models.TeacherInfo{}
//
//		err = rows3.Scan(&s.Title)
//		err = rows4.Scan(&s.Credits, &s.Goals, &s.SkillsCompetences, &s.Objectives, &s.LearningOutcomes, &s.Prerequisites, &s.Postrequisites, &s.Instructors)
//		err = rows5.Scan(&te.FullName, &te.Degree, &te.Rank, &te.Position, &te.Contacts, &te.Interests)
//		if err != nil {
//			return nil, nil, nil, nil, err
//		}
//
//		syllabus = append(syllabus, s)
//		teacher = append(teacher, te)
//	}
//
//	for rows1.Next() && rows2.Next() {
//		t := &models.TopicWeek{}
//		i := &models.StudentTopicWeek{}
//		err = rows1.Scan(&t.LectureTopic, &t.LectureHours, &t.PracticeTopic, &t.PracticeHours, &t.Assignment, &t.WeekNumber)
//		err = rows2.Scan(&i.WeekNumber, &i.Topics, &i.Hours, &i.RecommendedLiterature, &i.SubmissionForm)
//
//		topic = append(topic, t)
//		independent = append(independent, i)
//	}
//
//	if err = rows1.Err(); err != nil {
//		return nil, nil, nil, nil, err
//	}
//
//	return topic, independent, syllabus, teacher, nil
//}

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

func (m *PgModel) InsertSyllabusInfo(syllabus *models.Syllabus) (int, error) {
	var syllabusInfoId uint32
	row := m.Pool.QueryRow(context.Background(), insertSyllabusInfo,
		syllabus.Credits, syllabus.Goals, syllabus.SkillsCompetences, syllabus.Objectives,
		syllabus.LearningOutcomes, syllabus.Prerequisites, syllabus.Postrequisites, syllabus.Instructors)
	err := row.Scan(&syllabusInfoId)
	if err != nil {
		return 0, err
	}
	return int(syllabusInfoId), nil
}

func (m *PgModel) InsertSessionPlan(table1 []*models.TopicWeek, syllabusInfoId int) (int, error) {
	var planId uint32
	row := m.Pool.QueryRow(context.Background(), insertSessionPlan,
		0, syllabusInfoId)
	err := row.Scan(&planId)
	if err != nil {
		return 0, err
	}
	for _, topic := range table1 {
		var topicId uint32
		row := m.Pool.QueryRow(context.Background(), insertSessionTopic,
			topic.LectureTopic, topic.LectureHours, topic.PracticeTopic, topic.PracticeHours, topic.Assignment, topic.WeekNumber, planId)
		err := row.Scan(&topicId)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	return int(planId), nil
}

func (m *PgModel) InsertIndependentStudyPlan(table2 []*models.StudentTopicWeek, syllabusInfoId int) (int, error) {
	var independentStudyPlanId uint32
	row := m.Pool.QueryRow(context.Background(), insertIndependentStudyPlan,
		0, syllabusInfoId)
	err := row.Scan(&independentStudyPlanId)
	if err != nil {
		return 0, err
	}
	for _, topic := range table2 {
		var topicId uint32
		row := m.Pool.QueryRow(context.Background(), insertIndependentStudyPlanTopic,
			topic.WeekNumber, topic.Topics, topic.Hours, topic.RecommendedLiterature, topic.SubmissionForm, independentStudyPlanId)
		err := row.Scan(&topicId)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	return int(independentStudyPlanId), nil
}

func (m *PgModel) InsertSyllabus(syllabus *models.Syllabus, teacherId int, disciplineId int, name string) (int, error) {
	var syllabusId uint32
	syllabusInfoId, err := m.InsertSyllabusInfo(syllabus)
	_, err = m.InsertSessionPlan(syllabus.Table1, syllabusInfoId)
	_, err = m.InsertIndependentStudyPlan(syllabus.Table2, syllabusInfoId)
	row := m.Pool.QueryRow(context.Background(), insertSyllabus, teacherId, syllabusInfoId, disciplineId, name)
	err = row.Scan(&syllabusId)
	if err != nil {
		return 0, err
	}
	return int(syllabusId), nil
}

func (m *PgModel) GetStudentId() (int, error) {
	var id int
	err := m.Pool.QueryRow(context.Background(), GetStudentId, authID).
		Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Student id", id)
	teacherIDWithSyllabus = id
	return id, nil
}
func (m *PgModel) GetSyllabusById(id int) ([]*models.TopicWeek, []*models.StudentTopicWeek, []*models.Syllabus, []*models.TeacherInfo, error) {
	var topic []*models.TopicWeek
	var independent []*models.StudentTopicWeek
	var syllabus []*models.Syllabus
	var teacher []*models.TeacherInfo

	topic, err := m.selectTopicWithPlan(id)
	independent, err = m.selectIndependentStudyTopic(id)
	syllabus, err = m.selectSyllabusTableRow(id)
	teacher, err = m.selectTeacherInfo()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return topic, independent, syllabus, teacher, nil
}
func (m *PgModel) selectTopicWithPlan(id int) ([]*models.TopicWeek, error) {
	var topic []*models.TopicWeek
	rows1, err := m.Pool.Query(context.Background(), selectTopicWithPlan, id)
	if err != nil {
		return nil, err
	}
	for rows1.Next() {
		t := &models.TopicWeek{}
		err = rows1.Scan(&t.LectureTopic, &t.LectureHours, &t.PracticeTopic, &t.PracticeHours, &t.Assignment, &t.WeekNumber)
		topic = append(topic, t)
	}
	if err = rows1.Err(); err != nil {
		return nil, err
	}
	return topic, nil
}

func (m *PgModel) selectIndependentStudyTopic(id int) ([]*models.StudentTopicWeek, error) {
	var independent []*models.StudentTopicWeek
	rows2, err := m.Pool.Query(context.Background(), selectIndependentStudyTopic, id)
	if err != nil {
		return nil, err
	}
	for rows2.Next() {
		i := &models.StudentTopicWeek{}
		err = rows2.Scan(&i.WeekNumber, &i.Topics, &i.Hours, &i.RecommendedLiterature, &i.SubmissionForm)
		independent = append(independent, i)
	}
	if err = rows2.Err(); err != nil {
		return nil, err
	}
	return independent, nil
}

func (m *PgModel) selectSyllabusTableRow(id int) ([]*models.Syllabus, error) {
	var syllabus []*models.Syllabus
	rows3, err := m.Pool.Query(context.Background(), selectSyllabusTableRow, id)
	rows4, err := m.Pool.Query(context.Background(), selectSyllabusInfo, id)
	if err != nil {
		return nil, err
	}
	for rows3.Next() && rows4.Next() {

		s := &models.Syllabus{}

		err = rows3.Scan(&s.Title)
		err = rows4.Scan(&s.Credits, &s.Goals, &s.SkillsCompetences, &s.Objectives, &s.LearningOutcomes, &s.Prerequisites, &s.Postrequisites, &s.Instructors)

		if err != nil {
			return nil, err
		}

		syllabus = append(syllabus, s)
	}
	if err = rows3.Err(); err != nil {
		return nil, err
	}
	return syllabus, nil
}

func (m *PgModel) selectTeacherInfo() ([]*models.TeacherInfo, error) {
	var teacher []*models.TeacherInfo

	rows5, err := m.Pool.Query(context.Background(), selectTeacherInfo, authID)
	if err != nil {
		return nil, err
	}

	for rows5.Next() {
		te := &models.TeacherInfo{}

		err = rows5.Scan(&te.FullName, &te.Degree, &te.Rank, &te.Position, &te.Contacts, &te.Interests)
		if err != nil {
			return nil, err
		}

		teacher = append(teacher, te)
	}

	return teacher, nil
}
