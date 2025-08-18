package main

import (
	"encoding/json"
	"log"
	"net/http"

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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	defer conn.Close()
	log.Println("Client connected")

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
			}
		}
	}
}

func main() {
	redis.InitRedis("localhost:6379", "", 0)
	http.HandleFunc("/ws", wsHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
