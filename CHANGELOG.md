# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Release-please automation for version management
- GitHub Actions workflows for testing and releases
- Comprehensive offline-only architecture enforcement
- Channel-based streaming with Go concurrency primitives
- Direct GGUF model loading via llama.cpp bindings
- Advanced conversation management with context handling
- Token sanitization and output formatting
- Multi-agent framework with intelligent routing
- RAG (Retrieval-Augmented Generation) capabilities
- Evaluation tools for model performance testing
- Comprehensive test suite (78+ tests)

### Features
- **Pure Offline Operation**: Zero network dependencies, direct model file execution
- **Channel-Based Streaming**: Real-time token delivery via Go channels
- **Conversation Management**: Context handling, message truncation, auto-summarization
- **Agent Framework**: Multi-agent support with specialized roles
- **RAG Support**: Vector and keyword search capabilities
- **Evaluation Tools**: Model performance testing and comparison
- **Production Ready**: Comprehensive error handling and resource management

### Technical Architecture
- Direct llama.cpp integration via `github.com/ollama/ollama/llama`
- GGUF/GGML model file support with validation
- Go concurrency patterns for streaming and processing
- Comprehensive type system with fluent APIs
- Advanced memory management and context optimization

## [0.9.0] - 2024-09-11

### Added
- Initial release with offline LLM inference capabilities
- Core client and engine architecture
- Basic streaming and conversation support