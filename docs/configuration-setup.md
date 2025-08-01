# ⚙️ Configuration and Setup Guide

Complete configuration guide for chanollama - covering all settings for optimal offline LLM inference.

## 🎯 Overview

chanollama provides extensive configuration options to optimize performance, control resource usage, and customize behavior for your specific use case. This guide covers all configuration aspects from basic setup to advanced production deployment.

---

## 🔧 Basic Configuration

### Default Configuration

The simplest setup uses default values:

```go
// Minimal setup - uses all defaults
client := chanollama.New("/path/to/model.gguf")
defer client.Close()

// Check current configuration
config := client.GetConfig()
fmt.Printf("Temperature: %f\n", config.Temperature)     // 0.7
fmt.Printf("Max Tokens: %d\n", config.MaxTokens)       // 2048
fmt.Printf("Context Size: %d\n", config.ContextSize)   // 4096
```

### Custom Configuration

Override defaults with custom values:

```go
config := chanollama.Config{
    // Model settings
    ModelPath:   "/path/to/model.gguf",
    OfflineMode: true,                    // Explicit offline mode
    
    // Generation parameters
    Temperature: 0.8,                     // Creativity level (0.0-1.0)
    MaxTokens:   1500,                    // Response length limit
    Timeout:     45 * time.Second,        // Request timeout
    
    // Hardware settings
    NumThreads:  8,                       // CPU threads for inference
    ContextSize: 8192,                    // Context window size
    UseGPU:      false,                   // GPU acceleration
    
    // Logging
    Verbose:     false,                   // Suppress llama.cpp logs
}

client := chanollama.NewWithConfig(config.ModelPath, config)
defer client.Close()
```

### Factory Configuration

Use predefined configurations:

```go
// Default offline configuration
config := chanollama.DefaultOfflineConfig("/path/to/model.gguf")

// Modify as needed
config.Temperature = 0.5      // More deterministic
config.MaxTokens = 1000       // Shorter responses
config.NumThreads = 4         // Conservative CPU usage

client := chanollama.NewWithConfig(config.ModelPath, config)
```

---

## 🏗️ Model Configuration

### Model Selection and Loading

Configure how models are loaded and managed:

```go
config := chanollama.Config{
    ModelPath:   "/path/to/model.gguf",
    OfflineMode: true,
    
    // Model loading settings
    Verbose:     false,                   // Suppress loading logs
    ContextSize: 4096,                    // Adjust for model capabilities
}

client := chanollama.NewWithConfig(config.ModelPath, config)
defer client.Close()

// Explicit initialization with error handling
if err := client.Initialize(); err != nil {
    log.Fatalf("Model initialization failed: %v", err)
}

// Verify model is ready
if !client.IsHealthy() {
    log.Fatal("Model health check failed")
}

// Get model information
modelInfo, err := client.GetModelInfo()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Model loaded: %s\n", modelInfo.ModelPath)
fmt.Printf("Context size: %d\n", modelInfo.ContextSize)
fmt.Printf("Batch size: %d\n", modelInfo.BatchSize)
```

### Model Path Resolution

Handle different model path scenarios:

```go
func createClientWithPath(modelPath string) (*chanollama.Client, error) {
    // Validate model file exists
    if _, err := os.Stat(modelPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("model file not found: %s", modelPath)
    }
    
    // Get absolute path
    absPath, err := filepath.Abs(modelPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve model path: %w", err)
    }
    
    config := chanollama.DefaultOfflineConfig(absPath)
    
    // Auto-detect context size based on model filename
    if strings.Contains(filepath.Base(absPath), "7b") {
        config.ContextSize = 4096      // 7B models
    } else if strings.Contains(filepath.Base(absPath), "13b") {
        config.ContextSize = 4096      // 13B models  
    } else if strings.Contains(filepath.Base(absPath), "70b") {
        config.ContextSize = 2048      // Reduce for 70B models
    }
    
    return chanollama.NewWithConfig(absPath, config), nil
}
```

---

## 🎛️ Generation Parameters

### Temperature Settings

Control creativity and randomness:

