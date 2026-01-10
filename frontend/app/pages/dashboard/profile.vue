<template>
  <div class="mx-auto max-w-3xl space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-zinc-100">Profile Settings</h1>
      <p class="mt-1 text-zinc-400">Manage your personal information and resume defaults.</p>
    </div>

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
                    v-if="profile.avatarUrl"
                    :src="profile.avatarUrl"
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
                <UButton
                  color="neutral"
                  variant="solid"
                  size="xs"
                  icon="i-lucide-pencil"
                  class="absolute -bottom-1 -right-1"
                />
              </div>
              <div>
                <p class="font-medium text-zinc-100">Profile Photo</p>
                <p class="text-sm text-zinc-400">JPG, PNG or GIF. Max 2MB.</p>
              </div>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <UFormField
                label="First Name"
                name="firstName"
              >
                <UInput
                  v-model="profile.firstName"
                  placeholder="John"
                />
              </UFormField>

              <UFormField
                label="Last Name"
                name="lastName"
              >
                <UInput
                  v-model="profile.lastName"
                  placeholder="Developer"
                />
              </UFormField>
            </div>

            <UFormField
              label="Email"
              name="email"
            >
              <UInput
                v-model="profile.email"
                type="email"
                placeholder="john@example.com"
                icon="i-lucide-mail"
              />
            </UFormField>

            <UFormField
              label="Phone"
              name="phone"
            >
              <UInput
                v-model="profile.phone"
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
                v-model="profile.location"
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
                  v-model="profile.linkedin"
                  placeholder="linkedin.com/in/username"
                  icon="i-simple-icons-linkedin"
                />
              </UFormField>

              <UFormField
                label="GitHub"
                name="github"
              >
                <UInput
                  v-model="profile.github"
                  placeholder="github.com/username"
                  icon="i-simple-icons-github"
                />
              </UFormField>
            </div>

            <UFormField
              label="Website"
              name="website"
            >
              <UInput
                v-model="profile.website"
                placeholder="https://yoursite.com"
                icon="i-lucide-globe"
              />
            </UFormField>

            <div class="flex justify-end pt-4">
              <UButton
                type="submit"
                color="primary"
                :loading="isSaving"
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
            @submit.prevent="saveResumeDefaults"
          >
            <UFormField
              label="Professional Summary"
              name="summary"
            >
              <UTextarea
                v-model="resumeDefaults.summary"
                :rows="4"
                placeholder="Write a brief professional summary that will be used as the default for your resumes..."
              />
            </UFormField>

            <UFormField
              label="Default Template"
              name="template"
            >
              <USelectMenu
                v-model="resumeDefaults.template"
                :items="templateOptions"
                class="w-full"
              />
            </UFormField>

            <div class="grid gap-4 sm:grid-cols-2">
              <UFormField
                label="Font"
                name="font"
              >
                <USelectMenu
                  v-model="resumeDefaults.font"
                  :items="fontOptions"
                  class="w-full"
                />
              </UFormField>

              <UFormField
                label="Color Accent"
                name="color"
              >
                <USelectMenu
                  v-model="resumeDefaults.accentColor"
                  :items="colorOptions"
                  class="w-full"
                />
              </UFormField>
            </div>

            <div class="space-y-4">
              <h3 class="font-medium text-zinc-100">Section Visibility</h3>
              <div class="grid gap-3 sm:grid-cols-2">
                <UCheckbox
                  v-model="resumeDefaults.showPhoto"
                  label="Show profile photo"
                />
                <UCheckbox
                  v-model="resumeDefaults.showSummary"
                  label="Show summary section"
                />
                <UCheckbox
                  v-model="resumeDefaults.showSkills"
                  label="Show skills section"
                />
                <UCheckbox
                  v-model="resumeDefaults.showEducation"
                  label="Show education section"
                />
              </div>
            </div>

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
          <!-- Connected Accounts -->
          <UCard>
            <template #header>
              <h2 class="font-semibold text-zinc-100">Connected Accounts</h2>
            </template>

            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-800">
                    <UIcon
                      name="i-simple-icons-google"
                      class="h-5 w-5 text-zinc-400"
                    />
                  </div>
                  <div>
                    <p class="font-medium text-zinc-100">Google</p>
                    <p class="text-sm text-zinc-400">john@gmail.com</p>
                  </div>
                </div>
                <UButton
                  color="error"
                  variant="ghost"
                  size="sm"
                >
                  Disconnect
                </UButton>
              </div>

              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-800">
                    <UIcon
                      name="i-simple-icons-github"
                      class="h-5 w-5 text-zinc-400"
                    />
                  </div>
                  <div>
                    <p class="font-medium text-zinc-100">GitHub</p>
                    <p class="text-sm text-zinc-500">Not connected</p>
                  </div>
                </div>
                <UButton
                  color="neutral"
                  variant="outline"
                  size="sm"
                >
                  Connect
                </UButton>
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
                  <p class="font-medium text-zinc-100">Export Data</p>
                  <p class="text-sm text-zinc-400">Download all your data in JSON format.</p>
                </div>
                <UButton
                  color="neutral"
                  variant="outline"
                  size="sm"
                >
                  Export
                </UButton>
              </div>

              <div class="flex items-center justify-between">
                <div>
                  <p class="font-medium text-zinc-100">Delete Account</p>
                  <p class="text-sm text-zinc-400">Permanently delete your account and all data.</p>
                </div>
                <UButton
                  color="error"
                  variant="outline"
                  size="sm"
                >
                  Delete Account
                </UButton>
              </div>
            </div>
          </UCard>
        </div>
      </template>
    </UTabs>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard'
})

