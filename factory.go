package chanollama

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewEngine creates an appropriate engine based on configuration
func NewEngine(config Config) (Engine, error) {
	// Auto-detect offline mode if model path looks like a file path
	if !config.OfflineMode && config.ModelPath == "" {
		// Check if Model looks like a file path (contains / or ends with .gguf)
		if strings.Contains(config.Model, "/") || strings.HasSuffix(config.Model, ".gguf") {
			if _, err := os.Stat(config.Model); err == nil {
				config.OfflineMode = true
				config.ModelPath = config.Model
			}
		}
	}
	
	// Use offline mode if configured or model path is provided
	if config.OfflineMode || config.ModelPath != "" {
		return NewOfflineEngine(config)
	}
	
	// Default to offline client with model path
	return NewWithConfig(config.ModelPath, config), nil
}

// NewOfflineEngine creates an offline engine with the given configuration
func NewOfflineEngine(config Config) (Engine, error) {
	modelPath := config.ModelPath
	if modelPath == "" {
		modelPath = config.Model
	}
	
	// Validate model path
	if err := validateModelPath(modelPath); err != nil {
		return nil, fmt.Errorf("invalid model path: %w", err)
	}
	
	// Create offline client
	client := NewWithConfig(modelPath, config)
	
	return client, nil
}


// validateModelPath checks if the model path is valid
func validateModelPath(modelPath string) error {
	if modelPath == "" {
		return fmt.Errorf("model path cannot be empty")
	}
	
	// Check if file exists
	info, err := os.Stat(modelPath)
	if err != nil {
		return fmt.Errorf("model file not found: %s", modelPath)
	}
	
	// Check if it's a directory
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}
	
	// Check if it's a supported format or an Ollama blob
	ext := strings.ToLower(filepath.Ext(modelPath))
	supportedFormats := []string{".gguf", ".bin", ".ggml"}
	
	// Check for supported extensions
	for _, format := range supportedFormats {
		if ext == format {
			return nil
		}
	}
	
	// Check if it's an Ollama blob file (no extension but in .ollama/models/blobs/)
	if ext == "" && strings.Contains(modelPath, ".ollama/models/blobs/") {
		// Additional validation for Ollama blobs - check file size
		if info.Size() > 1024 { // At least 1KB
			return nil
		}
		return fmt.Errorf("ollama blob file appears to be empty or too small")
	}
	
	return fmt.Errorf("unsupported model format: %s (supported: %v, or Ollama blob files)", ext, supportedFormats)
}

// IsModelFile checks if the given string looks like a model file path
func IsModelFile(model string) bool {
	// Check for file path indicators
	if strings.Contains(model, "/") || strings.Contains(model, "\\") {
		return true
	}
	
	// Check for model file extensions
	ext := strings.ToLower(filepath.Ext(model))
	return ext == ".gguf" || ext == ".bin" || ext == ".ggml"
}

// AutoDetectMode determines the best mode for the given model parameter
func AutoDetectMode(model string) (bool, string) {
	// Returns (isOffline, modelPath)
	
	if IsModelFile(model) {
		if _, err := os.Stat(model); err == nil {
			return true, model
		}
	}
	
	// Default to HTTP mode
	return false, ""
}

