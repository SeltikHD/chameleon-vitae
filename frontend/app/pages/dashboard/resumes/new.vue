<template>
  <div class="mx-auto max-w-2xl space-y-8">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-2xl font-bold text-zinc-100">Create New Resume</h1>
      <p class="mt-2 text-zinc-400">Paste a job description to generate a tailored resume.</p>
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

    <!-- Step 1: Job Description -->
    <UCard v-if="currentStep === 0">
      <template #header>
        <h2 class="font-semibold text-zinc-100">Job Description</h2>
      </template>

      <form
        class="space-y-4"
        @submit.prevent="analyzeJob"
      >
        <UFormField
          label="Job URL (optional)"
          name="jobUrl"
        >
          <UInput
            v-model="form.jobUrl"
            placeholder="https://careers.example.com/job/123"
            icon="i-lucide-link"
            :disabled="isAnalyzing"
          />
        </UFormField>

        <div class="flex items-center gap-4">
          <div class="h-px flex-1 bg-zinc-800" />
          <span class="text-xs text-zinc-500">or paste the description</span>
          <div class="h-px flex-1 bg-zinc-800" />
        </div>

        <UFormField
          label="Job Description"
          name="description"
          required
        >
          <UTextarea
            v-model="form.description"
            :rows="12"
            placeholder="Paste the full job description here..."
            :disabled="isAnalyzing"
          />
        </UFormField>

        <div class="flex justify-end">
          <UButton
            type="submit"
            color="primary"
            :disabled="!form.description || isAnalyzing"
            :loading="isAnalyzing"
          >
            {{ isAnalyzing ? 'Analyzing...' : 'Analyze Job' }}
            <UIcon
              v-if="!isAnalyzing"
              name="i-lucide-arrow-right"
              class="ml-2 h-4 w-4"
            />
          </UButton>
        </div>
      </form>
    </UCard>

    <!-- Step 2: AI Analysis -->
    <UCard v-if="currentStep === 1">
      <template #header>
        <h2 class="font-semibold text-zinc-100">AI Analysis</h2>
      </template>

      <!-- Analysis Results -->
      <div class="space-y-6">
        <div>
          <h3 class="mb-3 text-sm font-medium text-zinc-400">Detected Company & Role</h3>
          <div class="flex items-center gap-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-zinc-800">
              <UIcon
                name="i-lucide-building-2"
                class="h-6 w-6 text-emerald-400"
              />
            </div>
            <div>
              <UInput
                v-model="analysis.company"
                placeholder="Company name"
                class="mb-1"
              />
              <UInput
                v-model="analysis.position"
                placeholder="Position title"
              />
            </div>
          </div>
        </div>

        <div>
          <h3 class="mb-3 text-sm font-medium text-zinc-400">Detected Requirements</h3>
          <div class="flex flex-wrap gap-2">
            <UBadge
              v-for="skill in analysis.requiredSkills"
              :key="skill"
              color="primary"
              variant="subtle"
            >
              {{ skill }}
            </UBadge>
            <UBadge
              v-if="analysis.requiredSkills.length === 0"
              color="neutral"
              variant="subtle"
            >
              No specific skills detected
            </UBadge>
          </div>
        </div>

        <div>
          <h3 class="mb-3 text-sm font-medium text-zinc-400">Template Selection</h3>
          <USelectMenu
            v-model="selectedTemplateId"
            :items="templateOptions"
            value-key="value"
            placeholder="Select a template"
            class="w-full"
          />
        </div>

        <div class="flex justify-between pt-4">
          <UButton
            color="neutral"
            variant="ghost"
            @click="currentStep = 0"
          >
            <UIcon
              name="i-lucide-arrow-left"
              class="mr-2 h-4 w-4"
            />
            Back
          </UButton>
          <UButton
            color="primary"
            :disabled="!selectedTemplateId || isGenerating"
            :loading="isGenerating"
            @click="generateResume"
          >
            {{ isGenerating ? 'Generating...' : 'Generate Resume' }}
            <UIcon
              v-if="!isGenerating"
              name="i-lucide-sparkles"
              class="ml-2 h-4 w-4"
            />
          </UButton>
        </div>
      </div>
    </UCard>

    <!-- Step 3: Generating -->
    <UCard v-if="currentStep === 2">
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
            <span>{{ generationProgress }}%</span>
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
import type { ResumeResponse, AnalyzeJobRequest, AnalyzeJobResponse } from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

