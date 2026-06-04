package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	username string
	socket   *websocket.Conn
	send     chan []byte
}

type Message struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	SenderID  string `json:"senderId"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	UserCount int    `json:"userCount"`
}

type ClientManager struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var manager = &ClientManager{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte, 256),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (m *ClientManager) start() {
	for {
		select {

		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			count := len(m.clients)
			m.mu.Unlock()

			log.Printf("✅ CONNECT id=%s user=%s total=%d", client.id, client.username, count)
			m.sendSystemMessage(fmt.Sprintf("%s joined", client.username), count)

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
			count := len(m.clients)
			m.mu.Unlock()

			log.Printf("❌ DISCONNECT id=%s user=%s total=%d", client.id, client.username, count)
			m.sendSystemMessage(fmt.Sprintf("%s left", client.username), count)

		case data := <-m.broadcast:
			m.mu.RLock()
			for client := range m.clients {
				select {
				case client.send <- data:
				default:
					go func(c *Client) { m.unregister <- c }(client)
				}
			}
			m.mu.RUnlock()
		}
	}
}

func (m *ClientManager) sendSystemMessage(content string, count int) {
	msg := Message{
		Type:      "system",
		Sender:    "SERVER",
		SenderID:  "server",
		Content:   content,
		Timestamp: time.Now().Format("15:04:05"),
		UserCount: count,
	}

	data, _ := json.Marshal(msg)
	m.broadcast <- data
}

func (m *ClientManager) userCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

func (c *Client) readPump() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, raw, err := c.socket.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var incoming Message

		if jsonErr := json.Unmarshal(raw, &incoming); jsonErr != nil {
			log.Printf("bad JSON from %s: %v", c.id, jsonErr)
			continue
		}

		if incoming.Content == "" {
			continue
		}

		outgoing := Message{
			Type:      incoming.Type,
			Sender:    c.username,
			SenderID:  c.id,
			Content:   incoming.Content,
			Timestamp: time.Now().Format("15:04:05"),
			UserCount: manager.userCount(),
		}

		data, _ := json.Marshal(outgoing)
		manager.broadcast <- data
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)

	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()

	for {
		select {

		case msg, ok := <-c.send:
			c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous_" + uuid.New().String()[:4]
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		id:       uuid.New().String(),
		username: username,
		socket:   conn,
		send:     make(chan []byte, 256),
	}

	manager.register <- client

	go client.writePump()
	go client.readPump()
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./public")).ServeHTTP(w, r)
}
func main() {
	go manager.start()

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", staticHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	log.Printf("🚀 NeonChat server running → http://localhost%s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
