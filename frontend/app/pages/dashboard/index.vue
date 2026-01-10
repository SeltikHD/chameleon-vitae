<template>
  <div class="space-y-8">
    <!-- Welcome Header -->
    <div>
      <h1 class="text-2xl font-bold text-zinc-100">Welcome back, {{ displayName }}!</h1>
      <p class="mt-1 text-zinc-400">Here's what's happening with your resumes this week.</p>
    </div>

    <!-- Loading State -->
    <template v-if="isLoading">
      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <USkeleton
          v-for="i in 4"
          :key="i"
          class="h-24 w-full rounded-lg"
        />
      </div>
      <USkeleton class="h-64 w-full rounded-lg" />
    </template>

    <template v-else>
      <!-- Stats Grid -->
      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <DashboardStatsCard
          v-for="stat in stats"
          :key="stat.title"
          :title="stat.title"
          :value="stat.value"
          :icon="stat.icon"
          :color="stat.color"
          :change="stat.change"
          :change-type="stat.changeType"
        />
      </div>

      <!-- Recent Resumes -->
      <div>
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-zinc-100">Recent Resumes</h2>
          <UButton
            to="/dashboard/resumes"
            color="neutral"
            variant="ghost"
            size="sm"
          >
            View all
            <UIcon
              name="i-lucide-arrow-right"
              class="ml-1 h-4 w-4"
            />
          </UButton>
        </div>

        <UCard v-if="recentResumes.length > 0">
          <UTable
            :data="recentResumes"
            :columns="resumeColumns"
          >
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
                <div class="h-2 w-16 overflow-hidden rounded-full bg-zinc-800">
                  <div
                    class="h-full rounded-full bg-emerald-500"
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
              <UDropdownMenu
                :items="[
                  [
                    {
                      label: 'View',
                      icon: 'i-lucide-eye',
                      click: () => navigateTo(`/dashboard/resumes/${row.original.id}`)
                    },
                    {
                      label: 'Download',
                      icon: 'i-lucide-download',
                      click: () => downloadResume(row.original.id)
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
            </template>
          </UTable>
        </UCard>

        <!-- Empty State -->
        <SharedEmptyState
          v-else
          icon="i-lucide-file-text"
          title="No resumes yet"
          description="Create your first tailored resume to get started."
          action-label="Create Resume"
          action-to="/dashboard/resumes/new"
        />
      </div>

      <!-- Quick Actions -->
      <div>
        <h2 class="mb-4 text-lg font-semibold text-zinc-100">Quick Actions</h2>

        <div class="grid gap-4 sm:grid-cols-3">
          <UCard
            class="cursor-pointer transition-colors hover:border-emerald-500/40"
            @click="navigateTo('/dashboard/resumes/new')"
          >
            <div class="flex items-center gap-4">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-emerald-500/10">
                <UIcon
                  name="i-lucide-plus"
                  class="h-6 w-6 text-emerald-400"
                />
              </div>
              <div>
                <h3 class="font-medium text-zinc-100">Create Resume</h3>
                <p class="text-sm text-zinc-400">Tailor for a new job</p>
              </div>
            </div>
          </UCard>

          <UCard
            class="cursor-pointer transition-colors hover:border-violet-500/40"
            @click="navigateTo('/dashboard/experiences')"
          >
            <div class="flex items-center gap-4">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-violet-500/10">
                <UIcon
                  name="i-lucide-briefcase"
                  class="h-6 w-6 text-violet-400"
                />
              </div>
              <div>
                <h3 class="font-medium text-zinc-100">Add Experience</h3>
                <p class="text-sm text-zinc-400">Expand your bullet library</p>
              </div>
            </div>
          </UCard>

          <UCard
            class="cursor-pointer transition-colors hover:border-sky-500/40"
            @click="navigateTo('/dashboard/skills')"
          >
            <div class="flex items-center gap-4">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-sky-500/10">
                <UIcon
                  name="i-lucide-tags"
                  class="h-6 w-6 text-sky-400"
                />
              </div>
              <div>
                <h3 class="font-medium text-zinc-100">Manage Skills</h3>
                <p class="text-sm text-zinc-400">Update your skill matrix</p>
              </div>
            </div>
          </UCard>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { apiFetch } from '~/composables/useApiFetch'
import type { ResumeResponse, ExperienceResponse, SkillResponse } from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

const authStore = useAuthStore()
const toast = useToast()

// Display name from auth store
const displayName = computed(() => {
  const name = authStore.user?.name || authStore.user?.email || 'there'
  // Get first name only
  const parts = name.split(' ')
  const firstName = parts[0] ?? name
  return firstName.split('@')[0] ?? firstName
})

// Data refs
const resumes = ref<ResumeResponse[]>([])
const experiences = ref<ExperienceResponse[]>([])
const skills = ref<SkillResponse[]>([])
const isLoading = ref(true)

// Computed stats
const stats = computed(() => {
  // Count total bullets across all experiences
  const bulletCount = experiences.value.reduce(
    (total, exp) => total + (exp.bullets?.length || 0),
    0
  )

  // Calculate average match score
  const resumesWithScore = resumes.value.filter((r) => r.score && r.score > 0)
  const avgMatchScore =
    resumesWithScore.length > 0
      ? Math.round(
          resumesWithScore.reduce((sum, r) => sum + (r.score || 0), 0) / resumesWithScore.length
        )
      : 0

  // Count submitted resumes
  const submittedCount = resumes.value.filter((r) => r.status === 'submitted').length

  return [
    {
      title: 'Total Resumes',
      value: String(resumes.value.length),
      icon: 'i-lucide-file-text',
      color: 'primary' as const,
      change: 0,
      changeType: 'neutral' as const
    },
    {
      title: 'Experience Bullets',
      value: String(bulletCount),
      icon: 'i-lucide-layers',
      color: 'secondary' as const,
      change: 0,
      changeType: 'neutral' as const
    },
    {
      title: 'Avg. Match Score',
      value: avgMatchScore > 0 ? `${avgMatchScore}%` : '-',
      icon: 'i-lucide-target',
      color: 'primary' as const,
      change: 0,
      changeType: 'neutral' as const
    },
    {
      title: 'Submitted',
      value: String(submittedCount),
      icon: 'i-lucide-send',
      color: 'primary' as const,
      change: 0,
      changeType: 'neutral' as const
    }
  ]
})

// Recent resumes (last 5)
const recentResumes = computed(() => {
  return [...resumes.value]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5)
})

// Resume columns
const resumeColumns = [
  { accessorKey: 'company_name', header: 'Company' },
  { accessorKey: 'job_title', header: 'Position' },
  { accessorKey: 'status', header: 'Status' },
  { accessorKey: 'score', header: 'Match Score' },
  { accessorKey: 'created_at', header: 'Created' },
  { accessorKey: 'actions', header: '' }
]

// Fetch data on mount
onMounted(async () => {
  await fetchDashboardData()
})

async function fetchDashboardData() {
  isLoading.value = true

  try {
    // Fetch all data in parallel
    const [resumesRes, experiencesRes, skillsRes] = await Promise.all([
      apiFetch<ResumeResponse[]>('/resumes'),
      apiFetch<ExperienceResponse[]>('/experiences'),
      apiFetch<SkillResponse[]>('/skills')
    ])

    resumes.value = resumesRes
    experiences.value = experiencesRes
    skills.value = skillsRes
  } catch (error) {
    console.error('[Dashboard] Failed to fetch data:', error)
    toast.add({
      title: 'Error',
      description: 'Failed to load dashboard data. Please try again.',
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

async function downloadResume(id: string) {
  try {
    toast.add({
      title: 'Downloading...',
      description: 'Preparing your PDF for download.',
      color: 'info',
      icon: 'i-heroicons-arrow-down-tray'
    })

    // TODO: Implement actual download when backend endpoint is ready
    console.warn('[Dashboard] Download not implemented yet, resume:', id)
  } catch (error) {
    console.error('[Dashboard] Download failed:', error)
    toast.add({
      title: 'Download Failed',
      description: 'Could not download the resume. Please try again.',
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
    console.error('[Dashboard] Delete failed:', error)
    toast.add({
      title: 'Delete Failed',
      description: 'Could not delete the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  }
}
</script>
