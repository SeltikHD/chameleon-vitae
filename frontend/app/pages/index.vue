<script setup lang="ts">
const config = useRuntimeConfig()
const siteUrl = config.public.siteUrl

definePageMeta({
  layout: 'default'
})

// Enhanced SEO for Homepage
useHead({
  title: 'Chameleon Vitae — AI-Powered Resume Builder That Beats ATS',
  link: [{ rel: 'canonical', href: siteUrl }]
})

useSeoMeta({
  title: 'Chameleon Vitae — AI-Powered Resume Builder That Beats ATS',
  description:
    'Create ATS-friendly resumes with AI. Chameleon Vitae maintains atomic experience bullets and intelligently rewrites them for each job application. Get past automated screening systems.',
  ogTitle: 'Chameleon Vitae — AI Resume Builder That Beats ATS',
  ogDescription:
    'AI-powered resume engineering system that beats ATS. Maintain atomic experience bullets, let AI select and rewrite them for each job application.',
  ogType: 'website',
  ogUrl: siteUrl,
  ogImage: `${siteUrl}/og-image.png`,
  twitterCard: 'summary_large_image'
})

// JSON-LD Structured Data for Homepage
useJsonld({
  '@context': 'https://schema.org',
  '@type': 'WebSite',
  name: 'Chameleon Vitae',
  url: siteUrl,
  potentialAction: {
    '@type': 'SearchAction',
    target: {
      '@type': 'EntryPoint',
      urlTemplate: `${siteUrl}/search?q={search_term_string}`
    }
  }
})

// FAQ JSON-LD
useJsonld({
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'What is Chameleon Vitae?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Chameleon Vitae is an AI-powered resume engineering system that helps you create ATS-friendly resumes by maintaining atomic experience bullets and intelligently rewriting them for each job application.'
      }
    },
    {
      '@type': 'Question',
      name: 'How does Chameleon Vitae beat ATS systems?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Chameleon Vitae uses AI to analyze job descriptions and automatically select and rewrite your experience bullets to match the specific requirements, ensuring maximum ATS match scores.'
      }
    },
    {
      '@type': 'Question',
      name: 'Is Chameleon Vitae free to use?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes, Chameleon Vitae is an open-source project available for free. You can self-host it or use our hosted version.'
      }
    }
  ]
})

const features = [
  {
    icon: 'i-lucide-brain',
    title: 'AI-Powered Tailoring',
    description:
      'Our AI analyzes job descriptions and rewrites your experience bullets to match exactly what recruiters are looking for.',
    color: 'secondary' as const
  },
  {
    icon: 'i-lucide-file-check',
    title: 'ATS-Optimized PDFs',
    description:
      'Generated with Headless Chrome for pixel-perfect, machine-readable documents that pass through any Applicant Tracking System.',
    color: 'primary' as const
  },
  {
    icon: 'i-lucide-layers',
    title: 'Atomic Experience Bullets',
    description:
      'Your experiences are stored as modular bullets, allowing the AI to select and combine them uniquely for each application.',
    color: 'primary' as const
  }
]

const steps = [
  {
    icon: 'i-lucide-upload',
    title: 'Import Your Experiences',
    description:
      'Add your work history, education, and projects. Break them down into achievement-focused bullets.'
  },
  {
    icon: 'i-lucide-clipboard-paste',
    title: 'Paste the Job Description',
    description:
      'Just paste the URL or text of any job posting. Our AI extracts key requirements automatically.'
  },
  {
    icon: 'i-lucide-download',
    title: 'Download Tailored Resume',
    description:
      'Get a professionally formatted PDF with bullets selected and rewritten to maximize your match score.'
  }
]

const stats = [
  { value: '95%', label: 'ATS Pass Rate' },
  { value: '3x', label: 'More Interviews' },
  { value: '<30s', label: 'Generation Time' },
  { value: '100%', label: 'Open Source' }
]
</script>

