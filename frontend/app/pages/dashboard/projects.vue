<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-zinc-100">Projects</h1>
        <p class="mt-1 text-zinc-400">Showcase your side projects and technical work.</p>
      </div>
      <UButton
        color="primary"
        @click="openAddModal"
      >
        <UIcon
          name="i-lucide-plus"
          class="mr-2 h-4 w-4"
        />
        Add Project
      </UButton>
    </div>

    <!-- Stats -->
    <div class="grid gap-4 sm:grid-cols-3">
      <DashboardStatsCard
        title="Total Projects"
        :value="String(projectList.length)"
        icon="i-lucide-folder-code"
        color="primary"
      />
      <DashboardStatsCard
        title="Technologies"
        :value="String(uniqueTechnologies.length)"
        icon="i-lucide-code-2"
        color="secondary"
      />
      <DashboardStatsCard
        title="Achievements"
        :value="String(totalBullets)"
        icon="i-lucide-list"
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
          <p class="font-medium text-zinc-100">Failed to load projects</p>
          <p class="mt-1 text-sm text-zinc-400">{{ error }}</p>
        </div>
        <UButton
          color="primary"
          @click="fetchProjects()"
        >
          Try Again
        </UButton>
      </div>
    </UCard>

    <!-- Empty State -->
    <SharedEmptyState
      v-else-if="projectList.length === 0"
      icon="i-lucide-folder-code"
      title="No projects yet"
      description="Add your side projects to showcase your technical skills."
      action-label="Add Project"
      @action="openAddModal"
    />

    <!-- Project Cards -->
    <div
      v-else
      class="grid gap-4 md:grid-cols-2"
    >
      <UCard
        v-for="project in projectList"
        :key="project.id"
      >
        <template #header>
          <div class="flex items-start justify-between">
            <div class="flex items-start gap-4">
              <div
                class="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-500/10"
              >
                <UIcon
                  name="i-lucide-folder-code"
                  class="h-5 w-5 text-emerald-400"
                />
              </div>
              <div>
                <h3 class="font-semibold text-zinc-100">
                  {{ project.name }}
                </h3>
                <p
                  v-if="project.description"
                  class="mt-1 text-sm text-zinc-400 line-clamp-2"
                >
                  {{ project.description }}
                </p>
              </div>
            </div>
            <UDropdown :items="getDropdownItems(project)">
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-more-vertical"
                size="sm"
              />
            </UDropdown>
          </div>
        </template>

        <div class="space-y-4">
          <!-- Tech Stack -->
          <div
            v-if="project.tech_stack && project.tech_stack.length > 0"
            class="flex flex-wrap gap-2"
          >
            <UBadge
              v-for="tech in project.tech_stack"
              :key="tech"
              color="primary"
              variant="subtle"
              size="sm"
            >
              {{ tech }}
            </UBadge>
          </div>

          <!-- Date Range -->
          <div
            v-if="project.start_date"
            class="flex items-center gap-2 text-sm text-zinc-400"
          >
            <UIcon
              name="i-lucide-calendar"
              class="h-4 w-4"
            />
            {{ formatDateRange(project) }}
          </div>

          <!-- Links -->
          <div class="flex gap-3">
            <a
              v-if="project.url"
              :href="project.url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-1 text-sm text-emerald-400 hover:text-emerald-300"
            >
              <UIcon
                name="i-lucide-external-link"
                class="h-4 w-4"
              />
              Demo
            </a>
            <a
              v-if="project.repository_url"
              :href="project.repository_url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-1 text-sm text-emerald-400 hover:text-emerald-300"
            >
              <UIcon
                name="i-lucide-github"
                class="h-4 w-4"
              />
              Source
            </a>
          </div>

          <!-- Bullets -->
          <div
            v-if="project.bullets && project.bullets.length > 0"
            class="space-y-2"
          >
            <div
              v-for="bullet in project.bullets"
              :key="bullet.id"
              class="group relative rounded-lg border border-zinc-800 bg-zinc-900 p-2 text-sm"
            >
              <div class="flex items-start gap-2">
                <UIcon
                  name="i-lucide-chevron-right"
                  class="mt-0.5 h-4 w-4 text-emerald-400"
                />
                <span class="flex-1 text-zinc-300">{{ bullet.content }}</span>
                <UButton
                  color="neutral"
                  variant="ghost"
                  icon="i-lucide-trash-2"
                  size="xs"
                  class="opacity-0 transition-opacity group-hover:opacity-100"
                  @click="deleteBullet(project.id, bullet.id)"
                />
              </div>
            </div>
          </div>

          <!-- Add Bullet -->
          <UButton
            color="neutral"
            variant="outline"
            block
            size="sm"
            @click="openAddBulletModal(project)"
          >
            <UIcon
              name="i-lucide-plus"
              class="mr-2 h-4 w-4"
            />
            Add Achievement
          </UButton>
        </div>
      </UCard>
    </div>

    <!-- Add/Edit Project Modal -->
    <UModal v-model:open="isModalOpen">
      <template #content>
        <UCard class="max-h-[80vh] overflow-y-auto">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-zinc-100">
                {{ isEditing ? 'Edit Project' : 'Add Project' }}
              </h3>
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
            @submit.prevent="saveProject"
          >
            <!-- Name -->
            <UFormField
              label="Project Name"
              required
            >
              <UInput
                v-model="form.name"
                placeholder="Chameleon Vitae"
                :color="formErrors.name ? 'error' : undefined"
              />
              <template
                v-if="formErrors.name"
                #error
              >
                {{ formErrors.name }}
              </template>
            </UFormField>

            <!-- Description -->
            <UFormField label="Description">
              <UTextarea
                v-model="form.description"
                placeholder="AI-powered resume engineering tool..."
                :rows="3"
              />
            </UFormField>

            <!-- Tech Stack -->
            <UFormField label="Tech Stack">
              <div class="space-y-2">
                <div class="flex flex-wrap gap-2">
                  <UBadge
                    v-for="(tech, index) in form.techStack"
                    :key="index"
                    color="primary"
                    variant="subtle"
                    class="cursor-pointer"
                    @click="removeTech(index)"
                  >
                    {{ tech }}
                    <UIcon
                      name="i-lucide-x"
                      class="ml-1 h-3 w-3"
                    />
                  </UBadge>
                </div>
                <div class="flex gap-2">
                  <UInput
                    v-model="newTech"
                    placeholder="Add technology (e.g., Go, Vue.js)"
                    class="flex-1"
                    @keydown.enter.prevent="addTech"
                  />
                  <UButton
                    color="neutral"
                    variant="outline"
                    @click="addTech"
                  >
                    Add
                  </UButton>
                </div>
              </div>
            </UFormField>

            <!-- URLs -->
            <div class="grid grid-cols-2 gap-4">
              <UFormField label="Demo URL">
                <UInput
                  v-model="form.url"
                  placeholder="https://project.demo.com"
                  type="url"
                />
              </UFormField>
              <UFormField label="Repository URL">
                <UInput
                  v-model="form.repositoryUrl"
                  placeholder="https://github.com/..."
                  type="url"
                />
              </UFormField>
            </div>

            <!-- Date Range -->
            <div class="grid grid-cols-2 gap-4">
              <UFormField label="Start Date">
                <UInput
                  v-model="form.startDate"
                  type="date"
                />
              </UFormField>
              <UFormField label="End Date">
                <UInput
                  v-model="form.endDate"
                  type="date"
                  placeholder="Leave empty if ongoing"
                />
              </UFormField>
            </div>

            <!-- Initial Bullets (for new projects) -->
            <UFormField
              v-if="!isEditing"
              label="Initial Achievements"
            >
              <div class="space-y-2">
                <div
                  v-for="(bullet, index) in form.bullets"
                  :key="index"
                  class="flex gap-2"
                >
                  <UInput
                    v-model="form.bullets[index]"
                    class="flex-1"
                    placeholder="Developed feature X..."
                  />
                  <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-lucide-trash-2"
                    @click="form.bullets.splice(index, 1)"
                  />
                </div>
                <UButton
                  color="neutral"
                  variant="outline"
                  block
                  size="sm"
                  @click="form.bullets.push('')"
                >
                  <UIcon
                    name="i-lucide-plus"
                    class="mr-2 h-4 w-4"
                  />
                  Add Achievement
                </UButton>
              </div>
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
                @click="saveProject"
              >
                {{ isEditing ? 'Update' : 'Create' }}
              </UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Add Bullet Modal -->
    <UModal v-model:open="isBulletModalOpen">
      <template #content>
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-zinc-100">Add Achievement</h3>
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                size="sm"
                @click="isBulletModalOpen = false"
              />
            </div>
          </template>

          <UFormField
            label="Achievement"
            required
          >
            <UTextarea
              v-model="bulletContent"
              placeholder="Implemented AI-driven bullet selection algorithm that increased relevance scores by 40%"
              :rows="3"
            />
          </UFormField>

          <template #footer>
            <div class="flex justify-end gap-3">
              <UButton
                color="neutral"
                variant="outline"
                @click="isBulletModalOpen = false"
              >
                Cancel
              </UButton>
              <UButton
                color="primary"
                :loading="savingBullet"
                @click="saveBullet"
              >
                Add
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
            <h3 class="text-lg font-semibold text-zinc-100">Delete Project</h3>
          </template>

          <p class="text-zinc-400">
            Are you sure you want to delete
            <strong class="text-zinc-100">{{ deleteTarget?.name }}</strong
            >? This will also remove all achievements. This action cannot be undone.
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
  ProjectResponse,
  CreateProjectRequest,
  UpdateProjectRequest,
  ListProjectsResponse
} from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

