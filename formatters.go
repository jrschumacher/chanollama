package chanollama

import (
	"fmt"
	"strings"
	"time"
)

// OutputFormatter interface for different output formats
type OutputFormatter interface {
	FormatSingleResult(result EvalResult) string
	FormatSummary(summary EvalSummary) string
	FormatProgress(results []EvalResult) string
}

// DetailedFormatter provides comprehensive output formatting
type DetailedFormatter struct {
	Title string
}

// NewDetailedFormatter creates a new detailed formatter
func NewDetailedFormatter(title string) *DetailedFormatter {
	return &DetailedFormatter{Title: title}
}

// FormatSingleResult formats a single evaluation result with detailed information
func (f *DetailedFormatter) FormatSingleResult(result EvalResult) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("🔍 Query: %s\n", result.Query))
	if result.SelectedAgent != "" {
		output.WriteString(fmt.Sprintf("🤖 Agent: %s\n", result.SelectedAgent))
	}
	output.WriteString(fmt.Sprintf("📏 Length: %d chars, %d words\n",
		result.Metrics.ResponseLength, result.Metrics.WordCount))
	output.WriteString(fmt.Sprintf("⏱️  Time: %v\n", result.ResponseTime))
	output.WriteString(fmt.Sprintf("📊 Quality Score: %.2f/10\n", result.QualityScore))
	output.WriteString(fmt.Sprintf("   Helpfulness: %.1f | Completeness: %.1f | Actionable: %.1f\n",
		result.Metrics.HelpfulnessScore, result.Metrics.CompletenessScore, result.Metrics.ActionabilityScore))

	return output.String()
}

// FormatSummary formats evaluation summary with comprehensive statistics
func (f *DetailedFormatter) FormatSummary(summary EvalSummary) string {
	var output strings.Builder

	title := f.Title
	if title == "" {
		title = "LLM Evaluation Summary"
	}

	output.WriteString(fmt.Sprintf("\n📊 %s\n", title))
	output.WriteString(strings.Repeat("═", len(title)+4) + "\n")
	output.WriteString(fmt.Sprintf("Tests: %d total, %d passed, %d failed\n",
		summary.TotalTests, summary.PassedTests, summary.FailedTests))
	output.WriteString(fmt.Sprintf("Average Quality Score: %.2f/10\n", summary.AverageScore))
	output.WriteString(fmt.Sprintf("Average Response Length: %.0f characters\n", summary.AverageLength))
	output.WriteString(fmt.Sprintf("Average Response Time: %v\n", summary.AverageTime))

	if summary.AgentAccuracy > 0 {
		output.WriteString(fmt.Sprintf("Agent Routing Accuracy: %.1f%%\n", summary.AgentAccuracy))
	}

	passRate := float64(summary.PassedTests) / float64(summary.TotalTests) * 100
	output.WriteString(fmt.Sprintf("Pass Rate: %.1f%%\n", passRate))

	// Quality assessment
	output.WriteString("\n🎯 Quality Assessment:\n")
	switch {
	case summary.AverageScore >= 8.0:
		output.WriteString("✅ Excellent - Ready for production use\n")
	case summary.AverageScore >= 6.0:
		output.WriteString("👍 Good - Minor improvements recommended\n")
	case summary.AverageScore >= 4.0:
		output.WriteString("⚠️  Fair - Significant improvements needed\n")
	default:
		output.WriteString("❌ Poor - Major issues require attention\n")
	}

	return output.String()
}

// FormatProgress formats progress during batch evaluation
func (f *DetailedFormatter) FormatProgress(results []EvalResult) string {
	var output strings.Builder

	for i, result := range results {
		status := "✅"
		if result.QualityScore < 6.0 {
			status = "⚠️ "
		}
		output.WriteString(fmt.Sprintf("  [%d] %s %.1f/10 - %s\n", i+1, status, result.QualityScore,
			truncateString(result.Query, 60)))
	}

	return output.String()
}

