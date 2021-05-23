package main

import (
	"net/http"
)

func (app *application) getMainPage(w http.ResponseWriter, r *http.Request) {

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "home.page.tmpl", &templateData{
		Flash: flash,
	})
}