<template>
  <div>
    <!-- Hero Section -->
    <section class="relative overflow-hidden py-24 md:py-32">
      <!-- Background Gradient -->
      <div class="absolute inset-0 -z-10">
        <div
          class="absolute left-1/2 top-0 h-125 w-200 -translate-x-1/2 rounded-full bg-emerald-500/10 blur-3xl"
        />
        <div
          class="absolute left-1/4 top-1/2 h-100 w-150 -translate-x-1/2 rounded-full bg-violet-500/10 blur-3xl"
        />
      </div>

      <div class="container mx-auto px-4 text-center">
        <!-- Badge -->
        <UBadge
          color="secondary"
          variant="subtle"
          size="lg"
          class="mb-6"
        >
          <UIcon
            name="i-lucide-sparkles"
            class="mr-1 h-4 w-4"
          />
          AI-Powered Resume Engineering
        </UBadge>

        <!-- Headline -->
        <h1
          class="mx-auto max-w-4xl text-4xl font-bold tracking-tight text-zinc-100 md:text-6xl lg:text-7xl"
        >
          The CV That
          <span
            class="bg-linear-to-r from-emerald-400 to-emerald-600 bg-clip-text text-transparent"
          >
            Adapts
          </span>
        </h1>

        <!-- Subheadline -->
        <p class="mx-auto mt-6 max-w-2xl text-lg text-zinc-400 md:text-xl">
          Stop sending the same resume everywhere. <strong>Chameleon Vitae</strong> uses AI to
          tailor your experience bullets to each job description, maximizing your ATS match score.
        </p>

        <!-- CTA Buttons -->
        <div class="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
          <UButton
            to="/login"
            color="primary"
            size="xl"
          >
            Get Started Free
            <UIcon
              name="i-lucide-arrow-right"
              class="ml-1 h-5 w-5"
            />
          </UButton>
          <UButton
            to="https://github.com/SeltikHD/chameleon-vitae"
            target="_blank"
            color="neutral"
            variant="outline"
            size="xl"
            icon="i-simple-icons-github"
          >
            View on GitHub
          </UButton>
        </div>

        <!-- Social Proof + Stats -->
        <div class="mt-16 grid grid-cols-2 gap-8 md:grid-cols-4">
          <div
            v-for="stat in stats"
            :key="stat.label"
          >
            <div class="text-3xl font-bold text-emerald-400 md:text-4xl">
              {{ stat.value }}
            </div>
            <div class="mt-2 text-sm text-zinc-400">
              {{ stat.label }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section
      id="features"
      class="border-t border-zinc-800 py-24"
    >
      <div class="container mx-auto px-4">
        <div class="text-center">
          <h2 class="text-3xl font-bold text-zinc-100 md:text-4xl">Why Chameleon Vitae?</h2>
          <p class="mx-auto mt-4 max-w-2xl text-zinc-400">
            Traditional resumes are static. Ours evolves with every application.
          </p>
        </div>

        <div class="mt-16 grid gap-8 md:grid-cols-3">
          <LandingFeatureCard
            v-for="feature in features"
            :key="feature.title"
            :icon="feature.icon"
            :title="feature.title"
            :description="feature.description"
            :color="feature.color"
          />
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section
      id="how-it-works"
      class="border-t border-zinc-800 bg-zinc-900/50 py-24"
    >
      <div class="container mx-auto px-4">
        <div class="text-center">
          <h2 class="text-3xl font-bold text-zinc-100 md:text-4xl">How It Works</h2>
          <p class="mx-auto mt-4 max-w-2xl text-zinc-400">
            Three simple steps to your perfect resume
          </p>
        </div>

        <div class="mt-16 grid gap-8 md:grid-cols-3">
          <LandingStepCard
            v-for="(step, index) in steps"
            :key="step.title"
            :number="index + 1"
            :title="step.title"
            :description="step.description"
            :icon="step.icon"
          />
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="border-t border-zinc-800 py-24">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl font-bold text-zinc-100 md:text-4xl">Ready to Land Your Dream Job?</h2>
        <p class="mx-auto mt-4 max-w-xl text-zinc-400">
          Join thousands of professionals who've increased their interview rates with tailored
          resumes.
        </p>
        <div class="mt-10">
          <UButton
            to="/login"
            color="primary"
            size="xl"
          >
            Create Your First Resume
            <UIcon
              name="i-lucide-arrow-right"
              class="ml-1 h-5 w-5"
            />
          </UButton>
        </div>
      </div>
    </section>
  </div>
</template>
