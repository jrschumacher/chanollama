# ⚙️ Configuration

Configuration options for chanollama clients, engines, and advanced features.

## 🔧 Client Configuration

### Basic Configuration

```go
config := chanollama.Config{
    Model:       "llama2:7b",               // Ollama model name
    BaseURL:     "http://localhost:11434",  // Ollama server URL
    Timeout:     30 * time.Second,          // Request timeout
    Temperature: 0.7,                       // Sampling temperature (0.0-1.0)
    MaxTokens:   1000,                      // Maximum tokens in response
    Headers: map[string]string{             // Custom HTTP headers
        "Authorization": "Bearer token",
        "X-Custom": "value",
    },
}

client := chanollama.NewWithConfig(config)
```

### Default Configuration

```go
// Get default config
config := chanollama.DefaultConfig("llama2:7b")

// Modify as needed
config.Temperature = 0.5
config.MaxTokens = 2048

client := chanollama.NewWithConfig(config)
```

## 🎛️ Engine Configuration

### Advanced Engine

```go
engineConfig := chanollama.EngineConfig{
    Model:        "llama2:7b",
    Temperature:  0.7,
    MaxTokens:    1000,
    SystemPrompt: "You are a helpful assistant.",
    
    // Channel configurations
    RequestChannelSize:   100,
    StreamingChannelSize: 1000,
    
    // Timeouts
    OperationTimeout: 30 * time.Second,
    ShutdownTimeout:  10 * time.Second,
    
    // Concurrency
    MaxConcurrentOps: 10,
}

engine := chanollama.NewAdvancedEngine(client, &engineConfig)
```

### Default Engine Config

```go
config := chanollama.DefaultEngineConfig()
// Defaults:
// - RequestChannelSize: 100
// - StreamingChannelSize: 1000
// - OperationTimeout: 30s
// - MaxConcurrentOps: 10
```

## 🔍 Stopwords Configuration

### Default Stopwords

```go
// Uses embedded 1,298 stopwords
keywords := chanollama.ExtractKeywords(text)
```

### Custom Stopwords

```go
// Add additional stopwords
config := chanollama.StopwordsConfig{
    UseDefault: true,
    Additional: []string{"api", "endpoint", "service"},
}

// Override completely
config = chanollama.StopwordsConfig{
    Override: []string{"the", "a", "an", "and", "or"},
}

// Disable stopwords
config = chanollama.StopwordsConfig{
    UseDefault: false,
}

keywords := chanollama.ExtractKeywordsWithConfig(text, config)
```

## 🤖 Agent Configuration

### Agent Manager

```go
agentManager := chanollama.NewAgentManager()

// Configure agent directory
agentManager = agentManager.WithAgentDirectory("./agents")

// Enable auto-refresh
agentManager = agentManager.WithAutoRefresh(5 * time.Minute)

// Add event handlers
agentManager = agentManager.WithEventHandler(&MyEventHandler{})
```

### Agent File Format

```yaml
---
name: "Development Expert"
description: "Expert in software development"
type: "development"
version: "1.0.0"
author: "Your Name"
tags: ["programming", "development", "coding"]
capabilities: ["code-generation", "debugging", "testing"]
custom:
  expertise_level: "expert"
  languages: "go,javascript,python"
---

Your system prompt here...
```

## 🏪 RAG Configuration

### Vector Store

```go
vectorStore := chanollama.NewVectorStore("./vector_index.json")

// Load existing index
err := vectorStore.LoadIndex()
if err != nil {
    log.Printf("No existing index: %v", err)
}

// Add documents with auto-embedding generation
docs := []chanollama.Document{
    {
        ID:      "doc1",
        Title:   "Go Programming",
        Content: "Go is a programming language...",
        Metadata: map[string]string{
            "category": "programming",
            "level":    "beginner",
        },
    },
}

err = vectorStore.AddDocuments(docs)
if err != nil {
    log.Fatal(err)
}

// Save index
err = vectorStore.SaveIndex()
if err != nil {
    log.Fatal(err)
}
```

### Simple RAG Store

```go
simpleRAG := chanollama.NewSimpleRAGStore("./simple_index.json")

// Configure with custom stopwords
stopwordsConfig := chanollama.StopwordsConfig{
    UseDefault: true,
    Additional: []string{"framework", "library", "module"},
}

// Add documents
doc := chanollama.Document{
    ID:      "guide1",
    Title:   "API Development Guide",
    Content: "How to build REST APIs...",
    Tags:    []string{"api", "rest", "development"},
}

err := simpleRAG.AddDocument(doc)
if err != nil {
    log.Fatal(err)
}
```

## 🌍 Environment Variables

chanollama respects these environment variables:

```bash
# Ollama server configuration
export OLLAMA_HOST="http://localhost:11434"
export OLLAMA_TIMEOUT="60s"

# Default model
export CHANOLLAMA_DEFAULT_MODEL="llama2:7b"

# Agent directory
export CHANOLLAMA_AGENT_DIR="./agents"

# RAG indices
export CHANOLLAMA_VECTOR_INDEX="./data/vector_index.json"
export CHANOLLAMA_SIMPLE_INDEX="./data/simple_index.json"
```

### Using Environment Variables

```go
func loadConfigFromEnv() chanollama.Config {
    config := chanollama.DefaultConfig("")
    
    if host := os.Getenv("OLLAMA_HOST"); host != "" {
        config.BaseURL = host
    }
    
    if model := os.Getenv("CHANOLLAMA_DEFAULT_MODEL"); model != "" {
        config.Model = model
    }
    
    if timeoutStr := os.Getenv("OLLAMA_TIMEOUT"); timeoutStr != "" {
        if timeout, err := time.ParseDuration(timeoutStr); err == nil {
            config.Timeout = timeout
        }
    }
    
    return config
}

client := chanollama.NewWithConfig(loadConfigFromEnv())
```

