package postgres

//
//
//import (
//"context"
//"examFortune/pkg/models"
//	"fmt"
//	"github.com/jackc/pgx/v4/pgxpool"
//)
//
//
//const (
//	getSyllabusByTeacherId = "select syllabus_info_id from syllabus where teacher_id=$1"
//	getSyllabusInfoBySyllabusId = "select * from syllabus s inner join syllabus_info si " +
//		"on s.syllabus_info_id=si.syllabus_info_id where s.syllabus_id=$1"
//	getTable1 = "select * from session_plan sp inner join topic t " +
//		"on sp.plan_id=t.plan_id where sp.syllabus_info_id=$1 order by t.week_number"
//	getTable2 = "select * from independent_study_plan isp inner join independent_study_topic t " +
//		"on isp.independent_study_plan_id=t.independent_study_plan_id where isp.syllabus_info_id=&1 " +
//		"order by t.week_numbers"
//	getTeacherInfoByTeacherId = "select * from teacher where teacher_id=$1"
//)
//
//type SyllabusModel struct {
//	Pool *pgxpool.Pool
//}
//
//
//func (m *SyllabusModel) GetAllSyllabuses(teacherId int) ([]*models.Syllabus, error) {
//	var syllabuses []*models.Syllabus
//	var syllabusIds []int
//
//	fmt.Println("aLL SYLLABUS")
//	rows, err := m.Pool.Query(context.Background(), getSyllabusByTeacherId, teacherId)
//	fmt.Println("aLL SYLLABUS")
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		var s int
//		err = rows.Scan(&s)
//		if err != nil {
//			return nil, err
//		}
//
//		syllabusIds = append(syllabusIds, s)
//	}
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	for _, s := range syllabusIds{
//		syllabuses = append(syllabuses, m.GetSyllabus(s))
//	}
//	return syllabuses, nil
//}
//
//func (m * SyllabusModel) GetSyllabus(syllabusId int) (*models.Syllabus) {
//	syllabus := &models.Syllabus{}
//	var s string
//	err := m.Pool.QueryRow(context.Background(), getSyllabusInfoBySyllabusId, syllabusId).
//		Scan(&syllabus.ID, &syllabus.Teacher.ID, &syllabus.SyllabusInfoID, &s,
//			&syllabus.Title, &syllabus.Credits, &syllabus.Goals, &syllabus.SkillsCompetences,
//			&syllabus.Objectives, &syllabus.LearningOutcomes, &syllabus.Prerequisites,
//			&syllabus.Postrequisites, &syllabus.Instructors)
//	if err != nil {
//		if err.Error() == "no rows in result set" {
//			return nil
//		} else {
//			return nil
//		}
//	}
//	fmt.Println("Sylabus")
//	syllabus.Teacher, _ = m.GetTeacherInfo(syllabus.Teacher.ID)
//	syllabus.Table1, _ = m.GetTable1(syllabus.SyllabusInfoID)
//	syllabus.Table2, _ = m.GetTable2(syllabus.SyllabusInfoID)
//
//	return syllabus
//}
//
//func (m * SyllabusModel) GetTeacherInfo(teacherId int) (*models.TeacherInfo, error) {
//	teacher := &models.TeacherInfo{}
//	err := m.Pool.QueryRow(context.Background(), getTeacherInfoByTeacherId, teacherId).
//		Scan(&teacher.ID, &teacher.FullName, &teacher.Degree, &teacher.Rank,
//			&teacher.Position, &teacher.Contacts, &teacher.Interests)
//	if err != nil {
//		if err.Error() == "no rows in result set" {
//			return nil, models.ErrNoRecord
//		} else {
//			return nil, err
//		}
//	}
//	fmt.Println("Teacher")
//	return teacher, nil
//}
//
//func (m *SyllabusModel) GetTable1(syllabus_info_id int) ([]*models.TopicWeek, error) {
//	topicWeek := []*models.TopicWeek{}
//
//	rows, err := m.Pool.Query(context.Background(), getTable1, syllabus_info_id)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		tw := &models.TopicWeek{}
//		var s int
//		err = rows.Scan(s, &tw.LectureTopic, &tw.LectureHours, &tw.PracticeTopic, &tw.PracticeHours)
//		if err != nil {
//			return nil, err
//		}
//
//		topicWeek = append(topicWeek, tw)
//	}
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//	fmt.Println("T1")
//	return topicWeek, nil
//}
//
//func (m *SyllabusModel) GetTable2(syllabus_info_id int) ([]*models.StudentTopicWeek, error) {
//	topicWeek := []*models.StudentTopicWeek{}
//
//	rows, err := m.Pool.Query(context.Background(), getTable2, syllabus_info_id)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		tw := &models.StudentTopicWeek{}
//		var s int
//		err = rows.Scan(&s, &s, &s, &s, &tw.Topics, &tw.Hours, &tw.RecommendedLiterature,
//			&tw.SubmissionForm)
//		if err != nil {
//			return nil, err
//		}
//
//		topicWeek = append(topicWeek, tw)
//	}
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//	fmt.Println("T2")
//	return topicWeek, nil
//}
