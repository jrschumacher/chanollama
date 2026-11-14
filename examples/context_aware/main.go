package main

// TEST_METADATA:
// name: context_awareness
// description: Context-aware conversation with message history
// goal: Verify conversation context is maintained across multiple turns
// expected_output: Multi-turn conversation with context statistics
// timeout: 90s

import (
	"fmt"
	"log"
	"os"

	"github.com/jrschumacher/chanollama"
)

func main() {
	fmt.Println("🧠 Context-Aware Conversation Example")
	fmt.Println("=====================================")

	// Get model path from environment variable
	modelPath := os.Getenv("CHANOLLAMA_MODEL_PATH")
	if modelPath == "" {
		log.Fatal("CHANOLLAMA_MODEL_PATH environment variable is required")
	}

	// Create client
	client := chanollama.New(modelPath)
	defer func() { _ = client.Close() }()

	// Create a conversation with history management
	conversation := chanollama.NewConversation()

	// Set system prompt for the conversation
	conversation.AddSystemMessage(`You are a helpful programming tutor. Give brief, clear answers.`)

	fmt.Println("🎓 Programming Tutor Session")
	fmt.Print("Let's have a conversation that builds upon itself...\n\n")

	// Turn 1: Ask about variables
	fmt.Println("👤 Student: What are variables?")
	conversation.AddUserMessage("What are variables?")

	response, err := client.Chat(conversation.Messages)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}

	fmt.Printf("🎓 Tutor: %s\n\n", response.Content)
	conversation.AddAssistantMessage(response.Content)

	// Turn 2: Follow up about types (context-aware)
	fmt.Println("👤 Student: What about types?")
	conversation.AddUserMessage("What about types?")

	response, err = client.Chat(conversation.Messages)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}

	fmt.Printf("🎓 Tutor: %s\n\n", response.Content)
	conversation.AddAssistantMessage(response.Content)

	// Turn 3: Ask for a code example (should remember context)
	fmt.Println("👤 Student: Show me Go example?")
	conversation.AddUserMessage("Show me Go example?")

	response, err = client.Chat(conversation.Messages)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}

	fmt.Printf("🎓 Tutor: %s\n\n", response.Content)
	conversation.AddAssistantMessage(response.Content)

	// Turn 4: Reference previous example (context continuity)
	fmt.Println("👤 Student: How to use pointers?")
	conversation.AddUserMessage("How to use pointers?")

	response, err = client.Chat(conversation.Messages)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}

	fmt.Printf("🎓 Tutor: %s\n\n", response.Content)
	conversation.AddAssistantMessage(response.Content)

	// Show conversation statistics
	fmt.Printf("📊 Conversation Statistics:\n")
	fmt.Printf("   Messages: %d\n", len(conversation.Messages))
	fmt.Printf("   Estimated tokens: %d\n", conversation.EstimateTokens())
	fmt.Printf("   Max messages limit: %d\n", conversation.MaxMessages)
	fmt.Printf("   Token limit: %d\n", conversation.TokenLimit)

	// Example with streaming and context
	fmt.Println("\n🔄 Streaming with context:")
	fmt.Print("🎓 Tutor: ")

	conversation.AddUserMessage("Explain the difference?")

	response = client.ChatStream(conversation.Messages, func(token string) {
		fmt.Print(token)
	})

	if response.Error != nil {
		log.Printf("\n❌ Error: %v", response.Error)
	} else {
		conversation.AddAssistantMessage(response.Content)
	}

	fmt.Printf("\n\n📈 Final conversation stats:\n")
	fmt.Printf("   Total messages: %d\n", len(conversation.Messages))
	fmt.Printf("   Total estimated tokens: %d\n", conversation.EstimateTokens())

	// Demonstrate automatic trimming by adding many messages
	fmt.Println("\n🔄 Testing automatic message trimming...")

	for i := 0; i < 25; i++ {
		conversation.AddUserMessage(fmt.Sprintf("Test message %d", i))
		conversation.AddAssistantMessage(fmt.Sprintf("Response %d", i))
	}

	fmt.Printf("📊 After adding 50 test messages:\n")
	fmt.Printf("   Messages kept: %d (should be ≤ %d)\n", len(conversation.Messages), conversation.MaxMessages)
	fmt.Printf("   System messages preserved: ")

	systemCount := 0
	for _, msg := range conversation.Messages {
		if msg.Role == "system" {
			systemCount++
		}
	}
	fmt.Printf("%d\n", systemCount)

	fmt.Println("\n✅ Context-aware conversation example completed!")
}
