# 🚀 Offline Model Inference Guide

A comprehensive guide to using chanollama for truly offline LLM inference without network dependencies.

## 🎯 Overview

chanollama is designed for **pure offline operation** - loading and running GGUF/GGML models directly in-process using Go's concurrency primitives. No Ollama server, no HTTP calls, no localhost communication required.

### Key Benefits

- **🔌 True Offline**: No network dependencies whatsoever
- **📦 Embeddable**: Ship applications with models included
- **⚡ Native Performance**: Direct model execution via llama.cpp bindings
- **🔄 Go-Native Streaming**: Real-time token delivery through channels
- **🛡️ Production Ready**: Built-in error recovery and resource management

---

## 📋 Prerequisites

### System Requirements

- **Go 1.24+** - Uses modern Go concurrency features
- **Sufficient RAM** - Model size + 2-4GB overhead (varies by model)
- **GGUF/GGML Model File** - Downloaded from Hugging Face or similar

### Model Files

Download GGUF models from sources like:
- [Hugging Face](https://huggingface.co/models?filter=gguf)
- [TheBloke's Models](https://huggingface.co/TheBloke)
- [Microsoft Models](https://huggingface.co/microsoft)

```bash
# Example: Download a 7B model (~4GB)
wget https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf
```

---

## 🚀 Quick Start

### Basic Offline Client

```go
package main

import (
    "fmt"
    "log"
    "github.com/jrschumacher/chanollama"
)

func main() {
    // Create offline client - model loads automatically on first use
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    // Simple conversation
    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain quantum computing in simple terms"),
    }

    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Assistant: %s\n", response.Content)
}
```

### Custom Configuration

```go
func main() {
    // Advanced offline configuration
    config := chanollama.Config{
        ModelPath:   "/path/to/model.gguf",
        OfflineMode: true,           // Explicitly enable offline mode
        Temperature: 0.7,            // Creativity level (0.0-1.0)
        MaxTokens:   2048,           // Max response length
        ContextSize: 4096,           // Context window size
        NumThreads:  8,              // CPU threads for inference
        UseGPU:      false,          // GPU acceleration (if supported)
        Verbose:     false,          // Suppress llama.cpp logs
    }

    client := chanollama.NewWithConfig("/path/to/model.gguf", config)
    defer client.Close()

    // Model initializes on first use
    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## 🔧 Model Management

### Model Loading

Models are loaded lazily on first inference call:

```go
client := chanollama.New("/path/to/model.gguf")

// Model not loaded yet - instant
fmt.Println("Client created")

// Model loads here - may take a few seconds
response, err := client.Chat(messages)
if err != nil {
    log.Fatal(err)
}
```

### Explicit Initialization

For better control over loading:

```go
client := chanollama.New("/path/to/model.gguf")
defer client.Close()

// Explicitly load model
if err := client.Initialize(); err != nil {
    log.Fatalf("Failed to load model: %v", err)
}

// Check if model is ready
if !client.IsHealthy() {
    log.Fatal("Model not ready for inference")
}

// Now inference is immediate
response, err := client.Chat(messages)
```

### Model Information

Get details about the loaded model:

```go
client := chanollama.New("/path/to/model.gguf")
defer client.Close()

// Initialize first
if err := client.Initialize(); err != nil {
    log.Fatal(err)
}

// Get model information
info, err := client.GetModelInfo()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Model Path: %s\n", info.ModelPath)
fmt.Printf("Context Size: %d\n", info.ContextSize)
fmt.Printf("Batch Size: %d\n", info.BatchSize)
```

---

## 💬 Conversation Patterns

### System Prompts

Guide model behavior with system prompts:

```go
func main() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    messages := []chanollama.Message{
        chanollama.NewSystemMessage(`You are a helpful coding assistant specialized in Go. 
- Provide clear, working code examples
- Explain best practices and common pitfalls  
- Focus on clean, idiomatic Go code`),
        chanollama.NewUserMessage("How do I handle HTTP requests in Go?"),
    }

    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Content)
}
```

### Context-Aware Conversations

Maintain conversation history:

```go
func interactiveChat() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    // Use conversation helper
    conv := chanollama.NewConversation()
    conv.AddSystemMessage("You are a helpful assistant.")

    scanner := bufio.NewScanner(os.Stdin)
    
    for {
        fmt.Print("You: ")
        if !scanner.Scan() {
            break
        }

        userInput := strings.TrimSpace(scanner.Text())
        if userInput == "quit" || userInput == "exit" {
            break
        }

        // Add user message
        conv.AddUserMessage(userInput)

        // Get response with full context
        response, err := client.Chat(conv.Messages)
        if err != nil {
            log.Printf("Error: %v", err)
            continue
        }

        // Add assistant response to conversation
        conv.AddResponse(response)

        fmt.Printf("Assistant: %s\n\n", response.Content)
    }
}
```

### Automatic Context Management

Prevent context overflow with automatic truncation:

```go
func managedConversation() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    // Configure automatic conversation management
    config := chanollama.ConversationConfig{
        MaxTokens:          3200,                            // 80% of context size
        MaxMessages:        50,                              // Fallback limit
        TruncationStrategy: chanollama.TruncationOldest,     // Remove oldest first
        PreserveSystem:     true,                            // Keep system messages
        AutoDetectLimits:   true,                            // Auto-detect from model
    }

    client.SetConversationConfig(config)

    // Now conversations are automatically managed
    messages := []chanollama.Message{
        chanollama.NewSystemMessage("You are a helpful assistant."),
        // ... many messages that would exceed limits ...
    }

    // Automatically truncated if needed
    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }

    // Check if truncation occurred  
    if response.Metadata["truncated_messages"] != "" {
        fmt.Printf("Truncated %s messages due to context limits\n", 
            response.Metadata["truncated_messages"])
    }
}
```

---

## 🔄 Streaming Inference

### Callback-Based Streaming

Real-time token delivery with callbacks:

```go
func streamingChat() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Write a short story about artificial intelligence"),
    }

    fmt.Print("Assistant: ")
    
    // Stream tokens in real-time
    response := client.ChatStream(messages, func(token string) {
        fmt.Print(token) // Each token appears immediately
    })

    if response.Error != nil {
        log.Fatal(response.Error)
    }

    fmt.Println() // New line after complete response
}
```

### Channel-Based Streaming

Go-native channel streaming for advanced control:

```go
func channelStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    // Create streaming engine
    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain Go channels and select statements"),
    }

    // Get channels for tokens and errors
    tokenCh, errorCh := engine.ChatStreamChan(messages)

    fmt.Print("Assistant: ")
    
    // Process streaming response with select
    for {
        select {
        case token, ok := <-tokenCh:
            if !ok {
                fmt.Println() // New line when complete
                return        // Channel closed, streaming complete
            }
            fmt.Print(token)

        case err := <-errorCh:
            if err != nil {
                log.Printf("\nError during streaming: %v", err)
                return
            }
        }
    }
}
```

### Streaming with Timeout

Add timeouts to streaming operations:

```go
func streamingWithTimeout() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Write a detailed technical explanation"),
    }

    // Set timeout for streaming operation
    timeout := 30 * time.Second
    tokenCh, errorCh := engine.ChatStreamChanWithTimeout(messages, timeout)

    fmt.Print("Assistant: ")
    
    for {
        select {
        case token, ok := <-tokenCh:
            if !ok {
                fmt.Println()
                return
            }
            fmt.Print(token)

        case err := <-errorCh:
            if err != nil {
                if errors.Is(err, context.DeadlineExceeded) {
                    fmt.Printf("\nStreaming timed out after %v\n", timeout)
                } else {
                    fmt.Printf("\nStreaming error: %v\n", err)
                }
                return
            }
        }
    }
}
```

---

## ⚡ Performance Optimization

### CPU Thread Configuration

Optimize for your hardware:

```go
func optimizedClient() *chanollama.Client {
    // Detect CPU cores
    numCPU := runtime.NumCPU()
    
    config := chanollama.Config{
        ModelPath:   "/path/to/model.gguf",
        OfflineMode: true,
        NumThreads:  numCPU,                    // Use all available cores
        ContextSize: 4096,                      // Adjust based on use case
        MaxTokens:   1024,                      // Balance speed vs length
        Temperature: 0.7,
    }

    return chanollama.NewWithConfig("/path/to/model.gguf", config)
}
```

### Memory Management

Monitor and optimize memory usage:

```go
func memoryOptimizedInference() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    // Check available memory before large operations
    monitor := chanollama.NewResourceMonitor()
    
    if !monitor.IsMemoryAvailable(4096) { // 4GB requirement
        log.Fatal("Insufficient memory for model inference")
    }

    // Use smaller context for memory-constrained environments
    config := chanollama.Config{
        ModelPath:   "/path/to/model.gguf", 
        ContextSize: 2048,     // Smaller context = less memory
        MaxTokens:   512,      // Shorter responses = faster inference
    }

    client.UpdateConfig(config)

    // Monitor memory during inference
    before, _ := monitor.GetMemoryUsage()
    
    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }

    after, _ := monitor.GetMemoryUsage()
    fmt.Printf("Memory used: %d MB\n", (after-before)/(1024*1024))
}
```

### Concurrent Processing

Handle multiple conversations safely:

```go
func concurrentInference() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    questions := []string{
        "What is machine learning?",
        "Explain neural networks",
        "How does natural language processing work?",
    }

    var wg sync.WaitGroup
    results := make(chan string, len(questions))

    // Process questions concurrently
    for i, question := range questions {
        wg.Add(1)
        go func(id int, q string) {
            defer wg.Done()

            messages := []chanollama.Message{
                chanollama.NewUserMessage(q),
            }

            response, err := client.Chat(messages)
            if err != nil {
                results <- fmt.Sprintf("Question %d failed: %v", id, err)
                return
            }

            results <- fmt.Sprintf("Question %d: %s", id, response.Content)
        }(i, question)
    }

    // Wait for all responses
    wg.Wait()
    close(results)

    // Print results
    for result := range results {
        fmt.Println(result)
        fmt.Println("---")
    }
}
```

---

## 🛡️ Error Handling & Recovery

### Robust Error Handling

chanollama includes built-in panic recovery:

```go
func robustInference() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain error handling in Go"),
    }

    // Built-in panic recovery converts panics to errors
    response, err := client.Chat(messages)
    if err != nil {
        // Handle specific error types
        switch {
        case errors.Is(err, chanollama.ErrModelNotLoaded):
            log.Fatal("Model failed to load - check model file path")
        case errors.Is(err, chanollama.ErrInferenceFailed):
            log.Printf("Inference failed, retrying: %v", err)
            // Implement retry logic
        case errors.Is(err, chanollama.ErrContextOverflow):
            log.Printf("Context too large, truncating: %v", err)
            // Handle context overflow
        default:
            log.Printf("Unexpected error: %v", err)
        }
        return
    }

    fmt.Println(response.Content)
}
```

### Custom Error Recovery

Add additional safety layers:

```go
func safeInference(client *chanollama.Client, messages []chanollama.Message) (*chanollama.Response, error) {
    // Additional panic recovery layer
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Additional panic recovery: %v", r)
        }
    }()

    // Check client health before inference
    if !client.IsHealthy() {
        return nil, fmt.Errorf("client not healthy, reinitializing")
    }

    // Perform inference with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // Note: context cancellation would need to be implemented in a wrapper
    response, err := client.Chat(messages)
    
    return response, err
}
```

### Automatic Recovery Strategies

Implement retry logic with backoff:

```go
func resilientInference(client *chanollama.Client, messages []chanollama.Message) (*chanollama.Response, error) {
    maxRetries := 3
    baseDelay := time.Second

    for attempt := 1; attempt <= maxRetries; attempt++ {
        response, err := client.Chat(messages)
        
        if err == nil {
            return response, nil
        }

        // Don't retry certain errors
        if errors.Is(err, chanollama.ErrModelNotLoaded) {
            return nil, err // Fatal error, don't retry
        }

        log.Printf("Attempt %d failed: %v", attempt, err)

        if attempt < maxRetries {
            delay := time.Duration(attempt) * baseDelay
            log.Printf("Retrying in %v...", delay)
            time.Sleep(delay)
        }
    }

    return nil, fmt.Errorf("all %d attempts failed", maxRetries)
}
```

---

## 🔧 Advanced Configuration

### Complete Configuration Example

```go
func advancedConfiguration() {
    config := chanollama.Config{
        // Model settings
        Model:       "/path/to/model.gguf",
        ModelPath:   "/path/to/model.gguf",  
        OfflineMode: true,

        // Generation parameters
        Temperature: 0.7,                      // Creativity (0.0-1.0)
        MaxTokens:   2048,                     // Max response length
        Timeout:     60 * time.Second,         // Request timeout

        // Hardware settings
        NumThreads:  runtime.NumCPU(),         // CPU threads
        ContextSize: 4096,                     // Context window
        UseGPU:      false,                    // GPU acceleration

        // Logging
        Verbose:     false,                    // Suppress verbose logs
    }

    client := chanollama.NewWithConfig(config.ModelPath, config)
    defer client.Close()

    // Conversation management
    convConfig := chanollama.ConversationConfig{
        MaxTokens:          3200,                            // 80% of context
        MaxMessages:        50,                              // Fallback limit
        TruncationStrategy: chanollama.TruncationOldest,     // Truncation method
        PreserveSystem:     true,                            // Keep system prompts
        AutoDetectLimits:   true,                            // Auto-detect from model
        EnableLogging:      false,                           // Log truncation events
    }

    client.SetConversationConfig(convConfig)

    // The client is now fully configured
    fmt.Println("Client configured successfully")
}
```

### Dynamic Configuration Updates

Update configuration at runtime:

```go
func dynamicConfiguration() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    // Initial configuration
    response1, _ := client.Chat(messages)
    fmt.Printf("Initial response length: %d tokens\n", response1.Tokens)

    // Update for longer responses
    newConfig := client.GetConfig()
    newConfig.MaxTokens = 4096
    newConfig.Temperature = 0.9  // More creative
    client.UpdateConfig(newConfig)

    // Note: Some changes require reinitialization
    response2, _ := client.Chat(messages)
    fmt.Printf("Updated response length: %d tokens\n", response2.Tokens)
}
```

---

## 🏭 Production Deployment

### Production-Ready Setup

```go
func productionClient(modelPath string) (*chanollama.Client, error) {
    // Validate model file exists
    if _, err := os.Stat(modelPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("model file not found: %s", modelPath)
    }

    // Production configuration
    config := chanollama.Config{
        ModelPath:   modelPath,
        OfflineMode: true,
        
        // Conservative settings for stability
        Temperature: 0.3,                     // More deterministic
        MaxTokens:   1024,                    // Reasonable length
        ContextSize: 4096,                    // Standard context
        NumThreads:  runtime.NumCPU() / 2,   // Leave resources for other processes
        Timeout:     30 * time.Second,       // Reasonable timeout
        
        // Production logging
        Verbose:     false,                   // Clean logs
    }

    client := chanollama.NewWithConfig(modelPath, config)

    // Pre-initialize for consistent performance
    if err := client.Initialize(); err != nil {
        return nil, fmt.Errorf("failed to initialize client: %w", err)
    }

    // Verify health
    if !client.IsHealthy() {
        client.Close()
        return nil, fmt.Errorf("client health check failed")
    }

    return client, nil
}
```

### Health Monitoring

Implement health checks for production:

```go
func healthCheck(client *chanollama.Client) error {
    // Basic health check
    if !client.IsHealthy() {
        return fmt.Errorf("client unhealthy")
    }

    // Test inference to verify functionality
    testMessages := []chanollama.Message{
        chanollama.NewUserMessage("Hello"),
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Simple test inference
    response, err := client.Chat(testMessages)
    if err != nil {
        return fmt.Errorf("inference test failed: %w", err)
    }

    if response.IsError() {
        return fmt.Errorf("inference returned error: %v", response.Error)
    }

    if len(response.Content) == 0 {
        return fmt.Errorf("inference returned empty response")
    }

    return nil
}
```

### Graceful Shutdown

Implement proper cleanup:

```go
func main() {
    client := chanollama.New("/path/to/model.gguf")
    
    // Setup graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigChan
        log.Println("Shutting down gracefully...")
        
        // Close client and cleanup resources
        if err := client.Close(); err != nil {
            log.Printf("Error during shutdown: %v", err)
        }
        
        os.Exit(0)
    }()

    // Your application logic here
    runApplication(client)
}
```

---

## 🎯 Best Practices

### 1. Model Selection

- **7B models**: Good balance of quality and speed, ~4GB RAM
- **13B models**: Better quality, slower inference, ~8GB RAM  
- **70B models**: Highest quality, very slow, ~40GB RAM
- **Quantized models** (Q4_K_M, Q5_K_M): Smaller size, minimal quality loss

### 2. Performance Optimization

- Use appropriate `NumThreads` for your hardware
- Set realistic `MaxTokens` to balance speed and completeness
- Enable `UseGPU` if you have compatible hardware
- Monitor memory usage and adjust `ContextSize` accordingly

### 3. Error Handling

- Always check `client.IsHealthy()` before critical operations
- Implement retry logic for transient failures
- Use timeouts for all inference operations
- Handle specific error types appropriately

### 4. Resource Management

- Always call `client.Close()` to free resources
- Use `defer client.Close()` immediately after creation
- Monitor memory usage in long-running applications
- Consider using connection pooling for high-concurrency scenarios

### 5. Configuration Management

- Use `DefaultOfflineConfig()` as a starting point
- Validate model files exist before creating clients
- Configure conversation limits to prevent context overflow
- Use appropriate logging levels for your environment

---

This guide covers the essential aspects of offline model inference with chanollama. For more advanced topics, see the [API Reference](./api-reference.md) and [examples directory](../examples/).