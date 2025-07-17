package chanollama

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	modelPath := "/path/to/model.gguf"
	client := New(modelPath)
	
	if client == nil {
		t.Fatal("New returned nil")
	}
	
	if client.GetModelPath() != modelPath {
		t.Errorf("Expected model path %s, got %s", modelPath, client.GetModelPath())
	}
}

func TestNewClientWithConfig(t *testing.T) {
	config := Config{
		ModelPath:    "/path/to/model.gguf",
		Temperature:  0.8,
		MaxTokens:    500,
	}
	
	client := NewWithConfig("/path/to/model.gguf", config)
	if client == nil {
		t.Fatal("NewWithConfig returned nil")
	}
	
	clientConfig := client.GetConfig()
	
	if clientConfig.Temperature != 0.8 {
		t.Errorf("Expected temperature 0.8, got %f", clientConfig.Temperature)
	}
	
	if clientConfig.MaxTokens != 500 {
		t.Errorf("Expected max tokens 500, got %d", clientConfig.MaxTokens)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig("/path/to/model.gguf")
	
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", config.Timeout)
	}
	
	if config.Temperature != 0.7 {
		t.Errorf("Expected default temperature 0.7, got %f", config.Temperature)
	}
}

func TestDefaultEngineConfig(t *testing.T) {
	config := DefaultEngineConfig()
	
	if config.RequestChannelSize <= 0 {
		t.Errorf("Expected positive request channel size, got %d", config.RequestChannelSize)
	}
	
	if config.StreamingChannelSize <= 0 {
		t.Errorf("Expected positive streaming channel size, got %d", config.StreamingChannelSize)
	}
}