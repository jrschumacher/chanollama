package chanollama

import (
	"fmt"
	"strings"
	"time"
)

// RealisticMockEngine provides deterministic but varied responses for testing
type RealisticMockEngine struct {
	ResponseDatabase map[string][]string
	currentIndex     map[string]int
	latency          time.Duration
}

// NewRealisticMockEngine creates a mock engine with pre-defined realistic responses
func NewRealisticMockEngine() *RealisticMockEngine {
	return &RealisticMockEngine{
		ResponseDatabase: map[string][]string{
			"ecommerce_platform": {
				`E-commerce platforms provide comprehensive inventory management through various tracking mechanisms and automated workflows.

**Core Components:**
- **Inventory Tracking**: Products are tracked from warehouse to delivery
- **Order Management**: Automated order processing and fulfillment systems
- **Customer Management**: Centralized customer data and relationship management
- **Analytics Dashboard**: Real-time monitoring of sales metrics and trends

**How it Works:**
1. Products are cataloged using standardized SKU systems
2. Inventory levels are automatically updated with each transaction
3. Customer orders are processed through automated workflows
4. Analytics provide insights for business optimization

**Benefits:**
- Streamlined operations from order to delivery
- Real-time inventory visibility across all channels
- Comprehensive customer journey tracking
- Data-driven business decision making

Modern e-commerce platforms enable businesses to scale operations while maintaining customer satisfaction.`,

				`Cloud-native e-commerce platforms integrate inventory management directly with sales channels. Unlike traditional isolated systems, these solutions synchronize data across all touchpoints, ensuring consistency regardless of sales channel.

**Key Features:**
- **Multi-Channel Integration**: Seamless connection between online, mobile, and in-store systems
- **Multiple Formats**: Support for various product catalogs and data standards
- **API Integration**: Support for Java, JavaScript, Go, Python, and other languages
- **Enterprise Ready**: Scalable architecture with comprehensive reporting and analytics

**Use Cases:**
- Unified inventory across multiple sales channels
- Real-time product availability with automated reordering
- Customer behavior analytics across all touchpoints
- Automated marketing campaigns based on purchase history

These platforms ensure that product data remains synchronized throughout the entire sales lifecycle, from listing to delivery.`,
			},

			"content_management": {
				`To create content workflows, use a structured approach with clear categorization:

**Step 1: Define the Workflow**
` + "```bash\ncontent-cli workflow create \\\n  --name \"blog.publishing\" \\\n  --type SEQUENTIAL\n```" + `

**Step 2: Verify Creation**
` + "```bash\ncontent-cli workflow list\n```" + `

**Best Practices:**
- Use hierarchical naming (e.g., "publication.section.category")
- Choose appropriate workflow types (SEQUENTIAL, PARALLEL, CONDITIONAL)
- Consider future content types when designing workflow structure

**Example JSON Configuration:**
` + "```json\n{\n  \"name\": \"blog.content-type\",\n  \"type\": \"SEQUENTIAL\",\n  \"stages\": [\n    {\"stage\": \"draft\"},\n    {\"stage\": \"review\"},\n    {\"stage\": \"approved\"},\n    {\"stage\": \"published\"}\n  ]\n}\n```" + `

**Performance Considerations:**
- Workflow design affects content processing performance
- Use meaningful names that align with editorial structure
- Regular review of workflow usage and stage transitions`,

				`Creating content workflows is fundamental to publishing system design. Here's a comprehensive approach:

**Command Syntax:**
` + "```bash\ncontent-cli workflow create \\\n  --name \"your.workflow.name\" \\\n  --type SEQUENTIAL \\\n  --active\n```" + `

**Workflow Types:**
- **SEQUENTIAL**: Content moves through stages in order
- **PARALLEL**: Multiple reviewers can work simultaneously
- **CONDITIONAL**: Dynamic routing based on content attributes

**Planning Your Workflows:**
1. **Editorial Alignment**: Mirror your editorial structure
2. **Content Classification**: Align with content types and sensitivity
3. **User Roles**: Consider how editors and reviewers will interact
4. **Future Growth**: Plan for new content types and processes

**Adding Workflow Stages:**
` + "```bash\ncontent-cli workflow stages create \\\n  --workflow-name \"blog.publishing\" \\\n  --stage \"peer-review\" \\\n  --active\n```" + `

Remember: Good workflow design is crucial for maintainable and efficient publishing systems.`,
			},

			"golang_integration": {
				`Here's a complete Go example for processing orders with an e-commerce SDK:

**Go Module:**
` + "```bash\ngo mod init your-ecommerce-app\ngo get github.com/example/ecommerce-sdk@v1.2.0\n```" + `

**Order Processing Code:**
` + "```go\npackage main\n\nimport (\n    \"fmt\"\n    \"log\"\n    \"github.com/example/ecommerce-sdk\"\n)\n\nfunc main() {\n    // Configure SDK\n    config := ecommerce.Config{\n        PlatformEndpoint: \"https://your-platform.com\",\n        APIKey:          \"your-api-key\",\n        MerchantID:      \"your-merchant-id\",\n    }\n    \n    client, err := ecommerce.New(config)\n    if err != nil {\n        log.Fatal(err)\n    }\n    \n    // Create order with customer details\n    order := &ecommerce.Order{\n        CustomerID: \"customer-123\",\n        Items: []ecommerce.Item{\n            {SKU: \"PROD-001\", Quantity: 2, Price: 29.99},\n            {SKU: \"PROD-002\", Quantity: 1, Price: 15.50},\n        },\n        ShippingAddress: ecommerce.Address{\n            Street: \"123 Main St\",\n            City:   \"Springfield\",\n            State:  \"IL\",\n            ZIP:    \"62701\",\n        },\n    }\n    \n    // Process the order\n    result, err := client.ProcessOrder(order)\n    if err != nil {\n        log.Fatal(err)\n    }\n    \n    fmt.Printf(\"Order processed successfully! Order ID: %s\\n\", result.OrderID)\n}\n```" + `

**Key Development Practices:**
- Store API credentials securely (environment variables or vault)
- Use appropriate order attributes for your business logic
- Implement proper error handling and logging
- Validate order processing before deployment

**Testing the Order Processing:**
` + "```go\n// Verify the order was processed\nif result.Status == \"confirmed\" {\n    fmt.Printf(\"Order confirmed with total: $%.2f\\n\", result.Total)\n    fmt.Printf(\"Estimated delivery: %s\\n\", result.EstimatedDelivery)\n}\n```",

				`Complete Go order processing implementation with error handling:

**Configuration Package:**
` + "```go\npackage config\n\nimport (\n    \"fmt\"\n    \"os\"\n    \"github.com/example/ecommerce-sdk\"\n)\n\ntype Config struct {\n    PlatformURL string\n    APIKey      string\n    MerchantID  string\n    Environment string\n}\n\nfunc LoadFromEnv() (*Config, error) {\n    platformURL := os.Getenv(\"ECOMMERCE_PLATFORM_URL\")\n    apiKey := os.Getenv(\"ECOMMERCE_API_KEY\")\n    merchantID := os.Getenv(\"ECOMMERCE_MERCHANT_ID\")\n    \n    if platformURL == \"\" || apiKey == \"\" || merchantID == \"\" {\n        return nil, fmt.Errorf(\"required environment variables not set\")\n    }\n    \n    return &Config{\n        PlatformURL: platformURL,\n        APIKey:      apiKey,\n        MerchantID:  merchantID,\n        Environment: getEnvOrDefault(\"ECOMMERCE_ENV\", \"production\"),\n    }, nil\n}\n\nfunc getEnvOrDefault(key, defaultValue string) string {\n    if value := os.Getenv(key); value != \"\" {\n        return value\n    }\n    return defaultValue\n}\n```" + `

**Order Processing Service:**
` + "```go\npackage service\n\nimport (\n    \"fmt\"\n    \"log\"\n    \"your-app/config\"\n    \"github.com/example/ecommerce-sdk\"\n)\n\ntype OrderService struct {\n    client *ecommerce.Client\n}\n\nfunc NewOrderService() (*OrderService, error) {\n    cfg, err := config.LoadFromEnv()\n    if err != nil {\n        return nil, err\n    }\n    \n    client, err := ecommerce.New(ecommerce.Config{\n        PlatformEndpoint: cfg.PlatformURL,\n        APIKey:          cfg.APIKey,\n        MerchantID:      cfg.MerchantID,\n    })\n    if err != nil {\n        return nil, err\n    }\n    \n    return &OrderService{client: client}, nil\n}\n\nfunc (s *OrderService) ProcessOrder(customerID string, items []ecommerce.Item, address ecommerce.Address) (*ecommerce.OrderResult, error) {\n    order := &ecommerce.Order{\n        CustomerID:      customerID,\n        Items:          items,\n        ShippingAddress: address,\n    }\n    \n    result, err := s.client.ProcessOrder(order)\n    if err != nil {\n        log.Printf(\"Order processing failed: %v\", err)\n        return nil, fmt.Errorf(\"order processing failed: %w\", err)\n    }\n    \n    log.Printf(\"Successfully processed order: %s -> %s\", customerID, result.OrderID)\n    return result, nil\n}\n```",
			},

			"container_deploy": {
				`Here's a comprehensive guide to deploying web applications with Docker:

**Docker Compose Setup:**
` + "```yaml\nversion: '3.8'\nservices:\n  webapp:\n    image: example/webapp:latest\n    ports:\n      - \"8080:8080\"\n    environment:\n      - DATABASE_URL=postgresql://postgres:password@db:5432/webapp\n      - REDIS_URL=redis://cache:6379\n      - API_ENDPOINT=https://api.example.com\n    depends_on:\n      - db\n      - cache\n      \n  api-service:\n    image: example/api-service:latest\n    ports:\n      - \"8000:8000\"\n    environment:\n      - DATABASE_URL=postgresql://postgres:password@db:5432/webapp\n      - WEBAPP_ENDPOINT=http://webapp:8080\n    depends_on:\n      - db\n      \n  db:\n    image: postgres:15\n    environment:\n      - POSTGRES_DB=webapp\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=password\n    volumes:\n      - postgres_data:/var/lib/postgresql/data\n      \n  cache:\n    image: redis:7-alpine\n    ports:\n      - \"6379:6379\"\n      \nvolumes:\n  postgres_data:\n```" + `

**Quick Start Commands:**
` + "```bash\n# Clone and start\ngit clone https://github.com/example/webapp.git\ncd webapp\n\n# Start all services\ndocker-compose up -d\n\n# Check service health\ndocker-compose ps\n\n# View logs\ndocker-compose logs -f webapp\n```" + `

**Configuration:**
1. **Environment Variables**: Update API_ENDPOINT for your backend services
2. **Networking**: Ensure services can communicate
3. **Persistence**: Database data is persisted in volumes
4. **Performance**: Configure Redis cache for optimal performance

**Verification:**
` + "```bash\n# Test webapp health\ncurl http://localhost:8080/health\n\n# Test API service health\ncurl http://localhost:8000/health\n\n# Check database connection\ndocker-compose exec db psql -U postgres -d webapp -c \"\\dt\"\n\n# Test Redis cache\ndocker-compose exec cache redis-cli ping\n```" + `

**Production Considerations:**
- Use environment-specific configuration files
- Implement proper secrets management
- Set up monitoring and logging
- Configure backup strategies
- Use reverse proxy (nginx) for SSL termination and load balancing`,
			},

			"debugging": {
				`When debugging API integration errors, follow this systematic approach:

**1. Verify Authentication:**
` + "```bash\n# Check if API credentials are valid\napi-cli auth test --endpoint https://api.your-platform.com\n\n# Verify API key status\napi-cli auth status --api-key your-api-key\n```" + `

**2. Check Request Configuration:**
` + "```bash\n# Test basic connectivity\ncurl -H \"Authorization: Bearer your-token\" https://api.your-platform.com/health\n\n# Verify endpoint configuration\napi-cli config get --endpoint\n```" + `

**3. Validate Request Format:**
` + "```bash\n# Check request schema\napi-cli schema validate --request-type orders\n\n# Test with minimal payload\napi-cli orders create --test-mode --customer-id test-123\n```" + `

**4. Debug API Request Process:**
- **Check request syntax**: Ensure JSON structure matches API schema
- **Verify permissions**: API key must have required scopes
- **Test with minimal payload**: Start with simple request structure
- **Check network connectivity**: Ensure API endpoint is reachable

**5. Common Issues & Solutions:**

**Authentication Failed:**
- API key expired or invalid
- Missing required authentication headers
- Incorrect endpoint URL in request

**Network Errors:**
- API endpoint unreachable
- Certificate validation failures
- Rate limiting or throttling

**Configuration Issues:**
- Missing environment variables
- Incorrect API endpoint URL
- Invalid request format or content-type

**Debugging Commands:**
` + "```bash\n# Enable verbose logging\nexport API_CLI_LOG_LEVEL=debug\n\n# Test simple API call\napi-cli orders list --limit 1 --verbose\n\n# Check API service health\ncurl https://api.your-platform.com/health\n\n# Test connectivity\nping api.your-platform.com\n```" + `

**Log Analysis:**
Look for specific error patterns in logs:
- Authentication failures (401, 403)
- Request format errors (400, 422)
- Network timeouts (504, connection refused)
- Rate limiting errors (429)`,
			},

			"default": {
				`This is a comprehensive response demonstrating the capabilities of modern application platforms and integration systems.

**Application Platform Overview:**
Modern application platforms provide comprehensive workflow management through APIs, integrations, and automated processing mechanisms.

**Key Features:**
- RESTful APIs that connect with various data sources and services
- Flexible configuration based on business requirements and workflows
- Multiple format support for different integration patterns
- Cross-platform SDK support (Java, Go, JavaScript, Python, .NET)
- Enterprise-ready with monitoring, logging, and analytics support

**Common Use Cases:**
1. **E-commerce Integration**: Connect inventory, payment, and shipping systems
2. **Content Management**: Automated publishing workflows and content distribution
3. **Data Processing**: ETL pipelines and real-time data transformation
4. **Healthcare Systems**: Patient data management and clinical workflows
5. **Financial Services**: Transaction processing and regulatory reporting

**Getting Started:**
1. Install platform CLI tool or SDK
2. Configure authentication with your platform endpoint
3. Create workflow definitions and business rules
4. Process data with appropriate transformations
5. Monitor and optimize system performance

**Best Practices:**
- Design workflows that align with business processes
- Implement proper error handling and retry mechanisms
- Regular monitoring of system performance and usage patterns
- Use environment-specific configurations for different deployment stages
- Implement comprehensive logging and observability

For specific implementation guidance, please ask about particular aspects like workflow design, API integration, or deployment strategies.`,
			},
		},
		currentIndex: make(map[string]int),
		latency:      100 * time.Millisecond, // Simulate realistic response time
	}
}

