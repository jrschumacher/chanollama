package chanollama

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ollama/ollama/llama"
)

// Client implements the Engine interface for direct model execution
type Client struct {
	config              Config
	modelPath           string
	model               *llama.Model
	context             *llama.Context
	sanitizer           *TokenSanitizer
	conversationManager *ConversationManager
	mu                  sync.RWMutex
	initialized         bool
	ctx                 context.Context
	cancel              context.CancelFunc
}

// New creates a new client with the specified model path
func New(modelPath string) *Client {
	return NewWithConfig(modelPath, DefaultConfig(modelPath))
}

// NewWithConfig creates a new client with custom configuration
func NewWithConfig(modelPath string, config Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	// Create conversation manager with appropriate batch size
	conversationConfig := DefaultConversationConfig()
	if config.ContextSize > 0 {
		// Set token limit to 80% of context size to prevent batch size issues
		conversationConfig.MaxTokens = int(float64(config.ContextSize) * 0.8)
		conversationConfig.BatchSize = config.ContextSize
		// Disable auto-detection since we're setting explicit values
		conversationConfig.AutoDetectLimits = false
	}

	return &Client{
		config:              config,
		modelPath:           modelPath,
		sanitizer:           NewTokenSanitizer(),
		conversationManager: NewConversationManagerWithConfig(conversationConfig),
		ctx:                 ctx,
		cancel:              cancel,
		initialized:         false,
	}
}

// Initialize loads the model and sets up the context
func (c *Client) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return nil
	}

	// Verify model file exists
	if _, err := os.Stat(c.modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s", c.modelPath)
	}

	// Suppress verbose loading logs for cleaner chat experience (unless verbose mode is enabled)
	var oldStderr *os.File
	var devNull *os.File

	if !c.config.Verbose {
		// Redirect verbose llama.cpp logs to prevent noise in chat
		oldStderr = os.Stderr
		devNull, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		os.Stderr = devNull
	}

	// Initialize llama backend
	llama.BackendInit()

	// Restore stderr for our application logs
	if !c.config.Verbose {
		_ = devNull.Close()
		os.Stderr = oldStderr
	}

	// Set up model parameters
	modelParams := llama.ModelParams{
		NumGpuLayers: 0, // TODO: Add GPU support detection
		UseMmap:      true,
		VocabOnly:    false,
	}

	// Load model (redirect stderr to suppress verbose output)
	if !c.config.Verbose {
		oldStderr = os.Stderr
		devNull, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		os.Stderr = devNull
	}

	model, err := llama.LoadModelFromFile(c.modelPath, modelParams)

	// Restore stderr
	if !c.config.Verbose {
		_ = devNull.Close()
		os.Stderr = oldStderr
	}

	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}
	c.model = model

	// Create context with configuration
	contextSize := 4096
	if c.config.MaxTokens > 0 {
		contextSize = c.config.MaxTokens
	}

	contextParams := llama.NewContextParams(
		contextSize, // numCtx
		512,         // batchSize
		1,           // numSeqMax
		4,           // threads
		false,       // flashAttention
		"",          // kvCacheType
	)

	// Suppress context creation logs
	if !c.config.Verbose {
		oldStderr = os.Stderr
		devNull, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		os.Stderr = devNull
	}

	context, err := llama.NewContextWithModel(model, contextParams)

	// Restore stderr
	if !c.config.Verbose {
		_ = devNull.Close()
		os.Stderr = oldStderr
	}

	if err != nil {
		llama.FreeModel(model)
		return fmt.Errorf("failed to create context: %w", err)
	}
	c.context = context

	c.initialized = true
	// Model loaded successfully - suppress verbose logs for cleaner chat experience

	// Auto-detect conversation limits if enabled
	if c.conversationManager.config.AutoDetectLimits {
		if err := c.AutoDetectConversationLimits(); err != nil {
			if c.config.Verbose {
				log.Printf("Failed to auto-detect conversation limits: %v", err)
			}
		}
	}

	return nil
}

// Chat performs a non-streaming chat request
func (c *Client) Chat(messages []Message) (*Response, error) {
	return c.safeInference(messages, nil)
}

// ChatStream performs a streaming chat request
func (c *Client) ChatStream(messages []Message, callback StreamCallback) *Response {
	response, err := c.safeInference(messages, callback)
	if err != nil {
		return &Response{Error: err}
	}
	return response
}

