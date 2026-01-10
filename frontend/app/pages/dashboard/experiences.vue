<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-zinc-100">Experiences</h1>
        <p class="mt-1 text-zinc-400">
          Your career journey broken down into powerful, reusable bullets.
        </p>
      </div>
      <UButton
        color="primary"
        @click="openAddModal"
      >
        <UIcon
          name="i-lucide-plus"
          class="mr-2 h-4 w-4"
        />
        Add Experience
      </UButton>
    </div>

    <!-- Stats -->
    <div class="grid gap-4 sm:grid-cols-3">
      <DashboardStatsCard
        title="Total Experiences"
        :value="String(experiences.length)"
        icon="i-lucide-briefcase"
        color="primary"
      />
      <DashboardStatsCard
        title="Total Bullets"
        :value="String(totalBullets)"
        icon="i-lucide-list"
        color="secondary"
      />
      <DashboardStatsCard
        title="Keywords"
        :value="String(uniqueKeywords.length)"
        icon="i-lucide-tags"
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
          <p class="font-medium text-zinc-100">Failed to load experiences</p>
          <p class="mt-1 text-sm text-zinc-400">{{ error }}</p>
        </div>
        <UButton
          color="primary"
          @click="fetchExperiences()"
        >
          Try Again
        </UButton>
      </div>
    </UCard>

    <!-- Empty State -->
    <SharedEmptyState
      v-else-if="experiences.length === 0"
      icon="i-lucide-briefcase"
      title="No experiences yet"
      description="Add your work experiences to build your bullet library."
      action-label="Add Experience"
      @action="openAddModal"
    />

    <!-- Timeline View -->
    <div
      v-else
      class="relative"
    >
      <!-- Timeline Line -->
      <div class="absolute bottom-0 left-6 top-0 w-px bg-zinc-800" />

      <!-- Experience Cards -->
      <div class="space-y-8">
        <div
          v-for="experience in experiences"
          :key="experience.id"
          class="relative pl-16"
        >
          <!-- Timeline Dot -->
          <div
            class="absolute left-0 flex h-12 w-12 items-center justify-center rounded-full border-4 border-zinc-950 bg-zinc-800"
          >
            <UIcon
              :name="getExperienceIcon(experience.type)"
              class="h-5 w-5 text-emerald-400"
            />
          </div>

          <UCard>
            <template #header>
              <div class="flex items-start justify-between">
                <div>
                  <h3 class="font-semibold text-zinc-100">
                    {{ experience.title }}
                  </h3>
                  <p class="text-sm text-zinc-400">{{ experience.organization }}</p>
                </div>
                <div class="text-right">
                  <UBadge
                    color="neutral"
                    variant="subtle"
                  >
                    {{ formatDateRange(experience) }}
                  </UBadge>
                  <p class="mt-1 text-xs text-zinc-500">
                    {{ experience.bullets?.length || 0 }} bullets
                  </p>
                </div>
              </div>
            </template>

            <!-- Bullets -->
            <div class="space-y-3">
              <div
                v-for="bullet in experience.bullets"
                :key="bullet.id"
                class="group relative rounded-lg border border-zinc-800 bg-zinc-900 p-3 transition-colors hover:border-zinc-700"
              >
                <div class="flex items-start gap-3">
                  <UIcon
                    name="i-lucide-chevron-right"
                    class="mt-1 h-4 w-4 text-emerald-400"
                  />
                  <div class="flex-1">
                    <p class="text-sm text-zinc-300">{{ bullet.content }}</p>
                    <div class="mt-2 flex flex-wrap gap-1">
                      <UBadge
                        v-for="keyword in bullet.keywords"
                        :key="keyword"
                        color="primary"
                        variant="subtle"
                        size="xs"
                      >
                        {{ keyword }}
                      </UBadge>
                      <UBadge
                        v-if="bullet.impact_score"
                        color="secondary"
                        variant="subtle"
                        size="xs"
                      >
                        Score: {{ bullet.impact_score }}
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
                      @click="openEditBulletModal(experience.id, bullet)"
                    />
                    <UButton
                      color="neutral"
                      variant="ghost"
                      icon="i-lucide-trash-2"
                      size="xs"
                      @click="deleteBullet(bullet.id)"
                    />
                  </div>
                </div>
              </div>

              <UButton
                color="neutral"
                variant="outline"
                block
                size="sm"
                @click="openAddBulletModal(experience.id)"
              >
                <UIcon
                  name="i-lucide-plus"
                  class="mr-2 h-4 w-4"
                />
                Add Bullet
              </UButton>
            </div>

            <template #footer>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-4 text-sm text-zinc-500">
                  <span v-if="experience.location">{{ experience.location }}</span>
                  <UBadge
                    :color="experience.is_current ? 'primary' : 'neutral'"
                    variant="subtle"
                    size="xs"
                  >
                    {{ experience.type }}
                  </UBadge>
                </div>
                <div class="flex items-center gap-2">
                  <UButton
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    icon="i-lucide-pencil"
                    @click="openEditModal(experience)"
                  >
                    Edit
                  </UButton>
                  <UButton
                    color="error"
                    variant="ghost"
                    size="xs"
                    icon="i-lucide-trash-2"
                    @click="deleteExperience(experience.id)"
                  >
                    Delete
                  </UButton>
                </div>
              </div>
            </template>
          </UCard>
        </div>
      </div>
    </div>

    <!-- Add/Edit Experience Modal -->
    <UModal v-model:open="showExperienceModal">
      <template #content>
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-zinc-100">
                {{ editingExperience ? 'Edit Experience' : 'Add Experience' }}
              </h3>
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                size="xs"
                @click="showExperienceModal = false"
              />
            </div>
          </template>

          <form
            class="space-y-4"
            @submit.prevent="saveExperience"
          >
            <UFormField
              label="Type"
              name="type"
              required
            >
              <USelectMenu
                v-model="experienceForm.type"
                :items="experienceTypes"
                value-key="value"
                class="w-full"
              />
            </UFormField>

            <UFormField
              label="Title"
              name="title"
              required
            >
              <UInput
                v-model="experienceForm.title"
                placeholder="Senior Software Engineer"
              />
            </UFormField>

            <UFormField
              label="Organization"
              name="organization"
              required
            >
              <UInput
                v-model="experienceForm.organization"
                placeholder="Company Name"
              />
            </UFormField>

            <UFormField
              label="Location"
              name="location"
            >
              <UInput
                v-model="experienceForm.location"
                placeholder="San Francisco, CA"
              />
            </UFormField>

            <div class="grid gap-4 sm:grid-cols-2">
              <UFormField
                label="Start Date"
                name="start_date"
                required
              >
                <UInput
                  v-model="experienceForm.start_date"
                  type="date"
                />
              </UFormField>

              <UFormField
                label="End Date"
                name="end_date"
              >
                <UInput
                  v-model="experienceForm.end_date"
                  type="date"
                  :disabled="experienceForm.is_current"
                />
              </UFormField>
            </div>

            <UCheckbox
              v-model="experienceForm.is_current"
              label="I currently work here"
            />

            <UFormField
              label="Description"
              name="description"
            >
              <UTextarea
                v-model="experienceForm.description"
                :rows="3"
                placeholder="Brief description of your role..."
              />
            </UFormField>

            <div class="flex justify-end gap-3 pt-4">
              <UButton
                color="neutral"
                variant="ghost"
                @click="showExperienceModal = false"
              >
                Cancel
              </UButton>
              <UButton
                type="submit"
                color="primary"
                :loading="isSaving"
              >
                {{ editingExperience ? 'Save Changes' : 'Add Experience' }}
              </UButton>
            </div>
          </form>
        </UCard>
      </template>
    </UModal>

    <!-- Add/Edit Bullet Modal -->
    <UModal v-model:open="showBulletModal">
      <template #content>
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-zinc-100">
                {{ editingBullet ? 'Edit Bullet' : 'Add Bullet' }}
              </h3>
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                size="xs"
                @click="showBulletModal = false"
              />
            </div>
          </template>

          <form
            class="space-y-4"
            @submit.prevent="saveBullet"
          >
            <UFormField
              label="Content"
              name="content"
              required
            >
              <UTextarea
                v-model="bulletForm.content"
                :rows="3"
                placeholder="Describe your achievement or responsibility..."
              />
            </UFormField>

            <UFormField
              label="Keywords"
              name="keywords"
            >
              <UInput
                v-model="bulletKeywordsInput"
                placeholder="React, TypeScript, Leadership (comma-separated)"
              />
              <template #hint>
                <span class="text-xs text-zinc-500">
                  Separate keywords with commas. These help the AI match your bullets to job
                  descriptions.
                </span>
              </template>
            </UFormField>

            <div class="flex justify-end gap-3 pt-4">
              <UButton
                color="neutral"
                variant="ghost"
                @click="showBulletModal = false"
              >
                Cancel
              </UButton>
              <UButton
                type="submit"
                color="primary"
                :loading="isSaving"
              >
                {{ editingBullet ? 'Save Changes' : 'Add Bullet' }}
              </UButton>
            </div>
          </form>
        </UCard>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { apiFetch } from '~/composables/useApiFetch'
