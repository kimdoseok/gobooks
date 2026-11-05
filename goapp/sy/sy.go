package sy

import (
	sv "goapp/sv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type (
	Menu struct {
		app.Compo
		AppService sv.AppService
	}
)

func (m *Menu) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Welcome to Rbooks"),
		app.P().Text("Yet another quickbooks like app"),
		app.A().Href("/login").Text("Login"),
	)
}