// State
const projectList = ref<ProjectResponse[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const saving = ref(false)
const savingBullet = ref(false)
const deleting = ref(false)

// Modal state
const isModalOpen = ref(false)
const isBulletModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)
const deleteTarget = ref<ProjectResponse | null>(null)
const bulletProjectId = ref<string | null>(null)

// Form state
const form = ref({
  name: '',
  description: '',
  techStack: [] as string[],
  url: '',
  repositoryUrl: '',
  startDate: '',
  endDate: '',
  bullets: [] as string[]
})

const formErrors = ref<Record<string, string>>({})
const newTech = ref('')
const bulletContent = ref('')

// Composables
const toast = useToast()

// Computed
const uniqueTechnologies = computed(() => {
  const techs = new Set<string>()
  projectList.value.forEach((project) => {
    project.tech_stack?.forEach((tech) => techs.add(tech))
  })
  return Array.from(techs)
})

const totalBullets = computed(() => {
  return projectList.value.reduce((sum, project) => sum + (project.bullets?.length || 0), 0)
})

// Methods
async function fetchProjects() {
  loading.value = true
  error.value = null

  try {
    const response = await apiFetch<ListProjectsResponse>('/projects')
    projectList.value = response.data || []
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load projects'
  } finally {
    loading.value = false
  }
}