import type {
  ExperienceResponse,
  ExperienceType,
  BulletResponse,
  CreateExperienceRequest,
  UpdateExperienceRequest,
  CreateBulletRequest,
  UpdateBulletRequest,
  ListExperiencesResponse
} from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const toast = useToast()

// State
const loading = ref(true)
const error = ref<string | null>(null)
const isSaving = ref(false)
const experiences = ref<ExperienceResponse[]>([])

// Modals
const showExperienceModal = ref(false)
const showBulletModal = ref(false)
const editingExperience = ref<ExperienceResponse | null>(null)
const editingBullet = ref<BulletResponse | null>(null)
const currentExperienceId = ref<string | null>(null)

// Form data
const experienceForm = reactive<CreateExperienceRequest>({
  type: 'work',
  title: '',
  organization: '',
  location: null,
  start_date: '',
  end_date: null,
  is_current: false,
  description: null
})

const bulletForm = reactive<CreateBulletRequest>({
  content: '',
  keywords: []
})

const bulletKeywordsInput = ref('')

// Experience types
const experienceTypes: Array<{ label: string; value: ExperienceType }> = [
  { label: 'Work', value: 'work' },
  { label: 'Education', value: 'education' },
  { label: 'Project', value: 'project' },
  { label: 'Volunteer', value: 'volunteer' },
  { label: 'Certification', value: 'certification' },
  { label: 'Award', value: 'award' }
]

