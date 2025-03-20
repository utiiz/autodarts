package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/internal/handler"
	ws "server/internal/websocket"
	"server/views/pages"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()
	var cm = ws.NewConnectionManager()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/ws", func(e *core.RequestEvent) error {
			return handler.WS(e, cm)
		})

		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		se.Router.GET("/avatar/{path...}", apis.Static(os.DirFS("./pb_data/storage/_pb_users_auth_"), false))

		se.Router.GET("/", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.HomePage())
		})

		se.Router.GET("/dashboard", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.DashboardPage())
		})

		se.Router.GET("/game", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.GamePage())
		})

		se.Router.GET("/login", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.LoginPage())
		})

		se.Router.GET("/signup", func(e *core.RequestEvent) error {
			return handler.Render(e, http.StatusOK, pages.SignupPage())
		})

		test := se.Router.Group("/test")
		test.POST("/game", func(e *core.RequestEvent) error {
			body := struct {
				Dartboard string `json:"dartboard"`
			}{}
			if err := e.BindBody(&body); err != nil {
				return e.BadRequestError("Invalid request body", err)
			}
			log.Print(body.Dartboard)

			client, exists := cm.GetClient(body.Dartboard)
			if !exists {
				return e.BadRequestError("Client not found", nil)
			}

			event := ws.Message{Type: "GAME_START"}
			data, err := json.Marshal(event)
			if err != nil {
				log.Print(err)
			}

			cm.Send([]byte(data), client.UUID)

			return nil
		})

		return se.Next()
	})

	err := app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
