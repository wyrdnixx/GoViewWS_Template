package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func processMessage(conn *websocket.Conn, msg Message) {

	log.Printf("Processing Message from Client %s : %s", conn.RemoteAddr(), msg)

	switch msg.Type {
	case "textMessage":
		log.Printf("Text message received: %s", msg.Content)
		err := conn.WriteMessage(websocket.TextMessage, []byte("{\"text\":\"got your message\"}"))
		if err != nil {
			log.Println("Write:", err)
			break
		}

	case "echoTest":
		log.Printf("Echo request message received: %s", msg.Content)
		responseJSON, err := json.Marshal(msg)
		if err != nil {
			log.Println("Write:", err)
			break
		}
		err = conn.WriteMessage(websocket.TextMessage, responseJSON)
		if err != nil {
			log.Println("Write:", err)
			break
		}
	}
}
