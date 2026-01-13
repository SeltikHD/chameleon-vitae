<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-zinc-100">Languages</h1>
        <p class="mt-1 text-zinc-400">Add the languages you speak and your proficiency level.</p>
      </div>
      <UButton
        color="primary"
        @click="openAddModal"
      >
        <UIcon
          name="i-lucide-plus"
          class="mr-2 h-4 w-4"
        />
        Add Language
      </UButton>
    </div>

    <!-- Stats -->
    <div class="grid gap-4 sm:grid-cols-3">
      <DashboardStatsCard
        title="Total Languages"
        :value="String(languageList.length)"
        icon="i-lucide-languages"
        color="primary"
      />
      <DashboardStatsCard
        title="Native Languages"
        :value="String(nativeCount)"
        icon="i-lucide-star"
        color="secondary"
      />
      <DashboardStatsCard
        title="Fluent+"
        :value="String(fluentPlusCount)"
        icon="i-lucide-message-circle"
        color="primary"
      />
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
          <p class="font-medium text-zinc-100">Failed to load languages</p>
          <p class="mt-1 text-sm text-zinc-400">{{ error }}</p>
        </div>
        <UButton
          color="primary"
          @click="fetchLanguages()"
        >
          Try Again
        </UButton>
      </div>
    </UCard>

    <!-- Empty State -->
    <SharedEmptyState
      v-else-if="languageList.length === 0"
      icon="i-lucide-languages"
      title="No languages yet"
      description="Add the languages you speak to enhance your resume."
      action-label="Add Language"
      @action="openAddModal"
    />

    <!-- Language Cards -->
    <div
      v-else
      class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
    >
      <UCard
        v-for="lang in languageList"
        :key="lang.id"
        class="group"
      >
        <div class="flex items-start justify-between">
          <div class="flex items-start gap-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-500/10">
              <UIcon
                :name="getLanguageIcon(lang.proficiency)"
                class="h-6 w-6 text-emerald-400"
              />
            </div>
            <div>
              <h3 class="font-semibold text-zinc-100">{{ lang.language }}</h3>
              <div class="mt-1 flex items-center gap-2">
                <UBadge
                  :color="getProficiencyColor(lang.proficiency)"
                  variant="subtle"
                  size="sm"
                >
                  {{ formatProficiency(lang.proficiency) }}
                </UBadge>
              </div>
            </div>
          </div>
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-trash-2"
            size="sm"
            class="opacity-0 transition-opacity group-hover:opacity-100"
            @click="openDeleteModal(lang)"
          />
        </div>
      </UCard>
    </div>

    <!-- Add Language Modal -->
    <UModal v-model:open="isModalOpen">
      <template #content>
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-zinc-100">Add Language</h3>
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                size="sm"
                @click="closeModal"
              />
            </div>
          </template>

          <form
            class="space-y-4"
            @submit.prevent="saveLanguage"
          >
            <!-- Language Name -->
            <UFormField
              label="Language"
              required
            >
              <UInput
                v-model="form.language"
                placeholder="e.g., English, Portuguese, Spanish"
                :color="formErrors.language ? 'error' : undefined"
              />
              <template
                v-if="formErrors.language"
                #error
              >
                {{ formErrors.language }}
              </template>
            </UFormField>

            <!-- Proficiency Level -->
            <UFormField
              label="Proficiency Level"
              required
            >
              <USelectMenu
                v-model="form.proficiency"
                :items="proficiencyOptions"
                value-key="value"
                placeholder="Select proficiency level"
              />
              <template
                v-if="formErrors.proficiency"
                #error
              >
                {{ formErrors.proficiency }}
              </template>
            </UFormField>
          </form>

          <template #footer>
            <div class="flex justify-end gap-3">
              <UButton
                color="neutral"
                variant="outline"
                @click="closeModal"
              >
                Cancel
              </UButton>
              <UButton
                color="primary"
                :loading="saving"
                @click="saveLanguage"
              >
                Add Language
              </UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Delete Confirmation Modal -->
    <UModal v-model:open="isDeleteModalOpen">
      <template #content>
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-zinc-100">Delete Language</h3>
          </template>

          <p class="text-zinc-400">
            Are you sure you want to remove
            <strong class="text-zinc-100">{{ deleteTarget?.language }}</strong>
            from your profile? This action cannot be undone.
          </p>

          <template #footer>
            <div class="flex justify-end gap-3">
              <UButton
                color="neutral"
                variant="outline"
                @click="isDeleteModalOpen = false"
              >
                Cancel
              </UButton>
              <UButton
                color="error"
                :loading="deleting"
                @click="confirmDelete"
              >
                Delete
              </UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { apiFetch } from '~/composables/useApiFetch'
