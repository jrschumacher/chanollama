# 🔧 Troubleshooting Guide

Complete troubleshooting guide for chanollama - solutions to common issues, debugging techniques, and best practices.

## 🚨 Common Issues

### Model Loading Problems

#### "Model file not found"
```
Error: model file not found: /path/to/model.gguf
```

**Solutions:**
```go
// Check if file exists
if _, err := os.Stat(modelPath); os.IsNotExist(err) {
    log.Fatalf("Model file not found: %s", modelPath)
}

// Use absolute path
absPath, err := filepath.Abs(modelPath)
if err != nil {
    log.Fatal("Failed to resolve path:", err)
}
client := chanollama.New(absPath)
```

**Verification:**
```bash
# Check file exists and permissions
ls -la /path/to/model.gguf

# Verify it's a valid GGUF file
file /path/to/model.gguf
```

#### "Failed to load model"
```
Error: failed to load model: llama_model_load_internal: failed to load model
```

**Causes & Solutions:**

1. **Corrupted Model File:**
```bash
# Check file integrity
shasum -a 256 model.gguf
# Compare with expected hash from download source
```

2. **Insufficient Memory:**
```go
// Check available memory
monitor := chanollama.NewResourceMonitor()
memMB, err := monitor.GetMemoryUsage()
if err == nil {
    fmt.Printf("Available memory: %d MB\n", memMB/(1024*1024))
}

// Use smaller context size
config := chanollama.DefaultOfflineConfig(modelPath)
config.ContextSize = 2048  // Reduce from 4096
client := chanollama.NewWithConfig(modelPath, config)
```

3. **Invalid GGUF Format:**
```go
// Validate model file before loading
func validateGGUFFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Check GGUF magic bytes
    magic := make([]byte, 4)
    _, err = file.Read(magic)
    if err != nil {
        return err
    }
    
    if string(magic) != "GGUF" {
        return fmt.Errorf("not a valid GGUF file")
    }
    
    return nil
}
```

---

### Inference Problems

#### "Inference operation failed"
```
Error: inference operation failed: context decode failed
```

**Solutions:**

1. **Context Overflow:**
```go
// Enable automatic conversation management
config := chanollama.ConversationConfig{
    MaxTokens:          2048,                          // Reduce context usage
    TruncationStrategy: chanollama.TruncationOldest,   // Auto-truncate
    PreserveSystem:     true,
}
client.SetConversationConfig(config)
```

2. **Batch Size Issues:**
```go
// Check and configure conversation limits
if err := client.AutoDetectConversationLimits(); err != nil {
    log.Printf("Auto-detection failed: %v", err)
    
    // Manual configuration
    config := chanollama.ConversationConfig{
        MaxTokens:   1600,  // 80% of typical batch size
        BatchSize:   2048,  // Conservative batch size
    }
    client.SetConversationConfig(config)
}
```

3. **Memory Exhaustion:**
```go
// Monitor memory during inference
func safeInference(client *chanollama.Client, messages []chanollama.Message) (*chanollama.Response, error) {
    monitor := chanollama.NewResourceMonitor()
    
    beforeMem, _ := monitor.GetMemoryUsage()
    
    // Check available memory
    if !monitor.IsMemoryAvailable(1024) { // 1GB minimum
        return nil, fmt.Errorf("insufficient memory for inference")
    }
    
    response, err := client.Chat(messages)
    
    afterMem, _ := monitor.GetMemoryUsage()
    memUsed := (afterMem - beforeMem) / (1024 * 1024)
    
    log.Printf("Inference used %d MB", memUsed)
    return response, err
}
```

#### "Empty or Invalid Response"
```
Error: model returned empty response
```

