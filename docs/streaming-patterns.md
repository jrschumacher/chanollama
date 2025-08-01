# 🔄 Streaming and Channel Patterns

Advanced guide to real-time streaming patterns in chanollama using Go channels and goroutines.

## 🎯 Overview

chanollama provides multiple streaming approaches leveraging Go's native concurrency primitives:

- **Callback Streaming** - Simple function-based token delivery
- **Channel Streaming** - Go-native channel-based token delivery  
- **Concurrent Streaming** - Multiple simultaneous streaming operations
- **Batched Streaming** - Token batching for efficient processing
- **Advanced Patterns** - Custom streaming workflows and error handling

---

## 🔧 Streaming Fundamentals

### Understanding Token Streaming

In LLM inference, responses are generated token by token. Streaming delivers these tokens in real-time rather than waiting for the complete response:

```go
// Without streaming - wait for complete response
response, err := client.Chat(messages)
fmt.Println(response.Content) // All tokens delivered at once

// With streaming - tokens delivered as generated
client.ChatStream(messages, func(token string) {
    fmt.Print(token) // Each token appears immediately
})
```

### Token Types

Different types of tokens in streaming:

```go
func analyzeTokens() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain Go programming"),
    }

    var tokenCount int
    var wordCount int
    var currentWord strings.Builder

    client.ChatStream(messages, func(token string) {
        tokenCount++

        // Tokens can be words, punctuation, or partial words
        if strings.Contains(token, " ") || strings.ContainsAny(token, ".,!?;") {
            if currentWord.Len() > 0 {
                wordCount++
                currentWord.Reset()
            }
        } else {
            currentWord.WriteString(token)
        }

        fmt.Print(token)
    })

    fmt.Printf("\nTokens: %d, Words: ~%d\n", tokenCount, wordCount)
}
```

---

## 📞 Callback-Based Streaming

### Basic Callback Streaming

The simplest streaming approach using function callbacks:

```go
func basicCallbackStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Write a short story about AI"),
    }

    fmt.Print("Story: ")
    
    response := client.ChatStream(messages, func(token string) {
        fmt.Print(token) // Real-time token display
    })

    if response.Error != nil {
        log.Printf("\nStreaming error: %v", response.Error)
    }
    
    fmt.Println() // Newline after complete response
}
```

### Advanced Callback Patterns

Collect streaming data while displaying:

```go
func advancedCallbackStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain machine learning algorithms"),
    }

    var fullResponse strings.Builder
    var tokenCount int
    startTime := time.Now()

    fmt.Print("Response: ")

    response := client.ChatStream(messages, func(token string) {
        tokenCount++
        fullResponse.WriteString(token)
        fmt.Print(token)

        // Progress indicator every 10 tokens
        if tokenCount%10 == 0 {
            elapsed := time.Since(startTime)
            fmt.Printf(" [%d tokens, %.1fs]", tokenCount, elapsed.Seconds())
        }
    })

    duration := time.Since(startTime)
    
    if response.Error != nil {
        log.Printf("\nError: %v", response.Error)
        return
    }

    fmt.Printf("\n\nStreaming complete:")
    fmt.Printf("- Tokens: %d\n", tokenCount)
    fmt.Printf("- Duration: %v\n", duration)
    fmt.Printf("- Tokens/sec: %.1f\n", float64(tokenCount)/duration.Seconds())
    fmt.Printf("- Characters: %d\n", fullResponse.Len())
}
```

### Callback with UI Updates

Integration with user interfaces:

