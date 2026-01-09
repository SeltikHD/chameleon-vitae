<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-zinc-100">
          Skills
        </h1>
        <p class="mt-1 text-zinc-400">
          Manage your skill matrix and see how they're used across experiences.
        </p>
      </div>
      <UButton color="primary" @click="showAddModal = true">
        <UIcon name="i-lucide-plus" class="mr-2 h-4 w-4" />
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
          placeholder="All categories"
          class="w-full sm:w-48"
        />
      </div>
    </UCard>

    <!-- Skills by Category -->
    <div class="space-y-6">
      <div v-for="category in filteredCategories" :key="category.name">
        <div class="mb-4 flex items-center gap-2">
          <UIcon :name="category.icon" class="h-5 w-5" :class="category.color" />
          <h2 class="text-lg font-semibold text-zinc-100">{{ category.name }}</h2>
          <UBadge color="neutral" variant="subtle" size="sm">
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
                  <UBadge :color="getLevelColor(skill.level)" variant="subtle" size="xs">
                    {{ skill.level }}
                  </UBadge>
                </div>
                <p class="mt-1 text-sm text-zinc-400">
                  Used in {{ skill.bulletCount }} bullets
                </p>

                <!-- Proficiency Bar -->
                <div class="mt-3">
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-zinc-500">Proficiency</span>
                    <span class="text-zinc-400">{{ skill.proficiency }}%</span>
                  </div>
                  <div class="mt-1 h-1.5 overflow-hidden rounded-full bg-zinc-800">
                    <div
                      class="h-full rounded-full bg-emerald-500"
                      :style="{ width: `${skill.proficiency}%` }"
                    />
                  </div>
                </div>

                <!-- Years of Experience -->
                <div class="mt-2 flex items-center gap-1 text-xs text-zinc-500">
                  <UIcon name="i-lucide-calendar" class="h-3 w-3" />
                  <span>{{ skill.yearsOfExperience }} years</span>
                </div>
              </div>

              <div class="flex flex-col items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100">
                <UButton color="neutral" variant="ghost" icon="i-lucide-pencil" size="xs" />
                <UButton color="neutral" variant="ghost" icon="i-lucide-trash-2" size="xs" />
              </div>
            </div>
          </UCard>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <SharedEmptyState
      v-if="filteredCategories.length === 0"
      icon="i-lucide-tags"
      title="No skills found"
      description="Add skills to build your expertise matrix."
      action-label="Add Skill"
      @action="showAddModal = true"
    />

    <!-- Skill Matrix Summary -->
    <UCard>
      <template #header>
        <h2 class="font-semibold text-zinc-100">Skill Matrix Summary</h2>
      </template>

      <div class="grid gap-6 md:grid-cols-3">
        <div class="text-center">
          <p class="text-3xl font-bold text-emerald-400">{{ totalSkills }}</p>
          <p class="text-sm text-zinc-400">Total Skills</p>
        </div>
        <div class="text-center">
          <p class="text-3xl font-bold text-violet-400">{{ expertSkills }}</p>
          <p class="text-sm text-zinc-400">Expert Level</p>
        </div>
        <div class="text-center">
          <p class="text-3xl font-bold text-sky-400">{{ avgProficiency }}%</p>
          <p class="text-sm text-zinc-400">Avg. Proficiency</p>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard'
})

const showAddModal = ref(false)
const searchQuery = ref('')
const categoryFilter = ref('')

const categoryOptions = [
  { label: 'All categories', value: '' },
  { label: 'Frontend', value: 'Frontend' },
  { label: 'Backend', value: 'Backend' },
  { label: 'DevOps', value: 'DevOps' },
  { label: 'Soft Skills', value: 'Soft Skills' }
]

// Mock skills data organized by category.
const categories = [
  {
    name: 'Frontend',
    icon: 'i-lucide-layout',
    color: 'text-emerald-400',
    skills: [
      { id: '1', name: 'React', level: 'Expert', proficiency: 95, yearsOfExperience: 5, bulletCount: 12 },
      { id: '2', name: 'TypeScript', level: 'Expert', proficiency: 90, yearsOfExperience: 4, bulletCount: 10 },
      { id: '3', name: 'Vue.js', level: 'Advanced', proficiency: 80, yearsOfExperience: 3, bulletCount: 6 },
      { id: '4', name: 'Tailwind CSS', level: 'Expert', proficiency: 92, yearsOfExperience: 3, bulletCount: 8 },
      { id: '5', name: 'Next.js', level: 'Advanced', proficiency: 85, yearsOfExperience: 2, bulletCount: 5 }
    ]
  },
  {
    name: 'Backend',
    icon: 'i-lucide-server',
    color: 'text-violet-400',
    skills: [
      { id: '6', name: 'Node.js', level: 'Advanced', proficiency: 85, yearsOfExperience: 5, bulletCount: 8 },
      { id: '7', name: 'Go', level: 'Intermediate', proficiency: 65, yearsOfExperience: 1, bulletCount: 3 },
      { id: '8', name: 'PostgreSQL', level: 'Advanced', proficiency: 80, yearsOfExperience: 4, bulletCount: 6 },
      { id: '9', name: 'GraphQL', level: 'Advanced', proficiency: 78, yearsOfExperience: 3, bulletCount: 4 }
    ]
  },
  {
    name: 'DevOps',
    icon: 'i-lucide-cloud',
    color: 'text-sky-400',
    skills: [
      { id: '10', name: 'Docker', level: 'Advanced', proficiency: 82, yearsOfExperience: 4, bulletCount: 5 },
      { id: '11', name: 'AWS', level: 'Advanced', proficiency: 75, yearsOfExperience: 3, bulletCount: 4 },
      { id: '12', name: 'CI/CD', level: 'Expert', proficiency: 88, yearsOfExperience: 4, bulletCount: 6 }
    ]
  },
  {
    name: 'Soft Skills',
    icon: 'i-lucide-users',
    color: 'text-amber-400',
    skills: [
      { id: '13', name: 'Leadership', level: 'Advanced', proficiency: 85, yearsOfExperience: 3, bulletCount: 4 },
      { id: '14', name: 'Communication', level: 'Expert', proficiency: 90, yearsOfExperience: 7, bulletCount: 3 },
      { id: '15', name: 'Problem Solving', level: 'Expert', proficiency: 92, yearsOfExperience: 7, bulletCount: 5 }
    ]
  }
]

const filteredCategories = computed(() => {
  return categories
    .filter(category =>
      categoryFilter.value === '' || category.name === categoryFilter.value
    )
    .map(category => ({
      ...category,
      skills: category.skills.filter(skill =>
        searchQuery.value === ''
        || skill.name.toLowerCase().includes(searchQuery.value.toLowerCase())
      )
    }))
    .filter(category => category.skills.length > 0)
})

const totalSkills = computed(() =>
  categories.reduce((acc, cat) => acc + cat.skills.length, 0)
)

const expertSkills = computed(() =>
  categories.reduce((acc, cat) =>
    acc + cat.skills.filter(s => s.level === 'Expert').length, 0
  )
)

const avgProficiency = computed(() => {
  const all = categories.flatMap(c => c.skills)
  const sum = all.reduce((acc, s) => acc + s.proficiency, 0)
  return Math.round(sum / all.length)
})

function getLevelColor(level: string) {
  const colors: Record<string, 'primary' | 'secondary' | 'warning'> = {
    Expert: 'primary',
    Advanced: 'secondary',
    Intermediate: 'warning'
  }
  return colors[level] || 'neutral'
}
</script>
