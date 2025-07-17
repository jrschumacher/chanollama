#!/bin/bash

# Chanollama Examples Test Runner
# This script runs all example tests and provides detailed output

set -e

echo "🧪 Chanollama Examples Test Suite"
echo "================================="

# Check if model path is provided
if [ -z "$CHANOLLAMA_MODEL_PATH" ]; then
    echo "⚠️  Warning: CHANOLLAMA_MODEL_PATH not set"
    echo "   Some examples will be skipped or may fail"
    echo "   Usage: CHANOLLAMA_MODEL_PATH=/path/to/model.gguf ./run_tests.sh"
    echo ""
    
    # Ask if user wants to continue
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Tests cancelled"
        exit 1
    fi
else
    echo "📁 Using model: $CHANOLLAMA_MODEL_PATH"
    
    # Verify model file exists
    if [ ! -f "$CHANOLLAMA_MODEL_PATH" ]; then
        echo "❌ Model file not found: $CHANOLLAMA_MODEL_PATH"
        exit 1
    fi
fi

echo ""

# Change to examples directory
cd "$(dirname "$0")"

# Clean up any previous test artifacts
echo "🧹 Cleaning up previous test artifacts..."
rm -f test_examples

# Build the test runner
echo "🔨 Building test runner..."
go build -o test_examples test_examples.go

# Run the tests
echo "🚀 Running tests..."
echo ""

# Export the model path for child processes
export CHANOLLAMA_MODEL_PATH

# Run the actual tests
./test_examples

# Clean up
echo ""
echo "🧹 Cleaning up..."
rm -f test_examples

echo "✅ Test run completed!"