**Debugging:**
```go
func debugInference(client *chanollama.Client, messages []chanollama.Message) {
    // Check client health
    if !client.IsHealthy() {
        log.Fatal("Client is not healthy")
    }
    
    // Check message format
    for i, msg := range messages {
        if msg.Role == "" || msg.Content == "" {
            log.Printf("Invalid message %d: role='%s' content='%s'", i, msg.Role, msg.Content)
        }
    }
    
    // Try with simple test message
    testMessages := []chanollama.Message{
        chanollama.NewUserMessage("Hello"),
    }
    
    response, err := client.Chat(testMessages)
    if err != nil {
        log.Printf("Test inference failed: %v", err)
    } else if len(response.Content) == 0 {
        log.Printf("Test inference returned empty content")
    } else {
        log.Printf("Test inference successful: %s", response.Content)
    }
}
```

---

### Performance Issues

#### "Slow Response Times"
```
Inference taking >30 seconds
```

**Optimization Strategies:**

1. **CPU Thread Optimization:**
```go
func optimizeForSpeed() chanollama.Config {
    config := chanollama.DefaultOfflineConfig("/path/to/model.gguf")
    
    // Use all available CPU cores
    config.NumThreads = runtime.NumCPU()
    
    // Reduce response length for speed
    config.MaxTokens = 512
    
    // Lower temperature for more direct responses
    config.Temperature = 0.3
    
    return config
}
```

2. **Model Selection:**
```go
// Use smaller quantized models for speed
models := map[string]string{
    "fastest":  "/path/to/model.Q3_K_M.gguf",    // Smallest, fastest
    "balanced": "/path/to/model.Q4_K_M.gguf",    // Good balance
    "quality":  "/path/to/model.Q5_K_M.gguf",    // Higher quality, slower
}

// Choose based on requirements
modelPath := models["balanced"]
```

3. **Context Management:**
```go
// Aggressive context truncation for speed
config := chanollama.ConversationConfig{
    MaxTokens:          1024,                          // Small context
    MaxMessages:        10,                            // Few messages
    TruncationStrategy: chanollama.TruncationOldest,
    PreserveSystem:     true,
}
```

#### "High Memory Usage"
```
Application consuming excessive RAM
```

**Memory Optimization:**
```go
func memoryOptimizedSetup() *chanollama.Client {
    config := chanollama.Config{
        ModelPath:   "/path/to/model.gguf",
        OfflineMode: true,
        
        // Reduce memory usage
        ContextSize: 2048,        // Smaller context
        MaxTokens:   256,         // Shorter responses
        NumThreads:  2,           // Fewer threads
    }
    
    client := chanollama.NewWithConfig(config.ModelPath, config)
    
    // Aggressive conversation limits
    convConfig := chanollama.ConversationConfig{
        MaxTokens:   800,         // Very conservative
        MaxMessages: 5,           // Minimal history
    }
    client.SetConversationConfig(convConfig)
    
    return client
}

// Memory monitoring
func monitorMemoryUsage(client *chanollama.Client) {
    monitor := chanollama.NewResourceMonitor()
    
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            memUsage, err := monitor.GetMemoryUsage()
            if err == nil {
                memGB := float64(memUsage) / (1024 * 1024 * 1024)
                log.Printf("Current memory usage: %.2f GB", memGB)
                
                if memGB > 8.0 { // Alert if over 8GB
                    log.Printf("HIGH MEMORY USAGE WARNING: %.2f GB", memGB)
                }
            }
        }
    }()
}
```

---

### Streaming Issues

#### "Streaming Stops or Hangs"
```
Streaming callback stops receiving tokens
```

**Debugging Streaming:**
```go
func debugStreaming(client *chanollama.Client, messages []chanollama.Message) {
    engine := chanollama.NewStreamingEngine(client)
    
    // Use channels for better control
    tokenCh, errorCh := engine.ChatStreamChanWithTimeout(messages, 30*time.Second)
    
    tokenCount := 0
    lastTokenTime := time.Now()
    
    for {
        select {
        case token, ok := <-tokenCh:
            if !ok {
                log.Printf("Streaming completed after %d tokens", tokenCount)
                return
            }
            
            tokenCount++
            fmt.Print(token)
            
            // Check for long gaps between tokens
            now := time.Now()
            gap := now.Sub(lastTokenTime)
            if gap > 5*time.Second {
                log.Printf("Long gap between tokens: %v", gap)
            }
            lastTokenTime = now
            
        case err := <-errorCh:
            if err != nil {
                log.Printf("Streaming error after %d tokens: %v", tokenCount, err)
                return
            }
            
        case <-time.After(10 * time.Second):
            log.Printf("No tokens received for 10 seconds (total: %d)", tokenCount)
            // Continue waiting or implement retry logic
        }
    }
}
```