```go
type UIUpdater struct {
    content strings.Builder
    mu      sync.Mutex
}

func (ui *UIUpdater) AppendToken(token string) {
    ui.mu.Lock()
    defer ui.mu.Unlock()
    ui.content.WriteString(token)
    // Update your UI here - send to web socket, update GUI, etc.
}

func (ui *UIUpdater) GetContent() string {
    ui.mu.Lock()
    defer ui.mu.Unlock()
    return ui.content.String()
}

func streamToUI() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    ui := &UIUpdater{}
    
    messages := []chanollama.Message{
        chanollama.NewUserMessage("Describe the future of AI"),
    }

    // Stream to UI updater
    response := client.ChatStream(messages, ui.AppendToken)

    if response.Error != nil {
        log.Printf("Streaming error: %v", response.Error)
        return
    }

    // Final content available
    finalContent := ui.GetContent()
    fmt.Printf("Final response (%d chars): %s\n", len(finalContent), finalContent)
}
```

---

## 📡 Channel-Based Streaming

### Basic Channel Streaming

Go-native channel streaming for maximum control:

```go
func basicChannelStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain Go channels and select statements"),
    }

    // Get channels for tokens and errors
    tokenCh, errorCh := engine.ChatStreamChan(messages)

    fmt.Print("Response: ")

    // Process streaming with select statement
    for {
        select {
        case token, ok := <-tokenCh:
            if !ok {
                fmt.Println() // Newline when complete
                return        // Channel closed, streaming finished
            }
            fmt.Print(token)

        case err := <-errorCh:
            if err != nil {
                log.Printf("\nStreaming error: %v", err)
                return
            }
            // No error, continue streaming
        }
    }
}
```

### Channel Streaming with Timeout

Add timeout control to streaming operations:

```go
func channelStreamingWithTimeout() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Write a detailed technical analysis"),
    }

    // Set streaming timeout
    timeout := 30 * time.Second
    tokenCh, errorCh := engine.ChatStreamChanWithTimeout(messages, timeout)

    fmt.Print("Response: ")

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
                    fmt.Printf("\n⏰ Streaming timed out after %v\n", timeout)
                } else {
                    fmt.Printf("\n❌ Streaming error: %v\n", err)
                }
                return
            }

        case <-time.After(timeout):
            fmt.Printf("\n⏰ Overall timeout after %v\n", timeout)
            return
        }
    }
}
```

### Multi-Channel Processing

Process multiple streams simultaneously:

```go
func multiChannelStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    questions := []string{
        "What is artificial intelligence?",
        "How do neural networks work?",
        "Explain machine learning types",
    }

    // Create channels for each question
    var tokenChannels []<-chan string
    var errorChannels []<-chan error

    // Start all streams
    for _, question := range questions {
        messages := []chanollama.Message{
            chanollama.NewUserMessage(question),
        }
        tokenCh, errorCh := engine.ChatStreamChan(messages)
        tokenChannels = append(tokenChannels, tokenCh)
        errorChannels = append(errorChannels, errorCh)
    }

    // Process all streams
    activeStreams := len(questions)
    responses := make([]strings.Builder, len(questions))

    for activeStreams > 0 {
        // Check each stream
        for i := 0; i < len(questions); i++ {
            if tokenChannels[i] == nil {
                continue // Stream already closed
            }

            select {
            case token, ok := <-tokenChannels[i]:
                if !ok {
                    tokenChannels[i] = nil // Mark stream as closed
                    activeStreams--
                    fmt.Printf("\n✅ Stream %d complete\n", i+1)
                    continue
                }
                responses[i].WriteString(token)
                fmt.Printf("[%d] %s", i+1, token)

            case err := <-errorChannels[i]:
                if err != nil {
                    fmt.Printf("\n❌ Stream %d error: %v\n", i+1, err)
                    tokenChannels[i] = nil
                    activeStreams--
                }

            default:
                // Non-blocking, check next stream
            }
        }
    }

    // Print final responses
    for i, response := range responses {
        fmt.Printf("\nFinal Response %d:\n%s\n", i+1, response.String())
    }
}
```

---

## 🔀 Concurrent Streaming

### Concurrent Streamer

Handle multiple simultaneous streaming requests:

```go
func concurrentStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    // Create concurrent streamer with max 5 concurrent operations
    streamer := chanollama.NewConcurrentStreamer(client, 5)
    defer streamer.Close()

    questions := []string{
        "Explain quantum computing",
        "What is blockchain technology?", 
        "How does machine learning work?",
        "Describe cloud computing",
        "What is artificial intelligence?",
    }

    var wg sync.WaitGroup
    
    // Start concurrent streams
    for i, question := range questions {
        wg.Add(1)
        go func(id int, q string) {
            defer wg.Done()

            messages := []chanollama.Message{
                chanollama.NewUserMessage(q),
            }

            fmt.Printf("🚀 Starting stream %d: %s\n", id+1, q)

            // Get streaming channels
            tokenCh, errorCh := streamer.Stream(messages)

            var response strings.Builder
            fmt.Printf("[%d] ", id+1)

            // Process stream
            for {
                select {
                case token, ok := <-tokenCh:
                    if !ok {
                        fmt.Printf("\n✅ Stream %d complete (%d chars)\n", 
                            id+1, response.Len())
                        return
                    }
                    response.WriteString(token)
                    fmt.Print(token)

                case err := <-errorCh:
                    if err != nil {
                        fmt.Printf("\n❌ Stream %d error: %v\n", id+1, err)
                        return
                    }
                }
            }
        }(i, question)
    }

    // Wait for all streams to complete
    wg.Wait()
    fmt.Println("🏁 All concurrent streams complete")
}
```

### Rate-Limited Concurrent Streaming

Control the rate of concurrent operations:

```go
func rateLimitedStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    streamer := chanollama.NewConcurrentStreamer(client, 3) // Max 3 concurrent
    defer streamer.Close()

    questions := []string{
        "Explain AI ethics",
        "What is deep learning?",
        "How do transformers work?",
        "Describe reinforcement learning",
        "What is computer vision?",
        "Explain natural language processing",
    }

    // Rate limiter - 1 new stream every 2 seconds
    rateLimiter := time.NewTicker(2 * time.Second)
    defer rateLimiter.Stop()

    var wg sync.WaitGroup

    for i, question := range questions {
        // Wait for rate limiter
        <-rateLimiter.C

        wg.Add(1)
        go func(id int, q string) {
            defer wg.Done()
            
            fmt.Printf("🚀 Starting rate-limited stream %d\n", id+1)

            messages := []chanollama.Message{
                chanollama.NewUserMessage(q),
            }

            tokenCh, errorCh := streamer.Stream(messages)

            // Process stream (similar to previous example)
            for {
                select {
                case token, ok := <-tokenCh:
                    if !ok {
                        fmt.Printf("✅ Rate-limited stream %d complete\n", id+1)
                        return
                    }
                    fmt.Printf("[%d] %s", id+1, token)

                case err := <-errorCh:
                    if err != nil {
                        fmt.Printf("❌ Rate-limited stream %d error: %v\n", id+1, err)
                        return
                    }
                }
            }
        }(i, question)
    }

    wg.Wait()
    fmt.Println("🏁 All rate-limited streams complete")
}
```

---

## 📦 Token Batching

### Basic Token Batching

Batch tokens for efficient processing:

```go
func basicTokenBatching() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Write a long explanation of distributed systems"),
    }

    // Get token stream
    tokenCh, errorCh := engine.ChatStreamChan(messages)

    // Create token batcher - batch 5 tokens every 100ms
    batcher := chanollama.NewTokenBatcher(5, 100*time.Millisecond)
    batchCh := batcher.BatchTokens(tokenCh)

    fmt.Print("Batched Response: ")

    // Process batched tokens
    for {
        select {
        case batch, ok := <-batchCh:
            if !ok {
                fmt.Println()
                return // Batching complete
            }

            // Process batch of tokens
            fmt.Printf("[Batch of %d] ", len(batch))
            for _, token := range batch {
                fmt.Print(token)
            }

        case err := <-errorCh:
            if err != nil {
                fmt.Printf("\nBatching error: %v\n", err)
                return
            }
        }
    }
}
```

