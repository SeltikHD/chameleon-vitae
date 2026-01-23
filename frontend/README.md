# Chameleon Vitae Frontend

[![Nuxt](https://img.shields.io/badge/Nuxt-4.2.2-00DC82?logo=nuxt&labelColor=020420)](https://nuxt.com)
[![Vue](https://img.shields.io/badge/Vue-3.5.26-42b883?logo=vue.js&labelColor=35495e)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9.3-3178c6?logo=typescript&labelColor=f0f0f0)](https://www.typescriptlang.org)

AI-powered resume engineering system frontend built with Nuxt 4, Vue 3, and Nuxt UI.

## 🚀 Features

- ⚡️ Nuxt 4 with full TypeScript support
- 🎨 Nuxt UI for beautiful, accessible components
- 🔍 Comprehensive SEO Implementation
- 📊 Vercel Analytics & Speed Insights integration
- 🎯 Firebase Authentication ready
- 📱 Responsive design with Tailwind CSS
- 🦎 Chameleon-themed branding (Emerald + Violet gradient)

## Quick Start

### Prerequisites

- [Bun](https://bun.sh) 1.3.5 or higher
- Node.js 18+ (if not using Bun)

### Installation

```bash
# Install dependencies
bun install
```

## Development Server

Start the development server on `http://localhost:3000`:

```bash
bun dev
```

## Production

Build the application for production:

```bash
bun run build
```

Preview the production build:

```bash
bun run preview
```

## SEO

This project includes a comprehensive SEO implementation with:

- ✅ Meta tags (title, description, keywords)
- ✅ Open Graph tags (Facebook, LinkedIn)
- ✅ Twitter Card tags
- ✅ JSON-LD structured data (Organization, SoftwareApplication, FAQPage, etc.)
- ✅ XML sitemap generation
- ✅ robots.txt configuration
- ✅ Favicon set (all sizes)
- ✅ Apple touch icon
- ✅ Web manifest (PWA support)
- ✅ Vercel Analytics integration
- ✅ Vercel Speed Insights integration

### SEO Documentation

- **[SEO_IMPLEMENTATION.md](./SEO_IMPLEMENTATION.md)** - Complete implementation details and results expectations
- **[SEO_CHECKLIST.md](./SEO_CHECKLIST.md)** - Testing checklist and maintenance tasks
- **[SEO_TESTING.md](./SEO_TESTING.md)** - Testing commands and validation tools
- **[SEO_IMAGES.md](./SEO_IMAGES.md)** - Image generation guidelines

### Generate SEO Images

All SEO images are included, but you can regenerate them:

```bash
# Requires ImageMagick or Inkscape
./generate-images.sh
```

This generates:

- favicon.ico, favicon-16x16.png, favicon-32x32.png
- apple-touch-icon.png
- android-chrome-192x192.png, android-chrome-512x512.png
- og-image.png (1200x630 for social media)
- logo.png (512x512)

### Test SEO Implementation

```bash
# View meta tags
curl -s http://localhost:3000 | grep -E '<title>|<meta'

# Check JSON-LD schemas
curl -s http://localhost:3000 | grep -A 50 'application/ld+json'

# Verify sitemap
curl http://localhost:3000/sitemap.xml

# Check robots.txt
curl http://localhost:3000/robots.txt
```

**Online validators:**

- [Google Rich Results Test](https://search.google.com/test/rich-results)
- [Facebook Sharing Debugger](https://developers.facebook.com/tools/debug/)
- [Twitter Card Validator](https://cards-dev.twitter.com/validator)
- [Schema Markup Validator](https://validator.schema.org/)

## Code Quality

```bash
# Lint
bun run lint

# Lint and fix
bun run lint:fix

# Format check
bun run format

# Format and fix
bun run format:fix

# Type check
bun run typecheck
```

## Project Structure

```text
frontend/
├── app/
│   ├── app.config.ts          # App configuration
│   ├── app.vue                # Root component with global SEO
│   ├── assets/                # Styles and assets
│   ├── components/            # Vue components
│   └── pages/                 # File-based routing
│       ├── index.vue          # Homepage (SEO optimized)
│       └── login.vue          # Login page (noindex)
├── plugins/                   # Nuxt plugins
│   ├── vercel-analytics.client.ts
│   └── vercel-speed-insights.client.ts
├── public/                    # Static assets
│   ├── *.png                  # SEO images
│   ├── *.ico                  # Favicons
│   ├── robots.txt             # Search engine directives
│   └── site.webmanifest       # PWA manifest
├── seo/                       # SEO related scripts and configs
│   ├── SEO_CHECKLIST.md       # SEO testing checklist
│   ├── SEO_IMPLEMENTATION.md  # SEO implementation details
│   ├── SEO_IMAGES.md          # SEO image generation guide
│   └── SEO_TESTING.md         # SEO testing instructions
├── nuxt.config.ts             # Nuxt configuration (SEO modules)
└── generate-images.sh         # Image generation script
```

## Design System

- **Colors**:
  - Primary: Emerald (`#10b981`)
  - Secondary: Violet (`#8b5cf6`)
  - Background: Zinc (`#09090b`, `#18181b`)
- **Icons**: Lucide Icons, Simple Icons, Heroicons
- **Typography**: System fonts with `sans-serif` stack

## Related Documentation

- Main project README: [../README.md](../README.md)
- Architecture decisions: [../DECISIONS.md](../DECISIONS.md)
- Backend documentation: [../internal/README.md](../internal/README.md)

## License

See the main project [LICENSE](../LICENSE)

## Deployment

Build the application for production:

```bash
bun build
```

Locally preview production build:

```bash
bun preview
```

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.
