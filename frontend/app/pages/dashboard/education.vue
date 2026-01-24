<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-zinc-100">Education</h1>
        <p class="mt-1 text-zinc-400">Your academic background for Jake's Resume format.</p>
      </div>
      <UButton
        color="primary"
        @click="openAddModal"
      >
        <UIcon
          name="i-lucide-plus"
          class="mr-2 h-4 w-4"
        />
        Add Education
      </UButton>
    </div>

    <!-- Stats -->
    <div class="grid gap-4 sm:grid-cols-2">
      <DashboardStatsCard
        title="Total Entries"
        :value="String(educationList.length)"
        icon="i-lucide-graduation-cap"
        color="primary"
      />
      <DashboardStatsCard
        title="Honors & Awards"
        :value="String(totalHonors)"
        icon="i-lucide-award"
        color="secondary"
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
          <p class="font-medium text-zinc-100">Failed to load education</p>
          <p class="mt-1 text-sm text-zinc-400">{{ error }}</p>
        </div>
        <UButton
          color="primary"
          @click="fetchEducation()"
        >
          Try Again
        </UButton>
      </div>
    </UCard>

    <!-- Empty State -->
    <SharedEmptyState
      v-else-if="educationList.length === 0"
      icon="i-lucide-graduation-cap"
      title="No education entries yet"
      description="Add your academic background for your resume."
      action-label="Add Education"
      @action="openAddModal"
    />

    <!-- Education Cards -->
    <div
      v-else
      class="space-y-4"
    >
      <UCard
        v-for="edu in educationList"
        :key="edu.id"
      >
        <template #header>
          <div class="flex items-start justify-between">
            <div class="flex items-start gap-4">
              <div
                class="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-500/10"
              >
                <UIcon
                  name="i-lucide-graduation-cap"
                  class="h-5 w-5 text-emerald-400"
                />
              </div>
              <div>
                <h3 class="font-semibold text-zinc-100">
                  {{ edu.institution }}
                </h3>
                <p class="text-sm text-zinc-400">
                  {{ edu.degree }}
                  <span v-if="edu.field_of_study"> in {{ edu.field_of_study }}</span>
                </p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <UBadge
                color="neutral"
                variant="subtle"
              >
                {{ formatDateRange(edu) }}
              </UBadge>
              <UPopover :popper="{ placement: 'bottom-end' }">
                <UButton
                  color="neutral"
                  variant="ghost"
                  icon="i-lucide-more-vertical"
                  size="sm"
                />

                <template #content="{ close }">
                  <div class="p-1 min-w-35 flex flex-col gap-0.5">
                    <UButton
                      color="neutral"
                      variant="ghost"
                      icon="i-lucide-pencil"
                      label="Edit"
                      class="justify-start w-full"
                      @click="openEditModal(edu, close)"
                    />

                    <UButton
                      color="error"
                      variant="ghost"
                      icon="i-lucide-trash-2"
                      label="Delete"
                      class="justify-start w-full text-red-500 hover:bg-red-50 dark:hover:bg-red-950/30"
                      @click="openDeleteModal(edu, close)"
                    />
                  </div>
                </template>
              </UPopover>
            </div>
          </div>
        </template>

        <div class="space-y-3">
          <!-- Location -->
          <div
            v-if="edu.location"
            class="flex items-center gap-2 text-sm text-zinc-400"
          >
            <UIcon
              name="i-lucide-map-pin"
              class="h-4 w-4"
            />
            {{ edu.location }}
          </div>

          <!-- GPA -->
          <div
            v-if="edu.gpa"
            class="flex items-center gap-2 text-sm text-zinc-400"
          >
            <UIcon
              name="i-lucide-bar-chart"
              class="h-4 w-4"
            />
            GPA: {{ edu.gpa }}
          </div>

          <!-- Honors -->
          <div
            v-if="edu.honors && edu.honors.length > 0"
            class="flex flex-wrap gap-2"
          >
            <UBadge
              v-for="honor in edu.honors"
              :key="honor"
              color="primary"
              variant="subtle"
              size="sm"
            >
              {{ honor }}
            </UBadge>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Add/Edit Modal -->
    <UModal v-model:open="isModalOpen">
      <template #content>
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-zinc-100">
                {{ isEditing ? 'Edit Education' : 'Add Education' }}
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
            @submit.prevent="saveEducation"
          >
            <!-- Institution -->
            <UFormField
              label="Institution"
              required
            >
              <UInput
                v-model="form.institution"
                placeholder="Massachusetts Institute of Technology"
                :color="formErrors.institution ? 'error' : undefined"
              />
              <template
                v-if="formErrors.institution"
                #error
              >
                {{ formErrors.institution }}
              </template>
            </UFormField>

            <!-- Degree -->
            <UFormField
              label="Degree"
              required
            >
              <UInput
                v-model="form.degree"
                placeholder="Bachelor of Science"
                :color="formErrors.degree ? 'error' : undefined"
              />
              <template
                v-if="formErrors.degree"
                #error
              >
                {{ formErrors.degree }}
              </template>
            </UFormField>

            <!-- Field of Study -->
            <UFormField label="Field of Study">
              <UInput
                v-model="form.fieldOfStudy"
                placeholder="Computer Science"
              />
            </UFormField>

            <!-- Location -->
            <UFormField label="Location">
              <UInput
                v-model="form.location"
                placeholder="Cambridge, MA"
              />
            </UFormField>

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

            <!-- GPA -->
            <UFormField label="GPA">
              <UInput
                v-model="form.gpa"
                placeholder="3.85/4.0"
              />
            </UFormField>

            <!-- Honors -->
            <UFormField label="Honors & Awards">
              <div class="space-y-2">
                <div class="flex flex-wrap gap-2">
                  <UBadge
                    v-for="(honor, index) in form.honors"
                    :key="index"
                    color="primary"
                    variant="subtle"
                    class="cursor-pointer"
                    @click="removeHonor(index)"
                  >
                    {{ honor }}
                    <UIcon
                      name="i-lucide-x"
                      class="ml-1 h-3 w-3"
                    />
                  </UBadge>
                </div>
                <div class="flex gap-2">
                  <UInput
                    v-model="newHonor"
                    placeholder="Add honor (e.g., Magna Cum Laude)"
                    class="flex-1"
                    @keydown.enter.prevent="addHonor"
                  />
                  <UButton
                    color="neutral"
                    variant="outline"
                    @click="addHonor"
                  >
                    Add
                  </UButton>
                </div>
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
                @click="saveEducation"
              >
                {{ isEditing ? 'Update' : 'Create' }}
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
            <h3 class="text-lg font-semibold text-zinc-100">Delete Education</h3>
          </template>

          <p class="text-zinc-400">
            Are you sure you want to delete
            <strong class="text-zinc-100">{{ deleteTarget?.institution }}</strong
            >? This action cannot be undone.
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
  EducationResponse,
  CreateEducationRequest,
  UpdateEducationRequest,
  ListEducationResponse
} from '~/types/api'