function formatDateRange(project: ProjectResponse): string {
  const formatDate = (dateStr: string | null | undefined) => {
    if (!dateStr) return null
    const date = new Date(dateStr)
    return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
  }

  const start = formatDate(project.start_date)
  const end = formatDate(project.end_date)

  if (!start && !end) return ''
  if (!start) return `Until ${end}`
  if (!end) return `${start} — Present`
  return `${start} — ${end}`
}

function openAddModal() {
  isEditing.value = false
  editingId.value = null
  resetForm()
  isModalOpen.value = true
}

function openEditModal(project: ProjectResponse) {
  isEditing.value = true
  editingId.value = project.id
  form.value = {
    name: project.name,
    description: project.description || '',
    techStack: [...(project.tech_stack || [])],
    url: project.url || '',
    repositoryUrl: project.repository_url || '',
    startDate: project.start_date?.split('T')[0] || '',
    endDate: project.end_date?.split('T')[0] || '',
    bullets: []
  }
  formErrors.value = {}
  isModalOpen.value = true
}

function closeModal() {
  isModalOpen.value = false
  resetForm()
}

function resetForm() {
  form.value = {
    name: '',
    description: '',
    techStack: [],
    url: '',
    repositoryUrl: '',
    startDate: '',
    endDate: '',
    bullets: []
  }
  formErrors.value = {}
  newTech.value = ''
}

function addTech() {
  const tech = newTech.value.trim()
  if (tech && !form.value.techStack.includes(tech)) {
    form.value.techStack.push(tech)
    newTech.value = ''
  }
}

