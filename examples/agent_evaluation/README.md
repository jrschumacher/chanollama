# Agent Evaluation Example

This example demonstrates how to use chanollama's flexible agent evaluation system to create domain-specific intelligent routing for your application.

## Key Features

- **No Hardcoded Agents**: Define agents specific to your application domain
- **Metadata-Driven Routing**: Uses agent capabilities, tags, and descriptions for intelligent routing
- **Multiple Domains**: Shows how the same system works for different applications

## Running the Example

```bash
cd examples/agent_evaluation
go run main.go
```

## How It Works

### 1. Create Agent Evaluator
```go
evaluator := chanollama.NewAgentEvaluator(nil) // No hardcoded agents
```

### 2. Define Domain-Specific Agents
```go
err := evaluator.LoadAgent("product", systemPrompt, chanollama.AgentMetadata{
    Name:         "Product Manager",
    Description:  "Handles product catalog, pricing, inventory, and merchandising",
    Capabilities: []string{"catalog_management", "pricing", "inventory", "merchandising"},
    Tags:         []string{"product", "catalog", "inventory", "pricing", "ecommerce"},
})
```

### 3. Intelligent Query Routing
```go
agent, err := evaluator.EvaluateQuery(ctx, "How do I update product pricing?")
// Returns: "product" - automatically routed based on keywords and metadata
```

## Example Domains

### E-commerce Platform
- **Product Agent**: Catalog, pricing, inventory management
- **Customer Agent**: Support, orders, returns, service
- **Operations Agent**: Fulfillment, logistics, warehousing
- **Analytics Agent**: Sales metrics, reporting, insights

### Content Management System
- **Content Agent**: Writing, editing, publishing, strategy
- **Media Agent**: Images, videos, digital assets, optimization
- **Workflow Agent**: Approvals, collaboration, governance

## Routing Logic

The system uses **LLM-first intelligent routing**:

1. **LLM-Based Evaluation (Primary)**: Uses the LLM engine to analyze the query and agent descriptions for intelligent routing
2. **Simple Fallback**: If no LLM engine is provided or LLM evaluation fails, returns first available agent

This approach is much more reliable and intelligent than keyword-based matching, as the LLM can understand context, synonyms, and intent naturally.

## Benefits

- **Flexible**: Works with any application domain
- **Intelligent**: Uses both LLM-based and keyword-based evaluation
- **Maintainable**: Single system handles routing for your entire application
- **No Domain Lock-in**: Not limited to security, e-commerce, or any specific domain