definePageMeta({
  layout: 'dashboard'
})

// State
const educationList = ref<EducationResponse[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const saving = ref(false)
const deleting = ref(false)

// Modal state
const isModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)
const deleteTarget = ref<EducationResponse | null>(null)

// Form state
const form = ref({
  institution: '',
  degree: '',
  fieldOfStudy: '',
  location: '',
  startDate: '',
  endDate: '',
  gpa: '',
  honors: [] as string[]
})

const formErrors = ref<Record<string, string>>({})
const newHonor = ref('')

// Composables
const toast = useToast()

// Computed
const totalHonors = computed(() => {
  return educationList.value.reduce((sum, edu) => sum + (edu.honors?.length || 0), 0)
})

// Methods
async function fetchEducation() {
  loading.value = true
  error.value = null

  try {
    const response = await apiFetch<ListEducationResponse>('/education')
    educationList.value = response.data || []
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load education'
  } finally {
    loading.value = false
  }
}

function formatDateRange(edu: EducationResponse): string {
  const formatDate = (dateStr: string | null | undefined) => {
    if (!dateStr) return null
    const date = new Date(dateStr)
    return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
  }

  const start = formatDate(edu.start_date)
  const end = formatDate(edu.end_date)

  if (!start && !end) return 'Date not set'
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

function openEditModal(edu: EducationResponse, closePopover: () => void) {
  isEditing.value = true
  editingId.value = edu.id
  form.value = {
    institution: edu.institution,
    degree: edu.degree,
    fieldOfStudy: edu.field_of_study || '',
    location: edu.location || '',
    startDate: edu.start_date?.split('T')[0] || '',
    endDate: edu.end_date?.split('T')[0] || '',
    gpa: edu.gpa || '',
    honors: [...(edu.honors || [])]
  }
  formErrors.value = {}
  isModalOpen.value = true
  closePopover()
}

function closeModal() {
  isModalOpen.value = false
  resetForm()
}

function resetForm() {
  form.value = {
    institution: '',
    degree: '',
    fieldOfStudy: '',
    location: '',
    startDate: '',
    endDate: '',
    gpa: '',
    honors: []
  }
  formErrors.value = {}
  newHonor.value = ''
}

function addHonor() {
  const honor = newHonor.value.trim()
  if (honor && !form.value.honors.includes(honor)) {
    form.value.honors.push(honor)
    newHonor.value = ''
  }
}

function removeHonor(index: number) {
  form.value.honors.splice(index, 1)
}

function validateForm(): boolean {
  formErrors.value = {}

  if (!form.value.institution.trim()) {
    formErrors.value.institution = 'Institution is required'
  }
  if (!form.value.degree.trim()) {
    formErrors.value.degree = 'Degree is required'
  }

  return Object.keys(formErrors.value).length === 0
}

async function saveEducation() {
  if (!validateForm()) return

  saving.value = true

  try {
    if (isEditing.value && editingId.value) {
      // Update existing
      const payload: UpdateEducationRequest = {
        institution: form.value.institution,
        degree: form.value.degree,
        field_of_study: form.value.fieldOfStudy || null,
        location: form.value.location || null,
        start_date: form.value.startDate || null,
        end_date: form.value.endDate || null,
        gpa: form.value.gpa || null,
        honors: form.value.honors
      }

      await apiFetch(`/education/${editingId.value}`, {
        method: 'PUT',
        body: payload
      })

      toast.add({
        title: 'Education Updated',
        description: 'Your education entry has been updated.',
        color: 'success'
      })
    } else {
      // Create new
      const payload: CreateEducationRequest = {
        institution: form.value.institution,
        degree: form.value.degree,
        field_of_study: form.value.fieldOfStudy || '',
        location: form.value.location || undefined,
        start_date: form.value.startDate || undefined,
        end_date: form.value.endDate || undefined,
        gpa: form.value.gpa || undefined,
        honors: form.value.honors.length > 0 ? form.value.honors : undefined
      }

      await apiFetch('/education', {
        method: 'POST',
        body: payload
      })

      toast.add({
        title: 'Education Added',
        description: 'Your education entry has been created.',
        color: 'success'
      })
    }

    closeModal()
    await fetchEducation()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to save education',
      color: 'error'
    })
  } finally {
    saving.value = false
  }
}

function openDeleteModal(edu: EducationResponse, closePopover: () => void) {
  deleteTarget.value = edu
  isDeleteModalOpen.value = true
  closePopover()
}

async function confirmDelete() {
  if (!deleteTarget.value) return

  deleting.value = true

  try {
    await apiFetch(`/education/${deleteTarget.value.id}`, {
      method: 'DELETE'
    })

    toast.add({
      title: 'Education Deleted',
      description: 'The education entry has been removed.',
      color: 'success'
    })

    isDeleteModalOpen.value = false
    deleteTarget.value = null
    await fetchEducation()
  } catch (e) {
    toast.add({
      title: 'Error',
      description: e instanceof Error ? e.message : 'Failed to delete education',
      color: 'error'
    })
  } finally {
    deleting.value = false
  }
}

// Lifecycle
onMounted(() => {
  fetchEducation()
})
</script>