import type {
  SpokenLanguageResponse,
  CreateSpokenLanguageRequest,
  ListSpokenLanguagesResponse,
  LanguageProficiency
} from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

// State
const languageList = ref<SpokenLanguageResponse[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const saving = ref(false)
const deleting = ref(false)

// Modal state
const isModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const deleteTarget = ref<SpokenLanguageResponse | null>(null)

// Form state
const form = ref<{ language: string; proficiency: LanguageProficiency }>({
  language: '',
  proficiency: 'intermediate'
})
const formErrors = ref<Record<string, string>>({})

// Composables
const toast = useToast()

// Proficiency options
const proficiencyOptions = [
  { label: 'Native / Bilingual', value: 'native' },
  { label: 'Fluent', value: 'fluent' },
  { label: 'Advanced', value: 'advanced' },
  { label: 'Intermediate', value: 'intermediate' },
  { label: 'Basic', value: 'basic' }
]

// Computed
const nativeCount = computed(() => {
  return languageList.value.filter((l) => l.proficiency === 'native').length
})

const fluentPlusCount = computed(() => {
  return languageList.value.filter((l) => ['native', 'fluent', 'advanced'].includes(l.proficiency))
    .length
})

// Methods
async function fetchLanguages() {
  loading.value = true
  error.value = null

  try {
    const response = await apiFetch<ListSpokenLanguagesResponse>('/languages')
    languageList.value = response.data || []
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load languages'
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  resetForm()
  isModalOpen.value = true
}

function closeModal() {
  isModalOpen.value = false
  resetForm()
}

function resetForm() {
  form.value = {
    language: '',
    proficiency: 'intermediate'
  }
  formErrors.value = {}
}

function validateForm(): boolean {
  formErrors.value = {}

  if (!form.value.language.trim()) {
    formErrors.value.language = 'Language name is required'
  }

  if (!form.value.proficiency) {
    formErrors.value.proficiency = 'Proficiency level is required'
  }

  return Object.keys(formErrors.value).length === 0
}

async function saveLanguage() {
  if (!validateForm()) return

  saving.value = true

  try {
    const payload: CreateSpokenLanguageRequest = {
      language: form.value.language.trim(),
      proficiency: form.value.proficiency
    }

    await apiFetch('/languages', {
      method: 'POST',
      body: payload
    })

    toast.add({
      title: 'Language Added',
      description: `${form.value.language} has been added to your profile.`,
      color: 'success'
    })

    closeModal()
    await fetchLanguages()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to add language',
      color: 'error'
    })
  } finally {
    saving.value = false
  }
}

function openDeleteModal(lang: SpokenLanguageResponse) {
  deleteTarget.value = lang
  isDeleteModalOpen.value = true
}

async function confirmDelete() {
  if (!deleteTarget.value) return

  deleting.value = true

  try {
    await apiFetch(`/languages/${deleteTarget.value.id}`, {
      method: 'DELETE'
    })

    toast.add({
      title: 'Language Deleted',
      description: `${deleteTarget.value.language} has been removed.`,
      color: 'success'
    })

    isDeleteModalOpen.value = false
    deleteTarget.value = null
    await fetchLanguages()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to delete language',
      color: 'error'
    })
  } finally {
    deleting.value = false
  }
}

function getLanguageIcon(proficiency: string): string {
  switch (proficiency) {
    case 'native':
      return 'i-lucide-star'
    case 'fluent':
      return 'i-lucide-message-circle'
    case 'advanced':
      return 'i-lucide-trending-up'
    case 'intermediate':
      return 'i-lucide-book-open'
    default:
      return 'i-lucide-book'
  }
}

function getProficiencyColor(
  proficiency: string
): 'primary' | 'secondary' | 'warning' | 'info' | 'neutral' {
  switch (proficiency) {
    case 'native':
      return 'primary'
    case 'fluent':
      return 'secondary'
    case 'advanced':
      return 'info'
    case 'intermediate':
      return 'warning'
    default:
      return 'neutral'
  }
}

function formatProficiency(proficiency: string): string {
  const option = proficiencyOptions.find((o) => o.value === proficiency)
  return option?.label || proficiency
}

// Lifecycle
onMounted(() => {
  fetchLanguages()
})
</script>
