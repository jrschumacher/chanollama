package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jrschumacher/chanollama"
)

func main() {
	// Create an agent evaluator (no hardcoded agents)
	evaluator := chanollama.NewAgentEvaluator(nil) // nil engine for keyword-only evaluation

	// Example: E-commerce application agents
	loadEcommerceAgents(evaluator)

	// Test queries
	testQueries := []string{
		"How do I update product pricing?",
		"Customer wants to return an order",
		"Set up inventory tracking system",
		"What are the current sales metrics?",
	}

	fmt.Println("🛒 E-commerce Agent Routing Examples:")
	fmt.Println("=====================================")

	for _, query := range testQueries {
		agent, err := evaluator.EvaluateQuery(context.Background(), query)
		if err != nil {
			log.Printf("Error evaluating query '%s': %v", query, err)
			continue
		}

		fmt.Printf("Query: \"%s\"\n", query)
		fmt.Printf("→ Routed to: %s agent\n\n", agent)
	}

	// Example: Different domain - Content Management System
	fmt.Println("📝 Content Management System Example:")
	fmt.Println("====================================")

	cmsEvaluator := chanollama.NewAgentEvaluator(nil)
	loadCMSAgents(cmsEvaluator)

	cmsQueries := []string{
		"How do I create a new blog post?",
		"Optimize images for web performance",
		"Set up content approval workflow",
	}

	for _, query := range cmsQueries {
		agent, err := cmsEvaluator.EvaluateQuery(context.Background(), query)
		if err != nil {
			log.Printf("Error evaluating query '%s': %v", query, err)
			continue
		}

		fmt.Printf("Query: \"%s\"\n", query)
		fmt.Printf("→ Routed to: %s agent\n\n", agent)
	}
}

// loadEcommerceAgents demonstrates loading domain-specific agents for e-commerce
func loadEcommerceAgents(evaluator *chanollama.AgentEvaluator) {
	// Product Management Agent
	err := evaluator.LoadAgent("product", 
		`You are a product management specialist for e-commerce platforms.
		Focus on catalog management, pricing, inventory, and product optimization.
		Provide practical guidance on product data, merchandising, and catalog operations.`, 
		chanollama.AgentMetadata{
			Name:         "Product Manager",
			Description:  "Handles product catalog, pricing, inventory, and merchandising",
			Capabilities: []string{"catalog_management", "pricing", "inventory", "merchandising", "product_data"},
			Tags:         []string{"product", "catalog", "inventory", "pricing", "ecommerce"},
		})
	if err != nil {
		log.Fatalf("Failed to load product agent: %v", err)
	}

	// Customer Service Agent
	err = evaluator.LoadAgent("customer",
		`You are a customer service specialist for e-commerce platforms.
		Focus on order management, returns, shipping, and customer satisfaction.
		Provide guidance on customer support processes and order fulfillment.`,
		chanollama.AgentMetadata{
			Name:         "Customer Service Rep",
			Description:  "Handles customer support, orders, returns, and service issues",
			Capabilities: []string{"customer_support", "order_management", "returns", "shipping", "satisfaction"},
			Tags:         []string{"customer", "support", "orders", "returns", "service"},
		})
	if err != nil {
		log.Fatalf("Failed to load customer agent: %v", err)
	}

	// Operations Agent  
	err = evaluator.LoadAgent("operations",
		`You are an operations specialist for e-commerce platforms.
		Focus on fulfillment, logistics, warehousing, and operational efficiency.
		Provide guidance on supply chain, automation, and operational processes.`,
		chanollama.AgentMetadata{
			Name:         "Operations Manager",
			Description:  "Handles fulfillment, logistics, warehousing, and operations",
			Capabilities: []string{"fulfillment", "logistics", "warehousing", "automation", "supply_chain"},
			Tags:         []string{"operations", "fulfillment", "logistics", "warehouse", "automation"},
		})
	if err != nil {
		log.Fatalf("Failed to load operations agent: %v", err)
	}

	// Analytics Agent
	err = evaluator.LoadAgent("analytics",
		`You are an analytics specialist for e-commerce platforms.
		Focus on sales metrics, performance analysis, and business intelligence.
		Provide insights on revenue, conversion, customer behavior, and growth.`,
		chanollama.AgentMetadata{
			Name:         "Analytics Specialist",
			Description:  "Handles sales metrics, analytics, reporting, and business intelligence",
			Capabilities: []string{"analytics", "reporting", "metrics", "insights", "performance"},
			Tags:         []string{"analytics", "metrics", "sales", "reporting", "insights"},
		})
	if err != nil {
		log.Fatalf("Failed to load analytics agent: %v", err)
	}
}

// loadCMSAgents demonstrates loading agents for a different domain - Content Management
func loadCMSAgents(evaluator *chanollama.AgentEvaluator) {
	// Content Creator Agent
	err := evaluator.LoadAgent("content",
		`You are a content creation specialist for content management systems.
		Focus on writing, editing, publishing, and content strategy.
		Provide guidance on content creation workflows and editorial processes.`,
		chanollama.AgentMetadata{
			Name:         "Content Creator",
			Description:  "Handles content creation, editing, publishing, and content strategy",
			Capabilities: []string{"content_creation", "editing", "publishing", "content_strategy", "writing"},
			Tags:         []string{"content", "writing", "editing", "publishing", "cms"},
		})
	if err != nil {
		log.Fatalf("Failed to load content agent: %v", err)
	}

	// Media Agent
	err = evaluator.LoadAgent("media",
		`You are a media management specialist for content management systems.
		Focus on image optimization, video processing, and digital asset management.
		Provide guidance on media workflows and performance optimization.`,
		chanollama.AgentMetadata{
			Name:         "Media Manager",
			Description:  "Handles media assets, optimization, and digital asset management",
			Capabilities: []string{"media_management", "image_optimization", "video_processing", "assets", "performance"},
			Tags:         []string{"media", "images", "video", "assets", "optimization"},
		})
	if err != nil {
		log.Fatalf("Failed to load media agent: %v", err)
	}

	// Workflow Agent
	err = evaluator.LoadAgent("workflow",
		`You are a workflow management specialist for content management systems.
		Focus on approval processes, collaboration, and content governance.
		Provide guidance on editorial workflows and content approval systems.`,
		chanollama.AgentMetadata{
			Name:         "Workflow Manager",
			Description:  "Handles content workflows, approvals, and collaboration processes",
			Capabilities: []string{"workflow_management", "approvals", "collaboration", "governance", "processes"},
			Tags:         []string{"workflow", "approval", "collaboration", "process", "governance"},
		})
	if err != nil {
		log.Fatalf("Failed to load workflow agent: %v", err)
	}
}