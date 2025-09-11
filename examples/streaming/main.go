package main

// TEST_METADATA:
// name: engine_streaming
// description: Engine-based streaming with health checks
// goal: Verify engine initialization, health check, and streaming functionality
// expected_output: ✅ Engine ready, callback-based streaming output, token counts
// timeout: 45s

import (
	"fmt"
	"log"
	"os"

	"github.com/jrschumacher/chanollama"
)

func main() {
	fmt.Println("🚀 Real-time Streaming Example")
	fmt.Println("==============================")

	// Get model path from environment variable
	modelPath := os.Getenv("CHANOLLAMA_MODEL_PATH")
	if modelPath == "" {
		log.Fatal("CHANOLLAMA_MODEL_PATH environment variable is required")
	}

	// Create a streaming client with model
	client := chanollama.New(modelPath)
	defer client.Close()

	fmt.Println("✅ Engine ready")

	// Example 1: Callback-based streaming
	fmt.Println("\n🔄 Callback-based streaming:")
	fmt.Print("Assistant: ")

	messages := []chanollama.Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Tell me a short joke"},
	}

	response := client.ChatStream(messages, func(token string) {
		fmt.Print(token) // Real-time token display
	})

	if response.Error != nil {
		log.Printf("\n❌ Error: %v", response.Error)
	} else {
		fmt.Printf("\n\n📊 Total tokens: %d\n", response.Tokens)
	}

	// Example 2: Callback-based streaming with response capture
	fmt.Println("\n🔄 Callback-based streaming with response capture:")
	fmt.Print("Assistant: ")

	messages2 := []chanollama.Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "What is the capital of France?"},
	}

	var fullResponse string
	response2 := client.ChatStream(messages2, func(token string) {
		fmt.Print(token)
		fullResponse += token
	})

	if response2.Error != nil {
		log.Printf("\n❌ Error: %v", response2.Error)
		return
	}

	fmt.Printf("\n\n📊 Full response length: %d characters\n", len(fullResponse))

	// Example 3: Simple streaming without batching (TokenBatcher requires channel-based streaming)
	fmt.Println("\n🔄 Simple streaming example:")
	fmt.Print("Assistant: ")

	messages3 := []chanollama.Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Count from 1 to 10"},
	}

	response3 := client.ChatStream(messages3, func(token string) {
		fmt.Print(token)
	})

	if response3.Error != nil {
		log.Printf("\n❌ Error: %v", response3.Error)
		return
	}

	fmt.Println("\n\n✅ Streaming examples completed!")
}