// Computed
const totalBullets = computed(() =>
  experiences.value.reduce((acc, exp) => acc + (exp.bullets?.length || 0), 0)
)

const uniqueKeywords = computed(() => {
  const keywords = new Set<string>()
  experiences.value.forEach((exp) => {
    exp.bullets?.forEach((bullet) => {
      bullet.keywords?.forEach((kw) => keywords.add(kw))
    })
  })
  return Array.from(keywords)
})

// Fetch experiences
async function fetchExperiences() {
  loading.value = true
  error.value = null

  try {
    const response = await apiFetch<ListExperiencesResponse>('/experiences')
    experiences.value = response.data
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load experiences'
  } finally {
    loading.value = false
  }
}

// Experience CRUD
function openAddModal() {
  editingExperience.value = null
  Object.assign(experienceForm, {
    type: 'work',
    title: '',
    organization: '',
    location: null,
    start_date: '',
    end_date: null,
    is_current: false,
    description: null
  })
  showExperienceModal.value = true
}

function openEditModal(experience: ExperienceResponse) {
  editingExperience.value = experience
  Object.assign(experienceForm, {
    type: experience.type,
    title: experience.title,
    organization: experience.organization,
    location: experience.location,
    start_date: experience.start_date,
    end_date: experience.end_date,
    is_current: experience.is_current,
    description: experience.description
  })
  showExperienceModal.value = true
}