func (m *RealisticMockEngine) Chat(messages []Message) (*Response, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	// Simulate processing time
	time.Sleep(m.latency)

	query := messages[len(messages)-1].Content
	response := m.generateContextualResponse(query)

	return &Response{
		Content:      response,
		Metadata:     make(map[string]string),
		Error:        nil,
		Tokens:       len(response) / 4, // Rough token estimate
		Model:        "offline-mock-engine",
		FinishReason: "stop",
	}, nil
}

func (m *RealisticMockEngine) ChatStream(messages []Message, callback StreamCallback) *Response {
	response, err := m.Chat(messages)
	if err != nil {
		return &Response{Error: err}
	}

	// Simulate streaming by sending chunks
	words := strings.Fields(response.Content)
	for i, word := range words {
		if i > 0 {
			callback(" ")
		}
		callback(word)
		time.Sleep(10 * time.Millisecond) // Simulate streaming delay
	}

	return response
}

func (m *RealisticMockEngine) IsHealthy() bool {
	return true
}

func (m *RealisticMockEngine) Close() error {
	return nil
}

// generateContextualResponse selects appropriate response based on query content
func (m *RealisticMockEngine) generateContextualResponse(query string) string {
	queryLower := strings.ToLower(query)

	// Determine response category based on query content
	var category string
	switch {
	case strings.Contains(queryLower, "what is") && (strings.Contains(queryLower, "ecommerce") || strings.Contains(queryLower, "platform") || strings.Contains(queryLower, "inventory")):
		category = "ecommerce_platform"
	case strings.Contains(queryLower, "workflow") && strings.Contains(queryLower, "create"):
		category = "content_management"
	case strings.Contains(queryLower, "go") && (strings.Contains(queryLower, "order") || strings.Contains(queryLower, "code") || strings.Contains(queryLower, "integration")):
		category = "golang_integration"
	case strings.Contains(queryLower, "docker") && strings.Contains(queryLower, "deploy"):
		category = "container_deploy"
	case (strings.Contains(queryLower, "error") || strings.Contains(queryLower, "debug")) && (strings.Contains(queryLower, "api") || strings.Contains(queryLower, "integration")):
		category = "debugging"
	default:
		category = "default"
	}

	// Get responses for category
	responses, exists := m.ResponseDatabase[category]
	if !exists {
		responses = m.ResponseDatabase["default"]
	}

	// Rotate through responses to provide variety
	if m.currentIndex[category] >= len(responses) {
		m.currentIndex[category] = 0
	}

	response := responses[m.currentIndex[category]]
	m.currentIndex[category]++

	return response
}

// SetLatency allows adjusting response time for testing
func (m *RealisticMockEngine) SetLatency(duration time.Duration) {
	m.latency = duration
}

// AddResponse allows adding custom responses for testing
func (m *RealisticMockEngine) AddResponse(category string, response string) {
	if m.ResponseDatabase[category] == nil {
		m.ResponseDatabase[category] = []string{}
	}
	m.ResponseDatabase[category] = append(m.ResponseDatabase[category], response)
}
