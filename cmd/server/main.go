package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected")

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("cb-%d", i)
		val, _ := redis.Get(key)
		if val == "true" {
			msg := Message{Cmd: "SET", Key: key, Value: "true"}
			data, _ := json.Marshal(msg)
			conn.WriteMessage(websocket.TextMessage, data)
		}
	}

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}

		if msg.Cmd == "SET" {
			log.Printf("Updating %s -> %s\n", msg.Key, msg.Value)
			redis.Set(msg.Key, msg.Value, 24*time.Hour) // store for a day
		}
	}
}

func main() {
	// Initialize Redis connection
	redis.InitRedis("localhost:6379", "", 0)

	// Set up WebSocket handler
	http.HandleFunc("/ws", wsHandler)

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
