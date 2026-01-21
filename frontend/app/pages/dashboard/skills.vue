<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-zinc-100">Skills</h1>
        <p class="mt-1 text-zinc-400">Manage your skill matrix and proficiency levels.</p>
      </div>
      <UButton
        color="primary"
        @click="openAddModal"
      >
        <UIcon
          name="i-lucide-plus"
          class="mr-2 h-4 w-4"
        />
        Add Skill
      </UButton>
    </div>

    <!-- Search & Filter -->
    <UCard>
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center">
        <UInput
          v-model="searchQuery"
          placeholder="Search skills..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
        />

        <USelectMenu
          v-model="categoryFilter"
          :items="categoryOptions"
          value-key="value"
          placeholder="All categories"
          class="w-full sm:w-48"
        />
      </div>
    </UCard>

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
          <p class="font-medium text-zinc-100">Failed to load skills</p>
          <p class="mt-1 text-sm text-zinc-400">{{ error }}</p>
        </div>
        <UButton
          color="primary"
          @click="fetchSkills()"
        >
          Try Again
        </UButton>
      </div>
    </UCard>

    <!-- Empty State -->
    <SharedEmptyState
      v-else-if="skills.length === 0"
      icon="i-lucide-tags"
      title="No skills yet"
      description="Add your skills to build your expertise matrix."
      action-label="Add Skill"
      @action="openAddModal"
    />

    <!-- Skills Grid -->
    <template v-else>
      <!-- Skills by Category -->
      <div class="space-y-6">
        <div
          v-for="category in filteredCategories"
          :key="category.name"
        >
          <div class="mb-4 flex items-center gap-2">
            <UIcon
              :name="getCategoryIcon(category.name)"
              class="h-5 w-5"
              :class="getCategoryColor(category.name)"
            />
            <h2 class="text-lg font-semibold text-zinc-100">{{ category.name }}</h2>
            <UBadge
              color="neutral"
              variant="subtle"
              size="sm"
            >
              {{ category.skills.length }}
            </UBadge>
          </div>

          <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <UCard
              v-for="skill in category.skills"
              :key="skill.id"
              class="group cursor-pointer transition-colors hover:border-emerald-500/40"
            >
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center gap-2">
                    <h3 class="font-medium text-zinc-100">{{ skill.name }}</h3>
                    <UBadge
                      :color="getSkillLevelColor(skill.proficiency_level)"
                      variant="subtle"
                      size="xs"
                    >
                      {{ getSkillLevel(skill.proficiency_level) }}
                    </UBadge>
                    <UIcon
                      v-if="skill.is_highlighted"
                      name="i-lucide-star"
                      class="h-4 w-4 text-amber-400"
                    />
                  </div>

                  <!-- Proficiency Bar -->
                  <div class="mt-3">
                    <div class="flex items-center justify-between text-xs">
                      <span class="text-zinc-500">Proficiency</span>
                      <span class="text-zinc-400">{{ skill.proficiency_level }}%</span>
                    </div>
                    <div class="mt-1 h-1.5 overflow-hidden rounded-full bg-zinc-800">
                      <div
                        class="h-full rounded-full bg-emerald-500"
                        :style="{ width: `${skill.proficiency_level}%` }"
                      />
                    </div>
                  </div>

                  <!-- Years of Experience -->
                  <div
                    v-if="skill.years_of_experience"
                    class="mt-2 flex items-center gap-1 text-xs text-zinc-500"
                  >
                    <UIcon
                      name="i-lucide-calendar"
                      class="h-3 w-3"
                    />
                    <span>{{ skill.years_of_experience }} years</span>
                  </div>
                </div>

                <div
                  class="flex flex-col items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100"
                >
                  <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-lucide-pencil"
                    size="xs"
                    @click="openEditModal(skill)"
                  />
                  <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-lucide-trash-2"
                    size="xs"
                    @click="deleteSkill(skill.id)"
                  />
                </div>
              </div>
            </UCard>
          </div>
        </div>
      </div>

      <!-- Empty filtered state -->
      <SharedEmptyState
        v-if="filteredCategories.length === 0 && skills.length > 0"
        icon="i-lucide-search"
        title="No skills found"
        description="Try adjusting your search or filter."
      />

      <!-- Skill Matrix Summary -->
      <UCard>
        <template #header>
          <h2 class="font-semibold text-zinc-100">Skill Matrix Summary</h2>
        </template>

        <div class="grid gap-6 md:grid-cols-3">
          <div class="text-center">
            <p class="text-3xl font-bold text-emerald-400">{{ skills.length }}</p>
            <p class="text-sm text-zinc-400">Total Skills</p>
          </div>
          <div class="text-center">
            <p class="text-3xl font-bold text-violet-400">{{ expertSkillsCount }}</p>
            <p class="text-sm text-zinc-400">Expert Level</p>
          </div>
          <div class="text-center">
            <p class="text-3xl font-bold text-sky-400">{{ avgProficiency }}%</p>
            <p class="text-sm text-zinc-400">Avg. Proficiency</p>
          </div>
        </div>
      </UCard>
    </template>

    <!-- Add/Edit Skill Modal -->
    <UModal v-model:open="showSkillModal">
      <template #content>
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-zinc-100">
                {{ editingSkill ? 'Edit Skill' : 'Add Skill' }}
              </h3>
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                size="xs"
                @click="showSkillModal = false"
              />
            </div>
          </template>

          <form
            class="space-y-4"
            @submit.prevent="saveSkill"
          >
            <UFormField
              label="Skill Name"
              name="name"
              required
            >
              <UInput
                v-model="skillForm.name"
                placeholder="React, Python, Leadership..."
              />
            </UFormField>

            <UFormField
              label="Category"
              name="category"
              hint="Select or type a new category"
            >
              <UInputMenu
                v-model="categoryModel"
                v-model:query="categoryQuery"
                :items="skillCategoryOptions"
                class="w-full"
                placeholder="Select or type a category..."
                creatable
                searchable
                @keydown.enter.prevent.stop="handleManualCreate"
              >
                <template #empty="{ searchTerm }">
                  <div
                    v-if="searchTerm"
                    class="px-3 py-2 text-sm text-zinc-400"
                  >
                    Press Enter to create "<span class="text-emerald-400 font-medium">{{
                      searchTerm
                    }}</span
                    >"
                  </div>
                  <div
                    v-else
                    class="px-3 py-2 text-sm text-zinc-400"
                  >
                    Start typing to create a new category
                  </div>
                </template>
              </UInputMenu>
            </UFormField>

            <UFormField
              label="Proficiency Level"
              name="proficiency_level"
            >
              <div class="flex items-center gap-4">
                <UInput
                  v-model.number="skillForm.proficiency_level"
                  type="range"
                  min="0"
                  max="100"
                  class="flex-1"
                />
                <span class="w-12 text-right text-sm text-zinc-400"
                  >{{ skillForm.proficiency_level }}%</span
                >
              </div>
            </UFormField>

            <UFormField
              label="Years of Experience"
              name="years_of_experience"
            >
              <UInput
                v-model.number="skillForm.years_of_experience"
                type="number"
                min="0"
                step="0.5"
                placeholder="3.5"
              />
            </UFormField>

            <UCheckbox
              v-model="skillForm.is_highlighted"
              label="Highlight this skill"
            />

            <div class="flex justify-end gap-3 pt-4">
              <UButton
                color="neutral"
                variant="ghost"
                @click="showSkillModal = false"
              >
                Cancel
              </UButton>
              <UButton
                type="submit"
                color="primary"
                :loading="isSaving"
              >
                {{ editingSkill ? 'Save Changes' : 'Add Skill' }}
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
import type { SkillResponse, SkillInput, ListSkillsResponse } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const toast = useToast()

