package postgres

import (
	"context"
	"database/sql"
	"errors"
	"examFortune/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"math/rand"
	"time"
)

const (
	insertTeacher = "INSERT INTO"
	insertSql     = "INSERT INTO student (username, password, group_name, subject_name, life_time, is_last)" +
		" VALUES ($1,$2,$3,$4,$5,$6)"
	getNameSyllabus   = "SELECT * FROM student"
	getStudentById    = "SELECT * FROM student WHERE student_id=$1"
	getRoleByUsername = "SELECT authorization_id, role FROM auth WHERE username=$1"
	deleteStudentById = "DELETE FROM  student WHERE student_id=$1"
	updateStudent     = "UPDATE student SET " +
		"username=$1, password=$2, group_name=$3, subject_name=$4, life_time=$5, is_last=$6 " +
		"WHERE student_id = $7"
	auth = "SELECT authorization_id, password FROM auth WHERE username = $1"
)

type StudentModel struct {
	Pool *pgxpool.Pool
}

func (m *StudentModel) InsertSyllabus(username, password, group_name, subject_name string) (int, error) {

	var id uint64
	row := m.Pool.QueryRow(context.Background(), insertSql,
		username, password, group_name, subject_name, 3, false)
	err := row.Scan(id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *StudentModel) GetNameSyllabus() ([]*models.Student, error) {
	var students []*models.Student
	rows, err := m.Pool.Query(context.Background(), getNameSyllabus)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := &models.Student{}
		err = rows.Scan(&s.ID, &s.Username, &s.Password, &s.GroupName, &s.SubjectName, &s.LifeTime, &s.IsLast)
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

func (m *StudentModel) DeleteStudentById(id int) error {
	_, err := m.Pool.Exec(context.Background(), deleteStudentById, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *StudentModel) Authenticate(username, password string) (int, error) {
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
	return id, nil
}
func (m *StudentModel) Get(id int) (*models.Student, error) {
	return nil, nil
}

func (m *StudentModel) GetStudentById(id int) (*models.Student, error) {
	s := &models.Student{}
	err := m.Pool.QueryRow(context.Background(), getStudentById, id).
		Scan(&s.ID, &s.Username, &s.Password, &s.GroupName, &s.SubjectName, &s.LifeTime, &s.IsLast)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *StudentModel) GetRoleByUsername(username string) (*models.Student, error) {
	s := &models.Student{}
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

/*func (m *StudentModel) DeleteStudentByUsername(username string) error {
	_, err := m.Pool.Exec(context.Background(), deleteStudentByUsername, username)
	if err != nil {
		return err
	}
	return nil
}*/

func (m *StudentModel) UpdateStudent(s *models.Student) error {
	_, err := m.Pool.Exec(context.Background(), updateStudent, s.Username, s.Password, s.GroupName, s.SubjectName, s.LifeTime, s.IsLast, s.ID)
	if err != nil {
		return err
	}
	return nil
}

func init() {

	rand.Seed(time.Now().UnixNano())
}
