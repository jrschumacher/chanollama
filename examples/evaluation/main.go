package main

// TEST_METADATA:
// name: evaluation_framework
// description: Evaluation framework demonstration with mock engine
// goal: Verify evaluation framework functionality without requiring real model
// expected_output: Evaluation examples with scores and metrics
// timeout: 30s

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jrschumacher/chanollama"
)

func main() {
	fmt.Println("🧪 Chanollama Evaluation Framework Example")
	fmt.Println("==========================================")
	
	// Example 1: Single Query Evaluation
	fmt.Println("\n📝 Example 1: Single Query Evaluation")
	fmt.Println("-------------------------------------")
	
	if err := demonstrateSingleQueryEvaluation(); err != nil {
		log.Printf("Single query evaluation failed: %v", err)
	}
	
	// Example 2: Custom Scoring Function
	fmt.Println("\n🎯 Example 2: Custom Scoring Function")
	fmt.Println("------------------------------------")
	
	if err := demonstrateCustomScoring(); err != nil {
		log.Printf("Custom scoring demonstration failed: %v", err)
	}
	
	// Example 3: Test Suite Evaluation
	fmt.Println("\n📊 Example 3: Test Suite Evaluation")
	fmt.Println("-----------------------------------")
	
	if err := demonstrateTestSuiteEvaluation(); err != nil {
		log.Printf("Test suite evaluation failed: %v", err)
	}
	
	// Example 4: Baseline Comparison
	fmt.Println("\n📈 Example 4: Baseline Comparison")
	fmt.Println("---------------------------------")
	
	demonstrateBaselineComparison()
	
	fmt.Println("\n✅ All examples completed!")
	fmt.Println("\n💡 Key Takeaways:")
	fmt.Println("   • Use chanollama.NewEvaluator() for evaluation setup")
	fmt.Println("   • Customize scoring with WithScoringFunction()")
	fmt.Println("   • Create test suites for systematic evaluation")
	fmt.Println("   • Compare results with baselines for regression detection")
}

// demonstrateSingleQueryEvaluation shows basic evaluation usage
func demonstrateSingleQueryEvaluation() error {
	// Create a simple mock engine for demonstration
	engine := &MockEngine{}
	
	// Create evaluator with offline configuration
	evaluator := chanollama.NewEvaluator(engine, false).
		WithSystemInfo(chanollama.SystemInfo{
			Version:   "example-v1.0",
			Model:     "offline-model",
			AgentMode: false,
			EvalMode:  "offline-demonstration",
		})
	defer evaluator.Close()
	
	// Evaluate a query
	query := "How do I encrypt data securely?"
	result, err := evaluator.EvaluateQuery(query)
	if err != nil {
		return fmt.Errorf("evaluation failed: %w", err)
	}
	
	// Display results
	fmt.Printf("Query: %s\n", result.Query)
	fmt.Printf("Response Length: %d characters\n", result.Metrics.ResponseLength)
	fmt.Printf("Word Count: %d\n", result.Metrics.WordCount)
	fmt.Printf("Quality Score: %.2f/10\n", result.QualityScore)
	fmt.Printf("Response Time: %v\n", result.ResponseTime)
	fmt.Printf("Helpfulness: %.1f | Completeness: %.1f | Actionability: %.1f\n",
		result.Metrics.HelpfulnessScore, result.Metrics.CompletenessScore, result.Metrics.ActionabilityScore)
	
	return nil
}

// demonstrateCustomScoring shows how to implement domain-specific scoring
func demonstrateCustomScoring() error {
	engine := &MockEngine{}
	
	// Custom scoring function for security-focused evaluation
	securityScoringFunc := func(response string, metrics chanollama.EvalMetrics, testCase *chanollama.TestCase) float64 {
		// Start with default scoring
		score := chanollama.DefaultScoringFunction(response, metrics, testCase)
		
		// Add security-specific bonuses
		securityTerms := []string{"encryption", "authentication", "authorization", "secure", "certificate"}
		for _, term := range securityTerms {
			if containsIgnoreCase(response, term) {
				score += 0.3
			}
		}
		
		// Bonus for practical security guidance
		if containsIgnoreCase(response, "best practice") { score += 0.5 }
		if containsIgnoreCase(response, "vulnerability") { score += 0.3 }
		if metrics.CodeExampleCount > 0 { score += 0.4 }
		
		// Penalty for missing security context in security queries
		if testCase != nil && testCase.Category == "security" {
			if !containsIgnoreCase(response, "security") { score -= 1.0 }
		}
		
		// Clamp to 0-10 range
		if score < 0 { score = 0 }
		if score > 10 { score = 10 }
		
		return score
	}
	
	// Create evaluator with custom scoring
	evaluator := chanollama.NewEvaluator(engine, false).
		WithScoringFunction(securityScoringFunc)
	defer evaluator.Close()
	
	// Test with security-related query
	query := "What are the best practices for API authentication?"
	result, err := evaluator.EvaluateQuery(query)
	if err != nil {
		return fmt.Errorf("custom scoring evaluation failed: %w", err)
	}
	
	fmt.Printf("Security-Focused Evaluation:\n")
	fmt.Printf("Query: %s\n", result.Query)
	fmt.Printf("Enhanced Quality Score: %.2f/10 (with security bonuses)\n", result.QualityScore)
	
	return nil
}

