<template>
  <div>
    <!-- Hero Section -->
    <section class="relative overflow-hidden py-24 md:py-32">
      <!-- Background Gradient -->
      <div class="absolute inset-0 -z-10">
        <div
          class="absolute left-1/2 top-0 h-[500px] w-[800px] -translate-x-1/2 rounded-full bg-emerald-500/10 blur-3xl"
        />
        <div
          class="absolute left-1/4 top-1/2 h-[400px] w-[600px] -translate-x-1/2 rounded-full bg-violet-500/10 blur-3xl"
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
            class="bg-gradient-to-r from-emerald-400 to-emerald-600 bg-clip-text text-transparent"
          >
            Adapts
          </span>
        </h1>

        <!-- Subheadline -->
        <p class="mx-auto mt-6 max-w-2xl text-lg text-zinc-400 md:text-xl">
          Stop sending the same resume everywhere. Chameleon Vitae uses AI to tailor your experience
          bullets to each job description, maximizing your ATS match score.
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
            to="#how-it-works"
            color="neutral"
            variant="outline"
            size="xl"
          >
            See How It Works
          </UButton>
        </div>

        <!-- Social Proof -->
        <p class="mt-8 text-sm text-zinc-500">
          Trusted by <span class="text-zinc-300">2,500+</span> job seekers worldwide
        </p>
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

<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

// Feature cards data.
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

// How it works steps.
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
</script>