### Adaptive Token Batching

Adjust batch size based on token rate:

```go
func adaptiveTokenBatching() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain complex algorithms in detail"),
    }

    tokenCh, errorCh := engine.ChatStreamChan(messages)

    // Start with small batches
    batchSize := 3
    batchInterval := 200 * time.Millisecond
    
    fmt.Print("Adaptive Batched Response: ")

    var tokenCount int
    var lastBatchTime time.Time = time.Now()

    // Custom batching logic
    var currentBatch []string
    ticker := time.NewTicker(batchInterval)
    defer ticker.Stop()

    for {
        select {
        case token, ok := <-tokenCh:
            if !ok {
                // Flush final batch
                if len(currentBatch) > 0 {
                    fmt.Printf("[Final batch: %d] ", len(currentBatch))
                    for _, t := range currentBatch {
                        fmt.Print(t)
                    }
                }
                fmt.Println()
                return
            }

            currentBatch = append(currentBatch, token)
            tokenCount++

            // Flush batch when full
            if len(currentBatch) >= batchSize {
                processBatch(currentBatch, batchSize)
                currentBatch = nil
                
                // Adapt batch size based on token rate
                elapsed := time.Since(lastBatchTime)
                if elapsed < 100*time.Millisecond {
                    // Tokens coming fast, increase batch size
                    if batchSize < 10 {
                        batchSize++
                    }
                } else if elapsed > 500*time.Millisecond {
                    // Tokens coming slow, decrease batch size
                    if batchSize > 1 {
                        batchSize--
                    }
                }
                lastBatchTime = time.Now()
            }

        case <-ticker.C:
            // Flush partial batch on interval
            if len(currentBatch) > 0 {
                processBatch(currentBatch, batchSize)
                currentBatch = nil
                lastBatchTime = time.Now()
            }

        case err := <-errorCh:
            if err != nil {
                fmt.Printf("\nAdaptive batching error: %v\n", err)
                return
            }
        }
    }
}

func processBatch(batch []string, expectedSize int) {
    fmt.Printf("[Batch %d/%d] ", len(batch), expectedSize)
    for _, token := range batch {
        fmt.Print(token)
    }
}
```

---

## 🔍 Stream Collection and Analysis

### Stream Collector

Collect streaming response into a complete Response object:

```go
func streamCollection() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain the benefits of streaming responses"),
    }

    // Start streaming
    tokenCh, errorCh := engine.ChatStreamChan(messages)

    // Collect streaming response
    collector := chanollama.NewStreamCollector()
    
    fmt.Print("Streaming (collecting): ")

    // Display tokens while collecting
    go func() {
        for token := range tokenCh {
            fmt.Print(token)
        }
    }()

    // Collect complete response
    response, err := collector.Collect(tokenCh, errorCh)
    
    if err != nil {
        log.Printf("\nCollection error: %v", err)
        return
    }

    fmt.Printf("\n\nCollected Response:")
    fmt.Printf("- Content length: %d characters\n", len(response.Content))
    fmt.Printf("- Estimated tokens: %d\n", response.Tokens)
    fmt.Printf("- Duration: %v\n", response.Duration)
    fmt.Printf("- First 100 chars: %.100s...\n", response.Content)
}
```

### Stream Metrics and Analysis

Analyze streaming performance:

