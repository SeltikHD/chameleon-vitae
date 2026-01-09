<template>
  <div class="min-h-screen bg-zinc-950">
    <!-- Public Header -->
    <header class="fixed top-0 left-0 right-0 z-50 border-b border-zinc-800 bg-zinc-950/80 backdrop-blur-sm">
      <div class="container mx-auto flex h-16 items-center justify-between px-4">
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center gap-2">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500">
            <UIcon name="i-lucide-file-text" class="h-5 w-5 text-zinc-950" />
          </div>
          <span class="text-lg font-bold text-zinc-100">Chameleon Vitae</span>
        </NuxtLink>

        <!-- Navigation -->
        <nav class="hidden items-center gap-6 md:flex">
          <NuxtLink to="/#features" class="text-sm text-zinc-400 transition hover:text-zinc-100">
            Features
          </NuxtLink>
          <NuxtLink to="/#how-it-works" class="text-sm text-zinc-400 transition hover:text-zinc-100">
            How It Works
          </NuxtLink>
        </nav>

        <!-- Auth Button -->
        <div class="flex items-center gap-4">
          <UButton to="/login" variant="ghost" color="neutral" size="sm">
            Sign In
          </UButton>
          <UButton to="/login" color="primary" size="sm">
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
      <div class="container mx-auto px-4 text-center">
        <p class="text-sm text-zinc-500">
          © 2026 Chameleon Vitae. Built with ☕ and a healthy disdain for HR forms.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
// Default layout for public pages.
</script>