```go
configs := map[string]chanollama.Config{
    "creative": {
        ModelPath:   "/path/to/model.gguf",
        Temperature: 0.9,              // Very creative, varied responses
        OfflineMode: true,
    },
    "balanced": {
        ModelPath:   "/path/to/model.gguf", 
        Temperature: 0.7,              // Default balanced setting
        OfflineMode: true,
    },
    "focused": {
        ModelPath:   "/path/to/model.gguf",
        Temperature: 0.3,              // More deterministic
        OfflineMode: true,
    },
    "deterministic": {
        ModelPath:   "/path/to/model.gguf",
        Temperature: 0.1,              // Very consistent responses
        OfflineMode: true,
    },
}

// Use appropriate configuration for your use case
client := chanollama.NewWithConfig(configs["focused"].ModelPath, configs["focused"])
defer client.Close()
```

### Response Length Control

Configure response length limits:

```go
configs := map[string]chanollama.Config{
    "short": {
        ModelPath:   "/path/to/model.gguf",
        MaxTokens:   256,              // Brief responses
        OfflineMode: true,
    },
    "medium": {
        ModelPath:   "/path/to/model.gguf",
        MaxTokens:   1024,             // Standard responses
        OfflineMode: true,
    },
    "long": {
        ModelPath:   "/path/to/model.gguf",
        MaxTokens:   2048,             // Detailed responses
        OfflineMode: true,
    },
    "extended": {
        ModelPath:   "/path/to/model.gguf",
        MaxTokens:   4096,             // Very detailed responses
        OfflineMode: true,
    },
}

// Dynamic configuration based on query type
func getConfigForQuery(query string, modelPath string) chanollama.Config {
    baseConfig := chanollama.DefaultOfflineConfig(modelPath)
    
    queryLower := strings.ToLower(query)
    
    if strings.Contains(queryLower, "brief") || strings.Contains(queryLower, "summary") {
        baseConfig.MaxTokens = 256
        baseConfig.Temperature = 0.3
    } else if strings.Contains(queryLower, "detailed") || strings.Contains(queryLower, "explain") {
        baseConfig.MaxTokens = 2048
        baseConfig.Temperature = 0.7
    } else if strings.Contains(queryLower, "creative") || strings.Contains(queryLower, "story") {
        baseConfig.MaxTokens = 1024
        baseConfig.Temperature = 0.9
    }
    
    return baseConfig
}
```

### Timeout Configuration

Configure request timeouts:

```go
configs := map[string]chanollama.Config{
    "fast": {
        ModelPath:   "/path/to/model.gguf",
        Timeout:     15 * time.Second,  // Quick responses only
        MaxTokens:   512,               // Shorter to meet timeout
        OfflineMode: true,
    },
    "standard": {
        ModelPath:   "/path/to/model.gguf",
        Timeout:     30 * time.Second,  // Default timeout
        MaxTokens:   1024,
        OfflineMode: true,
    },
    "patient": {
        ModelPath:   "/path/to/model.gguf",
        Timeout:     120 * time.Second, // Allow longer processing
        MaxTokens:   4096,
        OfflineMode: true,
    },
}

// Use context for additional timeout control
func chatWithTimeout(client *chanollama.Client, messages []chanollama.Message, timeout time.Duration) (*chanollama.Response, error) {
    // Note: This example shows the pattern, but chanollama doesn't currently accept context
    // The timeout is configured at the client level
    
    start := time.Now()
    response, err := client.Chat(messages)
    elapsed := time.Since(start)
    
    if elapsed > timeout {
        return nil, fmt.Errorf("operation exceeded timeout of %v (took %v)", timeout, elapsed)
    }
    
    return response, err
}
```

---

## 🖥️ Hardware Configuration

### CPU Thread Optimization

Configure CPU usage for optimal performance:

```go
func getOptimalThreadConfig() chanollama.Config {
    numCPU := runtime.NumCPU()
    
    config := chanollama.Config{
        ModelPath:   "/path/to/model.gguf",
        OfflineMode: true,
    }
    
    // Configure threads based on system capabilities
    switch {
    case numCPU >= 16:
        config.NumThreads = numCPU - 2    // Leave 2 cores for system
    case numCPU >= 8:
        config.NumThreads = numCPU - 1    // Leave 1 core for system
    case numCPU >= 4:
        config.NumThreads = numCPU        // Use all cores
    default:
        config.NumThreads = 2             // Minimum reasonable threads
    }
    
    fmt.Printf("Configured %d threads (system has %d CPUs)\n", config.NumThreads, numCPU)
    return config
}

// Environment-specific configuration
func getConfigForEnvironment() chanollama.Config {
    config := chanollama.DefaultOfflineConfig("/path/to/model.gguf")
    
    // Detect environment
    if os.Getenv("CONTAINER_ENV") == "true" {
        // Container environment - be conservative
        config.NumThreads = 2
        config.ContextSize = 2048
        config.MaxTokens = 512
    } else if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
        // Apple Silicon - optimize for efficiency cores
        config.NumThreads = 6
        config.ContextSize = 4096
    } else {
        // Default server environment
        config.NumThreads = runtime.NumCPU()
        config.ContextSize = 4096
    }
    
    return config
}
```

