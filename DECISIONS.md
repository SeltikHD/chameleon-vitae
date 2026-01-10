
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

## ADR-013: Use Jina Reader API for Job Description Parsing

**Date:** 2026-01-08  
**Status:** Accepted

### Context

Chameleon Vitae needs to parse job descriptions from various job boards (LinkedIn, Gupy, Indeed, etc.) to feed our LLM for resume tailoring. This presents several challenges:

1. **Modern Web Complexity:** Job posting pages are heavily JavaScript-rendered, making traditional scraping difficult.
2. **Anti-Scraping Measures:** Major job boards employ sophisticated bot detection (CAPTCHAs, rate limiting, IP blocking).
3. **Diverse Page Structures:** Each job board has different HTML structures, requiring custom parsers.
4. **Maintenance Burden:** Job boards frequently update their UI, breaking scrapers.
5. **Legal Concerns:** Direct scraping may violate terms of service.
6. **LLM Input Requirements:** Raw HTML is noisy and token-inefficient for LLM processing.

#### Alternatives Considered

|Option|Pros|Cons|
|--------|------|------|
|**Custom Scraper (Playwright/Puppeteer)**|Full control, no dependencies|High maintenance, bot detection, expensive hosting|
|**Browserless.io**|Managed headless browser|Expensive at scale ($0.01+ per page)|
|**ScrapingBee/Bright Data**|Residential proxies, high success rate|Very expensive, complex pricing|
|**Jina Reader API**|Free/cheap, LLM-optimized output, simple API|Third-party dependency, potential rate limits|
|**Ask users to paste text**|Simple, no scraping|Poor UX, inconsistent formatting|

### Decision

We will use **Jina Reader API** (`r.jina.ai`) as our primary job description parsing solution.

#### What is Jina Reader?

Jina Reader is an open-source service by Jina AI that converts any URL to LLM-friendly Markdown. It:

- Handles JavaScript rendering automatically
- Strips away navigation, ads, and irrelevant content
- Outputs clean Markdown optimized for LLM consumption
- Supports PDFs and images (with optional VLM captioning)
- Provides both free and paid tiers

#### Usage Pattern

```text
# Simply prepend r.jina.ai to any URL
GET https://r.jina.ai/https://linkedin.com/jobs/view/12345

# Returns clean Markdown:
## Senior Backend Engineer

**Company:** Awesome Corp
**Location:** Remote, USA

### About the Role
We are looking for an experienced backend engineer...

### Requirements
- 5+ years of experience with Go or similar languages
- Strong understanding of distributed systems
...
```

#### Integration Architecture

```text
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Frontend  │────▶│   Backend   │────▶│ Jina Reader │
│  (Job URL)  │      │   /parse    │      │  r.jina.ai  │
└─────────────┘      └─────────────┘      └─────────────┘
                            │
                            ▼
                     ┌─────────────┐
                     │   Groq AI   │
                     │  (Analysis) │
                     └─────────────┘
```

### Consequences

#### Positive

1. **Zero Maintenance:** No need to maintain scrapers or handle anti-bot measures.
2. **LLM-Optimized Output:** Markdown is token-efficient and preserves structure.
3. **Cost-Effective:** Free tier sufficient for MVP; paid tier is affordable (~$0.001/page).
4. **Simple Integration:** Single HTTP GET request with URL prefix.
5. **Rapid Development:** Allows us to focus on core resume tailoring logic.
6. **Scalable Infrastructure:** Jina maintains the infrastructure.

#### Negative

1. **Third-Party Dependency:** Service outage would affect job parsing functionality.
2. **Rate Limits:** Free tier has rate limits (may need API key for production).
3. **No Offline Mode:** Requires internet connectivity.
4. **Data Privacy:** Job URLs are sent to Jina's servers.
5. **Output Variability:** Markdown quality depends on page structure.

#### Mitigations

|Risk|Mitigation|
|------|------------|
|Service outage|Implement fallback: ask user to paste job description manually|
|Rate limits|Obtain API key for production; implement request queuing|
|Privacy concerns|Only public job URLs are parsed; no user data sent|
|Output quality|Validate and sanitize Markdown before LLM processing|

