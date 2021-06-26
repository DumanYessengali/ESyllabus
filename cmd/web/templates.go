package main

import (
	"examFortune/pkg/forms"
	"examFortune/pkg/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Form             *forms.Form
	Flash            string
	Prediction       string
	Time             bool
	AssessmentType   *models.Assessment
	syllabus         *models.Syllabus
	Syllabus         []*models.Syllabus
	Admin            []*models.User
	TopicOneRow      *models.SessionWeek
	IndepTopicOneRow *models.StudentWeek
	Topic            []*models.TopicWeek
	Independent      []*models.StudentTopicWeek
	Teacher          []*models.TeacherInfo
	Discipline       []*models.Discipline

	TempTable1 []*models.SessionWeek
	TempTable2 []*models.StudentWeek

	TempGoals      []*models.TempFields
	TempObjectives []*models.TempFields
	TempOutcomes   []*models.TempFields

	IsTeacher       bool
	IsStudent       bool
	IsCoordinator   bool
	IsDean          bool
	IsAuthenticated bool
	IsDied          bool
	IsAdmin         bool
	IsNewTeacher    bool
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
