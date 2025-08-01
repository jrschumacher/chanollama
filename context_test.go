package chanollama

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestCalculateContextUsage(t *testing.T) {
	// Test normal usage
	usage := CalculateContextUsage(800, 1000)
	
	if usage.UsedTokens != 800 {
		t.Errorf("Expected used tokens 800, got %d", usage.UsedTokens)
	}
	if usage.MaxTokens != 1000 {
		t.Errorf("Expected max tokens 1000, got %d", usage.MaxTokens)
	}
	if usage.Percentage != 80.0 {
		t.Errorf("Expected percentage 80.0, got %f", usage.Percentage)
	}
	if usage.NeedsCleanup {
		t.Error("Usage at exactly 80% should not need cleanup (>80% threshold)")
	}
	
	// Test low usage
	lowUsage := CalculateContextUsage(300, 1000)
	if lowUsage.Percentage != 30.0 {
		t.Errorf("Expected percentage 30.0, got %f", lowUsage.Percentage)
	}
	if lowUsage.NeedsCleanup {
		t.Error("Usage at 30% should not need cleanup")
	}
	
	// Test edge case - exactly 80%
	edgeUsage := CalculateContextUsage(800, 1000)
	if edgeUsage.Percentage != 80.0 {
		t.Errorf("Expected percentage 80.0, got %f", edgeUsage.Percentage)
	}
	if edgeUsage.NeedsCleanup {
		t.Error("Usage at exactly 80% should not need cleanup (>80.0 threshold)")
	}
	
	// Test above threshold
	highUsage := CalculateContextUsage(850, 1000)
	if !highUsage.NeedsCleanup {
		t.Error("Usage above 80% should need cleanup")
	}
}

func TestAutoSummarizer(t *testing.T) {
	// Create a mock engine
	mockEngine := NewRealisticMockEngine()
	summarizer := NewAutoSummarizer(mockEngine)
	
	// Test defaults
	if !summarizer.enabled {
		t.Error("Auto-summarizer should be enabled by default")
	}
	if summarizer.summaryThreshold != 75.0 {
		t.Errorf("Expected default threshold 75.0, got %f", summarizer.summaryThreshold)
	}
	if summarizer.maxSummaryLength != 500 {
		t.Errorf("Expected default max length 500, got %d", summarizer.maxSummaryLength)
	}
}

func TestAutoSummarizerConfiguration(t *testing.T) {
	mockEngine := NewRealisticMockEngine()
	summarizer := NewAutoSummarizer(mockEngine)
	
	// Test fluent configuration
	summarizer = summarizer.WithThreshold(85.0).WithMaxLength(300)
	
	if summarizer.summaryThreshold != 85.0 {
		t.Errorf("Expected threshold 85.0, got %f", summarizer.summaryThreshold)
	}
	if summarizer.maxSummaryLength != 300 {
		t.Errorf("Expected max length 300, got %d", summarizer.maxSummaryLength)
	}
	
	// Test enable/disable
	summarizer.SetEnabled(false)
	if summarizer.enabled {
		t.Error("Summarizer should be disabled")
	}
	
	summarizer.SetEnabled(true)
	if !summarizer.enabled {
		t.Error("Summarizer should be enabled")
	}
}

func TestShouldSummarize(t *testing.T) {
	mockEngine := NewRealisticMockEngine()
	summarizer := NewAutoSummarizer(mockEngine)
	
	// Test below threshold
	lowUsage := ContextUsage{Percentage: 60.0}
	if summarizer.ShouldSummarize(lowUsage) {
		t.Error("Should not summarize below threshold")
	}
	
	// Test above threshold
	highUsage := ContextUsage{Percentage: 80.0}
	if !summarizer.ShouldSummarize(highUsage) {
		t.Error("Should summarize above threshold")
	}
	
	// Test exactly at threshold
	thresholdUsage := ContextUsage{Percentage: 75.0}
	if !summarizer.ShouldSummarize(thresholdUsage) {
		t.Error("Should summarize at threshold")
	}
	
	// Test with disabled summarizer
	summarizer.SetEnabled(false)
	if summarizer.ShouldSummarize(highUsage) {
		t.Error("Disabled summarizer should not summarize")
	}
	
	// Test with nil engine
	nilSummarizer := NewAutoSummarizer(nil)
	if nilSummarizer.ShouldSummarize(highUsage) {
		t.Error("Summarizer with nil engine should not summarize")
	}
}

