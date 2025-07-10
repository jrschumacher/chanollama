package chanollama

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

// BaselineComparison represents the result of comparing current results with a baseline
type BaselineComparison struct {
	BaselineFile    string               `json:"baseline_file"`
	ComparisonTime  time.Time            `json:"comparison_time"`
	OverallChange   ComparisonChange     `json:"overall_change"`
	ScoreChange     ScoreComparison      `json:"score_change"`
	PerformanceChange PerformanceComparison `json:"performance_change"`
	TestChanges     []TestCaseComparison `json:"test_changes,omitempty"`
	Summary         ComparisonSummary    `json:"summary"`
}

// ComparisonChange represents the direction and magnitude of change
type ComparisonChange struct {
	Direction string  `json:"direction"` // "improved", "degraded", "stable"
	Magnitude float64 `json:"magnitude"` // Absolute change value
	Percentage float64 `json:"percentage"` // Percentage change
}

// ScoreComparison compares quality scores
type ScoreComparison struct {
	BaselineAverage float64          `json:"baseline_average"`
	CurrentAverage  float64          `json:"current_average"`
	Change          ComparisonChange `json:"change"`
	Significance    string           `json:"significance"` // "significant", "minor", "negligible"
}

// PerformanceComparison compares performance metrics
type PerformanceComparison struct {
	ResponseTime TimeComparison `json:"response_time"`
	Throughput   ThroughputComparison `json:"throughput"`
}

// TimeComparison compares response times
type TimeComparison struct {
	BaselineAverage time.Duration    `json:"baseline_average"`
	CurrentAverage  time.Duration    `json:"current_average"`
	Change          ComparisonChange `json:"change"`
}

// ThroughputComparison compares throughput metrics
type ThroughputComparison struct {
	BaselineRate float64          `json:"baseline_rate"` // requests per second
	CurrentRate  float64          `json:"current_rate"`
	Change       ComparisonChange `json:"change"`
}

// TestCaseComparison compares individual test case results
type TestCaseComparison struct {
	TestID         string           `json:"test_id"`
	TestQuery      string           `json:"test_query"`
	BaselineScore  float64          `json:"baseline_score"`
	CurrentScore   float64          `json:"current_score"`
	Change         ComparisonChange `json:"change"`
	Notable        bool             `json:"notable"` // Significant change worth highlighting
}

// ComparisonSummary provides high-level comparison insights
type ComparisonSummary struct {
	TotalTests      int     `json:"total_tests"`
	ImprovedTests   int     `json:"improved_tests"`
	DegradedTests   int     `json:"degraded_tests"`
	StableTests     int     `json:"stable_tests"`
	RegressionRisk  string  `json:"regression_risk"` // "low", "medium", "high"
	Recommendation  string  `json:"recommendation"`
}

// BaselineManager handles baseline operations
type BaselineManager struct {
	toleranceThreshold float64 // Threshold for considering changes significant
	stabilityWindow    float64 // Window for considering scores "stable"
}

// NewBaselineManager creates a new baseline manager with default thresholds
func NewBaselineManager() *BaselineManager {
	return &BaselineManager{
		toleranceThreshold: 0.5, // 0.5 point difference is significant
		stabilityWindow:    0.2, // ±0.2 points is considered stable
	}
}

// NewBaselineManagerWithThresholds creates a baseline manager with custom thresholds
func NewBaselineManagerWithThresholds(tolerance, stability float64) *BaselineManager {
	return &BaselineManager{
		toleranceThreshold: tolerance,
		stabilityWindow:    stability,
	}
}

// SaveBaseline saves evaluation results as a baseline file
func (bm *BaselineManager) SaveBaseline(summary EvalSummary, filePath string) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal baseline data: %w", err)
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write baseline file: %w", err)
	}
	
	return nil
}

// LoadBaseline loads a baseline from file
func (bm *BaselineManager) LoadBaseline(filePath string) (EvalSummary, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return EvalSummary{}, fmt.Errorf("failed to read baseline file: %w", err)
	}
	
	var baseline EvalSummary
	if err := json.Unmarshal(data, &baseline); err != nil {
		return EvalSummary{}, fmt.Errorf("failed to parse baseline JSON: %w", err)
	}
	
	return baseline, nil
}

