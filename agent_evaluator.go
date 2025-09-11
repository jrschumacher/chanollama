package chanollama

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AgentEvaluator handles intelligent agent selection based on user queries
type AgentEvaluator struct {
	engine            Engine
	agents            map[string]*Agent
	defaultAgent      string
	evaluationTimeout time.Duration
	disableLLMEval    bool // Fallback to keyword-only evaluation
}

// NewAgentEvaluator creates a new agent evaluator
func NewAgentEvaluator(engine Engine) *AgentEvaluator {
	return &AgentEvaluator{
		engine:            engine,
		agents:           make(map[string]*Agent),
		defaultAgent:     "general",
		evaluationTimeout: 10 * time.Second,
		disableLLMEval:    false,
	}
}

// LoadAgent loads an agent with the given token and system prompt
func (ae *AgentEvaluator) LoadAgent(token string, systemPrompt string, metadata AgentMetadata) error {
	agent := NewAgent(metadata.Name, metadata.Description, systemPrompt)
	agent.Metadata = metadata
	ae.agents[token] = agent
	return nil
}

// LoadAgentsFromDirectory loads agents from a directory path
func (ae *AgentEvaluator) LoadAgentsFromDirectory(dirPath string) error {
	// Try to load from file system first, fallback to programmatic loading
	if dirPath != "" {
		return ae.loadAgentsFromFiles(dirPath)
	}
	return ae.loadCoreAgents()
}

// loadAgentsFromFiles loads agents from markdown files in the given directory
func (ae *AgentEvaluator) loadAgentsFromFiles(dirPath string) error {
	// This method would use the pkg/llm agent loader
	// For now, fallback to programmatic loading
	return ae.loadCoreAgents()
}

// loadCoreAgents loads the four core security platform agents
func (ae *AgentEvaluator) loadCoreAgents() error {
	agents := map[string]AgentMetadata{
		"policy": {
			Name:        "Policy Agent",
			Description: "Security policy specialist expert in ABAC design and implementation",
			Capabilities: []string{"policy_design", "attributes", "access_control", "abac"},
			Tags:        []string{"policy", "security", "access_control"},
		},
		"development": {
			Name:        "Development Agent",
			Description: "Multi-language SDK specialist for security platform integration",
			Capabilities: []string{"sdk_integration", "api_usage", "code_examples", "troubleshooting"},
			Tags:        []string{"development", "sdk", "integration", "api"},
		},
		"deployment": {
			Name:        "Deployment Agent",
			Description: "Installation and configuration specialist for security platforms",
			Capabilities: []string{"installation", "configuration", "docker", "kubernetes", "troubleshooting"},
			Tags:        []string{"deployment", "installation", "configuration", "devops"},
		},
		"general": {
			Name:        "General Agent",
			Description: "Security platform concepts and architecture specialist",
			Capabilities: []string{"concepts", "architecture", "explanations", "guidance"},
			Tags:        []string{"general", "concepts", "education", "overview"},
		},
	}

	for token, metadata := range agents {
		// System prompts will be loaded from /pkg/llm
		systemPrompt := ae.getSystemPromptForAgent(token)
		if err := ae.LoadAgent(token, systemPrompt, metadata); err != nil {
			return fmt.Errorf("failed to load agent %s: %w", token, err)
		}
	}

	return nil
}

// getSystemPromptForAgent returns the system prompt for a given agent token
func (ae *AgentEvaluator) getSystemPromptForAgent(token string) string {
	// These will be loaded from files in /pkg/llm, for now return basic prompts
	switch token {
	case "policy":
		return `You are a security policy specialist expert in Attribute-Based Access Control (ABAC) design and implementation.

**Core Knowledge:**
- **Attributes**: Namespace design, hierarchical structures, value management
- **Subject Mappings**: OIDC integration, role mapping, dynamic attributes  
- **Resource Mappings**: Data classification, application integration
- **Registered Resources**: Onboarding, metadata, lifecycle management
- Policy design patterns and best practices

Provide strategic guidance on policy management, ABAC design, and security platform policies.`

	case "development":
		return `You are a multi-language SDK specialist expert in security platform integration and application development.

**Core Knowledge:**
- **SDKs**: Java, Go, JavaScript/TypeScript implementations
- **API Integration**: REST APIs, authentication, error handling
- **Code Examples**: Working implementations, best practices
- **Troubleshooting**: Common issues, debugging, performance optimization

Provide practical development guidance, code examples, and integration solutions.`

	case "deployment":
		return `You are an installation and configuration specialist expert in security platform deployment.

**Core Knowledge:**
- **Installation**: Docker, Kubernetes, bare metal deployment
- **Configuration**: Environment setup, networking, certificates
- **Troubleshooting**: Common deployment issues, debugging
- **DevOps**: CI/CD integration, monitoring, maintenance

Provide practical deployment guidance, configuration help, and troubleshooting solutions.`

	case "general":
		return `You are a security platform concepts and architecture specialist providing educational guidance.

**Core Knowledge:**
- **Security Concepts**: Encryption formats, access control models, data protection mechanisms
- **Architecture**: Platform components, data flows, security model
- **Use Cases**: Industry applications, compliance scenarios
- **Best Practices**: Security patterns, implementation strategies

Provide clear explanations of security platform concepts, architecture, and strategic guidance.`

	default:
		return `You are a helpful security platform assistant. Provide accurate, practical guidance on security platform concepts, implementation, and best practices.`
	}
}

