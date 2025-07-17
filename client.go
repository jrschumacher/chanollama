package chanollama

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ollama/ollama/llama"
)

// Client implements the Engine interface for direct model execution
type Client struct {
	config       Config
	modelPath    string
	model        *llama.Model
	context      *llama.Context
	mu           sync.RWMutex
	initialized  bool
	ctx          context.Context
	cancel       context.CancelFunc
}

// New creates a new client with the specified model path
func New(modelPath string) *Client {
	return NewWithConfig(modelPath, DefaultConfig(modelPath))
}

// NewWithConfig creates a new client with custom configuration
func NewWithConfig(modelPath string, config Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Client{
		config:      config,
		modelPath:   modelPath,
		ctx:         ctx,
		cancel:      cancel,
		initialized: false,
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
		devNull.Close()
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
		devNull.Close()
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
		devNull.Close()
		os.Stderr = oldStderr
	}
	
	if err != nil {
		llama.FreeModel(model)
		return fmt.Errorf("failed to create context: %w", err)
	}
	c.context = context
	
	c.initialized = true
	// Model loaded successfully - suppress verbose logs for cleaner chat experience
	
	return nil
}

// Chat performs a non-streaming chat request
func (c *Client) Chat(messages []Message) (*Response, error) {
	if !c.initialized {
		if err := c.Initialize(); err != nil {
			return &Response{Error: err}, err
		}
	}
	
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	start := time.Now()
	
	// Build prompt from messages
	prompt := c.buildPrompt(messages)
	
	// Perform inference
	content, err := c.performInference(prompt)
	if err != nil {
		return &Response{Error: err}, err
	}
	
	duration := time.Since(start)
	
	return &Response{
		Content:      content,
		Tokens:       c.estimateTokens(content),
		Duration:     duration,
		Model:        c.config.Model,
		Timestamp:    time.Now(),
		FinishReason: "stop",
	}, nil
}