#### "Streaming Performance Issues"
```
Tokens arriving slowly or irregularly
```

**Streaming Optimization:**
```go
func optimizedStreaming() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()
    
    // Configure for optimal streaming
    config := client.GetConfig()
    config.NumThreads = runtime.NumCPU()  // Max CPU for generation
    config.MaxTokens = 1024               // Reasonable length
    client.UpdateConfig(config)
    
    engine := chanollama.NewStreamingEngine(client)
    
    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain Go programming"),
    }
    
    // Use batching for smoother display
    tokenCh, errorCh := engine.ChatStreamChan(messages)
    batcher := chanollama.NewTokenBatcher(3, 100*time.Millisecond)
    batchCh := batcher.BatchTokens(tokenCh)
    
    for {
        select {
        case batch, ok := <-batchCh:
            if !ok {
                return
            }
            
            // Display batch smoothly
            for _, token := range batch {
                fmt.Print(token)
                time.Sleep(10 * time.Millisecond) // Smooth display
            }
            
        case err := <-errorCh:
            if err != nil {
                log.Printf("Streaming error: %v", err)
                return
            }
        }
    }
}
```

---

### Concurrency Issues

#### "Goroutine Leaks"
```
Warning: high number of goroutines detected
```

**Prevention:**
```go
func properGoroutineManagement() {
    client := chanollama.New("/path/to/model.gguf")
    defer client.Close()
    
    engine := chanollama.NewStreamingEngine(client)
    
    messages := []chanollama.Message{
        chanollama.NewUserMessage("Test message"),
    }
    
    // Always use context for cancellation
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    tokenCh, errorCh := engine.ChatStreamChan(messages)
    
    // Ensure cleanup with done channel
    done := make(chan bool, 1)
    
    go func() {
        defer func() { done <- true }()
        
        for {
            select {
            case token, ok := <-tokenCh:
                if !ok {
                    return
                }
                fmt.Print(token)
                
            case err := <-errorCh:
                if err != nil {
                    log.Printf("Error: %v", err)
                }
                return
                
            case <-ctx.Done():
                log.Printf("Context cancelled")
                return
            }
        }
    }()
    
    // Wait for completion or timeout
    select {
    case <-done:
        log.Printf("Streaming completed normally")
    case <-ctx.Done():
        log.Printf("Streaming timed out")
    }
}
```

#### "Race Conditions"
```
Concurrent access causing data races
```

**Thread-Safe Operations:**
```go
type SafeClient struct {
    client *chanollama.Client
    mu     sync.RWMutex
}

func NewSafeClient(modelPath string) *SafeClient {
    return &SafeClient{
        client: chanollama.New(modelPath),
    }
}

func (sc *SafeClient) Chat(messages []chanollama.Message) (*chanollama.Response, error) {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    
    return sc.client.Chat(messages)
}

func (sc *SafeClient) IsHealthy() bool {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    
    return sc.client.IsHealthy()
}

func (sc *SafeClient) Close() error {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    
    return sc.client.Close()
}

// Usage with concurrent access
func concurrentUsage() {
    safeClient := NewSafeClient("/path/to/model.gguf")
    defer safeClient.Close()
    
    var wg sync.WaitGroup
    
    // Multiple concurrent requests
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            messages := []chanollama.Message{
                chanollama.NewUserMessage(fmt.Sprintf("Question %d", id)),
            }
            
            response, err := safeClient.Chat(messages)
            if err != nil {
                log.Printf("Request %d failed: %v", id, err)
            } else {
                log.Printf("Request %d: %s", id, response.Content)
            }
        }(i)
    }
    
    wg.Wait()
}
```

