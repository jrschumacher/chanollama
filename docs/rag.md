# 📚 RAG (Retrieval-Augmented Generation)

RAG enhances language model responses by retrieving relevant context from a knowledge base before generating answers. chanollama provides both vector-based and keyword-based RAG implementations.

## 📖 Table of Contents

- [Overview](#overview)
- [Vector Store RAG](#vector-store-rag)
- [Simple RAG Store](#simple-rag-store)
- [Stopwords Configuration](#stopwords-configuration)
- [Document Management](#document-management)
- [Best Practices](#best-practices)
- [Examples](#examples)

## 🔍 Overview

chanollama offers two RAG approaches:

| Type | Best For | Features |
|------|----------|----------|
| **Vector Store** | Semantic similarity | Embedding-based search, cosine similarity |
| **Simple RAG** | Keyword matching | Fast, lightweight, no embeddings required |

## 🚀 Vector Store RAG

### Basic Setup

```go
import "github.com/jrschumacher/chanollama"

// Create vector store
vectorStore := chanollama.NewVectorStore("./vector_index.json")

// Create documents with embeddings
docs := []chanollama.Document{
    {
        ID:      "doc1",
        Title:   "Go Programming Guide",
        Content: "Go is a statically typed, compiled programming language...",
        URL:     "https://example.com/go-guide",
    },
    {
        ID:      "doc2", 
        Title:   "REST API Best Practices",
        Content: "When building REST APIs, follow these principles...",
        URL:     "https://example.com/api-guide",
    },
}

// Add documents (embeddings will be generated)
err := vectorStore.AddDocuments(docs)
if err != nil {
    log.Fatal(err)
}

// Save the index
err = vectorStore.SaveIndex()
if err != nil {
    log.Fatal(err)
}
```

### Querying with Vector Search

```go
// Search for relevant documents
query := "How do I build REST APIs in Go?"
results, err := vectorStore.SearchSimilar(query, 3) // Top 3 results
if err != nil {
    log.Fatal(err)
}

// Build context from results
var context strings.Builder
for _, result := range results {
    context.WriteString(fmt.Sprintf("Document: %s\n", result.Document.Title))
    context.WriteString(fmt.Sprintf("Content: %s\n", result.Document.Content))
    context.WriteString(fmt.Sprintf("Similarity: %.3f\n\n", result.Similarity))
}

// Use in chat
messages := []chanollama.Message{
    {
        Role: "system",
        Content: fmt.Sprintf("Use this context to answer questions:\n\n%s", context.String()),
    },
    {Role: "user", Content: query},
}

client := chanollama.New("llama2:7b")
response, err := client.Chat(messages)
```

### Vector Store API

```go
type VectorStore struct {
    documents []Document
    indexPath string
}

// Core methods
func NewVectorStore(indexPath string) *VectorStore
func (vs *VectorStore) AddDocument(doc Document) error
func (vs *VectorStore) AddDocuments(docs []Document) error
func (vs *VectorStore) SearchSimilar(query string, topK int) ([]SimilarityResult, error)
func (vs *VectorStore) LoadIndex() error
func (vs *VectorStore) SaveIndex() error
func (vs *VectorStore) GetDocumentCount() int

// Results
type SimilarityResult struct {
    Document   Document
    Similarity float32
}
```

## 🎯 Simple RAG Store

For keyword-based retrieval without embeddings:

### Basic Setup

```go
// Create simple RAG store
simpleRAG := chanollama.NewSimpleRAGStore("./simple_index.json")

// Add documents
documents := []chanollama.Document{
    {
        ID:      "go-basics",
        Title:   "Go Basics",
        Content: "Go programming language fundamentals, syntax, and best practices",
        Tags:    []string{"programming", "go", "tutorial"},
    },
    {
        ID:      "api-design",
        Title:   "API Design",
        Content: "RESTful API design principles, HTTP methods, and status codes",
        Tags:    []string{"api", "rest", "http", "design"},
    },
}

for _, doc := range documents {
    err := simpleRAG.AddDocument(doc)
    if err != nil {
        log.Fatal(err)
    }
}

// Save index
err := simpleRAG.SaveIndex()
if err != nil {
    log.Fatal(err)
}
```

### Querying with Keywords

```go
// Build RAG context for a query
query := "How to design REST APIs?"
ragContext := simpleRAG.BuildRAGContext(query, 3) // Top 3 relevant docs

// Use in conversation
messages := []chanollama.Message{
    {
        Role:    "system",
        Content: fmt.Sprintf("Context:\n%s\n\nUse this context to answer questions.", ragContext),
    },
    {Role: "user", Content: query},
}

client := chanollama.New("llama2:7b")
response, err := client.Chat(messages)
```

### Simple RAG API

```go
type SimpleRAGStore struct {
    documents []Document
    indexPath string
}

// Core methods
func NewSimpleRAGStore(indexPath string) *SimpleRAGStore
func (s *SimpleRAGStore) AddDocument(doc Document) error
func (s *SimpleRAGStore) SearchRelevant(query string, topK int) []Document
func (s *SimpleRAGStore) BuildRAGContext(query string, maxDocs int) string
func (s *SimpleRAGStore) LoadIndex() error
func (s *SimpleRAGStore) SaveIndex() error
func (s *SimpleRAGStore) GetDocumentCount() int
```

## 🔍 Stopwords Configuration

Customize keyword extraction for better RAG performance:

### Default Stopwords

```go
// Uses embedded 1,298 stopwords
keywords := chanollama.ExtractKeywords("The quick brown fox jumps over the lazy dog")
// Output: [quick, brown, fox, jumps, lazy, dog]
```

### Custom Stopwords

```go
// Add domain-specific stopwords
config := chanollama.StopwordsConfig{
    UseDefault: true,
    Additional: []string{"api", "endpoint", "service", "framework"},
}
keywords := chanollama.ExtractKeywordsWithConfig(text, config)

// Override with completely custom list
config = chanollama.StopwordsConfig{
    Override: []string{"the", "a", "an", "and", "or"},
}
keywords = chanollama.ExtractKeywordsWithConfig(text, config)

// Disable stopwords filtering
config = chanollama.StopwordsConfig{
    UseDefault: false,
}
keywords = chanollama.ExtractKeywordsWithConfig(text, config)
```

### Stopwords API

```go
type StopwordsConfig struct {
    UseDefault bool     // Use embedded default stopwords
    Additional []string // Additional words to filter
    Override   []string // Completely replace default stopwords
}

func DefaultStopwordsConfig() StopwordsConfig
func ExtractKeywords(text string) []string
func ExtractKeywordsWithConfig(text string, config StopwordsConfig) []string
```

## 📋 Document Management

### Document Structure

```go
type Document struct {
    ID          string            `json:"id"`
    Title       string            `json:"title"`
    Content     string            `json:"content"`
    URL         string            `json:"url"`
    Tags        []string          `json:"tags,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    Embedding   []float32         `json:"embedding,omitempty"`
    CreatedAt   time.Time         `json:"created_at,omitempty"`
    UpdatedAt   time.Time         `json:"updated_at,omitempty"`
}
```

### Chunking Large Documents

```go
// Split large documents into chunks
text := "Very long document content..."
chunks := chanollama.ChunkText(text, 200, 20) // 200 words, 20 word overlap

// Create documents from chunks
for i, chunk := range chunks {
    doc := chanollama.Document{
        ID:      fmt.Sprintf("doc-%d-chunk-%d", docID, i),
        Title:   fmt.Sprintf("%s (Part %d)", originalTitle, i+1),
        Content: chunk,
        Metadata: map[string]string{
            "parent_doc": originalDocID,
            "chunk_index": strconv.Itoa(i),
        },
    }
    
    err := vectorStore.AddDocument(doc)
    if err != nil {
        log.Printf("Failed to add chunk %d: %v", i, err)
    }
}
```

### Bulk Operations

```go
// Load documents from files
documents := make([]chanollama.Document, 0)

err := filepath.Walk("./docs", func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    
    if strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".txt") {
        content, err := os.ReadFile(path)
        if err != nil {
            return err
        }
        
        doc := chanollama.Document{
            ID:      strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
            Title:   filepath.Base(path),
            Content: string(content),
            URL:     path,
            Metadata: map[string]string{
                "file_path": path,
                "file_size": strconv.FormatInt(info.Size(), 10),
            },
            CreatedAt: info.ModTime(),
        }
        
        documents = append(documents, doc)
    }
    return nil
})

if err != nil {
    log.Fatal(err)
}

// Add all documents
err = vectorStore.AddDocuments(documents)
if err != nil {
    log.Fatal(err)
}
```

## 🎯 Best Practices

### 1. Choosing RAG Type

```go
// Use Vector Store for:
// - Semantic similarity search
// - Large document collections
// - Nuanced understanding needed

vectorStore := chanollama.NewVectorStore("./semantic_index.json")

// Use Simple RAG for:
// - Keyword-based search
// - Fast lookups
// - Resource-constrained environments

simpleRAG := chanollama.NewSimpleRAGStore("./keyword_index.json")
```

### 2. Document Preparation

```go
// Clean and normalize document content
func prepareDocument(rawContent string) string {
    // Remove extra whitespace
    content := strings.TrimSpace(rawContent)
    content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
    
    // Remove markdown formatting for cleaner text
    content = regexp.MustCompile(`[#*_`+"`"+`\[\]()]+`).ReplaceAllString(content, "")
    
    return content
}

doc := chanollama.Document{
    ID:      "clean-doc",
    Content: prepareDocument(rawContent),
    Title:   extractTitle(rawContent),
}
```

### 3. Context Window Management

```go
// Estimate token usage and manage context
func buildRAGResponse(query string, store *chanollama.SimpleRAGStore, maxTokens int) []chanollama.Message {
    // Reserve tokens for query and response
    availableTokens := maxTokens - chanollama.EstimateTokens(query) - 500 // 500 for response
    
    // Get relevant documents
    docs := store.SearchRelevant(query, 10) // Get more, then filter
    
    // Build context within token limit
    var context strings.Builder
    currentTokens := 0
    
    for _, doc := range docs {
        docTokens := chanollama.EstimateTokens(doc.Content)
        if currentTokens + docTokens > availableTokens {
            break
        }
        
        context.WriteString(fmt.Sprintf("Source: %s\n%s\n\n", doc.Title, doc.Content))
        currentTokens += docTokens
    }
    
    return []chanollama.Message{
        {
            Role:    "system",
            Content: fmt.Sprintf("Context (use for answering):\n%s", context.String()),
        },
        {Role: "user", Content: query},
    }
}
```

### 4. Performance Optimization

```go
// Preload indices for better performance
func initializeRAGSystem() (*chanollama.VectorStore, *chanollama.SimpleRAGStore, error) {
    // Load vector store
    vectorStore := chanollama.NewVectorStore("./vector_index.json")
    if err := vectorStore.LoadIndex(); err != nil {
        return nil, nil, fmt.Errorf("failed to load vector index: %v", err)
    }
    
    // Load simple store
    simpleRAG := chanollama.NewSimpleRAGStore("./simple_index.json")
    if err := simpleRAG.LoadIndex(); err != nil {
        return nil, nil, fmt.Errorf("failed to load simple index: %v", err)
    }
    
    log.Printf("RAG system initialized: %d vector docs, %d simple docs",
        vectorStore.GetDocumentCount(), simpleRAG.GetDocumentCount())
    
    return vectorStore, simpleRAG, nil
}

// Use connection pooling for multiple queries
func handleMultipleQueries(queries []string, store *chanollama.VectorStore) {
    var wg sync.WaitGroup
    resultsChan := make(chan struct{query string; results []chanollama.SimilarityResult}, len(queries))
    
    for _, query := range queries {
        wg.Add(1)
        go func(q string) {
            defer wg.Done()
            results, err := store.SearchSimilar(q, 3)
            if err != nil {
                log.Printf("Query failed: %v", err)
                return
            }
            resultsChan <- struct{query string; results []chanollama.SimilarityResult}{q, results}
        }(query)
    }
    
    wg.Wait()
    close(resultsChan)
    
    for result := range resultsChan {
        log.Printf("Query: %s, Found: %d documents", result.query, len(result.results))
    }
}
```

## 📝 Examples

### Complete RAG-powered Chatbot

```go
package main

import (
    "fmt"
    "log"
    "github.com/jrschumacher/chanollama"
)

func main() {
    // Initialize RAG system
    ragStore := chanollama.NewSimpleRAGStore("./knowledge.json")
    
    // Add knowledge base
    docs := []chanollama.Document{
        {
            ID: "go-error-handling",
            Title: "Go Error Handling",
            Content: "In Go, errors are values. Functions return error values to indicate an abnormal state...",
        },
        {
            ID: "go-concurrency",
            Title: "Go Concurrency",
            Content: "Go provides goroutines and channels for concurrent programming...",
        },
    }
    
    for _, doc := range docs {
        ragStore.AddDocument(doc)
    }
    ragStore.SaveIndex()
    
    // Initialize LLM client
    client := chanollama.New("llama2:7b")
    defer client.Close()
    
    // RAG-enhanced chat
    query := "How do I handle errors in Go?"
    
    // Get relevant context
    ragContext := ragStore.BuildRAGContext(query, 2)
    
    // Create messages with context
    messages := []chanollama.Message{
        {
            Role: "system",
            Content: fmt.Sprintf(`You are a Go programming expert. Use this context to answer questions:

%s

Provide detailed, practical answers based on the context above.`, ragContext),
        },
        {Role: "user", Content: query},
    }
    
    // Get response
    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Question: %s\n", query)
    fmt.Printf("Answer: %s\n", response.Content)
}
```

### Streaming RAG Responses

```go
func streamRAGResponse(query string, ragStore *chanollama.SimpleRAGStore, client *chanollama.Client) {
    // Build context
    ragContext := ragStore.BuildRAGContext(query, 3)
    
    messages := []chanollama.Message{
        {
            Role:    "system", 
            Content: fmt.Sprintf("Context:\n%s", ragContext),
        },
        {Role: "user", Content: query},
    }
    
    // Stream response
    fmt.Printf("🤖 Streaming response for: %s\n\n", query)
    
    response := client.ChatStream(messages, func(token string) {
        fmt.Print(token)
    })
    
    if response.Error != nil {
        log.Printf("Streaming error: %v", response.Error)
        return
    }
    
    fmt.Println("\n\n✅ Response complete")
}
```

## 🔧 Advanced Configuration

### Custom Similarity Functions

```go
// Custom cosine similarity with weights
func WeightedCosineSimilarity(a, b []float32, weights []float32) float32 {
    if len(a) != len(b) || len(a) != len(weights) {
        return 0.0
    }
    
    var dotProduct, normA, normB float32
    for i := range a {
        weightedA := a[i] * weights[i]
        weightedB := b[i] * weights[i]
        
        dotProduct += weightedA * weightedB
        normA += weightedA * weightedA
        normB += weightedB * weightedB
    }
    
    if normA == 0.0 || normB == 0.0 {
        return 0.0
    }
    
    return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
```

### Hybrid RAG (Vector + Keyword)

```go
type HybridRAG struct {
    vectorStore *chanollama.VectorStore
    simpleStore *chanollama.SimpleRAGStore
    vectorWeight float32
    keywordWeight float32
}

func (h *HybridRAG) Search(query string, topK int) ([]chanollama.Document, error) {
    // Get vector results
    vectorResults, err := h.vectorStore.SearchSimilar(query, topK*2)
    if err != nil {
        return nil, err
    }
    
    // Get keyword results  
    keywordResults := h.simpleStore.SearchRelevant(query, topK*2)
    
    // Combine and score
    scored := make(map[string]float32)
    
    // Add vector scores
    for _, result := range vectorResults {
        scored[result.Document.ID] += result.Similarity * h.vectorWeight
    }
    
    // Add keyword scores (simplified)
    for _, doc := range keywordResults {
        scored[doc.ID] += 0.8 * h.keywordWeight // Assume 0.8 keyword match score
    }
    
    // Sort by combined score and return top results
    // Implementation details...
    
    return finalResults, nil
}
```

For more examples, see the [examples directory](../examples/) and [API reference](./api/).