### Implementation Notes

#### API Endpoint

```go
// POST /tools/parse-job
type ParseJobRequest struct {
    URL string `json:"url" validate:"required,url"`
}

type ParseJobResponse struct {
    URL      string            `json:"url"`
    Title    string            `json:"title"`
    Markdown string            `json:"markdown"`
    Metadata map[string]string `json:"metadata"`
}
```

#### Headers for Production

```http
GET https://r.jina.ai/https://example.com/job
Authorization: Bearer <jina_api_key>
Accept: text/markdown
X-With-Generated-Alt: false
```

#### Error Handling

- **429 Too Many Requests:** Implement exponential backoff and retry.
- **Timeout (>30s):** Return partial content or fallback message.
- **Empty Content:** Prompt user to paste job description manually.

### Related Decisions

- [ADR-003: AI Strategy (LLM Provider)](./DECISIONS.md#adr-003-ai-strategy-llm-provider) - Jina output feeds into Groq/Llama
- [ADR-012: Multilanguage Resume Generation](./DECISIONS.md#adr-012-multilanguage-resume-generation) - Parsed content may be in different languages

### References

- [Jina Reader GitHub Repository](https://github.com/jina-ai/reader)
- [Jina AI Official Documentation](https://jina.ai/reader)
- [Jina Reader API Rate Limits](https://jina.ai/reader#pricing)

---

## ADR-014: Frontend Design System — "Cyber Chameleon"

**Date:** 2026-01-09  
**Status:** Accepted

---

### Decision

Adopt **Nuxt UI v4** as the primary component library with a custom **"Cyber Chameleon"** theme featuring a dark-mode-first design language.

---

### Context

The frontend needs a cohesive, professional design system that:

1. Reflects the project's identity (adaptability, AI-powered, modern)
2. Provides accessible, production-ready components out of the box
3. Minimizes custom CSS while maximizing consistency
4. Supports dark mode as the primary experience (developers and power users preference)

We evaluated several options:

|Option|Pros|Cons|
|------|----|----|
|**Nuxt UI v4**|Official Nuxt module, 125+ components, Tailwind v4, free & open-source|Opinionated styling|
|Vuetify 3|Material Design, extensive components|Heavy bundle, not Nuxt-native|
|PrimeVue|Large component library|Not Tailwind-based|
|Headless UI + custom|Maximum flexibility|Significant development effort|

**Nuxt UI v4** was chosen because it:

- Is built on **Reka UI** (accessible, WAI-ARIA compliant)
- Uses **Tailwind CSS v4** with CSS-first configuration
- Provides **Tailwind Variants** for type-safe component variants
- Unifies Nuxt UI + Nuxt UI Pro into a single open-source package

---

### Theme: "Cyber Chameleon"

#### Design Philosophy

- **Dark Mode First:** The primary experience is dark mode, optimized for developer tools and focused work sessions.
- **High Contrast:** Text and interactive elements have sufficient contrast ratios for accessibility.
- **AI Accent:** Violet/purple tones represent AI-powered features, creating visual distinction for "magic" moments.
- **Nature Meets Tech:** Emerald green evokes the chameleon identity while serving as the primary action color.

#### Color Palette

|Role|Color Name|Hex Code|Tailwind Shade|Usage|
|----|----------|--------|--------------|-----|
|**Primary**|Emerald|`#10b981`|`emerald-500`|Buttons, links, primary actions, success states|
|**Secondary**|Violet|`#8b5cf6`|`violet-500`|AI features, highlights, premium actions|
|**Background**|Zinc|`#09090b`|`zinc-950`|Page background, modals|
|**Surface**|Zinc|`#18181b`|`zinc-900`|Cards, sidebars, elevated surfaces|
|**Text**|Zinc|`#f4f4f5`|`zinc-100`|Primary text, headings|
|**Muted**|Zinc|`#a1a1aa`|`zinc-400`|Secondary text, placeholders, disabled states|
|**Border**|Zinc|`#27272a`|`zinc-800`|Dividers, card borders|
|**Error**|Red|`#ef4444`|`red-500`|Error states, destructive actions|
|**Warning**|Amber|`#f59e0b`|`amber-500`|Warning states, caution|
|**Info**|Sky|`#0ea5e9`|`sky-500`|Informational states|

#### Visual Hierarchy

```text
┌────────────────────────────────────────────────────────────┐
│ Background (zinc-950)                                      │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Surface (zinc-900)                                   │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ Card (zinc-900 + border zinc-800)              │  │  │
│  │  │                                                │  │  │
│  │  │  [Primary Button] ← emerald-500                │  │  │
│  │  │  [Secondary Button] ← violet-500               │  │  │
│  │  │                                                │  │  │
│  │  │  Text (zinc-100)                               │  │  │
│  │  │  Muted text (zinc-400)                         │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘
```

---

### Implementation

#### 1. Nuxt Configuration (`nuxt.config.ts`)

The color system is registered in the Nuxt configuration:

```typescript
export default defineNuxtConfig({
  modules: ['@nuxt/ui'],
  ui: {
    theme: {
      colors: ['primary', 'secondary', 'success', 'warning', 'error', 'info']
    }
  }
})
```

#### 2. App Configuration (`app.config.ts`)

Runtime color mapping is defined in `app.config.ts`:

```typescript
export default defineAppConfig({
  ui: {
    colors: {
      primary: 'emerald',    // Actions, links, buttons
      secondary: 'violet',   // AI features, highlights
      neutral: 'zinc',       // Backgrounds, borders, text
      success: 'emerald',    // Success states
      warning: 'amber',      // Warning states
      error: 'red',          // Error states
      info: 'sky'            // Informational states
    }
  }
})
```

#### 3. CSS Theme Layer (`app/assets/css/main.css`)

Additional theme customizations are applied via Tailwind's `@theme` directive:

```css
@import "tailwindcss";
@import "@nuxt/ui";

@theme {
  /* Typography */
  --font-sans: 'Inter', ui-sans-serif, system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', ui-monospace, monospace;
}
```

---

### Component Guidelines

#### Button Variants

|Variant|Color|Use Case|
|-------|-----|--------|
|`solid` (default)|`primary`|Primary actions (Submit, Save, Create)|
|`solid`|`secondary`|AI-powered actions (Generate, Analyze)|
|`outline`|`primary`|Secondary actions (Cancel, Back)|
|`ghost`|`neutral`|Tertiary actions (Edit, View)|
|`soft`|`error`|Destructive actions (Delete, Remove)|

#### Cards

- Use `UCard` with default styling (inherits surface background)
- Add `divide-y divide-zinc-800` for sectioned cards
- Use `shadow-lg` sparingly for elevated modals

#### Forms

- Use `UFormField` for consistent label/error placement
- Use `UInput`, `UTextarea`, `USelect` from Nuxt UI
- Apply `size="lg"` for primary inputs

---

### Consequences

#### Positive

- **Consistency:** Single source of truth for colors and component styling
- **Accessibility:** Reka UI base ensures WCAG 2.1 AA compliance
- **Developer Experience:** Tailwind Variants provide type-safe component customization
- **Performance:** Tailwind CSS v4 with CSS-first config reduces bundle size
- **Maintainability:** Theme changes propagate automatically via CSS variables

#### Negative

- **Opinionated:** Nuxt UI's styling may require overrides for unique designs
- **Learning Curve:** Developers must learn Nuxt UI's component API
- **Dark Mode Focus:** Light mode is secondary, may need additional testing

#### Mitigations

- Document component usage patterns in this ADR
- Create Storybook-like examples in `/docs/components/`
- Test light mode toggle before each release

---

### References

- [Nuxt UI Documentation](https://ui.nuxt.com/)
- [Tailwind CSS v4](https://tailwindcss.com/blog/tailwindcss-v4)
- [Reka UI (Accessibility)](https://reka-ui.com/)
- [WCAG 2.1 Contrast Guidelines](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html)

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

Last Updated: 2026-01-10
