package chanollama

import (
	"errors"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	modelPath := "/path/to/model.gguf"
	client := New(modelPath)
	
	if client == nil {
		t.Fatal("New returned nil")
	}
	
	if client.GetModelPath() != modelPath {
		t.Errorf("Expected model path %s, got %s", modelPath, client.GetModelPath())
	}
}

func TestNewClientWithConfig(t *testing.T) {
	config := Config{
		ModelPath:    "/path/to/model.gguf",
		Temperature:  0.8,
		MaxTokens:    500,
	}
	
	client := NewWithConfig("/path/to/model.gguf", config)
	if client == nil {
		t.Fatal("NewWithConfig returned nil")
	}
	
	clientConfig := client.GetConfig()
	
	if clientConfig.Temperature != 0.8 {
		t.Errorf("Expected temperature 0.8, got %f", clientConfig.Temperature)
	}
	
	if clientConfig.MaxTokens != 500 {
		t.Errorf("Expected max tokens 500, got %d", clientConfig.MaxTokens)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig("/path/to/model.gguf")
	
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", config.Timeout)
	}
	
	if config.Temperature != 0.7 {
		t.Errorf("Expected default temperature 0.7, got %f", config.Temperature)
	}
}

func TestDefaultEngineConfig(t *testing.T) {
	config := DefaultEngineConfig()
	
	if config.RequestChannelSize <= 0 {
		t.Errorf("Expected positive request channel size, got %d", config.RequestChannelSize)
	}
	
	if config.StreamingChannelSize <= 0 {
		t.Errorf("Expected positive streaming channel size, got %d", config.StreamingChannelSize)
	}
}

// Tests for conversation management

func TestDefaultConversationConfig(t *testing.T) {
	config := DefaultConversationConfig()
	
	if config.MaxTokens != 0 {
		t.Errorf("Expected default max tokens 0 (auto-detect), got %d", config.MaxTokens)
	}
	
	if config.MaxMessages != 50 {
		t.Errorf("Expected default max messages 50, got %d", config.MaxMessages)
	}
	
	if !config.PreserveSystem {
		t.Errorf("Expected PreserveSystem to be true by default")
	}
	
	if config.TruncationStrategy != TruncationOldest {
		t.Errorf("Expected TruncationStrategy to be TruncationOldest by default")
	}
	
	if !config.AutoDetectLimits {
		t.Errorf("Expected AutoDetectLimits to be true by default")
	}
}

func TestConversationManager(t *testing.T) {
	config := ConversationConfig{
		MaxTokens:          100,
		MaxMessages:        3,
		TruncationStrategy: TruncationOldest,
		PreserveSystem:     true,
		AutoDetectLimits:   false,
		BatchSize:          4096,
		EnableLogging:      false,
	}
	
	manager := NewConversationManagerWithConfig(config)
	
	// Test with messages under limit
	messages := []Message{
		NewSystemMessage("System prompt"),
		NewUserMessage("User message"),
		NewAssistantMessage("Assistant response"),
	}
	
	result := manager.TruncateMessages(messages)
	if len(result) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(result))
	}
	
	// Test with messages over count limit
	messages = []Message{
		NewSystemMessage("System prompt"),
		NewUserMessage("User message 1"),
		NewAssistantMessage("Assistant response 1"),
		NewUserMessage("User message 2"),
		NewAssistantMessage("Assistant response 2"),
		NewUserMessage("User message 3"),
	}
	
	result = manager.TruncateMessages(messages)
	
	// Should keep system message + 2 newest messages (within limit of 3)
	if len(result) != 3 {
		t.Errorf("Expected 3 messages after truncation, got %d", len(result))
	}
	
	// Check that system message is preserved
	hasSystem := false
	for _, msg := range result {
		if msg.IsSystem() {
			hasSystem = true
			break
		}
	}
	if !hasSystem {
		t.Error("System message should be preserved during truncation")
	}
}

