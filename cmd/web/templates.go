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
	syllabus         *models.Syllabus
	Syllabus         []*models.Syllabus
	TopicOneRow      *models.TopicWeek
	IndepTopicOneRow *models.StudentTopicWeek
	Topic            []*models.TopicWeek
	Independent      []*models.StudentTopicWeek
	Teacher          []*models.TeacherInfo
	IsAdmin          bool
	IsAuthenticated  bool
	IsDied           bool
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
