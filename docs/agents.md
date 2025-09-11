# 🤖 Agents Framework

chanollama's agent framework enables building specialized AI assistants with intelligent routing and context-aware responses.

## 📖 Overview

Agents in chanollama are specialized AI assistants that:
- Have specific expertise and capabilities
- Use custom system prompts and configurations
- Can be automatically selected based on user queries
- Support metadata-driven routing and management

## 🏗️ Agent Structure

### Agent Definition

```go
type Agent struct {
    Type         AgentType     `yaml:"type" json:"type"`
    SystemPrompt string        `yaml:"system_prompt" json:"system_prompt"`
    Capabilities []string      `yaml:"capabilities" json:"capabilities"`
    Metadata     AgentMetadata `yaml:"metadata" json:"metadata"`
    // Internal fields
    LastUsed    time.Time `json:"last_used,omitempty"`
    UsageCount  int       `json:"usage_count,omitempty"`
}

type AgentMetadata struct {
    Name         string            `yaml:"name" json:"name"`
    Description  string            `yaml:"description" json:"description"`
    Version      string            `yaml:"version" json:"version"`
    Author       string            `yaml:"author" json:"author"`
    Tags         []string          `yaml:"tags" json:"tags"`
    FilePath     string            `json:"file_path,omitempty"`
    Custom       map[string]string `yaml:"custom" json:"custom,omitempty"`
}
```

### Agent Types

```go
type AgentType string

const (
    AgentTypeGeneral     AgentType = "general"
    AgentTypePolicy      AgentType = "policy"
    AgentTypeDevelopment AgentType = "development"
    AgentTypeDeployment  AgentType = "deployment"
    // Add custom types as needed
)
```

## 📄 Agent File Format

Agents are defined using YAML frontmatter in markdown files:

```yaml
---
name: "Go Development Assistant"
description: "Expert in Go programming, best practices, and ecosystem"
type: "development"
version: "1.0.0"
author: "Your Name"
tags: ["go", "programming", "development", "backend"]
capabilities: ["code-generation", "debugging", "optimization", "testing"]
custom:
  expertise_level: "expert"
  specializations: "web-apis,concurrency,performance"
---

You are an expert Go developer with deep knowledge of:

## Core Expertise
- Go language fundamentals and advanced features
- Concurrency patterns with goroutines and channels  
- HTTP/REST API development with standard library and frameworks
- Database integration and ORM usage
- Testing strategies and best practices
- Performance optimization and profiling

## Response Style
- Provide working, production-ready code examples
- Explain the reasoning behind design decisions
- Highlight potential gotchas and edge cases
- Suggest testing approaches
- Include performance considerations

## Code Quality
- Follow Go best practices and idioms
- Use appropriate error handling patterns
- Include relevant documentation and comments
- Consider security implications
- Optimize for readability and maintainability

Always provide practical, actionable advice with concrete examples.
```

## 🎯 Agent Manager API

### Creating an Agent Manager

```go
// Create agent manager
agentManager := chanollama.NewAgentManager()

// Configure agent directory
agentManager = agentManager.WithAgentDirectory("./agents")

// Load agents from directory
ctx := context.Background()
err := agentManager.LoadAgents(ctx)
if err != nil {
    log.Fatal(err)
}

// Check loaded agents
agents := agentManager.ListAgents()
fmt.Printf("Loaded %d agents\n", len(agents))
```

### Agent Selection

```go
// Create agent context for selection
agentContext := chanollama.NewAgentContext("How do I optimize Go code performance?")

// Select best agent
agent, err := agentManager.SelectAgent(agentContext)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Selected agent: %s (%s)\n", agent.Metadata.Name, agent.Type)
```

### Processing with Agents

```go
// Process query with selected agent
messages := []chanollama.Message{
    chanollama.NewUserMessage("How do I handle errors in Go?"),
}

// Create chat engine
chatEngine := chanollama.NewWithConfig(chanollama.DefaultConfig("llama2:7b"))

// Process with agent
response, err := agentManager.ProcessWithAgent(agent, messages, chatEngine)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Agent response: %s\n", response.Content)
```

## 🔄 Agent Selection Logic

The agent manager uses keyword matching to select the most appropriate agent:

```go
// Selection scoring
type AgentScore struct {
    Agent *Agent
    Score float64
}

// Keywords by agent type
typeKeywords := map[AgentType][]string{
    AgentTypePolicy:     {"policy", "attribute", "abac", "permission", "access", "control"},
    AgentTypeDevelopment: {"code", "sdk", "api", "development", "programming", "bug"},
    AgentTypeDeployment: {"install", "deploy", "setup", "configure", "docker", "kubernetes"},
    AgentTypeGeneral:    {"what", "explain", "help", "overview", "concept"},
}
```

