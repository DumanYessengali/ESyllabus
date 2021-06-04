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
	getNameSyllabusWithTeacher = "select syllabus_id,name,syllabus_info_id from syllabus where teacher_id=$1"
	getNameSyllabusWithStudent = "select syllabus_id,name,syllabus_info_id from syllabus where syllabus_id=(select syllabus_id from student_syllabus where student_id=$1)"
	getTeacherId               = "SELECT teacher_id from teacher where authorization_id=$1"
	getRoleByUsername          = "SELECT authorization_id, role FROM auth WHERE username=$1"
	auth                       = "SELECT authorization_id, password FROM auth WHERE username = $1"

	deleteTopicWithPlan                 = "delete from topic where plan_id=(select plan_id from session_plan where syllabus_info_id=$1)"
	deleteSessionPlanWithSyllabusInfoId = "DELETE FROM session_plan WHERE syllabus_info_id=$1"
	deleteIndependentStudyTopic         = "delete from independent_study_topic where independent_study_plan_id=(select independent_study_plan_id from independent_study_plan where syllabus_info_id=$1)"
	deleteIndependentStudyPlan          = "DELETE FROM independent_study_plan WHERE syllabus_info_id=$1"
	deleteStudentSyllabus               = "delete from student_syllabus where syllabus_id=(select syllabus_id from syllabus where syllabus_info_id=$1)"
	deleteSyllabusTableRow              = "DELETE FROM  syllabus WHERE syllabus_info_id=$1"
	deleteSyllabusInfo                  = "DELETE FROM syllabus_info WHERE syllabus_info_id=$1"

	selectOnlyOneTopic  = "select topic_id, lecture,lecture_hours,practice,practice_hours,assignment,week_number from topic where topic_id=$1 "
	selectOnlyOneIndep  = "select independent_study_topic_id, week_numbers,topics,hours,recommended_literature,sudmission_form from independent_study_topic where independent_study_topic_id=$1 "
	selectTopicWithPlan = "select topic_id, lecture,lecture_hours,practice,practice_hours,assignment,week_number " +
		"from topic where plan_id=(select plan_id from session_plan where syllabus_info_id=$1)"
	selectIndependentStudyTopic = "select independent_study_topic_id, week_numbers,topics,hours,recommended_literature,sudmission_form " +
		"from independent_study_topic where independent_study_plan_id=(select independent_study_plan_id from independent_study_plan where syllabus_info_id=$1)"
	selectSyllabusTableRow = "select name FROM  syllabus WHERE syllabus_info_id=$1"
	selectSyllabusInfo     = "select syllabus_info_id,credits_num,goals,skills_competences,objectives,learning_outcomes,prerequisites,postrequisites,instructors " +
		"FROM syllabus_info WHERE syllabus_info_id=$1"
	selectTeacherInfo = "select fullname, degree, rank, position, contacts, interests from teacher where authorization_id = $1"

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

	updateTopicWeek = "update topic set lecture=$1, lecture_hours=$2, " +
		"practice=$3, practice_hours=$4, assignment=$5, week_number=$6 where topic_id=$7"

	updateSyllabusInfo = "update syllabus_info " +
		"set credits_num=$1, goals=$2, skills_competences=$3, objectives=$4, learning_outcomes=$5," +
		"prerequisites=$6, postrequisites=$7, instructors=$8 " +
		"where syllabus_info_id=$9"

	updateStudentTopicWeek = "update independent_study_topic set week_numbers=$1, topics=$2, " +
		"hours=$3, recommended_literature=$4, sudmission_form=$5 where independent_study_topic_id=$6"
)

type PgModel struct {
	Pool *pgxpool.Pool
}

var authID int
var iDFromSyllabus int
var Role string

func (m *PgModel) GetRole() string {
	return Role
}

