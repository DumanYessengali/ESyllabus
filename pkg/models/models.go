package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Syllabus struct {
	ID                int
	Title             string
	Teacher           *TeacherInfo
	Credits           int
	Status            string
	Feedback          string
	Discipline        string
	Goals             string
	SkillsCompetences string
	Objectives        string
	LearningOutcomes  string
	Prerequisites     string
	Postrequisites    string
	Instructors       string
	SyllabusInfoID    int
	Assessment        int
	Table1            []*TopicWeek
	Table2            []*StudentTopicWeek
	TempGoals         []string
	TempObjectives    []string
	TempOutcomes      []string
}
type Discipline struct {
	DisciplineId int
	Title        string
}

type TeacherInfo struct {
	ID        int
	FullName  string
	Degree    string
	Rank      string
	Position  string
	Contacts  string
	Interests string
}

type TopicWeek struct {
	TopicWeekID    int
	WeekNumber     int
	SyllabusInfoId int
	TeacherId      int
	TeacherName    string
	LectureTopic   string
	LectureHours   int
	PracticeTopic  string
	PracticeHours  int
	Assignment     string
}

type StudentTopicWeek struct {
	StudentTopicWeekID    int
	SyllabusInfoId        int
	WeekNumber            int
	TeacherId             int
	TeacherName           string
	Topics                string
	Hours                 int
	RecommendedLiterature string
	SubmissionForm        string
}

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

type Assessment struct {
	AssessmentId int
	Assignment1  []string
	PointsNum1   []string
	Assignment2  []string
	PointsNum2   []string
}

//----------------------FOR TEMP

type SessionWeek struct {
	SyllabusInfoId int
	WeekNumber     int
	LectureTopic   []*TempFields
	LectureHours   []*TempFields
	PracticeTopic  []*TempFields
	PracticeHours  []*TempFields
	Assignment     []*TempFields
}

type StudentWeek struct {
	SyllabusInfoId        int
	WeekNumber            int
	Topics                []*TempFields
	Hours                 []*TempFields
	RecommendedLiterature []*TempFields
	SubmissionForm        []*TempFields
}

type TempFields struct {
	TeacherId      int
	TeacherName    string
	SyllabusInfoId int
	Content        string
	ContentInt     int
}

type TempSessionTopic struct {
	SyllabusInfoId int
	TeacherID      int
	TeacherName    string
	WeekNumber     int
	Topic          string
	Hours          int
}

type TempStudentTopic struct {
	SyllabusInfoId        int
	TeacherId             int
	TeacherName           string
	WeekNumber            int
	Topic                 string
	Hours                 int
	RecommendedLiterature string
	SubmissionForm        string
}
