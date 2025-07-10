# Stopwords Configuration Examples

This example demonstrates the flexible stopwords configuration in chanollama's RAG system.

## Features

1. **Embedded Default Stopwords**: Uses a comprehensive list of 1,298 stopwords embedded in the binary
2. **Additional Stopwords**: Append domain-specific words to the default list
3. **Custom Override**: Replace the entire stopwords list with your own
4. **No Stopwords**: Disable stopword filtering entirely

## Usage

```go
import "github.com/jrschumacher/chanollama"

// Use default embedded stopwords
keywords := chanollama.ExtractKeywords(text)

// Add additional stopwords
config := chanollama.StopwordsConfig{
    UseDefault: true,
    Additional: []string{"domain", "specific", "terms"},
}
keywords := chanollama.ExtractKeywordsWithConfig(text, config)

// Use completely custom stopwords
config := chanollama.StopwordsConfig{
    Override: []string{"the", "a", "an"},
}
keywords := chanollama.ExtractKeywordsWithConfig(text, config)
```

## Benefits

- **Performance**: Embedded stopwords load instantly without file I/O
- **Flexibility**: Easy to customize for domain-specific use cases
- **Comprehensive**: 1,298 stopwords covering common English words, contractions, and technical terms
- **Backward Compatible**: Existing code using `ExtractKeywords()` continues to work