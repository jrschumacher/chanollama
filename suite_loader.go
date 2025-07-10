package chanollama

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SuiteLoader handles loading and validation of test suites
type SuiteLoader struct {
	searchPaths []string
}

// NewSuiteLoader creates a new suite loader with default search paths
func NewSuiteLoader() *SuiteLoader {
	return &SuiteLoader{
		searchPaths: []string{
			".", // Current directory
			"evals", // Common evals directory
			"tests", // Common tests directory
			"test_suites", // Alternative name
		},
	}
}

// NewSuiteLoaderWithPaths creates a suite loader with custom search paths
func NewSuiteLoaderWithPaths(paths []string) *SuiteLoader {
	return &SuiteLoader{
		searchPaths: paths,
	}
}

// LoadSuite loads a test suite from file with automatic path resolution
func (sl *SuiteLoader) LoadSuite(filePath string) (EvalSuite, error) {
	resolvedPath, err := sl.resolveSuitePath(filePath)
	if err != nil {
		return EvalSuite{}, err
	}
	
	return sl.loadSuiteFromPath(resolvedPath)
}

// LoadSuiteFromPath loads a test suite from a specific file path
func (sl *SuiteLoader) LoadSuiteFromPath(filePath string) (EvalSuite, error) {
	return sl.loadSuiteFromPath(filePath)
}

// resolveSuitePath resolves the full path to a test suite file
func (sl *SuiteLoader) resolveSuitePath(filePath string) (string, error) {
	// If it's already an absolute path, use it directly
	if filepath.IsAbs(filePath) {
		if _, err := os.Stat(filePath); err != nil {
			return "", fmt.Errorf("test suite file not found: %s", filePath)
		}
		return filePath, nil
	}
	
	// Search in configured paths
	for _, searchPath := range sl.searchPaths {
		candidatePath := filepath.Join(searchPath, filePath)
		if _, err := os.Stat(candidatePath); err == nil {
			absPath, err := filepath.Abs(candidatePath)
			if err != nil {
				return "", fmt.Errorf("failed to resolve absolute path: %w", err)
			}
			return absPath, nil
		}
	}
	
	// Also try common extensions if not provided
	baseFilePath := filePath
	if !strings.HasSuffix(filePath, ".json") {
		baseFilePath = filePath + ".json"
		
		for _, searchPath := range sl.searchPaths {
			candidatePath := filepath.Join(searchPath, baseFilePath)
			if _, err := os.Stat(candidatePath); err == nil {
				absPath, err := filepath.Abs(candidatePath)
				if err != nil {
					return "", fmt.Errorf("failed to resolve absolute path: %w", err)
				}
				return absPath, nil
			}
		}
	}
	
	return "", fmt.Errorf("test suite file not found in search paths: %s (searched: %v)", 
		filePath, sl.searchPaths)
}

// loadSuiteFromPath loads and validates a test suite from a specific path
func (sl *SuiteLoader) loadSuiteFromPath(filePath string) (EvalSuite, error) {
	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return EvalSuite{}, fmt.Errorf("failed to read test suite file %s: %w", filePath, err)
	}
	
	// Parse JSON
	var suite EvalSuite
	if err := json.Unmarshal(data, &suite); err != nil {
		return EvalSuite{}, fmt.Errorf("failed to parse test suite JSON in %s: %w", filePath, err)
	}
	
	// Validate suite structure
	if err := sl.validateSuite(suite); err != nil {
		return EvalSuite{}, fmt.Errorf("test suite validation failed for %s: %w", filePath, err)
	}
	
	return suite, nil
}

// validateSuite validates the structure and content of a test suite
func (sl *SuiteLoader) validateSuite(suite EvalSuite) error {
	if suite.Name == "" {
		return fmt.Errorf("test suite must have a name")
	}
	
	if len(suite.TestCases) == 0 {
		return fmt.Errorf("test suite must contain at least one test case")
	}
	
	// Validate each test case
	seenIDs := make(map[string]bool)
	for i, testCase := range suite.TestCases {
		if err := sl.validateTestCase(testCase, i); err != nil {
			return fmt.Errorf("test case %d validation failed: %w", i+1, err)
		}
		
		// Check for duplicate IDs
		if testCase.ID != "" {
			if seenIDs[testCase.ID] {
				return fmt.Errorf("duplicate test case ID: %s", testCase.ID)
			}
			seenIDs[testCase.ID] = true
		}
	}
	
	return nil
}

