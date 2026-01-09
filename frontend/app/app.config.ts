/**
 * App Configuration for Chameleon Vitae.
 *
 * Theme: "Cyber Chameleon" — Dark Mode First
 *
 * Color Palette:
 * - Primary (Emerald #10b981): Actions, links, buttons, success
 * - Secondary (Violet #8b5cf6): AI features, highlights, premium
 * - Neutral (Zinc): Backgrounds, borders, text
 *
 * @see docs/decisions/ADR-012-frontend-design-system.md
 */
export default defineAppConfig({
  ui: {
    colors: {
      // Primary actions and interactive elements.
      primary: 'emerald',
      // AI-powered features and highlights.
      secondary: 'violet',
      // Backgrounds, borders, and neutral text.
      neutral: 'zinc',
      // Semantic colors for feedback states.
      success: 'emerald',
      warning: 'amber',
      error: 'red',
      info: 'sky'
    }
  }
})
