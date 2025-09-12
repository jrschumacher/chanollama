# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1](https://github.com/jrschumacher/chanollama/compare/v1.0.0...v1.0.1) (2025-09-12)


### Bug Fixes

* remove domain bias ([72a4fd8](https://github.com/jrschumacher/chanollama/commit/72a4fd8bb6a3addb760a1b9bf340f65e3b19aeb3))

## 1.0.0 (2025-09-12)


### Features

* first commit ([13ababe](https://github.com/jrschumacher/chanollama/commit/13ababe26e90a8a088cc2daad0f1fe0da09018b6))
* first release of the library ([d0007a8](https://github.com/jrschumacher/chanollama/commit/d0007a878483733ae7644ac084a91c10803c9c35))


### Bug Fixes

* agent evaluation removing keywords ([83800b8](https://github.com/jrschumacher/chanollama/commit/83800b883122169d2ddd7182aed1ab73ac5edbb1))
* bump ([effc6c5](https://github.com/jrschumacher/chanollama/commit/effc6c599cb737ce8fa74ccd738f44a46fc679b3))
* lint errors ([a859b10](https://github.com/jrschumacher/chanollama/commit/a859b1015285b41889b59d47c3fb253dbe3e3643))
* lint issues ([e17b70b](https://github.com/jrschumacher/chanollama/commit/e17b70b8dc2804cd82d5f0abcd61542e43c4caaa))
* lint issues ([f5464af](https://github.com/jrschumacher/chanollama/commit/f5464af61c0471c5890b2c8e42bb77a104b4947d))
* lint issues ([16367e3](https://github.com/jrschumacher/chanollama/commit/16367e3940aca4ec90252e8c94fe8b340b444f9f))
* linter issue ([60e20d5](https://github.com/jrschumacher/chanollama/commit/60e20d54d02e5a9ca45f6e7c6a1891465a808686))
* remove custom lint file ([8bb4489](https://github.com/jrschumacher/chanollama/commit/8bb448984335634e38fb2c8deb896e44c99debfd))

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
