# Chanollama Examples Test Suite

This directory contains examples demonstrating various features of the chanollama library, along with an automated test suite to verify functionality.

## Running Examples

All examples now require the `CHANOLLAMA_MODEL_PATH` environment variable to be set to a valid GGUF model file path.

### Individual Example

```bash
export CHANOLLAMA_MODEL_PATH="/path/to/your/model.gguf"
cd basic && go run main.go
```

### All Examples (Automated Testing)

```bash
# Run the test suite
export CHANOLLAMA_MODEL_PATH="/path/to/your/model.gguf"
./run_tests.sh
```

Or run the test runner directly:

```bash
export CHANOLLAMA_MODEL_PATH="/path/to/your/model.gguf"
go run test_examples.go
```

## Test System Architecture

### Test Metadata

Each example includes embedded test metadata in comments at the top of the file:

```go
// TEST_METADATA:
// name: basic_chat
// description: Basic chat functionality test
// goal: Verify client can initialize and get a non-empty response
// expected_output: Response: [non-empty string]
// timeout: 30s
```

### Functional Goal Evaluation

The test system evaluates **functional goals** rather than response quality:

- ✅ **Functional**: Does the example run without errors?
- ✅ **Goal Met**: Does the output indicate the intended functionality worked?
- ❌ **Not Evaluated**: Response quality, accuracy, or content correctness

### Test Categories

1. **Model-Required Tests** - Need `CHANOLLAMA_MODEL_PATH` set
   - `basic_chat` - Basic chat functionality
   - `streaming_chat` - Streaming responses
   - `engine_streaming` - Engine-based streaming
   - `system_prompts` - System prompt handling
   - `context_awareness` - Context-aware conversations
   - `concurrent_streaming` - Concurrent request handling
   - `advanced_config` - Advanced configuration
   - `conversation_history` - Conversation history management

2. **Model-Independent Tests** - Don't require external model
   - `stopwords_extraction` - Keyword extraction utilities
   - `evaluation_framework` - Evaluation framework demo (uses mock)

## Available Examples

| Example | Description | Features Demonstrated |
|---------|-------------|----------------------|
| `basic/` | Simple chat request | Basic client usage |
| `streaming/` | Streaming responses | Callback-based streaming |
| `02_streaming/` | Engine streaming | Engine interface, health checks |
| `03_system_prompts/` | System prompts | Role-based interactions |
| `04_context_aware/` | Context awareness | Conversation management |
| `06_concurrent_streaming/` | Concurrent requests | Concurrent streaming |
| `advanced/` | Advanced config | Custom engine configuration |
| `conversation/` | Conversation history | Multi-turn conversations |
| `evaluation/` | Evaluation framework | Testing and evaluation |
| `stopwords/` | Keyword extraction | Text processing utilities |

## Adding New Examples

1. Create a new directory with `main.go`
2. Add test metadata comments at the top:
   ```go
   // TEST_METADATA:
   // name: example_name
   // description: Brief description
   // goal: What functionality should be verified
   // expected_output: What output pattern to look for
   // timeout: 30s
   ```

3. Use environment variable for model path:
   ```go
   modelPath := os.Getenv("CHANOLLAMA_MODEL_PATH")
   if modelPath == "" {
       log.Fatal("CHANOLLAMA_MODEL_PATH environment variable is required")
   }
   ```

4. Add evaluation logic to `test_examples.go` if needed

## Test Output Format

```
🧪 Chanollama Examples Test Suite
=================================
📁 Using model: /path/to/model.gguf
📋 Found 10 examples to test

[1/10] Running basic_chat...
✅ basic_chat - Basic chat response received (2.34s)

[2/10] Running streaming_chat...
✅ streaming_chat - Streaming callback executed (3.12s)

📊 Test Results Summary
======================
Total: 10 | Passed: 10 | Failed: 0
Success Rate: 100.0%

🎉 All tests passed!
```

## Troubleshooting

### Common Issues

1. **Model Path Not Set**
   ```
   ❌ CHANOLLAMA_MODEL_PATH not set (required for this example)
   ```
   Solution: Set the environment variable to your model file path

2. **Model File Not Found**
   ```
   ❌ Model file not found: /path/to/model.gguf
   ```
   Solution: Verify the model file exists and path is correct

3. **Timeout Errors**
   ```
   ❌ command timed out after 30s
   ```
   Solution: Use a faster model or increase timeout in test metadata

4. **Functional Goal Not Met**
   ```
   ❌ basic_chat - No valid response found
   ```
   Solution: Check if the model is working properly or adjust evaluation criteria

5. **Token Batch Size Errors**
   ```
   GGML_ASSERT(n_tokens_all <= cparams.n_batch) failed
   ```
   Solution: Reduce conversation length or use a model with larger batch size

6. **Memory Issues**
   ```
   SIGABRT: abort during cgo execution
   ```
   Solution: Use a smaller model or reduce concurrent requests

### Performance Tips

- Use quantized models (Q4, Q5) for faster inference
- Consider smaller models for testing (7B vs 70B parameters)
- Set appropriate timeouts based on your hardware
- Run tests with GPU acceleration if available

## Design Philosophy

This test system focuses on **functional verification** rather than output quality:

- **Fast Execution**: Tests run quickly to enable frequent validation
- **Reliable Results**: Consistent pass/fail based on functional goals
- **Easy Maintenance**: Simple to add new examples with minimal setup
- **Model Agnostic**: Works with any GGUF model that provides basic functionality

The goal is to ensure all examples **work as intended** without evaluating the **quality** of LLM responses, which can vary significantly between models and configurations.

## Future Roadmap

### High Priority Improvements

**Better Error Handling**
- Graceful handling of batch size limits with automatic context trimming
- Clear error messages for configuration issues
- Automatic retry logic for transient failures
- Model compatibility validation before execution

**Enhanced Configuration Management**
- Auto-detect optimal batch size based on model capabilities
- Model-specific configuration profiles (context length, temperature)
- Intelligent context window management for long conversations
- Hardware-aware optimization (CPU vs GPU, memory limits)
- Configuration validation and warnings for suboptimal settings

### Medium Priority Features

**Improved Test Infrastructure**
- Performance benchmarking across different models
- Memory usage monitoring and alerts
- Automated model compatibility testing
- CI/CD integration for continuous validation

**Advanced Model Management**
- Model metadata detection (context length, batch size, capabilities)
- Automatic model format conversion utilities
- Model performance profiling and optimization suggestions
- Support for model quantization detection and handling

### Low Priority Enhancements

**Developer Experience**
- Interactive test runner with model selection
- Real-time performance monitoring dashboard
- Example templates for common use cases
- Integration with popular model repositories

**Documentation and Tooling**
- Model compatibility matrix
- Performance comparison guides
- Best practices for different hardware configurations
- Troubleshooting decision trees

---

*This roadmap prioritizes reliability and developer experience while maintaining the core principle of functional verification over response quality evaluation.*