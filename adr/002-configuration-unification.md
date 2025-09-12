# ADR-002: Configuration Structure Unification

## Status
**Proposed** - Awaiting review and implementation

## Context

chanollama currently has three overlapping configuration structures that create confusion, maintenance burden, and subtle bugs:

### Current Configuration Structures

#### 1. `Config` (types.go) - Primary Configuration
```go
type Config struct {
    Model       string        // Model identifier
    ModelPath   string        // Path to GGUF model file
    MaxTokens   int          // Maximum tokens to generate
    Temperature float64      // Sampling temperature
    UseGPU      bool         // Enable GPU acceleration
    NumThreads  int          // CPU thread count
    ContextSize int          // Model context window
    Timeout     time.Duration // Operation timeout
    OfflineMode bool         // Always true, deprecated
    Verbose     bool         // Enable verbose logging
}
```

#### 2. `EngineConfig` (engine.go) - Engine-Specific
```go
type EngineConfig struct {
    RequestChannelSize   int           // Request queue size
    StreamingChannelSize int           // Streaming buffer size
    OperationTimeout     time.Duration // Engine operation timeout
    ShutdownTimeout      time.Duration // Graceful shutdown timeout
    MaxConcurrentOps     int           // Concurrent operation limit
    MetricsInterval      time.Duration // Metrics collection interval
}
```

#### 3. `ConversationConfig` (types.go) - Conversation Management
```go
type ConversationConfig struct {
    MaxTokens            int                // Token limit for conversation
    MaxMessages          int                // Message count limit
    BatchSize            int                // Model batch size
    TruncationStrategy   TruncationStrategy // How to truncate
    PreserveSystemMessages bool             // Keep system messages
    AutoDetectLimits     bool               // Auto-detect model limits
    EnableLogging        bool               // Enable conversation logging
    ContextPadding       int                // Context padding size
}
```

### Critical Problems Identified

#### 1. Field Name Collisions
- **`MaxTokens` collision**: `Config.MaxTokens` (generation limit) vs `ConversationConfig.MaxTokens` (conversation limit)
- **Different semantics, same name**: Leads to silent bugs and developer confusion

#### 2. Timeout Ambiguity
- **`Config.Timeout`** vs **`EngineConfig.OperationTimeout`**: Which takes precedence?
- **No clear precedence rules**: Undefined behavior when both are set

#### 3. Context/Batch Confusion
- **`Config.ContextSize`** vs **`ConversationConfig.BatchSize`**: Both control context windows
- **Inconsistent naming**: Same concept, different field names

#### 4. Usage Complexity
- **15+ constructors** accepting different config combinations
- **No precedence documentation**: Developers must guess which config wins
- **Maintenance burden**: Changes require updates across multiple structures

### Constructor Proliferation
```go
// Current complexity
DefaultConfig(model)
DefaultEngineConfig()
DefaultConversationConfig()

NewWithConfig(path, Config)
NewAdvancedEngineWithConfig(Config, EngineConfig)
NewConversationManagerWithConfig(ConversationConfig)

// Different config combinations throughout codebase
Client{config: Config, conversationManager: ConversationManager{config: ConversationConfig}}
```

### Impact Assessment
- **Severity**: Medium-High - Creates subtle bugs and maintenance burden
- **Bug Risk**: Silent field conflicts cause unexpected behavior
- **Developer Confusion**: Must understand 3 different config systems
- **Maintenance Cost**: 3x configuration documentation and testing overhead

## Decision

**Unify all configuration into a single, comprehensive `Config` structure** with clear field semantics and no naming conflicts.

### Proposed Unified Configuration

