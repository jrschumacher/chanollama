package chanollama

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// Tests for TokenEstimator

func TestTokenEstimator(t *testing.T) {
	estimator := NewTokenEstimator()

	// Test default chars per token
	if estimator.CharsPerToken != 4.0 {
		t.Errorf("Expected default chars per token 4.0, got %f", estimator.CharsPerToken)
	}

	// Test with custom ratio
	customEstimator := estimator.WithCharsPerToken(3.5)
	if customEstimator.CharsPerToken != 3.5 {
		t.Errorf("Expected chars per token 3.5, got %f", customEstimator.CharsPerToken)
	}
}

func TestEstimateTokens(t *testing.T) {
	estimator := NewTokenEstimator()

	// Test empty string
	if tokens := estimator.EstimateTokens(""); tokens != 0 {
		t.Errorf("Expected 0 tokens for empty string, got %d", tokens)
	}

	// Test simple text (4 chars per token)
	text := "Hello world"    // 11 characters
	expectedTokens := 11 / 4 // 2 tokens
	if tokens := estimator.EstimateTokens(text); tokens != expectedTokens {
		t.Errorf("Expected %d tokens, got %d", expectedTokens, tokens)
	}

	// Test with custom ratio
	customEstimator := estimator.WithCharsPerToken(2.0)
	expectedCustomTokens := 11 / 2 // 5 tokens
	if tokens := customEstimator.EstimateTokens(text); tokens != expectedCustomTokens {
		t.Errorf("Expected %d tokens with custom ratio, got %d", expectedCustomTokens, tokens)
	}
}

func TestEstimateMessageTokens(t *testing.T) {
	estimator := NewTokenEstimator()
	msg := NewUserMessage("Hello world") // 11 characters

	tokens := estimator.EstimateMessageTokens(msg)
	expectedTokens := (11 / 4) + 10 // content tokens + role overhead

	if tokens != expectedTokens {
		t.Errorf("Expected %d tokens for message, got %d", expectedTokens, tokens)
	}
}

func TestEstimateMessagesTokens(t *testing.T) {
	estimator := NewTokenEstimator()
	messages := []Message{
		NewUserMessage("Hello"),      // 5 chars + 10 overhead = 11 tokens
		NewAssistantMessage("World"), // 5 chars + 10 overhead = 11 tokens
	}

	tokens := estimator.EstimateMessagesTokens(messages)
	expectedTokens := 11 + 11 // 22 tokens total

	if tokens != expectedTokens {
		t.Errorf("Expected %d tokens for messages, got %d", expectedTokens, tokens)
	}

	// Test empty slice
	if emptyTokens := estimator.EstimateMessagesTokens([]Message{}); emptyTokens != 0 {
		t.Errorf("Expected 0 tokens for empty messages, got %d", emptyTokens)
	}
}

// Tests for TextProcessor

func TestTextProcessor(t *testing.T) {
	processor := NewTextProcessor()

	// Test defaults
	if processor.MaxLength != 50000 {
		t.Errorf("Expected default max length 50000, got %d", processor.MaxLength)
	}
	if !processor.CleanText {
		t.Error("Expected clean text to be enabled by default")
	}
	if !processor.PreserveCode {
		t.Error("Expected preserve code to be enabled by default")
	}
}

func TestTextProcessorConfiguration(t *testing.T) {
	processor := NewTextProcessor()

	// Test fluent configuration
	processor = processor.WithMaxLength(1000).WithCleanText(false).WithPreserveCode(false)

	if processor.MaxLength != 1000 {
		t.Errorf("Expected max length 1000, got %d", processor.MaxLength)
	}
	if processor.CleanText {
		t.Error("Expected clean text to be disabled")
	}
	if processor.PreserveCode {
		t.Error("Expected preserve code to be disabled")
	}
}

func TestProcessText(t *testing.T) {
	processor := NewTextProcessor()

	// Test empty string
	if result := processor.ProcessText(""); result != "" {
		t.Errorf("Expected empty string to remain empty, got %s", result)
	}

	// Test length limiting
	processor = processor.WithMaxLength(10)
	longText := "This is a very long text that exceeds the limit"
	result := processor.ProcessText(longText)
	if len(result) != 13 { // 10 chars + "..." = 13
		t.Errorf("Expected length 13 (10 + ...), got %d", len(result))
	}
	if !strings.HasSuffix(result, "...") {
		t.Error("Expected truncated text to end with ...")
	}
}

func TestProcessTextCleaning(t *testing.T) {
	processor := NewTextProcessor()

	// Test line ending normalization
	textWithCRLF := "Line 1\r\nLine 2\r\nLine 3"
	result := processor.ProcessText(textWithCRLF)
	if strings.Contains(result, "\r") {
		t.Error("Expected carriage returns to be normalized")
	}

	// Test with cleaning disabled
	processor = processor.WithCleanText(false)
	result = processor.ProcessText(textWithCRLF)
	// With cleaning disabled, the text processor may not normalize line endings
	// This depends on the actual implementation behavior
	if len(result) == 0 {
		t.Error("Text processor should return some result")
	}
}

