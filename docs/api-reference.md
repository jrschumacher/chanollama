# 📚 chanollama API Reference

Complete API documentation for the chanollama library - a Go library for truly offline LLM inference using channels and goroutines.

## 🎯 Core Interfaces

### Engine Interface

The primary interface for all LLM operations in chanollama.

```go
type Engine interface {
    Chat(messages []Message) (*Response, error)
    ChatStream(messages []Message, callback StreamCallback) *Response
    IsHealthy() bool
    Close() error
}
```

**Implementations:**
- `*Client` - Direct offline model execution
- `*AdvancedEngine` - Message-passing architecture with advanced features

---

## 🏗️ Core Types

### Message

Represents a single message in a conversation.

```go
type Message struct {
    Role      string            `json:"role"`      // "system", "user", "assistant"  
    Content   string            `json:"content"`   // Message content
    Timestamp time.Time         `json:"timestamp,omitempty"`
    Metadata  map[string]string `json:"metadata,omitempty"`
}
```

**Constructor Functions:**
```go
func NewMessage(role, content string) Message
func NewUserMessage(content string) Message
func NewAssistantMessage(content string) Message  
func NewSystemMessage(content string) Message
```

**Methods:**
```go
func (m Message) WithMetadata(key, value string) Message
func (m Message) GetMetadata(key string) (string, bool)
func (m Message) IsRole(role string) bool
func (m Message) IsUser() bool
func (m Message) IsAssistant() bool
func (m Message) IsSystem() bool
func (m Message) Age() time.Duration
func (m Message) String() string
```

### Response

Represents a complete response from the model.

```go
type Response struct {
    Content      string            `json:"content"`
    Error        error             `json:"error,omitempty"`
    Tokens       int               `json:"tokens,omitempty"`
    Duration     time.Duration     `json:"duration,omitempty"`
    Model        string            `json:"model,omitempty"`
    Timestamp    time.Time         `json:"timestamp,omitempty"`
    Metadata     map[string]string `json:"metadata,omitempty"`
    FinishReason string            `json:"finish_reason,omitempty"` // "stop", "length", "error"
}
```

**Constructor & Methods:**
```go
func NewResponse(content string) *Response
func (r *Response) WithError(err error) *Response
func (r *Response) WithTokens(tokens int) *Response
func (r *Response) WithDuration(duration time.Duration) *Response
func (r *Response) WithModel(model string) *Response
func (r *Response) WithFinishReason(reason string) *Response
func (r *Response) WithMetadata(key, value string) *Response
func (r *Response) IsSuccess() bool
func (r *Response) IsError() bool
func (r *Response) ToMessage() Message
```

### Config

Configuration for the LLM engine.

```go
type Config struct {
    Model       string        `json:"model"`       // Model name or path to GGUF file
    Timeout     time.Duration `json:"timeout"`     // Default: 30s
    MaxTokens   int          `json:"max_tokens"`  // Default: 2048
    Temperature float64      `json:"temperature"` // Default: 0.7

    // Offline mode configuration
    OfflineMode bool   `json:"offline_mode"` // Enable direct model execution
    ModelPath   string `json:"model_path"`   // Path to GGUF model file

    // Hardware settings (offline mode)
    UseGPU      bool `json:"use_gpu"`      // Enable GPU acceleration
    NumThreads  int  `json:"num_threads"`  // Number of CPU threads
    ContextSize int  `json:"context_size"` // Context window size

    // Logging settings
    Verbose     bool `json:"verbose"`      // Enable verbose logging
}
```

**Constructor Functions:**
```go
func DefaultConfig(model string) Config
func DefaultOfflineConfig(modelPath string) Config
```

---

## 🚀 Client API

### Client

Direct offline model execution client.

```go
type Client struct {
    // Private fields
}
```

**Constructor Functions:**
```go
func New(modelPath string) *Client
func NewWithConfig(modelPath string, config Config) *Client
```

**Core Methods:**
```go
func (c *Client) Initialize() error
func (c *Client) Chat(messages []Message) (*Response, error)
func (c *Client) ChatStream(messages []Message, callback StreamCallback) *Response
func (c *Client) IsHealthy() bool
func (c *Client) Close() error
```

