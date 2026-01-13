<template>
  <div class="mx-auto max-w-2xl space-y-8">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-2xl font-bold text-zinc-100">Create New Resume</h1>
      <p class="mt-2 text-zinc-400">
        Enter a job posting URL to automatically fetch and generate a tailored resume.
      </p>
    </div>

    <!-- Steps Indicator -->
    <div class="flex items-center justify-center gap-4">
      <div
        v-for="(step, index) in steps"
        :key="step.id"
        class="flex items-center"
      >
        <div
          class="flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium transition-colors"
          :class="getStepClasses(index)"
        >
          <UIcon
            v-if="currentStep > index"
            name="i-lucide-check"
            class="h-4 w-4"
          />
          <span v-else>{{ index + 1 }}</span>
        </div>
        <span
          class="ml-2 hidden text-sm sm:block"
          :class="currentStep >= index ? 'text-zinc-100' : 'text-zinc-500'"
        >
          {{ step.label }}
        </span>
        <UIcon
          v-if="index < steps.length - 1"
          name="i-lucide-chevron-right"
          class="mx-4 h-4 w-4 text-zinc-600"
        />
      </div>
    </div>

    <!-- Step 1: Job Details -->
    <UCard v-if="currentStep === 0">
      <template #header>
        <h2 class="font-semibold text-zinc-100">Job Details</h2>
      </template>

      <form
        class="space-y-4"
        @submit.prevent="handleSubmitJobDetails"
      >
        <!-- Job URL (Required) -->
        <UFormField
          label="Job URL"
          name="jobUrl"
          required
          :error="urlError ?? undefined"
        >
          <div class="flex gap-2">
            <UInput
              v-model="form.jobUrl"
              placeholder="https://linkedin.com/jobs/view/12345"
              icon="i-lucide-link"
              class="flex-1"
              :disabled="isParsing"
              @blur="handleUrlBlur"
            />
            <UButton
              type="button"
              color="neutral"
              variant="outline"
              :loading="isParsing"
              :disabled="!form.jobUrl || isParsing"
              @click="parseJobUrl"
            >
              <UIcon
                v-if="!isParsing"
                name="i-lucide-download"
                class="mr-1 h-4 w-4"
              />
              {{ isParsing ? 'Fetching...' : 'Fetch' }}
            </UButton>
          </div>
          <template #hint>
            <span class="text-xs text-zinc-500">
              Enter a job posting URL and click Fetch to auto-fill the description.
            </span>
          </template>
        </UFormField>

        <!-- Additional Info (Optional, shown after parsing) -->
        <div
          v-if="parsedJobTitle"
          class="rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-4"
        >
          <div class="flex items-center gap-2">
            <UIcon
              name="i-lucide-check-circle"
              class="h-5 w-5 text-emerald-400"
            />
            <span class="font-medium text-emerald-400">Job Parsed Successfully</span>
          </div>
          <p class="mt-1 text-sm text-zinc-300">{{ parsedJobTitle }}</p>
        </div>

        <div class="grid gap-4 sm:grid-cols-3">
          <UFormField
            label="Company Name (Optional)"
            name="companyName"
          >
            <UInput
              v-model="form.companyName"
              placeholder="e.g. Google, Microsoft..."
              icon="i-lucide-building-2"
              :disabled="isParsing"
            />
          </UFormField>

          <UFormField
            label="Job Title (Optional)"
            name="jobTitle"
          >
            <UInput
              v-model="form.jobTitle"
              placeholder="e.g. Senior Software Engineer"
              icon="i-lucide-briefcase"
              :disabled="isParsing"
            />
          </UFormField>

          <UFormField
            label="Target Language"
            name="targetLanguage"
          >
            <USelectMenu
              :v-model="form.targetLanguage"
              :items="[
                { label: 'English', value: 'en' },
                { label: 'Portuguese', value: 'pt-br' },
                { label: 'Spanish', value: 'es' },
                { label: 'French', value: 'fr' },
                { label: 'German', value: 'de' }
              ]"
              :default-value="{ label: 'English', value: 'en' }"
              option-attribute="label"
              value-attribute="value"
              class="w-full"
              :disabled="isParsing"
            />
          </UFormField>
        </div>

        <!-- Job Description (Large, Required) -->
        <UFormField
          label="Job Description"
          name="description"
          required
        >
          <UTextarea
            v-model="form.description"
            :rows="12"
            class="font-mono text-sm"
            placeholder="The job description will appear here after fetching from the URL..."
            :disabled="isParsing"
          />
          <template #hint>
            <span class="text-xs text-zinc-500"> {{ form.description.length }} characters </span>
          </template>
        </UFormField>

        <div class="flex justify-end pt-4">
          <UButton
            type="submit"
            color="primary"
            :disabled="!canGenerate"
          >
            Generate Resume
            <UIcon
              name="i-lucide-sparkles"
              class="ml-2 h-4 w-4"
            />
          </UButton>
        </div>
      </form>
    </UCard>

    <!-- Step 2: Generating -->
    <UCard v-if="currentStep === 1">
      <div class="flex flex-col items-center py-12">
        <div class="relative">
          <div class="h-20 w-20 animate-pulse rounded-full bg-emerald-500/20" />
          <div class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
            <UIcon
              name="i-lucide-file-text"
              class="h-10 w-10 text-emerald-400"
            />
          </div>
        </div>
        <h3 class="mt-6 text-lg font-semibold text-zinc-100">Generating Your Resume</h3>
        <p class="mt-2 text-zinc-400">{{ generationStatusText }}</p>

        <div class="mt-6 w-full max-w-xs">
          <div class="flex justify-between text-sm text-zinc-400">
            <span>Progress</span>
            <span>{{ formattedProgress }}%</span>
          </div>
          <div class="mt-2 h-2 overflow-hidden rounded-full bg-zinc-800">
            <div
              class="h-full rounded-full bg-emerald-500 transition-all duration-300"
              :style="{ width: `${generationProgress}%` }"
            />
          </div>
        </div>

        <!-- Error State -->
        <div
          v-if="generationError"
          class="mt-6 text-center"
        >
          <UAlert
            color="error"
            variant="subtle"
            icon="i-heroicons-exclamation-triangle"
            :title="generationError"
          />
          <UButton
            color="primary"
            class="mt-4"
            @click="generateResume"
          >
            Retry Generation
          </UButton>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { apiFetch } from '~/composables/useApiFetch'