const isSaving = ref(false)

const tabs = [
  { label: 'Personal Info', slot: 'personal' as const, icon: 'i-lucide-user' },
  { label: 'Resume Defaults', slot: 'resume' as const, icon: 'i-lucide-file-text' },
  { label: 'Account', slot: 'account' as const, icon: 'i-lucide-settings' }
]

// Mock profile data.
const profile = reactive({
  firstName: 'John',
  lastName: 'Developer',
  email: 'john@example.com',
  phone: '+1 (555) 123-4567',
  location: 'San Francisco, CA',
  linkedin: 'linkedin.com/in/johndeveloper',
  github: 'github.com/johndeveloper',
  website: 'https://johndeveloper.com',
  avatarUrl: ''
})

// Mock resume defaults.
const resumeDefaults = reactive({
  summary:
    'Senior Software Engineer with 7+ years of experience building scalable web applications. Expert in React, TypeScript, and Node.js with a passion for clean code and user experience.',
  template: 'modern',
  font: 'inter',
  accentColor: 'emerald',
  showPhoto: false,
  showSummary: true,
  showSkills: true,
  showEducation: true
})

const templateOptions = [
  { label: 'Modern', value: 'modern' },
  { label: 'Classic', value: 'classic' },
  { label: 'Minimal', value: 'minimal' },
  { label: 'Professional', value: 'professional' }
]

const fontOptions = [
  { label: 'Inter', value: 'inter' },
  { label: 'Roboto', value: 'roboto' },
  { label: 'Open Sans', value: 'opensans' },
  { label: 'Lato', value: 'lato' }
]

const colorOptions = [
  { label: 'Emerald', value: 'emerald' },
  { label: 'Blue', value: 'blue' },
  { label: 'Violet', value: 'violet' },
  { label: 'Neutral', value: 'neutral' }
]

function saveProfile() {
  isSaving.value = true
  // Mock API call.
  setTimeout(() => {
    isSaving.value = false
  }, 1000)
}

function saveResumeDefaults() {
  isSaving.value = true
  // Mock API call.
  setTimeout(() => {
    isSaving.value = false
  }, 1000)
}
</script>
