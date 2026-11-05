package cm

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type (
	Dashboard struct {
		app.Compo
	}
)

func (d *Dashboard) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Welcome to Rbooks"),
		app.P().Text("Yet another quickbooks like app"),
		app.A().Href("/login").Text("Login"),
	)
}
