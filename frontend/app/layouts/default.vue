<template>
  <div class="min-h-screen bg-zinc-950">
    <SpeedInsights />
    <!-- Public Header -->
    <header
      class="fixed left-0 right-0 top-0 z-50 border-b border-zinc-800 bg-zinc-950/80 backdrop-blur-sm"
    >
      <div class="container mx-auto flex h-16 items-center justify-between px-4">
        <!-- Logo -->
        <SharedLogo
          to="/"
          size="md"
        />

        <!-- Navigation -->
        <nav class="hidden items-center gap-6 md:flex">
          <NuxtLink
            to="/#features"
            class="text-sm text-zinc-400 transition hover:text-zinc-100"
          >
            Features
          </NuxtLink>
          <NuxtLink
            to="/#how-it-works"
            class="text-sm text-zinc-400 transition hover:text-zinc-100"
          >
            How It Works
          </NuxtLink>
        </nav>

        <!-- Auth Buttons -->
        <div class="flex items-center gap-4">
          <UButton
            to="/login"
            variant="ghost"
            color="neutral"
            size="sm"
          >
            Sign In
          </UButton>
          <UButton
            to="/login"
            color="primary"
            size="sm"
          >
            Get Started
          </UButton>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="pt-16">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="border-t border-zinc-800 bg-zinc-950 py-8">
      <div class="container mx-auto px-4">
        <div class="flex flex-col items-center justify-between gap-4 md:flex-row">
          <SharedLogo
            to="/"
            size="sm"
          />
          <p class="text-sm text-zinc-500">
            © {{ new Date().getFullYear() }} Chameleon Vitae. Built with ☕ and a healthy disdain
            for HR forms.
          </p>
          <div class="flex items-center gap-4">
            <UButton
              to="https://github.com/SeltikHD/chameleon-vitae"
              target="_blank"
              icon="i-simple-icons-github"
              aria-label="GitHub"
              color="neutral"
              variant="ghost"
              size="sm"
            />
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { SpeedInsights } from '@vercel/speed-insights/vue'
</script>
