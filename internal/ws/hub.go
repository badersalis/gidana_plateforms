package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type client struct {
	userID uint
	conn   *websocket.Conn
	send   chan []byte
}

func (cl *client) writePump() {
	defer cl.conn.Close()
	for msg := range cl.send {
		if err := cl.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

// Hub routes real-time events to connected WebSocket clients.
type Hub struct {
	mu      sync.RWMutex
	clients map[uint]*client
}

var H = &Hub{clients: make(map[uint]*client)}

// Connect registers a new WebSocket connection for the given user.
// If the user already has a connection it is closed (last-write-wins per device).
func (h *Hub) Connect(userID uint, conn *websocket.Conn) *client {
	cl := &client{userID: userID, conn: conn, send: make(chan []byte, 64)}
	h.mu.Lock()
	if prev, ok := h.clients[userID]; ok {
		close(prev.send)
	}
	h.clients[userID] = cl
	h.mu.Unlock()
	go cl.writePump()
	return cl
}

// Disconnect removes the user's connection from the hub.
func (h *Hub) Disconnect(userID uint) {
	h.mu.Lock()
	delete(h.clients, userID)
	h.mu.Unlock()
}

// IsOnline reports whether the user has an active WebSocket connection.
func (h *Hub) IsOnline(userID uint) bool {
	h.mu.RLock()
	_, ok := h.clients[userID]
	h.mu.RUnlock()
	return ok
}

// Emit sends an event to the user if they are currently connected.
func (h *Hub) Emit(userID uint, event Event) {
	h.mu.RLock()
	cl, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	data, _ := json.Marshal(event)
	select {
	case cl.send <- data:
	default:
		log.Printf("ws: send buffer full for user %d, dropping event", userID)
	}
}