func TestConversationManagerTokenTruncation(t *testing.T) {
	config := ConversationConfig{
		MaxTokens:          30,  // Very low limit to force truncation
		MaxMessages:        100, // High message limit
		TruncationStrategy: TruncationOldest,
		PreserveSystem:     true,
		AutoDetectLimits:   false,
		BatchSize:          4096,
		EnableLogging:      false,
	}
	
	manager := NewConversationManagerWithConfig(config)
	
	// Create messages that will exceed token limit
	messages := []Message{
		NewSystemMessage("System"),
		NewUserMessage("Long user message that exceeds token limits"),
		NewAssistantMessage("Long assistant response that also exceeds limits"),
		NewUserMessage("Another message"),
		NewAssistantMessage("Final response"),
	}
	
	result := manager.TruncateMessages(messages)
	
	// Should have fewer messages due to token limits
	if len(result) >= len(messages) {
		t.Errorf("Expected fewer messages after token truncation, got %d (original: %d)", len(result), len(messages))
	}
	
	// Check that system message is preserved
	hasSystem := false
	for _, msg := range result {
		if msg.IsSystem() {
			hasSystem = true
			break
		}
	}
	if !hasSystem {
		t.Error("System message should be preserved during token truncation")
	}
	
	// Check that total tokens are reasonable (may not be exactly under limit due to estimation)
	totalTokens := manager.EstimateTokens(result)
	if totalTokens > config.MaxTokens*2 { // Allow some tolerance for estimation inaccuracy
		t.Errorf("Total tokens %d significantly exceeds limit %d", totalTokens, config.MaxTokens)
	}
}

func TestTruncationInfo(t *testing.T) {
	config := ConversationConfig{
		MaxTokens:          50,
		MaxMessages:        2,
		TruncationStrategy: TruncationOldest,
		PreserveSystem:     true,
		AutoDetectLimits:   false,
		BatchSize:          4096,
		EnableLogging:      false,
	}
	
	manager := NewConversationManagerWithConfig(config)
	
	messages := []Message{
		NewSystemMessage("System prompt"),
		NewUserMessage("User message 1"),
		NewAssistantMessage("Assistant response 1"),
		NewUserMessage("User message 2"),
		NewAssistantMessage("Assistant response 2"),
	}
	
	info := manager.GetTruncationInfo(messages)
	
	if info.OriginalCount != 5 {
		t.Errorf("Expected original count 5, got %d", info.OriginalCount)
	}
	
	if info.TruncatedCount >= info.OriginalCount {
		t.Errorf("Expected truncated count to be less than original count")
	}
	
	if !info.WillTruncate {
		t.Error("Expected WillTruncate to be true")
	}
	
	if info.RemovedMessages != (info.OriginalCount - info.TruncatedCount) {
		t.Errorf("RemovedMessages calculation incorrect: %d vs %d", info.RemovedMessages, info.OriginalCount-info.TruncatedCount)
	}
}

func TestClientConversationIntegration(t *testing.T) {
	config := DefaultConfig("/path/to/model.gguf")
	config.ContextSize = 1000 // Set a context size for testing
	
	client := NewWithConfig("/path/to/model.gguf", config)
	
	// Test that conversation manager is properly initialized
	manager := client.GetConversationManager()
	if manager == nil {
		t.Fatal("Conversation manager should not be nil")
	}
	
	// Test that conversation config can be updated
	newConfig := ConversationConfig{
		MaxTokens:          200,
		MaxMessages:        10,
		TruncationStrategy: TruncationNewest,
		PreserveSystem:     false,
		AutoDetectLimits:   false,
		BatchSize:          1000,
		EnableLogging:      true,
	}
	
	client.SetConversationConfig(newConfig)
	
	// Verify config was updated (we can't directly test this without accessing private fields,
	// but we can test that the manager behaves differently)
	messages := make([]Message, 15) // More than the new limit
	for i := range messages {
		messages[i] = NewUserMessage("Test message")
	}
	
	truncated := manager.TruncateMessages(messages)
	if len(truncated) != 10 {
		t.Errorf("Expected 10 messages after applying new config, got %d", len(truncated))
	}
}

