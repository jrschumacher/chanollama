# 🚀 Getting Started with chanollama

A quick guide to get you up and running with chanollama for **truly offline** LLM inference.

## 📦 Installation

```bash
go get github.com/jrschumacher/chanollama
```

## 🔧 Prerequisites

### System Requirements

- **Go 1.24+** - Uses modern Go concurrency features
- **Sufficient RAM** - Model size + 2-4GB overhead
- **GGUF/GGML Model File** - No Ollama server required!

### Download a Model

Download GGUF models directly from sources like:

```bash
# Example: Download Llama 2 7B model (~4GB)
wget https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf

# Or download from other sources:
# - https://huggingface.co/models?filter=gguf  
# - https://huggingface.co/TheBloke (popular quantized models)
# - https://huggingface.co/microsoft (official Microsoft models)
```

**No Ollama Server Required** - chanollama loads models directly from GGUF files!

## 🎯 Quick Start

### Basic Offline Chat

```go
package main

import (
    "fmt"
    "log"
    "github.com/jrschumacher/chanollama"
)

func main() {
    // Create offline client - no server needed!
    client := chanollama.New("/path/to/llama-2-7b-chat.Q4_K_M.gguf")
    defer client.Close()
    
    // Simple chat
    messages := []chanollama.Message{
        chanollama.NewUserMessage("Hello! What can you help me with?"),
    }
    
    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Assistant:", response.Content)
}
```

### Real-time Streaming

```go
func main() {
    client := chanollama.New("/path/to/llama-2-7b-chat.Q4_K_M.gguf")
    defer client.Close()
    
    messages := []chanollama.Message{
        chanollama.NewUserMessage("Tell me a short story about Go programming"),
    }
    
    fmt.Print("Assistant: ")
    response := client.ChatStream(messages, func(token string) {
        fmt.Print(token) // Each token appears in real-time!
    })
    
    if response.Error != nil {
        log.Fatal(response.Error)
    }
    
    fmt.Println() // New line after response
}
```

## 📚 Core Concepts

### Messages

chanollama uses a simple message structure:

```go
type Message struct {
    Role      string    // "system", "user", "assistant"
    Content   string    // The message content
    Timestamp time.Time // When the message was created
    Metadata  map[string]string // Optional metadata
}

// Helper constructors
systemMsg := chanollama.NewSystemMessage("You are a helpful assistant")
userMsg := chanollama.NewUserMessage("Hello!")
assistantMsg := chanollama.NewAssistantMessage("Hi there!")
```

### Configuration

```go
// Custom configuration
config := chanollama.Config{
    Model:       "llama2:7b",
    BaseURL:     "http://localhost:11434", // Default
    Timeout:     30 * time.Second,         // Default
    Temperature: 0.7,                      // Default
    MaxTokens:   1000,                     // Default
    Headers: map[string]string{
        "Authorization": "Bearer your-token", // If needed
    },
}

client := chanollama.NewWithConfig(config)
```

### Error Handling

```go
response, err := client.Chat(messages)
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        fmt.Println("Request timed out")
    case errors.Is(err, chanollama.ErrModelNotFound):
        fmt.Println("Model not available")
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
    return
}

// Check response
if response.Error != nil {
    fmt.Printf("Response error: %v\n", response.Error)
    return
}
```

## 🎭 System Prompts

Guide the AI's behavior with system prompts:

```go
messages := []chanollama.Message{
    {
        Role: "system",
        Content: `You are a Go programming expert. Provide:
- Clear, concise explanations
- Working code examples
- Best practices and common pitfalls
- Performance considerations`,
    },
    {
        Role: "user",
        Content: "How do I handle JSON in Go?",
    },
}

response, err := client.Chat(messages)
```

## 💬 Conversation Management

### Basic Conversation

```go
var conversation []chanollama.Message

// Add system prompt
conversation = append(conversation, chanollama.NewSystemMessage(
    "You are a helpful coding assistant specialized in Go.",
))

// Interactive loop
scanner := bufio.NewScanner(os.Stdin)
for {
    fmt.Print("You: ")
    if !scanner.Scan() {
        break
    }
    
    userInput := scanner.Text()
    if userInput == "quit" {
        break
    }
    
    // Add user message
    conversation = append(conversation, chanollama.NewUserMessage(userInput))
    
    // Get response
    response, err := client.Chat(conversation)
    if err != nil {
        log.Printf("Error: %v", err)
        continue
    }
    
    // Add assistant response to conversation
    conversation = append(conversation, chanollama.NewAssistantMessage(response.Content))
    
    fmt.Printf("Assistant: %s\n\n", response.Content)
}
```

### Context Window Management

