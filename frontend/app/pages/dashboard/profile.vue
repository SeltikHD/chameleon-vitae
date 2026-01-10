<template>
  <div class="mx-auto max-w-3xl space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-zinc-100">Profile Settings</h1>
      <p class="mt-1 text-zinc-400">Manage your personal information and resume defaults.</p>
    </div>

    <!-- Loading State -->
    <div
      v-if="loading"
      class="flex items-center justify-center py-12"
    >
      <UIcon
        name="i-lucide-loader-2"
        class="h-8 w-8 animate-spin text-emerald-400"
      />
    </div>

    <!-- Error State -->
    <UCard
      v-else-if="error"
      class="border-red-500/20"
    >
      <div class="flex flex-col items-center gap-4 py-8 text-center">
        <UIcon
          name="i-lucide-alert-circle"
          class="h-12 w-12 text-red-400"
        />
        <div>
          <p class="font-medium text-zinc-100">Failed to load profile</p>
          <p class="mt-1 text-sm text-zinc-400">{{ error }}</p>
        </div>
        <UButton
          color="primary"
          @click="fetchProfile()"
        >
          Try Again
        </UButton>
      </div>
    </UCard>

    <!-- Content -->
    <template v-else-if="profile">
      <!-- Tabs -->
      <UTabs
        :items="tabs"
        class="w-full"
      >
        <template #personal>
          <UCard class="mt-6">
            <form
              class="space-y-6"
              @submit.prevent="saveProfile"
            >
              <!-- Avatar Section -->
              <div class="flex items-center gap-6">
                <div class="relative">
                  <div class="h-20 w-20 overflow-hidden rounded-full bg-zinc-800">
                    <img
                      v-if="profile.picture_url"
                      :src="profile.picture_url"
                      alt="Profile"
                      class="h-full w-full object-cover"
                    />
                    <div
                      v-else
                      class="flex h-full w-full items-center justify-center"
                    >
                      <UIcon
                        name="i-lucide-user"
                        class="h-8 w-8 text-zinc-500"
                      />
                    </div>
                  </div>
                </div>
                <div>
                  <p class="font-medium text-zinc-100">Profile Photo</p>
                  <p class="text-sm text-zinc-400">Synced from your login provider.</p>
                </div>
              </div>

              <UFormField
                label="Full Name"
                name="name"
              >
                <UInput
                  v-model="formData.name"
                  placeholder="John Doe"
                />
              </UFormField>

              <UFormField
                label="Professional Headline"
                name="headline"
              >
                <UInput
                  v-model="formData.headline"
                  placeholder="Senior Software Engineer"
                />
              </UFormField>

              <UFormField
                label="Email"
                name="email"
              >
                <UInput
                  :model-value="profile.email || ''"
                  type="email"
                  placeholder="john@example.com"
                  icon="i-lucide-mail"
                  disabled
                />
                <template #hint>
                  <span class="text-xs text-zinc-500"
                    >Email is managed by your login provider.</span
                  >
                </template>
              </UFormField>

              <UFormField
                label="Phone"
                name="phone"
              >
                <UInput
                  v-model="formData.phone"
                  type="tel"
                  placeholder="+1 (555) 123-4567"
                  icon="i-lucide-phone"
                />
              </UFormField>

              <UFormField
                label="Location"
                name="location"
              >
                <UInput
                  v-model="formData.location"
                  placeholder="San Francisco, CA"
                  icon="i-lucide-map-pin"
                />
              </UFormField>

              <div class="grid gap-4 sm:grid-cols-2">
                <UFormField
                  label="LinkedIn"
                  name="linkedin"
                >
                  <UInput
                    v-model="formData.linkedin_url"
                    placeholder="https://linkedin.com/in/username"
                    icon="i-simple-icons-linkedin"
                  />
                </UFormField>

                <UFormField
                  label="GitHub"
                  name="github"
                >
                  <UInput
                    v-model="formData.github_url"
                    placeholder="https://github.com/username"
                    icon="i-simple-icons-github"
                  />
                </UFormField>
              </div>

              <UFormField
                label="Website"
                name="website"
              >
                <UInput
                  v-model="formData.website"
                  placeholder="https://yoursite.com"
                  icon="i-lucide-globe"
                />
              </UFormField>

              <UFormField
                label="Portfolio"
                name="portfolio"
              >
                <UInput
                  v-model="formData.portfolio_url"
                  placeholder="https://portfolio.yoursite.com"
                  icon="i-lucide-folder"
                />
              </UFormField>

              <div class="flex justify-end gap-3 pt-4">
                <UButton
                  color="neutral"
                  variant="ghost"
                  :disabled="!hasChanges"
                  @click="resetForm"
                >
                  Discard Changes
                </UButton>
                <UButton
                  type="submit"
                  color="primary"
                  :loading="isSaving"
                  :disabled="!hasChanges"
                >
                  Save Changes
                </UButton>
              </div>
            </form>
          </UCard>
        </template>

        <template #resume>
          <UCard class="mt-6">
            <form
              class="space-y-6"
              @submit.prevent="saveSummary"
            >
              <UFormField
                label="Professional Summary"
                name="summary"
              >
                <UTextarea
                  v-model="formData.summary"
                  :rows="6"
                  placeholder="Write a brief professional summary that will be used as the default for your resumes..."
                />
              </UFormField>

              <UFormField
                label="Preferred Language"
                name="preferred_language"
              >
                <USelectMenu
                  v-model="preferredLanguageModel"
                  :items="languageOptions"
                  value-key="value"
                  class="w-full sm:w-64"
                />
              </UFormField>

              <div class="flex justify-end pt-4">
                <UButton
                  type="submit"
                  color="primary"
                  :loading="isSaving"
                >
                  Save Defaults
                </UButton>
              </div>
            </form>
          </UCard>
        </template>

        <template #account>
          <div class="mt-6 space-y-6">
            <!-- Connected Account -->
            <UCard>
              <template #header>
                <h2 class="font-semibold text-zinc-100">Account Information</h2>
              </template>

              <div class="space-y-4">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-3">
                    <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-800">
                      <UIcon
                        :name="authProviderIcon"
                        class="h-5 w-5 text-zinc-400"
                      />
                    </div>
                    <div>
                      <p class="font-medium text-zinc-100">{{ authProviderName }}</p>
                      <p class="text-sm text-zinc-400">{{ profile.email }}</p>
                    </div>
                  </div>
                  <UBadge
                    color="primary"
                    variant="subtle"
                  >
                    Connected
                  </UBadge>
                </div>

                <div class="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
                  <div class="flex items-center gap-2 text-sm text-zinc-400">
                    <UIcon
                      name="i-lucide-info"
                      class="h-4 w-4"
                    />
                    <span>Account created {{ formatDate(profile.created_at) }}</span>
                  </div>
                </div>
              </div>
            </UCard>

            <!-- Danger Zone -->
            <UCard class="border-red-500/20">
              <template #header>
                <h2 class="font-semibold text-red-400">Danger Zone</h2>
              </template>

              <div class="space-y-4">
                <div class="flex items-center justify-between">
                  <div>
                    <p class="font-medium text-zinc-100">Sign Out</p>
                    <p class="text-sm text-zinc-400">Sign out from your account on this device.</p>
                  </div>
                  <UButton
                    color="error"
                    variant="outline"
                    size="sm"
                    @click="handleSignOut"
                  >
                    Sign Out
                  </UButton>
                </div>
              </div>
            </UCard>
          </div>
        </template>
      </UTabs>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { apiFetch } from '~/composables/useApiFetch'