// validateTestCase validates an individual test case
func (sl *SuiteLoader) validateTestCase(testCase TestCase, index int) error {
	if testCase.Query == "" {
		return fmt.Errorf("test case must have a query")
	}
	
	if len(testCase.Query) > 10000 {
		return fmt.Errorf("test case query too long (max 10000 characters)")
	}
	
	// Validate constraints
	if testCase.MinLength < 0 {
		return fmt.Errorf("min_length cannot be negative")
	}
	
	// Note: MaxLength validation removed as TestCase doesn't have MaxLength field
	// This could be added in the future if needed
	
	// Validate category if provided
	if testCase.Category != "" && !sl.isValidCategory(testCase.Category) {
		// Warning only - don't fail validation for unknown categories
		// This allows for domain-specific categories
	}
	
	return nil
}

// isValidCategory checks if a category is in the standard set
func (sl *SuiteLoader) isValidCategory(category string) bool {
	standardCategories := []string{
		"general", "technical", "policy", "development", 
		"deployment", "troubleshooting", "security",
		"code_example", "policy_design", "agent_routing",
	}
	
	for _, validCategory := range standardCategories {
		if category == validCategory {
			return true
		}
	}
	
	return false
}

// AddSearchPath adds a new search path for test suite files
func (sl *SuiteLoader) AddSearchPath(path string) {
	sl.searchPaths = append(sl.searchPaths, path)
}

// GetSearchPaths returns the current search paths
func (sl *SuiteLoader) GetSearchPaths() []string {
	return append([]string{}, sl.searchPaths...) // Return a copy
}

// LoadSuiteFromFile is a convenience function that uses the default loader
func LoadSuiteFromFile(filePath string) (EvalSuite, error) {
	loader := NewSuiteLoader()
	return loader.LoadSuite(filePath)
}

// ValidateSuiteFile validates a test suite file without fully loading it
func ValidateSuiteFile(filePath string) error {
	loader := NewSuiteLoader()
	_, err := loader.LoadSuite(filePath)
	return err
}

// ListAvailableSuites lists all JSON files in the search paths
func (sl *SuiteLoader) ListAvailableSuites() ([]string, error) {
	var suites []string
	seen := make(map[string]bool)
	
	for _, searchPath := range sl.searchPaths {
		if _, err := os.Stat(searchPath); err != nil {
			continue // Skip non-existent directories
		}
		
		entries, err := os.ReadDir(searchPath)
		if err != nil {
			continue // Skip unreadable directories
		}
		
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			
			if strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
				suitePath := filepath.Join(searchPath, entry.Name())
				if !seen[suitePath] {
					suites = append(suites, suitePath)
					seen[suitePath] = true
				}
			}
		}
	}
	
	return suites, nil
}

// GetSuiteInfo returns basic information about a test suite without loading all test cases
func (sl *SuiteLoader) GetSuiteInfo(filePath string) (SuiteInfo, error) {
	resolvedPath, err := sl.resolveSuitePath(filePath)
	if err != nil {
		return SuiteInfo{}, err
	}
	
	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return SuiteInfo{}, fmt.Errorf("failed to read test suite file: %w", err)
	}
	
	// Parse just the top-level fields
	var partial struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Version     string     `json:"version"`
		TestCases   []TestCase `json:"test_cases"`
	}
	
	if err := json.Unmarshal(data, &partial); err != nil {
		return SuiteInfo{}, fmt.Errorf("failed to parse test suite JSON: %w", err)
	}
	
	return SuiteInfo{
		Name:        partial.Name,
		Description: partial.Description,
		Version:     partial.Version,
		TestCount:   len(partial.TestCases),
		FilePath:    resolvedPath,
	}, nil
}

// SuiteInfo contains basic information about a test suite
type SuiteInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	TestCount   int    `json:"test_count"`
	FilePath    string `json:"file_path"`
}