// safeInference wraps inference calls with panic recovery and error handling
func (c *Client) safeInference(messages []Message, callback StreamCallback) (*Response, error) {
	defer func() {
		if r := recover(); r != nil {
			// Log the panic and return graceful error
			if c.config.Verbose {
				log.Printf("LLM inference panic recovered: %v", r)
			}
		}
	}()

	if !c.initialized {
		if err := c.Initialize(); err != nil {
			wrappedErr := fmt.Errorf("%w: %v", ErrModelNotLoaded, err)
			return &Response{Error: wrappedErr}, wrappedErr
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	start := time.Now()

	// Truncate messages to prevent batch size crashes
	truncatedMessages := c.conversationManager.TruncateMessages(messages)

	// Check if truncation occurred due to batch size limits
	if len(truncatedMessages) < len(messages) {
		estimatedTokens := c.conversationManager.EstimateTokens(messages)
		if estimatedTokens > c.config.ContextSize {
			if c.config.Verbose {
				log.Printf("Context overflow detected: %d tokens > %d limit, truncated to %d messages",
					estimatedTokens, c.config.ContextSize, len(truncatedMessages))
			}
		}
	}

	// Build prompt from truncated messages
	prompt := c.buildPrompt(truncatedMessages)

	// Perform inference based on whether streaming is requested
	var content string
	var err error

	if callback != nil {
		// Streaming inference
		content, err = c.performStreamingInference(prompt, callback)
	} else {
		// Non-streaming inference
		content, err = c.performInference(prompt)
	}

	if err != nil {
		wrappedErr := fmt.Errorf("%w: %v", ErrInferenceFailed, err)
		return &Response{Error: wrappedErr}, wrappedErr
	}

	duration := time.Since(start)
	now := time.Now()

	response := &Response{
		Content:      content,
		Tokens:       c.estimateTokens(content),
		Duration:     duration,
		Model:        c.config.Model,
		Timestamp:    &now,
		FinishReason: "stop",
	}

	// Add truncation info to metadata if messages were truncated
	if len(truncatedMessages) < len(messages) {
		if response.Metadata == nil {
			response.Metadata = make(map[string]string)
		}
		response.Metadata["truncated_messages"] = fmt.Sprintf("%d", len(messages)-len(truncatedMessages))
		response.Metadata["original_message_count"] = fmt.Sprintf("%d", len(messages))
		response.Metadata["truncated_message_count"] = fmt.Sprintf("%d", len(truncatedMessages))
		response.Metadata["truncation_reason"] = "batch_size_limit"
	}

	return response, nil
}

// IsHealthy checks if the client is ready for inference
func (c *Client) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.initialized && c.model != nil && c.context != nil
}

// Close cleans up resources
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}

	if c.model != nil {
		llama.FreeModel(c.model)
		c.model = nil
	}

	c.context = nil
	c.initialized = false

	return nil
}

// performInference runs actual model inference using Ollama's llama bindings
func (c *Client) performInference(prompt string) (string, error) {
	// Tokenize the prompt
	tokens, err := c.model.Tokenize(prompt, true, true)
	if err != nil {
		return "", fmt.Errorf("tokenization failed: %w", err)
	}

	// Create batch for processing
	batch, err := llama.NewBatch(len(tokens), 1, 0)
	if err != nil {
		return "", fmt.Errorf("batch creation failed: %w", err)
	}
	defer batch.Free()

	// Add tokens to batch
	for i, token := range tokens {
		batch.Add(token, nil, i, i == len(tokens)-1, 0) // Only get logits for last token
	}

	// Process the batch
	err = c.context.Decode(batch)
	if err != nil {
		return "", fmt.Errorf("context decode failed: %w", err)
	}

	// Set up sampling parameters
	samplingParams := llama.SamplingParams{
		TopK:           40,
		TopP:           0.9,
		MinP:           0.1,
		Temp:           float32(c.config.Temperature),
		RepeatLastN:    64,
		PenaltyRepeat:  1.1,
		PenaltyFreq:    0.0,
		PenaltyPresent: 0.0,
		PenalizeNl:     true,
		Seed:           0,
	}

	// Create sampling context
	sampler, err := c.createSamplingContext(samplingParams)
	if err != nil {
		return "", fmt.Errorf("sampling context creation failed: %w", err)
	}

	var response strings.Builder
	maxTokens := 1024 // Increased default for longer responses
	if c.config.MaxTokens > 0 {
		maxTokens = c.config.MaxTokens
	}

	// Generate tokens iteratively
	for i := 0; i < maxTokens; i++ {
		// Sample next token
		token := sampler.Sample(c.context, batch.NumTokens()-1)

		// Check for end of generation
		if c.model.TokenIsEog(token) {
			break
		}

		// Convert token to text
		piece := c.model.TokenToPiece(token)
		response.WriteString(piece)

		// Check for stop sequences to prevent self-conversation
		if c.sanitizer.ShouldStop(response.String()) {
			break
		}

		// Accept the token for grammar/repetition tracking
		sampler.Accept(token, true)

		// Prepare for next iteration - add token to batch
		batch.Clear()
		batch.Add(token, nil, len(tokens)+i, true, 0)

		// Decode for next iteration
		err = c.context.Decode(batch)
		if err != nil {
			// Decode failed during generation - suppress verbose logs
			break
		}
	}

	return c.sanitizer.CleanResponse(response.String()), nil
}