// Tests for error recovery

func TestErrorTypes(t *testing.T) {
	// Test that all error types are properly defined
	errors := []error{
		ErrBatchSizeExceeded,
		ErrTokenCountFailed,
		ErrModelNotLoaded,
		ErrInferenceFailed,
		ErrContextOverflow,
	}
	
	for _, err := range errors {
		if err == nil {
			t.Error("Error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Error should have a message")
		}
	}
}

func TestSafeInferenceWithoutModel(t *testing.T) {
	// Test that safeInference handles uninitialized model gracefully
	client := New("/nonexistent/model.gguf")
	
	messages := []Message{
		NewUserMessage("Test message"),
	}
	
	// This should not crash, but return an error
	response, err := client.safeInference(messages, nil)
	
	if err == nil {
		t.Error("Expected error for nonexistent model")
	}
	
	if response == nil {
		t.Error("Response should not be nil even on error")
	}
	
	if response.Error == nil {
		t.Error("Response should contain error")
	}
	
	// Check that it's the right type of error
	if !errors.Is(err, ErrModelNotLoaded) {
		t.Errorf("Expected ErrModelNotLoaded, got %T", err)
	}
}

func TestTruncationMetadata(t *testing.T) {
	client := New("/path/to/model.gguf")
	
	// Create messages that will be truncated
	messages := make([]Message, 100)
	for i := range messages {
		messages[i] = NewUserMessage("Test message that will be truncated")
	}
	
	// Test truncation info
	manager := client.GetConversationManager()
	info := manager.GetTruncationInfo(messages)
	
	if !info.WillTruncate {
		t.Error("Expected truncation for 100 messages")
	}
	
	if info.RemovedMessages <= 0 {
		t.Error("Expected some messages to be removed")
	}
	
	if info.TruncatedCount >= info.OriginalCount {
		t.Error("Truncated count should be less than original count")
	}
}

// Tests for new truncation strategies

func TestTruncationStrategies(t *testing.T) {
	messages := []Message{
		NewUserMessage("First message"),
		NewAssistantMessage("First response"),
		NewUserMessage("Second message"),
		NewAssistantMessage("Second response"),
		NewUserMessage("Third message"),
		NewAssistantMessage("Third response"),
	}
	
	// Test oldest truncation
	config := ConversationConfig{
		MaxTokens:          50,
		MaxMessages:        3,
		TruncationStrategy: TruncationOldest,
		PreserveSystem:     true,
		AutoDetectLimits:   false,
		BatchSize:          4096,
		EnableLogging:      false,
	}
	
	manager := NewConversationManagerWithConfig(config)
	truncated := manager.TruncateMessages(messages)
	
	if len(truncated) != 3 {
		t.Errorf("Expected 3 messages after oldest truncation, got %d", len(truncated))
	}
	
	// Should keep the last 3 messages
	if truncated[0].Content != "Second response" {
		t.Error("Oldest truncation should keep newest messages")
	}
	
	// Test newest truncation
	config.TruncationStrategy = TruncationNewest
	manager.SetConfig(config)
	truncated = manager.TruncateMessages(messages)
	
	if len(truncated) != 3 {
		t.Errorf("Expected 3 messages after newest truncation, got %d", len(truncated))
	}
	
	// Should keep the first 3 messages
	if truncated[0].Content != "First message" {
		t.Error("Newest truncation should keep oldest messages")
	}
	
	// Test middle truncation
	config.TruncationStrategy = TruncationMiddle
	manager.SetConfig(config)
	truncated = manager.TruncateMessages(messages)
	
	if len(truncated) != 3 {
		t.Errorf("Expected 3 messages after middle truncation, got %d", len(truncated))
	}
	
	// Should keep first and last messages
	if truncated[0].Content != "First message" {
		t.Error("Middle truncation should keep first message")
	}
	if truncated[len(truncated)-1].Content != "Third response" {
		t.Error("Middle truncation should keep last message")
	}
}

