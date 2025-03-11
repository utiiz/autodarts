package main

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("super-secret-key")) // Replace with a secure key
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Game struct {
	Players []Player
}

type Player struct {
	Name    string
	Score   int
	Darts   [3]Dart
	History []Dart
}

type Dart struct {
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Score Score   `json:"score"`
}

type Score struct {
	Bed     string `json:"bed"`
	Segment int    `json:"segment"`
	Score   int    `json:"score"`
}

func (s Score) String() string {
	return s.Bed + strconv.Itoa(s.Segment)
}

func init() {
	gob.Register(Game{}) // Needed to store Game struct in session
}

type Templates struct {
	Templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *Templates {
	return &Templates{
		Templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	templates := newTemplates()
	e.Renderer = templates

	e.GET("/", indexHandler)
	e.GET("/ws", func(c echo.Context) error {
		return websocketHandler(c, templates)
	})

	log.Println("Server running on :8080")
	e.Start(":8080")
}

func indexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func websocketHandler(c echo.Context, renderer *Templates) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Simple echo websocket server
	for {
		// Read message from browser
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Print the message to the console
		log.Printf("Received: %s", msg)

		// Write message back to browser
		if err = ws.WriteMessage(msgType, msg); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
		var dart Dart
		err = json.Unmarshal(msg, &dart)
		if err != nil {
			log.Println("Invalid dart data:", err)
			continue
		}

		htmlContent, err := renderer.Render(http.StatusOK, "partial.html", data)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			continue
		}

		if err = ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
	return nil
}