### Custom Selection Strategy

```go
// Implement custom agent selection
func selectAgentByCapability(agents []*chanollama.Agent, requiredCapability string) *chanollama.Agent {
    for _, agent := range agents {
        for _, capability := range agent.Capabilities {
            if strings.Contains(capability, requiredCapability) {
                return agent
            }
        }
    }
    return nil // Fallback to default selection
}

// Use custom selection
agent := selectAgentByCapability(agentManager.ListAgents(), "debugging")
if agent != nil {
    // Process with debugging specialist
}
```

## 📁 Agent Directory Structure

Organize agents in a directory structure:

```
agents/
├── general/
│   └── general_assistant.md
├── development/
│   ├── go_expert.md
│   ├── javascript_expert.md
│   └── python_expert.md
├── deployment/
│   ├── docker_specialist.md
│   └── kubernetes_expert.md
└── policy/
    ├── abac_expert.md
    └── security_specialist.md
```

## 🎛️ Agent Configuration Examples

### General Assistant

```yaml
---
name: "General Assistant"
description: "Helpful general-purpose AI assistant"
type: "general"
tags: ["general", "helper", "assistant"]
capabilities: ["question-answering", "explanation", "general-help"]
---

You are a helpful, knowledgeable assistant. Provide clear, accurate answers to user questions across a wide range of topics. Be concise but thorough, and ask for clarification when needed.
```

### Policy Expert

```yaml
---
name: "ABAC Policy Expert"
description: "Expert in Attribute-Based Access Control and policy management"
type: "policy"
tags: ["abac", "policy", "security", "access-control"]
capabilities: ["policy-design", "attribute-mapping", "access-control", "compliance"]
custom:
  frameworks: "ABAC,RBAC"
  compliance: "GDPR,HIPAA,SOX"
---

You are an expert in Attribute-Based Access Control (ABAC) and policy management systems. You help users:

- Design and implement ABAC policies
- Map attributes and create subject mappings
- Understand access control patterns
- Ensure compliance with regulations
- Optimize policy performance

Provide practical examples using standard ABAC policy syntax when relevant.
```

### Development Specialist

```yaml
---
name: "API Development Expert"
description: "Specialist in REST API design and implementation"
type: "development"
tags: ["api", "rest", "web-services", "backend"]
capabilities: ["api-design", "rest", "http", "web-development"]
custom:
  protocols: "REST,GraphQL,gRPC"
  languages: "Go,JavaScript,Python"
---

You are an expert in API development and web services. Focus on:

- RESTful API design principles
- HTTP best practices and status codes
- Authentication and authorization patterns
- API documentation and testing
- Performance and scalability considerations

Provide working code examples and follow industry standards.
```

## 🔧 Advanced Agent Features

### Agent Events and Lifecycle

```go
// Agent event handler
type AgentEventHandler interface {
    OnAgentLoaded(agent *Agent)
    OnAgentUnloaded(agent *Agent)
    OnAgentSelected(agent *Agent, context *AgentContext)
    OnAgentSwitched(from *Agent, to *Agent, reason string)
    OnAgentError(agent *Agent, err error)
}

// Custom event handler
type MyEventHandler struct{}

func (h *MyEventHandler) OnAgentSelected(agent *Agent, context *AgentContext) {
    log.Printf("Agent selected: %s for query: %s", 
        agent.Metadata.Name, context.UserQuery)
}

// Register handler
agentManager.WithEventHandler(&MyEventHandler{})
```

### Agent Statistics

```go
// Get agent usage statistics
stats := agentManager.GetAgentStats()
fmt.Printf("Agent statistics: %+v\n", stats)

// Example output:
// map[
//   "total_agents": 5,
//   "agents_by_type": map[string]int{
//     "general": 1,
//     "development": 2,
//     "policy": 1,
//     "deployment": 1,
//   },
//   "last_refresh": "2024-01-15T10:30:00Z",
// ]
```

### Dynamic Agent Loading

```go
// Refresh agents from directory
err := agentManager.RefreshAgents(ctx)
if err != nil {
    log.Printf("Failed to refresh agents: %v", err)
}

// Hot-reload agents when files change
func watchAgentDirectory(agentManager *chanollama.AgentManager, dir string) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    log.Printf("Agent file modified: %s", event.Name)
                    agentManager.RefreshAgents(context.Background())
                }
            case err := <-watcher.Errors:
                log.Printf("Watcher error: %v", err)
            }
        }
    }()

    err = watcher.Add(dir)
    if err != nil {
        log.Fatal(err)
    }
}
```

