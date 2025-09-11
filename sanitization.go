package chanollama

import (
	"strings"
)

// TokenSanitizer provides centralized token sanitization logic
type TokenSanitizer struct {
	// ChatML tokens that should be removed from responses
	chatMLTokens []string

	// Role indicators that should be filtered out
	roleIndicators []string

	// Stop patterns that indicate the model should stop generating
	stopPatterns []string

	// Conversation artifacts that should be cleaned up
	conversationArtifacts []string
}

// NewTokenSanitizer creates a new token sanitizer with default patterns
func NewTokenSanitizer() *TokenSanitizer {
	return &TokenSanitizer{
		chatMLTokens: []string{
			"<|im_start|>",
			"<|im_end|>",
			"<|eot_id|>",
			"<|begin_of_text|>",
			"<|end_of_text|>",
			"<|eom_id|>",
		},
		roleIndicators: []string{
			"user\n",
			"assistant\n",
			"system\n",
			"User:\n",
			"Assistant:\n",
			"System:\n",
			"User: ",
			"Assistant: ",
			"System: ",
		},
		stopPatterns: []string{
			"<|eot_id|>",
			"<|end_of_text|>",
			"<|im_start|>",
			"<|im_end|>",
		},
		conversationArtifacts: []string{
			"User",
			"Assistant",
			"Human",
			"AI",
			"User:",
			"Assistant:",
			"Human:",
			"AI:",
		},
	}
}

// ShouldStop checks if the response contains patterns that indicate generation should stop
func (s *TokenSanitizer) ShouldStop(response string) bool {
	// Check for explicit stop patterns
	for _, pattern := range s.stopPatterns {
		if strings.Contains(response, pattern) {
			return true
		}
	}

	// Check for conversation patterns that indicate self-conversation
	lines := strings.Split(response, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Only stop if we see a clear conversation pattern after some content
		if i > 0 && (strings.HasPrefix(trimmed, "User:") ||
			strings.HasPrefix(trimmed, "Assistant:") ||
			strings.HasPrefix(trimmed, "Human:") ||
			strings.HasPrefix(trimmed, "AI:")) {
			// Check if there's actual content before this line
			hasContent := false
			for j := 0; j < i; j++ {
				if strings.TrimSpace(lines[j]) != "" {
					hasContent = true
					break
				}
			}
			if hasContent {
				return true
			}
		}
	}

	// Stop if we see standalone conversation indicators at the end
	if strings.Contains(response, "\n\nUser\n") ||
		strings.Contains(response, "\n\nAssistant\n") ||
		strings.HasSuffix(response, "\n\nUser") ||
		strings.HasSuffix(response, "\n\nAssistant") ||
		strings.HasSuffix(response, "\nUser") ||
		strings.HasSuffix(response, "\nAssistant") {
		return true
	}

	return false
}

// CleanStreamingToken cleans individual tokens during streaming
func (s *TokenSanitizer) CleanStreamingToken(token string) string {
	// Filter out ChatML tokens as they're generated
	for _, chatToken := range s.chatMLTokens {
		if strings.Contains(token, chatToken) {
			return ""
		}
	}

	// Filter out role indicators that might appear at the start of tokens
	for _, indicator := range s.roleIndicators {
		if strings.HasPrefix(token, indicator) {
			return ""
		}
	}

	// Filter out partial ChatML tokens that might appear across multiple pieces
	if strings.Contains(token, "<|") || strings.Contains(token, "|>") {
		return ""
	}

	return token
}

// CleanResponse removes ChatML tokens and cleans up the full response
func (s *TokenSanitizer) CleanResponse(response string) string {
	cleaned := response

	// Remove ChatML tokens
	for _, token := range s.chatMLTokens {
		cleaned = strings.ReplaceAll(cleaned, token, "")
	}

	// Remove role indicators
	for _, indicator := range s.roleIndicators {
		// Only remove if it appears at the beginning or after a newline
		cleaned = strings.TrimPrefix(cleaned, indicator)
		cleaned = strings.ReplaceAll(cleaned, "\n"+indicator, "\n")
	}

	// Clean up extra whitespace and newlines
	cleaned = strings.TrimSpace(cleaned)

	// Remove trailing conversation artifacts
	cleanedLines := []string{}
	lines := strings.Split(cleaned, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and lines that look like conversation artifacts
		if trimmed == "" || strings.HasPrefix(trimmed, "<|") || strings.HasSuffix(trimmed, "|>") {
			continue
		}

		// Skip standalone role indicators that appear at the end
		if i == len(lines)-1 {
			isArtifact := false
			for _, artifact := range s.conversationArtifacts {
				if trimmed == artifact {
					isArtifact = true
					break
				}
			}
			if isArtifact {
				continue
			}
		}

		cleanedLines = append(cleanedLines, line)
	}

	result := strings.Join(cleanedLines, "\n")

	// Remove trailing conversation artifacts
	for _, artifact := range s.conversationArtifacts {
		result = strings.TrimSuffix(result, "\n"+artifact)
		result = strings.TrimSuffix(result, "\n\n"+artifact)
	}

	return strings.TrimSpace(result)
}

// AddChatMLToken adds a new ChatML token to the sanitizer
func (s *TokenSanitizer) AddChatMLToken(token string) {
	s.chatMLTokens = append(s.chatMLTokens, token)
}

// AddRoleIndicator adds a new role indicator to the sanitizer
func (s *TokenSanitizer) AddRoleIndicator(indicator string) {
	s.roleIndicators = append(s.roleIndicators, indicator)
}

// AddStopPattern adds a new stop pattern to the sanitizer
func (s *TokenSanitizer) AddStopPattern(pattern string) {
	s.stopPatterns = append(s.stopPatterns, pattern)
}

// AddConversationArtifact adds a new conversation artifact to the sanitizer
func (s *TokenSanitizer) AddConversationArtifact(artifact string) {
	s.conversationArtifacts = append(s.conversationArtifacts, artifact)
}