**Configuration Methods:**
```go
func (c *Client) GetModelPath() string
func (c *Client) GetConfig() Config
func (c *Client) UpdateConfig(config Config)
func (c *Client) GetModelInfo() (*ModelInfo, error)
```

**Conversation Management:**
```go
func (c *Client) GetConversationManager() *ConversationManager
func (c *Client) SetConversationConfig(config ConversationConfig)
func (c *Client) SetConversationLimits(config ConversationConfig)
func (c *Client) AutoDetectConversationLimits() error
```

**Embedding Support:**
```go
func (c *Client) GenerateEmbedding(text string, embeddingModel string) (*EmbeddingResponse, error)
func (c *Client) GenerateEmbeddings(texts []string, embeddingModel string) ([]*EmbeddingResponse, error)
```

---

## ⚡ Advanced Engine API

### AdvancedEngine

Message-passing based engine for complex operations.

```go
type AdvancedEngine struct {
    // Private fields
}
```

**Constructor Functions:**
```go
func NewAdvancedEngine(config Config) *AdvancedEngine
func NewAdvancedEngineWithConfig(config Config, engineConfig EngineConfig) *AdvancedEngine
```

**Lifecycle Methods:**
```go
func (e *AdvancedEngine) Start(ctx context.Context) error
func (e *AdvancedEngine) Stop() error
func (e *AdvancedEngine) IsHealthy() bool
func (e *AdvancedEngine) Close() error
```

**Core Methods:**
```go
func (e *AdvancedEngine) Chat(messages []Message) (*Response, error)
func (e *AdvancedEngine) ChatStream(messages []Message, callback StreamCallback) *Response
func (e *AdvancedEngine) SubmitRequest(req EngineRequest) error
func (e *AdvancedEngine) GetStreamingChannel() <-chan StreamingResponse
```

### Operation Types

```go
type OperationType int

const (
    OpInference OperationType = iota
    OpStream
    OpGetMetrics
    OpHealthCheck
    OpConfigureGeneration
    OpTokenCount
    OpShutdown
)
```

### Engine Request/Response

```go
type EngineRequest struct {
    ID       string
    Type     OperationType
    Payload  interface{}
    Response chan<- EngineResponse
    Context  context.Context
}

type EngineResponse struct {
    RequestID string
    Success   bool
    Data      interface{}
    Error     error
    Metadata  map[string]interface{}
}
```

---

## 🔄 Streaming API

### StreamCallback

Function type for streaming callbacks.

```go
type StreamCallback func(token string)
```

### StreamingEngine

Channel-based streaming wrapper.

```go
type StreamingEngine struct {
    // Private fields
}
```

**Constructor:**
```go
func NewStreamingEngine(client Engine) *StreamingEngine
```

**Methods:**
```go
func (s *StreamingEngine) ChatStreamChan(messages []Message) (<-chan string, <-chan error)
func (s *StreamingEngine) ChatStreamChanWithTimeout(messages []Message, timeout time.Duration) (<-chan string, <-chan error)
```

### ConcurrentStreamer

Manages multiple concurrent streaming operations.

```go
type ConcurrentStreamer struct {
    // Private fields
}
```

**Constructor:**
```go
func NewConcurrentStreamer(engine Engine, maxConnections int) *ConcurrentStreamer
```

**Methods:**
```go
func (cs *ConcurrentStreamer) Stream(messages []Message) (<-chan string, <-chan error)
func (cs *ConcurrentStreamer) Close() error
```

### StreamCollector

Collects streaming tokens into a complete response.

```go
type StreamCollector struct {
    // Private fields
}

func NewStreamCollector() *StreamCollector
func (sc *StreamCollector) Collect(tokenCh <-chan string, errorCh <-chan error) (*Response, error)
```

### TokenBatcher

Batches streaming tokens for efficient processing.

```go
type TokenBatcher struct {
    BatchSize int
    Interval  time.Duration
}

func NewTokenBatcher(batchSize int, interval time.Duration) *TokenBatcher
func (tb *TokenBatcher) BatchTokens(tokenCh <-chan string) <-chan []string
```

---

## 💬 Conversation Management

### Conversation

Helper for managing chat history.

```go
type Conversation struct {
    Messages    []Message `json:"messages"`
    MaxMessages int       `json:"max_messages"`
    TokenLimit  int       `json:"token_limit"`
}
```

