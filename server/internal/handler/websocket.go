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

func WS(e *core.RequestEvent, cm *ws.ConnectionManager) error {
	conn, err := upgrader.Upgrade(e.Response, e.Request, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

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
	}
	return nil
}
