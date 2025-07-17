package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TestMetadata represents the metadata extracted from each example
type TestMetadata struct {
	Name            string
	Description     string
	Goal            string
	ExpectedOutput  string
	Timeout         time.Duration
	FilePath        string
	RequiresModel   bool
}

// TestResult represents the result of running an example
type TestResult struct {
	Metadata       TestMetadata
	Success        bool
	Output         string
	Error          string
	Duration       time.Duration
	FunctionalGoal string
}

func main() {
	fmt.Println("🧪 Chanollama Examples Test Suite")
	fmt.Println("=================================")

	// Check if model path is set (required for most examples)
	modelPath := os.Getenv("CHANOLLAMA_MODEL_PATH")
	if modelPath == "" {
		fmt.Println("⚠️  Warning: CHANOLLAMA_MODEL_PATH not set - some examples may fail")
	} else {
		fmt.Printf("📁 Using model: %s\n", modelPath)
	}

	// Find all example files
	exampleFiles, err := findExampleFiles()
	if err != nil {
		log.Fatalf("Failed to find example files: %v", err)
	}

	fmt.Printf("📋 Found %d examples to test\n\n", len(exampleFiles))

	// Parse metadata for each example
	var tests []TestMetadata
	for _, file := range exampleFiles {
		metadata, err := parseMetadata(file)
		if err != nil {
			log.Printf("⚠️  Failed to parse metadata for %s: %v", file, err)
			continue
		}
		tests = append(tests, metadata)
	}

	// Run all tests
	results := make([]TestResult, 0, len(tests))
	passed := 0
	failed := 0

	for i, test := range tests {
		fmt.Printf("[%d/%d] Running %s...\n", i+1, len(tests), test.Name)
		
		result := runTest(test)
		results = append(results, result)
		
		if result.Success {
			fmt.Printf("✅ %s - %s (%.2fs)\n", test.Name, result.FunctionalGoal, result.Duration.Seconds())
			passed++
		} else {
			fmt.Printf("❌ %s - %s\n", test.Name, result.Error)
			failed++
		}
		fmt.Println()
	}

	// Print summary
	fmt.Println("📊 Test Results Summary")
	fmt.Println("======================")
	fmt.Printf("Total: %d | Passed: %d | Failed: %d\n", len(results), passed, failed)
	fmt.Printf("Success Rate: %.1f%%\n", float64(passed)/float64(len(results))*100)

	if failed > 0 {
		fmt.Println("\n❌ Failed Tests:")
		for _, result := range results {
			if !result.Success {
				fmt.Printf("  - %s: %s\n", result.Metadata.Name, result.Error)
			}
		}
		os.Exit(1)
	}

	fmt.Println("\n🎉 All tests passed!")
}

// findExampleFiles finds all main.go files in example directories
func findExampleFiles() ([]string, error) {
	var files []string
	
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
			mainFile := filepath.Join(entry.Name(), "main.go")
			if _, err := os.Stat(mainFile); err == nil {
				files = append(files, mainFile)
			}
		}
	}

	return files, nil
}

// parseMetadata extracts test metadata from example file comments
func parseMetadata(filePath string) (TestMetadata, error) {
	metadata := TestMetadata{
		FilePath:      filePath,
		Timeout:       30 * time.Second, // default timeout
		RequiresModel: true,             // assume most examples require model
	}

	file, err := os.Open(filePath)
	if err != nil {
		return metadata, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inMetadata := false
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if strings.Contains(line, "TEST_METADATA:") {
			inMetadata = true
			continue
		}
		
		if inMetadata && !strings.HasPrefix(line, "//") {
			break
		}
		
		if inMetadata && strings.HasPrefix(line, "//") {
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimSpace(line)
			
			if strings.HasPrefix(line, "name:") {
				metadata.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
			} else if strings.HasPrefix(line, "description:") {
				metadata.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
			} else if strings.HasPrefix(line, "goal:") {
				metadata.Goal = strings.TrimSpace(strings.TrimPrefix(line, "goal:"))
			} else if strings.HasPrefix(line, "expected_output:") {
				metadata.ExpectedOutput = strings.TrimSpace(strings.TrimPrefix(line, "expected_output:"))
			} else if strings.HasPrefix(line, "timeout:") {
				timeoutStr := strings.TrimSpace(strings.TrimPrefix(line, "timeout:"))
				if timeout, err := time.ParseDuration(timeoutStr); err == nil {
					metadata.Timeout = timeout
				}
			}
		}
	}

	// Special handling for examples that don't require model
	if metadata.Name == "stopwords_extraction" || metadata.Name == "evaluation_framework" {
		metadata.RequiresModel = false
	}

	return metadata, scanner.Err()
}

// runTest executes an example and evaluates its functional goal
func runTest(metadata TestMetadata) TestResult {
	result := TestResult{
		Metadata: metadata,
		Success:  false,
	}

	// Check if model is required and available
	if metadata.RequiresModel {
		modelPath := os.Getenv("CHANOLLAMA_MODEL_PATH")
		if modelPath == "" {
			result.Error = "CHANOLLAMA_MODEL_PATH not set (required for this example)"
			return result
		}
		if _, err := os.Stat(modelPath); os.IsNotExist(err) {
			result.Error = fmt.Sprintf("Model file not found: %s", modelPath)
			return result
		}
	}

	// Run the example with timeout
	start := time.Now()
	
	cmd := exec.Command("go", "run", metadata.FilePath)
	cmd.Dir = "."
	
	// Capture output
	output, err := runWithTimeout(cmd, metadata.Timeout)
	result.Duration = time.Since(start)
	result.Output = output

	if err != nil {
		result.Error = fmt.Sprintf("Execution failed: %v", err)
		return result
	}

	// Evaluate functional goal
	success, goalResult := evaluateFunctionalGoal(metadata, output)
	result.Success = success
	result.FunctionalGoal = goalResult

	return result
}