func TestAutoDetectLimits(t *testing.T) {
	config := DefaultConfig("/path/to/model.gguf")
	config.ContextSize = 2048
	
	client := NewWithConfig("/path/to/model.gguf", config)
	
	// Test that auto-detection is disabled when context size is set
	manager := client.GetConversationManager()
	conversationConfig := manager.config
	
	if conversationConfig.AutoDetectLimits {
		t.Error("Auto-detection should be disabled when context size is explicitly set")
	}
	
	expectedMaxTokens := 1638 // 80% of 2048
	if conversationConfig.MaxTokens != expectedMaxTokens {
		t.Errorf("Expected max tokens to be 80%% of context size (%d), got %d", expectedMaxTokens, conversationConfig.MaxTokens)
	}
}

func TestSetConversationLimits(t *testing.T) {
	client := New("/path/to/model.gguf")
	
	// Test the convenience method
	config := ConversationConfig{
		MaxTokens:          1000,
		MaxMessages:        20,
		TruncationStrategy: TruncationMiddle,
		PreserveSystem:     false,
		AutoDetectLimits:   false,
		BatchSize:          2048,
		EnableLogging:      true,
	}
	
	client.SetConversationLimits(config)
	
	// Verify the configuration was set
	manager := client.GetConversationManager()
	actualConfig := manager.config
	
	if actualConfig.MaxTokens != 1000 {
		t.Errorf("Expected max tokens 1000, got %d", actualConfig.MaxTokens)
	}
	
	if actualConfig.TruncationStrategy != TruncationMiddle {
		t.Errorf("Expected truncation strategy middle, got %s", actualConfig.TruncationStrategy)
	}
}

// Tests for Message and Response fluent API

func TestMessageMetadata(t *testing.T) {
	msg := NewUserMessage("Test message")
	
	// Test metadata operations
	msg = msg.WithMetadata("key1", "value1")
	msg = msg.WithMetadata("key2", "value2")
	
	value1, exists1 := msg.GetMetadata("key1")
	if !exists1 || value1 != "value1" {
		t.Errorf("Expected metadata key1=value1, got %s (exists: %v)", value1, exists1)
	}
	
	value2, exists2 := msg.GetMetadata("key2")
	if !exists2 || value2 != "value2" {
		t.Errorf("Expected metadata key2=value2, got %s (exists: %v)", value2, exists2)
	}
	
	// Test non-existent key
	_, exists3 := msg.GetMetadata("nonexistent")
	if exists3 {
		t.Error("Expected nonexistent key to return false")
	}
}

func TestMessageRoleCheckers(t *testing.T) {
	userMsg := NewUserMessage("User message")
	assistantMsg := NewAssistantMessage("Assistant message")
	systemMsg := NewSystemMessage("System message")
	
	// Test user message
	if !userMsg.IsUser() {
		t.Error("User message should return true for IsUser()")
	}
	if userMsg.IsAssistant() {
		t.Error("User message should return false for IsAssistant()")
	}
	if userMsg.IsSystem() {
		t.Error("User message should return false for IsSystem()")
	}
	if !userMsg.IsRole("user") {
		t.Error("User message should return true for IsRole('user')")
	}
	
	// Test assistant message
	if assistantMsg.IsUser() {
		t.Error("Assistant message should return false for IsUser()")
	}
	if !assistantMsg.IsAssistant() {
		t.Error("Assistant message should return true for IsAssistant()")
	}
	if assistantMsg.IsSystem() {
		t.Error("Assistant message should return false for IsSystem()")
	}
	if !assistantMsg.IsRole("assistant") {
		t.Error("Assistant message should return true for IsRole('assistant')")
	}
	
	// Test system message
	if systemMsg.IsUser() {
		t.Error("System message should return false for IsUser()")
	}
	if systemMsg.IsAssistant() {
		t.Error("System message should return false for IsAssistant()")
	}
	if !systemMsg.IsSystem() {
		t.Error("System message should return true for IsSystem()")
	}
	if !systemMsg.IsRole("system") {
		t.Error("System message should return true for IsRole('system')")
	}
}