import type { UserResponse, UpdateUserRequest } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const authStore = useAuthStore()
const toast = useToast()

const loading = ref(true)
const error = ref<string | null>(null)
const isSaving = ref(false)
const profile = ref<UserResponse | null>(null)

const tabs = [
  { label: 'Personal Info', slot: 'personal' as const, icon: 'i-lucide-user' },
  { label: 'Resume Defaults', slot: 'resume' as const, icon: 'i-lucide-file-text' },
  { label: 'Account', slot: 'account' as const, icon: 'i-lucide-settings' }
]

const languageOptions = [
  { label: 'English', value: 'en' },
  { label: 'Portuguese (Brazil)', value: 'pt-br' }
]

// Form data for editing.
const formData = reactive<UpdateUserRequest>({
  name: null,
  headline: null,
  summary: null,
  location: null,
  phone: null,
  website: null,
  linkedin_url: null,
  github_url: null,
  portfolio_url: null,
  preferred_language: null
})

// Computed wrapper for USelectMenu - converts null to undefined
const preferredLanguageModel = computed({
  get: () => formData.preferred_language ?? undefined,
  set: (value: string | undefined) => {
    formData.preferred_language = value ?? null
  }
})

// Check if form has changes.
const hasChanges = computed(() => {
  if (!profile.value) return false
  return (
    formData.name !== (profile.value.name ?? '') ||
    formData.headline !== (profile.value.headline ?? '') ||
    formData.location !== (profile.value.location ?? '') ||
    formData.phone !== (profile.value.phone ?? '') ||
    formData.website !== (profile.value.website ?? '') ||
    formData.linkedin_url !== (profile.value.linkedin_url ?? '') ||
    formData.github_url !== (profile.value.github_url ?? '') ||
    formData.portfolio_url !== (profile.value.portfolio_url ?? '')
  )
})