// demonstrateTestSuiteEvaluation shows systematic evaluation with test suites
func demonstrateTestSuiteEvaluation() error {
	// Create example test suite
	suite := chanollama.EvalSuite{
		Name:        "Example Security Test Suite",
		Version:     "1.0.0",
		Description: "Demonstrates test suite evaluation for security topics",
		TestCases: []chanollama.TestCase{
			{
				ID:               "security_001",
				Query:            "How do I implement secure authentication?",
				Category:         "security",
				MinLength:        200,
				MustInclude:      []string{"authentication", "secure"},
				ExpectedKeywords: []string{"token", "password", "certificate"},
			},
			{
				ID:               "security_002", 
				Query:            "What is data encryption?",
				Category:         "concepts",
				MinLength:        150,
				MustInclude:      []string{"encryption", "data"},
				ExpectedKeywords: []string{"algorithm", "key", "cipher"},
			},
			{
				ID:               "practical_001",
				Query:            "Show me code to hash a password",
				Category:         "code_example",
				MinLength:        100,
				MustInclude:      []string{"password", "hash"},
				ExpectedKeywords: []string{"code", "example", "function"},
			},
		},
	}
	
	// Create evaluator
	engine := &MockEngine{}
	evaluator := chanollama.NewEvaluator(engine, false)
	defer evaluator.Close()
	
	// Run test suite
	summary, err := evaluator.EvaluateSuite(suite)
	if err != nil {
		return fmt.Errorf("test suite evaluation failed: %w", err)
	}
	
	// Display summary
	fmt.Printf("Test Suite: %s\n", suite.Name)
	fmt.Printf("Total Tests: %d\n", summary.TotalTests)
	fmt.Printf("Passed: %d | Failed: %d\n", summary.PassedTests, summary.FailedTests)
	fmt.Printf("Average Score: %.2f/10\n", summary.AverageScore)
	fmt.Printf("Pass Rate: %.1f%%\n", float64(summary.PassedTests)/float64(summary.TotalTests)*100)
	fmt.Printf("Average Response Time: %v\n", summary.AverageTime)
	
	// Show individual results
	fmt.Println("\nIndividual Test Results:")
	for i, result := range summary.Results {
		status := "✅ PASS"
		if result.QualityScore < 6.0 {
			status = "❌ FAIL"
		}
		fmt.Printf("  [%d] %s - Score: %.1f/10 - %s\n", 
			i+1, result.TestID, result.QualityScore, status)
	}
	
	return nil
}

// demonstrateBaselineComparison shows regression detection capabilities
func demonstrateBaselineComparison() {
	// Create mock baseline and current results
	baseline := chanollama.EvalSummary{
		TotalTests:    3,
		PassedTests:   2,
		FailedTests:   1,
		AverageScore:  7.2,
		AverageLength: 350.0,
		AverageTime:   5000000000, // 5 seconds in nanoseconds
	}
	
	current := chanollama.EvalSummary{
		TotalTests:    3,
		PassedTests:   3,
		FailedTests:   0,
		AverageScore:  8.1,
		AverageLength: 420.0,
		AverageTime:   4500000000, // 4.5 seconds in nanoseconds
	}
	
	// Compare results
	comparison := chanollama.CompareResults(baseline, current)
	
	fmt.Printf("Baseline vs Current Comparison:\n")
	fmt.Printf("Score Change: %+.2f points\n", comparison.ScoreDelta)
	fmt.Printf("Pass Rate Change: %+.1f%%\n", comparison.PassRateDelta)
	fmt.Printf("Response Time Change: %v\n", comparison.TimeDelta)
	
	if comparison.IsRegression {
		fmt.Printf("🚨 REGRESSION DETECTED!\n")
	} else if comparison.ScoreDelta > 0.5 {
		fmt.Printf("🎉 SIGNIFICANT IMPROVEMENT!\n")
	} else {
		fmt.Printf("✅ Quality maintained\n")
	}
}

// MockEngine provides a simple mock implementation for demonstration
type MockEngine struct{}

func (m *MockEngine) Chat(messages []chanollama.Message) (*chanollama.Response, error) {
	// Generate different mock responses based on query content
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}
	
	query := messages[len(messages)-1].Content
	response := generateMockResponse(query)
	
	return &chanollama.Response{
		Content:      response,
		Metadata:     make(map[string]string),
		Error:        nil,
		Tokens:       len(response) / 4, // Rough token estimate
		Model:        "offline-mock-model",
		FinishReason: "stop",
	}, nil
}

