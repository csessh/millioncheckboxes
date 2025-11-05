package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/csessh/1M-backend/internal/protocol"
	"github.com/csessh/1M-backend/internal/redis"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Cmd   string `json:"cmd"`
	Index int    `json:"index,omitempty"`
	Value string `json:"value,omitempty"`
}

// Connection registry
var (
	connections = make(map[*websocket.Conn]bool)
	connMutex   sync.RWMutex
)

// Register a new connection
func registerConnection(conn *websocket.Conn) {
	connMutex.Lock()
	defer connMutex.Unlock()
	connections[conn] = true
	log.Printf("Client registered. Total connections: %d", len(connections))
}

// Unregister a connection
func unregisterConnection(conn *websocket.Conn) {
	connMutex.Lock()
	defer connMutex.Unlock()
	delete(connections, conn)
	log.Printf("Client unregistered. Total connections: %d", len(connections))
}

// Broadcast message to all connected clients except sender
func broadcast(sender *websocket.Conn, msg Message) {
	connMutex.RLock()
	defer connMutex.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	for conn := range connections {
		if conn == sender {
			continue
		}

		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error broadcasting to client: %v", err)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	defer conn.Close()
	defer unregisterConnection(conn)

	registerConnection(conn)
	log.Println("Client connected")

	// Send initial state to the new client
	for i := 0; i < 100; i++ {
		isChecked, err := redis.GetCheckbox(i)
		if err == nil && isChecked {
			initMsg := Message{
				Cmd:   "SET",
				Index: i,
				Value: "true",
			}
			data, _ := json.Marshal(initMsg)
			conn.WriteMessage(websocket.TextMessage, data)
		}
	}

	for {
		messageType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		if messageType == websocket.BinaryMessage {
			binaryMsg, err := protocol.DecodeBinaryMessage(msgBytes)
			if err != nil {
				log.Printf("Binary decode error: %v", err)
				continue
			}

			if binaryMsg.Command == protocol.CmdSet {
				log.Printf("Binary: Updating index %d -> %v", binaryMsg.CheckboxIndex, binaryMsg.IsChecked)
				redis.SetCheckbox(int(binaryMsg.CheckboxIndex), binaryMsg.IsChecked)

				// Broadcast to other clients
				broadcastMsg := Message{
					Cmd:   "SET",
					Index: int(binaryMsg.CheckboxIndex),
					Value: func() string {
						if binaryMsg.IsChecked {
							return "true"
						}
						return "false"
					}(),
				}
				broadcast(conn, broadcastMsg)
			}
		} else if messageType == websocket.TextMessage {
			log.Printf("JSON message received: %q", string(msgBytes))

			var msg Message
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Printf("JSON parse error: %v - Raw data: %q", err, string(msgBytes))

				errorMsg := Message{Cmd: "ERROR", Value: "Invalid JSON format"}
				errorData, _ := json.Marshal(errorMsg)
				conn.WriteMessage(websocket.TextMessage, errorData)
				continue
			}

			if msg.Cmd == "SET" {
				value := msg.Value == "true"
				log.Printf("JSON: Updating index %d -> %v", msg.Index, value)
				redis.SetCheckbox(msg.Index, value)

				// Broadcast to other clients
				broadcast(conn, msg)
			}
		}
	}
}

func main() {
	redis.InitRedis("localhost:6379", "", 0)

	// Serve static files from public directory
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// WebSocket endpoint
	http.HandleFunc("/ws", wsHandler)

	log.Println("Server starting on :8080")
	log.Println("Serving static files from ../public")
	log.Println("WebSocket endpoint: ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