// Auth provider info - use actual Firebase provider from authStore.
const authProviderIcon = computed(() => {
  const provider = authStore.providerName
  switch (provider) {
    case 'Google':
      return 'i-simple-icons-google'
    case 'GitHub':
      return 'i-simple-icons-github'
    default:
      return 'i-lucide-mail'
  }
})

const authProviderName = computed(() => {
  return authStore.providerName
})

// Fetch profile data.
async function fetchProfile() {
  loading.value = true
  error.value = null

  try {
    const data = await apiFetch<UserResponse>('/me')
    profile.value = data
    syncFormData()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load profile'
  } finally {
    loading.value = false
  }
}

// Sync form data from profile.
function syncFormData() {
  if (!profile.value) return
  formData.name = profile.value.name ?? ''
  formData.headline = profile.value.headline ?? ''
  formData.summary = profile.value.summary ?? ''
  formData.location = profile.value.location ?? ''
  formData.phone = profile.value.phone ?? ''
  formData.website = profile.value.website ?? ''
  formData.linkedin_url = profile.value.linkedin_url ?? ''
  formData.github_url = profile.value.github_url ?? ''
  formData.portfolio_url = profile.value.portfolio_url ?? ''
  formData.preferred_language = profile.value.preferred_language ?? 'en'
}

// Reset form to original values.
function resetForm() {
  syncFormData()
}

// Save profile changes.
async function saveProfile() {
  isSaving.value = true

  try {
    const updateData: UpdateUserRequest = {
      name: formData.name || null,
      headline: formData.headline || null,
      location: formData.location || null,
      phone: formData.phone || null,
      website: formData.website || null,
      linkedin_url: formData.linkedin_url || null,
      github_url: formData.github_url || null,
      portfolio_url: formData.portfolio_url || null
    }

    const updated = await apiFetch<UserResponse>('/me', {
      method: 'PATCH',
      body: updateData
    })

    profile.value = updated
    syncFormData()

    toast.add({
      title: 'Profile Updated',
      description: 'Your profile has been saved successfully.',
      color: 'success'
    })
  } catch (e) {
    console.error('Failed to save profile:', e)
  } finally {
    isSaving.value = false
  }
}

// Save summary and language preferences.
async function saveSummary() {
  isSaving.value = true

  try {
    const updateData: UpdateUserRequest = {
      summary: formData.summary || null,
      preferred_language: formData.preferred_language || 'en'
    }

    const updated = await apiFetch<UserResponse>('/me', {
      method: 'PATCH',
      body: updateData
    })

    profile.value = updated
    syncFormData()

    toast.add({
      title: 'Defaults Updated',
      description: 'Your resume defaults have been saved.',
      color: 'success'
    })
  } catch (e) {
    console.error('Failed to save defaults:', e)
  } finally {
    isSaving.value = false
  }
}

// Sign out handler.
async function handleSignOut() {
  try {
    await authStore.signOut()
    await navigateTo('/login')
  } catch (e) {
    console.error('Failed to sign out:', e)
  }
}

// Format date for display.
function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// Fetch on mount.
onMounted(() => {
  fetchProfile()
})
</script>
