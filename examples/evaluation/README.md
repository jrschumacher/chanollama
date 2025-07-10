# 🧪 Chanollama Evaluation Framework Example

This example demonstrates how to use the chanollama evaluation framework to assess and improve LLM response quality in your own projects.

## 🎯 **What This Example Shows**

### **1. Single Query Evaluation**
- Basic evaluation setup with `chanollama.NewEvaluator()`
- Quality metrics calculation (helpfulness, completeness, actionability)
- Response analysis (length, word count, code examples)

### **2. Custom Scoring Functions**
- Implementing domain-specific quality scoring
- Adding bonuses for relevant terminology and best practices
- Penalties for missing critical content

### **3. Test Suite Evaluation**
- Creating structured test cases with constraints
- Batch evaluation with automatic result aggregation
- Pass/fail determination with quality thresholds

### **4. Baseline Comparison**
- Regression detection by comparing evaluation runs
- Tracking quality improvements or degradations over time
- CI-ready quality gates and automated alerts

## 🚀 **Running the Example**

```bash
# From the chanollama directory
go run examples/evaluation/main.go
```

### **Expected Output**
```
🧪 Chanollama Evaluation Framework Example
==========================================

📝 Example 1: Single Query Evaluation
-------------------------------------
Query: How do I encrypt data securely?
Response Length: 782 characters
Word Count: 125
Quality Score: 8.2/10
Response Time: 1.2ms
Helpfulness: 8.5 | Completeness: 7.8 | Actionability: 8.4

🎯 Example 2: Custom Scoring Function
------------------------------------
Security-Focused Evaluation:
Query: What are the best practices for API authentication?
Enhanced Quality Score: 9.1/10 (with security bonuses)

📊 Example 3: Test Suite Evaluation
-----------------------------------
Test Suite: Example Security Test Suite
Total Tests: 3
Passed: 3 | Failed: 0
Average Score: 8.4/10
Pass Rate: 100.0%
Average Response Time: 1.1ms

Individual Test Results:
  [1] security_001 - Score: 8.7/10 - ✅ PASS
  [2] security_002 - Score: 8.2/10 - ✅ PASS
  [3] practical_001 - Score: 8.3/10 - ✅ PASS

📈 Example 4: Baseline Comparison
---------------------------------
Baseline vs Current Comparison:
Score Change: +0.90 points
Pass Rate Change: +33.3%
Response Time Change: -500ms
🎉 SIGNIFICANT IMPROVEMENT!

✅ All examples completed!
```

## 🏗️ **Key Components Demonstrated**

### **MockEngine Implementation**
Shows how to create a simple engine for testing:
```go
type MockEngine struct{}

func (m *MockEngine) Chat(messages []chanollama.Message) (*chanollama.Response, error) {
    // Generate contextual mock responses based on query content
    query := messages[len(messages)-1].Content
    response := generateMockResponse(query)
    return &chanollama.Response{Content: response}, nil
}
```

### **Custom Scoring Function**
Example of domain-specific quality enhancement:
```go
securityScoringFunc := func(response string, metrics chanollama.EvalMetrics, testCase *chanollama.TestCase) float64 {
    score := chanollama.DefaultScoringFunction(response, metrics, testCase)
    
    // Add domain-specific bonuses
    securityTerms := []string{"encryption", "authentication", "secure"}
    for _, term := range securityTerms {
        if containsIgnoreCase(response, term) {
            score += 0.3
        }
    }
    
    return clampScore(score)
}
```

### **Test Suite Structure**
Demonstrates comprehensive test case creation:
```go
suite := chanollama.EvalSuite{
    Name: "Example Security Test Suite",
    TestCases: []chanollama.TestCase{
        {
            ID:               "security_001",
            Query:            "How do I implement secure authentication?",
            Category:         "security",
            MinLength:        200,
            MustInclude:      []string{"authentication", "secure"},
            ExpectedKeywords: []string{"token", "password", "certificate"},
        },
    },
}
```