func (m *MockEngine) ChatStream(messages []chanollama.Message, callback chanollama.StreamCallback) *chanollama.Response {
	response, err := m.Chat(messages)
	if err != nil {
		return &chanollama.Response{Error: err}
	}
	
	// Simulate streaming by calling callback with the full response
	callback(response.Content)
	
	return response
}

func (m *MockEngine) IsHealthy() bool {
	return true
}

func (m *MockEngine) Close() error {
	return nil
}

// generateMockResponse creates realistic mock responses for demonstration
func generateMockResponse(query string) string {
	
	switch {
	case containsIgnoreCase(query, "encrypt") && containsIgnoreCase(query, "data"):
		return `Data encryption is the process of converting readable data into an encoded format that can only be read by someone who has the decryption key. Here are the key concepts:

**Symmetric Encryption**: Uses the same key for encryption and decryption. Examples include AES (Advanced Encryption Standard) and ChaCha20.

**Asymmetric Encryption**: Uses a public-private key pair. The public key encrypts data, and the private key decrypts it. RSA and ECC are common examples.

**Best Practices**:
1. Always use established encryption algorithms (never roll your own)
2. Use appropriate key lengths (AES-256, RSA-2048 minimum)
3. Implement proper key management and rotation
4. Use authenticated encryption to prevent tampering

Example with AES:
` + "```python\nfrom cryptography.fernet import Fernet\nkey = Fernet.generate_key()\ncipher = Fernet(key)\nencrypted = cipher.encrypt(b\"sensitive data\")\n```"

	case containsIgnoreCase(query, "authentication"):
		return `API authentication ensures that only authorized users can access your services. Here are the main approaches:

**Token-Based Authentication**:
- JWT (JSON Web Tokens): Self-contained tokens with claims
- OAuth 2.0: Industry standard for authorization
- API Keys: Simple but less secure for user authentication

**Best Practices**:
1. Use HTTPS always to protect credentials in transit
2. Implement proper token expiration and refresh mechanisms
3. Store secrets securely (environment variables, key vaults)
4. Implement rate limiting to prevent brute force attacks
5. Use strong password policies and multi-factor authentication

**Example JWT Implementation**:
` + "```javascript\nconst jwt = require('jsonwebtoken');\nconst token = jwt.sign({ userId: 123 }, process.env.JWT_SECRET, { expiresIn: '1h' });\n```" + `

**Security Considerations**:
- Never store passwords in plain text
- Use bcrypt or similar for password hashing
- Implement account lockout after failed attempts
- Log authentication events for monitoring`

	case containsIgnoreCase(query, "password") && containsIgnoreCase(query, "hash"):
		return `Password hashing is crucial for secure user authentication. Here's how to implement it properly:

**Use bcrypt** - it's designed specifically for password hashing with built-in salting:

` + "```python\nimport bcrypt\n\n# Hash a password\npassword = \"user_password\".encode('utf-8')\nhashed = bcrypt.hashpw(password, bcrypt.gensalt())\n\n# Verify password\nis_valid = bcrypt.checkpw(password, hashed)\n```" + `

**Key Security Principles**:
1. **Never store plain text passwords** - always hash them
2. **Use a salt** - bcrypt handles this automatically
3. **Use appropriate work factor** - bcrypt default is usually good
4. **Never try to decrypt** - only verify by hashing input and comparing

**Why bcrypt?**
- Adaptive hashing - can increase difficulty over time
- Built-in salt generation
- Resistant to rainbow table attacks
- Time-tested and widely adopted

**Additional Security Measures**:
- Implement account lockout after failed attempts
- Use HTTPS to protect passwords in transit
- Consider multi-factor authentication for sensitive accounts`

	default:
		return `This is a mock response for demonstration purposes. In a real implementation, this would be generated by an actual LLM engine.

The query was: "` + query + `"

This mock response includes:
- Basic explanation of the topic
- Some practical guidance
- Example code or commands where relevant
- Security considerations when applicable

For a complete evaluation framework demonstration, you would:
1. Connect to a real LLM engine
2. Implement domain-specific scoring functions
3. Create comprehensive test suites
4. Set up baseline comparison workflows`
	}
}

// containsIgnoreCase checks if a string contains a substring (case-insensitive)
func containsIgnoreCase(text, substr string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}

// Example of saving results to file for baseline comparison
func saveResultsToFile(summary chanollama.EvalSummary, filename string) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Example of loading baseline from file
func loadBaselineFromFile(filename string) (chanollama.EvalSummary, error) {
	var summary chanollama.EvalSummary
	data, err := os.ReadFile(filename)
	if err != nil {
		return summary, err
	}
	err = json.Unmarshal(data, &summary)
	return summary, err
}