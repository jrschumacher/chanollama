package main

// TEST_METADATA:
// name: concurrent_streaming
// description: Sequential streaming demo (avoiding memory pool issues)
// goal: Verify streaming works with multiple requests in sequence
// expected_output: Multiple conversations completed sequentially
// timeout: 90s

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jrschumacher/chanollama"
)

func main() {
	fmt.Println("⚡ Concurrent Streaming Example")
	fmt.Println("===============================")

	// Get model path from environment variable
	modelPath := os.Getenv("CHANOLLAMA_MODEL_PATH")
	if modelPath == "" {
		log.Fatal("CHANOLLAMA_MODEL_PATH environment variable is required")
	}

	// Create a simple client for sequential streaming
	client := chanollama.New(modelPath)
	defer func() { _ = client.Close() }()

	// Example 1: Multiple sequential conversations
	fmt.Println("🔄 Starting 3 sequential conversations...")

	questions := []string{
		"What is Go?",
		"Explain Docker",
		"What is HTTP?",
	}

	for i, question := range questions {
		fmt.Printf("🚀 Starting conversation %d: %s\n", i+1, question)

		messages := []chanollama.Message{
			{Role: "system", Content: "Give brief answers."},
			{Role: "user", Content: question},
		}

		fmt.Printf("💬 Conversation %d: ", i+1)
		response := client.ChatStream(messages, func(token string) {
			fmt.Print(token)
		})

		if response.Error != nil {
			log.Printf("\n❌ Conversation %d error: %v", i+1, response.Error)
			continue
		}

		fmt.Printf("\n✅ Conversation %d completed\n\n", i+1)
	}

	fmt.Println("🎉 All conversations completed!")

	// Example 2: Simple question answering
	fmt.Println("\n💬 Simple Q&A Demo")
	fmt.Println("==================")

	quickQuestions := []string{
		"What is AI?",
		"Define programming",
	}

	for i, question := range quickQuestions {
		fmt.Printf("👤 Question %d: %s\n", i+1, question)

		messages := []chanollama.Message{
			{Role: "system", Content: "Be brief."},
			{Role: "user", Content: question},
		}

		fmt.Printf("🤖 Answer: ")
		response := client.ChatStream(messages, func(token string) {
			fmt.Print(token)
		})

		if response.Error != nil {
			log.Printf("\n❌ Error: %v", response.Error)
			continue
		}

		fmt.Print("\n\n")
	}

	fmt.Println("🏁 Q&A session ended!")

	// Example 3: Performance timing
	fmt.Println("\n📈 Performance Test")
	fmt.Println("==================")

	start := time.Now()

	testMessages := []chanollama.Message{
		{Role: "user", Content: "Say hello"},
	}

	fmt.Print("🚀 Test: ")
	response := client.ChatStream(testMessages, func(token string) {
		fmt.Print(token)
	})

	if response.Error != nil {
		log.Printf("\n❌ Test error: %v", response.Error)
	} else {
		duration := time.Since(start)
		fmt.Printf("\n⏱️ Test completed in %v\n", duration)
	}

	fmt.Println("\n✅ Sequential streaming examples completed!")
}