// CompareWithBaseline compares current results with a baseline
func (bm *BaselineManager) CompareWithBaseline(current EvalSummary, baselineFile string) (BaselineComparison, error) {
	baseline, err := bm.LoadBaseline(baselineFile)
	if err != nil {
		return BaselineComparison{}, err
	}
	
	return bm.Compare(baseline, current, baselineFile), nil
}

// Compare compares two evaluation summaries
func (bm *BaselineManager) Compare(baseline, current EvalSummary, baselineFile string) BaselineComparison {
	scoreComparison := bm.compareScores(baseline, current)
	perfComparison := bm.comparePerformance(baseline, current)
	testComparisons := bm.compareTestCases(baseline.Results, current.Results)
	
	overallChange := bm.determineOverallChange(scoreComparison, perfComparison)
	summary := bm.generateSummary(testComparisons, overallChange)
	
	return BaselineComparison{
		BaselineFile:      baselineFile,
		ComparisonTime:    time.Now(),
		OverallChange:     overallChange,
		ScoreChange:       scoreComparison,
		PerformanceChange: perfComparison,
		TestChanges:       testComparisons,
		Summary:           summary,
	}
}

// compareScores compares quality scores between baseline and current
func (bm *BaselineManager) compareScores(baseline, current EvalSummary) ScoreComparison {
	change := bm.calculateChange(baseline.AverageScore, current.AverageScore)
	significance := bm.determineSignificance(change.Magnitude)
	
	return ScoreComparison{
		BaselineAverage: baseline.AverageScore,
		CurrentAverage:  current.AverageScore,
		Change:          change,
		Significance:    significance,
	}
}

// comparePerformance compares performance metrics
func (bm *BaselineManager) comparePerformance(baseline, current EvalSummary) PerformanceComparison {
	timeChange := bm.calculateChange(
		float64(baseline.AverageTime.Nanoseconds()),
		float64(current.AverageTime.Nanoseconds()),
	)
	
	// Calculate throughput (rough approximation)
	baselineThroughput := 1.0 / baseline.AverageTime.Seconds()
	currentThroughput := 1.0 / current.AverageTime.Seconds()
	throughputChange := bm.calculateChange(baselineThroughput, currentThroughput)
	
	return PerformanceComparison{
		ResponseTime: TimeComparison{
			BaselineAverage: baseline.AverageTime,
			CurrentAverage:  current.AverageTime,
			Change:          timeChange,
		},
		Throughput: ThroughputComparison{
			BaselineRate: baselineThroughput,
			CurrentRate:  currentThroughput,
			Change:       throughputChange,
		},
	}
}

// compareTestCases compares individual test case results
func (bm *BaselineManager) compareTestCases(baseline, current []EvalResult) []TestCaseComparison {
	var comparisons []TestCaseComparison
	
	// Create maps for easier lookup
	baselineMap := make(map[string]EvalResult)
	currentMap := make(map[string]EvalResult)
	
	for _, result := range baseline {
		key := bm.getTestKey(result)
		baselineMap[key] = result
	}
	
	for _, result := range current {
		key := bm.getTestKey(result)
		currentMap[key] = result
	}
	
	// Compare matching test cases
	for key, currentResult := range currentMap {
		if baselineResult, exists := baselineMap[key]; exists {
			change := bm.calculateChange(baselineResult.QualityScore, currentResult.QualityScore)
			notable := math.Abs(change.Magnitude) >= bm.toleranceThreshold
			
			comparisons = append(comparisons, TestCaseComparison{
				TestID:        currentResult.TestID,
				TestQuery:     currentResult.Query,
				BaselineScore: baselineResult.QualityScore,
				CurrentScore:  currentResult.QualityScore,
				Change:        change,
				Notable:       notable,
			})
		}
	}
	
	return comparisons
}