```go
// Keep conversation within token limits
func trimConversation(messages []chanollama.Message, maxTokens int) []chanollama.Message {
    if len(messages) == 0 {
        return messages
    }
    
    // Always keep system message if present
    var result []chanollama.Message
    if messages[0].IsSystem() {
        result = append(result, messages[0])
        messages = messages[1:]
    }
    
    // Count tokens and trim from the beginning
    totalTokens := 0
    for i := len(messages) - 1; i >= 0; i-- {
        tokens := chanollama.EstimateTokens(messages[i].Content)
        if totalTokens+tokens > maxTokens {
            break
        }
        
        totalTokens += tokens
        result = append([]chanollama.Message{messages[i]}, result...)
    }
    
    return result
}

// Usage
conversation = trimConversation(conversation, 2000) // Keep within 2K tokens
```

## 🔄 Streaming Examples

### Real-time UI Updates

```go
func streamToUI(query string, updateFunc func(string)) {
    client := chanollama.New("llama2:7b")
    defer client.Close()
    
    messages := []chanollama.Message{
        {Role: "user", Content: query},
    }
    
    var response strings.Builder
    client.ChatStream(messages, func(token string) {
        response.WriteString(token)
        updateFunc(response.String()) // Update UI with complete response so far
    })
}
```

### Concurrent Streaming

```go
func handleMultipleQueries(queries []string) {
    client := chanollama.New("llama2:7b")
    defer client.Close()
    
    var wg sync.WaitGroup
    
    for i, query := range queries {
        wg.Add(1)
        go func(id int, q string) {
            defer wg.Done()
            
            messages := []chanollama.Message{
                {Role: "user", Content: q},
            }
            
            fmt.Printf("Query %d: %s\n", id, q)
            fmt.Printf("Response %d: ", id)
            
            response := client.ChatStream(messages, func(token string) {
                fmt.Print(token)
            })
            
            if response.Error != nil {
                fmt.Printf("\nError %d: %v\n", id, response.Error)
            } else {
                fmt.Println() // New line
            }
        }(i, query)
    }
    
    wg.Wait()
}
```

## 🧪 Testing Your Setup

Create a simple test to verify everything works:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/jrschumacher/chanollama"
)

func main() {
    fmt.Println("🧪 Testing chanollama setup...")
    
    // Test 1: Client creation
    client := chanollama.New("llama2:7b")
    defer client.Close()
    fmt.Println("✅ Client created")
    
    // Test 2: Health check
    if !client.IsHealthy() {
        log.Fatal("❌ Ollama service not available")
    }
    fmt.Println("✅ Ollama service healthy")
    
    // Test 3: Simple chat
    start := time.Now()
    messages := []chanollama.Message{
        {Role: "user", Content: "Say 'Hello from chanollama!' and nothing else."},
    }
    
    response, err := client.Chat(messages)
    if err != nil {
        log.Fatalf("❌ Chat failed: %v", err)
    }
    
    duration := time.Since(start)
    fmt.Printf("✅ Chat successful in %v\n", duration)
    fmt.Printf("📝 Response: %s\n", response.Content)
    
    // Test 4: Streaming
    fmt.Print("🔄 Testing streaming: ")
    response = client.ChatStream(messages, func(token string) {
        fmt.Print(".")
    })
    
    if response.Error != nil {
        log.Fatalf("❌ Streaming failed: %v", response.Error)
    }
    
    fmt.Println("\n✅ All tests passed! 🎉")
}
```

## 🔗 Next Steps

Now that you have chanollama working:

1. **Learn about RAG**: [RAG Documentation](./rag.md) - Add document search to your LLM
2. **Build Agents**: [Agents Documentation](./agents.md) - Create specialized AI assistants
3. **Explore Examples**: [Examples](./examples/) - See real-world use cases
4. **Advanced Features**: [Configuration](./configuration.md) - Customize for your needs

## 🆘 Troubleshooting

### Common Issues

**"connection refused"**
```bash
# Check if Ollama is running
curl http://localhost:11434/api/version

# Start Ollama if needed
ollama serve
```

**"model not found"**
```bash
# List available models
ollama list

# Pull the model
ollama pull llama2:7b
```

**Slow responses**
```go
// Increase timeout
config := chanollama.Config{
    Timeout: 60 * time.Second, // Increase from 30s default
}
client := chanollama.NewWithConfig(config)
```

**Memory issues**
```go
// Use smaller model or reduce context
config := chanollama.Config{
    Model:     "llama2:7b",    // Instead of 13b or 70b
    MaxTokens: 500,            // Reduce from 1000 default
}
```

For more troubleshooting, see [Troubleshooting Guide](./troubleshooting.md).