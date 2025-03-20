package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	UUID       string
	Connection *websocket.Conn
	LastActive time.Time
	// Add more client-specific fields as needed
}

// ConnectionManager manages active WebSocket connections
type ConnectionManager struct {
	Clients map[string]*Client
	mutex   sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Clients: make(map[string]*Client),
	}
}

// Register adds a new client connection
func (cm *ConnectionManager) Register(uuid string, conn *websocket.Conn) *Client {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	client := &Client{
		UUID:       uuid,
		Connection: conn,
		LastActive: time.Now(),
	}
	cm.Clients[uuid] = client
	return client
}

// Unregister removes a client connection
func (cm *ConnectionManager) Unregister(uuid string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if client, exists := cm.Clients[uuid]; exists {
		client.Connection.Close()
		delete(cm.Clients, uuid)
	}
}

// GetClient retrieves a client by UUID
func (cm *ConnectionManager) GetClient(uuid string) (*Client, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	client, exists := cm.Clients[uuid]
	return client, exists
}

// BroadcastMessage sends a message to all connected clients
func (cm *ConnectionManager) Broadcast(message []byte) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	for _, client := range cm.Clients {
		client.Connection.WriteMessage(websocket.TextMessage, message)
	}
}

// ClientCount returns the number of connected clients
func (cm *ConnectionManager) ClientCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.Clients)
}

// GetAllClients returns a slice of all connected clients
func (cm *ConnectionManager) GetClients() []*Client {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	clients := make([]*Client, 0, len(cm.Clients))
	for _, client := range cm.Clients {
		clients = append(clients, client)
	}
	return clients
}

// ---------------------------- //
// ---- MESSAGES STRUCTS ----- //
// ---------------------------- //

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Data struct {
}

type Init struct {
	Data
	UUID string `json:"uuid"`
}