### Memory Configuration

Configure memory usage and context size:

```go
func getMemoryOptimizedConfig(modelPath string) chanollama.Config {
    config := chanollama.DefaultOfflineConfig(modelPath)
    
    // Get available memory
    monitor := chanollama.NewResourceMonitor()
    availableMemMB, err := monitor.GetMemoryUsage()
    if err != nil {
        log.Printf("Could not determine memory usage, using defaults")
        return config
    }
    
    availableMemGB := availableMemMB / (1024 * 1024)
    
    fmt.Printf("Available memory: %d GB\n", availableMemGB)
    
    // Adjust context size based on available memory
    switch {
    case availableMemGB >= 32:
        config.ContextSize = 8192      // Large context for high-memory systems
        config.MaxTokens = 2048
    case availableMemGB >= 16:
        config.ContextSize = 4096      // Standard context
        config.MaxTokens = 1024
    case availableMemGB >= 8:
        config.ContextSize = 2048      // Reduced context for medium memory
        config.MaxTokens = 512
    default:
        config.ContextSize = 1024      // Minimal context for low memory
        config.MaxTokens = 256
    }
    
    fmt.Printf("Configured context size: %d, max tokens: %d\n", 
        config.ContextSize, config.MaxTokens)
    
    return config
}

// Memory monitoring during inference
func monitoredInference(client *chanollama.Client, messages []chanollama.Message) (*chanollama.Response, error) {
    monitor := chanollama.NewResourceMonitor()
    
    // Check memory before inference
    beforeMem, _ := monitor.GetMemoryUsage()
    
    // Perform inference
    response, err := client.Chat(messages)
    
    // Check memory after inference
    afterMem, _ := monitor.GetMemoryUsage()
    
    memUsedMB := (afterMem - beforeMem) / (1024 * 1024)
    fmt.Printf("Memory used for inference: %d MB\n", memUsedMB)
    
    return response, err
}
```

### GPU Configuration

Configure GPU acceleration (when supported):

```go
func getGPUConfig(modelPath string) chanollama.Config {
    config := chanollama.DefaultOfflineConfig(modelPath)
    
    // GPU configuration (placeholder for future GPU support)
    config.UseGPU = false  // Currently not supported in chanollama
    
    // When GPU support is added, this would detect GPU availability
    // if hasCompatibleGPU() {
    //     config.UseGPU = true
    //     config.NumThreads = 4  // Fewer CPU threads when using GPU
    // }
    
    return config
}

// Future GPU detection function
func hasCompatibleGPU() bool {
    // This would detect CUDA, Metal, or other GPU acceleration
    // Currently returns false as GPU support is not implemented
    return false
}
```

---

## 💬 Conversation Management

### Basic Conversation Configuration

Configure automatic conversation management:

```go
func configureConversationManagement(client *chanollama.Client) {
    config := chanollama.ConversationConfig{
        // Token limits
        MaxTokens:          3200,                          // 80% of context size
        MaxMessages:        50,                            // Fallback message limit
        
        // Truncation behavior
        TruncationStrategy: chanollama.TruncationOldest,   // Remove oldest messages
        PreserveSystem:     true,                          // Always keep system messages
        
        // Auto-detection
        AutoDetectLimits:   true,                          // Auto-detect from model
        EnableLogging:      false,                         // Log truncation events
        
        // Summarization (if enabled)
        SummaryPrompt: "Summarize the following conversation concisely:",
    }
    
    client.SetConversationConfig(config)
    
    // Test auto-detection
    if err := client.AutoDetectConversationLimits(); err != nil {
        log.Printf("Auto-detection failed: %v", err)
    }
}
```

### Advanced Conversation Strategies

Configure different truncation strategies:

