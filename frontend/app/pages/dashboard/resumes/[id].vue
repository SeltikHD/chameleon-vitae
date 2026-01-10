<template>
  <div class="space-y-6">
    <!-- Loading State -->
    <template v-if="isLoading">
      <div class="flex items-center gap-4">
        <USkeleton class="h-10 w-10 rounded" />
        <div class="flex-1">
          <USkeleton class="mb-2 h-6 w-64" />
          <USkeleton class="h-4 w-32" />
        </div>
      </div>
      <USkeleton class="h-24 w-full rounded-lg" />
      <div class="grid gap-6 lg:grid-cols-2">
        <USkeleton class="h-96 w-full rounded-lg" />
        <USkeleton class="h-96 w-full rounded-lg" />
      </div>
    </template>

    <!-- Error State -->
    <template v-else-if="error">
      <UAlert
        color="error"
        variant="subtle"
        icon="i-heroicons-exclamation-triangle"
        :title="error"
        :description="'Could not load the resume. Please try again.'"
      />
      <div class="flex gap-4">
        <UButton
          color="primary"
          @click="fetchResume"
        >
          Retry
        </UButton>
        <UButton
          to="/dashboard/resumes"
          color="neutral"
          variant="ghost"
        >
          Back to Resumes
        </UButton>
      </div>
    </template>

    <!-- Resume Content -->
    <template v-else-if="resume">
      <!-- Header -->
      <div class="flex items-center gap-4">
        <UButton
          color="neutral"
          variant="ghost"
          icon="i-lucide-arrow-left"
          to="/dashboard/resumes"
        />
        <div class="flex-1">
          <h1 class="text-2xl font-bold text-zinc-100">
            {{ resume.company_name || 'Untitled Resume' }} - {{ resume.job_title || 'Position' }}
          </h1>
          <p class="mt-1 text-zinc-400">Created {{ formatDate(resume.created_at) }}</p>
        </div>
        <div class="flex items-center gap-2">
          <UBadge
            :color="getStatusColor(resume.status)"
            variant="subtle"
            size="lg"
          >
            {{ formatStatus(resume.status) }}
          </UBadge>
          <UButton
            color="neutral"
            variant="outline"
            icon="i-lucide-download"
            :loading="isDownloading"
            @click="downloadResume"
          >
            Download PDF
          </UButton>
          <UButton
            color="primary"
            icon="i-lucide-refresh-cw"
            :loading="isRegenerating"
            @click="regenerateResume"
          >
            Regenerate
          </UButton>
        </div>
      </div>

      <!-- Match Score Banner -->
      <UCard class="border-emerald-500/20 bg-emerald-500/5">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="flex h-14 w-14 items-center justify-center rounded-full bg-emerald-500/20">
              <UIcon
                name="i-lucide-target"
                class="h-7 w-7 text-emerald-400"
              />
            </div>
            <div>
              <p class="text-sm text-zinc-400">ATS Match Score</p>
              <p class="text-3xl font-bold text-emerald-400">{{ resume.score || 0 }}%</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm text-zinc-400">{{ selectedBullets.length }} bullets selected</p>
            <p class="text-sm text-zinc-400">{{ matchedSkillsCount }} skills matched</p>
          </div>
        </div>
      </UCard>

      <!-- Split View -->
      <div class="grid gap-6 lg:grid-cols-2">
        <!-- Left: Selected Bullets -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold text-zinc-100">Selected Experience Bullets</h2>
              <UButton
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-plus"
                @click="showAddBulletModal = true"
              >
                Add Bullet
              </UButton>
            </div>
          </template>

          <div
            v-if="selectedBullets.length > 0"
            class="space-y-4"
          >
            <div
              v-for="(group, idx) in groupedBullets"
              :key="idx"
              class="group"
            >
              <div class="mb-2 flex items-center justify-between">
                <span class="text-sm font-medium text-zinc-400">{{ group.organization }}</span>
                <span class="text-xs text-zinc-500">{{ group.period }}</span>
              </div>

              <div
                v-for="bullet in group.bullets"
                :key="bullet.id"
                class="relative mb-3 rounded-lg border border-zinc-800 bg-zinc-900 p-4 transition-colors hover:border-zinc-700"
              >
                <div class="flex items-start gap-3">
                  <UIcon
                    name="i-lucide-grip-vertical"
                    class="mt-1 h-4 w-4 cursor-grab text-zinc-600"
                  />
                  <div class="flex-1">
                    <p class="text-sm text-zinc-300">{{ bullet.content }}</p>
                    <div
                      v-if="bullet.keywords && bullet.keywords.length > 0"
                      class="mt-2 flex flex-wrap gap-2"
                    >
                      <UBadge
                        v-for="keyword in bullet.keywords"
                        :key="keyword"
                        color="primary"
                        variant="subtle"
                        size="xs"
                      >
                        {{ keyword }}
                      </UBadge>
                    </div>
                  </div>
                  <div
                    class="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100"
                  >
                    <UButton
                      color="neutral"
                      variant="ghost"
                      icon="i-lucide-pencil"
                      size="xs"
                      @click="editBullet(bullet)"
                    />
                    <UButton
                      color="neutral"
                      variant="ghost"
                      icon="i-lucide-trash-2"
                      size="xs"
                      @click="removeBullet(bullet.id)"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div
            v-else
            class="py-8 text-center text-zinc-500"
          >
            <UIcon
              name="i-lucide-file-text"
              class="mx-auto h-12 w-12 text-zinc-600"
            />
            <p class="mt-2">No bullets selected yet</p>
            <UButton
              color="primary"
              variant="ghost"
              class="mt-4"
              @click="showAddBulletModal = true"
            >
              Add Bullet
            </UButton>
          </div>
        </UCard>

        <!-- Right: PDF Preview -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h2 class="font-semibold text-zinc-100">PDF Preview</h2>
              <div class="flex items-center gap-2">
                <UButton
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  icon="i-lucide-zoom-out"
                  @click="zoomLevel = Math.max(50, zoomLevel - 25)"
                />
                <span class="text-sm text-zinc-400">{{ zoomLevel }}%</span>
                <UButton
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  icon="i-lucide-zoom-in"
                  @click="zoomLevel = Math.min(150, zoomLevel + 25)"
                />
              </div>
            </div>
          </template>

          <!-- PDF Preview Placeholder -->
          <div
            class="aspect-[8.5/11] overflow-hidden rounded-lg border border-zinc-800 bg-white p-8"
            :style="{ transform: `scale(${zoomLevel / 100})`, transformOrigin: 'top center' }"
          >
            <div class="space-y-4">
              <div class="border-b-2 border-zinc-900 pb-4">
                <h3 class="text-xl font-bold text-zinc-900">{{ userName }}</h3>
                <p class="text-sm text-zinc-600">{{ userEmail }}</p>
              </div>

              <div v-if="resume.job_description">
                <h4 class="mb-2 font-semibold text-zinc-900">Summary</h4>
                <p class="text-xs text-zinc-700">
                  {{ truncateText(resume.job_description, 200) }}
                </p>
              </div>

              <div>
                <h4 class="mb-2 font-semibold text-zinc-900">Experience</h4>
                <div class="space-y-3">
                  <div
                    v-for="(group, idx) in groupedBullets"
                    :key="idx"
                  >
                    <div class="flex items-baseline justify-between">
                      <span class="text-sm font-medium text-zinc-900">
                        {{ group.title }} @ {{ group.organization }}
                      </span>
                      <span class="text-xs text-zinc-500">{{ group.period }}</span>
                    </div>
                    <ul class="ml-4 mt-1 list-disc space-y-1">
                      <li
                        v-for="bullet in group.bullets.slice(0, 3)"
                        :key="bullet.id"
                        class="text-xs text-zinc-700"
                      >
                        {{ truncateText(bullet.content, 150) }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>

              <div v-if="skills.length > 0">
                <h4 class="mb-2 font-semibold text-zinc-900">Skills</h4>
                <p class="text-xs text-zinc-700">
                  {{ skills.map((s) => s.name).join(', ') }}
                </p>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </template>

    <!-- Add Bullet Modal -->
    <UModal v-model:open="showAddBulletModal">
      <template #content>
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="font-semibold text-zinc-100">Add Bullet from Library</h3>
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                size="xs"
                @click="showAddBulletModal = false"
              />
            </div>
          </template>

          <div class="space-y-4">
            <UInput
              v-model="bulletSearchQuery"
              placeholder="Search bullets..."
              icon="i-lucide-search"
            />

            <div class="max-h-64 space-y-2 overflow-y-auto">
              <div
                v-for="bullet in availableBullets"
                :key="bullet.id"
                class="cursor-pointer rounded-lg border border-zinc-800 p-3 transition-colors hover:border-emerald-500/40"
                @click="addBullet(bullet)"
              >
                <p class="text-sm text-zinc-300">{{ bullet.content }}</p>
                <div class="mt-2 flex flex-wrap gap-1">
                  <UBadge
                    v-for="keyword in bullet.keywords?.slice(0, 3)"
                    :key="keyword"
                    color="neutral"
                    variant="subtle"
                    size="xs"
                  >
                    {{ keyword }}
                  </UBadge>
                </div>
              </div>

              <div
                v-if="availableBullets.length === 0"
                class="py-4 text-center text-zinc-500"
              >
                No bullets found. Add some in the Experiences section.
              </div>
            </div>
          </div>
        </UCard>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { apiFetch } from '~/composables/useApiFetch'
import type { ResumeResponse, BulletResponse, ExperienceResponse, SkillResponse } from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

const route = useRoute()
const toast = useToast()
const authStore = useAuthStore()

const resumeId = computed(() => route.params.id as string)

// State
const resume = ref<ResumeResponse | null>(null)
const experiences = ref<ExperienceResponse[]>([])
const skills = ref<SkillResponse[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const isDownloading = ref(false)
const isRegenerating = ref(false)
const zoomLevel = ref(100)
const showAddBulletModal = ref(false)
const bulletSearchQuery = ref('')

// User info from auth store
const userName = computed(() => authStore.user?.name || 'Your Name')
const userEmail = computed(() => authStore.user?.email || 'email@example.com')

// Build a map of all bullets from experiences for quick lookup
const allBulletsMap = computed(() => {
  const map = new Map<string, { bullet: BulletResponse; experience: ExperienceResponse }>()
  experiences.value.forEach((exp) => {
    if (exp.bullets) {
      exp.bullets.forEach((bullet) => {
        map.set(bullet.id, { bullet, experience: exp })
      })
    }
  })
  return map
})

// Selected bullet IDs from resume
const selectedBulletIds = computed(() => {
  return new Set(resume.value?.selected_bullets || [])
})

// Selected bullets with full data
const selectedBullets = computed(() => {
  const bullets: BulletResponse[] = []
  selectedBulletIds.value.forEach((id) => {
    const data = allBulletsMap.value.get(id)
    if (data) {
      bullets.push(data.bullet)
    }
  })
  return bullets
})

// Available bullets for adding (not already selected)
const availableBullets = computed(() => {
  const allBullets: BulletResponse[] = []

  experiences.value.forEach((exp) => {
    if (exp.bullets) {
      exp.bullets.forEach((bullet) => {
        if (!selectedBulletIds.value.has(bullet.id)) {
          // Filter by search query
          if (
            bulletSearchQuery.value === '' ||
            bullet.content.toLowerCase().includes(bulletSearchQuery.value.toLowerCase())
          ) {
            allBullets.push(bullet)
          }
        }
      })
    }
  })

  return allBullets
})

// Group bullets by experience
const groupedBullets = computed(() => {
  const groups: {
    organization: string
    title: string
    period: string
    bullets: BulletResponse[]
  }[] = []

  selectedBullets.value.forEach((bullet) => {
    const data = allBulletsMap.value.get(bullet.id)

    if (data) {
      const existingGroup = groups.find((g) => g.organization === data.experience.organization)
      if (existingGroup) {
        existingGroup.bullets.push(bullet)
      } else {
        groups.push({
          organization: data.experience.organization,
          title: data.experience.title,
          period: `${formatYear(data.experience.start_date)} - ${data.experience.end_date ? formatYear(data.experience.end_date) : 'Present'}`,
          bullets: [bullet]
        })
      }
    } else {
      // Orphan bullet
      const orphanGroup = groups.find((g) => g.organization === 'Other')
      if (orphanGroup) {
        orphanGroup.bullets.push(bullet)
      } else {
        groups.push({
          organization: 'Other',
          title: '',
          period: '',
          bullets: [bullet]
        })
      }
    }
  })

  return groups
})

// Count matched skills
const matchedSkillsCount = computed(() => {
  return skills.value.length
})

// Fetch data on mount
onMounted(async () => {
  await fetchResume()
})

async function fetchResume() {
  isLoading.value = true
  error.value = null

  try {
    const [resumeRes, experiencesRes, skillsRes] = await Promise.all([
      apiFetch<ResumeResponse>(`/resumes/${resumeId.value}`),
      apiFetch<ExperienceResponse[]>('/experiences'),
      apiFetch<SkillResponse[]>('/skills')
    ])

    resume.value = resumeRes
    experiences.value = experiencesRes
    skills.value = skillsRes
  } catch (err) {
    console.error('[ResumeDetail] Failed to fetch:', err)
    error.value = err instanceof Error ? err.message : 'Failed to load resume'
  } finally {
    isLoading.value = false
  }
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
  return date.toLocaleDateString()
}

function formatYear(dateString: string): string {
  return new Date(dateString).getFullYear().toString()
}

function formatStatus(status: string): string {
  return status.charAt(0).toUpperCase() + status.slice(1)
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

function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text
  return text.slice(0, maxLength) + '...'
}

async function downloadResume() {
  isDownloading.value = true

  try {
    toast.add({
      title: 'Downloading...',
      description: 'Preparing your PDF for download.',
      color: 'info',
      icon: 'i-heroicons-arrow-down-tray'
    })

    // TODO: Implement actual download when backend endpoint is ready
    console.warn('[ResumeDetail] Download not implemented yet')

    // Download successful - could update status to 'reviewed' if needed
  } catch (err) {
    console.error('[ResumeDetail] Download failed:', err)
    toast.add({
      title: 'Download Failed',
      description: 'Could not download the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isDownloading.value = false
  }
}

async function regenerateResume() {
  isRegenerating.value = true

  try {
    await apiFetch(`/resumes/${resumeId.value}/generate`, {
      method: 'POST'
    })

    await fetchResume()

    toast.add({
      title: 'Resume Regenerated',
      description: 'Your resume has been regenerated with updated content.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })
  } catch (err) {
    console.error('[ResumeDetail] Regenerate failed:', err)
    toast.add({
      title: 'Regeneration Failed',
      description: 'Could not regenerate the resume. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  } finally {
    isRegenerating.value = false
  }
}

async function addBullet(bullet: BulletResponse) {
  try {
    await apiFetch(`/resumes/${resumeId.value}/bullets`, {
      method: 'POST',
      body: { bullet_id: bullet.id }
    })

    // Add to local state - selected_bullets is string[] of IDs
    if (resume.value) {
      if (!resume.value.selected_bullets) {
        resume.value.selected_bullets = []
      }
      resume.value.selected_bullets.push(bullet.id)
    }

    showAddBulletModal.value = false

    toast.add({
      title: 'Bullet Added',
      description: 'The bullet has been added to your resume.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })
  } catch (err) {
    console.error('[ResumeDetail] Add bullet failed:', err)
    toast.add({
      title: 'Add Failed',
      description: 'Could not add the bullet. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  }
}

async function removeBullet(bulletId: string) {
  if (!confirm('Remove this bullet from the resume?')) return

  try {
    await apiFetch(`/resumes/${resumeId.value}/bullets/${bulletId}`, {
      method: 'DELETE'
    })

    // Remove from local state
    if (resume.value && resume.value.selected_bullets) {
      resume.value.selected_bullets = resume.value.selected_bullets.filter((id) => id !== bulletId)
    }

    toast.add({
      title: 'Bullet Removed',
      description: 'The bullet has been removed from your resume.',
      color: 'success',
      icon: 'i-heroicons-check-circle'
    })
  } catch (err) {
    console.error('[ResumeDetail] Remove bullet failed:', err)
    toast.add({
      title: 'Remove Failed',
      description: 'Could not remove the bullet. Please try again.',
      color: 'error',
      icon: 'i-heroicons-exclamation-circle'
    })
  }
}

function editBullet(_bullet: BulletResponse) {
  // TODO: Implement bullet editing modal
  toast.add({
    title: 'Coming Soon',
    description: 'Bullet editing will be available soon.',
    color: 'info',
    icon: 'i-heroicons-information-circle'
  })
}
</script>
