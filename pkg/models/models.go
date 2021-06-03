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
	Goals             string
	SkillsCompetences string
	Objectives        string
	LearningOutcomes  string
	Prerequisites     string
	Postrequisites    string
	Instructors       string
	SyllabusInfoID    int
	Table1            []*TopicWeek
	Table2            []*StudentTopicWeek
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
