package handler

import (
	"encoding/json"
	"log"
	"net/http"

	ws "server/internal/websocket"

	"github.com/gorilla/websocket"
	"github.com/pocketbase/pocketbase/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WS(e *core.RequestEvent) error {
	conn, err := upgrader.Upgrade(e.Response, e.Request, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	var cm = ws.NewConnectionManager()

	for {
		var message ws.Message
		if err := conn.ReadJSON(&message); err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("Received: %s", message)

		switch message.Type {
		case "INIT":
			rawData, _ := json.Marshal(message.Data)
			var data ws.Init
			if err := json.Unmarshal(rawData, &data); err != nil {
				log.Println("read error:", err)
				break
			}
			log.Printf("Received: %s", data.UUID)
			cm.Register(data.UUID, conn)
		}

		cm.Broadcast([]byte("Hello, World!"))
	}
	return nil
}