---

## 🔍 Debugging Techniques

### Enable Verbose Logging

```go
func enableDebugLogging() {
    config := chanollama.DefaultOfflineConfig("/path/to/model.gguf")
    config.Verbose = true  // Enable llama.cpp verbose logging
    
    client := chanollama.NewWithConfig(config.ModelPath, config)
    defer client.Close()
    
    // Also enable conversation logging
    convConfig := chanollama.ConversationConfig{
        EnableLogging: true,
        MaxTokens:     2048,
    }
    client.SetConversationConfig(convConfig)
}
```

### Health Check Function

```go
func comprehensiveHealthCheck(client *chanollama.Client) error {
    // Basic health check
    if !client.IsHealthy() {
        return fmt.Errorf("client health check failed")
    }
    
    // Model info check
    modelInfo, err := client.GetModelInfo()
    if err != nil {
        return fmt.Errorf("failed to get model info: %w", err)
    }
    
    log.Printf("Model info: %+v", modelInfo)
    
    // Simple inference test
    testMessages := []chanollama.Message{
        chanollama.NewUserMessage("Test"),
    }
    
    start := time.Now()
    response, err := client.Chat(testMessages)
    duration := time.Since(start)
    
    if err != nil {
        return fmt.Errorf("inference test failed: %w", err)
    }
    
    if len(response.Content) == 0 {
        return fmt.Errorf("inference test returned empty response")
    }
    
    log.Printf("Health check passed: %d chars in %v", len(response.Content), duration)
    return nil
}
```

### Performance Profiling

```go
func profilePerformance(client *chanollama.Client) {
    messages := []chanollama.Message{
        chanollama.NewUserMessage("Explain machine learning in detail"),
    }
    
    // Memory profiling
    var m1, m2 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    start := time.Now()
    response, err := client.Chat(messages)
    duration := time.Since(start)
    
    runtime.ReadMemStats(&m2)
    
    if err != nil {
        log.Printf("Performance test failed: %v", err)
        return
    }
    
    // Calculate metrics
    memUsed := m2.Alloc - m1.Alloc
    tokensPerSecond := float64(response.Tokens) / duration.Seconds()
    
    log.Printf("Performance Metrics:")
    log.Printf("- Duration: %v", duration)
    log.Printf("- Tokens: %d", response.Tokens)
    log.Printf("- Tokens/sec: %.1f", tokensPerSecond)
    log.Printf("- Memory used: %d bytes", memUsed)
    log.Printf("- Response length: %d chars", len(response.Content))
}
```

---

## 🛠️ Diagnostic Tools

### Model Validator

```go
func validateModel(modelPath string) error {
    // Check file exists
    info, err := os.Stat(modelPath)
    if err != nil {
        return fmt.Errorf("cannot access model file: %w", err)
    }
    
    // Check file size
    sizeGB := float64(info.Size()) / (1024 * 1024 * 1024)
    log.Printf("Model size: %.2f GB", sizeGB)
    
    if sizeGB < 0.1 {
        return fmt.Errorf("model file too small (%.2f GB), possibly corrupted", sizeGB)
    }
    
    // Try to create client
    client := chanollama.New(modelPath)
    defer client.Close()
    
    // Initialize and test
    if err := client.Initialize(); err != nil {
        return fmt.Errorf("model initialization failed: %w", err)
    }
    
    if !client.IsHealthy() {
        return fmt.Errorf("model health check failed")
    }
    
    log.Printf("Model validation successful")
    return nil
}
```

### System Requirements Check

