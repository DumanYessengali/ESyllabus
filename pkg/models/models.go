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
	TopicWeekID   int
	WeekNumber    int
	LectureTopic  string
	LectureHours  int
	PracticeTopic string
	PracticeHours int
	Assignment    string
}

type StudentTopicWeek struct {
	StudentTopicWeekID    int
	WeekNumber            int
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
