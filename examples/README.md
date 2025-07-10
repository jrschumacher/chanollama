# 🚀 Chanollama Examples

This directory contains examples demonstrating various features of the chanollama library.

## 📋 **Examples Overview**

- **`basic/`** - Simple setup and usage patterns
- **`streaming/`** - Real-time token streaming demonstrations  
- **`conversation/`** - Multi-turn conversation management
- **`advanced/`** - Complex use cases and configurations
- **`evaluation/`** - LLM quality assessment and testing framework ⭐
- **`stopwords/`** - Keyword extraction and text processing utilities

## 🧪 **Featured: Evaluation Framework**

The **`evaluation/`** example demonstrates chanollama's comprehensive LLM evaluation utilities:

- **Quality Metrics**: Helpfulness, completeness, actionability scoring (0-10 scale)
- **Custom Scoring**: Domain-specific quality assessment functions
- **Test Suites**: Systematic evaluation with multiple test cases and constraints
- **Baseline Comparison**: Regression detection and quality tracking over time
- **CI Integration**: Automated quality assurance workflows for production

Perfect for building quality assurance into any LLM-powered project!

## 🚀 **Quick Start**

```bash
# Run the evaluation framework example
go run examples/evaluation/main.go

# Run other examples
go run examples/basic/main.go
go run examples/streaming/main.go
```

## 🎯 **Use Cases**

### **Development**
- **Prompt Engineering**: Test and optimize system prompts systematically
- **Model Comparison**: Evaluate different LLMs on standardized test suites
- **Quality Assurance**: Ensure consistent response quality during development

### **Production**
- **Regression Detection**: Catch quality drops in CI/CD pipelines  
- **Performance Monitoring**: Track response time and success rates
- **Baseline Management**: Maintain quality standards over time

### **Research**
- **A/B Testing**: Compare different configurations and approaches
- **Benchmarking**: Establish quality standards for your domain
- **Optimization**: Measure improvements from fine-tuning or prompt changes

## 💡 **Best Practices**

When using these examples in your projects:

- ✅ **Start with basic examples** to understand core concepts
- ✅ **Adapt evaluation framework** for your specific domain and quality requirements
- ✅ **Use realistic test cases** that match your actual user queries
- ✅ **Establish baselines** before making significant changes to your system
- ✅ **Integrate quality checks** into your CI/CD pipeline for continuous assurance

Each example includes detailed documentation and comments to help you understand and adapt the code for your specific needs.