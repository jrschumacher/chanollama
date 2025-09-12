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
}

// NewAgentEvaluator creates a new agent evaluator
func NewAgentEvaluator(engine Engine) *AgentEvaluator {
	return &AgentEvaluator{
		engine:            engine,
		agents:            make(map[string]*Agent),
		defaultAgent:      "", // No default agent - developers must define their own
		evaluationTimeout: 10 * time.Second,
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
	if dirPath == "" {
		return fmt.Errorf("directory path cannot be empty - developers must define their own agents")
	}
	return ae.loadAgentsFromFiles(dirPath)
}

// loadAgentsFromFiles loads agents from markdown files in the given directory
func (ae *AgentEvaluator) loadAgentsFromFiles(dirPath string) error {
	// TODO: Implement agent loading from directory
	// This would scan for agent definition files (e.g., *.md, *.yaml, *.json)
	// and load agent metadata and system prompts
	return fmt.Errorf("loading agents from files not yet implemented - use LoadAgent() to define agents programmatically")
}


// EvaluateQuery analyzes a user query and returns the best agent token
func (ae *AgentEvaluator) EvaluateQuery(ctx context.Context, query string) (string, error) {
	// Check if any agents are loaded
	if len(ae.agents) == 0 {
		return "", fmt.Errorf("no agents loaded - use LoadAgent() to define agents for your application")
	}
	
	// Use LLM-based evaluation if engine is available
	if ae.engine != nil {
		agent, err := ae.llmBasedEvaluation(ctx, query)
		if err == nil && agent != "" {
			return agent, nil
		}
		// If LLM evaluation fails, fall through to simple fallback
	}
	
	// Simple fallback: return first available agent
	// This ensures the system always works, even without LLM
	for token := range ae.agents {
		return token, nil
	}
	
	return "", fmt.Errorf("no agents available")
}

// llmBasedEvaluation uses the LLM engine to intelligently select the best agent
func (ae *AgentEvaluator) llmBasedEvaluation(_ context.Context, query string) (string, error) {
	prompt := ae.buildEvaluationPrompt(query)
	messages := []Message{NewUserMessage(prompt)}
	
	// For now, use the basic Chat method since we have an Engine interface
	// TODO: In a more complete implementation, this would use a context-aware method with timeout
	response, err := ae.engine.Chat(messages)
	if err != nil {
		return "", err // Return error, let caller handle fallback
	}
	
	// Parse the agent token from response
	token := ae.parseAgentToken(response.Content)
	if token == "" {
		return "", fmt.Errorf("could not parse agent from LLM response: %s", response.Content)
	}
	
	return token, nil
}

// buildEvaluationPrompt creates the prompt for agent evaluation using loaded agent metadata
func (ae *AgentEvaluator) buildEvaluationPrompt(query string) string {
	var agentDescriptions []string
	
	// Generate agent descriptions from metadata
	for token, agent := range ae.agents {
		if agent.Metadata.Description != "" {
			desc := fmt.Sprintf("- %s: %s", token, agent.Metadata.Description)
			agentDescriptions = append(agentDescriptions, desc)
		} else {
			// Fallback to tags if no description
			tags := strings.Join(agent.Metadata.Tags, ", ")
			desc := fmt.Sprintf("- %s: %s", token, tags)
			agentDescriptions = append(agentDescriptions, desc)
		}
	}
	
	// If no agents loaded, return simple prompt
	if len(agentDescriptions) == 0 {
		return fmt.Sprintf(`Analyze this query and respond with the most appropriate category:

Query: "%s"
Category:`, query)
	}
	
	return fmt.Sprintf(`Route this query to the best agent. Return only the agent name:

Available Agents:
%s

Query: "%s"
Agent:`, strings.Join(agentDescriptions, "\n"), query)
}

// parseAgentToken extracts the agent token from the LLM response
func (ae *AgentEvaluator) parseAgentToken(response string) string {
	// Clean up the response
	token := strings.TrimSpace(strings.ToLower(response))

	// Remove common prefixes/suffixes
	token = strings.TrimPrefix(token, "response:")
	token = strings.TrimPrefix(token, "agent:")
	token = strings.TrimPrefix(token, "category:")
	token = strings.TrimSpace(token)

	// First, try exact match with loaded agents
	if _, exists := ae.agents[token]; exists {
		return token
	}

	// If no exact match, look for partial matches in the response
	for agentToken := range ae.agents {
		if strings.Contains(token, agentToken) {
			return agentToken
		}
	}

	// If still no match, check if any agent names are mentioned in the response
	for agentToken, agent := range ae.agents {
		agentName := strings.ToLower(agent.Metadata.Name)
		if strings.Contains(token, agentName) || strings.Contains(response, agentName) {
			return agentToken
		}
	}

	return ""
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

