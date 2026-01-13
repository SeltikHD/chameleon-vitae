<script setup lang="ts">
import { apiFetch } from '~/composables/useApiFetch'
import type { ResumeResponse, TailorResumeRequest, UpdateResumeContentRequest } from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

const route = useRoute()
const router = useRouter()
const toast = useToast()

// State
const resume = ref<ResumeResponse | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)
const isDownloading = ref(false)
const isRegenerating = ref(false)
const isRetailoring = ref(false)
const isUpdatingStatus = ref(false)
const showDeleteModal = ref(false)
const isDeleting = ref(false)
const notes = ref('')

const resumeId = computed(() => route.params.id as string)

// Computed
const statusColor = computed(() => {
  switch (resume.value?.status) {
    case 'draft':
      return 'warning'
    case 'reviewed':
      return 'success'
    default:
      return 'neutral'
  }
})

const formattedDate = computed(() => {
  if (!resume.value?.created_at) return ''
  return new Date(resume.value.created_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})

const formattedUpdatedDate = computed(() => {
  if (!resume.value?.updated_at) return ''
  return new Date(resume.value.updated_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})

// Lifecycle
onMounted(async () => {
  await fetchResume()
})

// Methods
async function fetchResume() {
  isLoading.value = true
  error.value = null

  try {
    const response = await apiFetch<ResumeResponse>(`/resumes/${resumeId.value}`)
    resume.value = response
    notes.value = response.notes || ''
  } catch (err) {
    console.error('[ResumeDetail] Failed to fetch resume:', err)
    error.value = err instanceof Error ? err.message : 'Failed to load resume'
  } finally {
    isLoading.value = false
  }
}

async function retailorResume() {
  if (!resume.value) return

  isRetailoring.value = true
  try {
    const request: TailorResumeRequest = {
      max_bullets_per_job: 5
    }

    const updated = await apiFetch<ResumeResponse>(`/resumes/${resumeId.value}/tailor`, {
      method: 'POST',
      body: request
    })

    resume.value = updated

    toast.add({
      title: 'Resume Re-tailored',
      description: 'Your resume has been regenerated with fresh AI-optimized content.',
      color: 'success',
      icon: 'i-heroicons-sparkles'
    })
  } catch (err) {
    console.error('[ResumeDetail] Failed to retailor:', err)
    toast.add({
      title: 'Tailoring Failed',
      description: 'Could not regenerate the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isRetailoring.value = false
  }
}

async function updateStatus(newStatus: 'reviewed') {
  if (!resume.value) return

  isUpdatingStatus.value = true
  try {
    const request: UpdateResumeContentRequest = {
      status: newStatus,
      notes: notes.value || undefined
    }

    const updated = await apiFetch<ResumeResponse>(`/resumes/${resumeId.value}/content`, {
      method: 'PATCH',
      body: request
    })

    resume.value = updated

    toast.add({
      title: 'Status Updated',
      description: `Resume marked as ${newStatus}.`,
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })
  } catch (err) {
    console.error('[ResumeDetail] Failed to update status:', err)
    toast.add({
      title: 'Update Failed',
      description: 'Could not update the resume status.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isUpdatingStatus.value = false
  }
}

async function downloadPDF() {
  if (!resume.value) return

  isDownloading.value = true
  try {
    // Fetch PDF - the endpoint returns a blob
    const response = await apiFetch<Blob>(`/resumes/${resumeId.value}/pdf`, {
      responseType: 'blob'
    })

    // Create download link
    const url = globalThis.URL.createObjectURL(response)
    const link = document.createElement('a')
    link.href = url

    // Generate filename
    const jobTitle = resume.value.job_title || 'resume'
    const company = resume.value.company_name || ''
    const filename = company
      ? `${jobTitle.replaceAll(/\s+/g, '-')}_${company.replaceAll(/\s+/g, '-')}.pdf`
      : `${jobTitle.replaceAll(/\s+/g, '-')}.pdf`

    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    globalThis.URL.revokeObjectURL(url)

    toast.add({
      title: 'Downloaded',
      description: 'Your resume PDF has been downloaded.',
      color: 'success',
      icon: 'i-heroicons-arrow-down-tray'
    })
  } catch (err) {
    console.error('[ResumeDetail] Failed to download PDF:', err)
    toast.add({
      title: 'Download Failed',
      description: 'Could not generate PDF. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isDownloading.value = false
  }
}

async function regeneratePDF() {
  if (!resume.value) return

  isRegenerating.value = true
  try {
    // Force regenerate PDF by ignoring cache.
    const response = await apiFetch<Blob>(`/resumes/${resumeId.value}/pdf?force_regenerate=true`, {
      responseType: 'blob'
    })

    // Create download link.
    const url = globalThis.URL.createObjectURL(response)
    const link = document.createElement('a')
    link.href = url

    // Generate filename.
    const jobTitle = resume.value.job_title || 'resume'
    const company = resume.value.company_name || ''
    const filename = company
      ? `${jobTitle.replaceAll(/\s+/g, '-')}_${company.replaceAll(/\s+/g, '-')}.pdf`
      : `${jobTitle.replaceAll(/\s+/g, '-')}.pdf`

    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    globalThis.URL.revokeObjectURL(url)

    toast.add({
      title: 'PDF Regenerated',
      description: 'A fresh PDF has been generated and downloaded.',
      color: 'success',
      icon: 'i-heroicons-arrow-path'
    })
  } catch (err) {
    console.error('[ResumeDetail] Failed to regenerate PDF:', err)
    toast.add({
      title: 'Regeneration Failed',
      description: 'Could not regenerate PDF. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isRegenerating.value = false
  }
}

async function deleteResume() {
  if (!resume.value) return

  isDeleting.value = true
  try {
    await apiFetch(`/resumes/${resumeId.value}`, {
      method: 'DELETE'
    })

    toast.add({
      title: 'Resume Deleted',
      description: 'The resume has been permanently deleted.',
      color: 'success',
      icon: 'i-heroicons-trash'
    })

    await router.push('/dashboard/resumes')
  } catch (err) {
    console.error('[ResumeDetail] Failed to delete resume:', err)
    toast.add({
      title: 'Delete Failed',
      description: 'Could not delete the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isDeleting.value = false
    showDeleteModal.value = false
  }
}

function formatStatus(status: string): string {
  return status.charAt(0).toUpperCase() + status.slice(1)
}
</script>

<template>
  <div class="min-h-screen p-6 lg:p-8">
    <!-- Loading State -->
    <div
      v-if="isLoading"
      class="flex items-center justify-center min-h-100"
    >
      <div class="text-center">
        <UIcon
          name="i-heroicons-arrow-path"
          class="w-8 h-8 text-emerald-400 animate-spin mb-4"
        />
        <p class="text-zinc-400">Loading resume...</p>
      </div>
    </div>

    <!-- Error State -->
    <div
      v-else-if="error"
      class="max-w-2xl mx-auto text-center py-16"
    >
      <UIcon
        name="i-heroicons-exclamation-triangle"
        class="w-12 h-12 text-red-400 mx-auto mb-4"
      />
      <h2 class="text-xl font-semibold text-zinc-100 mb-2">Failed to Load Resume</h2>
      <p class="text-zinc-400 mb-6">{{ error }}</p>
      <div class="flex justify-center gap-4">
        <UButton
          variant="solid"
          color="primary"
          @click="fetchResume"
        >
          Retry
        </UButton>
        <UButton
          variant="ghost"
          color="neutral"
          to="/dashboard/resumes"
        >
          Back to Resumes
        </UButton>
      </div>
    </div>

    <!-- Resume Content -->
    <div
      v-else-if="resume"
      class="max-w-6xl mx-auto"
    >
      <!-- Header -->
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-8">
        <div class="flex items-start gap-4">
          <!-- Back Button -->
          <UButton
            variant="ghost"
            color="neutral"
            size="sm"
            icon="i-heroicons-arrow-left"
            to="/dashboard/resumes"
          />

          <div>
            <h1 class="text-2xl font-bold text-zinc-100">
              {{ resume.job_title || 'Untitled Resume' }}
            </h1>
            <div class="flex flex-wrap items-center gap-3 mt-2 text-sm text-zinc-400">
              <span
                v-if="resume.company_name"
                class="flex items-center gap-1"
              >
                <UIcon
                  name="i-heroicons-building-office-2"
                  class="w-4 h-4"
                />
                {{ resume.company_name }}
              </span>
              <span class="flex items-center gap-1">
                <UIcon
                  name="i-heroicons-calendar"
                  class="w-4 h-4"
                />
                {{ formattedDate }}
              </span>
              <UBadge
                :color="statusColor"
                variant="subtle"
                size="sm"
              >
                {{ formatStatus(resume.status) }}
              </UBadge>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2 flex-wrap">
          <UButton
            variant="soft"
            color="primary"
            size="sm"
            icon="i-heroicons-sparkles"
            :loading="isRetailoring"
            @click="retailorResume"
          >
            Re-tailor
          </UButton>
          <UButton
            variant="soft"
            color="primary"
            size="sm"
            icon="i-heroicons-arrow-down-tray"
            :loading="isDownloading"
            @click="downloadPDF"
          >
            Download PDF
          </UButton>
          <UButton
            variant="ghost"
            color="neutral"
            size="sm"
            icon="i-heroicons-arrow-path"
            :loading="isRegenerating"
            title="Force regenerate PDF (ignoring cache)"
            @click="regeneratePDF"
          >
            Regenerate PDF
          </UButton>
          <UButton
            v-if="resume.status !== 'reviewed'"
            variant="ghost"
            color="neutral"
            size="sm"
            icon="i-heroicons-check"
            :loading="isUpdatingStatus"
            @click="updateStatus('reviewed')"
          >
            Mark Reviewed
          </UButton>
          <UButton
            variant="ghost"
            color="error"
            size="sm"
            icon="i-heroicons-trash"
            @click="showDeleteModal = true"
          />
        </div>
      </div>

      <!-- Main Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Left: Job Description -->
        <div class="bg-zinc-900 rounded-xl border border-zinc-800 p-6">
          <div class="flex items-center gap-2 mb-4">
            <UIcon
              name="i-heroicons-briefcase"
              class="w-5 h-5 text-emerald-400"
            />
            <h2 class="text-lg font-semibold text-zinc-100">Job Description</h2>
          </div>
          <div class="prose prose-invert prose-sm max-w-none">
            <pre
              class="whitespace-pre-wrap text-zinc-300 text-sm leading-relaxed font-mono bg-zinc-950 p-4 rounded-lg overflow-auto max-h-125"
              >{{ resume.job_description || 'No job description available.' }}</pre
            >
          </div>
          <div
            v-if="resume.job_url"
            class="mt-4 pt-4 border-t border-zinc-800"
          >
            <a
              :href="resume.job_url"
              target="_blank"
              rel="noopener noreferrer"
              class="text-sm text-emerald-400 hover:text-emerald-300 flex items-center gap-1"
            >
              <UIcon
                name="i-heroicons-link"
                class="w-4 h-4"
              />
              View Original Posting
            </a>
          </div>
        </div>

        <!-- Right: Generated Content -->
        <div class="bg-zinc-900 rounded-xl border border-zinc-800 p-6">
          <div class="flex items-center gap-2 mb-4">
            <UIcon
              name="i-heroicons-document-text"
              class="w-5 h-5 text-emerald-400"
            />
            <h2 class="text-lg font-semibold text-zinc-100">Generated Resume</h2>
          </div>

          <!-- Generated Content Display -->
          <div
            v-if="resume.generated_content"
            class="prose prose-invert prose-sm max-w-none"
          >
            <pre
              class="whitespace-pre-wrap text-zinc-300 text-sm leading-relaxed font-mono bg-zinc-950 p-4 rounded-lg overflow-auto max-h-125"
              >{{ resume.generated_content }}</pre
            >
          </div>
          <div
            v-else
            class="text-center py-8"
          >
            <UIcon
              name="i-heroicons-document"
              class="w-12 h-12 text-zinc-600 mx-auto mb-4"
            />
            <p class="text-zinc-400">No content generated yet.</p>
            <UButton
              variant="soft"
              color="primary"
              size="sm"
              icon="i-heroicons-sparkles"
              class="mt-4"
              :loading="isRetailoring"
              @click="retailorResume"
            >
              Generate Content
            </UButton>
          </div>
        </div>
      </div>

      <!-- Notes Section -->
      <div class="mt-6 bg-zinc-900 rounded-xl border border-zinc-800 p-6">
        <div class="flex items-center gap-2 mb-4">
          <UIcon
            name="i-heroicons-pencil-square"
            class="w-5 h-5 text-emerald-400"
          />
          <h2 class="text-lg font-semibold text-zinc-100">Notes</h2>
        </div>
        <UTextarea
          v-model="notes"
          :rows="3"
          placeholder="Add notes about this resume..."
          class="w-full"
        />
        <p class="mt-2 text-xs text-zinc-500">Notes are saved when you update the resume status.</p>
      </div>

      <!-- Metadata Footer -->
      <div
        class="mt-6 p-4 bg-zinc-900/50 rounded-xl border border-zinc-800/50 text-sm text-zinc-500"
      >
        <div class="flex flex-wrap gap-6">
          <span>Created: {{ formattedDate }}</span>
          <span v-if="resume.updated_at !== resume.created_at"
            >Updated: {{ formattedUpdatedDate }}</span
          >
          <span>ID: {{ resume.id }}</span>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <UModal v-model:open="showDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-red-500/10 rounded-full">
              <UIcon
                name="i-heroicons-exclamation-triangle"
                class="w-6 h-6 text-red-400"
              />
            </div>
            <h3 class="text-lg font-semibold text-zinc-100">Delete Resume</h3>
          </div>

          <p class="text-zinc-400 mb-6">
            Are you sure you want to delete this resume? This action cannot be undone.
          </p>

          <div class="flex justify-end gap-3">
            <UButton
              variant="ghost"
              color="neutral"
              @click="showDeleteModal = false"
            >
              Cancel
            </UButton>
            <UButton
              variant="solid"
              color="error"
              :loading="isDeleting"
              @click="deleteResume"
            >
              Delete Resume
            </UButton>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>
