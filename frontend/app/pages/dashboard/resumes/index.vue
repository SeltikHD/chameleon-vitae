<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-zinc-100">Resumes</h1>
        <p class="mt-1 text-zinc-400">Manage all your tailored resumes in one place.</p>
      </div>
      <UButton
        to="/dashboard/resumes/new"
        color="primary"
      >
        <UIcon
          name="i-lucide-plus"
          class="mr-2 h-4 w-4"
        />
        Create Resume
      </UButton>
    </div>

    <!-- Filters -->
    <UCard>
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center">
        <UInput
          v-model="searchQuery"
          placeholder="Search resumes..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
        />

        <USelectMenu
          v-model="statusFilter"
          :items="statusOptions"
          value-key="value"
          placeholder="All statuses"
          class="w-full sm:w-40"
        />

        <div class="ml-auto flex items-center gap-2">
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-refresh-cw"
            size="sm"
            :loading="isLoading"
            @click="fetchResumes"
          />
          <UButton
            :color="viewMode === 'table' ? 'primary' : 'neutral'"
            :variant="viewMode === 'table' ? 'solid' : 'ghost'"
            icon="i-lucide-list"
            size="sm"
            @click="viewMode = 'table'"
          />
          <UButton
            :color="viewMode === 'grid' ? 'primary' : 'neutral'"
            :variant="viewMode === 'grid' ? 'solid' : 'ghost'"
            icon="i-lucide-grid-3x3"
            size="sm"
            @click="viewMode = 'grid'"
          />
        </div>
      </div>
    </UCard>

    <!-- Loading State -->
    <template v-if="isLoading && resumes.length === 0">
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <USkeleton
          v-for="i in 6"
          :key="i"
          class="h-48 w-full rounded-lg"
        />
      </div>
    </template>

    <template v-else>
      <!-- Table View -->
      <UCard v-if="viewMode === 'table' && filteredResumes.length > 0">
        <UTable
          :data="filteredResumes"
          :columns="columns"
        >
          <template #company_name-cell="{ row }">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-800">
                <UIcon
                  name="i-lucide-building-2"
                  class="h-5 w-5 text-zinc-400"
                />
              </div>
              <div>
                <p class="font-medium text-zinc-100">
                  {{ row.original.company_name || 'Unknown Company' }}
                </p>
                <p class="text-sm text-zinc-400">
                  {{ row.original.job_title || 'Unknown Position' }}
                </p>
              </div>
            </div>
          </template>

          <template #status-cell="{ row }">
            <UBadge
              :color="getStatusColor(row.original.status)"
              variant="subtle"
              size="sm"
            >
              {{ formatStatus(row.original.status) }}
            </UBadge>
          </template>

          <template #score-cell="{ row }">
            <div class="flex items-center gap-2">
              <div class="h-2 w-20 overflow-hidden rounded-full bg-zinc-800">
                <div
                  class="h-full rounded-full"
                  :class="getScoreColor(row.original.score || 0)"
                  :style="{ width: `${row.original.score || 0}%` }"
                />
              </div>
              <span class="text-sm text-zinc-400">{{ row.original.score || 0 }}%</span>
            </div>
          </template>

          <template #created_at-cell="{ row }">
            <span class="text-sm text-zinc-400">{{ formatDate(row.original.created_at) }}</span>
          </template>

          <template #actions-cell="{ row }">
            <div class="flex items-center gap-1">
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-eye"
                size="xs"
                @click="navigateTo(`/dashboard/resumes/${row.original.id}`)"
              />
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-download"
                size="xs"
                @click="downloadResume(row.original.id)"
              />
              <UDropdownMenu
                :items="[
                  [
                    {
                      label: 'Duplicate',
                      icon: 'i-lucide-copy',
                      click: () => duplicateResume(row.original.id)
                    },
                    {
                      label: 'Edit',
                      icon: 'i-lucide-pencil',
                      click: () => navigateTo(`/dashboard/resumes/${row.original.id}`)
                    }
                  ],
                  [
                    {
                      label: 'Delete',
                      icon: 'i-lucide-trash-2',
                      color: 'error',
                      click: () => deleteResume(row.original.id)
                    }
                  ]
                ]"
              >
                <UButton
                  color="neutral"
                  variant="ghost"
                  icon="i-lucide-more-horizontal"
                  size="xs"
                />
              </UDropdownMenu>
            </div>
          </template>
        </UTable>
      </UCard>

      <!-- Grid View -->
      <div
        v-else-if="viewMode === 'grid' && filteredResumes.length > 0"
        class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
      >
        <UCard
          v-for="resume in filteredResumes"
          :key="resume.id"
          class="cursor-pointer transition-colors hover:border-emerald-500/40"
          @click="navigateTo(`/dashboard/resumes/${resume.id}`)"
        >
          <div class="space-y-4">
            <div class="flex items-start justify-between">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-zinc-800">
                <UIcon
                  name="i-lucide-building-2"
                  class="h-6 w-6 text-zinc-400"
                />
              </div>
              <UBadge
                :color="getStatusColor(resume.status)"
                variant="subtle"
                size="sm"
              >
                {{ formatStatus(resume.status) }}
              </UBadge>
            </div>

            <div>
              <h3 class="font-semibold text-zinc-100">
                {{ resume.company_name || 'Unknown Company' }}
              </h3>
              <p class="text-sm text-zinc-400">{{ resume.job_title || 'Unknown Position' }}</p>
            </div>

            <div class="space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="text-zinc-400">Match Score</span>
                <span class="font-medium text-zinc-100">{{ resume.score || 0 }}%</span>
              </div>
              <div class="h-2 overflow-hidden rounded-full bg-zinc-800">
                <div
                  class="h-full rounded-full"
                  :class="getScoreColor(resume.score || 0)"
                  :style="{ width: `${resume.score || 0}%` }"
                />
              </div>
            </div>

            <div class="flex items-center justify-between text-sm text-zinc-500">
              <span>{{ formatDate(resume.created_at) }}</span>
              <UButton
                color="primary"
                variant="ghost"
                size="xs"
                icon="i-lucide-download"
                @click.stop="downloadResume(resume.id)"
              />
            </div>
          </div>
        </UCard>
      </div>

      <!-- Empty State -->
      <SharedEmptyState
        v-if="filteredResumes.length === 0 && !isLoading"
        icon="i-lucide-file-text"
        title="No resumes found"
        :description="
          searchQuery || statusFilter
            ? 'Try adjusting your filters.'
            : 'Create your first tailored resume to get started.'
        "
        action-label="Create Resume"
        action-to="/dashboard/resumes/new"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { apiFetch } from '~/composables/useApiFetch'
