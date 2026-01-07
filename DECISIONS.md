
# Architecture Decision Records (ADR)

This document catalogs the fundamental technical decisions made during the conception of **Chameleon Vitae**, their context, and consequences.

---

## ADR-001: Project Identity and Purpose

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

The project will be named **"Chameleon Vitae"**.

### Context

We needed a name that evokes adaptation (Chameleon) while maintaining the professional acronym "CV" (Curriculum Vitae).

### Consequences

- **Positive:** Visual identity aligned with the idea of "mimicry" to pass through HR/ATS filters.
- **Positive:** Memorable and unique branding.
- **Negative:** May require explanation in non-English speaking markets.

---

## ADR-002: Hexagonal Architecture (Ports and Adapters)

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

Adopt **Hexagonal Architecture (Ports and Adapters)** as the foundational architectural pattern.

### Context

The system depends heavily on unstable or swappable services (AI APIs, PDF Engines, Databases). We need to isolate the core business logic from external dependencies to enable:

- Easy testing of domain logic
- Swapping providers without refactoring business rules
- Clear separation of concerns

### Consequences

- **Positive:** The Core Domain (rules about how a resume is assembled) doesn't know that the AI is Llama or the database is PostgreSQL. Facilitates unit testing and future swap from Groq to Ollama (Local) without refactoring business rules.
- **Negative:** Initial increase in boilerplate (interfaces/adapters) compared to a simple MVC architecture.

### Directory Structure

```text
internal/
├── core/
│   ├── domain/     # Entities, Value Objects, Domain Errors
│   ├── ports/      # Interfaces (Input & Output Ports)
│   └── services/   # Application Services / Use Cases
└── adapters/
    ├── primary/    # Input Adapters (HTTP handlers)
    └── secondary/  # Output Adapters (DB, AI, PDF)
```

---

## ADR-003: AI Strategy (LLM Provider)

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

- **Provider:** Use **Groq API** (Free Tier) initially.
- **Models:**
  - **Llama 3.3 70b-versatile:** For final text generation (creative rewriting).
  - **Llama 4 Scout 17b:** For job analysis and data extraction (due to high TPM).

### Context

We need high-quality text output (70b model) and high token throughput (TPM) without local infrastructure costs for the MVP. Future migration to Ollama Local is considered.

### Consequences

- **Positive:** High-quality output without infrastructure costs.
- **Positive:** Extremely fast inference via Groq's LPU technology.
- **Negative:** Dependency on internet connection for MVP functionality.
- **Mitigation:** Data model prepared for fallback or engine swap.

---

## ADR-004: Data Modeling (Atomic Experience Bullets)

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

Adopt **Atomic Granularity** for experiences using "Experience Bullets".

### Context

Traditional resumes store giant text blocks. This prevents customization and fine-grained selection.

### Consequences

- **Positive:** Experiences are broken into relational tables (`experiences` 1:N `bullets`).
- **Positive:** The AI selects *which* lines to include, rather than just summarizing a large text.
- **Positive:** Increases precision and relevance of generated CVs.
- **Negative:** Requires user to enter experiences in a more structured way.

### Schema Design

```sql
experiences (id, user_id, type, title, organization, ...)
    └── bullets (id, experience_id, content, impact_score, keywords[], ...)
```

---

## ADR-005: Document Generation (PDF Engine)

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

Use **Gotenberg 8** (Headless Chromium) for PDF rendering.

### Context

Native PDF libraries (e.g., `gofpdf`) are difficult to style. HTML+CSS is universal and designers can work with it.

### Consequences

- **Positive:** Go backend generates only HTML. Gotenberg container converts to PDF.
- **Positive:** Final PDF is machine-readable (selectable text), a fundamental criterion for ATS.
- **Positive:** Easy to create and modify resume templates using standard web technologies.
- **Negative:** Additional container dependency in infrastructure.

---

## ADR-006: Container Runtime

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

Full containerization via **Podman** (rootless mode).

### Context

Preference for security and daemonless/rootless execution. Compatible with Docker images.

### Consequences

- **Positive:** Development environment is identical to production.
- **Positive:** Enhanced security through rootless containers.
- **Positive:** Use of `podman-compose` for orchestrating Database and PDF Engine.
- **Negative:** Some users may need to learn Podman commands.

---

## ADR-007: Frontend Framework

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

Use **Vue.js 3 + Nuxt 3** with Static Site Generation (SSG).

### Context

We need a reactive interface for managing the bullets database, but one that is lightweight and easy to host.

### Consequences

- **Positive:** Use of TailwindCSS for quick styling of resume templates.
- **Positive:** Utility classes are easier to manipulate via code/AI if necessary.
- **Positive:** SSG allows for easy static hosting.
- **Negative:** Adds Node.js to the development toolchain.

---

## ADR-008: Database Selection

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

Use **PostgreSQL 17** as the primary database.