import type {
  ResumeResponse,
  CreateResumeRequest,
  ParseJobURLRequest,
  ParseJobURLResponse
} from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

const router = useRouter()
const toast = useToast()

// State
const currentStep = ref(0)
const isParsing = ref(false)
const isGenerating = ref(false)
const generationProgress = ref(0)
const generationError = ref<string | null>(null)
const urlError = ref<string | null>(null)
const parsedJobTitle = ref<string | null>(null)
const createdResumeId = ref<string | null>(null)

// Steps - simplified to 2 steps (removed template selection)
const steps = [
  { id: 'details', label: 'Job Details' },
  { id: 'generate', label: 'Generate' }
]

// Form data
const form = reactive({
  jobUrl: '',
  description: '',
  companyName: '',
  jobTitle: '',
  targetLanguage: 'en'
})

// Computed
const canGenerate = computed(() => {
  return form.jobUrl.trim().length > 0 && form.description.trim().length > 0
})

const formattedProgress = computed(() => {
  return generationProgress.value.toFixed(2)
})

const generationStatusText = computed(() => {
  if (generationProgress.value < 25) return 'Creating resume draft...'
  if (generationProgress.value < 50) return 'Analyzing job requirements...'
  if (generationProgress.value < 75) return 'Selecting the best experience bullets...'
  if (generationProgress.value < 95) return 'Tailoring content for maximum impact...'
  return 'Finalizing your resume...'
})

// Methods
function getStepClasses(index: number) {
  if (currentStep.value > index) {
    return 'bg-emerald-500 text-zinc-950'
  }
  if (currentStep.value === index) {
    return 'bg-emerald-500/20 text-emerald-400 ring-2 ring-emerald-500'
  }
  return 'bg-zinc-800 text-zinc-500'
}

/**
 * Handle URL input blur - auto-parse if URL is valid.
 */