// performStreamingInference runs inference with streaming callback
func (c *Client) performStreamingInference(prompt string, callback StreamCallback) (string, error) {
	// Tokenize the prompt
	tokens, err := c.model.Tokenize(prompt, true, true)
	if err != nil {
		return "", fmt.Errorf("tokenization failed: %w", err)
	}

	// Create batch for processing
	batch, err := llama.NewBatch(len(tokens), 1, 0)
	if err != nil {
		return "", fmt.Errorf("batch creation failed: %w", err)
	}
	defer batch.Free()

	// Add tokens to batch
	for i, token := range tokens {
		batch.Add(token, nil, i, i == len(tokens)-1, 0) // Only get logits for last token
	}

	// Process the batch
	err = c.context.Decode(batch)
	if err != nil {
		return "", fmt.Errorf("context decode failed: %w", err)
	}

	// Set up sampling parameters
	samplingParams := llama.SamplingParams{
		TopK:           40,
		TopP:           0.9,
		MinP:           0.1,
		Temp:           float32(c.config.Temperature),
		RepeatLastN:    64,
		PenaltyRepeat:  1.1,
		PenaltyFreq:    0.0,
		PenaltyPresent: 0.0,
		PenalizeNl:     true,
		Seed:           0,
	}

	// Create sampling context
	sampler, err := c.createSamplingContext(samplingParams)
	if err != nil {
		return "", fmt.Errorf("sampling context creation failed: %w", err)
	}

	var response strings.Builder
	maxTokens := 1024 // Increased default for longer responses
	if c.config.MaxTokens > 0 {
		maxTokens = c.config.MaxTokens
	}

	// Generate tokens iteratively with streaming
	for i := 0; i < maxTokens; i++ {
		// Sample next token
		token := sampler.Sample(c.context, batch.NumTokens()-1)

		// Check for end of generation
		if c.model.TokenIsEog(token) {
			break
		}

		// Convert token to text
		piece := c.model.TokenToPiece(token)
		response.WriteString(piece)

		// Check for stop sequences to prevent self-conversation
		if c.sanitizer.ShouldStop(response.String()) {
			break
		}

		// Clean and call streaming callback
		if callback != nil {
			cleanedPiece := c.sanitizer.CleanStreamingToken(piece)
			if cleanedPiece != "" {
				callback(cleanedPiece)
			}
		}

		// Accept the token for grammar/repetition tracking
		sampler.Accept(token, true)

		// Prepare for next iteration - add token to batch
		batch.Clear()
		batch.Add(token, nil, len(tokens)+i, true, 0)

		// Decode for next iteration
		err = c.context.Decode(batch)
		if err != nil {
			// Decode failed during generation - suppress verbose logs
			break
		}
	}

	return c.sanitizer.CleanResponse(response.String()), nil
}

// createSamplingContext creates a sampling context with optional log suppression
func (c *Client) createSamplingContext(samplingParams llama.SamplingParams) (*llama.SamplingContext, error) {
	if !c.config.Verbose {
		oldStderr := os.Stderr
		devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		os.Stderr = devNull

		sampler, err := llama.NewSamplingContext(c.model, samplingParams)

		// Restore stderr
		_ = devNull.Close()
		os.Stderr = oldStderr

		return sampler, err
	} else {
		return llama.NewSamplingContext(c.model, samplingParams)
	}
}