// runWithTimeout executes a command with a timeout
func runWithTimeout(cmd *exec.Cmd, timeout time.Duration) (string, error) {
	output := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		out, err := cmd.CombinedOutput()
		if err != nil {
			errChan <- err
			return
		}
		output <- string(out)
	}()

	select {
	case out := <-output:
		return out, nil
	case err := <-errChan:
		return "", err
	case <-time.After(timeout):
		cmd.Process.Kill()
		return "", fmt.Errorf("command timed out after %v", timeout)
	}
}

// evaluateFunctionalGoal evaluates if the example met its functional goal
func evaluateFunctionalGoal(metadata TestMetadata, output string) (bool, string) {
	switch metadata.Name {
	case "basic_chat":
		return evaluateBasicChat(output)
	case "streaming_chat":
		return evaluateStreamingChat(output)
	case "engine_streaming":
		return evaluateEngineStreaming(output)
	case "system_prompts":
		return evaluateSystemPrompts(output)
	case "context_awareness":
		return evaluateContextAwareness(output)
	case "concurrent_streaming":
		return evaluateConcurrentStreaming(output)
	case "advanced_config":
		return evaluateAdvancedConfig(output)
	case "conversation_history":
		return evaluateConversationHistory(output)
	case "evaluation_framework":
		return evaluateEvaluationFramework(output)
	case "stopwords_extraction":
		return evaluateStopwordsExtraction(output)
	default:
		return evaluateGeneric(output)
	}
}

// Individual evaluation functions for each example type

func evaluateBasicChat(output string) (bool, string) {
	if strings.Contains(output, "Response:") && len(strings.TrimSpace(output)) > 20 {
		return true, "Basic chat response received"
	}
	return false, "No valid response found"
}

func evaluateStreamingChat(output string) (bool, string) {
	if strings.Contains(output, "Streaming response:") {
		return true, "Streaming callback executed"
	}
	return false, "No streaming output detected"
}

func evaluateEngineStreaming(output string) (bool, string) {
	checks := []string{"Engine ready", "callback-based streaming", "Total tokens:"}
	passed := 0
	for _, check := range checks {
		if strings.Contains(output, check) {
			passed++
		}
	}
	if passed >= 2 {
		return true, fmt.Sprintf("Engine streaming functional (%d/3 checks)", passed)
	}
	return false, "Engine streaming checks failed"
}

func evaluateSystemPrompts(output string) (bool, string) {
	roles := []string{"Coding Assistant", "Documentation Expert", "Creative Writer", "Debugging Assistant"}
	found := 0
	for _, role := range roles {
		if strings.Contains(output, role) {
			found++
		}
	}
	if found >= 2 {
		return true, fmt.Sprintf("System prompts functional (%d/4 roles)", found)
	}
	return false, "System prompts failed"
}

func evaluateContextAwareness(output string) (bool, string) {
	checks := []string{"Conversation Statistics", "Messages:", "Estimated tokens:"}
	passed := 0
	for _, check := range checks {
		if strings.Contains(output, check) {
			passed++
		}
	}
	if passed >= 2 {
		return true, "Context awareness functional"
	}
	return false, "Context awareness checks failed"
}

func evaluateConcurrentStreaming(output string) (bool, string) {
	checks := []string{"conversations completed", "Q&A Demo", "Performance Test"}
	passed := 0
	for _, check := range checks {
		if strings.Contains(output, check) {
			passed++
		}
	}
	if passed >= 2 {
		return true, fmt.Sprintf("Sequential streaming functional (%d/3 checks)", passed)
	}
	return false, "Sequential streaming checks failed"
}

func evaluateAdvancedConfig(output string) (bool, string) {
	if strings.Contains(output, "advanced configuration example completed") {
		return true, "Advanced configuration functional"
	}
	// Check for streaming output and completion
	if strings.Contains(output, "Streaming response") && len(strings.TrimSpace(output)) > 200 {
		return true, "Advanced configuration streaming functional"
	}
	return false, "Advanced configuration failed"
}

func evaluateConversationHistory(output string) (bool, string) {
	// Look for multiple "Assistant:" responses indicating conversation turns
	assistantCount := strings.Count(output, "Assistant:")
	if assistantCount >= 3 {
		return true, fmt.Sprintf("Conversation history functional (%d turns)", assistantCount)
	}
	return false, "Insufficient conversation turns"
}

func evaluateEvaluationFramework(output string) (bool, string) {
	checks := []string{"Single Query Evaluation", "Custom Scoring", "Test Suite", "Baseline Comparison"}
	passed := 0
	for _, check := range checks {
		if strings.Contains(output, check) {
			passed++
		}
	}
	if passed >= 3 {
		return true, "Evaluation framework functional"
	}
	return false, "Evaluation framework checks failed"
}

func evaluateStopwordsExtraction(output string) (bool, string) {
	checks := []string{"Default Stopwords", "Additional Stopwords", "Custom Stopwords", "Keywords:"}
	passed := 0
	for _, check := range checks {
		if strings.Contains(output, check) {
			passed++
		}
	}
	if passed >= 3 {
		return true, "Stopwords extraction functional"
	}
	return false, "Stopwords extraction checks failed"
}

func evaluateGeneric(output string) (bool, string) {
	// Generic evaluation: check for reasonable output length and no obvious errors
	if len(strings.TrimSpace(output)) > 50 && !strings.Contains(output, "panic:") && !strings.Contains(output, "fatal error") {
		return true, "Generic output validation passed"
	}
	return false, "Generic output validation failed"
}