async function saveExperience() {
  isSaving.value = true

  try {
    if (editingExperience.value) {
      // Update existing
      const updateData: UpdateExperienceRequest = {
        type: experienceForm.type,
        title: experienceForm.title,
        organization: experienceForm.organization,
        location: experienceForm.location || null,
        start_date: experienceForm.start_date,
        end_date: experienceForm.is_current ? null : experienceForm.end_date,
        is_current: experienceForm.is_current,
        description: experienceForm.description || null
      }

      await apiFetch(`/experiences/${editingExperience.value.id}`, {
        method: 'PUT',
        body: updateData
      })

      toast.add({
        title: 'Experience Updated',
        description: 'Your experience has been saved.',
        color: 'success'
      })
    } else {
      // Create new
      await apiFetch('/experiences', {
        method: 'POST',
        body: {
          ...experienceForm,
          end_date: experienceForm.is_current ? null : experienceForm.end_date
        }
      })

      toast.add({
        title: 'Experience Added',
        description: 'Your new experience has been created.',
        color: 'success'
      })
    }

    showExperienceModal.value = false
    await fetchExperiences()
  } catch (e) {
    console.error('Failed to save experience:', e)
  } finally {
    isSaving.value = false
  }
}

async function deleteExperience(id: string) {
  if (!confirm('Are you sure you want to delete this experience and all its bullets?')) return

  try {
    await apiFetch(`/experiences/${id}`, { method: 'DELETE' })

    toast.add({
      title: 'Experience Deleted',
      description: 'The experience has been removed.',
      color: 'success'
    })

    await fetchExperiences()
  } catch (e) {
    console.error('Failed to delete experience:', e)
  }
}

// Bullet CRUD
function openAddBulletModal(experienceId: string) {
  currentExperienceId.value = experienceId
  editingBullet.value = null
  bulletForm.content = ''
  bulletForm.keywords = []
  bulletKeywordsInput.value = ''
  showBulletModal.value = true
}

function openEditBulletModal(experienceId: string, bullet: BulletResponse) {
  currentExperienceId.value = experienceId
  editingBullet.value = bullet
  bulletForm.content = bullet.content
  bulletForm.keywords = bullet.keywords || []
  bulletKeywordsInput.value = (bullet.keywords || []).join(', ')
  showBulletModal.value = true
}

async function saveBullet() {
  if (!currentExperienceId.value) return

  isSaving.value = true

  // Parse keywords from input.
  const keywords = bulletKeywordsInput.value
    .split(',')
    .map((k) => k.trim())
    .filter(Boolean)

  try {
    if (editingBullet.value) {
      // Update existing
      const updateData: UpdateBulletRequest = {
        content: bulletForm.content,
        keywords
      }

      await apiFetch(`/bullets/${editingBullet.value.id}`, {
        method: 'PUT',
        body: updateData
      })

      toast.add({
        title: 'Bullet Updated',
        description: 'Your bullet has been saved.',
        color: 'success'
      })
    } else {
      // Create new
      const createData: CreateBulletRequest = {
        content: bulletForm.content,
        keywords
      }

      await apiFetch(`/experiences/${currentExperienceId.value}/bullets`, {
        method: 'POST',
        body: createData
      })

      toast.add({
        title: 'Bullet Added',
        description: 'Your new bullet has been created.',
        color: 'success'
      })
    }

    showBulletModal.value = false
    await fetchExperiences()
  } catch (e) {
    console.error('Failed to save bullet:', e)
  } finally {
    isSaving.value = false
  }
}

async function deleteBullet(id: string) {
  if (!confirm('Are you sure you want to delete this bullet?')) return

  try {
    await apiFetch(`/bullets/${id}`, { method: 'DELETE' })

    toast.add({
      title: 'Bullet Deleted',
      description: 'The bullet has been removed.',
      color: 'success'
    })

    await fetchExperiences()
  } catch (e) {
    console.error('Failed to delete bullet:', e)
  }
}

// Helpers
function getExperienceIcon(type: ExperienceType): string {
  const icons: Record<ExperienceType, string> = {
    work: 'i-lucide-briefcase',
    education: 'i-lucide-graduation-cap',
    project: 'i-lucide-folder',
    volunteer: 'i-lucide-heart',
    certification: 'i-lucide-award',
    award: 'i-lucide-trophy'
  }
  return icons[type] || 'i-lucide-briefcase'
}

function formatDateRange(exp: ExperienceResponse): string {
  const start = new Date(exp.start_date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short'
  })
  if (exp.is_current) return `${start} - Present`
  if (exp.end_date) {
    const end = new Date(exp.end_date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short'
    })
    return `${start} - ${end}`
  }
  return start
}

// Fetch on mount
onMounted(() => {
  fetchExperiences()
})
</script>
