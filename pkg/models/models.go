package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Student struct {
	ID          int
	Username    string
	Password    string
	GroupName   string
	SubjectName string
	LifeTime    int
	IsLast      bool
}
