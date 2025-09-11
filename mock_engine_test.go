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
			"data_security_concept": {
				`Data security platforms provide comprehensive protection through various access control mechanisms and encryption techniques.

**Core Components:**
- **Encryption**: Data is encrypted at rest and in transit
- **Access Control**: Role-based or attribute-based access control systems
- **Key Management**: Centralized key management and policy enforcement
- **Audit Trail**: Logging and monitoring of all data access events

**How it Works:**
1. Data is encrypted using strong cryptographic algorithms
2. Access policies are defined and enforced
3. User credentials and attributes are validated
4. Keys are provided only when policies allow access

**Benefits:**
- Data remains protected even when shared externally
- Fine-grained access control based on user attributes
- Comprehensive audit trails for compliance
- Zero-trust security architecture

Modern data security platforms enable organizations to protect sensitive information while maintaining operational flexibility.`,

				`Zero-trust data security platforms embed access policies directly with encrypted data. Unlike traditional perimeter-based security, these solutions travel with the data, ensuring protection regardless of storage location.

**Key Features:**
- **Attribute-Based Access Control (ABAC)**: Fine-grained permissions based on user and data attributes
- **Multiple Formats**: Support for various encryption formats and standards
- **SDK Integration**: Support for Java, JavaScript, Go, Python, and other languages
- **Enterprise Ready**: Scalable architecture with comprehensive audit trails

**Use Cases:**
- Secure document sharing across organizations
- Data lake protection with granular access controls
- Compliance-driven data governance (GDPR, HIPAA, SOX)
- Enterprise data protection and classification

These platforms ensure that sensitive data remains protected throughout its entire lifecycle, from creation to deletion.`,
			},

			"policy_management": {
				`To create access control policies, use a structured approach with clear naming conventions:

**Step 1: Define the Policy**
` + "```bash\nsecurity-cli policy create \\\n  --name \"organization.department\" \\\n  --rule ANY_OF\n```" + `

**Step 2: Verify Creation**
` + "```bash\nsecurity-cli policy list\n```" + `

**Best Practices:**
- Use hierarchical naming (e.g., "company.division.team")
- Choose appropriate rules (ANY_OF, ALL_OF, HIERARCHY)
- Consider future expansion when designing policy structure

**Example JSON Definition:**
` + "```json\n{\n  \"name\": \"organization.security-level\",\n  \"rule\": \"ANY_OF\",\n  \"values\": [\n    {\"value\": \"public\"},\n    {\"value\": \"internal\"},\n    {\"value\": \"confidential\"},\n    {\"value\": \"restricted\"}\n  ]\n}\n```" + `

**Security Considerations:**
- Policy design affects access control performance
- Use meaningful names that align with organizational structure
- Regular audit of policy usage and attribute assignments`,

				`Creating access control policies is fundamental to security system design. Here's a comprehensive approach:

**Command Syntax:**
` + "```bash\nsecurity-cli policy create \\\n  --name \"your.policy.name\" \\\n  --rule ANY_OF \\\n  --active\n```" + `

**Policy Rules:**
- **ANY_OF**: User needs any attribute from the policy
- **ALL_OF**: User needs all attributes from the policy  
- **HIERARCHY**: Supports hierarchical attribute relationships

**Planning Your Policies:**
1. **Organizational Alignment**: Mirror your organizational structure
2. **Data Classification**: Align with data sensitivity levels
3. **User Roles**: Consider how users will be assigned attributes
4. **Future Growth**: Plan for expansion and changes

**Adding Policy Values:**
` + "```bash\nsecurity-cli policy values create \\\n  --policy-name \"organization.department\" \\\n  --value \"engineering\" \\\n  --active\n```" + `

Remember: Good policy design is crucial for maintainable and performant access control systems.`,
			},

			"java_integration": {
				`Here's a complete Java example for encrypting files with a security SDK:

**Maven Dependency:**
` + "```xml\n<dependency>\n    <groupId>com.example</groupId>\n    <artifactId>security-sdk</artifactId>\n    <version>1.0.0</version>\n</dependency>\n```" + `

**Encryption Code:**
` + "```java\nimport com.example.security.sdk.*;\nimport java.io.*;\nimport java.nio.file.*;\n\npublic class FileEncryption {\n    public static void main(String[] args) throws Exception {\n        // Configure SDK\n        SDKConfig config = SDKConfig.builder()\n            .platformEndpoint(\"https://your-platform.com\")\n            .clientId(\"your-client-id\")\n            .clientSecret(\"your-client-secret\")\n            .build();\n            \n        SDK sdk = new SDK(config);\n        \n        // Create encryption config with attributes\n        EncryptionConfig encConfig = EncryptionConfig.builder()\n            .addDataAttribute(\"organization.department\", \"engineering\")\n            .addDataAttribute(\"security.classification\", \"internal\")\n            .build();\n            \n        // Encrypt file\n        Path inputFile = Paths.get(\"sensitive-data.txt\");\n        Path outputFile = Paths.get(\"sensitive-data.enc\");\n        \n        try (InputStream input = Files.newInputStream(inputFile);\n             OutputStream output = Files.newOutputStream(outputFile)) {\n            \n            sdk.encryptFile(input, output, encConfig);\n            System.out.println(\"File encrypted successfully!\");\n        }\n    }\n}\n```" + `

**Key Security Practices:**
- Store credentials securely (environment variables or key vault)
- Use appropriate data attributes for your use case
- Implement proper error handling and logging
- Validate encryption before deployment

**Testing the Encryption:**
` + "```java\n// Verify the encrypted file was created\nFile encFile = new File(\"sensitive-data.enc\");\nif (encFile.exists()) {\n    System.out.println(\"Encrypted file size: \" + encFile.length() + \" bytes\");\n}\n```",

				`Complete Java encryption implementation with error handling:

**Configuration Class:**
` + "```java\nimport com.example.security.sdk.*;\n\npublic class SecurityConfig {\n    private static final String PLATFORM_URL = System.getenv(\"SECURITY_PLATFORM_URL\");\n    private static final String CLIENT_ID = System.getenv(\"SECURITY_CLIENT_ID\");\n    private static final String CLIENT_SECRET = System.getenv(\"SECURITY_CLIENT_SECRET\");\n    \n    public static SDK createSDK() throws Exception {\n        if (PLATFORM_URL == null || CLIENT_ID == null || CLIENT_SECRET == null) {\n            throw new IllegalStateException(\"Required environment variables not set\");\n        }\n        \n        SDKConfig config = SDKConfig.builder()\n            .platformEndpoint(PLATFORM_URL)\n            .clientId(CLIENT_ID)\n            .clientSecret(CLIENT_SECRET)\n            .build();\n            \n        return new SDK(config);\n    }\n}\n```" + `

**Encryption Service:**
` + "```java\nimport java.util.List;\nimport java.util.Map;\n\npublic class EncryptionService {\n    private final SDK sdk;\n    \n    public EncryptionService() throws Exception {\n        this.sdk = SecurityConfig.createSDK();\n    }\n    \n    public void encryptFile(String inputPath, String outputPath, Map<String, String> attributes) {\n        try {\n            EncryptionConfig.Builder configBuilder = EncryptionConfig.builder();\n            \n            // Add data attributes\n            attributes.forEach(configBuilder::addDataAttribute);\n            \n            EncryptionConfig encConfig = configBuilder.build();\n            \n            // Perform encryption\n            try (InputStream input = Files.newInputStream(Paths.get(inputPath));\n                 OutputStream output = Files.newOutputStream(Paths.get(outputPath))) {\n                \n                sdk.encryptFile(input, output, encConfig);\n                System.out.println(\"Successfully encrypted: \" + inputPath + \" -> \" + outputPath);\n            }\n        } catch (Exception e) {\n            System.err.println(\"Encryption failed: \" + e.getMessage());\n            throw new RuntimeException(\"Encryption failed\", e);\n        }\n    }\n}\n```",
			},

			"docker_deploy": {
				`Here's a comprehensive guide to deploying security platforms with Docker:

**Docker Compose Setup:**
` + "```yaml\nversion: '3.8'\nservices:\n  security-platform:\n    image: example/security-platform:latest\n    ports:\n      - \"8080:8080\"\n    environment:\n      - DATABASE_URL=postgresql://postgres:password@db:5432/security\n      - OIDC_ISSUER=https://your-identity-provider.com\n      - KEY_SERVICE_ENDPOINT=https://your-key-service.com\n    depends_on:\n      - db\n      - key-service\n      \n  key-service:\n    image: example/key-service:latest\n    ports:\n      - \"8000:8000\"\n    environment:\n      - DATABASE_URL=postgresql://postgres:password@db:5432/keys\n      - PLATFORM_ENDPOINT=http://security-platform:8080\n    depends_on:\n      - db\n      \n  db:\n    image: postgres:15\n    environment:\n      - POSTGRES_DB=security\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=password\n    volumes:\n      - postgres_data:/var/lib/postgresql/data\n      \nvolumes:\n  postgres_data:\n```" + `

**Quick Start Commands:**
` + "```bash\n# Clone and start\ngit clone https://github.com/example/security-platform.git\ncd security-platform\n\n# Start all services\ndocker-compose up -d\n\n# Check service health\ndocker-compose ps\n\n# View logs\ndocker-compose logs -f security-platform\n```" + `

**Configuration:**
1. **Environment Variables**: Update OIDC_ISSUER for your identity provider
2. **Networking**: Ensure services can communicate
3. **Persistence**: Database data is persisted in volumes
4. **Security**: Change default passwords in production

**Verification:**
` + "```bash\n# Test platform health\ncurl http://localhost:8080/health\n\n# Test key service health\ncurl http://localhost:8000/health\n\n# Check database connection\ndocker-compose exec db psql -U postgres -d security -c \"\\dt\"\n```" + `

**Production Considerations:**
- Use environment-specific configuration files
- Implement proper secrets management
- Set up monitoring and logging
- Configure backup strategies
- Use reverse proxy (nginx) for SSL termination`,
			},

			"troubleshooting": {
				`When troubleshooting encryption permission errors, follow this systematic approach:

**1. Verify Authentication:**
` + "```bash\n# Check if credentials are valid\nsecurity-cli auth login --platform-endpoint https://your-platform.com\n\n# Verify token\nsecurity-cli auth print-access-token\n```" + `

**2. Check Attribute Assignments:**
` + "```bash\n# List user's attributes\nsecurity-cli policy user-attributes list --user-id your-user-id\n\n# Verify required attributes exist\nsecurity-cli policy attributes list --policy-name organization.department\n```" + `

**3. Validate Policy Configuration:**
` + "```bash\n# Check policy rules\nsecurity-cli policy get --name organization.department\n\n# Verify attribute values are active\nsecurity-cli policy values get --policy-name organization.department --value engineering\n```" + `

**4. Debug Encryption Process:**
- **Check attribute syntax**: Ensure attribute names match exactly
- **Verify permissions**: User must have attributes required by policy
- **Test with minimal attributes**: Start with simple attribute set
- **Check network connectivity**: Ensure key service is reachable

**5. Common Issues & Solutions:**

**Permission Denied:**
- User lacks required attributes
- Policy doesn't exist
- Incorrect attribute syntax in encryption request

**Network Errors:**
- Key service endpoint unreachable
- Certificate validation failures
- Firewall blocking connections

**Configuration Issues:**
- Missing environment variables
- Incorrect platform endpoint
- Invalid client credentials

**Debugging Commands:**
` + "```bash\n# Enable verbose logging\nexport SECURITY_CLI_LOG_LEVEL=debug\n\n# Test simple encryption\nsecurity-cli encrypt --file test.txt --output test.enc --attributes \"organization.department=engineering\"\n\n# Check key service health\ncurl https://your-key-service.com/health\n```" + `

**Log Analysis:**
Look for specific error patterns in logs:
- Authentication failures
- Attribute resolution errors  
- Network timeouts
- Policy evaluation denials`,
			},

			"default": {
				`This is a comprehensive response demonstrating the capabilities of modern security platforms and evaluation systems.

**Security Platform Overview:**
Modern security platforms provide comprehensive data protection through encryption, access controls, and policy enforcement mechanisms.

**Key Features:**
- Embedded access policies that travel with encrypted data
- Fine-grained permissions based on user and data attributes
- Multiple format support for different use cases
- Cross-platform SDK support (Java, Go, JavaScript, Python)
- Enterprise-ready with audit trails and compliance support

**Common Use Cases:**
1. **Secure Document Sharing**: Share sensitive documents across organizational boundaries
2. **Data Lake Protection**: Implement granular access controls in big data environments
3. **Compliance**: Meet regulatory requirements (GDPR, HIPAA, SOX)
4. **Zero-Trust Architecture**: Implement data-centric security models

**Getting Started:**
1. Install security CLI tool
2. Configure authentication with your security platform
3. Create access control policies
4. Encrypt data with appropriate attributes
5. Share encrypted data securely

**Best Practices:**
- Design attribute hierarchies that align with organizational structure
- Implement least-privilege access principles
- Regular audit of attribute assignments and policy effectiveness
- Use environment-specific configurations for different deployment stages
- Monitor access patterns and policy evaluation performance

For specific implementation guidance, please ask about particular aspects like policy design, SDK integration, or deployment strategies.`,
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
	case strings.Contains(queryLower, "what is") && (strings.Contains(queryLower, "security") || strings.Contains(queryLower, "data protection")):
		category = "data_security_concept"
	case strings.Contains(queryLower, "policy") && strings.Contains(queryLower, "create"):
		category = "policy_management"
	case strings.Contains(queryLower, "java") && (strings.Contains(queryLower, "encrypt") || strings.Contains(queryLower, "code")):
		category = "java_integration"
	case strings.Contains(queryLower, "docker") && strings.Contains(queryLower, "deploy"):
		category = "docker_deploy"
	case strings.Contains(queryLower, "permission") && strings.Contains(queryLower, "error"):
		category = "troubleshooting"
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