**Constructor:**
```go
func NewConversation() *Conversation
```

**Message Management:**
```go
func (c *Conversation) AddUserMessage(content string)
func (c *Conversation) AddAssistantMessage(content string)
func (c *Conversation) AddSystemMessage(content string)
func (c *Conversation) AddMessage(msg Message)
func (c *Conversation) AddResponse(response *Response)
```

**Query Methods:**
```go
func (c *Conversation) GetLastMessage() (Message, bool)
func (c *Conversation) GetLastUserMessage() (Message, bool)
func (c *Conversation) GetSystemMessages() []Message
func (c *Conversation) FilterMessagesByRole(role string) []Message
func (c *Conversation) GetMessageCount() int
func (c *Conversation) GetMessagesByAge(maxAge time.Duration) []Message
```

**Management Methods:**
```go
func (c *Conversation) Clear()
func (c *Conversation) ClearExceptSystem()
func (c *Conversation) EstimateTokens() int
```

### ConversationConfig

Configuration for conversation management.

```go
type ConversationConfig struct {
    MaxTokens          int                `json:"max_tokens"`
    MaxMessages        int                `json:"max_messages"`
    TruncationStrategy TruncationStrategy `json:"truncation_strategy"`
    PreserveSystem     bool               `json:"preserve_system"`
    AutoDetectLimits   bool               `json:"auto_detect_limits"`
    BatchSize          int                `json:"batch_size"`
    EnableLogging      bool               `json:"enable_logging"`
    SummaryPrompt      string             `json:"summary_prompt"`
}
```

**Constructor:**
```go
func DefaultConversationConfig() ConversationConfig
```

### TruncationStrategy

```go
type TruncationStrategy string

const (
    TruncationOldest    TruncationStrategy = "oldest"
    TruncationNewest    TruncationStrategy = "newest"
    TruncationMiddle    TruncationStrategy = "middle"
    TruncationSummarize TruncationStrategy = "summarize"
)
```

---

## 🏭 Factory Functions

### Engine Creation

```go
func NewEngine(config Config) (Engine, error)
func NewOfflineEngine(config Config) (Engine, error)
```

**Auto-Detection Logic:**
- Detects if model path looks like a file path (contains `/` or ends with `.gguf`)
- Automatically enables offline mode for file paths
- Falls back to offline client by default

---

## 🔧 Utility APIs

### TokenEstimator

Provides token estimation utilities.

```go
type TokenEstimator struct {
    CharsPerToken float64
}

func NewTokenEstimator() *TokenEstimator
func (te *TokenEstimator) WithCharsPerToken(ratio float64) *TokenEstimator
func (te *TokenEstimator) EstimateTokens(text string) int
func (te *TokenEstimator) EstimateMessageTokens(msg Message) int
func (te *TokenEstimator) EstimateMessagesTokens(messages []Message) int
```

### TextProcessor

Text processing utilities for LLM inputs.

```go
type TextProcessor struct {
    MaxLength    int
    CleanText    bool
    PreserveCode bool
}

func NewTextProcessor() *TextProcessor
func (tp *TextProcessor) WithMaxLength(length int) *TextProcessor
func (tp *TextProcessor) WithCleanText(clean bool) *TextProcessor
func (tp *TextProcessor) WithPreserveCode(preserve bool) *TextProcessor
func (tp *TextProcessor) ProcessText(text string) string
```

### ResourceMonitor

System resource monitoring for LLM operations.

```go
type ResourceMonitor struct {
    // Private fields
}

func NewResourceMonitor() *ResourceMonitor
func (rm *ResourceMonitor) GetMemoryUsage() (uint64, error)
func (rm *ResourceMonitor) GetCPUUsage() (float64, error)
func (rm *ResourceMonitor) IsMemoryAvailable(requiredMB uint64) bool
```

---

## 📊 Metrics & Monitoring

### EngineMetrics

Performance and resource metrics.

