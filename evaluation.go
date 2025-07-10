package chanollama

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// EvalResult represents the result of a single evaluation
type EvalResult struct {
	Query           string            `json:"query"`
	Response        string            `json:"response"`
	SelectedAgent   string            `json:"selected_agent,omitempty"`
	ExpectedAgent   string            `json:"expected_agent,omitempty"`
	Metrics         EvalMetrics       `json:"metrics"`
	QualityScore    float64           `json:"quality_score"`
	Timestamp       time.Time         `json:"timestamp"`
	ResponseTime    time.Duration     `json:"response_time"`
	TestID          string            `json:"test_id,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// EvalMetrics contains detailed evaluation metrics
type EvalMetrics struct {
	ResponseLength     int     `json:"response_length"`
	WordCount         int     `json:"word_count"`
	SentenceCount     int     `json:"sentence_count"`
	AgentAccuracy     bool    `json:"agent_accuracy"`
	HelpfulnessScore  float64 `json:"helpfulness_score"`
	CompletenessScore float64 `json:"completeness_score"`
	ActionabilityScore float64 `json:"actionability_score"`
	CodeExampleCount  int     `json:"code_example_count"`
	CommandCount      int     `json:"command_count"`
	LinkCount         int     `json:"link_count"`
}

// TestCase represents a single test case for evaluation
type TestCase struct {
	ID               string            `json:"id"`
	Query            string            `json:"query"`
	ExpectedAgent    string            `json:"expected_agent,omitempty"`
	Category         string            `json:"category"`
	MinLength        int               `json:"min_length,omitempty"`
	MustInclude      []string          `json:"must_include,omitempty"`
	MustNotInclude   []string          `json:"must_not_include,omitempty"`
	ExpectedKeywords []string          `json:"expected_keywords,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// EvalSuite represents a collection of test cases
type EvalSuite struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	TestCases   []TestCase `json:"test_cases"`
}

// EvalSummary provides aggregate results
type EvalSummary struct {
	TotalTests       int             `json:"total_tests"`
	PassedTests      int             `json:"passed_tests"`
	FailedTests      int             `json:"failed_tests"`
	AverageScore     float64         `json:"average_score"`
	AverageLength    float64         `json:"average_length"`
	AverageTime      time.Duration   `json:"average_response_time"`
	AgentAccuracy    float64         `json:"agent_accuracy"`
	Results          []EvalResult    `json:"results"`
	FailedTestIDs    []string        `json:"failed_test_ids,omitempty"`
	Timestamp        time.Time       `json:"timestamp"`
	SystemInfo       SystemInfo      `json:"system_info"`
}

// SystemInfo captures system context for evaluation
type SystemInfo struct {
	Version        string `json:"version"`
	Model          string `json:"model"`
	AgentMode      bool   `json:"agent_mode"`
	EvalMode       string `json:"eval_mode"`
	ConfigHash     string `json:"config_hash,omitempty"`
}

// Evaluator provides comprehensive LLM evaluation capabilities
type Evaluator struct {
	engine       Engine
	scoringFunc  ScoringFunction
	agentMode    bool
	systemInfo   SystemInfo
}

// ScoringFunction allows customizable quality scoring
type ScoringFunction func(response string, metrics EvalMetrics, testCase *TestCase) float64

// NewEvaluator creates a new LLM evaluator
func NewEvaluator(engine Engine, agentMode bool) *Evaluator {
	return &Evaluator{
		engine:      engine,
		scoringFunc: DefaultScoringFunction,
		agentMode:   agentMode,
		systemInfo: SystemInfo{
			AgentMode: agentMode,
			EvalMode:  "comprehensive",
		},
	}
}

// WithScoringFunction sets a custom scoring function
func (e *Evaluator) WithScoringFunction(fn ScoringFunction) *Evaluator {
	e.scoringFunc = fn
	return e
}

// WithSystemInfo sets system information for evaluation context
func (e *Evaluator) WithSystemInfo(info SystemInfo) *Evaluator {
	e.systemInfo = info
	return e
}

// EvaluateQuery evaluates a single query and returns detailed results
func (e *Evaluator) EvaluateQuery(query string) (EvalResult, error) {
	start := time.Now()
	
	// Prepare messages
	messages := []Message{
		NewUserMessage(query),
	}
	
	// Get response
	response, err := e.engine.Chat(messages)
	if err != nil {
		return EvalResult{}, err
	}
	
	responseTime := time.Since(start)
	
	// Extract agent info if available
	selectedAgent := "default"
	if response.Metadata != nil {
		if agent, exists := response.Metadata["agent_used"]; exists {
			selectedAgent = agent
		}
	}
	
	// Calculate metrics
	metrics := e.calculateMetrics(response.Content, selectedAgent, "")
	qualityScore := e.scoringFunc(response.Content, metrics, nil)
	
	return EvalResult{
		Query:        query,
		Response:     response.Content,
		SelectedAgent: selectedAgent,
		Metrics:      metrics,
		QualityScore: qualityScore,
		Timestamp:    time.Now(),
		ResponseTime: responseTime,
	}, nil
}