```go
func setupConversationStrategies() {
    modelPath := "/path/to/model.gguf"
    
    strategies := map[string]chanollama.ConversationConfig{
        "oldest_first": {
            MaxTokens:          2048,
            TruncationStrategy: chanollama.TruncationOldest,    // Remove oldest messages
            PreserveSystem:     true,
        },
        "newest_first": {
            MaxTokens:          2048,
            TruncationStrategy: chanollama.TruncationNewest,    // Remove newest messages
            PreserveSystem:     true,
        },
        "middle_out": {
            MaxTokens:          2048,
            TruncationStrategy: chanollama.TruncationMiddle,    // Keep first and last
            PreserveSystem:     true,
        },
        "summarize": {
            MaxTokens:          2048,
            TruncationStrategy: chanollama.TruncationSummarize, // Summarize old messages
            PreserveSystem:     true,
            SummaryPrompt:      "Provide a brief summary of the following conversation:",
        },
    }
    
    // Create clients with different strategies
    for name, config := range strategies {
        client := chanollama.New(modelPath)
        client.SetConversationConfig(config)
        
        fmt.Printf("Client '%s' configured with %s strategy\n", name, config.TruncationStrategy)
        // Use client for specific use case...
        client.Close()
    }
}
```

### Dynamic Conversation Management

Adjust conversation settings based on context:

```go
func dynamicConversationManagement(client *chanollama.Client, messages []chanollama.Message) {
    // Analyze conversation complexity
    totalTokens := 0
    for _, msg := range messages {
        totalTokens += len(msg.Content) / 4  // Rough token estimate
    }
    
    // Get current configuration
    manager := client.GetConversationManager()
    currentConfig := manager.GetConfig()
    
    // Adjust based on conversation size
    if totalTokens > 3000 {
        // Large conversation - use aggressive truncation
        currentConfig.MaxTokens = 2048
        currentConfig.TruncationStrategy = chanollama.TruncationOldest
        currentConfig.PreserveSystem = true
    } else if totalTokens > 1500 {
        // Medium conversation - moderate truncation
        currentConfig.MaxTokens = 2500
        currentConfig.TruncationStrategy = chanollama.TruncationMiddle
    } else {
        // Small conversation - minimal truncation
        currentConfig.MaxTokens = 3500
        currentConfig.TruncationStrategy = chanollama.TruncationOldest
    }
    
    // Apply updated configuration
    client.SetConversationConfig(currentConfig)
    
    fmt.Printf("Adjusted conversation limits for %d tokens\n", totalTokens)
}
```

---

## 📊 Logging and Monitoring

### Logging Configuration

Configure logging levels and output:

```go
func setupLogging() chanollama.Config {
    config := chanollama.DefaultOfflineConfig("/path/to/model.gguf")
    
    // Environment-based logging
    logLevel := os.Getenv("LOG_LEVEL")
    switch strings.ToLower(logLevel) {
    case "debug", "verbose":
        config.Verbose = true
    case "info", "normal":
        config.Verbose = false
    case "quiet", "silent":
        config.Verbose = false
    default:
        // Default to quiet in production
        config.Verbose = false
    }
    
    // Custom logging for conversation management
    convConfig := chanollama.ConversationConfig{
        EnableLogging: config.Verbose,
        MaxTokens:     3200,
        PreserveSystem: true,
    }
    
    client := chanollama.NewWithConfig(config.ModelPath, config)
    client.SetConversationConfig(convConfig)
    
    return config
}

// Custom logging wrapper
type LoggingClient struct {
    client *chanollama.Client
    logger *log.Logger
}

func NewLoggingClient(modelPath string, logger *log.Logger) *LoggingClient {
    config := chanollama.DefaultOfflineConfig(modelPath)
    config.Verbose = false  // We'll handle logging ourselves
    
    return &LoggingClient{
        client: chanollama.NewWithConfig(modelPath, config),
        logger: logger,
    }
}

func (lc *LoggingClient) Chat(messages []chanollama.Message) (*chanollama.Response, error) {
    start := time.Now()
    lc.logger.Printf("Starting chat with %d messages", len(messages))
    
    response, err := lc.client.Chat(messages)
    
    duration := time.Since(start)
    if err != nil {
        lc.logger.Printf("Chat failed after %v: %v", duration, err)
    } else {
        lc.logger.Printf("Chat completed in %v: %d tokens, %d chars", 
            duration, response.Tokens, len(response.Content))
    }
    
    return response, err
}

func (lc *LoggingClient) Close() error {
    lc.logger.Printf("Closing logging client")
    return lc.client.Close()
}
```

