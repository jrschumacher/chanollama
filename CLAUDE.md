# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

chanollama is a Go library for **truly offline LLM (Large Language Model) inference** using Go channels and goroutines. The project eliminates network dependencies by loading and running GGUF/GGML models directly in-process, enabling embedded AI applications that can be distributed with models included.

## Essential Commands

### Development
- `go mod tidy` - Manage dependencies
- `go build` - Build the library
- `go test` - Run unit tests
- `go test -v` - Run tests with verbose output
- `go run examples/*/main.go` - Run specific examples

### Testing & Examples
- `go test ./...` - Run all tests
- `go run examples/basic/main.go` - Simple chat example
- `go run examples/streaming/main.go` - Real-time streaming demo
- `go run examples/06_concurrent_streaming/main.go` - Concurrent operations

### Documentation
- `go doc github.com/jrschumacher/chanollama` - API reference
- `godoc -http=:6060` - Local documentation server

## Architecture Overview

### Core Vision
chanollama is designed for **pure offline operation** with zero network dependencies:

- **Offline-Only Architecture**: Direct GGUF/GGML model loading with llama.cpp bindings
- **No Ollama Server Dependency**: Eliminates all external service requirements
- **Key Innovation**: Channel-based streaming using Go concurrency primitives
- **Embeddable Distribution**: Models can be included directly in application binaries

### Model Loading
```go
// Direct GGUF model file loading (preferred approach)
engine, err := chanollama.NewEngineWithModel("/path/to/model.gguf")

// Model path configuration
config := chanollama.Config{
    ModelPath:   "/path/to/model.gguf",
    OfflineMode: true, // Always true for robust operation
}
```

### Core Components

#### Primary Files
- **`types.go`** - Core data structures (Message, Response, Config, Conversation)
- **`offline_client.go`** - Direct model execution using llama.cpp bindings
- **`engine.go`** - Advanced engine with message-passing architecture
- **`factory.go`** - Model loading and initialization logic
- **`streaming.go`** - Channel-based streaming utilities for real-time token delivery
- **`client.go`** - Legacy HTTP client (deprecated - avoid using)

#### Advanced Features
- **`rag.go`** - Retrieval-Augmented Generation with vector and keyword search
- **`agents.go`** - Multi-agent framework with intelligent routing
- **`context.go`** - Conversation management and history trimming
- **`evaluation.go`** - Model evaluation and testing tools
- **`formatters.go`** - Output formatting utilities

### Key Patterns

#### Channel-Based Streaming
Real-time token delivery through Go channels:
```go
tokenCh, errorCh := engine.ChatStreamChan(messages)
for token := range tokenCh {
    fmt.Print(token) // Real-time streaming
}
```

#### Interface-Based Design
`Engine` interface provides consistent API for offline model operations:
```go
type Engine interface {
    Chat(messages []Message) (*Response, error)
    ChatStreamChan(messages []Message) (<-chan string, <-chan error)
    IsHealthy() bool
    Close() error
}
```

#### Message-Passing Architecture
Engine operations use structured message passing for concurrent request handling:
```go
type Operation struct {
    Type    OperationType
    Request interface{}
    Result  chan OperationResult
}
```

### Dependencies

**Core Dependencies**:
- `github.com/google/uuid v1.6.0` - UUID generation
- `github.com/ollama/ollama v0.9.6` - llama.cpp bindings for offline mode

**Go Version**: 1.24.0+ (uses modern Go concurrency features)

## Configuration Patterns

### Offline Client Configuration
```go
config := chanollama.Config{
    ModelPath:   "/path/to/model.gguf",    // Direct GGUF model file
    OfflineMode: true,                     // Always true for robust operation
    Temperature: 0.7,
    MaxTokens:   2048,
    UseGPU:      false,                    // Hardware acceleration option
    NumThreads:  4,                        // CPU thread control
    ContextSize: 4096,                     // Context window size
}
```