func (m *PgModel) GetNameSyllabusWithStudent() ([]*models.Syllabus, error) {
	var students []*models.Syllabus
	rows, err := m.Pool.Query(context.Background(), getNameSyllabusWithStudent, iDFromSyllabus)
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

func (m *PgModel) GetNameSyllabus() ([]*models.Syllabus, error) {
	var students []*models.Syllabus
	rows, err := m.Pool.Query(context.Background(), getNameSyllabusWithTeacher, iDFromSyllabus)
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
	iDFromSyllabus = id
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
	fmt.Println(s.Role)
	Role = s.Role
	return s, nil
}

//
///*func (m *PgModel) DeleteStudentByUsername(username string) error {
//  _, err := m.Pool.Exec(context.Background(), deleteStudentByUsername, username)
//  if err != nil {
//      return err
//  }
//  return nil
//}*/
//
//func (m *PgModel) UpdateStudent(s *models.Student) error {
//  _, err := m.Pool.Exec(context.Background(), updateStudent, s.Username, s.Password, s.GroupName, s.SubjectName, s.LifeTime, s.IsLast, s.ID)
//  if err != nil {
//      return err
//  }
//  return nil
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
	iDFromSyllabus = id
	return id, nil
}
func (m *PgModel) GetSyllabusById(id int) ([]*models.TopicWeek, []*models.StudentTopicWeek, []*models.Syllabus, []*models.TeacherInfo, error) {
	var topic []*models.TopicWeek
	var independent []*models.StudentTopicWeek
	var syllabus []*models.Syllabus
	var teacher []*models.TeacherInfo

	topic, err := m.SelectTopicWithPlan(id)
	independent, err = m.selectIndependentStudyTopic(id)
	syllabus, err = m.SelectSyllabusTableRow(id)
	teacher, err = m.selectTeacherInfo()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return topic, independent, syllabus, teacher, nil
}
func (m *PgModel) SelectTopicWithPlan(id int) ([]*models.TopicWeek, error) {
	var topic []*models.TopicWeek
	rows1, err := m.Pool.Query(context.Background(), selectTopicWithPlan, id)
	if err != nil {
		return nil, err
	}
	for rows1.Next() {
		t := &models.TopicWeek{}
		err = rows1.Scan(&t.TopicWeekID, &t.LectureTopic, &t.LectureHours, &t.PracticeTopic, &t.PracticeHours, &t.Assignment, &t.WeekNumber)
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
		err = rows2.Scan(&i.StudentTopicWeekID, &i.WeekNumber, &i.Topics, &i.Hours, &i.RecommendedLiterature, &i.SubmissionForm)
		independent = append(independent, i)

	}
	if err = rows2.Err(); err != nil {
		return nil, err
	}
	return independent, nil
}

func (m *PgModel) SelectSyllabusTableRow(id int) ([]*models.Syllabus, error) {
	var syllabus []*models.Syllabus
	rows3, err := m.Pool.Query(context.Background(), selectSyllabusTableRow, id)
	rows4, err := m.Pool.Query(context.Background(), selectSyllabusInfo, id)
	if err != nil {
		return nil, err
	}
	for rows3.Next() && rows4.Next() {

		s := &models.Syllabus{}

		err = rows3.Scan(&s.Title)
		err = rows4.Scan(&s.SyllabusInfoID, &s.Credits, &s.Goals, &s.SkillsCompetences, &s.Objectives, &s.LearningOutcomes, &s.Prerequisites, &s.Postrequisites, &s.Instructors)

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

func (m *PgModel) UpdateSyllabusInfo(syllabus *models.Syllabus, id int) error {
	_, err := m.Pool.Exec(context.Background(), updateSyllabusInfo,
		syllabus.Credits, syllabus.Goals, syllabus.SkillsCompetences, syllabus.Objectives,
		syllabus.LearningOutcomes, syllabus.Prerequisites, syllabus.Postrequisites,
		syllabus.Instructors, id)

	fmt.Println(syllabus.Credits, syllabus.Prerequisites, syllabus.Postrequisites)
	if err != nil {
		return err
	}
	return nil
}

func (m *PgModel) UpdateTopicWeek(tw *models.TopicWeek, id int) error {

	_, err := m.Pool.Exec(context.Background(), updateTopicWeek,
		tw.LectureTopic, tw.LectureHours, tw.PracticeTopic, tw.PracticeHours,
		tw.Assignment, tw.WeekNumber, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PgModel) SelecOnlyOneTopic(id int) (*models.TopicWeek, error) {
	t := &models.TopicWeek{}
	err := m.Pool.QueryRow(context.Background(), selectOnlyOneTopic, id).
		Scan(&t.TopicWeekID, &t.LectureTopic, &t.LectureHours, &t.PracticeTopic, &t.PracticeHours, &t.Assignment, &t.WeekNumber)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return t, nil
}

func (m *PgModel) UpdateStudentTopicWeek(tw *models.StudentTopicWeek, id int) error {

	_, err := m.Pool.Exec(context.Background(), updateStudentTopicWeek,
		tw.WeekNumber, tw.Topics, tw.Hours, tw.RecommendedLiterature, tw.SubmissionForm, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PgModel) SelecOnlyOneIndep(id int) (*models.StudentTopicWeek, error) {
	i := &models.StudentTopicWeek{}
	err := m.Pool.QueryRow(context.Background(), selectOnlyOneIndep, id).
		Scan(&i.StudentTopicWeekID, &i.WeekNumber, &i.Topics, &i.Hours, &i.RecommendedLiterature, &i.SubmissionForm)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return i, nil
}
