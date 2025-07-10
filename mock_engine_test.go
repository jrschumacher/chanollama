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
			"opentdf_concept": {
				`OpenTDF (Open Trusted Data Format) is an open-source data protection platform that provides comprehensive data security through attribute-based access control (ABAC). 

**Core Components:**
- **TDF Format**: Encrypted data container with embedded access policies
- **Policy Engine**: ABAC-based decision making for data access
- **Key Access Service (KAS)**: Centralized key management and policy enforcement
- **Attribute Authority**: Manages user and data attributes for access decisions

**How it Works:**
1. Data is encrypted and wrapped in a TDF container
2. Access policies are embedded with the data
3. When accessing data, policies are evaluated against user attributes
4. KAS provides decryption keys only if policies allow access

**Benefits:**
- Data remains protected even when shared externally
- Fine-grained access control based on attributes
- Auditability and compliance support
- Zero-trust data security model

OpenTDF enables organizations to share sensitive data while maintaining control and visibility over who can access it and under what conditions.`,

				`OpenTDF is a zero-trust data security platform that embeds access policies directly with encrypted data. Unlike traditional perimeter-based security, OpenTDF travels with the data, ensuring protection regardless of where the data is stored or shared.

**Key Features:**
- **Attribute-Based Access Control (ABAC)**: Fine-grained permissions based on user, data, and environmental attributes
- **Format Support**: TDF (full-featured), NanoTDF (lightweight), and ZTDF (zero-trust) formats
- **SDK Integration**: Support for Java, JavaScript, Go, and Python
- **Enterprise Ready**: Scalable architecture with audit trails and compliance support

**Use Cases:**
- Secure document sharing across organizations
- Data lake protection with granular access controls
- Compliance-driven data governance (GDPR, HIPAA)
- Government and defense classified data handling

The platform ensures that sensitive data remains protected through its entire lifecycle, from creation to deletion.`,
			},

			"attribute_namespace": {
				`To create an attribute namespace in OpenTDF, use the otdfctl policy command with the following approach:

**Step 1: Create the Namespace**
` + "```bash\notdfctl policy attributes namespaces create \\\n  --name \"organization.department\" \\\n  --rule ANY_OF\n```" + `

**Step 2: Verify Creation**
` + "```bash\notdfctl policy attributes namespaces list\n```" + `

**Best Practices:**
- Use hierarchical naming (e.g., "company.division.team")
- Choose appropriate rules (ANY_OF, ALL_OF, HIERARCHY)
- Consider future attribute expansion when designing namespace structure

**Example JSON Definition:**
` + "```json\n{\n  \"name\": \"organization.security-clearance\",\n  \"rule\": \"ANY_OF\",\n  \"values\": [\n    {\"value\": \"public\"},\n    {\"value\": \"internal\"},\n    {\"value\": \"confidential\"},\n    {\"value\": \"secret\"}\n  ]\n}\n```" + `

**Security Considerations:**
- Namespace design affects policy evaluation performance
- Use meaningful names that align with organizational structure
- Regular audit of namespace usage and attribute assignments`,

				`Creating attribute namespaces is fundamental to OpenTDF policy design. Here's a comprehensive approach:

**Command Syntax:**
` + "```bash\notdfctl policy attributes namespaces create \\\n  --name \"your.namespace.here\" \\\n  --rule ANY_OF \\\n  --active\n```" + `

**Namespace Rules:**
- **ANY_OF**: User needs any attribute from the namespace
- **ALL_OF**: User needs all attributes from the namespace  
- **HIERARCHY**: Supports hierarchical attribute relationships

**Planning Your Namespace:**
1. **Organizational Alignment**: Mirror your org structure
2. **Data Classification**: Align with data sensitivity levels
3. **User Roles**: Consider how users will be assigned attributes
4. **Future Growth**: Plan for expansion and changes

**Adding Attribute Values:**
` + "```bash\notdfctl policy attributes values create \\\n  --namespace-name \"organization.department\" \\\n  --value \"engineering\" \\\n  --active\n```" + `

Remember: Good namespace design is crucial for maintainable and performant ABAC policies.`,
			},

			"java_encrypt": {
				`Here's a complete Java example for encrypting files with the OpenTDF SDK:

**Maven Dependency:**
` + "```xml\n<dependency>\n    <groupId>io.opentdf</groupId>\n    <artifactId>sdk</artifactId>\n    <version>0.4.0</version>\n</dependency>\n```" + `

**Encryption Code:**
` + "```java\nimport io.opentdf.platform.sdk.*;\nimport java.io.*;\nimport java.nio.file.*;\n\npublic class TDFEncryption {\n    public static void main(String[] args) throws Exception {\n        // Configure SDK\n        SDKConfig config = SDKConfig.builder()\n            .platformEndpoint(\"https://your-platform.com\")\n            .clientId(\"your-client-id\")\n            .clientSecret(\"your-client-secret\")\n            .build();\n            \n        SDK sdk = new SDK(config);\n        \n        // Create TDF with attributes\n        TDFConfig tdfConfig = TDFConfig.builder()\n            .addDataAttribute(\"organization.department\", \"engineering\")\n            .addDataAttribute(\"security.classification\", \"internal\")\n            .build();\n            \n        // Encrypt file\n        Path inputFile = Paths.get(\"sensitive-data.txt\");\n        Path outputFile = Paths.get(\"sensitive-data.tdf\");\n        \n        try (InputStream input = Files.newInputStream(inputFile);\n             OutputStream output = Files.newOutputStream(outputFile)) {\n            \n            sdk.createTDF(input, output, tdfConfig);\n            System.out.println(\"File encrypted successfully!\");\n        }\n    }\n}\n```" + `

**Key Security Practices:**
- Store credentials securely (environment variables or key vault)
- Use appropriate data attributes for your use case
- Implement proper error handling and logging
- Validate TDF creation before deployment

**Testing the Encryption:**
` + "```java\n// Verify the TDF was created\nFile tdfFile = new File(\"sensitive-data.tdf\");\nif (tdfFile.exists()) {\n    System.out.println(\"TDF size: \" + tdfFile.length() + \" bytes\");\n}\n```",

				`Complete OpenTDF Java encryption implementation with error handling:

**Configuration Class:**
` + "```java\nimport io.opentdf.platform.sdk.*;\n\npublic class OpenTDFConfig {\n    private static final String PLATFORM_URL = System.getenv(\"OPENTDF_PLATFORM_URL\");\n    private static final String CLIENT_ID = System.getenv(\"OPENTDF_CLIENT_ID\");\n    private static final String CLIENT_SECRET = System.getenv(\"OPENTDF_CLIENT_SECRET\");\n    \n    public static SDK createSDK() throws Exception {\n        if (PLATFORM_URL == null || CLIENT_ID == null || CLIENT_SECRET == null) {\n            throw new IllegalStateException(\"Required environment variables not set\");\n        }\n        \n        SDKConfig config = SDKConfig.builder()\n            .platformEndpoint(PLATFORM_URL)\n            .clientId(CLIENT_ID)\n            .clientSecret(CLIENT_SECRET)\n            .build();\n            \n        return new SDK(config);\n    }\n}\n```" + `

**Encryption Service:**
` + "```java\nimport java.util.List;\nimport java.util.Map;\n\npublic class TDFEncryptionService {\n    private final SDK sdk;\n    \n    public TDFEncryptionService() throws Exception {\n        this.sdk = OpenTDFConfig.createSDK();\n    }\n    \n    public void encryptFile(String inputPath, String outputPath, Map<String, String> attributes) {\n        try {\n            TDFConfig.Builder configBuilder = TDFConfig.builder();\n            \n            // Add data attributes\n            attributes.forEach(configBuilder::addDataAttribute);\n            \n            TDFConfig tdfConfig = configBuilder.build();\n            \n            // Perform encryption\n            try (InputStream input = Files.newInputStream(Paths.get(inputPath));\n                 OutputStream output = Files.newOutputStream(Paths.get(outputPath))) {\n                \n                sdk.createTDF(input, output, tdfConfig);\n                System.out.println(\"Successfully encrypted: \" + inputPath + \" -> \" + outputPath);\n            }\n        } catch (Exception e) {\n            System.err.println(\"Encryption failed: \" + e.getMessage());\n            throw new RuntimeException(\"TDF creation failed\", e);\n        }\n    }\n}\n```",
			},

			"docker_deploy": {
				`Here's a comprehensive guide to deploying OpenTDF with Docker:

**Docker Compose Setup:**
` + "```yaml\nversion: '3.8'\nservices:\n  opentdf-platform:\n    image: opentdf/platform:latest\n    ports:\n      - \"8080:8080\"\n    environment:\n      - DATABASE_URL=postgresql://postgres:password@db:5432/opentdf\n      - OIDC_ISSUER=https://your-identity-provider.com\n      - KAS_ENDPOINT=https://your-kas.com\n    depends_on:\n      - db\n      - kas\n      \n  kas:\n    image: opentdf/kas:latest\n    ports:\n      - \"8000:8000\"\n    environment:\n      - DATABASE_URL=postgresql://postgres:password@db:5432/kas\n      - PLATFORM_ENDPOINT=http://opentdf-platform:8080\n    depends_on:\n      - db\n      \n  db:\n    image: postgres:15\n    environment:\n      - POSTGRES_DB=opentdf\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=password\n    volumes:\n      - postgres_data:/var/lib/postgresql/data\n      \nvolumes:\n  postgres_data:\n```" + `

**Quick Start Commands:**
` + "```bash\n# Clone and start\ngit clone https://github.com/opentdf/platform.git\ncd platform\n\n# Start all services\ndocker-compose up -d\n\n# Check service health\ndocker-compose ps\n\n# View logs\ndocker-compose logs -f opentdf-platform\n```" + `

**Configuration:**
1. **Environment Variables**: Update OIDC_ISSUER for your identity provider
2. **Networking**: Ensure services can communicate
3. **Persistence**: Database data is persisted in volumes
4. **Security**: Change default passwords in production

**Verification:**
` + "```bash\n# Test platform health\ncurl http://localhost:8080/health\n\n# Test KAS health\ncurl http://localhost:8000/health\n\n# Check database connection\ndocker-compose exec db psql -U postgres -d opentdf -c \"\\dt\"\n```" + `

**Production Considerations:**
- Use environment-specific configuration files
- Implement proper secrets management
- Set up monitoring and logging
- Configure backup strategies
- Use reverse proxy (nginx) for SSL termination`,

				`Complete OpenTDF Docker deployment with production best practices:

**Production Docker Compose:**
` + "```yaml\nversion: '3.8'\nservices:\n  nginx:\n    image: nginx:alpine\n    ports:\n      - \"443:443\"\n      - \"80:80\"\n    volumes:\n      - ./nginx.conf:/etc/nginx/nginx.conf\n      - ./ssl:/etc/ssl/certs\n    depends_on:\n      - opentdf-platform\n      \n  opentdf-platform:\n    image: opentdf/platform:latest\n    environment:\n      - DATABASE_URL=${DATABASE_URL}\n      - OIDC_ISSUER=${OIDC_ISSUER}\n      - LOG_LEVEL=${LOG_LEVEL:-info}\n    volumes:\n      - ./config:/app/config\n    restart: unless-stopped\n    \n  kas:\n    image: opentdf/kas:latest\n    environment:\n      - DATABASE_URL=${KAS_DATABASE_URL}\n      - PLATFORM_ENDPOINT=${PLATFORM_ENDPOINT}\n    restart: unless-stopped\n    \n  db:\n    image: postgres:15-alpine\n    environment:\n      - POSTGRES_DB=${POSTGRES_DB}\n      - POSTGRES_USER=${POSTGRES_USER}\n      - POSTGRES_PASSWORD_FILE=/run/secrets/postgres_password\n    secrets:\n      - postgres_password\n    volumes:\n      - postgres_data:/var/lib/postgresql/data\n      - ./init.sql:/docker-entrypoint-initdb.d/init.sql\n      \nsecrets:\n  postgres_password:\n    file: ./secrets/postgres_password.txt\n    \nvolumes:\n  postgres_data:\n```" + `

**Environment Configuration (.env):**
` + "```bash\n# Platform Configuration\nDATABASE_URL=postgresql://opentdf:${POSTGRES_PASSWORD}@db:5432/opentdf\nOIDC_ISSUER=https://your-oidc-provider.com\nPLATFORM_ENDPOINT=https://platform.yourdomain.com\n\n# Database Configuration\nPOSTGRES_DB=opentdf\nPOSTGRES_USER=opentdf\n\n# Logging\nLOG_LEVEL=info\n```" + `

**Nginx Configuration (nginx.conf):**
` + "```nginx\nevents {}\nhttp {\n    upstream platform {\n        server opentdf-platform:8080;\n    }\n    \n    server {\n        listen 443 ssl;\n        server_name platform.yourdomain.com;\n        \n        ssl_certificate /etc/ssl/certs/cert.pem;\n        ssl_certificate_key /etc/ssl/certs/key.pem;\n        \n        location / {\n            proxy_pass http://platform;\n            proxy_set_header Host $host;\n            proxy_set_header X-Real-IP $remote_addr;\n        }\n    }\n}\n```",
			},

			"troubleshooting": {
				`When troubleshooting TDF encryption permission errors, follow this systematic approach:

**1. Verify Authentication:**
` + "```bash\n# Check if credentials are valid\notdfctl auth login --platform-endpoint https://your-platform.com\n\n# Verify token\notdfctl auth print-access-token\n```" + `

**2. Check Attribute Assignments:**
` + "```bash\n# List user's attributes\notdfctl policy subject-mappings list --subject-id your-user-id\n\n# Verify required attributes exist\notdfctl policy attributes values list --namespace-name organization.department\n```" + `

**3. Validate Policy Configuration:**
` + "```bash\n# Check namespace rules\notdfctl policy attributes namespaces get --name organization.department\n\n# Verify attribute values are active\notdfctl policy attributes values get --namespace-name organization.department --value engineering\n```" + `

**4. Debug Encryption Process:**
- **Check attribute syntax**: Ensure attribute names match exactly
- **Verify permissions**: User must have attributes required by policy
- **Test with minimal attributes**: Start with simple attribute set
- **Check network connectivity**: Ensure KAS is reachable

**5. Common Issues & Solutions:**

**Permission Denied:**
- User lacks required attributes
- Attribute namespace doesn't exist
- Incorrect attribute syntax in TDF creation

**Network Errors:**
- KAS endpoint unreachable
- Certificate validation failures
- Firewall blocking connections

**Configuration Issues:**
- Missing environment variables
- Incorrect platform endpoint
- Invalid client credentials

**Debugging Commands:**
` + "```bash\n# Enable verbose logging\nexport OTDFCTL_LOG_LEVEL=debug\n\n# Test simple encryption\notdfctl tdf encrypt --file test.txt --tdf-file test.tdf --attributes \"organization.department=engineering\"\n\n# Check KAS health\ncurl https://your-kas.com/health\n```" + `

**Log Analysis:**
Look for specific error patterns in logs:
- Authentication failures
- Attribute resolution errors  
- Network timeouts
- Policy evaluation denials`,
			},

			"default": {
				`This is a comprehensive response demonstrating the capabilities of the OpenTDF platform and evaluation system.

**OpenTDF Overview:**
OpenTDF (Open Trusted Data Format) is an open-source zero-trust data security platform that provides attribute-based access control (ABAC) for sensitive data protection.

**Key Features:**
- Embedded access policies that travel with encrypted data
- Fine-grained permissions based on user and data attributes
- Multiple format support (TDF, NanoTDF, ZTDF)
- Cross-platform SDK support (Java, Go, JavaScript, Python)
- Enterprise-ready with audit trails and compliance support

**Common Use Cases:**
1. **Secure Document Sharing**: Share sensitive documents across organizational boundaries
2. **Data Lake Protection**: Implement granular access controls in big data environments
3. **Compliance**: Meet regulatory requirements (GDPR, HIPAA, SOX)
4. **Zero-Trust Architecture**: Implement data-centric security models

**Getting Started:**
1. Install otdfctl CLI tool
2. Configure authentication with your OpenTDF platform
3. Create attribute namespaces and policies
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
		Model:        "realistic-mock-engine",
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
	case strings.Contains(queryLower, "what is opentdf") || strings.Contains(queryLower, "opentdf") && strings.Contains(queryLower, "work"):
		category = "opentdf_concept"
	case strings.Contains(queryLower, "attribute") && strings.Contains(queryLower, "namespace"):
		category = "attribute_namespace"
	case strings.Contains(queryLower, "java") && (strings.Contains(queryLower, "encrypt") || strings.Contains(queryLower, "code")):
		category = "java_encrypt"
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