### Engine Configuration
```go
engineConfig := chanollama.EngineConfig{
    ModelPath:    "/path/to/model.gguf",   // Direct model file path
    SystemPrompt: "You are a helpful assistant.",
    Temperature:  0.7,
    MaxTokens:    1000,
}
```

## Testing Architecture

### Test Structure
- **`chanollama_test.go`** - Core functionality tests (client, engine, configuration)
- **`mock_engine_test.go`** - Mock implementations for testing without real models
- **Examples** - 10+ working examples serving as integration tests

### Test Coverage Areas
- Client initialization and configuration validation
- Engine creation and system prompt handling
- Default configuration verification
- Basic API surface testing

### Running Tests
Tests run completely offline without any external dependencies or services:
```bash
go test -v                    # All tests (fully offline)
go test -run TestNewClient    # Specific test
go test ./...                 # All packages
```

## Integration Context

This library is being developed for integration with **OpenTDF's otdfctl CLI tool** (based on project path), likely to provide:
- Local LLM capabilities without external dependencies
- Embeddable models for secure offline operation
- AI-powered features aligned with OpenTDF's data protection mission

## Development Strategy & Architecture Goals

### Development Focus
**All new development should target offline-only operation:**

1. **Direct Model Loading** 🔄 In Progress - GGUF/GGML file loading via llama.cpp
2. **Channel-Based Streaming** 🔄 In Progress - Real-time token delivery through Go channels
3. **Production Features** 📋 Planned - Advanced model management and optimization
4. **Embedded Distribution** 📋 Planned - Models bundled with applications

### Key Technical Decisions
- **Zero Network Dependencies** - No HTTP, no localhost communication, no external services
- **Channel-First Architecture** - Go concurrency primitives for all streaming operations
- **Direct Model Execution** - In-process inference without external servers
- **Production Considerations** - Robust error handling, memory management, resource optimization

## Common Development Tasks

### Adding New Features
1. Define interfaces in `types.go` if needed
2. Implement in `offline_client.go` for direct model execution
3. Add comprehensive examples in `examples/` directory
4. Update tests to cover new functionality (no external dependencies)
5. **Avoid using `client.go`** - focus only on offline implementation

### Working with Streaming
- Use channels for real-time token delivery
- Implement proper error handling via error channels
- Consider goroutine safety and resource cleanup
- Test with concurrent scenarios using Go's testing framework
- Ensure streaming works without any network communication

### Model Integration
- **Focus only on file paths**: `/path/to/model.gguf` (no model names like `llama2:7b`)
- Implement direct GGUF/GGML file loading logic
- Handle model loading/unloading for offline mode
- Consider memory management for large models
- Support embedded model distribution patterns

## Output Formats & API Patterns

### Response Structures
```go
type Response struct {
    Content string
    Error   error
    Tokens  int
}

type Message struct {
    Role    string // "system", "user", "assistant"
    Content string
}
```

### Callback Patterns
```go
type StreamCallback func(token string)

// Usage
client.ChatStream(messages, func(token string) {
    fmt.Print(token)
})
```

## ⚠️ Critical Development Guidelines

### Offline-Only Requirements
**To build a robust solution, all development must focus exclusively on offline operation:**

1. **Never depend on external servers** - Avoid any HTTP communication or network calls
2. **Use only direct model files** - Accept file paths like `/path/to/model.gguf`, not model names
3. **Implement with llama.cpp bindings** - Direct in-process model execution
4. **Channel-based streaming only** - No HTTP response body scanning
5. **Zero external services** - Application must work without any running services

### What to Avoid
- `client.go` HTTP implementation (legacy code)
- Any `BaseURL` or HTTP configuration
- Model names that require external services (`llama2:7b`)
- Network timeouts or HTTP error handling
- Dependencies on any running external services

### What to Build
- Direct GGUF model loading via `offline_client.go`
- Channel-based token streaming using goroutines
- In-process model inference without network calls
- Embeddable model distribution capabilities
- Pure Go concurrency patterns for real-time AI

This architecture enables truly portable, offline AI applications with zero external dependencies while maintaining excellent developer ergonomics through Go's concurrency primitives.