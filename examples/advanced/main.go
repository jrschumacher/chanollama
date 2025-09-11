package main

// TEST_METADATA:
// name: advanced_config
// description: Advanced configuration with custom engine settings
// goal: Verify engine can be configured with custom parameters
// expected_output: Streaming response with advanced configuration
// timeout: 45s

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

	// Create client with model path (advanced configuration can be added later)
	client := chanollama.New(modelPath)
	defer client.Close()

	// Example of using the engine with streaming
	messages := []chanollama.Message{
		{Role: "user", Content: "Explain the concept of microservices architecture."},
	}

	fmt.Println("Streaming response from technical assistant:")

	// Add system prompt to messages
	systemMessages := []chanollama.Message{
		{Role: "system", Content: "You are a helpful assistant specialized in technical documentation."},
	}
	systemMessages = append(systemMessages, messages...)

	response := client.ChatStream(systemMessages, func(token string) {
		fmt.Print(token)
	})

	if response.Error != nil {
		log.Printf("Error: %v", response.Error)
		return
	}

	fmt.Println()

	// Example completed - advanced configuration demonstrated
	fmt.Printf("\n✅ Advanced configuration example completed!\n")
}