const router = useRouter()
const toast = useToast()

const currentStep = ref(0)
const isAnalyzing = ref(false)
const isGenerating = ref(false)
const generationProgress = ref(0)
const generationError = ref<string | null>(null)
const selectedTemplateId = ref('')
const createdResumeId = ref<string | null>(null)

const steps = [
  { id: 'description', label: 'Job Description' },
  { id: 'analysis', label: 'AI Analysis' },
  { id: 'generate', label: 'Generate' }
]

const form = reactive({
  jobUrl: '',
  description: ''
})

const analysis = reactive({
  company: '',
  position: '',
  requiredSkills: [] as string[],
  niceToHave: [] as string[],
  experienceLevel: ''
})

// Template options (these would come from API in production)
const templateOptions = ref([
  { label: 'Modern Professional', value: 'template-modern' },
  { label: 'Classic Traditional', value: 'template-classic' },
  { label: 'Minimalist', value: 'template-minimal' },
  { label: 'Creative', value: 'template-creative' }
])

const generationStatusText = computed(() => {
  if (generationProgress.value < 30) return 'Analyzing your experience bullets...'
  if (generationProgress.value < 60) return 'Selecting the best matches for this job...'
  if (generationProgress.value < 90) return 'Tailoring content for maximum impact...'
  return 'Finalizing your resume...'
})

function getStepClasses(index: number) {
  if (currentStep.value > index) {
    return 'bg-emerald-500 text-zinc-950'
  }
  if (currentStep.value === index) {
    return 'bg-emerald-500/20 text-emerald-400 ring-2 ring-emerald-500'
  }
  return 'bg-zinc-800 text-zinc-500'
}

async function analyzeJob() {
  if (!form.description.trim()) return

  isAnalyzing.value = true

  try {
    const request: AnalyzeJobRequest = {
      job_description: form.description,
      job_url: form.jobUrl || undefined
    }

    const response = await apiFetch<AnalyzeJobResponse>('/resumes/analyze', {
      method: 'POST',
      body: request
    })

    // Update analysis with response
    analysis.company = response.company || ''
    analysis.position = response.position || ''
    analysis.requiredSkills = response.required_skills || []
    analysis.niceToHave = response.nice_to_have || []
    analysis.experienceLevel = response.experience_level || ''

    currentStep.value = 1
  } catch (error) {
    console.error('[NewResume] Analysis failed:', error)
    toast.add({
      title: 'Analysis Failed',
      description: 'Could not analyze the job description. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isAnalyzing.value = false
  }
}

async function generateResume() {
  if (!selectedTemplateId.value) return

  isGenerating.value = true
  generationError.value = null
  currentStep.value = 2
  generationProgress.value = 0

  try {
    // Create the resume
    const resume = await apiFetch<ResumeResponse>('/resumes', {
      method: 'POST',
      body: {
        job_description: form.description,
        job_title: analysis.position || undefined,
        company_name: analysis.company || undefined
      }
    })

    createdResumeId.value = resume.id

    // Simulate progress while the backend generates
    await simulateProgress()

    // Trigger generation
    await apiFetch(`/resumes/${resume.id}/generate`, {
      method: 'POST'
    })

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

async function simulateProgress() {
  return new Promise<void>((resolve) => {
    const interval = setInterval(() => {
      if (generationProgress.value < 85) {
        generationProgress.value += Math.random() * 15
      }
      if (generationProgress.value >= 85) {
        clearInterval(interval)
        resolve()
      }
    }, 500)
  })
}
</script>