## 🔧 **Adapting for Your Project**

### **1. Replace MockEngine with Your LLM**
```go
// Instead of MockEngine, use your real LLM implementation
engine, err := yourproject.NewLLMEngine(config)
evaluator := chanollama.NewEvaluator(engine, false)
```

### **2. Create Domain-Specific Scoring**
```go
func YourDomainScoringFunction(response string, metrics chanollama.EvalMetrics, testCase *chanollama.TestCase) float64 {
    score := chanollama.DefaultScoringFunction(response, metrics, testCase)
    
    // Add your domain-specific logic
    if containsDomainTerms(response) { score += 0.5 }
    if missingCriticalInfo(response, testCase) { score -= 1.0 }
    
    return score
}
```

### **3. Build Your Test Suite**
```go
yourTestSuite := chanollama.EvalSuite{
    Name: "Your Project Quality Tests",
    TestCases: []chanollama.TestCase{
        {
            ID:           "domain_001", 
            Query:        "Your domain-specific question",
            Category:     "your_category",
            MinLength:    300,
            MustInclude:  []string{"required", "terms"},
        },
    },
}
```

### **4. Integrate with CI/CD**
```go
// Save results for CI comparison
summary, _ := evaluator.EvaluateSuite(yourTestSuite)
saveResultsToFile(summary, "evaluation-results.json")

// Compare with baseline
baseline, _ := loadBaselineFromFile("baseline.json") 
comparison := chanollama.CompareResults(baseline, summary)
if comparison.IsRegression {
    os.Exit(1) // Fail CI build
}
```

## 📊 **Quality Metrics Explained**

### **Core Metrics (0-10 scale)**
- **Helpfulness**: Practical value, actionable guidance, relevance to query
- **Completeness**: Thoroughness, context coverage, comprehensive explanations  
- **Actionability**: Code examples, commands, step-by-step instructions
- **Length**: Appropriate detail level without excessive verbosity

### **Content Analysis**
- **Response Length**: Character count for minimum quality thresholds
- **Word Count**: Vocabulary richness and response depth
- **Code Examples**: Practical implementation guidance
- **Commands**: Actionable instructions and tools usage

### **Performance Metrics**
- **Response Time**: Latency measurement for performance benchmarking
- **Token Usage**: Resource consumption tracking
- **Success Rate**: Percentage of queries that complete successfully

## 🎯 **Use Cases for This Framework**

### **Development Workflow**
- **Prompt Engineering**: Test prompt changes systematically
- **Model Comparison**: Evaluate different LLMs on same test suite
- **Quality Assurance**: Ensure consistent response quality over time

### **Production Monitoring**
- **Regression Detection**: Catch quality drops in CI/CD pipelines
- **Performance Tracking**: Monitor response time and success rates
- **User Satisfaction**: Correlate evaluation scores with user feedback

### **Research & Optimization**
- **A/B Testing**: Compare different system configurations
- **Fine-tuning Validation**: Measure improvement from model customization
- **Benchmark Creation**: Establish quality standards for your domain

## 💡 **Best Practices**

### **Test Design**
- ✅ **Cover edge cases** in addition to common queries
- ✅ **Use realistic queries** that match actual user needs
- ✅ **Set appropriate thresholds** based on your quality requirements
- ✅ **Include both positive and negative test cases**

### **Scoring Functions**
- ✅ **Start with default scoring** and add domain-specific enhancements
- ✅ **Weight bonuses and penalties** appropriately for your use case
- ✅ **Test scoring changes** against known good/bad responses
- ✅ **Document scoring rationale** for team understanding

### **Baseline Management**
- ✅ **Establish baselines** before making significant changes
- ✅ **Update baselines** when intentional improvements are made
- ✅ **Archive historical results** for trend analysis
- ✅ **Set regression thresholds** appropriate for your stability needs

This evaluation framework provides the foundation for systematic LLM quality assurance that scales with your project's needs!