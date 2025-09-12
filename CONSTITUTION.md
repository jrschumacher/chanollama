# chanollama Constitution

This document defines the core architectural principles, constraints, and decision-making frameworks that guide all development decisions for chanollama.

## Mission Statement

**chanollama is a truly offline LLM inference library for Go that enables developers across ALL industries to embed Large Language Models directly into their applications with zero network dependencies.**

## Core Principles

### 1. Offline-Only Architecture (Non-Negotiable)
- **NEVER add HTTP support or network communication**
- **Direct model execution** using llama.cpp bindings only
- **Zero external service dependencies** - applications must work completely offline
- **Embeddable distribution** - models can be included directly in application binaries
- **File-based models only** - accept file paths like `/path/to/model.gguf`, not model names

### 2. Domain Neutrality (Non-Negotiable)
- **No hardcoded industry assumptions** - library must work for ANY domain
- **Developer-defined semantics** - let developers define their own use cases, agent types, and workflows
- **Cross-industry examples** - demonstrate e-commerce, healthcare, finance, content management equally
- **Avoid domain-specific terminology** in APIs, documentation, and examples

### 3. Developer Experience First
- **Simplicity over cleverness** - prefer clear, obvious APIs over clever abstractions
- **Single clear path** - provide ONE obvious way to accomplish common tasks
- **Composable features** - optional features that can be added when needed
- **Discoverable APIs** - IDE autocompletion should reveal available options
- **Comprehensive examples** - real-world usage patterns for different industries

### 4. Go Idioms and Best Practices
- **Channel-based streaming** - use Go's concurrency primitives for real-time token delivery
- **Interface-based design** - depend on interfaces, not concrete implementations
- **Error handling** - explicit error returns, no panics in public APIs
- **Context support** - respect context.Context for cancellation and timeouts
- **Memory safety** - proper resource cleanup and lifecycle management

### 5. Backward Compatibility
- **Non-breaking changes preferred** - maintain API compatibility when possible
- **Deprecation process** - clear migration paths and deprecation warnings
- **Semantic versioning** - follow semver strictly for version management
- **Migration guides** - comprehensive documentation for breaking changes

## Architectural Constraints

### Code Organization
- **Single responsibility** - each package/type should have one clear purpose
- **Interface segregation** - small, focused interfaces over large ones
- **Dependency inversion** - depend on abstractions, not concretions
- **Composition over inheritance** - favor struct embedding and composition

### Performance Requirements
- **Memory efficiency** - careful management of model loading and context
- **Concurrent safety** - all public APIs must be goroutine-safe
- **Resource cleanup** - proper cleanup of models, contexts, and file handles
- **Streaming performance** - sub-100ms token delivery latency

### Testing Standards
- **>90% test coverage** - comprehensive testing of all public APIs
- **No external dependencies in tests** - tests must run completely offline
- **Deterministic tests** - no flaky tests that depend on timing or external state
- **Example validation** - all examples must be tested and functional

## Decision-Making Framework

### API Design Decisions
1. **Simplicity Check**: Is this the simplest way to accomplish the goal?
2. **Domain Neutrality Check**: Does this work equally well for all industries?
3. **Go Idioms Check**: Does this follow Go best practices and conventions?
4. **Offline Check**: Does this maintain zero network dependencies?
5. **Backward Compatibility Check**: Does this preserve existing functionality?

### Feature Addition Criteria
- **General Purpose**: Feature must be useful across multiple domains/industries
- **Offline Compatible**: Feature must work without network access
- **Composable**: Feature should be optional and composable with other features
- **Well-Tested**: Feature must include comprehensive tests and examples
- **Documented**: Feature must include clear documentation and usage examples

### Breaking Change Approval Process
1. **Necessity Justification**: Why is the breaking change necessary?
2. **Migration Path**: How will existing users migrate to the new API?
3. **Deprecation Timeline**: What is the timeline for deprecation warnings?
4. **Impact Assessment**: How many users will be affected?
5. **Alternative Evaluation**: Have non-breaking alternatives been considered?

