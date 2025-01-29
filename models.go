package main

import "github.com/gorilla/websocket"

// Message represents the structure of a WebSocket message
type Message struct {
	Type    string `json:"type"`    // Type of the message (e.g., "text", "notification")
	Content string `json:"content"` // Actual message content
}

type WSConnections struct {
	C []*websocket.Conn
}
