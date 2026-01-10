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
        @submit.prevent="nextStep"
      >
        <UFormField
          label="Job URL (optional)"
          name="jobUrl"
        >
          <UInput
            v-model="form.jobUrl"
            placeholder="https://careers.example.com/job/123"
            icon="i-lucide-link"
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
          />
        </UFormField>

        <div class="flex justify-end">
          <UButton
            type="submit"
            color="primary"
            :disabled="!form.description"
          >
            Analyze Job
            <UIcon
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

      <!-- Loading State -->
      <div
        v-if="isAnalyzing"
        class="flex flex-col items-center py-12"
      >
        <div class="relative">
          <div
            class="h-16 w-16 animate-spin rounded-full border-4 border-zinc-700 border-t-emerald-500"
          />
          <UIcon
            name="i-lucide-brain"
            class="absolute left-1/2 top-1/2 h-6 w-6 -translate-x-1/2 -translate-y-1/2 text-emerald-400"
          />
        </div>
        <p class="mt-4 text-zinc-400">Analyzing job requirements...</p>
        <p class="text-sm text-zinc-500">This may take a few seconds</p>
      </div>

      <!-- Analysis Results -->
      <div
        v-else
        class="space-y-6"
      >
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
              <p class="font-semibold text-zinc-100">{{ analysis.company }}</p>
              <p class="text-sm text-zinc-400">{{ analysis.position }}</p>
            </div>
          </div>
        </div>

        <div>
          <h3 class="mb-3 text-sm font-medium text-zinc-400">Required Skills</h3>
          <div class="flex flex-wrap gap-2">
            <UBadge
              v-for="skill in analysis.requiredSkills"
              :key="skill"
              color="primary"
              variant="subtle"
            >
              {{ skill }}
            </UBadge>
          </div>
        </div>

        <div>
          <h3 class="mb-3 text-sm font-medium text-zinc-400">Nice to Have</h3>
          <div class="flex flex-wrap gap-2">
            <UBadge
              v-for="skill in analysis.niceToHave"
              :key="skill"
              color="neutral"
              variant="subtle"
            >
              {{ skill }}
            </UBadge>
          </div>
        </div>

        <div>
          <h3 class="mb-3 text-sm font-medium text-zinc-400">Experience Level</h3>
          <UBadge
            color="secondary"
            variant="subtle"
          >
            {{ analysis.experienceLevel }}
          </UBadge>
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
            @click="nextStep"
          >
            Generate Resume
            <UIcon
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
        <p class="mt-2 text-zinc-400">Selecting and tailoring your best experience bullets...</p>

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
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard'
})

const currentStep = ref(0)
const isAnalyzing = ref(false)
const generationProgress = ref(0)

const steps = [
  { id: 'description', label: 'Job Description' },
  { id: 'analysis', label: 'AI Analysis' },
  { id: 'generate', label: 'Generate' }
]

const form = reactive({
  jobUrl: '',
  description: ''
})

// Mock analysis result.
const analysis = reactive({
  company: 'Acme Corp',
  position: 'Senior Frontend Engineer',
  requiredSkills: ['React', 'TypeScript', 'Next.js', 'Tailwind CSS', 'GraphQL'],
  niceToHave: ['Vue.js', 'Testing', 'CI/CD', 'AWS'],
  experienceLevel: '5+ years'
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

function nextStep() {
  if (currentStep.value === 0) {
    // Simulate analysis.
    currentStep.value = 1
    isAnalyzing.value = true
    setTimeout(() => {
      isAnalyzing.value = false
    }, 2000)
  } else if (currentStep.value === 1) {
    // Simulate generation.
    currentStep.value = 2
    simulateGeneration()
  }
}

function simulateGeneration() {
  generationProgress.value = 0
  const interval = setInterval(() => {
    generationProgress.value += 10
    if (generationProgress.value >= 100) {
      clearInterval(interval)
      // Navigate to the resume detail page.
      setTimeout(() => {
        navigateTo('/dashboard/resumes/new-resume-id')
      }, 500)
    }
  }, 300)
}
</script>