func TestMessageAge(t *testing.T) {
	// Create message and immediately check age
	msg := NewUserMessage("Test message")
	age := msg.Age()
	
	if age < 0 {
		t.Error("Message age should not be negative")
	}
	
	if age > time.Second {
		t.Error("Message age should be very small for newly created message")
	}
	
	// Test message with zero timestamp
	msgZero := Message{Role: "user", Content: "test"}
	ageZero := msgZero.Age()
	if ageZero != 0 {
		t.Errorf("Message with zero timestamp should return age 0, got %v", ageZero)
	}
}

func TestMessageString(t *testing.T) {
	msg := NewUserMessage("Hello world")
	expected := "[user] Hello world"
	if msg.String() != expected {
		t.Errorf("Expected %s, got %s", expected, msg.String())
	}
	
	systemMsg := NewSystemMessage("System prompt")
	expectedSystem := "[system] System prompt"
	if systemMsg.String() != expectedSystem {
		t.Errorf("Expected %s, got %s", expectedSystem, systemMsg.String())
	}
}

func TestResponseFluentAPI(t *testing.T) {
	resp := NewResponse("Test response")
	
	// Test fluent API chaining
	resp = resp.WithTokens(100).WithDuration(time.Millisecond * 500).WithModel("test-model").WithFinishReason("stop")
	resp = resp.WithMetadata("key", "value")
	
	if resp.Tokens != 100 {
		t.Errorf("Expected tokens 100, got %d", resp.Tokens)
	}
	
	if resp.Duration != time.Millisecond*500 {
		t.Errorf("Expected duration 500ms, got %v", resp.Duration)
	}
	
	if resp.Model != "test-model" {
		t.Errorf("Expected model 'test-model', got %s", resp.Model)
	}
	
	if resp.FinishReason != "stop" {
		t.Errorf("Expected finish reason 'stop', got %s", resp.FinishReason)
	}
	
	if resp.Metadata["key"] != "value" {
		t.Errorf("Expected metadata key=value, got %s", resp.Metadata["key"])
	}
}

func TestResponseStatusCheckers(t *testing.T) {
	// Test successful response
	resp := NewResponse("Success")
	if !resp.IsSuccess() {
		t.Error("Response without error should be successful")
	}
	if resp.IsError() {
		t.Error("Response without error should not be error")
	}
	
	// Test error response
	errorResp := NewResponse("Error").WithError(errors.New("test error"))
	if errorResp.IsSuccess() {
		t.Error("Response with error should not be successful")
	}
	if !errorResp.IsError() {
		t.Error("Response with error should be error")
	}
	if errorResp.FinishReason != "error" {
		t.Error("Response with error should have finish reason 'error'")
	}
}

func TestResponseToMessage(t *testing.T) {
	resp := NewResponse("Assistant response")
	resp = resp.WithModel("test-model").WithTokens(50).WithDuration(time.Millisecond * 200)
	
	msg := resp.ToMessage()
	
	if !msg.IsAssistant() {
		t.Error("Response should convert to assistant message")
	}
	
	if msg.Content != "Assistant response" {
		t.Errorf("Expected content 'Assistant response', got %s", msg.Content)
	}
	
	// Check metadata transfer
	model, exists := msg.GetMetadata("model")
	if !exists || model != "test-model" {
		t.Errorf("Expected model metadata 'test-model', got %s (exists: %v)", model, exists)
	}
	
	tokens, exists := msg.GetMetadata("tokens")
	if !exists || tokens != "50" {
		t.Errorf("Expected tokens metadata '50', got %s (exists: %v)", tokens, exists)
	}
	
	duration, exists := msg.GetMetadata("duration")
	if !exists || duration != "200ms" {
		t.Errorf("Expected duration metadata '200ms', got %s (exists: %v)", duration, exists)
	}
}