package chanollama

import "fmt"

// NewAgentAwareEngine creates an agent-aware engine from a config
func NewAgentAwareEngine(config Config) (*AgentEngine, error) {
	// Create base engine
	baseEngine, err := NewEngine(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create base engine: %w", err)
	}

	// Wrap with agent capabilities
	agentEngine := NewAgentEngine(baseEngine)

	return agentEngine, nil
}

// NewAgentAwareEngineWithDirectory creates an agent-aware engine and loads agents from directory
func NewAgentAwareEngineWithDirectory(config Config, agentDir string) (*AgentEngine, error) {
	engine, err := NewAgentAwareEngine(config)
	if err != nil {
		return nil, err
	}

	// Enable agents
	if err := engine.EnableAgents(agentDir); err != nil {
		// Log warning but don't fail - engine will work without agents
		fmt.Printf("⚠️  Warning: Failed to load agents from %s: %v\n", agentDir, err)
	}

	return engine, nil
}

// NewSimpleAgentEngine creates an agent engine with default settings and keyword-only evaluation
func NewSimpleAgentEngine(config Config) (*AgentEngine, error) {
	engine, err := NewAgentAwareEngine(config)
	if err != nil {
		return nil, err
	}

	// Load core agents programmatically (no file loading needed)
	if err := engine.EnableAgents(""); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load core agents: %v\n", err)
	}

	return engine, nil
}
