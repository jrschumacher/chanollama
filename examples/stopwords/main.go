package main

// TEST_METADATA:
// name: stopwords_extraction
// description: Keyword extraction with stopwords filtering
// goal: Verify stopwords functionality works without requiring model
// expected_output: Various keyword extraction results with different configurations
// timeout: 10s

import (
	"fmt"
	"github.com/jrschumacher/chanollama"
)

func main() {
	text := "The quick brown fox jumps over the lazy dog and runs through the forest"

	// Example 1: Default stopwords (from embedded file)
	fmt.Println("=== Default Stopwords ===")
	keywords1 := chanollama.ExtractKeywords(text)
	fmt.Printf("Input: %s\n", text)
	fmt.Printf("Keywords: %v\n\n", keywords1)

	// Example 2: Add additional stopwords
	fmt.Println("=== Default + Additional Stopwords ===")
	config2 := chanollama.StopwordsConfig{
		UseDefault: true,
		Additional: []string{"quick", "brown"},
	}
	keywords2 := chanollama.ExtractKeywordsWithConfig(text, config2)
	fmt.Printf("Keywords: %v\n\n", keywords2)

	// Example 3: Override with custom stopwords
	fmt.Println("=== Custom Stopwords Only ===")
	config3 := chanollama.StopwordsConfig{
		Override: []string{"the", "and", "over"},
	}
	keywords3 := chanollama.ExtractKeywordsWithConfig(text, config3)
	fmt.Printf("Keywords: %v\n\n", keywords3)

	// Example 4: No stopwords
	fmt.Println("=== No Stopwords ===")
	config4 := chanollama.StopwordsConfig{
		UseDefault: false,
	}
	keywords4 := chanollama.ExtractKeywordsWithConfig(text, config4)
	fmt.Printf("Keywords: %v\n\n", keywords4)
}