```go
type EngineMetrics struct {
    TotalRequests       int64         `json:"total_requests"`
    SuccessfulRequests  int64         `json:"successful_requests"`
    FailedRequests      int64         `json:"failed_requests"`
    AverageResponseTime time.Duration `json:"average_response_time"`
    ModelName           string        `json:"model_name"`
    TokensGenerated     int64         `json:"tokens_generated"`
    ActiveRequests      int           `json:"active_requests"`
    QueuedRequests      int           `json:"queued_requests"`
    LastActivity        time.Time     `json:"last_activity"`
    Uptime              time.Duration `json:"uptime"`
    StartTime           time.Time     `json:"start_time"`
    LastError           string        `json:"last_error,omitempty"`
    ErrorCount          map[string]int `json:"error_count"`
}
```

**Methods:**
```go
func NewEngineMetrics() *EngineMetrics
func (m *EngineMetrics) IncrementRequests()
func (m *EngineMetrics) IncrementSuccessful()
func (m *EngineMetrics) IncrementFailed(errorType string)
func (m *EngineMetrics) SetActiveRequests(count int)
func (m *EngineMetrics) SetQueuedRequests(count int)
func (m *EngineMetrics) UpdateResponseTime(duration time.Duration)
func (m *EngineMetrics) SetModelInfo(name string)
func (m *EngineMetrics) AddTokensGenerated(count int64)
func (m *EngineMetrics) GetSnapshot() EngineMetrics
```

### HealthStatus

```go
type HealthStatus struct {
    Status    string        `json:"status"`    // "healthy", "degraded", "unhealthy"
    Timestamp time.Time     `json:"timestamp"`
    Uptime    time.Duration `json:"uptime"`
    LastError string        `json:"last_error,omitempty"`
    Version   string        `json:"version,omitempty"`
}
```

---

## ❌ Error Types

### Predefined Errors

```go
var (
    ErrBatchSizeExceeded = errors.New("conversation too long for model batch size")
    ErrTokenCountFailed  = errors.New("unable to count tokens, using message limit fallback")
    ErrModelNotLoaded    = errors.New("model not loaded or initialization failed")
    ErrInferenceFailed   = errors.New("inference operation failed")
    ErrContextOverflow   = errors.New("context size exceeded, conversation truncated")
)
```

### Error Handling Patterns

```go
// Check for specific error types
if errors.Is(err, chanollama.ErrModelNotLoaded) {
    // Handle model loading error
}

// Check response errors
if response.IsError() {
    log.Printf("Response error: %v", response.Error)
}

// Handle panic recovery (built-in)
response, err := client.Chat(messages)
if err != nil {
    // Panics are automatically recovered and converted to errors
}
```

---

## 🎯 Usage Examples

### Basic Offline Usage

```go
// Create offline client
client := chanollama.New("/path/to/model.gguf")
defer client.Close()

// Initialize (optional - called automatically on first use)
if err := client.Initialize(); err != nil {
    log.Fatal(err)
}

// Simple chat
messages := []chanollama.Message{
    chanollama.NewUserMessage("Hello!"),
}

response, err := client.Chat(messages)
if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Content)
```

### Streaming with Channels

```go
// Create streaming engine
engine := chanollama.NewStreamingEngine(client)

// Get token and error channels
tokenCh, errorCh := engine.ChatStreamChan(messages)

// Process streaming response
for {
    select {
    case token, ok := <-tokenCh:
        if !ok {
            return // Streaming complete
        }
        fmt.Print(token)
        
    case err := <-errorCh:
        if err != nil {
            log.Printf("Error: %v", err)
            return
        }
    }
}
```

### Advanced Engine Operations

```go
// Create advanced engine
engine := chanollama.NewAdvancedEngine(chanollama.DefaultConfig("/path/to/model.gguf"))

// Start engine
ctx := context.Background()
if err := engine.Start(ctx); err != nil {
    log.Fatal(err)
}
defer engine.Stop()

// Submit request
responseCh := make(chan chanollama.EngineResponse, 1)
request := chanollama.EngineRequest{
    ID:       "req-1",
    Type:     chanollama.OpInference,
    Payload:  chanollama.InferencePayload{Messages: messages},
    Response: responseCh,
    Context:  context.Background(),
}

if err := engine.SubmitRequest(request); err != nil {
    log.Fatal(err)
}

// Wait for response
response := <-responseCh
if response.Success {
    fmt.Println(response.Data)
}
```

---

This API reference covers all public APIs in the chanollama library. For more examples and usage patterns, see the [examples directory](../examples/) and other documentation files.