### Context

We evaluated PostgreSQL 17 vs 18. While PostgreSQL 18.1 is available, PostgreSQL 17 offers:

- Better stability and wider community support
- All features we need (UUID, JSONB, pg_trgm for fuzzy search)
- LTS-like reliability for production use

### Consequences

- **Positive:** Robust relational database with excellent JSONB support.
- **Positive:** UUID extension for primary keys.
- **Positive:** pg_trgm extension for fuzzy skill matching.
- **Negative:** N/A - PostgreSQL is well-established and understood.

---

## ADR-009: Backend Framework

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

Use **Go with Chi Router** for the HTTP layer.

### Context

We needed a lightweight, idiomatic Go HTTP router that:

- Has minimal dependencies
- Follows standard `net/http` patterns
- Provides good middleware support

### Consequences

- **Positive:** Chi is lightweight and idiomatic.
- **Positive:** Follows `net/http` standards closely.
- **Positive:** Good middleware ecosystem.
- **Negative:** Less feature-rich than larger frameworks (intentional trade-off).

---

## ADR-010: Language Policy

**Date:** 2026-01-06  
**Status:** Accepted

### Decision

All code, comments, documentation, and commit messages must be in **English**.

### Context

As an open-source project intended for global contribution, a single language standard is necessary.

### Consequences

- **Positive:** Consistent codebase accessible to international contributors.
- **Positive:** Better integration with AI coding assistants.
- **Negative:** May slow down contributors whose primary language is not English.

---

## ADR-011: Firebase Authentication

**Date:** 2026-01-07  
**Status:** Accepted

### Decision

Use **Firebase Authentication** as the primary identity provider, with support for email/password and social logins (Google, GitHub).

### Context

We need a robust authentication system but want to:

- Avoid implementing complex auth flows from scratch (password reset, email verification, OAuth flows)
- Accelerate development by leveraging a battle-tested auth service
- Support multiple authentication methods without additional complexity
- Maintain flexibility for future migration if needed

Firebase Authentication provides:

- Pre-built UI components and SDKs
- Secure token-based authentication (JWT)
- Built-in support for popular OAuth providers
- Free tier sufficient for MVP and early growth

### Consequences

- **Positive:** Dramatically reduces auth implementation time.
- **Positive:** Battle-tested security from Google.
- **Positive:** Easy integration with frontend via Firebase SDK.
- **Positive:** JWT tokens can be verified server-side without Firebase SDK dependency.
- **Negative:** External dependency for authentication.
- **Negative:** Vendor lock-in (mitigated by storing only `firebase_uid` in DB).

### Implementation Notes

- Backend validates Firebase ID tokens using public keys (no Firebase Admin SDK needed for basic validation).
- Users table stores `firebase_uid` (VARCHAR 128) as the unique identifier.
- Email is optional in DB (can be fetched from token claims).
- Consider adding Firebase Admin SDK for advanced features (user management, custom claims).

---

## ADR-012: Multilanguage Resume Generation

**Date:** 2026-01-07  
**Status:** Accepted

### Decision

Support **multilanguage resume generation** with prompts stored in external JSON files.

### Context

- Resumes need to be generated in different languages (e.g., English, Brazilian Portuguese).
- AI prompts must be version-controlled and easily editable without code changes.
- The system should be extensible to add new languages without code modifications.

### Consequences

- **Positive:** Prompts are externalized and can be modified by non-developers.
- **Positive:** Easy to add new languages by adding new JSON files.
- **Positive:** Version control for prompts enables tracking changes over time.
- **Positive:** Separation of concerns between code logic and prompt content.
- **Negative:** Need to implement prompt loading and caching mechanism.

### File Structure

```text
internal/
└── adapters/
    └── secondary/
        └── groq/
            └── prompts/
                ├── en/
                │   ├── analyze_job.json
                │   ├── select_bullets.json
                │   └── rewrite_bullet.json
                └── pt-br/
                    ├── analyze_job.json
                    ├── select_bullets.json
                    └── rewrite_bullet.json
```

### Prompt JSON Format

```json
{
  "name": "analyze_job",
  "description": "Analyzes job description to extract requirements",
  "system_prompt": "You are an expert HR analyst...",
  "user_prompt_template": "Analyze the following job description:\n\n{{job_description}}",
  "temperature": 0.3,
  "max_tokens": 2000
}
```

---

## Template for New ADRs

```markdown
## ADR-XXX: [Title]

**Date:** [YYYY-MM-DD]  
**Status:** [Proposed | Accepted | Deprecated | Superseded]

### Decision
[What is the change that we're proposing and/or doing?]

### Context
[What is the issue that we're seeing that is motivating this decision?]

### Consequences
- **Positive:** [What becomes easier?]
- **Negative:** [What becomes harder?]
- **Risks:** [What could go wrong?]
```

---

Last Updated: 2026-01-07
