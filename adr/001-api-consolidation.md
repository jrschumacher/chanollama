# ADR-001: API Entry Point Consolidation

## Status
**Proposed** - Awaiting review and implementation

## Context

chanollama currently has multiple confusing entry points that create developer friction and adoption barriers:

### Current Entry Points
```go
// Basic client creation
client1 := chanollama.New(modelPath)
client2 := chanollama.NewWithConfig(modelPath, config)

// Factory pattern
engine1 := chanollama.NewEngine(config)
engine2 := chanollama.NewOfflineEngine(config)

// Advanced features
engine3 := chanollama.NewAdvancedEngine(config)
engine4 := chanollama.NewAdvancedEngineWithConfig(config, engineConfig)

// Feature-specific
ragEngine := chanollama.NewRAGEngine(engine, indexPath, useVector)
agentEngine := chanollama.NewAgentEngine(engine)
```

### Problems Identified
1. **Developer Confusion**: 15+ minutes researching "which constructor should I use?"
2. **Wrong Choices**: Developers pick suboptimal entry points for their needs
3. **Feature Coupling**: RAG, Agents, Streaming require separate constructors
4. **God Objects**: AdvancedEngine violates Single Responsibility Principle
5. **API Inconsistency**: Mix of factory pattern vs direct construction

### Impact Assessment
- **Severity**: Medium - Library is functional but harder to adopt
- **Adoption Friction**: New users bounce off due to API complexity
- **Support Burden**: Repeated questions about which API to use
- **Maintenance Cost**: Multiple code paths for similar functionality

## Decision

Consolidate to **Single Entry Point with Optional Feature Composition** pattern.

### Proposed API Design

```go
// Single, clear entry point
client := chanollama.New(modelPath, opts...)

// Optional feature composition using functional options
client := chanollama.New("/path/to/model.gguf",
    chanollama.WithMaxTokens(1024),
    chanollama.WithTemperature(0.7),
    chanollama.WithGPU(true),
    chanollama.WithRAG(indexPath),
    chanollama.WithAgents(agentDir),
    chanollama.WithStreaming(),
)

// Advanced configurations remain composable
client := chanollama.New("/path/to/model.gguf",
    chanollama.WithEngineConfig(chanollama.EngineConfig{
        RequestChannelSize: 100,
        MaxConcurrentOps:   10,
    }),
    chanollama.WithConversationLimits(4000, 20),
)
```

### Architecture Changes

#### 1. Unified Client Structure
```go
type Client struct {
    // Core components
    config   Config
    engine   Engine
    
    // Optional features (nil if not enabled)
    ragEngine    *RAGEngine
    agentManager *AgentManager
    streamer     *StreamingEngine
    
    // Internal state
    mu           sync.RWMutex
    initialized  bool
}
```

#### 2. Functional Options Pattern
```go
type ClientOption func(*Client) error

func WithMaxTokens(tokens int) ClientOption
func WithTemperature(temp float64) ClientOption
func WithGPU(enable bool) ClientOption
func WithRAG(indexPath string) ClientOption
func WithAgents(agentDir string) ClientOption
func WithStreaming() ClientOption
func WithEngineConfig(config EngineConfig) ClientOption
```

#### 3. Feature Composition
Features are added as optional components rather than requiring separate constructors:
- **RAG**: Wraps base engine with RAG capabilities when enabled
- **Agents**: Adds agent management layer when enabled  
- **Streaming**: Enhances with streaming capabilities when enabled
- **Advanced Engine**: Becomes configuration options rather than separate constructor

## Implementation Plan

### Phase 1: Add New API (Non-Breaking)
1. Implement new `New(modelPath string, opts ...ClientOption)` constructor
2. Add functional options for all current features
3. Create unified Client struct with optional feature composition
4. Add comprehensive tests for new API

### Phase 2: Update Examples & Documentation
1. Update all examples to use new API
2. Create migration guide showing old -> new patterns
3. Update documentation to emphasize new approach
4. Add feature composition examples

### Phase 3: Deprecation (Backward Compatible)
1. Add deprecation warnings to old constructors
2. Maintain full backward compatibility
3. Point to migration guide in deprecation messages
4. Monitor usage patterns in the wild

### Phase 4: Clean Removal (Future Version)
1. Remove deprecated constructors after 2-3 versions
2. Clean up internal code to use only unified approach
3. Remove feature-specific constructors
4. Simplify codebase architecture

## Benefits

### Developer Experience
- **Single Mental Model**: One way to create clients with optional features
- **Discoverability**: IDE autocompletion reveals all available options
- **Flexibility**: Compose exactly the features needed
- **Consistency**: All configuration follows same functional options pattern

### Maintenance Benefits
- **Reduced Complexity**: Single constructor with optional features vs multiple constructors
- **Better Testing**: One main code path with feature toggles vs multiple separate paths
- **Cleaner Architecture**: Features as optional components vs god objects
- **Easier Documentation**: One comprehensive API reference vs multiple constructor docs

### Migration Benefits
- **Non-Breaking**: Existing code continues to work during transition
- **Gradual Adoption**: Developers can migrate at their own pace
- **Clear Path**: Functional options map directly to existing configuration patterns
- **Feature Parity**: No loss of functionality during migration

## Risks & Mitigations

### Risk: Functional Options Complexity
- **Mitigation**: Provide clear examples and sensible defaults
- **Fallback**: Keep simple `New(modelPath)` for basic usage

### Risk: Breaking Changes During Implementation
- **Mitigation**: Implement alongside existing API, deprecate gradually
- **Testing**: Comprehensive test coverage for both old and new APIs

### Risk: Feature Composition Bugs
- **Mitigation**: Careful validation of option combinations
- **Testing**: Test all feature combination permutations

## Success Metrics

### Adoption Metrics
- **New User Onboarding**: Time from "hello world" to first successful usage
- **API Usage Patterns**: Percentage using new vs old constructors
- **Support Questions**: Reduction in "which constructor?" questions

### Code Quality Metrics
- **Lines of Code**: Reduction in constructor-related code
- **Test Coverage**: Maintain >90% coverage during transition
- **Documentation**: Single comprehensive API guide vs multiple constructor docs

## Examples

### Before (Current)
```go
// Confusion: Which one do I need?
client := chanollama.New(modelPath)
// or
engine := chanollama.NewEngine(config)
// or  
advanced := chanollama.NewAdvancedEngine(config)
// plus separate
ragEngine := chanollama.NewRAGEngine(engine, index, true)
```

### After (Proposed)
```go
// Clear, composable, discoverable
client := chanollama.New("/path/to/model.gguf",
    chanollama.WithMaxTokens(1024),
    chanollama.WithRAG("./index"),
    chanollama.WithAgents("./agents"),
    chanollama.WithStreaming(),
)
```

## Related ADRs
- ADR-002: Configuration Structure Unification (addresses overlapping config issues)

---
*This ADR addresses the API consolidation analysis from Phase 2A of the architectural cleanup initiative.*