## Quality Standards

### Code Quality
- **Linting**: All code must pass `go vet`, `golint`, and `gofmt`
- **Error Handling**: All errors must be properly handled and meaningful
- **Documentation**: All public APIs must have godoc comments
- **Examples**: All public APIs should have usage examples

### Security Standards
- **No secrets in code** - never hardcode API keys, tokens, or credentials
- **Input validation** - validate all inputs to public functions
- **Safe defaults** - default configurations should be secure
- **Dependency scanning** - monitor dependencies for security vulnerabilities

### Performance Standards
- **Benchmarking**: Critical paths must have benchmark tests
- **Memory profiling**: Monitor memory usage and prevent leaks
- **Latency targets**: Streaming operations <100ms, batch operations <5s
- **Resource limits**: Configurable limits on memory and concurrent operations

## Anti-Patterns (Forbidden)

### Network Dependencies
- ❌ HTTP clients or servers
- ❌ Network timeouts or connectivity checks
- ❌ External service integrations
- ❌ Model downloads or remote model references

### Domain Bias
- ❌ Hardcoded industry-specific terminology
- ❌ Assumptions about use cases or workflows
- ❌ Examples that only work for specific domains
- ❌ Agent types or classifications that assume specific industries

### API Complexity
- ❌ Multiple ways to do the same thing
- ❌ Unclear constructor proliferation
- ❌ Overlapping configuration structures
- ❌ God objects that do everything

### Go Anti-Patterns
- ❌ Panics in public APIs
- ❌ Ignoring context.Context
- ❌ Not handling errors
- ❌ Race conditions in public APIs

## Success Metrics

### Developer Experience Metrics
- **Time to First Success**: <15 minutes from installation to working example
- **API Discoverability**: IDE autocompletion reveals appropriate next steps
- **Documentation Quality**: Developers can accomplish tasks using only documentation
- **Cross-Domain Usage**: Examples and tests demonstrate usage across multiple industries

### Technical Metrics
- **Test Coverage**: >90% coverage for all packages
- **Performance**: <100ms streaming latency, <5s batch operations
- **Memory Usage**: Predictable memory usage with configurable limits
- **Build Time**: <30s for full build and test cycle

### Ecosystem Metrics
- **Issue Resolution**: <48h response time for issues
- **Adoption Growth**: Steady increase in usage across different domains
- **Community Contributions**: Regular community contributions and feedback
- **Backward Compatibility**: No unplanned breaking changes

## Evolution Guidelines

### Adding New Features
1. Create ADR (Architecture Decision Record) for significant changes
2. Ensure feature aligns with domain neutrality principle
3. Design with offline-only constraint in mind
4. Include comprehensive tests and examples
5. Update documentation and migration guides

### Deprecating Features
1. Announce deprecation with clear timeline (minimum 2 versions)
2. Provide migration guide with examples
3. Include deprecation warnings in code
4. Monitor usage before removal
5. Remove in major version with comprehensive changelog

### Breaking Changes
1. Must be approved through ADR process
2. Require migration guide and tooling where possible
3. Must provide significant value over existing approach
4. Should be batched into major version releases
5. Require comprehensive testing and validation

## Enforcement

### Code Review Requirements
- All changes must maintain constitutional compliance
- Performance impact must be assessed for critical paths
- Breaking changes require explicit approval
- Tests must demonstrate offline-only operation

### Release Criteria
- All tests pass in offline environment
- Documentation is complete and accurate
- Examples work across multiple domains
- Performance benchmarks meet targets
- Security scan passes with no critical issues

---

This constitution serves as the foundation for all architectural decisions in chanollama. Any code, feature, or change that violates these principles should be rejected or redesigned to align with the constitutional requirements.

**Last Updated**: 2025-01-20
**Version**: 1.0.0