// State
const loading = ref(true)
const error = ref<string | null>(null)
const isSaving = ref(false)
const skills = ref<SkillResponse[]>([])

// Filters
const searchQuery = ref('')
const categoryFilter = ref('')

// Modal
const showSkillModal = ref(false)
const editingSkill = ref<SkillResponse | null>(null)

// Form
const skillForm = reactive<SkillInput & { is_highlighted: boolean }>({
  name: '',
  category: null,
  proficiency_level: 50,
  years_of_experience: null,
  is_highlighted: false
})

// Query for UInputMenu creatable functionality
const categoryQuery = ref('')

// Computed wrapper for UInputMenu - handles object value for creatable menu
const categoryModel = computed({
  get: (): { label: string; value: string } | undefined => {
    if (!skillForm.category) return undefined
    return { label: skillForm.category, value: skillForm.category }
  },
  set: (value: string | { label: string; value: string } | undefined) => {
    // Handle both string (from creatable) and object (from selection)
    if (typeof value === 'string') {
      skillForm.category = value || null
    } else if (value && typeof value === 'object') {
      skillForm.category = value.value || null
    } else {
      skillForm.category = null
    }
  }
})

// Default category options
const defaultCategories = [
  'Frontend',
  'Backend',
  'DevOps',
  'Database',
  'Mobile',
  'AI/ML',
  'Soft Skills',
  'Tools',
  'Other'
]