// calculateChange calculates the direction and magnitude of change
func (bm *BaselineManager) calculateChange(baseline, current float64) ComparisonChange {
	if baseline == 0 {
		// Avoid division by zero
		return ComparisonChange{
			Direction:  "stable",
			Magnitude:  0,
			Percentage: 0,
		}
	}
	
	magnitude := current - baseline
	percentage := (magnitude / baseline) * 100
	
	var direction string
	if math.Abs(magnitude) <= bm.stabilityWindow {
		direction = "stable"
	} else if magnitude > 0 {
		direction = "improved"
	} else {
		direction = "degraded"
	}
	
	return ComparisonChange{
		Direction:  direction,
		Magnitude:  magnitude,
		Percentage: percentage,
	}
}

// determineSignificance determines if a change is significant
func (bm *BaselineManager) determineSignificance(magnitude float64) string {
	absMagnitude := math.Abs(magnitude)
	
	switch {
	case absMagnitude >= bm.toleranceThreshold:
		return "significant"
	case absMagnitude >= bm.stabilityWindow:
		return "minor"
	default:
		return "negligible"
	}
}

// determineOverallChange determines the overall direction of change
func (bm *BaselineManager) determineOverallChange(score ScoreComparison, perf PerformanceComparison) ComparisonChange {
	// Primary consideration is score change
	// Secondary consideration is performance (inverse relationship for time)
	
	scoreMagnitude := score.Change.Magnitude
	timeMagnitude := -perf.ResponseTime.Change.Magnitude // Negative because less time is better
	
	// Weighted combination (score is more important)
	overallMagnitude := (scoreMagnitude * 0.8) + (timeMagnitude * 0.2)
	
	var direction string
	if math.Abs(overallMagnitude) <= bm.stabilityWindow {
		direction = "stable"
	} else if overallMagnitude > 0 {
		direction = "improved"
	} else {
		direction = "degraded"
	}
	
	// Calculate percentage based on score change primarily
	percentage := score.Change.Percentage
	
	return ComparisonChange{
		Direction:  direction,
		Magnitude:  overallMagnitude,
		Percentage: percentage,
	}
}

// generateSummary creates a high-level summary of the comparison
func (bm *BaselineManager) generateSummary(testComparisons []TestCaseComparison, overallChange ComparisonChange) ComparisonSummary {
	var improved, degraded, stable int
	
	for _, test := range testComparisons {
		switch test.Change.Direction {
		case "improved":
			improved++
		case "degraded":
			degraded++
		case "stable":
			stable++
		}
	}
	
	// Determine regression risk
	totalTests := len(testComparisons)
	degradationRate := float64(degraded) / float64(totalTests)
	
	var regressionRisk string
	switch {
	case degradationRate >= 0.5:
		regressionRisk = "high"
	case degradationRate >= 0.25:
		regressionRisk = "medium"
	default:
		regressionRisk = "low"
	}
	
	// Generate recommendation
	var recommendation string
	switch overallChange.Direction {
	case "improved":
		recommendation = "Quality improvements detected. Consider updating baseline."
	case "degraded":
		if regressionRisk == "high" {
			recommendation = "Significant regression detected. Investigation recommended."
		} else {
			recommendation = "Minor degradation detected. Monitor trends."
		}
	case "stable":
		recommendation = "Performance stable. No action required."
	}
	
	return ComparisonSummary{
		TotalTests:     totalTests,
		ImprovedTests:  improved,
		DegradedTests:  degraded,
		StableTests:    stable,
		RegressionRisk: regressionRisk,
		Recommendation: recommendation,
	}
}

// getTestKey creates a unique key for test case identification
func (bm *BaselineManager) getTestKey(result EvalResult) string {
	if result.TestID != "" {
		return result.TestID
	}
	// Fallback to query hash or truncated query
	query := result.Query
	if len(query) > 100 {
		query = query[:100]
	}
	return fmt.Sprintf("query_%x", []byte(query))
}

// FormatComparison formats a baseline comparison for display
func FormatComparison(comparison BaselineComparison) string {
	// This would format the comparison results in a human-readable way
	// Implementation would depend on the desired output format
	return fmt.Sprintf("Comparison with baseline: %s (%s change of %.2f%%)",
		comparison.OverallChange.Direction,
		comparison.ScoreChange.Significance,
		comparison.OverallChange.Percentage)
}