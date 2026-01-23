// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ['@nuxt/eslint', '@nuxt/ui', '@pinia/nuxt', '@nuxtjs/sitemap', 'nuxt-jsonld'],

  devtools: {
    enabled: true
  },

  css: ['~/assets/css/main.css'],

  // SEO and Site Configuration
  site: {
    url: process.env.NUXT_PUBLIC_SITE_URL || 'http://localhost:3000',
    name: 'Chameleon Vitae',
    description:
      'AI-powered resume engineering system that adapts your CV to any job description. Beat ATS with atomic experience bullets and intelligent rewriting.',
    defaultLocale: 'en'
  },

  nitro: {
    prerender: {
      crawlLinks: true,
      routes: ['/', '/login', 'sitemap.xml']
    }
  },

  // Sitemap Configuration
  sitemap: {
    include: ['/', '/login'],
    exclude: ['/dashboard/**'],
    zeroRuntime: true
  },

  // Font configuration to avoid font decoding warnings.
  // @nuxt/ui includes @nuxt/fonts which handles Google Fonts.
  fonts: {
    families: [
      { name: 'Inter', provider: 'google', weights: [400, 500, 600, 700] },
      { name: 'JetBrains Mono', provider: 'google', weights: [400, 500, 600] }
    ],
    defaults: {
      weights: [400, 500, 600, 700]
    }
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api',
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'https://chameleon-vitae.vercel.app',
      // Firebase client-side configuration (safe to expose).
      firebaseApiKey: process.env.NUXT_PUBLIC_FIREBASE_API_KEY || '',
      firebaseAuthDomain: process.env.NUXT_PUBLIC_FIREBASE_AUTH_DOMAIN || '',
      firebaseProjectId: process.env.NUXT_PUBLIC_FIREBASE_PROJECT_ID || '',
      firebaseStorageBucket: process.env.NUXT_PUBLIC_FIREBASE_STORAGE_BUCKET || '',
      firebaseMessagingSenderId: process.env.NUXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID || '',
      firebaseAppId: process.env.NUXT_PUBLIC_FIREBASE_APP_ID || ''
    }
  },

  compatibilityDate: '2025-01-23',

  // ESLint configuration.
  // Formatting is handled by Prettier (see .prettierrc).
  // ESLint Stylistic is disabled to avoid conflicts with Prettier.
  eslint: {
    config: {
      stylistic: false
    }
  }
})
