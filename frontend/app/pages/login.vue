<template>
  <div
    class="flex min-h-screen items-center justify-center bg-linear-to-br from-zinc-950 via-zinc-900 to-zinc-950 px-4"
  >
    <div class="w-full max-w-md">
      <!-- Logo & Header -->
      <div class="mb-8 text-center">
        <div
          class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-linear-to-br from-emerald-500 to-emerald-600 shadow-lg shadow-emerald-500/20"
        >
          <UIcon
            name="i-heroicons-document-duplicate"
            class="h-8 w-8 text-white"
          />
        </div>
        <h1 class="text-2xl font-bold text-white">Chameleon Vitae</h1>
        <p class="mt-2 text-sm text-zinc-400">Sign in to manage your resumes</p>
      </div>

      <!-- Login Card -->
      <UCard class="border border-zinc-800/50 bg-zinc-900/90 p-6 backdrop-blur-sm sm:p-8">
        <!-- Error Alert -->
        <UAlert
          v-if="authError"
          color="error"
          variant="subtle"
          icon="i-heroicons-exclamation-triangle"
          :title="authError"
          class="mb-6"
          :close-button="{
            icon: 'i-heroicons-x-mark-20-solid',
            color: 'error',
            variant: 'link'
          }"
          @close="clearError"
        />

        <form @submit.prevent="handleSubmit">
          <!-- Email Field -->
          <UFormField
            label="Email"
            name="email"
            class="mb-4"
          >
            <UInput
              v-model="form.email"
              type="email"
              placeholder="you@example.com"
              icon="i-heroicons-envelope"
              size="lg"
              class="w-full"
              :disabled="isLoading"
              autocomplete="email"
            />
          </UFormField>

          <!-- Password Field -->
          <UFormField
            label="Password"
            name="password"
            class="mb-6"
          >
            <UInput
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••"
              icon="i-heroicons-lock-closed"
              size="lg"
              class="w-full"
              :disabled="isLoading"
              :ui="{ base: 'autocomplete-current-password' }"
              :input="{ autocomplete: 'current-password' }"
            >
              <template #trailing>
                <UButton
                  :icon="showPassword ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'"
                  variant="ghost"
                  color="neutral"
                  size="xs"
                  @click="showPassword = !showPassword"
                />
              </template>
            </UInput>
          </UFormField>

          <!-- Remember Me & Forgot Password -->
          <div class="mb-6 flex items-center justify-between">
            <UCheckbox
              v-model="form.rememberMe"
              label="Remember me"
              :disabled="isLoading"
            />
            <UButton
              variant="link"
              color="primary"
              size="xs"
              :disabled="isLoading"
            >
              Forgot password?
            </UButton>
          </div>

          <!-- Submit Button -->
          <UButton
            type="submit"
            color="primary"
            size="lg"
            block
            :loading="isLoading"
            :disabled="isLoading || !form.email || !form.password"
          >
            {{ isLoading ? 'Signing in...' : 'Sign in' }}
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
            :loading="socialLoading === 'google'"
            :disabled="isLoading || !!socialLoading"
            @click="handleSocialLogin('google')"
          >
            Google
          </UButton>
          <UButton
            color="neutral"
            variant="outline"
            icon="i-simple-icons-github"
            block
            :loading="socialLoading === 'github'"
            :disabled="isLoading || !!socialLoading"
            @click="handleSocialLogin('github')"
          >
            GitHub
          </UButton>
        </div>

        <template #footer>
          <p class="text-center text-sm text-zinc-400">
            Don't have an account?
            <UButton
              variant="link"
              color="primary"
              size="xs"
              :disabled="isLoading || !!socialLoading"
            >
              Sign up
            </UButton>
          </p>
        </template>
      </UCard>

      <!-- Terms -->
      <p class="mt-6 text-center text-xs text-zinc-500">
        By signing in, you agree to our
        <UButton
          variant="link"
          color="neutral"
          size="xs"
          >Terms of Service</UButton
        >
        and
        <UButton
          variant="link"
          color="neutral"
          size="xs"
          >Privacy Policy</UButton
        >
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

const config = useRuntimeConfig()

// SEO: Login pages should not be indexed by search engines
useHead({
  link: [
    {
      rel: 'canonical',
      href: `${config.public.siteUrl}/login`
    }
  ]
})

useSeoMeta({
  title: 'Login - Chameleon Vitae',
  description: 'Sign in to your Chameleon Vitae account to create AI-powered resumes.',
  robots: 'noindex, nofollow', // Prevent search engines from indexing this page
  ogTitle: 'Login - Chameleon Vitae',
  ogDescription: 'Sign in to your account',
  twitterTitle: 'Login - Chameleon Vitae',
  twitterDescription: 'Sign in to your account'
})

definePageMeta({
  layout: false
})

const authStore = useAuthStore()
const router = useRouter()
const toast = useToast()

const form = reactive({
  email: '',
  password: '',
  rememberMe: false
})

const showPassword = ref(false)
const socialLoading = ref<'google' | 'github' | null>(null)

// Computed state from auth store
const isLoading = computed(() => authStore.loading)
const authError = computed(() => authStore.error)

// Clear error when user starts typing
watch(
  () => [form.email, form.password],
  () => {
    if (authStore.error) {
      authStore.$patch({ error: null })
    }
  }
)

// Redirect if already authenticated
watch(
  () => authStore.isAuthenticated,
  (isAuth) => {
    if (isAuth) {
      router.push('/dashboard')
    }
  },
  { immediate: true }
)

function clearError() {
  authStore.$patch({ error: null })
}

async function handleSubmit() {
  if (!form.email || !form.password) {
    return
  }

  try {
    await authStore.signInWithEmail(form.email, form.password)

    toast.add({
      title: 'Welcome back!',
      description: 'You have successfully signed in.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })

    await router.push('/dashboard')
  } catch {
    // Error is already set in the store
  }
}

async function handleSocialLogin(provider: 'google' | 'github') {
  socialLoading.value = provider

  try {
    if (provider === 'google') {
      await authStore.signInWithGoogle()
    } else {
      await authStore.signInWithGitHub()
    }

    toast.add({
      title: 'Welcome!',
      description: 'You have successfully signed in.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })

    await router.push('/dashboard')
  } catch {
    // Error is already set in the store
  } finally {
    socialLoading.value = null
  }
}
</script>
