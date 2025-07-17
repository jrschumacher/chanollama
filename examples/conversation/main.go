package main

// TEST_METADATA:
// name: conversation_history
// description: Conversation with persistent message history
// goal: Verify engine maintains conversation context across multiple messages
// expected_output: Multi-turn conversation about Japan travel
// timeout: 60s

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

	// Create a new client for managing chat history
	client := chanollama.New(modelPath)
	defer client.Close()

	// Start a conversation with message history
	messages := []chanollama.Message{}
	
	// First message
	messages = append(messages, chanollama.Message{Role: "user", Content: "Hello, I'm planning a trip to Japan."})
	response1, err := client.Chat(messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Assistant: %s\n\n", response1.Content)
	messages = append(messages, chanollama.Message{Role: "assistant", Content: response1.Content})

	// Follow-up message (engine remembers context)
	messages = append(messages, chanollama.Message{Role: "user", Content: "What's the best season to visit?"})
	response2, err := client.Chat(messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Assistant: %s\n\n", response2.Content)
	messages = append(messages, chanollama.Message{Role: "assistant", Content: response2.Content})

	// Another follow-up
	messages = append(messages, chanollama.Message{Role: "user", Content: "What about food recommendations?"})
	response3, err := client.Chat(messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Assistant: %s\n", response3.Content)
}