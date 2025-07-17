package main

// TEST_METADATA:
// name: streaming_chat
// description: Streaming chat functionality test
// goal: Verify client can stream responses via callback
// expected_output: Streaming response: [tokens printed in real-time]
// timeout: 30s

import (
	"fmt"
	"log"
	"os"

	"github.com/jrschumacher/chanollama"
)

func main() {
	// Get model path from environment variable
	modelPath := os.Getenv("CHANOLLAMA_MODEL_PATH")
	if modelPath == "" {
		log.Fatal("CHANOLLAMA_MODEL_PATH environment variable is required")
	}

	// Create a new client with streaming enabled
	client := chanollama.New(modelPath)
	defer client.Close()

	// Make a streaming chat request
	fmt.Println("Streaming response:")
	response := client.ChatStream([]chanollama.Message{
		{Role: "user", Content: "Tell me a short story about a robot."},
	}, func(token string) {
		fmt.Print(token)
	})

	if response.Error != nil {
		log.Printf("Error: %v", response.Error)
	}
	fmt.Println()
}