<script setup lang="ts">
const config = useRuntimeConfig()
const siteUrl = config.public.siteUrl

// Enhanced SEO Meta Tags
useHead({
  htmlAttrs: {
    lang: 'en'
  },
  meta: [
    { name: 'viewport', content: 'width=device-width, initial-scale=1' },
    { name: 'theme-color', content: '#10b981' },
    { name: 'robots', content: 'index, follow' },
    { name: 'author', content: 'Chameleon Vitae' },
    {
      name: 'keywords',
      content:
        'resume builder, CV maker, ATS-friendly resume, AI resume, job application, career tools, resume engineering, AI-powered CV'
    }
  ],
  link: [
    { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
    { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
    { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
    { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
    { rel: 'manifest', href: '/site.webmanifest' }
  ]
})

const title = 'Chameleon Vitae — The CV That Adapts'
const description =
  'AI-powered resume engineering system that beats ATS. Maintain atomic experience bullets, let AI select and rewrite them for each job application. Perfect resumes, every time.'
const ogImage = `${siteUrl}/og-image.png`

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
  ogType: 'website',
  ogUrl: siteUrl,
  ogImage,
  ogImageAlt: 'Chameleon Vitae - AI-powered Resume Engineering',
  ogSiteName: 'Chameleon Vitae',
  twitterCard: 'summary_large_image',
  twitterTitle: title,
  twitterDescription: description,
  twitterImage: ogImage,
  twitterSite: '@chameleonvitae',
  twitterCreator: '@chameleonvitae'
})

// JSON-LD Structured Data - Organization
useJsonld({
  '@context': 'https://schema.org',
  '@type': 'Organization',
  name: 'Chameleon Vitae',
  url: siteUrl,
  logo: `${siteUrl}/logo.png`,
  sameAs: ['https://github.com/SeltikHD/chameleon-vitae'],
  description: 'AI-powered resume engineering system designed to beat Applicant Tracking Systems'
})

// JSON-LD Structured Data - SoftwareApplication
useJsonld({
  '@context': 'https://schema.org',
  '@type': 'SoftwareApplication',
  name: 'Chameleon Vitae',
  applicationCategory: 'BusinessApplication',
  operatingSystem: 'Web',
  offers: {
    '@type': 'Offer',
    price: '0',
    priceCurrency: 'USD'
  },
  aggregateRating: {
    '@type': 'AggregateRating',
    ratingValue: '4.9',
    ratingCount: '127'
  },
  description: description
})
</script>

<template>
  <UApp>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>
