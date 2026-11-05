package sy

import (
	sv "goapp/sv"
	"net/http"
	"net/url"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type (
	Login struct {
		app.Compo
		AppService sv.AppService

		UserID string
		Passwd string
	}
)

func (l *Login) OnLogin(ctx app.Context, e app.Event) {
	//l.name = ctx.JSSrc().Get("value").String()
	app.Log("UserID: ", l.UserID, "Password: ", l.Passwd)

	go func() {
		data := url.Values{}
		data.Set("yus_code", l.UserID)
		data.Set("yus_passwd", l.Passwd)

		resp, err := http.PostForm(
			"http://localhost:8080/api/v1/users/login",
			data,
		)
		if err != nil {
			app.Log(err)
			return
		}
		defer resp.Body.Close()
		app.Log("login successful %+v", resp)

		// handle response
	}()

}

func (l *Login) Render() app.UI {
	return app.Section().Class("hero is-fullheight").Body(
		app.Div().Class("hero-body is-justify-content-center is-align-items-center").Body(
			app.Div().Class("columns is-centered").Body(
				app.Div().Class("column is-4").Body(
					app.Div().Class("box").Style("width", "400px").Body(
						app.H1().Class("title has-text-centered").Text("Login"),

						app.Div().Class("field").Body(
							app.Label().Class("label").Text("UserID"),
							app.Div().Class("control").Body(
								app.Input().Class("input").Type("text").Value(l.UserID).Placeholder("Enter your UserID").OnChange(func(ctx app.Context, e app.Event) {
									l.UserID = strings.ToUpper(ctx.JSSrc().Get("value").String())
								}),
							),
						),

						app.Div().Class("field").Body(
							app.Label().Class("label").Text("Password"),
							app.Div().Class("control").Body(
								app.Input().Class("input").Type("password").Value(l.Passwd).Placeholder("Enter your password").OnChange(func(ctx app.Context, e app.Event) {
									l.Passwd = ctx.JSSrc().Get("value").String()
								}),
							),
						),
						/*
							app.Div().Class("field").Body(
								app.Label().Class("checkbox").Body(
									app.Input().Type("checkbox"),
									app.Text(" Remember me"),
								),
							),
						*/
						app.Div().Class("field").Body(
							app.Button().Class("button is-primary is-fullwidth").Text("Sign In").OnClick(l.OnLogin),
						),
						/*
							app.P().Class("has-text-centered").Body(
								app.A().Href("/forgot-password").Text("Forgot Password?"),
							),
						*/
					),
				),
			),
		),
	)
}