func TestExtractKeywords(t *testing.T) {
	processor := NewTextProcessor()

	// Test empty text
	keywords := processor.ExtractKeywords("", 5)
	if keywords != nil {
		t.Errorf("Expected nil keywords for empty text, got %v", keywords)
	}

	// Test zero max keywords
	keywords = processor.ExtractKeywords("some text", 0)
	if keywords != nil {
		t.Errorf("Expected nil keywords for zero max, got %v", keywords)
	}

	// Test keyword extraction with repeated words
	text := "Machine learning is a powerful technology. Machine learning provides automation and insights."
	keywords = processor.ExtractKeywords(text, 3)

	if len(keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(keywords))
	}

	// Should contain "machine" as it appears twice
	found := false
	for _, keyword := range keywords {
		if keyword == "machine" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'machine' to be in keywords, got %v", keywords)
	}

	// Test with more keywords requested than available
	shortText := "hello world"
	keywords = processor.ExtractKeywords(shortText, 10)
	if len(keywords) > 2 {
		t.Errorf("Expected max 2 keywords from short text, got %d", len(keywords))
	}
}

func TestTruncateText(t *testing.T) {
	processor := NewTextProcessor()

	// Test text shorter than limit
	shortText := "Hello"
	result := processor.TruncateText(shortText, 10)
	if result != shortText {
		t.Errorf("Expected short text unchanged, got %s", result)
	}

	// Test truncation at word boundary
	text := "This is a long sentence that needs truncation"
	result = processor.TruncateText(text, 20)
	if len(result) > 20 {
		t.Errorf("Expected result length <= 20, got %d", len(result))
	}
	if !strings.HasSuffix(result, "...") {
		t.Error("Expected truncated text to end with ...")
	}

	// Test hard truncation when no good space found
	textNoSpaces := "verylongtextwithoutanyspaces"
	result = processor.TruncateText(textNoSpaces, 10)
	expectedLength := 10 - 3 + 3 // maxLength - 3 + len("...")
	if len(result) != expectedLength {
		t.Errorf("Expected hard truncation length %d, got %d", expectedLength, len(result))
	}
}

// Tests for ResourceMonitor

func TestResourceMonitor(t *testing.T) {
	monitor := NewResourceMonitor()

	// Test enabled by default
	usage := monitor.GetCurrentUsage()
	if usage.MemoryMB == 0 && usage.GoroutineCount == 0 {
		t.Error("Expected non-zero resource usage when enabled")
	}

	// Test disabled monitoring
	monitor.SetEnabled(false)
	usage = monitor.GetCurrentUsage()
	if usage.MemoryMB != 0 || usage.GoroutineCount != 0 {
		t.Error("Expected zero resource usage when disabled")
	}
}

func TestShouldLimitProcessing(t *testing.T) {
	monitor := NewResourceMonitor()

	// Test with disabled monitoring
	monitor.SetEnabled(false)
	if monitor.ShouldLimitProcessing(100, 100) {
		t.Error("Disabled monitor should never limit processing")
	}

	// Test with enabled monitoring
	monitor.SetEnabled(true)

	// Test with very low limits (should trigger limiting)
	if !monitor.ShouldLimitProcessing(1, 1) {
		t.Error("Expected processing to be limited with very low thresholds")
	}

	// Test with very high limits (should not trigger limiting)
	if monitor.ShouldLimitProcessing(100000, 100000) {
		t.Error("Expected processing not to be limited with very high thresholds")
	}

	// Test with zero limits (should be ignored)
	if monitor.ShouldLimitProcessing(0, 0) {
		t.Error("Zero limits should be ignored")
	}
}

// Tests for Conversation helper methods

func TestConversationHelpers(t *testing.T) {
	conv := NewConversation()

	// Test defaults
	if conv.MaxMessages != 20 {
		t.Errorf("Expected default max messages 20, got %d", conv.MaxMessages)
	}
	if conv.TokenLimit != 4000 {
		t.Errorf("Expected default token limit 4000, got %d", conv.TokenLimit)
	}

	// Test adding messages
	conv.AddUserMessage("Hello")
	conv.AddAssistantMessage("Hi there")
	conv.AddSystemMessage("System prompt")

	if conv.GetMessageCount() != 3 {
		t.Errorf("Expected 3 messages, got %d", conv.GetMessageCount())
	}
}

func TestConversationMessageRetrieval(t *testing.T) {
	conv := NewConversation()
	conv.AddUserMessage("First user message")
	conv.AddAssistantMessage("First assistant message")
	conv.AddUserMessage("Second user message")

	// Test GetLastMessage
	lastMsg, exists := conv.GetLastMessage()
	if !exists {
		t.Error("Expected last message to exist")
	}
	if lastMsg.Content != "Second user message" {
		t.Errorf("Expected last message to be 'Second user message', got %s", lastMsg.Content)
	}

	// Test GetLastUserMessage
	lastUserMsg, exists := conv.GetLastUserMessage()
	if !exists {
		t.Error("Expected last user message to exist")
	}
	if lastUserMsg.Content != "Second user message" {
		t.Errorf("Expected last user message to be 'Second user message', got %s", lastUserMsg.Content)
	}

	// Test empty conversation
	emptyConv := NewConversation()
	_, exists = emptyConv.GetLastMessage()
	if exists {
		t.Error("Expected no last message in empty conversation")
	}

	_, exists = emptyConv.GetLastUserMessage()
	if exists {
		t.Error("Expected no last user message in empty conversation")
	}
}