```go
type StreamMetrics struct {
    TokenCount    int
    StartTime     time.Time
    LastTokenTime time.Time
    TokenTimes    []time.Duration
    TotalChars    int
    ErrorCount    int
    mu            sync.Mutex
}

func (sm *StreamMetrics) AddToken(token string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    now := time.Now()
    if sm.StartTime.IsZero() {
        sm.StartTime = now
    }
    
    if !sm.LastTokenTime.IsZero() {
        intervalSinceLastToken := now.Sub(sm.LastTokenTime)
        sm.TokenTimes = append(sm.TokenTimes, intervalSinceLastToken)
    }
    
    sm.TokenCount++
    sm.TotalChars += len(token)
    sm.LastTokenTime = now
}

func (sm *StreamMetrics) AddError() {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.ErrorCount++
}

func (sm *StreamMetrics) GetStats() (avgTokenInterval time.Duration, tokensPerSecond float64, totalDuration time.Duration) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    if len(sm.TokenTimes) == 0 {
        return 0, 0, 0
    }
    
    var totalInterval time.Duration
    for _, interval := range sm.TokenTimes {
        totalInterval += interval
    }
    
    avgTokenInterval = totalInterval / time.Duration(len(sm.TokenTimes))
    totalDuration = sm.LastTokenTime.Sub(sm.StartTime)
    tokensPerSecond = float64(sm.TokenCount) / totalDuration.Seconds()
    
    return
}

func streamAnalysis() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Write a comprehensive guide to Go programming"),
    }

    tokenCh, errorCh := engine.ChatStreamChan(messages)
    metrics := &StreamMetrics{}

    fmt.Print("Analyzing stream: ")

    // Process streaming with metrics
    for {
        select {
        case token, ok := <-tokenCh:
            if !ok {
                // Streaming complete, analyze results
                avgInterval, tokensPerSec, totalDuration := metrics.GetStats()
                
                fmt.Printf("\n\n📊 Stream Analysis:")
                fmt.Printf("- Total tokens: %d\n", metrics.TokenCount)
                fmt.Printf("- Total characters: %d\n", metrics.TotalChars)
                fmt.Printf("- Total duration: %v\n", totalDuration)
                fmt.Printf("- Average token interval: %v\n", avgInterval)
                fmt.Printf("- Tokens per second: %.2f\n", tokensPerSec)
                fmt.Printf("- Characters per token: %.2f\n", 
                    float64(metrics.TotalChars)/float64(metrics.TokenCount))
                fmt.Printf("- Errors encountered: %d\n", metrics.ErrorCount)
                return
            }

            metrics.AddToken(token)
            fmt.Print(token)

        case err := <-errorCh:
            if err != nil {
                metrics.AddError()
                fmt.Printf("\nStream error: %v\n", err)
                return
            }
        }
    }
}
```

---

## 🚨 Advanced Error Handling

### Stream Error Recovery

Implement robust error handling for streaming:

```go
func resilientStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain error handling in distributed systems"),
    }

    maxRetries := 3
    retryDelay := time.Second

    for attempt := 1; attempt <= maxRetries; attempt++ {
        fmt.Printf("🔄 Streaming attempt %d/%d\n", attempt, maxRetries)

        tokenCh, errorCh := engine.ChatStreamChan(messages)
        var response strings.Builder
        var tokenCount int
        streamError := false

        // Process stream with error detection
        for {
            select {
            case token, ok := <-tokenCh:
                if !ok {
                    if tokenCount > 0 {
                        fmt.Printf("\n✅ Streaming successful after %d attempts\n", attempt)
                        fmt.Printf("Response (%d tokens): %s\n", tokenCount, response.String())
                        return
                    } else {
                        fmt.Printf("❌ Empty response on attempt %d\n", attempt)
                        streamError = true
                        break
                    }
                }

                tokenCount++
                response.WriteString(token)
                fmt.Print(token)

                // Detect streaming issues
                if tokenCount == 1 && strings.TrimSpace(token) == "" {
                    fmt.Printf("⚠️ Detected empty first token on attempt %d\n", attempt)
                }

            case err := <-errorCh:
                if err != nil {
                    fmt.Printf("\n❌ Stream error on attempt %d: %v\n", attempt, err)
                    streamError = true
                    break
                }

            case <-time.After(30 * time.Second):
                fmt.Printf("\n⏰ Stream timeout on attempt %d\n", attempt)
                streamError = true
                break
            }
        }

        if !streamError {
            break // Success
        }

        if attempt < maxRetries {
            fmt.Printf("⏳ Retrying in %v...\n", retryDelay)
            time.Sleep(retryDelay)
            retryDelay *= 2 // Exponential backoff
        } else {
            fmt.Printf("💥 All streaming attempts failed\n")
        }
    }
}
```

