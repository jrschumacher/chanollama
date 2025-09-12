# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records documenting significant architectural decisions for the chanollama project.

## ADR Index

| ADR | Title | Status | Category |
|-----|-------|--------|----------|
| [001](./001-api-consolidation.md) | API Entry Point Consolidation | Proposed | Developer Experience |
| [002](./002-configuration-unification.md) | Configuration Structure Unification | Proposed | Developer Experience |

## Background

These ADRs emerged from a comprehensive architectural cleanup initiative that identified and addressed domain bias and API complexity issues in chanollama:

### ✅ Completed (Phase 1)
- **Domain Bias Elimination**: Transformed security-focused test data to demonstrate cross-industry applicability
- **Agent Type Cleanup**: Removed hardcoded domain-specific agent types in favor of developer-defined types

### 📋 Proposed (Phase 2)
- **ADR-001**: Consolidate multiple confusing entry points into a single, composable API
- **ADR-002**: Unify three overlapping configuration structures into one coherent system

## Implementation Priority

1. **ADR-002 (Configuration Unification)** should be implemented first, as it provides the foundation for ADR-001
2. **ADR-001 (API Consolidation)** builds on the unified configuration to provide clean functional options

Both ADRs are designed to be **non-breaking** during implementation, with gradual deprecation of old patterns.

## Success Metrics

### Current State (After Phase 1)
- ✅ **Domain Neutrality**: Library demonstrates cross-industry flexibility
- ✅ **Agent Flexibility**: Agent types completely developer-defined
- ✅ **Zero Regressions**: All 78+ tests pass

### Target State (After ADRs Implementation)
- 🎯 **API Simplicity**: Single entry point instead of 4+ confusing constructors
- 🎯 **Configuration Clarity**: One config structure instead of 3 overlapping ones
- 🎯 **Developer Experience**: 15-minute learning curve instead of research-heavy onboarding
- 🎯 **Maintenance Simplicity**: Unified codebase instead of multiple parallel patterns

## Review Process

1. **Review ADRs** for technical accuracy and architectural soundness
2. **Validate approach** against real-world use cases
3. **Approve implementation** plan and phasing strategy  
4. **Track implementation** progress against success metrics

---

These architectural improvements will transform chanollama from a complex, seemingly domain-specific library into a clean, general-purpose offline LLM solution that developers across all industries can confidently adopt.