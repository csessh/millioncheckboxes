package main

import (
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		text := string(msg)
		log.Printf("Received: %s\n", text)

		if err := redis.Set("last_message", text, 10*time.Minute); err != nil {
			log.Println("Redis SET error:", err)
		}

		val, err := redis.Get("last_message")
		if err != nil {
			log.Println("Redis GET error:", err)
		}

		reply := fmt.Sprintf("Echo (stored in Redis): %s", val)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(reply)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	redis.InitRedis("localhost:6379", "", 0)

	http.HandleFunc("/ws", wsHandler)

	addr := ":8080"
	fmt.Println("WebSocket server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
