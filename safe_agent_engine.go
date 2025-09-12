package chanollama

import (
	"fmt"
	"strings"
)

// SafeAgentEngine wraps AgentEngine with automatic error recovery
type SafeAgentEngine struct {
	*AgentEngine
	errorCount int
	maxErrors  int
}

// NewSafeAgentEngine creates a new safe agent engine with automatic fallback
func NewSafeAgentEngine(engine Engine) *SafeAgentEngine {
	agentEngine := NewAgentEngine(engine)
	return &SafeAgentEngine{
		AgentEngine: agentEngine,
		errorCount:  0,
		maxErrors:   3, // Disable LLM evaluation after 3 errors
	}
}

// Chat processes a conversation with error recovery
func (sae *SafeAgentEngine) Chat(messages []Message) (*Response, error) {
	response, err := sae.AgentEngine.Chat(messages)
	if err != nil {
		sae.handleError(err)
	}
	return response, err
}

// ChatStream processes a streaming conversation with error recovery
func (sae *SafeAgentEngine) ChatStream(messages []Message, callback StreamCallback) *Response {
	response := sae.AgentEngine.ChatStream(messages, callback)
	if response != nil && response.Error != nil {
		sae.handleError(response.Error)
	}
	return response
}

// handleError manages error recovery and automatic fallback
func (sae *SafeAgentEngine) handleError(err error) {
	if err == nil {
		return
	}

	// Check if this is a context/batch size error
	errorStr := strings.ToLower(err.Error())
	isContextError := strings.Contains(errorStr, "n_tokens_all <= cparams.n_batch") ||
		strings.Contains(errorStr, "ggml_assert") ||
		strings.Contains(errorStr, "context") ||
		strings.Contains(errorStr, "batch")

	if isContextError {
		sae.errorCount++

		// If we've had too many context errors, log warning
		if sae.errorCount >= sae.maxErrors {
			fmt.Printf("⚠️  Context overflow detected. Using fallback agent selection.\n")
		}
	}
}

// GetErrorCount returns the current error count
func (sae *SafeAgentEngine) GetErrorCount() int {
	return sae.errorCount
}

// ResetErrorCount resets the error counter
func (sae *SafeAgentEngine) ResetErrorCount() {
	sae.errorCount = 0
}

// GetCurrentAgent returns the agent that would be selected for a given query with error handling
func (sae *SafeAgentEngine) GetCurrentAgent(query string) (*Agent, error) {
	// If too many errors occurred, use simple fallback (first available agent)
	if sae.errorCount >= sae.maxErrors {
		if sae.evaluator != nil {
			// Simple fallback: return first available agent
			for token := range sae.evaluator.agents {
				return sae.evaluator.GetAgent(token)
			}
			return nil, fmt.Errorf("no agents available")
		}
		return nil, fmt.Errorf("agent evaluator not available")
	}

	return sae.AgentEngine.GetCurrentAgent(query)
}