// ChatStream performs a streaming chat request
func (c *Client) ChatStream(messages []Message, callback StreamCallback) *Response {
	if !c.initialized {
		if err := c.Initialize(); err != nil {
			return &Response{Error: err}
		}
	}
	
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	start := time.Now()
	
	// Build prompt from messages
	prompt := c.buildPrompt(messages)
	
	// Perform streaming inference
	content, err := c.performStreamingInference(prompt, callback)
	if err != nil {
		return &Response{Error: err}
	}
	
	duration := time.Since(start)
	
	return &Response{
		Content:      content,
		Tokens:       c.estimateTokens(content),
		Duration:     duration,
		Model:        c.config.Model,
		Timestamp:    time.Now(),
		FinishReason: "stop",
	}
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
	maxTokens := 1024  // Increased default for longer responses
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
		if c.shouldStop(response.String()) {
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
	
	return c.cleanResponse(response.String()), nil
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
	maxTokens := 1024  // Increased default for longer responses
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
		if c.shouldStop(response.String()) {
			break
		}
		
		// Clean and call streaming callback
		if callback != nil {
			cleanedPiece := c.cleanStreamingToken(piece)
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
	
	return c.cleanResponse(response.String()), nil
}

// cleanResponse removes ChatML tokens and cleans up the response
func (c *Client) cleanResponse(response string) string {
	// Remove ChatML tokens that might appear in the response
	cleaned := response
	
	// Remove any ChatML tokens that the model might have generated
	chatMLTokens := []string{
		"<|im_start|>",
		"<|im_end|>",
		"<|eot_id|>",
		"<|begin_of_text|>",
		"<|end_of_text|>",
		"<|eom_id|>",
	}
	
	for _, token := range chatMLTokens {
		cleaned = strings.ReplaceAll(cleaned, token, "")
	}
	
	// Remove role indicators that might leak through
	roleIndicators := []string{
		"user\n",
		"assistant\n",
		"system\n",
		"User:\n",
		"Assistant:\n",
		"System:\n",
		"User: ",
		"Assistant: ",
		"System: ",
	}
	
	for _, indicator := range roleIndicators {
		// Only remove if it appears at the beginning or after a newline
		if strings.HasPrefix(cleaned, indicator) {
			cleaned = strings.TrimPrefix(cleaned, indicator)
		}
		cleaned = strings.ReplaceAll(cleaned, "\n"+indicator, "\n")
	}
	
	// Clean up extra whitespace and newlines
	cleaned = strings.TrimSpace(cleaned)
	
	// Remove trailing conversation artifacts
	cleanedLines := []string{}
	lines := strings.Split(cleaned, "\n")
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip empty lines and lines that look like conversation artifacts
		if trimmed == "" || strings.HasPrefix(trimmed, "<|") || strings.HasSuffix(trimmed, "|>") {
			continue
		}
		
		// Skip standalone role indicators that appear at the end
		if i == len(lines)-1 && (trimmed == "User" || trimmed == "Assistant" || 
		                         trimmed == "Human" || trimmed == "AI" ||
		                         trimmed == "User:" || trimmed == "Assistant:" ||
		                         trimmed == "Human:" || trimmed == "AI:") {
			continue
		}
		
		cleanedLines = append(cleanedLines, line)
	}
	
	result := strings.Join(cleanedLines, "\n")
	
	// Remove trailing conversation artifacts that might be at the very end
	result = strings.TrimSuffix(result, "\nUser")
	result = strings.TrimSuffix(result, "\nAssistant")
	result = strings.TrimSuffix(result, "\nHuman")
	result = strings.TrimSuffix(result, "\nAI")
	result = strings.TrimSuffix(result, "\n\nUser")
	result = strings.TrimSuffix(result, "\n\nAssistant")
	result = strings.TrimSuffix(result, "\n\nHuman")
	result = strings.TrimSuffix(result, "\n\nAI")
	
	return strings.TrimSpace(result)
}

// cleanStreamingToken cleans individual tokens during streaming to prevent ChatML artifacts
func (c *Client) cleanStreamingToken(token string) string {
	// Filter out ChatML tokens as they're generated
	chatMLTokens := []string{
		"<|im_start|>",
		"<|im_end|>",
		"<|eot_id|>",
		"<|begin_of_text|>",
		"<|end_of_text|>",
		"<|eom_id|>",
	}
	
	for _, chatToken := range chatMLTokens {
		if strings.Contains(token, chatToken) {
			return ""
		}
	}
	
	// Filter out role indicators that might appear at the start of tokens
	roleIndicators := []string{
		"user\n", "assistant\n", "system\n",
		"User:\n", "Assistant:\n", "System:\n",
		"User: ", "Assistant: ", "System: ",
	}
	
	for _, indicator := range roleIndicators {
		if strings.HasPrefix(token, indicator) {
			return ""
		}
	}
	
	// Also filter out partial ChatML tokens that might appear across multiple pieces
	if strings.Contains(token, "<|") || strings.Contains(token, "|>") {
		return ""
	}
	
	return token
}

// createSamplingContext creates a sampling context with optional log suppression
func (c *Client) createSamplingContext(samplingParams llama.SamplingParams) (*llama.SamplingContext, error) {
	if !c.config.Verbose {
		oldStderr := os.Stderr
		devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		os.Stderr = devNull
		
		sampler, err := llama.NewSamplingContext(c.model, samplingParams)
		
		// Restore stderr
		devNull.Close()
		os.Stderr = oldStderr
		
		return sampler, err
	} else {
		return llama.NewSamplingContext(c.model, samplingParams)
	}
}

// shouldStop checks if the response contains patterns that indicate self-conversation
func (c *Client) shouldStop(response string) bool {
	// More specific stop patterns that clearly indicate model is starting a new conversation turn
	stopPatterns := []string{
		"<|eot_id|>",
		"<|end_of_text|>",
		"<|im_start|>",
		"<|im_end|>",
	}
	
	for _, pattern := range stopPatterns {
		if strings.Contains(response, pattern) {
			return true
		}
	}
	
	// Check for conversation patterns, but be more careful about context
	lines := strings.Split(response, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Only stop if we see a clear conversation pattern after some content
		if i > 0 && (strings.HasPrefix(trimmed, "User:") || 
		             strings.HasPrefix(trimmed, "Assistant:") ||
		             strings.HasPrefix(trimmed, "Human:") ||
		             strings.HasPrefix(trimmed, "AI:")) {
			// Check if there's actual content before this line
			hasContent := false
			for j := 0; j < i; j++ {
				if strings.TrimSpace(lines[j]) != "" {
					hasContent = true
					break
				}
			}
			if hasContent {
				return true
			}
		}
	}
	
	// Stop if we see standalone "User" or "Assistant" on its own line after content
	if strings.Contains(response, "\n\nUser\n") || 
	   strings.Contains(response, "\n\nAssistant\n") ||
	   strings.HasSuffix(response, "\n\nUser") ||
	   strings.HasSuffix(response, "\n\nAssistant") ||
	   strings.HasSuffix(response, "\nUser") ||
	   strings.HasSuffix(response, "\nAssistant") {
		return true
	}
	
	return false
}

// buildPrompt converts chat messages to a prompt string using a simple format
func (c *Client) buildPrompt(messages []Message) string {
	var prompt strings.Builder
	
	for _, msg := range messages {
		switch msg.Role {
		case "system":
			prompt.WriteString(fmt.Sprintf("System: %s\n\n", msg.Content))
		case "user":
			prompt.WriteString(fmt.Sprintf("User: %s\n\n", msg.Content))
		case "assistant":
			prompt.WriteString(fmt.Sprintf("Assistant: %s\n\n", msg.Content))
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

// UpdateConfig updates the configuration (requires reinitialization)
func (c *Client) UpdateConfig(config Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.config = config
	// Note: Changing config requires reinitialization for some parameters
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