```go
type Config struct {
    // Model Configuration (Core - Required)
    ModelPath   string `json:"model_path"`   // Path to GGUF model file
    Model       string `json:"model"`       // Model identifier (legacy support)
    
    // Generation Parameters (Core)
    MaxTokens   int     `json:"max_tokens"`  // Maximum tokens to generate per request
    Temperature float64 `json:"temperature"` // Sampling temperature (0.0-1.0)
    ContextSize int     `json:"context_size"` // Model context window size
    
    // Hardware Configuration
    UseGPU      bool `json:"use_gpu"`     // Enable GPU acceleration
    NumThreads  int  `json:"num_threads"` // CPU thread count (0 = auto)
    
    // Operation Timeouts (Resolved Naming)
    RequestTimeout  time.Duration `json:"request_timeout"`  // Individual request timeout
    OperationTimeout time.Duration `json:"operation_timeout"` // Engine operation timeout  
    ShutdownTimeout time.Duration `json:"shutdown_timeout"` // Graceful shutdown timeout
    
    // Conversation Management (Clear Semantics)
    ConversationMaxTokens   int                `json:"conversation_max_tokens"`   // Total conversation token limit
    ConversationMaxMessages int                `json:"conversation_max_messages"` // Total message count limit
    TruncationStrategy      TruncationStrategy `json:"truncation_strategy"`       // How to truncate when limits exceeded
    PreserveSystemMessages  bool               `json:"preserve_system_messages"`  // Keep system messages during truncation
    AutoDetectLimits        bool               `json:"auto_detect_limits"`        // Auto-detect model limits from model file
    
    // Engine Performance (Advanced)
    RequestChannelSize   int           `json:"request_channel_size"`   // Request queue buffer size
    StreamingChannelSize int           `json:"streaming_channel_size"` // Streaming response buffer size
    MaxConcurrentOps     int           `json:"max_concurrent_ops"`     // Maximum concurrent operations
    MetricsInterval      time.Duration `json:"metrics_interval"`       // Metrics collection frequency
    
    // Logging & Debugging
    Verbose             bool `json:"verbose"`              // Enable verbose output
    EnableLogging       bool `json:"enable_logging"`       // Enable conversation logging
    
    // Legacy/Deprecated (Phase Out)
    OfflineMode bool `json:"offline_mode"` // Always true, will be removed
}
```

### Field Semantic Resolution

| Problem | Current State | Unified Solution | Rationale |
|---------|---------------|------------------|-----------|
| `MaxTokens` collision | Config.MaxTokens vs ConversationConfig.MaxTokens | `MaxTokens` (per-request) + `ConversationMaxTokens` (total limit) | Clear distinction between per-request generation limit and total conversation limit |
| Timeout ambiguity | Config.Timeout vs EngineConfig.OperationTimeout | `RequestTimeout` + `OperationTimeout` + `ShutdownTimeout` | Specific timeouts for specific operations |
| Context confusion | Config.ContextSize vs ConversationConfig.BatchSize | `ContextSize` (unified) | Context is the standard LLM terminology |

### Constructor Simplification

```go
// Single configuration approach
func NewClient(modelPath string, opts ...ConfigOption) *Client
func NewAdvancedEngine(modelPath string, opts ...ConfigOption) *AdvancedEngine

// Functional options for clean configuration
type ConfigOption func(*Config)

func WithMaxTokens(tokens int) ConfigOption { 
    return func(c *Config) { c.MaxTokens = tokens } 
}

func WithConversationLimits(maxTokens, maxMessages int) ConfigOption {
    return func(c *Config) { 
        c.ConversationMaxTokens = maxTokens
        c.ConversationMaxMessages = maxMessages
    }
}

func WithTimeouts(request, operation, shutdown time.Duration) ConfigOption {
    return func(c *Config) {
        c.RequestTimeout = request
        c.OperationTimeout = operation  
        c.ShutdownTimeout = shutdown
    }
}

func WithEnginePerformance(reqChanSize, streamChanSize, maxOps int) ConfigOption {
    return func(c *Config) {
        c.RequestChannelSize = reqChanSize
        c.StreamingChannelSize = streamChanSize
        c.MaxConcurrentOps = maxOps
    }
}
```

### Usage Examples

#### Before (Current Confusion)
```go
// Developer confusion: which MaxTokens applies?
config := Config{MaxTokens: 1024}
engineConfig := EngineConfig{OperationTimeout: 30 * time.Second}
convConfig := ConversationConfig{MaxTokens: 2048} // Collision!

client := NewWithConfig(modelPath, config)
// How do I set engine and conversation config? Unclear...
```

#### After (Clear and Unified)
```go
// Clear, no confusion
client := NewClient("/path/to/model.gguf",
    WithMaxTokens(1024),                           // Per-request generation limit
    WithConversationLimits(4000, 20),              // Total conversation: 4K tokens, 20 messages
    WithTimeouts(5*time.Second, 30*time.Second, 10*time.Second), // Request, operation, shutdown
    WithEnginePerformance(100, 50, 5),             // Channel sizes and concurrency
    WithGPU(true),
    WithVerbose(true),
)
```

## Implementation Plan

### Phase 1: Add Unified Config (Non-Breaking)
1. **Add new unified `Config` struct** alongside existing structures
2. **Implement functional options** for all configuration aspects
3. **Create constructor overloads** that accept unified config
4. **Add comprehensive validation** for config field combinations
5. **Maintain backward compatibility** with existing configs

### Phase 2: Internal Migration
1. **Convert internal code** to use unified config exclusively
2. **Add adapter functions** to convert old configs to new format
3. **Update default config functions** to return unified structure
4. **Comprehensive testing** of config conversion and validation

### Phase 3: Deprecation Warnings
1. **Add deprecation warnings** to old config structures
2. **Update examples** to use unified config approach
3. **Create migration guide** with clear before/after examples
4. **Document new configuration semantics** comprehensively

