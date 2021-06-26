package main

import (
	"errors"
	"examFortune/pkg/models"
	"net/http"
)

func (app *application) AdminHomePage(w http.ResponseWriter, r *http.Request) {
	adminInfo, err := app.student.GetInfoForAdmin()

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	flash := app.session.PopString(r, "flash")

	app.render(w, r, "TrueAdmin.page.tmpl", &templateData{
		Flash: flash,
		Admin: adminInfo,
	})
}