// Function to handle manual creation of category on Enter key
const handleManualCreate = (e: KeyboardEvent) => {
  // Cast target to HTMLInputElement to access value
  const target = e.target as HTMLInputElement

  // Get the value directly from the DOM, ignoring Vue's reactivity which might be reset
  const newCategory = target.value

  if (newCategory && newCategory.trim() !== '') {
    // Set the category in the form
    skillForm.category = newCategory.trim()

    // Update the query to visually reflect (even though the modal might close or the input might clear)
    categoryQuery.value = newCategory.trim()
  }
}

// Mantenha seu categoryModel computed como estava, ele já trata string corretamente.

// Dynamic category options - combines defaults with user's custom categories
const skillCategoryOptions = computed(() => {
  // Collect all unique categories from existing skills
  const existingCategories = new Set(
    skills.value
      .map((s) => s.category)
      .filter((c): c is string => c !== null && c !== undefined && c !== '')
  )

  // Merge with defaults, keeping unique values
  const allCategories = new Set([...defaultCategories, ...existingCategories])

  return Array.from(allCategories)
    .sort((a, b) => {
      // Keep defaults at top, then alphabetical
      const aIsDefault = defaultCategories.includes(a)
      const bIsDefault = defaultCategories.includes(b)
      if (aIsDefault && !bIsDefault) return -1
      if (!aIsDefault && bIsDefault) return 1
      return a.localeCompare(b)
    })
    .map((cat) => ({ label: cat, value: cat }))
})

const categoryOptions = computed(() => [
  { label: 'All categories', value: '' },
  ...skillCategoryOptions.value
])

// Computed
const categories = computed(() => {
  const categoryMap = new Map<string, SkillResponse[]>()

  skills.value.forEach((skill) => {
    const cat = skill.category || 'Other'
    if (!categoryMap.has(cat)) {
      categoryMap.set(cat, [])
    }
    categoryMap.get(cat)!.push(skill)
  })

  return Array.from(categoryMap.entries())
    .map(([name, skills]) => ({ name, skills }))
    .sort((a, b) => a.name.localeCompare(b.name))
})

const filteredCategories = computed(() => {
  return categories.value
    .filter((cat) => categoryFilter.value === '' || cat.name === categoryFilter.value)
    .map((cat) => ({
      ...cat,
      skills: cat.skills.filter(
        (skill) =>
          searchQuery.value === '' ||
          skill.name.toLowerCase().includes(searchQuery.value.toLowerCase())
      )
    }))
    .filter((cat) => cat.skills.length > 0)
})

const expertSkillsCount = computed(
  () => skills.value.filter((s) => s.proficiency_level >= 80).length
)