async function handleUrlBlur() {
  if (form.jobUrl && !form.description) {
    await parseJobUrl()
  }
}

/**
 * Parse job URL using the backend tools endpoint.
 */
async function parseJobUrl() {
  if (!form.jobUrl.trim()) {
    urlError.value = 'Please enter a valid job URL'
    return
  }

  // Basic URL validation
  try {
    new URL(form.jobUrl)
  } catch {
    urlError.value = 'Please enter a valid URL'
    return
  }

  urlError.value = null
  isParsing.value = true
  parsedJobTitle.value = null

  try {
    const request: ParseJobURLRequest = {
      url: form.jobUrl.trim()
    }

    const response = await apiFetch<ParseJobURLResponse>('/tools/parse-job', {
      method: 'POST',
      body: request
    })

    // Fill the form with parsed data
    form.description = response.markdown || ''
    parsedJobTitle.value = response.title || 'Job description fetched'

    // Try to extract company/title from the parsed title
    if (response.title) {
      // Common patterns: "Title at Company", "Title - Company", etc.
      const atMatch = response.title.match(/^(.+?)\s+at\s+(.+)$/i)
      const dashMatch = response.title.match(/^(.+?)\s*[-|]\s*(.+)$/i)

      if (atMatch) {
        form.jobTitle = form.jobTitle || atMatch[1]?.trim() || ''
        form.companyName = form.companyName || atMatch[2]?.trim() || ''
      } else if (dashMatch) {
        form.jobTitle = form.jobTitle || dashMatch[1]?.trim() || ''
        form.companyName = form.companyName || dashMatch[2]?.trim() || ''
      }
    }

    toast.add({
      title: 'Job Fetched',
      description: 'The job description has been loaded successfully.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })
  } catch (error) {
    console.error('[NewResume] Failed to parse job URL:', error)
    urlError.value = error instanceof Error ? error.message : 'Failed to fetch job posting'
    toast.add({
      title: 'Fetch Failed',
      description: 'Could not fetch the job posting. Please paste the description manually.',
      color: 'warning',
      icon: 'i-heroicons-exclamation-triangle'
    })
  } finally {
    isParsing.value = false
  }
}

/**
 * Handle form submission - go directly to generation.
 */
function handleSubmitJobDetails() {
  if (!canGenerate.value) return
  generateResume()
}

/**
 * Generate the resume.
 */
async function generateResume() {
  isGenerating.value = true
  generationError.value = null
  currentStep.value = 1
  generationProgress.value = 0

  try {
    // Start progress simulation
    const progressPromise = simulateProgress()

    // Create the resume - the backend handles creation and initial tailoring
    const request: CreateResumeRequest = {
      job_description: form.description,
      job_title: form.jobTitle || undefined,
      company_name: form.companyName || undefined,
      job_url: form.jobUrl || undefined,
      target_language: form.targetLanguage || 'en'
    }

    const resume = await apiFetch<ResumeResponse>('/resumes', {
      method: 'POST',
      body: request
    })

    createdResumeId.value = resume.id

    // Wait for progress simulation to complete
    await progressPromise

    // If the resume is still in draft status, trigger tailoring
    if (resume.status === 'draft') {
      await apiFetch(`/resumes/${resume.id}/tailor`, {
        method: 'POST',
        body: { max_bullets_per_job: 5 }
      })
    }

    generationProgress.value = 100

    toast.add({
      title: 'Resume Generated!',
      description: 'Your tailored resume is ready to view.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })

    // Navigate to the resume detail page
    await router.push(`/dashboard/resumes/${resume.id}`)
  } catch (error) {
    console.error('[NewResume] Generation failed:', error)
    generationError.value = error instanceof Error ? error.message : 'Failed to generate resume'
    toast.add({
      title: 'Generation Failed',
      description: 'Could not generate the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
    isGenerating.value = false
  }
}

/**
 * Simulate progress for better UX.
 */
async function simulateProgress() {
  return new Promise<void>((resolve) => {
    const interval = setInterval(() => {
      if (generationProgress.value < 85) {
        generationProgress.value += Math.random() * 12
      }
      if (generationProgress.value >= 85) {
        clearInterval(interval)
        resolve()
      }
    }, 400)
  })
}
</script>