// EvaluateQuery analyzes a user query and returns the best agent token
func (ae *AgentEvaluator) EvaluateQuery(ctx context.Context, query string) (string, error) {
	// Always use keyword-based evaluation to avoid context issues
	// The agent selection doesn't need LLM-level accuracy and keyword matching is reliable
	return ae.keywordBasedEvaluation(query), nil
}

// buildEvaluationPrompt creates the prompt for agent evaluation
func (ae *AgentEvaluator) buildEvaluationPrompt(query string) string {
	return `Route this query to the best agent. Return only one word:

Agents:
- policy: ABAC, attributes, permissions, access control
- development: SDKs, APIs, code, integration
- deployment: install, setup, Docker, configuration  
- general: concepts, explanations, overview

Examples:
"create user policy" -> policy
"Java SDK encrypt" -> development  
"deploy Docker" -> deployment
"what is encryption" -> general

Query: "` + query + `"
Agent:`
}

// parseAgentToken extracts the agent token from the LLM response
func (ae *AgentEvaluator) parseAgentToken(response string) string {
	// Clean up the response
	token := strings.TrimSpace(strings.ToLower(response))
	
	// Remove common prefixes/suffixes
	token = strings.TrimPrefix(token, "response:")
	token = strings.TrimPrefix(token, "agent:")
	token = strings.TrimSpace(token)
	
	// Extract just the token if it's in a sentence
	validTokens := []string{"policy", "development", "deployment", "general"}
	for _, validToken := range validTokens {
		if strings.Contains(token, validToken) {
			return validToken
		}
	}
	
	return ""
}

// keywordBasedEvaluation provides robust evaluation using keyword matching
func (ae *AgentEvaluator) keywordBasedEvaluation(query string) string {
	query = strings.ToLower(query)
	
	// Score-based approach for better accuracy
	scores := map[string]int{
		"policy":      0,
		"development": 0,
		"deployment":  0,
		"general":     0,
	}
	
	// Policy keywords with weights
	policyKeywords := map[string]int{
		"policy": 5, "attribute": 4, "abac": 5, "access": 3, "permission": 4,
		"namespace": 4, "subject": 3, "resource": 3, "mapping": 3, "grant": 3,
		"authorization": 4, "control": 2, "security": 2, "role": 2,
	}
	
	// Development keywords with weights
	devKeywords := map[string]int{
		"sdk": 5, "api": 4, "code": 4, "java": 4, "javascript": 4, "go": 3,
		"integration": 4, "encrypt": 4, "decrypt": 4, "develop": 3, "client": 3,
		"library": 3, "example": 2, "implementation": 3, "typescript": 4,
	}
	
	// Deployment keywords with weights
	deployKeywords := map[string]int{
		"install": 5, "deploy": 5, "docker": 5, "kubernetes": 5, "k8s": 5,
		"config": 4, "setup": 4, "certificate": 4, "network": 3, "server": 3,
		"environment": 3, "production": 3, "staging": 2, "helm": 4,
	}
	
	// General keywords with weights (concepts, architecture)
	generalKeywords := map[string]int{
		"what": 2, "how": 2, "why": 2, "explain": 3, "understand": 3,
		"concept": 4, "architecture": 4, "overview": 3, "introduction": 3,
		"encryption": 3, "security": 3, "platform": 3, "data protection": 2,
	}
	
	// Score each category
	for keyword, weight := range policyKeywords {
		if strings.Contains(query, keyword) {
			scores["policy"] += weight
		}
	}
	
	for keyword, weight := range devKeywords {
		if strings.Contains(query, keyword) {
			scores["development"] += weight
		}
	}
	
	for keyword, weight := range deployKeywords {
		if strings.Contains(query, keyword) {
			scores["deployment"] += weight
		}
	}
	
	for keyword, weight := range generalKeywords {
		if strings.Contains(query, keyword) {
			scores["general"] += weight
		}
	}
	
	// Find the highest scoring agent
	maxScore := 0
	selectedAgent := "general" // default
	
	for agent, score := range scores {
		if score > maxScore {
			maxScore = score
			selectedAgent = agent
		}
	}
	
	return selectedAgent
}

// GetAgent returns the agent for a given token
func (ae *AgentEvaluator) GetAgent(token string) (*Agent, error) {
	agent, exists := ae.agents[token]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", token)
	}
	return agent, nil
}

// GetAvailableAgents returns all available agent tokens
func (ae *AgentEvaluator) GetAvailableAgents() []string {
	tokens := make([]string, 0, len(ae.agents))
	for token := range ae.agents {
		tokens = append(tokens, token)
	}
	return tokens
}

// SetDefaultAgent sets the default agent token
func (ae *AgentEvaluator) SetDefaultAgent(token string) error {
	if _, exists := ae.agents[token]; !exists {
		return fmt.Errorf("agent not found: %s", token)
	}
	ae.defaultAgent = token
	return nil
}

// DisableLLMEvaluation disables LLM-based agent evaluation and uses keyword-only matching
func (ae *AgentEvaluator) DisableLLMEvaluation() {
	ae.disableLLMEval = true
}

// EnableLLMEvaluation enables LLM-based agent evaluation
func (ae *AgentEvaluator) EnableLLMEvaluation() {
	ae.disableLLMEval = false
}

// IsLLMEvaluationEnabled returns whether LLM evaluation is enabled
func (ae *AgentEvaluator) IsLLMEvaluationEnabled() bool {
	return !ae.disableLLMEval
}