```go
func checkSystemRequirements(modelPath string) error {
    // Check Go version
    goVersion := runtime.Version()
    log.Printf("Go version: %s", goVersion)
    
    // Check CPU cores
    numCPU := runtime.NumCPU()
    log.Printf("CPU cores: %d", numCPU)
    
    if numCPU < 2 {
        return fmt.Errorf("minimum 2 CPU cores required, found %d", numCPU)
    }
    
    // Check memory
    monitor := chanollama.NewResourceMonitor()
    if !monitor.IsMemoryAvailable(2048) { // 2GB minimum
        return fmt.Errorf("insufficient memory (minimum 2GB required)")
    }
    
    // Check model file size vs available memory
    info, err := os.Stat(modelPath)
    if err != nil {
        return fmt.Errorf("cannot check model file: %w", err)
    }
    
    modelSizeMB := info.Size() / (1024 * 1024)
    recommendedMemMB := modelSizeMB * 2  // 2x model size recommended
    
    if !monitor.IsMemoryAvailable(uint64(recommendedMemMB)) {
        log.Printf("WARNING: Recommended memory (%d MB) may not be available", recommendedMemMB)
    }
    
    log.Printf("System requirements check passed")
    return nil
}
```

---

## 📋 Best Practices

### 1. **Error Handling**
- Always check `client.IsHealthy()` before critical operations
- Implement retry logic with exponential backoff
- Use specific error type checking with `errors.Is()`
- Log errors with sufficient context for debugging

### 2. **Resource Management**
- Always call `client.Close()` to free resources
- Use `defer client.Close()` immediately after creation
- Monitor memory usage in long-running applications
- Implement health checks for production systems

### 3. **Performance Optimization**
- Use appropriate model quantization for your hardware
- Configure `NumThreads` based on available CPU cores
- Set realistic `MaxTokens` to balance speed and completeness
- Enable conversation management to prevent context overflow

### 4. **Debugging**
- Enable verbose logging during development
- Use timeouts for all operations
- Implement comprehensive health checks
- Profile performance regularly

### 5. **Production Deployment**
- Validate models and system requirements before deployment
- Implement monitoring and alerting
- Use graceful shutdown procedures
- Test error scenarios and recovery procedures

---

## 🆘 Getting Help

### Collect Debug Information

When reporting issues, include:

```go
func collectDebugInfo(client *chanollama.Client) {
    log.Printf("=== chanollama Debug Information ===")
    
    // System info
    log.Printf("Go version: %s", runtime.Version())
    log.Printf("OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
    log.Printf("CPU cores: %d", runtime.NumCPU())
    
    // Client info
    if client != nil {
        config := client.GetConfig()
        log.Printf("Model path: %s", config.ModelPath)
        log.Printf("Context size: %d", config.ContextSize)
        log.Printf("Max tokens: %d", config.MaxTokens)
        log.Printf("Temperature: %f", config.Temperature)
        log.Printf("Num threads: %d", config.NumThreads)
        log.Printf("Client healthy: %v", client.IsHealthy())
        
        // Model info if available
        if modelInfo, err := client.GetModelInfo(); err == nil {
            log.Printf("Model context size: %d", modelInfo.ContextSize)
            log.Printf("Model batch size: %d", modelInfo.BatchSize)
        }
    }
    
    // Memory info
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    log.Printf("Allocated memory: %d MB", m.Alloc/(1024*1024))
    log.Printf("Total allocations: %d MB", m.TotalAlloc/(1024*1024))
    log.Printf("Sys memory: %d MB", m.Sys/(1024*1024))
    
    log.Printf("=== End Debug Information ===")
}
```

### Issue Template

When reporting bugs, provide:

1. **System Information**
   - Go version
   - Operating system and architecture
   - Available RAM and CPU cores

2. **Model Information**
   - Model file name and size
   - Where the model was downloaded from
   - File checksum if available

3. **Configuration**
   - chanollama configuration used
   - Any custom settings

4. **Error Details**
   - Complete error message
   - Stack trace if available
   - Steps to reproduce

5. **Expected vs Actual Behavior**
   - What you expected to happen
   - What actually happened

---

This troubleshooting guide covers the most common issues and solutions for chanollama. For additional help, see the [API Reference](./api-reference.md), [examples directory](../examples/), or create an issue in the project repository.