### Performance Monitoring

Monitor and track performance metrics:

```go
type PerformanceMonitor struct {
    requestCount    int64
    totalDuration   time.Duration
    errorCount      int64
    averageTokens   float64
    mu              sync.Mutex
}

func NewPerformanceMonitor() *PerformanceMonitor {
    return &PerformanceMonitor{}
}

func (pm *PerformanceMonitor) RecordRequest(duration time.Duration, tokens int, err error) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    pm.requestCount++
    pm.totalDuration += duration
    
    if err != nil {
        pm.errorCount++
    } else {
        // Update running average of tokens
        pm.averageTokens = (pm.averageTokens*float64(pm.requestCount-1) + float64(tokens)) / float64(pm.requestCount)
    }
}

func (pm *PerformanceMonitor) GetStats() (avgDuration time.Duration, errorRate float64, avgTokens float64, totalRequests int64) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    if pm.requestCount > 0 {
        avgDuration = pm.totalDuration / time.Duration(pm.requestCount)
        errorRate = float64(pm.errorCount) / float64(pm.requestCount)
    }
    
    return avgDuration, errorRate, pm.averageTokens, pm.requestCount
}

// Monitored client wrapper
type MonitoredClient struct {
    client  *chanollama.Client
    monitor *PerformanceMonitor
}

func NewMonitoredClient(modelPath string) *MonitoredClient {
    return &MonitoredClient{
        client:  chanollama.New(modelPath),
        monitor: NewPerformanceMonitor(),
    }
}

func (mc *MonitoredClient) Chat(messages []chanollama.Message) (*chanollama.Response, error) {
    start := time.Now()
    response, err := mc.client.Chat(messages)
    duration := time.Since(start)
    
    tokens := 0
    if response != nil {
        tokens = response.Tokens
    }
    
    mc.monitor.RecordRequest(duration, tokens, err)
    return response, err
}

func (mc *MonitoredClient) GetPerformanceStats() {
    avgDuration, errorRate, avgTokens, totalRequests := mc.monitor.GetStats()
    
    fmt.Printf("Performance Statistics:\n")
    fmt.Printf("- Total requests: %d\n", totalRequests)
    fmt.Printf("- Average duration: %v\n", avgDuration)
    fmt.Printf("- Error rate: %.2f%%\n", errorRate*100)
    fmt.Printf("- Average tokens: %.1f\n", avgTokens)
    
    if avgDuration > 0 {
        tokensPerSecond := avgTokens / avgDuration.Seconds()
        fmt.Printf("- Tokens per second: %.1f\n", tokensPerSecond)
    }
}

func (mc *MonitoredClient) Close() error {
    return mc.client.Close()
}
```

---

## 🏭 Production Configuration

### Production-Ready Setup

Complete production configuration:

```go
func createProductionClient(modelPath string) (*chanollama.Client, error) {
    // Validate environment
    if err := validateProductionEnvironment(modelPath); err != nil {
        return nil, fmt.Errorf("production validation failed: %w", err)
    }
    
    // Production configuration
    config := chanollama.Config{
        ModelPath:   modelPath,
        OfflineMode: true,
        
        // Conservative generation settings
        Temperature: 0.3,                     // More deterministic
        MaxTokens:   1024,                    // Reasonable response length
        Timeout:     30 * time.Second,        // Standard timeout
        
        // Hardware optimization
        NumThreads:  runtime.NumCPU() / 2,    // Leave resources for other processes
        ContextSize: 4096,                    // Standard context size
        UseGPU:      false,                   // Disable GPU for stability
        
        // Logging
        Verbose:     false,                   // Clean production logs
    }
    
    client := chanollama.NewWithConfig(modelPath, config)
    
    // Production conversation management
    convConfig := chanollama.ConversationConfig{
        MaxTokens:          3200,                          // 80% of context
        MaxMessages:        30,                            // Conservative limit
        TruncationStrategy: chanollama.TruncationOldest,   // Predictable truncation
        PreserveSystem:     true,                          // Keep system prompts
        AutoDetectLimits:   true,                          // Auto-optimize
        EnableLogging:      false,                         // Disable conv logging
    }
    
    client.SetConversationConfig(convConfig)
    
    // Pre-initialize for consistent performance
    if err := client.Initialize(); err != nil {
        client.Close()
        return nil, fmt.Errorf("client initialization failed: %w", err)
    }
    
    // Health check
    if !client.IsHealthy() {
        client.Close()
        return nil, fmt.Errorf("client health check failed")
    }
    
    return client, nil
}

func validateProductionEnvironment(modelPath string) error {
    // Check model file
    if _, err := os.Stat(modelPath); os.IsNotExist(err) {
        return fmt.Errorf("model file not found: %s", modelPath)
    }
    
    // Check available memory
    monitor := chanollama.NewResourceMonitor()
    if !monitor.IsMemoryAvailable(4096) { // 4GB minimum
        return fmt.Errorf("insufficient memory for production use")
    }
    
    // Check CPU cores
    if runtime.NumCPU() < 2 {
        return fmt.Errorf("minimum 2 CPU cores required for production")
    }
    
    return nil
}
```