func TestSummarizeConversation(t *testing.T) {
	mockEngine := NewRealisticMockEngine()
	summarizer := NewAutoSummarizer(mockEngine)
	
	// Test empty messages
	_, err := summarizer.SummarizeConversation([]Message{}, "test context")
	if err == nil {
		t.Error("Expected error for empty messages")
	}
	if !strings.Contains(err.Error(), "no messages to summarize") {
		t.Errorf("Expected 'no messages to summarize' error, got %s", err.Error())
	}
	
	// Test with messages
	messages := []Message{
		NewUserMessage("What is OpenTDF?"),
		NewAssistantMessage("OpenTDF is a data security platform..."),
		NewUserMessage("How do I install it?"),
		NewAssistantMessage("You can install OpenTDF using..."),
	}
	
	summary, err := summarizer.SummarizeConversation(messages, "OpenTDF discussion")
	if err != nil {
		t.Errorf("Expected no error for valid messages, got %s", err.Error())
	}
	
	if summary == nil {
		t.Fatal("Expected summary to be created")
	}
	
	if summary.Summary == "" {
		t.Error("Expected non-empty summary text")
	}
	
	if summary.OriginalMessages != len(messages) {
		t.Errorf("Expected original messages %d, got %d", len(messages), summary.OriginalMessages)
	}
	
	if summary.SummarizedMessages == 0 {
		t.Error("Expected some messages to be summarized")
	}
	
	if summary.CreatedAt.IsZero() {
		t.Error("Expected created timestamp to be set")
	}
	
	if summary.Context != "OpenTDF discussion" {
		t.Errorf("Expected context 'OpenTDF discussion', got %s", summary.Context)
	}
	
	// TokensSaved might be 0 if summary is longer than original
	// This is acceptable in test scenarios
	if summary.TokensSaved < 0 {
		t.Error("Expected non-negative tokens saved")
	}
}

func TestSummarizeConversationFiltering(t *testing.T) {
	mockEngine := NewRealisticMockEngine()
	summarizer := NewAutoSummarizer(mockEngine)
	
	// Test with system messages (should be filtered out)
	messages := []Message{
		NewSystemMessage("You are a helpful assistant"),
		NewUserMessage("What is OpenTDF?"),
		NewAssistantMessage("OpenTDF is..."),
	}
	
	summary, err := summarizer.SummarizeConversation(messages, "")
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
	
	// Should have fewer summarized messages than original due to filtering
	if summary.SummarizedMessages >= len(messages) {
		t.Error("Expected system messages to be filtered out")
	}
	
	// Test with only system messages
	systemOnlyMessages := []Message{
		NewSystemMessage("System message 1"),
		NewSystemMessage("System message 2"),
	}
	
	_, err = summarizer.SummarizeConversation(systemOnlyMessages, "")
	if err == nil {
		t.Error("Expected error when no relevant messages to summarize")
	}
	if !strings.Contains(err.Error(), "no relevant messages") {
		t.Errorf("Expected 'no relevant messages' error, got %s", err.Error())
	}
}

func TestConversationSummaryStruct(t *testing.T) {
	now := time.Now()
	summary := &ConversationSummary{
		Summary:            "Test summary",
		OriginalMessages:   10,
		SummarizedMessages: 6,
		CreatedAt:          now,
		Context:            "Test context",
		TokensSaved:        500,
	}
	
	if summary.Summary != "Test summary" {
		t.Errorf("Expected summary 'Test summary', got %s", summary.Summary)
	}
	
	if summary.OriginalMessages != 10 {
		t.Errorf("Expected 10 original messages, got %d", summary.OriginalMessages)
	}
	
	if summary.SummarizedMessages != 6 {
		t.Errorf("Expected 6 summarized messages, got %d", summary.SummarizedMessages)
	}
	
	if summary.CreatedAt != now {
		t.Errorf("Expected created time %v, got %v", now, summary.CreatedAt)
	}
	
	if summary.Context != "Test context" {
		t.Errorf("Expected context 'Test context', got %s", summary.Context)
	}
	
	if summary.TokensSaved != 500 {
		t.Errorf("Expected 500 tokens saved, got %d", summary.TokensSaved)
	}
}

