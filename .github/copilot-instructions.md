# GitHub Copilot Instructions for Chameleon Vitae

> **"The CV that adapts."** — AI-powered resume engineering using Hexagonal Architecture.

## Project Overview

Chameleon Vitae is a resume engineering tool that uses AI (Groq/Llama) to tailor resumes to specific job descriptions. The system maintains a database of atomic "experience bullets" and intelligently selects and rewrites them for each job application.

## Technology Stack

| Layer       | Technology         | Version     |
|-------------|--------------------|-------------|
| Backend     | Go + Chi Router    | Go 1.22+    |
| Frontend    | Vue.js + Nuxt      | Nuxt 3 / Vue 3 |
| Database    | PostgreSQL         | 17+         |
| PDF Engine  | Gotenberg          | 8           |
| AI Provider | Groq API           | Llama 3.3/4 |
| Containers  | Podman             | Rootless    |

---

## 🚨 CRITICAL RULES — READ FIRST

### Confidence Threshold

> **If you are not 95% confident in your solution, you MUST ask clarifying questions before generating code.**

Do not guess. Do not assume. If the requirements are ambiguous or you're unsure about architectural boundaries, stop and ask.

### Post-Action Summary

> **At the end of every response, provide a brief summary of:**
> 1. What was done
> 2. Any files created or modified
> 3. **Explicitly list any Business Rules or Domain Logic that were touched or modified**

This helps maintain traceability and ensures reviewers can quickly understand the impact.

### Language Policy

> **STRICTLY ENGLISH.** All code, comments, variable names, documentation, commit messages, and error messages must be in English. No exceptions.

---

## 🏛️ Hexagonal Architecture Enforcement

This project follows **Hexagonal Architecture (Ports and Adapters)**. This is non-negotiable.

### Directory Structure

```
chameleon-vitae/
├── cmd/
│   └── server/             # Application entrypoint
│       └── main.go
├── internal/
│   ├── core/               # 🔒 PURE DOMAIN — NO EXTERNAL DEPENDENCIES
│   │   ├── domain/         # Entities, Value Objects, Domain Errors
│   │   ├── ports/          # Interfaces (Input & Output Ports)
│   │   └── services/       # Application Services / Use Cases
│   └── adapters/
│       ├── primary/        # Input Adapters (HTTP handlers, CLI)
│       │   └── http/       # Chi router handlers
│       └── secondary/      # Output Adapters (implementations)
│           ├── postgres/   # Database adapter
│           ├── groq/       # AI provider adapter
│           └── gotenberg/  # PDF engine adapter
├── pkg/                    # Shared utilities (can be imported by adapters)
├── frontend/               # Nuxt.js application
└── deploy/                 # Infrastructure (Compose, Dockerfiles)
```

### 🚫 ABSOLUTE PROHIBITION

```go
// ❌ NEVER DO THIS — Importing adapter into core
package services

import (
    "github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/postgres" // FORBIDDEN!
)
```

```go
// ✅ CORRECT — Core only knows about ports (interfaces)
package services

import (
    "github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

type ResumeService struct {
    repo ports.ResumeRepository  // Interface, not implementation
    ai   ports.AIProvider        // Interface, not implementation
}
```

### Dependency Rules

1. **Core Domain (`internal/core/`)** — ZERO external dependencies
   - Only standard library
   - Only its own packages (`domain`, `ports`, `services`)
   - Defines interfaces (ports) that adapters must implement

2. **Adapters (`internal/adapters/`)** — Implement core interfaces
   - Import from `internal/core/ports`
   - Import external libraries (pgx, chi, etc.)
   - NEVER import other adapters

3. **Cmd (`cmd/`)** — Wires everything together
   - Creates adapter instances
   - Injects adapters into services
   - Starts the application

---

## 📐 Code Style Guidelines

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- Use `gofmt` and `golangci-lint`
- Error handling: Always handle errors explicitly, wrap with context
- Naming: Use clear, descriptive names; avoid abbreviations
- Comments: Document all exported types and functions
- Tests: Write table-driven tests

```go
// ✅ Good error handling
if err != nil {
    return fmt.Errorf("failed to create resume: %w", err)
}

// ✅ Good naming
type ResumeGenerationService struct { ... }
func (s *ResumeGenerationService) GenerateForJob(ctx context.Context, jobDesc string) (*Resume, error)

// ❌ Bad naming
type RGS struct { ... }
func (s *RGS) Gen(j string) (*R, error)
```

### TypeScript/Vue Code

- Use TypeScript strict mode
- Follow Vue 3 Composition API with `<script setup>`
- Use the project's ESLint configuration
- Components use PascalCase, composables use `use` prefix

### SQL

- Use parameterized queries (NEVER concatenate user input)
- Name constraints explicitly
- Use UUID for primary keys
- Include `created_at` and `updated_at` timestamps

---

## 🧪 Testing Requirements

### Backend

- Unit tests for all core domain logic
- Integration tests for adapters
- Use table-driven tests in Go
- Mock interfaces, not implementations

```go
// ✅ Good test structure
func TestResumeService_GenerateForJob(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Resume
        wantErr bool
    }{
        // test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

### Frontend

- Unit tests for composables and utilities
- Component tests for complex components
- E2E tests for critical user flows

---

## 📝 Documentation Standards

- README.md must stay up to date
- Document architectural decisions in DECISIONS.md (ADR format)
- API endpoints documented with examples
- Complex functions require doc comments

---

## 🔒 Security Considerations

- Never log sensitive data (API keys, passwords)
- Validate all user input
- Use parameterized SQL queries
- Environment variables for secrets
- Gotenberg runs isolated (no external network access)

---

## 🤖 AI Integration Guidelines

When working with Groq API:

- Handle rate limits gracefully
- Implement retry with exponential backoff
- Cache AI responses when appropriate
- Log token usage for monitoring
- Keep prompts in separate, version-controlled files

---

## 📊 Database Guidelines

- All schema changes go in `deploy/postgres/init/`
- Use migrations for production changes
- JSONB for flexible metadata
- Proper indexes for query patterns
- Foreign keys with appropriate cascades

---

## 🎯 Before Generating Code

Ask yourself:

1. Does this belong in core, adapters, or cmd?
2. Am I importing anything forbidden into core?
3. Am I handling all errors properly?
4. Are there tests for this functionality?
5. Is the code in English?
6. Would a reviewer understand this without explanation?

---

## 📋 Response Format

When providing code or making changes, structure your response as:

1. **Analysis** — Brief explanation of what you understand
2. **Implementation** — The actual code/changes
3. **Summary** — What was done, files modified
4. **Domain Impact** — List any business rules or domain logic affected

---

*Last Updated: January 2026*