function removeTech(index: number) {
  form.value.techStack.splice(index, 1)
}

function validateForm(): boolean {
  formErrors.value = {}

  if (!form.value.name.trim()) {
    formErrors.value.name = 'Project name is required'
  }

  return Object.keys(formErrors.value).length === 0
}

async function saveProject() {
  if (!validateForm()) return

  saving.value = true

  try {
    if (isEditing.value && editingId.value) {
      // Update existing
      const payload: UpdateProjectRequest = {
        name: form.value.name,
        description: form.value.description || null,
        tech_stack: form.value.techStack,
        url: form.value.url || null,
        repository_url: form.value.repositoryUrl || null,
        start_date: form.value.startDate || null,
        end_date: form.value.endDate || null
      }

      await apiFetch(`/projects/${editingId.value}`, {
        method: 'PUT',
        body: payload
      })

      toast.add({
        title: 'Project Updated',
        description: 'Your project has been updated.',
        color: 'success'
      })
    } else {
      // Create new
      const validBullets = form.value.bullets.filter((b) => b.trim())
      const payload: CreateProjectRequest = {
        name: form.value.name,
        description: form.value.description || undefined,
        tech_stack: form.value.techStack.length > 0 ? form.value.techStack : undefined,
        url: form.value.url || undefined,
        repository_url: form.value.repositoryUrl || undefined,
        start_date: form.value.startDate || undefined,
        end_date: form.value.endDate || undefined,
        bullets: validBullets.length > 0 ? validBullets : undefined
      }

      await apiFetch('/projects', {
        method: 'POST',
        body: payload
      })

      toast.add({
        title: 'Project Added',
        description: 'Your project has been created.',
        color: 'success'
      })
    }

    closeModal()
    await fetchProjects()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to save project',
      color: 'error'
    })
  } finally {
    saving.value = false
  }
}

function openAddBulletModal(project: ProjectResponse) {
  bulletProjectId.value = project.id
  bulletContent.value = ''
  isBulletModalOpen.value = true
}

async function saveBullet() {
  if (!bulletProjectId.value || !bulletContent.value.trim()) return

  savingBullet.value = true

  try {
    await apiFetch(`/projects/${bulletProjectId.value}/bullets`, {
      method: 'POST',
      body: { content: bulletContent.value.trim() }
    })

    toast.add({
      title: 'Achievement Added',
      description: 'The achievement has been added to your project.',
      color: 'success'
    })

    isBulletModalOpen.value = false
    bulletContent.value = ''
    bulletProjectId.value = null
    await fetchProjects()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to add achievement',
      color: 'error'
    })
  } finally {
    savingBullet.value = false
  }
}

async function deleteBullet(projectId: string, bulletId: string) {
  try {
    await apiFetch(`/projects/${projectId}/bullets/${bulletId}`, {
      method: 'DELETE'
    })

    toast.add({
      title: 'Achievement Deleted',
      description: 'The achievement has been removed.',
      color: 'success'
    })

    await fetchProjects()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to delete achievement',
      color: 'error'
    })
  }
}

function openDeleteModal(project: ProjectResponse) {
  deleteTarget.value = project
  isDeleteModalOpen.value = true
}

async function confirmDelete() {
  if (!deleteTarget.value) return

  deleting.value = true

  try {
    await apiFetch(`/projects/${deleteTarget.value.id}`, {
      method: 'DELETE'
    })

    toast.add({
      title: 'Project Deleted',
      description: 'The project has been removed.',
      color: 'success'
    })

    isDeleteModalOpen.value = false
    deleteTarget.value = null
    await fetchProjects()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to delete project',
      color: 'error'
    })
  } finally {
    deleting.value = false
  }
}

function getDropdownItems(project: ProjectResponse) {
  return [
    [
      {
        label: 'Edit',
        icon: 'i-lucide-pencil',
        click: () => openEditModal(project)
      }
    ],
    [
      {
        label: 'Delete',
        icon: 'i-lucide-trash-2',
        color: 'error' as const,
        click: () => openDeleteModal(project)
      }
    ]
  ]
}

// Lifecycle
onMounted(() => {
  fetchProjects()
})
</script>
