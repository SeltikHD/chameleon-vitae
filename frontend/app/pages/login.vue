<template>
  <div class="flex min-h-screen items-center justify-center py-12">
    <!-- Background Gradient -->
    <div class="absolute inset-0 -z-10">
      <div class="absolute right-1/4 top-1/4 h-[400px] w-[400px] rounded-full bg-emerald-500/5 blur-3xl" />
      <div class="absolute bottom-1/4 left-1/4 h-[300px] w-[300px] rounded-full bg-violet-500/5 blur-3xl" />
    </div>

    <div class="w-full max-w-md px-4">
      <!-- Logo -->
      <div class="mb-8 text-center">
        <NuxtLink to="/" class="inline-flex items-center gap-2">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500">
            <UIcon name="i-lucide-file-text" class="h-6 w-6 text-zinc-950" />
          </div>
          <span class="text-xl font-bold text-zinc-100">Chameleon Vitae</span>
        </NuxtLink>
      </div>

      <!-- Login Card -->
      <UCard>
        <template #header>
          <div class="text-center">
            <h1 class="text-2xl font-bold text-zinc-100">
              Welcome Back
            </h1>
            <p class="mt-2 text-sm text-zinc-400">
              Sign in to continue building your tailored resumes
            </p>
          </div>
        </template>

        <form class="space-y-4" @submit.prevent="handleSubmit">
          <!-- Email Input -->
          <UFormField label="Email" name="email">
            <UInput
              v-model="form.email"
              type="email"
              placeholder="you@example.com"
              icon="i-lucide-mail"
              class="w-full"
            />
          </UFormField>

          <!-- Password Input -->
          <UFormField label="Password" name="password">
            <UInput
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Enter your password"
              icon="i-lucide-lock"
              class="w-full"
              :ui="{ trailing: 'pe-1' }"
            >
              <template #trailing>
                <UButton
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                  @click="showPassword = !showPassword"
                />
              </template>
            </UInput>
          </UFormField>

          <!-- Remember Me & Forgot Password -->
          <div class="flex items-center justify-between">
            <UCheckbox v-model="form.rememberMe" label="Remember me" />
            <UButton variant="link" color="primary" size="xs">
              Forgot password?
            </UButton>
          </div>

          <!-- Submit Button -->
          <UButton type="submit" color="primary" size="lg" block :loading="isLoading">
            Sign In
          </UButton>
        </form>

        <!-- Divider -->
        <div class="my-6 flex items-center gap-4">
          <div class="h-px flex-1 bg-zinc-800" />
          <span class="text-xs text-zinc-500">or continue with</span>
          <div class="h-px flex-1 bg-zinc-800" />
        </div>

        <!-- Social Login Buttons -->
        <div class="grid grid-cols-2 gap-4">
          <UButton
            color="neutral"
            variant="outline"
            icon="i-simple-icons-google"
            block
            @click="handleSocialLogin('google')"
          >
            Google
          </UButton>
          <UButton
            color="neutral"
            variant="outline"
            icon="i-simple-icons-github"
            block
            @click="handleSocialLogin('github')"
          >
            GitHub
          </UButton>
        </div>

        <template #footer>
          <p class="text-center text-sm text-zinc-400">
            Don't have an account?
            <UButton variant="link" color="primary" size="xs">
              Sign up
            </UButton>
          </p>
        </template>
      </UCard>

      <!-- Terms -->
      <p class="mt-6 text-center text-xs text-zinc-500">
        By signing in, you agree to our
        <UButton variant="link" color="neutral" size="xs">Terms of Service</UButton>
        and
        <UButton variant="link" color="neutral" size="xs">Privacy Policy</UButton>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false
})

const form = reactive({
  email: '',
  password: '',
  rememberMe: false
})

const showPassword = ref(false)
const isLoading = ref(false)

function handleSubmit() {
  isLoading.value = true

  // Mock login - simulate API call.
  setTimeout(() => {
    isLoading.value = false
    navigateTo('/dashboard')
  }, 1500)
}

function handleSocialLogin(provider: 'google' | 'github') {
  // Mock social login.
  navigateTo('/dashboard')
}
</script>