### Environment-Specific Configurations

Handle different deployment environments:

```go
func getEnvironmentConfig(modelPath string) chanollama.Config {
    env := os.Getenv("DEPLOYMENT_ENV")
    
    switch env {
    case "development":
        return getDevelopmentConfig(modelPath)
    case "staging":
        return getStagingConfig(modelPath)
    case "production":
        return getProductionConfig(modelPath)
    default:
        log.Printf("Unknown environment '%s', using development config", env)
        return getDevelopmentConfig(modelPath)
    }
}

func getDevelopmentConfig(modelPath string) chanollama.Config {
    return chanollama.Config{
        ModelPath:   modelPath,
        OfflineMode: true,
        Temperature: 0.7,                     // Balanced for testing
        MaxTokens:   2048,                    // Longer responses for dev
        Timeout:     60 * time.Second,        // Longer timeout for debugging
        NumThreads:  runtime.NumCPU(),        // Use all CPU for dev
        ContextSize: 4096,
        Verbose:     true,                    // Verbose logging for debugging
    }
}

func getStagingConfig(modelPath string) chanollama.Config {
    return chanollama.Config{
        ModelPath:   modelPath,
        OfflineMode: true,
        Temperature: 0.5,                     // More consistent for testing
        MaxTokens:   1024,                    // Standard responses
        Timeout:     45 * time.Second,        // Moderate timeout
        NumThreads:  runtime.NumCPU() - 1,    // Leave 1 core for monitoring
        ContextSize: 4096,
        Verbose:     false,                   // Cleaner staging logs
    }
}

func getProductionConfig(modelPath string) chanollama.Config {
    return chanollama.Config{
        ModelPath:   modelPath,
        OfflineMode: true,
        Temperature: 0.3,                     // Most deterministic
        MaxTokens:   1024,                    // Predictable response length
        Timeout:     30 * time.Second,        // Fast response requirement
        NumThreads:  runtime.NumCPU() / 2,    // Conservative CPU usage
        ContextSize: 4096,
        Verbose:     false,                   // Minimal production logging
    }
}
```

---

## 🎯 Configuration Best Practices

### 1. **Model Selection**
- Use Q4_K_M quantization for balance of quality and size
- Consider Q5_K_M for better quality at larger size
- Use Q3_K_M only for extremely memory-constrained environments

### 2. **Hardware Optimization**
- Set `NumThreads` to 50-75% of available CPU cores
- Adjust `ContextSize` based on available memory
- Monitor memory usage and adjust accordingly

### 3. **Generation Parameters**
- Use `Temperature` 0.1-0.3 for deterministic responses
- Use `Temperature` 0.7-0.9 for creative responses
- Set `MaxTokens` based on your use case requirements

### 4. **Production Deployment**
- Always validate configuration before deployment
- Implement health checks and monitoring
- Use conservative settings for stability
- Log important metrics for troubleshooting

### 5. **Resource Management**
- Monitor memory usage in long-running applications
- Implement graceful shutdown procedures
- Use connection pooling for high-concurrency scenarios
- Consider horizontal scaling for very high loads

---

This configuration guide covers all aspects of setting up and optimizing chanollama. For specific use cases and advanced patterns, see the [API Reference](./api-reference.md) and [examples directory](../examples/).