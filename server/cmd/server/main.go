package main

import (
	"log"
	"net/http"
	"os"
	"server/internal/handler"
	"server/views/pages"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		se.Router.GET("/avatar/{path...}", apis.Static(os.DirFS("./pb_data/storage/_pb_users_auth_"), false))

		se.Router.GET("/", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.HomePage())
		})

		se.Router.GET("/dashboard", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.DashboardPage())
		})

		se.Router.GET("/login", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.LoginPage())
		})

		se.Router.GET("/signup", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.SignupPage())
		})

		return se.Next()
	})

	err := app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