// buildPrompt converts chat messages to a prompt string using a simple format
func (c *Client) buildPrompt(messages []Message) string {
	var prompt strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			fmt.Fprintf(&prompt, "System: %s\n\n", msg.Content)
		case "user":
			fmt.Fprintf(&prompt, "User: %s\n\n", msg.Content)
		case "assistant":
			fmt.Fprintf(&prompt, "Assistant: %s\n\n", msg.Content)
		}
	}

	// Add the assistant prompt to start generation
	prompt.WriteString("Assistant: ")

	return prompt.String()
}

// estimateTokens provides a rough token count estimate
func (c *Client) estimateTokens(content string) int {
	// Rough estimation: 1 token ≈ 4 characters
	return len(content) / 4
}

// GetModelPath returns the path to the model file
func (c *Client) GetModelPath() string {
	return c.modelPath
}

// GetConfig returns the current configuration
func (c *Client) GetConfig() Config {
	return c.config
}

// GetModelInfo returns information about the loaded model
func (c *Client) GetModelInfo() (*ModelInfo, error) {
	if !c.initialized {
		return nil, ErrModelNotLoaded
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.model == nil {
		return nil, ErrModelNotLoaded
	}

	// Get context size from configuration (as set during initialization)
	contextSize := c.config.ContextSize
	if contextSize == 0 {
		contextSize = 4096 // Default context size
	}

	// Estimate batch size based on context size and configuration
	// Most models use batch size that's a fraction of context size
	batchSize := contextSize / 8 // Conservative estimate
	if batchSize < 512 {
		batchSize = 512 // Minimum reasonable batch size
	}
	if batchSize > 2048 {
		batchSize = 2048 // Maximum typical batch size
	}

	return &ModelInfo{
		VocabSize:   -1, // Not available without proper API
		ContextSize: contextSize,
		BatchSize:   batchSize,
		ModelPath:   c.modelPath,
	}, nil
}

// AutoDetectConversationLimits automatically configures conversation limits based on model capabilities
func (c *Client) AutoDetectConversationLimits() error {
	modelInfo, err := c.GetModelInfo()
	if err != nil {
		return err
	}

	config := c.conversationManager.config

	// Auto-detect batch size if not set
	if config.BatchSize == 0 || config.AutoDetectLimits {
		config.BatchSize = modelInfo.BatchSize
	}

	// Auto-detect max tokens if not set (80% of batch size)
	if config.MaxTokens == 0 || config.AutoDetectLimits {
		config.MaxTokens = int(float64(modelInfo.BatchSize) * 0.8)
	}

	// Update the conversation manager
	c.conversationManager.SetConfig(config)

	if c.config.Verbose {
		log.Printf("Auto-detected conversation limits: MaxTokens=%d, BatchSize=%d",
			config.MaxTokens, config.BatchSize)
	}

	return nil
}

// UpdateConfig updates the configuration (requires reinitialization)
func (c *Client) UpdateConfig(config Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config = config
	// Note: Changing config requires reinitialization for some parameters
}

// GetConversationManager returns the conversation manager for advanced configuration
func (c *Client) GetConversationManager() *ConversationManager {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.conversationManager
}

// SetConversationConfig updates the conversation management configuration
func (c *Client) SetConversationConfig(config ConversationConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conversationManager.SetConfig(config)
}

// SetConversationLimits is a convenience method to set conversation limits
func (c *Client) SetConversationLimits(config ConversationConfig) {
	c.SetConversationConfig(config)
}

// GenerateEmbedding generates an embedding for the given text using the specified model
func (c *Client) GenerateEmbedding(text string, embeddingModel string) (*EmbeddingResponse, error) {
	// TODO: Implement offline embedding generation using llama.cpp
	// For now, return an error indicating this feature is not yet implemented
	return NewEmbeddingResponse(nil).WithError(fmt.Errorf("embedding generation not yet implemented for offline mode")),
		fmt.Errorf("embedding generation not yet implemented for offline mode")
}

// GenerateEmbeddings generates embeddings for multiple texts (batch operation)
func (c *Client) GenerateEmbeddings(texts []string, embeddingModel string) ([]*EmbeddingResponse, error) {
	// TODO: Implement offline batch embedding generation
	responses := make([]*EmbeddingResponse, len(texts))

	for i := range texts {
		responses[i] = NewEmbeddingResponse(nil).WithError(fmt.Errorf("embedding generation not yet implemented for offline mode"))
	}

	return responses, fmt.Errorf("embedding generation not yet implemented for offline mode")
}