## 🔒 Security Configuration

### Authentication

```go
config := chanollama.Config{
    BaseURL: "https://ollama.company.com",
    Headers: map[string]string{
        "Authorization": "Bearer " + os.Getenv("OLLAMA_TOKEN"),
        "X-API-Key":     os.Getenv("API_KEY"),
    },
}
```

### TLS Configuration

```go
import (
    "crypto/tls"
    "net/http"
)

// Custom HTTP client with TLS config
tlsConfig := &tls.Config{
    InsecureSkipVerify: false, // Set to true for self-signed certs
    MinVersion:         tls.VersionTLS12,
}

transport := &http.Transport{
    TLSClientConfig: tlsConfig,
}

httpClient := &http.Client{
    Transport: transport,
    Timeout:   30 * time.Second,
}

// Use with chanollama (requires custom client configuration)
```

## 📊 Performance Configuration

### Memory Management

```go
// For large document collections
config := chanollama.Config{
    Timeout:   60 * time.Second, // Longer timeout for large operations
    MaxTokens: 500,              // Smaller responses to conserve memory
}

// Batch processing for large datasets
const batchSize = 100
for i := 0; i < len(documents); i += batchSize {
    batch := documents[i:min(i+batchSize, len(documents))]
    err := vectorStore.AddDocuments(batch)
    // Handle error and continue
}
```

### Concurrency Tuning

```go
engineConfig := chanollama.EngineConfig{
    MaxConcurrentOps:     5,    // Reduce for lower memory usage
    RequestChannelSize:   50,   // Smaller buffer for memory efficiency
    StreamingChannelSize: 500,  // Adjust based on streaming needs
}
```

### Token Optimization

```go
// Estimate tokens before sending
func optimizeMessage(message string, maxTokens int) string {
    tokens := chanollama.EstimateTokens(message)
    if tokens <= maxTokens {
        return message
    }
    
    // Truncate to fit within token limit
    words := strings.Fields(message)
    targetWords := len(words) * maxTokens / tokens
    
    if targetWords > len(words) {
        targetWords = len(words)
    }
    
    return strings.Join(words[:targetWords], " ")
}

optimizedContent := optimizeMessage(longContent, 1000)
```

## 🐛 Debug Configuration

### Logging

```go
import "log"

// Enable debug logging
config := chanollama.Config{
    Model:   "llama2:7b",
    Debug:   true,  // If supported
    Timeout: 60 * time.Second,
}

// Custom logging
type DebugClient struct {
    *chanollama.Client
}

func (d *DebugClient) Chat(messages []chanollama.Message) (*chanollama.Response, error) {
    log.Printf("Sending %d messages to model", len(messages))
    start := time.Now()
    
    response, err := d.Client.Chat(messages)
    
    log.Printf("Response received in %v", time.Since(start))
    if err != nil {
        log.Printf("Error: %v", err)
    } else {
        log.Printf("Response: %d tokens", response.Tokens)
    }
    
    return response, err
}
```

### Health Monitoring

```go
func monitorHealth(client *chanollama.Client) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if !client.IsHealthy() {
            log.Printf("WARNING: Ollama service is not healthy")
            // Implement retry logic or alerting
        }
    }
}

go monitorHealth(client)
```

## 📋 Configuration Validation

```go
func validateConfig(config chanollama.Config) error {
    if config.Model == "" {
        return fmt.Errorf("model is required")
    }
    
    if config.Timeout <= 0 {
        return fmt.Errorf("timeout must be positive")
    }
    
    if config.Temperature < 0 || config.Temperature > 1 {
        return fmt.Errorf("temperature must be between 0 and 1")
    }
    
    if config.MaxTokens <= 0 {
        return fmt.Errorf("max tokens must be positive")
    }
    
    // Validate URL
    if _, err := url.Parse(config.BaseURL); err != nil {
        return fmt.Errorf("invalid base URL: %v", err)
    }
    
    return nil
}

// Usage
config := chanollama.Config{...}
if err := validateConfig(config); err != nil {
    log.Fatal(err)
}
```

## 🏭 Production Configuration

### Recommended Production Settings

```go
prodConfig := chanollama.Config{
    Model:       "llama2:7b",
    BaseURL:     os.Getenv("OLLAMA_URL"),
    Timeout:     60 * time.Second,     // Longer timeout for production
    Temperature: 0.3,                 // Lower temperature for consistency
    MaxTokens:   1000,
    Headers: map[string]string{
        "Authorization": "Bearer " + os.Getenv("OLLAMA_TOKEN"),
        "User-Agent":    "chanollama/1.0.0",
    },
}

// Agent configuration
agentManager := chanollama.NewAgentManager().
    WithAgentDirectory("/app/agents").
    WithAutoRefresh(10 * time.Minute). // Auto-refresh for config updates
    WithEventHandler(&ProductionEventHandler{})

// RAG configuration
vectorStore := chanollama.NewVectorStore("/data/vector_index.json")
simpleRAG := chanollama.NewSimpleRAGStore("/data/simple_index.json")
```

### Health Checks

```go
func healthCheck(client *chanollama.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !client.IsHealthy() {
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
            return
        }
        
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "healthy",
            "timestamp": time.Now().ISO8601(),
        })
    }
}
```

For more advanced configuration examples, see the [examples directory](../examples/).