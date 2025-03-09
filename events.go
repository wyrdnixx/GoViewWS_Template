package main

import (
	"encoding/json"
	"fmt"
)

// processMessage handles the incoming message based on its "type"
func processMessage(msg []byte) ([]byte, error) {
	var message Message

	// Try to unmarshal the incoming message into the Message struct
	if err := json.Unmarshal(msg, &message); err != nil {
		return generateErrorJSON("Invalid JSON format")
	}

	// Switch on the "type" field to handle different types of messages
	switch message.Type {
	case "echoTestx":
		// Example case: Echo the content back
		response := map[string]string{
			"type":    "echoResponse",
			"content": message.Content,
		}
		return json.Marshal(response)

	// You can add more cases here for different message types
	default:
		// Unsupported type
		return generateErrorJSON(fmt.Sprintf("Unsupported message type: %s", message.Type))
	}
}

// generateErrorJSON creates an error JSON response
func generateErrorJSON(errorMessage string) ([]byte, error) {
	errorResponse := ErrorResponse{
		Type:    "error",
		Message: errorMessage,
	}
	return json.Marshal(errorResponse)
}