## 💡 Best Practices

### 1. Agent Specialization

```yaml
# Good: Specific, focused agent
---
name: "Go Concurrency Expert"
description: "Specialist in Go goroutines, channels, and concurrent patterns"
type: "development"
capabilities: ["concurrency", "goroutines", "channels", "sync"]
---

# Less optimal: Too broad
---
name: "Programming Expert"
description: "Expert in all programming languages"
type: "development"
capabilities: ["programming"]
---
```

### 2. Clear Agent Boundaries

```go
// Structure agents by domain
agents/
├── security/          # Security-focused agents
│   ├── encryption_expert.md
│   └── auth_specialist.md
├── performance/       # Performance optimization
│   ├── go_performance.md
│   └── database_optimization.md
└── architecture/      # System design
    ├── microservices_expert.md
    └── scalability_specialist.md
```

### 3. Consistent Metadata

```yaml
# Use consistent tagging strategy
tags: ["domain", "technology", "level", "use-case"]
# Examples:
tags: ["security", "go", "intermediate", "authentication"]
tags: ["performance", "database", "advanced", "optimization"]
tags: ["general", "explanation", "beginner", "concepts"]
```

### 4. Agent Composition

```go
// Combine multiple agents for complex queries
func handleComplexQuery(query string, agentManager *chanollama.AgentManager) {
    // Step 1: Get architectural guidance
    archAgent, _ := agentManager.GetAgentByType(chanollama.AgentTypeArchitecture)
    archResponse, _ := processWithAgent(archAgent, query)
    
    // Step 2: Get implementation details
    devAgent, _ := agentManager.GetAgentByType(chanollama.AgentTypeDevelopment)
    implQuery := fmt.Sprintf("Based on this architecture: %s\nHow do I implement: %s", 
        archResponse.Content, query)
    implResponse, _ := processWithAgent(devAgent, implQuery)
    
    // Step 3: Get deployment guidance
    deployAgent, _ := agentManager.GetAgentByType(chanollama.AgentTypeDeployment)
    deployQuery := fmt.Sprintf("How do I deploy this solution: %s", implResponse.Content)
    deployResponse, _ := processWithAgent(deployAgent, deployQuery)
}
```

## 🧪 Testing Agents

### Agent Validation

```go
func validateAgent(agent *chanollama.Agent) error {
    if agent.Metadata.Name == "" {
        return fmt.Errorf("agent must have a name")
    }
    
    if agent.SystemPrompt == "" {
        return fmt.Errorf("agent must have a system prompt")
    }
    
    if len(agent.Capabilities) == 0 {
        return fmt.Errorf("agent should define capabilities")
    }
    
    return nil
}

// Test agent selection
func testAgentSelection(t *testing.T) {
    agentManager := chanollama.NewAgentManager().WithAgentDirectory("./test_agents")
    ctx := context.Background()
    
    err := agentManager.LoadAgents(ctx)
    assert.NoError(t, err)
    
    // Test policy query routes to policy agent
    context := chanollama.NewAgentContext("How do I create ABAC policies?")
    agent, err := agentManager.SelectAgent(context)
    assert.NoError(t, err)
    assert.Equal(t, chanollama.AgentTypePolicy, agent.Type)
}
```

## 🔗 Integration Examples

### Web API with Agent Routing

```go
func handleChatRequest(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Message string `json:"message"`
        AgentType string `json:"agent_type,omitempty"`
    }
    
    json.NewDecoder(r.Body).Decode(&req)
    
    var agent *chanollama.Agent
    var err error
    
    if req.AgentType != "" {
        // Use specific agent type
        agent, err = agentManager.GetAgentByType(chanollama.AgentType(req.AgentType))
    } else {
        // Auto-select agent
        context := chanollama.NewAgentContext(req.Message)
        agent, err = agentManager.SelectAgent(context)
    }
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process with selected agent
    messages := []chanollama.Message{
        chanollama.NewUserMessage(req.Message),
    }
    
    response, err := agentManager.ProcessWithAgent(agent, messages, chatEngine)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": response.Content,
        "agent": map[string]string{
            "name": agent.Metadata.Name,
            "type": string(agent.Type),
        },
    })
}
```

For more examples and detailed API documentation, see the [Agent API Reference](./api/agents.md).