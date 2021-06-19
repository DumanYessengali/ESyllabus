package postgres

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"examFortune/pkg/models"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

const (
	getNameSyllabusWithTeacher    = "select syllabus_id,name,syllabus_info_id,status from syllabus where teacher_id=$1"
	getNameSyllabusForCoordinator = "select s.syllabus_id,s.name,s.syllabus_info_id,s.status,s.feedback\nfrom syllabus s\ninner join syllabus_info si on s.syllabus_info_id = si.syllabus_info_id\ninner join discipline d on si.discipline_id = d.discipline_id\ninner join coordinator_discipline cd on d.discipline_id = cd.discipline_id\nwhere cd.coordinator_id=$1 and s.status=$2"
	getNameSyllabusWithStudent    = "select s.syllabus_id, s.name, s.syllabus_info_id\nfrom syllabus s\n    inner join student_syllabus ss on s.syllabus_id = ss.syllabus_id\nwhere ss.student_id=$1 and status = $2"
	getTeacherId                  = "SELECT teacher_id from teacher where authorization_id=$1"
	getTeacherIds                 = "SELECT teacher_id from teacher_discipline where discipline_id=$1"
	getTeacherName                = "SELECT username from auth a inner join teacher t on a.authorization_id = t.authorization_id where teacher_id=$1"
	getCoordinatorId              = "SELECT coordinator_id from coordinator where authid=$1"
	getRoleByUsername             = "SELECT authorization_id, role FROM auth WHERE username=$1"
	auth                          = "SELECT authorization_id, password FROM auth WHERE username = $1"
	selectTeacherInfoById         = "select fullname, degree, rank, position, contacts, interests from teacher where teacher_id=$1"

	deleteTopicWithPlan                 = "delete from topic where plan_id=(select plan_id from session_plan where syllabus_info_id=$1)"
	deleteSessionPlanWithSyllabusInfoId = "DELETE FROM session_plan WHERE syllabus_info_id=$1"
	deleteIndependentStudyTopic         = "delete from independent_study_topic where independent_study_plan_id=(select independent_study_plan_id from independent_study_plan where syllabus_info_id=$1)"
	deleteIndependentStudyPlan          = "DELETE FROM independent_study_plan WHERE syllabus_info_id=$1"
	deleteStudentSyllabus               = "delete from student_syllabus where syllabus_id=(select syllabus_id from syllabus where syllabus_info_id=$1)"
	deleteSyllabusTableRow              = "DELETE FROM  syllabus WHERE syllabus_info_id=$1"
	deleteSyllabusInfo                  = "DELETE FROM syllabus_info WHERE syllabus_info_id=$1"

	selectAllDiscipline = "select discipline_id,title from discipline"
	selectCredits       = "select title,credit from discipline where discipline_id=(select discipline_id from syllabus_info where syllabus_info_id=$1)"
	selectOnlyOneTopic  = "select topic_id, lecture,lecture_hours,practice,practice_hours,assignment,week_number from topic where topic_id=$1 "
	selectOnlyOneIndep  = "select independent_study_topic_id, week_numbers,topics,hours,recommended_literature,sudmission_form from independent_study_topic where independent_study_topic_id=$1 "
	selectTopicWithPlan = "select topic_id, lecture,lecture_hours,practice,practice_hours,assignment,week_number " +
		"from topic where plan_id=(select plan_id from session_plan where syllabus_info_id=$1) order by week_number"
	selectIndependentStudyTopic = "select independent_study_topic_id, week_numbers,topics,hours,recommended_literature,sudmission_form " +
		"from independent_study_topic where independent_study_plan_id=(select independent_study_plan_id from independent_study_plan where syllabus_info_id=$1) order by week_numbers"
	selectSyllabusTableRow = "select name,feedback FROM  syllabus WHERE syllabus_info_id=$1"
	selectSyllabusInfo     = "select syllabus_info_id,goals,skills_competences,objectives,learning_outcomes,prerequisites,postrequisites,instructors " +
		"FROM syllabus_info WHERE syllabus_info_id=$1"
	selectTeacherInfo     = "select t.fullname, t.degree, t.rank, t.position, t.contacts, t.interests\nfrom teacher t\ninner join syllabus s on t.teacher_id = s.teacher_id\nwhere s.syllabus_info_id=$1;"
	selectTeacherByAuthId = ""

	insertSyllabus     = "insert into syllabus (teacher_id, syllabus_info_id, name,status, feedback) values ($1, $2,$3,$4,$5) returning syllabus_id"
	insertDiscipline   = "update syllabus_info set discipline_id=$1 where syllabus_info_id=$2 returning discipline_id"
	insertSyllabusInfo = "insert into syllabus_info( goals, skills_competences, objectives, learning_outcomes," +
		"prerequisites, postrequisites, instructors,assessment_id) values($1, $2, $3, $4, $5, $6, $7, $8) returning syllabus_info_id"
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
		"set goals=$1, skills_competences=$2, objectives=$3, learning_outcomes=$4," +
		"prerequisites=$5, postrequisites=$6, instructors=$7 " +
		"where syllabus_info_id=$8"

	updateStudentTopicWeek = "update independent_study_topic set week_numbers=$1, topics=$2, " +
		"hours=$3, recommended_literature=$4, sudmission_form=$5 where independent_study_topic_id=$6"

	selectAssessment = "select assessment_id, assessment_title1, points_num1, assessment_title2, points_num2 from assessment " +
		"where assessment_id=(select assessment_id from syllabus_info where syllabus_info_id=$1)"

	updateStatus           = "update syllabus set status=$1 where syllabus_info_id=$2"
	updateFeedback         = "update syllabus set feedback = $1 where syllabus_info_id = $2 returning syllabus_id"
	updateSyllabusInfoTemp = "update syllabus_info set goals=$1, objectives=$2, learning_outcomes=$3 " +
		"where syllabus_info_id=$4"

	getSyllabusForDean = "select syllabus_id,syllabus_info_id,status,name, feedback from syllabus where status=$1"
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
	rows, err := m.Pool.Query(context.Background(), getNameSyllabusWithStudent, iDFromSyllabus, "ready")

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

func (m *PgModel) GetSyllabusForDean(status string) ([]*models.Syllabus, error) {
	var students []*models.Syllabus
	rows, err := m.Pool.Query(context.Background(), getSyllabusForDean, status)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &models.Syllabus{}
		err = rows.Scan(&s.ID, &s.SyllabusInfoID, &s.Status, &s.Title, &s.Feedback)
		if err != nil {
			return nil, err
		}
		if s.Status == status {
			students = append(students, s)
		}

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (m *PgModel) GetNameSyllabus(status string) ([]*models.Syllabus, error) {
	var students []*models.Syllabus
	rows, err := m.Pool.Query(context.Background(), getNameSyllabusWithTeacher, iDFromSyllabus)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := &models.Syllabus{}
		err = rows.Scan(&s.ID, &s.Title, &s.SyllabusInfoID, &s.Status)
		if err != nil {
			return nil, err
		}
		if s.Status == status {
			students = append(students, s)
		}

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (m *PgModel) GetFullInfoByTeacherId() (*models.TeacherInfo, error) {
	var Fullname string
	var Degree string
	var Rank string
	var Position string
	var Contacts string
	var Interests string
	err := m.Pool.QueryRow(context.Background(), selectTeacherInfoById, iDFromSyllabus).
		Scan(&Fullname, &Degree, &Rank, &Position, &Contacts, &Interests)

	if err != nil {
		fmt.Println(err.Error())
	}

	return &models.TeacherInfo{
		iDFromSyllabus,
		Fullname,
		Degree,
		Rank,
		Position,
		Contacts,
		Interests,
	}, nil
}

func (m *PgModel) GetNameSyllabusFromCoordinator(status string) ([]*models.Syllabus, error) {
	var students []*models.Syllabus
	rows, err := m.Pool.Query(context.Background(), getNameSyllabusForCoordinator, iDFromSyllabus, status)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &models.Syllabus{}
		err = rows.Scan(&s.ID, &s.Title, &s.SyllabusInfoID, &s.Status, &s.Feedback)
		if err != nil {
			return nil, err
		}
		if s.Status == status {
			students = append(students, s)
		}
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
	iDFromSyllabus = id
	return id, nil
}

func (m *PgModel) GetTeacherIds(dId int) ([]int, error) {
	var ids []int
	rows1, err := m.Pool.Query(context.Background(), getTeacherIds, dId)
	if err != nil {
		return nil, err
	}
	for rows1.Next() {
		var id int
		err = rows1.Scan(&id)
		ids = append(ids, id)
	}

	if err = rows1.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (m *PgModel) GetTeacherUsername(teacherId int) (string, error) {
	var username string
	err := m.Pool.QueryRow(context.Background(), getTeacherName, teacherId).
		Scan(&username)
	if err != nil {
		fmt.Println(err.Error())
	}
	return username, nil
}

func (m *PgModel) GetCoordinatorId() (int, error) {
	var id int
	err := m.Pool.QueryRow(context.Background(), getCoordinatorId, authID).
		Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
	}
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

func (m *PgModel) DeleteForRejection(id int) error {
	_, err := m.Pool.Exec(context.Background(), deleteTopicWithPlan, id)
	_, err = m.Pool.Exec(context.Background(), deleteSessionPlanWithSyllabusInfoId, id)
	_, err = m.Pool.Exec(context.Background(), deleteIndependentStudyTopic, id)
	_, err = m.Pool.Exec(context.Background(), deleteIndependentStudyPlan, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *PgModel) SendSyllabus(id int, status string) error {
	_, err := m.Pool.Exec(context.Background(), updateStatus, status, id)
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
	passwordWithoutHash := []byte(password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	md5HashInBytes := md5.Sum(passwordWithoutHash)
	md5HashInString := hex.EncodeToString(md5HashInBytes[:])

	if pass != md5HashInString {
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
	Role = s.Role
	return s, nil
}

func (m *PgModel) InsertSyllabusInfo(syllabus *models.Syllabus) (int, error) {
	var syllabusInfoId uint32
	row := m.Pool.QueryRow(context.Background(), insertSyllabusInfo,
		syllabus.Goals, syllabus.SkillsCompetences, syllabus.Objectives,
		syllabus.LearningOutcomes, syllabus.Prerequisites, syllabus.Postrequisites, syllabus.Instructors, syllabus.Assessment)
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
			return 0, err
		}
	}
	return int(independentStudyPlanId), nil
}

func (m *PgModel) InsertSyllabus(syllabus *models.Syllabus, teacherId int, name string) (int, int, error) {
	var syllabusId uint32
	syllabusInfoId, err := m.InsertSyllabusInfo(syllabus)
	//_, err = m.InsertSessionPlan(syllabus.Table1, syllabusInfoId)
	//_, err = m.InsertIndependentStudyPlan(syllabus.Table2, syllabusInfoId)
	row := m.Pool.QueryRow(context.Background(), insertSyllabus, teacherId, syllabusInfoId, name, "in_process", "")
	err = row.Scan(&syllabusId)
	if err != nil {
		return 0, 0, err
	}

	return int(syllabusId), syllabusInfoId, nil
}

func (m *PgModel) InsertSyllabusForOtherTeachers(teacherId int, name string, sId int) (int, error) {
	var syllabusId uint32
	fmt.Println(teacherId, name, sId)
	row := m.Pool.QueryRow(context.Background(), insertSyllabus, teacherId, sId, name, "in_process", "")
	err := row.Scan(&syllabusId)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(syllabusId), nil
}

func (m *PgModel) InsertDiscipline(dId, sId int) (int64, error) {
	var id uint32
	row := m.Pool.QueryRow(context.Background(), insertDiscipline,
		dId, sId)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

func (m *PgModel) InsertFeedback(feed string, sId int) (int64, error) {
	var id uint32
	row := m.Pool.QueryRow(context.Background(), updateFeedback, feed, sId)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

func (m *PgModel) GetStudentId() (int, error) {
	var id int
	err := m.Pool.QueryRow(context.Background(), GetStudentId, authID).
		Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
	}
	iDFromSyllabus = id
	return id, nil
}

func (m *PgModel) GetSyllabusById(id int) ([]*models.TopicWeek, []*models.StudentTopicWeek, []*models.Syllabus, []*models.TeacherInfo, *models.Assessment, error) {
	var topic []*models.TopicWeek
	var independent []*models.StudentTopicWeek
	var syllabus []*models.Syllabus
	var teacher []*models.TeacherInfo
	var assessment *models.Assessment

	topic, err := m.SelectTopicWithPlan(id)
	independent, err = m.selectIndependentStudyTopic(id)
	syllabus, err = m.SelectSyllabusTableRow(id)
	teacher, err = m.selectTeacherInfo(id)
	assessment, err = m.SelectAssesmentInfo(id)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return topic, independent, syllabus, teacher, assessment, nil
}

func (m *PgModel) GetSyllabusTeacherInProcess(id int) ([]*models.Syllabus, []*models.TeacherInfo, *models.Assessment, error) {
	var syllabus []*models.Syllabus
	var teacher []*models.TeacherInfo
	var assessment *models.Assessment

	syllabus, err := m.SelectSyllabusTableRow(id)
	teacher, err = m.selectTeacherInfo(id)
	assessment, err = m.SelectAssesmentInfo(id)

	if err != nil {
		return nil, nil, nil, err
	}

	return syllabus, teacher, assessment, nil
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

	rows3, err := m.Pool.Query(context.Background(), selectSyllabusTableRow, id)
	rows4, err := m.Pool.Query(context.Background(), selectSyllabusInfo, id)
	if err != nil {
		return nil, err
	}
	var syllabus []*models.Syllabus
	for rows3.Next() && rows4.Next() {

		s := &models.Syllabus{}
		err = rows3.Scan(&s.Title, &s.Feedback)
		err = rows4.Scan(&s.SyllabusInfoID, &s.Goals, &s.SkillsCompetences, &s.Objectives, &s.LearningOutcomes, &s.Prerequisites, &s.Postrequisites, &s.Instructors)
		if err != nil {
			return nil, err
		}
		s.Discipline, s.Credits, err = m.SelectSyllabusDiscipline(id)
		syllabus = append(syllabus, s)
	}
	if err = rows3.Err(); err != nil {
		return nil, err
	}
	return syllabus, nil
}

func (m *PgModel) SelectSyllabusDiscipline(id int) (string, int, error) {

	rows2, err := m.Pool.Query(context.Background(), selectCredits, id)
	if err != nil {
		return "", 0, err
	}
	var cred int
	var discipline string
	for rows2.Next() {
		err = rows2.Scan(&discipline, &cred)
		if err != nil {
			return "", 0, err
		}

	}
	return discipline, cred, nil
}

func (m *PgModel) selectTeacherInfo(id int) ([]*models.TeacherInfo, error) {
	var teacher []*models.TeacherInfo

	rows5, err := m.Pool.Query(context.Background(), selectTeacherInfo, id)
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
		syllabus.Goals, syllabus.SkillsCompetences, syllabus.Objectives,
		syllabus.LearningOutcomes, syllabus.Prerequisites, syllabus.Postrequisites,
		syllabus.Instructors, id)

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

func (m *PgModel) SelectAllDiscipline() ([]*models.Discipline, error) {
	var discipline []*models.Discipline

	rows5, err := m.Pool.Query(context.Background(), selectAllDiscipline)
	if err != nil {
		return nil, err
	}

	for rows5.Next() {
		d := &models.Discipline{}

		err = rows5.Scan(&d.DisciplineId, &d.Title)
		if err != nil {
			return nil, err
		}

		discipline = append(discipline, d)
	}
	return discipline, nil
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

func (m *PgModel) SelectAssesmentInfo(id int) (*models.Assessment, error) {
	a := &models.Assessment{}
	var ass1 string
	var point1 string
	var ass2 string
	var point2 string

	err := m.Pool.QueryRow(context.Background(), selectAssessment, id).
		Scan(&a.AssessmentId, &ass1, &point1, &ass2, &point2)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	a.Assignment1 = strings.Split(ass1, "//n")
	a.PointsNum1 = strings.Split(point1, "//n")
	a.Assignment2 = strings.Split(ass2, "//n")
	a.PointsNum2 = strings.Split(point2, "//n")

	return a, nil
}
func (m *PgModel) UpdateSyllabusInfoTemp(goals string, objectives string, outcomes string, id int) error {
	_, err := m.Pool.Exec(context.Background(), updateSyllabusInfoTemp,
		goals, objectives, outcomes, id)

	if err != nil {
		return err
	}
	return nil
}
