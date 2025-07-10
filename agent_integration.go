package chanollama

import (
	"context"
	"fmt"
)

// AgentEngine wraps a regular Engine with agent evaluation capabilities
type AgentEngine struct {
	engine    Engine
	evaluator *AgentEvaluator
	enabled   bool
}

// NewAgentEngine creates a new agent-powered engine
func NewAgentEngine(engine Engine) *AgentEngine {
	evaluator := NewAgentEvaluator(engine)
	return &AgentEngine{
		engine:    engine,
		evaluator: evaluator,
		enabled:   false,
	}
}

// EnableAgents loads agents and enables agent-based routing
func (ae *AgentEngine) EnableAgents(agentDirPath string) error {
	if err := ae.evaluator.LoadAgentsFromDirectory(agentDirPath); err != nil {
		return fmt.Errorf("failed to load agents: %w", err)
	}
	ae.enabled = true
	return nil
}

// DisableAgents disables agent-based routing
func (ae *AgentEngine) DisableAgents() {
	ae.enabled = false
}

// IsAgentModeEnabled returns whether agent mode is enabled
func (ae *AgentEngine) IsAgentModeEnabled() bool {
	return ae.enabled
}

// Chat processes a conversation with optional agent routing
func (ae *AgentEngine) Chat(messages []Message) (*Response, error) {
	if !ae.enabled || len(messages) == 0 {
		return ae.engine.Chat(messages)
	}

	// Get the last user message for agent evaluation
	var lastUserMessage string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserMessage = messages[i].Content
			break
		}
	}

	if lastUserMessage == "" {
		return ae.engine.Chat(messages)
	}

	// Evaluate which agent should handle this query
	agentToken, err := ae.evaluator.EvaluateQuery(context.Background(), lastUserMessage)
	if err != nil {
		// Fallback to normal chat if agent evaluation fails
		return ae.engine.Chat(messages)
	}

	// Get the selected agent
	agent, err := ae.evaluator.GetAgent(agentToken)
	if err != nil {
		// Fallback to normal chat if agent not found
		return ae.engine.Chat(messages)
	}

	// Prepare messages with agent's system prompt
	agentMessages := ae.prepareAgentMessages(messages, agent)

	// Process with the underlying engine
	response, err := ae.engine.Chat(agentMessages)
	if err != nil {
		return response, err
	}

	// Add agent metadata to response
	if response != nil {
		if response.Metadata == nil {
			response.Metadata = make(map[string]string)
		}
		response.Metadata["agent_used"] = agentToken
		response.Metadata["agent_name"] = agent.Metadata.Name
		response.Metadata["agent_type"] = string(agent.Type)
	}

	return response, nil
}

// ChatStream processes a streaming conversation with optional agent routing
func (ae *AgentEngine) ChatStream(messages []Message, callback StreamCallback) *Response {
	if !ae.enabled || len(messages) == 0 {
		return ae.engine.ChatStream(messages, callback)
	}

	// Get the last user message for agent evaluation
	var lastUserMessage string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserMessage = messages[i].Content
			break
		}
	}

	if lastUserMessage == "" {
		return ae.engine.ChatStream(messages, callback)
	}

	// Evaluate which agent should handle this query
	agentToken, err := ae.evaluator.EvaluateQuery(context.Background(), lastUserMessage)
	if err != nil {
		// Fallback to normal chat if agent evaluation fails
		return ae.engine.ChatStream(messages, callback)
	}

	// Get the selected agent
	agent, err := ae.evaluator.GetAgent(agentToken)
	if err != nil {
		// Fallback to normal chat if agent not found
		return ae.engine.ChatStream(messages, callback)
	}

	// Prepare messages with agent's system prompt
	agentMessages := ae.prepareAgentMessages(messages, agent)

	// Process with the underlying engine
	response := ae.engine.ChatStream(agentMessages, callback)

	// Add agent metadata to response
	if response != nil {
		if response.Metadata == nil {
			response.Metadata = make(map[string]string)
		}
		response.Metadata["agent_used"] = agentToken
		response.Metadata["agent_name"] = agent.Metadata.Name
		response.Metadata["agent_type"] = string(agent.Type)
	}

	return response
}

// prepareAgentMessages prepares the message array with the agent's system prompt
func (ae *AgentEngine) prepareAgentMessages(messages []Message, agent *Agent) []Message {
	agentMessages := []Message{
		NewSystemMessage(agent.SystemPrompt),
	}

	// Add conversation history, skipping any existing system messages
	for _, msg := range messages {
		if msg.Role != "system" {
			agentMessages = append(agentMessages, msg)
		}
	}

	return agentMessages
}

// GetCurrentAgent returns the agent that would be selected for a given query
func (ae *AgentEngine) GetCurrentAgent(query string) (*Agent, error) {
	if !ae.enabled {
		return nil, fmt.Errorf("agent mode is not enabled")
	}

	agentToken, err := ae.evaluator.EvaluateQuery(context.Background(), query)
	if err != nil {
		return nil, err
	}

	return ae.evaluator.GetAgent(agentToken)
}

// ListAvailableAgents returns all available agents
func (ae *AgentEngine) ListAvailableAgents() []string {
	if !ae.enabled {
		return []string{}
	}
	return ae.evaluator.GetAvailableAgents()
}

// GetAgentInfo returns detailed information about an agent
func (ae *AgentEngine) GetAgentInfo(token string) (*Agent, error) {
	if !ae.enabled {
		return nil, fmt.Errorf("agent mode is not enabled")
	}
	return ae.evaluator.GetAgent(token)
}

// DisableLLMEvaluation disables LLM-based agent evaluation to prevent context overflow
func (ae *AgentEngine) DisableLLMEvaluation() {
	if ae.evaluator != nil {
		ae.evaluator.DisableLLMEvaluation()
	}
}

// EnableLLMEvaluation enables LLM-based agent evaluation
func (ae *AgentEngine) EnableLLMEvaluation() {
	if ae.evaluator != nil {
		ae.evaluator.EnableLLMEvaluation()
	}
}

// IsLLMEvaluationEnabled returns whether LLM evaluation is enabled
func (ae *AgentEngine) IsLLMEvaluationEnabled() bool {
	if ae.evaluator != nil {
		return ae.evaluator.IsLLMEvaluationEnabled()
	}
	return false
}

// IsHealthy checks if the underlying engine is healthy
func (ae *AgentEngine) IsHealthy() bool {
	return ae.engine.IsHealthy()
}

// Close cleans up resources
func (ae *AgentEngine) Close() error {
	return ae.engine.Close()
}