### Phase 4: Clean Removal (Future Version)
1. **Remove deprecated config structures** after 2-3 versions
2. **Clean up adapter functions** and conversion code
3. **Simplify constructor signatures** to use only unified config
4. **Update documentation** to reflect simplified configuration model

## Benefits

### Developer Experience
- **Single Mental Model**: One configuration structure to learn and understand
- **Clear Semantics**: No field name collisions or ambiguous meanings
- **IDE Discoverability**: All configuration options visible in single structure
- **Validation**: Compile-time and runtime validation of config combinations
- **Documentation**: One comprehensive configuration reference

### Code Quality
- **Single Source of Truth**: All configuration in one place
- **Type Safety**: Compile-time checking of field usage
- **Reduced Complexity**: Eliminate config conversion and merging logic
- **Better Testing**: Test configuration as unified system vs separate concerns
- **Maintainability**: Changes require updates in only one place

### Migration Benefits
- **Backward Compatibility**: Existing code works during transition
- **Gradual Migration**: Convert components one at a time
- **Clear Mapping**: Old configs map directly to unified structure
- **No Feature Loss**: All current functionality preserved

## Risks & Mitigations

### Risk: Large Configuration Structure
- **Mitigation**: Use functional options for clean, focused configuration
- **Grouping**: Organize fields logically with clear documentation sections

### Risk: Breaking Changes During Migration
- **Mitigation**: Maintain full backward compatibility during transition
- **Adapters**: Automatic conversion from old configs to new structure

### Risk: Configuration Validation Complexity
- **Mitigation**: Implement comprehensive validation with clear error messages
- **Testing**: Extensive testing of valid and invalid configuration combinations

### Risk: Documentation Overhead
- **Mitigation**: Generate documentation from code comments and examples
- **Structure**: Organize docs by use case rather than field-by-field

## Validation Strategy

### Configuration Validation Rules
```go
func (c *Config) Validate() error {
    // Required fields
    if c.ModelPath == "" && c.Model == "" {
        return errors.New("either ModelPath or Model must be specified")
    }
    
    // Logical consistency
    if c.ConversationMaxTokens > 0 && c.ContextSize > 0 && c.ConversationMaxTokens > c.ContextSize {
        return errors.New("ConversationMaxTokens cannot exceed ContextSize")
    }
    
    // Range validation
    if c.Temperature < 0 || c.Temperature > 1 {
        return errors.New("Temperature must be between 0 and 1")
    }
    
    // Timeout validation
    if c.RequestTimeout > c.OperationTimeout {
        return errors.New("RequestTimeout cannot exceed OperationTimeout")
    }
    
    return nil
}
```

## Success Metrics

### Code Quality Metrics
- **Configuration Structures**: 3 → 1 (66% reduction)
- **Constructor Functions**: 15+ → ~3 (80% reduction)
- **Config-Related Code**: Measure LOC reduction in config handling
- **Test Coverage**: Maintain >90% coverage during migration

### Developer Experience Metrics  
- **Configuration Questions**: Reduction in support questions about config conflicts
- **Documentation**: Single config reference vs 3 separate sections
- **Example Clarity**: Before/after comparison of configuration examples

### Maintenance Metrics
- **Bug Rate**: Reduction in config-related bugs and issues
- **Change Impact**: Time to implement config changes (should decrease significantly)
- **Testing Complexity**: Reduction in config combination testing matrix

## Examples

### Migration Examples

#### Current: E-commerce Application
```go
// Before: Multiple configs, unclear precedence
config := Config{
    ModelPath:   "/models/llama2-7b.gguf",
    MaxTokens:   512,  // Generation limit
    Temperature: 0.7,
    UseGPU:     true,
}

convConfig := ConversationConfig{
    MaxTokens:   2048,  // Conversation limit - COLLISION!
    MaxMessages: 20,
}

client := NewWithConfig(config.ModelPath, config)
// How to apply conversation config? Unclear!
```

#### After: Unified, Clear Configuration  
```go
// After: Clear, unified, no confusion
client := chanollama.NewClient("/models/llama2-7b.gguf",
    chanollama.WithMaxTokens(512),                    // Per-request generation limit
    chanollama.WithConversationLimits(2048, 20),     // Total conversation: 2K tokens, 20 messages  
    chanollama.WithTemperature(0.7),
    chanollama.WithGPU(true),
    chanollama.WithTimeouts(5*time.Second, 30*time.Second, 10*time.Second),
)
```

## Related ADRs
- ADR-001: API Entry Point Consolidation (addresses multiple constructor confusion)

---
*This ADR addresses the configuration overlap analysis from Phase 2B of the architectural cleanup initiative.*