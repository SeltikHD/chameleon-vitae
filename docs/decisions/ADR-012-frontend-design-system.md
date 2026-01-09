# ADR-012: Frontend Design System — "Cyber Chameleon"

**Date:** 2026-01-09  
**Status:** Accepted

---

## Decision

Adopt **Nuxt UI v4** as the primary component library with a custom **"Cyber Chameleon"** theme featuring a dark-mode-first design language.

---

## Context

The frontend needs a cohesive, professional design system that:

1. Reflects the project's identity (adaptability, AI-powered, modern)
2. Provides accessible, production-ready components out of the box
3. Minimizes custom CSS while maximizing consistency
4. Supports dark mode as the primary experience (developers and power users preference)

We evaluated several options:

| Option | Pros | Cons |
|--------|------|------|
| **Nuxt UI v4** | Official Nuxt module, 125+ components, Tailwind v4, free & open-source | Opinionated styling |
| Vuetify 3 | Material Design, extensive components | Heavy bundle, not Nuxt-native |
| PrimeVue | Large component library | Not Tailwind-based |
| Headless UI + custom | Maximum flexibility | Significant development effort |

**Nuxt UI v4** was chosen because it:

- Is built on **Reka UI** (accessible, WAI-ARIA compliant)
- Uses **Tailwind CSS v4** with CSS-first configuration
- Provides **Tailwind Variants** for type-safe component variants
- Unifies Nuxt UI + Nuxt UI Pro into a single open-source package

---

## Theme: "Cyber Chameleon"

### Design Philosophy

- **Dark Mode First:** The primary experience is dark mode, optimized for developer tools and focused work sessions.
- **High Contrast:** Text and interactive elements have sufficient contrast ratios for accessibility.
- **AI Accent:** Violet/purple tones represent AI-powered features, creating visual distinction for "magic" moments.
- **Nature Meets Tech:** Emerald green evokes the chameleon identity while serving as the primary action color.

### Color Palette

| Role | Color Name | Hex Code | Tailwind Shade | Usage |
|------|------------|----------|----------------|-------|
| **Primary** | Emerald | `#10b981` | `emerald-500` | Buttons, links, primary actions, success states |
| **Secondary** | Violet | `#8b5cf6` | `violet-500` | AI features, highlights, premium actions |
| **Background** | Zinc | `#09090b` | `zinc-950` | Page background, modals |
| **Surface** | Zinc | `#18181b` | `zinc-900` | Cards, sidebars, elevated surfaces |
| **Text** | Zinc | `#f4f4f5` | `zinc-100` | Primary text, headings |
| **Muted** | Zinc | `#a1a1aa` | `zinc-400` | Secondary text, placeholders, disabled states |
| **Border** | Zinc | `#27272a` | `zinc-800` | Dividers, card borders |
| **Error** | Red | `#ef4444` | `red-500` | Error states, destructive actions |
| **Warning** | Amber | `#f59e0b` | `amber-500` | Warning states, caution |
| **Info** | Sky | `#0ea5e9` | `sky-500` | Informational states |

### Visual Hierarchy

```text
┌─────────────────────────────────────────────────────────────┐
│ Background (zinc-950)                                       │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Surface (zinc-900)                                    │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ Card (zinc-900 + border zinc-800)              │  │  │
│  │  │                                                 │  │  │
│  │  │  [Primary Button] ← emerald-500                │  │  │
│  │  │  [Secondary Button] ← violet-500               │  │  │
│  │  │                                                 │  │  │
│  │  │  Text (zinc-100)                               │  │  │
│  │  │  Muted text (zinc-400)                         │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## Implementation

### 1. Nuxt Configuration (`nuxt.config.ts`)

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

### 2. App Configuration (`app.config.ts`)

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

### 3. CSS Theme Layer (`app/assets/css/main.css`)

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

## Component Guidelines

### Button Variants

| Variant | Color | Use Case |
|---------|-------|----------|
| `solid` (default) | `primary` | Primary actions (Submit, Save, Create) |
| `solid` | `secondary` | AI-powered actions (Generate, Analyze) |
| `outline` | `primary` | Secondary actions (Cancel, Back) |
| `ghost` | `neutral` | Tertiary actions (Edit, View) |
| `soft` | `error` | Destructive actions (Delete, Remove) |

### Cards

- Use `UCard` with default styling (inherits surface background)
- Add `divide-y divide-zinc-800` for sectioned cards
- Use `shadow-lg` sparingly for elevated modals

### Forms

- Use `UFormField` for consistent label/error placement
- Use `UInput`, `UTextarea`, `USelect` from Nuxt UI
- Apply `size="lg"` for primary inputs

---

## Consequences

### Positive

- **Consistency:** Single source of truth for colors and component styling
- **Accessibility:** Reka UI base ensures WCAG 2.1 AA compliance
- **Developer Experience:** Tailwind Variants provide type-safe component customization
- **Performance:** Tailwind CSS v4 with CSS-first config reduces bundle size
- **Maintainability:** Theme changes propagate automatically via CSS variables

### Negative

- **Opinionated:** Nuxt UI's styling may require overrides for unique designs
- **Learning Curve:** Developers must learn Nuxt UI's component API
- **Dark Mode Focus:** Light mode is secondary, may need additional testing

### Mitigations

- Document component usage patterns in this ADR
- Create Storybook-like examples in `/docs/components/`
- Test light mode toggle before each release

---

## References

- [Nuxt UI Documentation](https://ui.nuxt.com/)
- [Tailwind CSS v4](https://tailwindcss.com/blog/tailwindcss-v4)
- [Reka UI (Accessibility)](https://reka-ui.com/)
- [WCAG 2.1 Contrast Guidelines](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html)

---

*Approved by: Architecture Team*  
*Last Updated: 2026-01-09*
