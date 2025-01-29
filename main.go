package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsconnections WSConnections

// Serve the homepage (this serves the HTML page)
func serveHome(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "frontend/dist/index.html")

	// Serve the static folder at the root URL ("/")
	fs := http.FileServer(http.Dir("./frontend/dist/"))

	// Handle the root URL ("/") to serve files from the "./static" directory
	http.Handle("/", fs)
	/*
		err := http.FileServer(http.Dir("frontend/dist/"))
		if err != nil {
			log.Fatalf("Error serving frontend: %s", err)
		} */
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP to WebSocket
	if err != nil {
		log.Println("Upgrade:", err)
		return
	} else {
		wsconnections.C = append(wsconnections.C, conn)
		log.Printf("Active Connections: %s", wsconnections.C[0].RemoteAddr())
	}
	defer conn.Close()

	for {
		// Read JSON message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read:", err)
			break
		}

		// Unmarshal the received message (JSON) into a struct
		var receivedMessage Message
		err = json.Unmarshal(msg, &receivedMessage)
		// Errorhandling JSON processing
		if err != nil {
			log.Println("Unmarshal:", err)
			//break
			// Write the JSON response back to the client
			responseMessage := Message{
				Type:    "err",
				Content: "failure on JSON Unmarshal client message",
			}
			// Marshal the response message to JSON
			responseJSON, err := json.Marshal(responseMessage)
			if err != nil {
				log.Println("Marshal:", err)
				break
			}
			err = conn.WriteMessage(websocket.TextMessage, responseJSON)
			if err != nil {
				log.Println("Write:", err)
				break
			}
		}

		processMessage(conn, receivedMessage)

	}
}

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if not specified
	}

	http.HandleFunc("/", serveHome)         // Serve the index page
	http.HandleFunc("/ws", handleWebSocket) // WebSocket endpoint

	// Serve static files (e.g., JavaScript)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("Server started on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