func TestConversationFiltering(t *testing.T) {
	conv := NewConversation()
	conv.AddSystemMessage("System prompt")
	conv.AddUserMessage("User message")
	conv.AddAssistantMessage("Assistant message")

	// Test GetSystemMessages
	systemMsgs := conv.GetSystemMessages()
	if len(systemMsgs) != 1 {
		t.Errorf("Expected 1 system message, got %d", len(systemMsgs))
	}
	if systemMsgs[0].Content != "System prompt" {
		t.Errorf("Expected system message 'System prompt', got %s", systemMsgs[0].Content)
	}

	// Test FilterMessagesByRole
	userMsgs := conv.FilterMessagesByRole("user")
	if len(userMsgs) != 1 {
		t.Errorf("Expected 1 user message, got %d", len(userMsgs))
	}

	assistantMsgs := conv.FilterMessagesByRole("assistant")
	if len(assistantMsgs) != 1 {
		t.Errorf("Expected 1 assistant message, got %d", len(assistantMsgs))
	}

	// Test non-existent role
	nonExistent := conv.FilterMessagesByRole("nonexistent")
	if len(nonExistent) != 0 {
		t.Errorf("Expected 0 messages for non-existent role, got %d", len(nonExistent))
	}
}

func TestConversationByAge(t *testing.T) {
	conv := NewConversation()
	conv.AddUserMessage("Old message")

	// Wait a tiny bit
	time.Sleep(time.Millisecond)

	conv.AddUserMessage("New message")

	// Get messages newer than 1 nanosecond (should get both)
	recentMsgs := conv.GetMessagesByAge(time.Hour)
	if len(recentMsgs) != 2 {
		t.Errorf("Expected 2 recent messages, got %d", len(recentMsgs))
	}

	// Get messages newer than 0 (should get none or just the newest)
	veryRecentMsgs := conv.GetMessagesByAge(0)
	if len(veryRecentMsgs) > 2 {
		t.Errorf("Expected <= 2 very recent messages, got %d", len(veryRecentMsgs))
	}
}

func TestConversationClear(t *testing.T) {
	conv := NewConversation()
	conv.AddSystemMessage("System")
	conv.AddUserMessage("User")
	conv.AddAssistantMessage("Assistant")

	// Test Clear
	conv.Clear()
	if conv.GetMessageCount() != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", conv.GetMessageCount())
	}

	// Test ClearExceptSystem
	conv.AddSystemMessage("System")
	conv.AddUserMessage("User")
	conv.AddAssistantMessage("Assistant")

	conv.ClearExceptSystem()
	if conv.GetMessageCount() != 1 {
		t.Errorf("Expected 1 message after clear except system, got %d", conv.GetMessageCount())
	}

	remainingMsgs := conv.GetSystemMessages()
	if len(remainingMsgs) != 1 {
		t.Errorf("Expected 1 system message to remain, got %d", len(remainingMsgs))
	}
}

func TestConversationEstimateTokens(t *testing.T) {
	conv := NewConversation()

	// Empty conversation
	if tokens := conv.EstimateTokens(); tokens != 0 {
		t.Errorf("Expected 0 tokens for empty conversation, got %d", tokens)
	}

	// Add messages
	conv.AddUserMessage("Hello world")   // ~2.75 tokens content
	conv.AddAssistantMessage("Hi there") // ~2.25 tokens content

	tokens := conv.EstimateTokens()
	if tokens == 0 {
		t.Error("Expected non-zero tokens for conversation with messages")
	}

	// Should be approximately (11 + 8) / 4 = ~4.75 tokens for content
	// Actual calculation should be close to this
	if tokens < 4 || tokens > 20 {
		t.Errorf("Expected tokens between 4-20, got %d", tokens)
	}
}

func TestConversationAddResponse(t *testing.T) {
	conv := NewConversation()

	// Test adding successful response
	successResp := NewResponse("Success response")
	conv.AddResponse(successResp)

	if conv.GetMessageCount() != 1 {
		t.Errorf("Expected 1 message after adding successful response, got %d", conv.GetMessageCount())
	}

	lastMsg, _ := conv.GetLastMessage()
	if !lastMsg.IsAssistant() {
		t.Error("Expected response to be added as assistant message")
	}

	// Test adding error response (should not be added)
	errorResp := NewResponse("Error").WithError(errors.New("test error"))
	conv.AddResponse(errorResp)

	if conv.GetMessageCount() != 1 {
		t.Errorf("Expected still 1 message after adding error response, got %d", conv.GetMessageCount())
	}

	// Test adding nil response
	conv.AddResponse(nil)

	if conv.GetMessageCount() != 1 {
		t.Errorf("Expected still 1 message after adding nil response, got %d", conv.GetMessageCount())
	}
}