import type { ResumeResponse } from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

const toast = useToast()

const searchQuery = ref('')
const statusFilter = ref('')
const viewMode = ref<'table' | 'grid'>('table')
const isLoading = ref(true)
const resumes = ref<ResumeResponse[]>([])

const statusOptions = [
  { label: 'All statuses', value: '' },
  { label: 'Draft', value: 'draft' },
  { label: 'Generated', value: 'generated' },
  { label: 'Reviewed', value: 'reviewed' },
  { label: 'Submitted', value: 'submitted' },
  { label: 'Interview', value: 'interview' },
  { label: 'Rejected', value: 'rejected' },
  { label: 'Accepted', value: 'accepted' }
]

const columns = [
  { accessorKey: 'company_name', header: 'Company' },
  { accessorKey: 'status', header: 'Status' },
  { accessorKey: 'score', header: 'Match Score' },
  { accessorKey: 'created_at', header: 'Created' },
  { accessorKey: 'actions', header: '' }
]

// Filtered resumes based on search and status
const filteredResumes = computed(() => {
  return resumes.value.filter((resume) => {
    const matchesSearch =
      searchQuery.value === '' ||
      (resume.company_name?.toLowerCase().includes(searchQuery.value.toLowerCase()) ?? false) ||
      (resume.job_title?.toLowerCase().includes(searchQuery.value.toLowerCase()) ?? false)

    const matchesStatus = statusFilter.value === '' || resume.status === statusFilter.value

    return matchesSearch && matchesStatus
  })
})

// Fetch resumes on mount
onMounted(async () => {
  await fetchResumes()
})

async function fetchResumes() {
  isLoading.value = true

  try {
    resumes.value = await apiFetch<ResumeResponse[]>('/resumes')
  } catch (error) {
    console.error('[Resumes] Failed to fetch:', error)
    toast.add({
      title: 'Error',
      description: 'Failed to load resumes. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isLoading.value = false
  }
}

function formatStatus(status: string): string {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffHours / 24)

  if (diffHours < 1) return 'Just now'
  if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`
  if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`
  if (diffDays < 30)
    return `${Math.floor(diffDays / 7)} week${Math.floor(diffDays / 7) > 1 ? 's' : ''} ago`
  return date.toLocaleDateString()
}

function getStatusColor(status: string) {
  const colors: Record<string, 'primary' | 'warning' | 'secondary' | 'success' | 'error'> = {
    draft: 'warning',
    generated: 'primary',
    reviewed: 'secondary',
    submitted: 'primary',
    interview: 'success',
    rejected: 'error',
    accepted: 'success'
  }
  return colors[status] || 'neutral'
}

function getScoreColor(score: number) {
  if (score >= 85) return 'bg-emerald-500'
  if (score >= 70) return 'bg-amber-500'
  return 'bg-red-500'
}

async function downloadResume(id: string) {
  try {
    toast.add({
      title: 'Downloading...',
      description: 'Preparing your PDF for download.',
      color: 'info',
      icon: 'i-heroicons-arrow-down-tray'
    })

    // TODO: Implement actual download when backend endpoint is ready
    console.warn('[Resumes] Download not implemented yet, resume:', id)
  } catch (error) {
    console.error('[Resumes] Download failed:', error)
    toast.add({
      title: 'Download Failed',
      description: 'Could not download the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  }
}

async function duplicateResume(id: string) {
  try {
    // Find the resume to duplicate
    const original = resumes.value.find((r) => r.id === id)
    if (!original) return

    const newResume = await apiFetch<ResumeResponse>('/resumes', {
      method: 'POST',
      body: {
        job_description: original.job_description,
        job_title: original.job_title,
        company_name: original.company_name ? `${original.company_name} (Copy)` : 'Copy',
        target_language: original.target_language
      }
    })

    resumes.value.unshift(newResume)

    toast.add({
      title: 'Resume Duplicated',
      description: 'The resume has been duplicated successfully.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })

    navigateTo(`/dashboard/resumes/${newResume.id}`)
  } catch (error) {
    console.error('[Resumes] Duplicate failed:', error)
    toast.add({
      title: 'Duplicate Failed',
      description: 'Could not duplicate the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  }
}

async function deleteResume(id: string) {
  if (!confirm('Are you sure you want to delete this resume?')) return

  try {
    await apiFetch(`/resumes/${id}`, { method: 'DELETE' })

    resumes.value = resumes.value.filter((r) => r.id !== id)

    toast.add({
      title: 'Resume Deleted',
      description: 'The resume has been successfully deleted.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })
  } catch (error) {
    console.error('[Resumes] Delete failed:', error)
    toast.add({
      title: 'Delete Failed',
      description: 'Could not delete the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  }
}
</script>
