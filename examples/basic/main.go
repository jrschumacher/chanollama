package main

// TEST_METADATA:
// name: basic_chat
// description: Basic chat functionality test
// goal: Verify client can initialize and get a non-empty response
// expected_output: Response: [non-empty string]
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

	// Create a new client with model file path
	client := chanollama.New(modelPath)
	defer func() { _ = client.Close() }()

	// Make a simple chat request
	response, err := client.Chat([]chanollama.Message{
		{Role: "user", Content: "What is the capital of France? Wait what?"},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response.Content)
}
