# 📚 chanollama Documentation

Welcome to the comprehensive documentation for chanollama - a powerful Go library for streaming interactions with Ollama models.

## 📋 Table of Contents

### 🚀 Getting Started
- [Installation & Quick Start](./getting-started.md)
- [Configuration](./configuration.md)
- [Basic Examples](./examples/basic.md)

### 🔧 Core Features
- [Streaming](./streaming.md) - Real-time token streaming with channels
- [Context Management](./context-management.md) - Conversation history and memory
- [Message Types](./message-types.md) - System, user, and assistant messages

### 🧠 Advanced Features
- **[RAG (Retrieval-Augmented Generation)](./rag.md)** - Document retrieval and context enhancement
- [Agents](./agents.md) - Multi-agent systems and intelligent routing
- [Stopwords](./stopwords.md) - Flexible keyword extraction configuration
- [Validation](./validation.md) - Input sanitization and security

### 📖 API Reference
Use Go's built-in documentation tools:
```bash
# View package documentation
go doc github.com/jrschumacher/chanollama

# View specific types/functions
go doc chanollama.VectorStore
go doc chanollama.ExtractKeywords

# Start local documentation server
godoc -http=:6060
# Then visit: http://localhost:6060/pkg/github.com/jrschumacher/chanollama/
```

### 💡 Examples & Tutorials
- [Building a Chatbot](./examples/chatbot.md)
- [RAG-powered Search](./examples/rag-search.md)
- [Multi-Agent System](./examples/multi-agent.md)
- [Streaming Dashboard](./examples/streaming-dashboard.md)

### 🔍 Deep Dives
- [Architecture Overview](./architecture.md)
- [Performance Optimization](./performance.md)
- [Security Best Practices](./security.md)
- [Troubleshooting](./troubleshooting.md)

## 🎯 Quick Navigation

| What you want to do | Where to look |
|---------------------|---------------|
| Get started quickly | [Getting Started](./getting-started.md) |
| Add document search | **[RAG Documentation](./rag.md)** |
| Build a multi-agent system | [Agents](./agents.md) |
| Stream responses in real-time | [Streaming](./streaming.md) |
| Configure custom stopwords | [Stopwords](./stopwords.md) |
| Secure your implementation | [Security](./security.md) |

## 🤝 Contributing

See the main [README](../README.md) for contribution guidelines and development setup.