// EvaluateTestCase evaluates a test case with constraints and expected outcomes
func (e *Evaluator) EvaluateTestCase(testCase TestCase) (EvalResult, error) {
	result, err := e.EvaluateQuery(testCase.Query)
	if err != nil {
		return result, err
	}
	
	// Add test case specific info
	result.TestID = testCase.ID
	result.ExpectedAgent = testCase.ExpectedAgent
	
	// Validate agent accuracy
	if testCase.ExpectedAgent != "" {
		result.Metrics.AgentAccuracy = (result.SelectedAgent == testCase.ExpectedAgent)
	}
	
	// Apply test case constraints and recalculate score
	result.QualityScore = e.scoringFunc(result.Response, result.Metrics, &testCase)
	
	return result, nil
}

// EvaluateSuite runs a complete evaluation suite
func (e *Evaluator) EvaluateSuite(suite EvalSuite) (EvalSummary, error) {
	var results []EvalResult
	
	for _, testCase := range suite.TestCases {
		result, err := e.EvaluateTestCase(testCase)
		if err != nil {
			// Log error but continue with other tests
			fmt.Printf("Warning: Failed to evaluate test %s: %v\n", testCase.ID, err)
			continue
		}
		results = append(results, result)
	}
	
	return e.generateSummary(results), nil
}

// LoadSuite loads an evaluation suite from a JSON file
func LoadSuite(filePath string) (EvalSuite, error) {
	var suite EvalSuite
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return suite, err
	}
	
	err = json.Unmarshal(data, &suite)
	return suite, err
}

// calculateMetrics computes detailed metrics for a response
func (e *Evaluator) calculateMetrics(response, selectedAgent, expectedAgent string) EvalMetrics {
	words := strings.Fields(response)
	sentences := strings.Split(response, ".")
	
	// Count different elements
	codeExamples := strings.Count(response, "```") / 2
	commands := strings.Count(response, "`") / 2
	links := strings.Count(response, "http")
	
	// Calculate subjective scores (0-10)
	helpfulness := e.calculateHelpfulnessScore(response)
	completeness := e.calculateCompletenessScore(response)
	actionability := e.calculateActionabilityScore(response, codeExamples, commands)
	
	return EvalMetrics{
		ResponseLength:     len(response),
		WordCount:         len(words),
		SentenceCount:     len(sentences),
		AgentAccuracy:     selectedAgent == expectedAgent,
		HelpfulnessScore:  helpfulness,
		CompletenessScore: completeness,
		ActionabilityScore: actionability,
		CodeExampleCount:  codeExamples,
		CommandCount:      commands,
		LinkCount:         links,
	}
}

// calculateHelpfulnessScore evaluates how helpful a response is
func (e *Evaluator) calculateHelpfulnessScore(response string) float64 {
	score := 5.0 // baseline
	
	// Positive indicators
	if strings.Contains(response, "example") { score += 1.0 }
	if strings.Contains(response, "step") { score += 0.5 }
	if strings.Contains(response, "```") { score += 1.0 }
	if len(response) > 200 { score += 1.0 }
	if strings.Contains(response, "documentation") { score += 0.5 }
	
	// Negative indicators
	if len(response) < 50 { score -= 3.0 }
	if strings.Contains(response, "I don't know") { score -= 2.0 }
	if strings.Contains(response, "use the") && len(response) < 100 { score -= 1.5 }
	
	return clampScore(score)
}

// calculateCompletenessScore evaluates response thoroughness
func (e *Evaluator) calculateCompletenessScore(response string) float64 {
	score := 5.0 // baseline
	
	// Check for comprehensive indicators
	if strings.Contains(response, "overview") || strings.Contains(response, "introduction") { score += 0.5 }
	if strings.Contains(response, "example") { score += 1.0 }
	if strings.Contains(response, "consider") || strings.Contains(response, "also") { score += 0.5 }
	if strings.Contains(response, "best practice") { score += 1.0 }
	if len(response) > 500 { score += 1.5 }
	
	// Penalize incomplete responses
	if len(response) < 100 { score -= 2.0 }
	if !strings.Contains(response, ".") { score -= 1.0 }
	
	return clampScore(score)
}

// calculateActionabilityScore evaluates how actionable a response is
func (e *Evaluator) calculateActionabilityScore(response string, codeExamples, commands int) float64 {
	score := 3.0 // baseline lower for actionability
	
	// Actionable elements
	score += float64(codeExamples) * 2.0
	score += float64(commands) * 1.0
	
	if strings.Contains(response, "step") { score += 1.0 }
	if strings.Contains(response, "run") || strings.Contains(response, "execute") { score += 0.5 }
	if strings.Contains(response, "example") { score += 0.5 }
	
	return clampScore(score)
}