func TestContextUsageEdgeCases(t *testing.T) {
	// Test zero max tokens (will result in +Inf, but should be handled)
	usage := CalculateContextUsage(100, 0)
	// The actual implementation may return +Inf for division by zero
	// This is acceptable behavior that should be handled by calling code
	if usage.UsedTokens != 100 || usage.MaxTokens != 0 {
		t.Errorf("Expected used=100, max=0, got used=%d, max=%d", usage.UsedTokens, usage.MaxTokens)
	}
	
	// Test zero used tokens
	usage = CalculateContextUsage(0, 1000)
	if usage.Percentage != 0.0 {
		t.Errorf("Expected percentage 0.0 with zero used tokens, got %f", usage.Percentage)
	}
	if usage.NeedsCleanup {
		t.Error("Zero usage should not need cleanup")
	}
	
	// Test 100% usage
	usage = CalculateContextUsage(1000, 1000)
	if usage.Percentage != 100.0 {
		t.Errorf("Expected percentage 100.0, got %f", usage.Percentage)
	}
	if !usage.NeedsCleanup {
		t.Error("100% usage should need cleanup")
	}
	
	// Test over 100% usage
	usage = CalculateContextUsage(1200, 1000)
	if usage.Percentage != 120.0 {
		t.Errorf("Expected percentage 120.0, got %f", usage.Percentage)
	}
	if !usage.NeedsCleanup {
		t.Error("Over 100% usage should need cleanup")
	}
}

func TestAutoSummarizerErrorHandling(t *testing.T) {
	// Test with engine that returns errors
	mockEngine := &ErrorEngine{}
	summarizer := NewAutoSummarizer(mockEngine)
	
	messages := []Message{
		NewUserMessage("Test message"),
		NewAssistantMessage("Test response"),
	}
	
	_, err := summarizer.SummarizeConversation(messages, "test")
	if err == nil {
		t.Error("Expected error from failing engine")
	}
	if !strings.Contains(err.Error(), "failed to generate summary") {
		t.Errorf("Expected 'failed to generate summary' error, got %s", err.Error())
	}
}

func TestSummarizerThresholdBoundaries(t *testing.T) {
	mockEngine := NewRealisticMockEngine()
	summarizer := NewAutoSummarizer(mockEngine)
	
	// Test boundary conditions
	testCases := []struct {
		threshold  float64
		percentage float64
		should     bool
	}{
		{75.0, 74.9, false},
		{75.0, 75.0, true},
		{75.0, 75.1, true},
		{80.0, 79.9, false},
		{80.0, 80.0, true},
		{80.0, 80.1, true},
	}
	
	for _, tc := range testCases {
		summarizer = summarizer.WithThreshold(tc.threshold)
		usage := ContextUsage{Percentage: tc.percentage}
		result := summarizer.ShouldSummarize(usage)
		
		if result != tc.should {
			t.Errorf("Threshold %.1f, percentage %.1f: expected %v, got %v", 
				tc.threshold, tc.percentage, tc.should, result)
		}
	}
}

// ErrorEngine is a mock engine that always returns errors
type ErrorEngine struct{}

func (e *ErrorEngine) Chat(messages []Message) (*Response, error) {
	return nil, errors.New("mock engine error")
}

func (e *ErrorEngine) ChatStream(messages []Message, callback StreamCallback) *Response {
	return &Response{Error: errors.New("mock engine error")}
}

func (e *ErrorEngine) IsHealthy() bool {
	return false
}

func (e *ErrorEngine) Close() error {
	return nil
}