// FormatDetailedResults formats detailed test results for comprehensive output
func (f *DetailedFormatter) FormatDetailedResults(results []EvalResult) string {
	var output strings.Builder

	output.WriteString("\n📋 Detailed Test Results\n")
	output.WriteString(strings.Repeat("═", 25) + "\n\n")

	for i, result := range results {
		output.WriteString(fmt.Sprintf("[%d] \n", i+1))
		output.WriteString(fmt.Sprintf("    Query: %s\n", result.Query))
		if result.SelectedAgent != "" {
			output.WriteString(fmt.Sprintf("    Agent: %s\n", result.SelectedAgent))
		}
		output.WriteString(fmt.Sprintf("    Score: %.1f/10 | Length: %d chars | Time: %v\n",
			result.QualityScore, result.Metrics.ResponseLength, result.ResponseTime))
		output.WriteString(fmt.Sprintf("    Quality: Help %.1f | Complete %.1f | Actionable %.1f\n",
			result.Metrics.HelpfulnessScore, result.Metrics.CompletenessScore, result.Metrics.ActionabilityScore))
		output.WriteString(fmt.Sprintf("    Content: %d code examples, %d commands, %d links\n\n",
			result.Metrics.CodeExampleCount, result.Metrics.CommandCount, result.Metrics.LinkCount))
	}

	return output.String()
}

// SummaryFormatter provides concise output formatting
type SummaryFormatter struct{}

// NewSummaryFormatter creates a new summary formatter
func NewSummaryFormatter() *SummaryFormatter {
	return &SummaryFormatter{}
}

// FormatSingleResult formats a single result concisely
func (f *SummaryFormatter) FormatSingleResult(result EvalResult) string {
	agent := ""
	if result.SelectedAgent != "" {
		agent = fmt.Sprintf(" (Agent: %s)", result.SelectedAgent)
	}
	return fmt.Sprintf("Score: %.1f/10 | Length: %d chars | Time: %v%s\n",
		result.QualityScore, result.Metrics.ResponseLength, result.ResponseTime, agent)
}

// FormatSummary formats summary concisely
func (f *SummaryFormatter) FormatSummary(summary EvalSummary) string {
	passRate := float64(summary.PassedTests) / float64(summary.TotalTests) * 100
	return fmt.Sprintf("Tests: %d | Avg Score: %.1f/10 | Pass Rate: %.1f%% | Avg Time: %v\n",
		summary.TotalTests, summary.AverageScore, passRate, summary.AverageTime)
}

// FormatProgress formats progress concisely
func (f *SummaryFormatter) FormatProgress(results []EvalResult) string {
	if len(results) == 0 {
		return ""
	}

	total := len(results)
	var totalScore float64
	for _, result := range results {
		totalScore += result.QualityScore
	}
	avgScore := totalScore / float64(total)

	return fmt.Sprintf("Progress: %d tests completed | Avg Score: %.1f/10\n", total, avgScore)
}

// JSONFormatter provides JSON output formatting
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// FormatSingleResult formats a single result as JSON
func (f *JSONFormatter) FormatSingleResult(result EvalResult) string {
	// This would return the result as JSON
	// Implementation would use json.Marshal
	return "" // Placeholder - actual JSON marshaling would be implemented
}

// FormatSummary formats summary as JSON
func (f *JSONFormatter) FormatSummary(summary EvalSummary) string {
	// This would return the summary as JSON
	// Implementation would use json.Marshal
	return "" // Placeholder - actual JSON marshaling would be implemented
}

// FormatProgress formats progress as JSON
func (f *JSONFormatter) FormatProgress(results []EvalResult) string {
	// This would return the results as JSON
	// Implementation would use json.Marshal
	return "" // Placeholder - actual JSON marshaling would be implemented
}

// Helper function to truncate strings for display
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatDuration formats duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d.Nanoseconds())/1e6)
	}
	return d.Round(time.Millisecond).String()
}

// GetQualityLevel returns a quality assessment based on score
func GetQualityLevel(score float64) (string, string) {
	switch {
	case score >= 8.0:
		return "excellent", "✅ Excellent"
	case score >= 6.0:
		return "good", "👍 Good"
	case score >= 4.0:
		return "fair", "⚠️  Fair"
	default:
		return "poor", "❌ Poor"
	}
}