// DefaultScoringFunction provides the standard quality scoring algorithm
func DefaultScoringFunction(response string, metrics EvalMetrics, testCase *TestCase) float64 {
	// Weighted average of different aspects
	weights := map[string]float64{
		"helpfulness":  0.4,
		"completeness": 0.3,
		"actionability": 0.2,
		"length":      0.1,
	}
	
	lengthScore := 5.0 // baseline
	if metrics.WordCount > 50 { lengthScore += 2.0 }
	if metrics.WordCount > 100 { lengthScore += 2.0 }
	if metrics.WordCount < 20 { lengthScore -= 3.0 }
	lengthScore = clampScore(lengthScore)
	
	score := weights["helpfulness"]*metrics.HelpfulnessScore +
		weights["completeness"]*metrics.CompletenessScore +
		weights["actionability"]*metrics.ActionabilityScore +
		weights["length"]*lengthScore
	
	// Apply test case constraints if provided
	if testCase != nil {
		score = adjustScoreForConstraints(score, response, metrics, *testCase)
	}
	
	return clampScore(score)
}

// adjustScoreForConstraints applies test case specific constraints
func adjustScoreForConstraints(score float64, response string, metrics EvalMetrics, testCase TestCase) float64 {
	// Apply constraints
	if testCase.MinLength > 0 && metrics.ResponseLength < testCase.MinLength {
		score -= 2.0 // Significant penalty for too short responses
	}
	
	responseText := strings.ToLower(response)
	
	// Must include constraints
	for _, mustInclude := range testCase.MustInclude {
		if !strings.Contains(responseText, strings.ToLower(mustInclude)) {
			score -= 1.0
		}
	}
	
	// Must not include constraints
	for _, mustNotInclude := range testCase.MustNotInclude {
		if strings.Contains(responseText, strings.ToLower(mustNotInclude)) {
			score -= 1.5
		}
	}
	
	// Agent accuracy
	if testCase.ExpectedAgent != "" && !metrics.AgentAccuracy {
		score -= 0.5
	}
	
	return score
}

// generateSummary creates evaluation summary from results
func (e *Evaluator) generateSummary(results []EvalResult) EvalSummary {
	totalTests := len(results)
	passedTests := 0
	var totalScore, totalLength float64
	var totalTime time.Duration
	agentCorrect := 0
	agentTotal := 0
	var failedIDs []string
	
	for _, result := range results {
		if result.QualityScore >= 6.0 {
			passedTests++
		} else {
			if result.TestID != "" {
				failedIDs = append(failedIDs, result.TestID)
			}
		}
		
		totalScore += result.QualityScore
		totalLength += float64(result.Metrics.ResponseLength)
		totalTime += result.ResponseTime
		
		if result.ExpectedAgent != "" {
			agentTotal++
			if result.Metrics.AgentAccuracy {
				agentCorrect++
			}
		}
	}
	
	var agentAccuracy float64
	if agentTotal > 0 {
		agentAccuracy = float64(agentCorrect) / float64(agentTotal) * 100
	}
	
	summary := EvalSummary{
		TotalTests:       totalTests,
		PassedTests:      passedTests,
		FailedTests:      totalTests - passedTests,
		AverageScore:     totalScore / float64(totalTests),
		AverageLength:    totalLength / float64(totalTests),
		AverageTime:      totalTime / time.Duration(totalTests),
		AgentAccuracy:    agentAccuracy,
		Results:          results,
		FailedTestIDs:    failedIDs,
		Timestamp:        time.Now(),
		SystemInfo:       e.systemInfo,
	}
	
	return summary
}

// Close cleans up the evaluator resources
func (e *Evaluator) Close() error {
	return e.engine.Close()
}

// clampScore ensures scores stay within 0-10 range
func clampScore(score float64) float64 {
	if score < 0 {
		return 0
	}
	if score > 10 {
		return 10
	}
	return score
}

// CompareResults compares two evaluation summaries for regression detection
func CompareResults(baseline, current EvalSummary) EvaluationComparison {
	return EvaluationComparison{
		ScoreDelta:    current.AverageScore - baseline.AverageScore,
		PassRateDelta: (float64(current.PassedTests)/float64(current.TotalTests) - 
			          float64(baseline.PassedTests)/float64(baseline.TotalTests)) * 100,
		TimeDelta:     current.AverageTime - baseline.AverageTime,
		IsRegression:  current.AverageScore < baseline.AverageScore-0.5, // 0.5 point drop = regression
		Baseline:      baseline,
		Current:       current,
	}
}

// EvaluationComparison represents comparison between two evaluation runs
type EvaluationComparison struct {
	ScoreDelta    float64     `json:"score_delta"`
	PassRateDelta float64     `json:"pass_rate_delta"`
	TimeDelta     time.Duration `json:"time_delta"`
	IsRegression  bool        `json:"is_regression"`
	Baseline      EvalSummary `json:"baseline"`
	Current       EvalSummary `json:"current"`
}