### Circuit Breaker Pattern

Implement circuit breaker for streaming operations:

```go
type StreamCircuitBreaker struct {
    maxFailures int
    failures    int
    lastFailure time.Time
    timeout     time.Duration
    state       string // "closed", "open", "half-open"
    mu          sync.Mutex
}

func NewStreamCircuitBreaker(maxFailures int, timeout time.Duration) *StreamCircuitBreaker {
    return &StreamCircuitBreaker{
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       "closed",
    }
}

func (cb *StreamCircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    // Check if circuit should reset
    if cb.state == "open" && time.Since(cb.lastFailure) > cb.timeout {
        cb.state = "half-open"
        cb.failures = 0
    }

    // Reject calls if circuit is open
    if cb.state == "open" {
        return fmt.Errorf("circuit breaker is open")
    }

    // Execute function
    err := fn()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }

    // Success - reset failures
    cb.failures = 0
    if cb.state == "half-open" {
        cb.state = "closed"
    }

    return nil
}

func circuitBreakerStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()

    engine := chanollama.NewStreamingEngine(client)
    
    // Circuit breaker: 3 failures, 5 second timeout
    cb := NewStreamCircuitBreaker(3, 5*time.Second)

    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain circuit breaker patterns"),
    }

    for i := 0; i < 10; i++ {
        fmt.Printf("🔄 Stream attempt %d: ", i+1)
        
        err := cb.Call(func() error {
            tokenCh, errorCh := engine.ChatStreamChan(messages)
            
            timeout := time.After(10 * time.Second)
            tokenCount := 0

            for {
                select {
                case token, ok := <-tokenCh:
                    if !ok {
                        if tokenCount > 10 {
                            fmt.Printf("✅ Success (%d tokens)\n", tokenCount)
                            return nil
                        }
                        return fmt.Errorf("insufficient tokens received")
                    }
                    tokenCount++

                case err := <-errorCh:
                    if err != nil {
                        return fmt.Errorf("stream error: %w", err)
                    }

                case <-timeout:
                    return fmt.Errorf("stream timeout")
                }
            }
        })

        if err != nil {
            fmt.Printf("❌ %v\n", err)
        }

        time.Sleep(1 * time.Second)
    }
}
```

---

## 🎯 Best Practices

### 1. **Choose the Right Streaming Pattern**

- **Callback streaming**: Simple UI updates, logging, progress indication
- **Channel streaming**: Complex processing, timeouts, concurrent operations
- **Batched streaming**: High-throughput scenarios, network transmission
- **Concurrent streaming**: Multiple simultaneous requests

### 2. **Error Handling Strategy**

- Always handle both token and error channels
- Implement timeouts for streaming operations
- Use circuit breakers for unreliable connections
- Provide graceful degradation for streaming failures

### 3. **Performance Considerations**

- Buffer channels appropriately (100-1000 tokens typically)
- Use goroutines for non-blocking UI updates
- Implement backpressure handling for slow consumers
- Monitor memory usage in long streaming operations

### 4. **Resource Management**

- Always close clients and streamers
- Use context cancellation for cleanup
- Implement proper goroutine lifecycle management
- Monitor for goroutine leaks in long-running applications

### 5. **Concurrency Safety**

- Use mutexes when sharing data across goroutines
- Be careful with shared state in streaming callbacks
- Use channels for communication between goroutines
- Avoid blocking operations in callback functions

---

This guide covers the essential streaming patterns in chanollama. For more advanced topics and API details, see the [API Reference](./api-reference.md) and [examples directory](../examples/).