const avgProficiency = computed(() => {
  if (skills.value.length === 0) return 0
  const sum = skills.value.reduce((acc, s) => acc + s.proficiency_level, 0)
  return Math.round(sum / skills.value.length)
})

// Fetch skills
async function fetchSkills() {
  loading.value = true
  error.value = null

  try {
    const response = await apiFetch<ListSkillsResponse>('/skills')
    skills.value = response.data
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load skills'
  } finally {
    loading.value = false
  }
}

// Skill CRUD
function openAddModal() {
  editingSkill.value = null
  Object.assign(skillForm, {
    name: '',
    category: null,
    proficiency_level: 50,
    years_of_experience: null,
    is_highlighted: false
  })
  showSkillModal.value = true
}

function openEditModal(skill: SkillResponse) {
  editingSkill.value = skill
  Object.assign(skillForm, {
    name: skill.name,
    category: skill.category,
    proficiency_level: skill.proficiency_level,
    years_of_experience: skill.years_of_experience,
    is_highlighted: skill.is_highlighted
  })
  showSkillModal.value = true
}

async function saveSkill() {
  isSaving.value = true

  try {
    // Use batch upsert API for both create and update.
    const skillData: SkillInput = {
      name: skillForm.name,
      category: skillForm.category || null,
      proficiency_level: skillForm.proficiency_level,
      years_of_experience: skillForm.years_of_experience || null,
      is_highlighted: skillForm.is_highlighted
    }

    await apiFetch('/skills/batch', {
      method: 'POST',
      body: { skills: [skillData] }
    })

    toast.add({
      title: editingSkill.value ? 'Skill Updated' : 'Skill Added',
      description: `${skillForm.name} has been saved.`,
      color: 'success'
    })

    showSkillModal.value = false
    await fetchSkills()
  } catch (e) {
    console.error('Failed to save skill:', e)
  } finally {
    isSaving.value = false
  }
}

async function deleteSkill(id: string) {
  if (!confirm('Are you sure you want to delete this skill?')) return

  try {
    await apiFetch(`/skills/${id}`, { method: 'DELETE' })

    toast.add({
      title: 'Skill Deleted',
      description: 'The skill has been removed.',
      color: 'success'
    })

    await fetchSkills()
  } catch (e) {
    console.error('Failed to delete skill:', e)
  }
}

// Helpers
function getCategoryIcon(category: string): string {
  const icons: Record<string, string> = {
    Frontend: 'i-lucide-layout',
    Backend: 'i-lucide-server',
    DevOps: 'i-lucide-cloud',
    Database: 'i-lucide-database',
    Mobile: 'i-lucide-smartphone',
    'AI/ML': 'i-lucide-brain',
    'Soft Skills': 'i-lucide-users',
    Tools: 'i-lucide-wrench',
    Other: 'i-lucide-tag'
  }
  return icons[category] || 'i-lucide-tag'
}

function getCategoryColor(category: string): string {
  const colors: Record<string, string> = {
    Frontend: 'text-emerald-400',
    Backend: 'text-violet-400',
    DevOps: 'text-sky-400',
    Database: 'text-orange-400',
    Mobile: 'text-pink-400',
    'AI/ML': 'text-purple-400',
    'Soft Skills': 'text-amber-400',
    Tools: 'text-cyan-400',
    Other: 'text-zinc-400'
  }
  return colors[category] || 'text-zinc-400'
}

function getSkillLevel(proficiency: number): string {
  if (proficiency >= 80) return 'Expert'
  if (proficiency >= 60) return 'Advanced'
  if (proficiency >= 40) return 'Intermediate'
  return 'Beginner'
}

function getSkillLevelColor(proficiency: number): 'primary' | 'secondary' | 'warning' | 'neutral' {
  if (proficiency >= 80) return 'primary'
  if (proficiency >= 60) return 'secondary'
  if (proficiency >= 40) return 'warning'
  return 'neutral'
}

// Fetch on mount
onMounted(() => {
  fetchSkills()
})
</script>
