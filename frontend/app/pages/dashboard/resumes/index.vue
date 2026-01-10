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
          placeholder="All statuses"
          class="w-full sm:w-40"
        />

        <div class="ml-auto flex items-center gap-2">
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

    <!-- Table View -->
    <UCard v-if="viewMode === 'table'">
      <UTable
        :data="filteredResumes"
        :columns="columns"
      >
        <template #company-cell="{ row }">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-800">
              <UIcon
                name="i-lucide-building-2"
                class="h-5 w-5 text-zinc-400"
              />
            </div>
            <div>
              <p class="font-medium text-zinc-100">{{ row.original.company }}</p>
              <p class="text-sm text-zinc-400">{{ row.original.position }}</p>
            </div>
          </div>
        </template>

        <template #status-cell="{ row }">
          <UBadge
            :color="getStatusColor(row.original.status)"
            variant="subtle"
            size="sm"
          >
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #matchScore-cell="{ row }">
          <div class="flex items-center gap-2">
            <div class="h-2 w-20 overflow-hidden rounded-full bg-zinc-800">
              <div
                class="h-full rounded-full"
                :class="getScoreColor(row.original.matchScore)"
                :style="{ width: `${row.original.matchScore}%` }"
              />
            </div>
            <span class="text-sm text-zinc-400">{{ row.original.matchScore }}%</span>
          </div>
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
            />
            <UDropdownMenu
              :items="[
                [
                  { label: 'Duplicate', icon: 'i-lucide-copy' },
                  { label: 'Edit', icon: 'i-lucide-pencil' }
                ],
                [{ label: 'Delete', icon: 'i-lucide-trash-2', color: 'error' }]
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
      v-else
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
              {{ resume.status }}
            </UBadge>
          </div>

          <div>
            <h3 class="font-semibold text-zinc-100">{{ resume.company }}</h3>
            <p class="text-sm text-zinc-400">{{ resume.position }}</p>
          </div>

          <div class="space-y-2">
            <div class="flex items-center justify-between text-sm">
              <span class="text-zinc-400">Match Score</span>
              <span class="font-medium text-zinc-100">{{ resume.matchScore }}%</span>
            </div>
            <div class="h-2 overflow-hidden rounded-full bg-zinc-800">
              <div
                class="h-full rounded-full"
                :class="getScoreColor(resume.matchScore)"
                :style="{ width: `${resume.matchScore}%` }"
              />
            </div>
          </div>

          <div class="flex items-center justify-between text-sm text-zinc-500">
            <span>{{ resume.createdAt }}</span>
            <UButton
              color="primary"
              variant="ghost"
              size="xs"
              icon="i-lucide-download"
            />
          </div>
        </div>
      </UCard>
    </div>

    <!-- Empty State -->
    <SharedEmptyState
      v-if="filteredResumes.length === 0"
      icon="i-lucide-file-text"
      title="No resumes found"
      description="Create your first tailored resume to get started."
      action-label="Create Resume"
      action-to="/dashboard/resumes/new"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard'
})

const searchQuery = ref('')
const statusFilter = ref('')
const viewMode = ref<'table' | 'grid'>('table')

const statusOptions = [
  { label: 'All statuses', value: '' },
  { label: 'Draft', value: 'Draft' },
  { label: 'Generated', value: 'Generated' },
  { label: 'Downloaded', value: 'Downloaded' }
]

const columns = [
  { accessorKey: 'company', header: 'Company' },
  { accessorKey: 'status', header: 'Status' },
  { accessorKey: 'matchScore', header: 'Match Score' },
  { accessorKey: 'createdAt', header: 'Created' },
  { accessorKey: 'actions', header: '' }
]

// Mock resumes data.
const resumes = [
  {
    id: '1',
    company: 'Google',
    position: 'Senior Software Engineer',
    status: 'Generated',
    matchScore: 92,
    createdAt: '2 hours ago'
  },
  {
    id: '2',
    company: 'Meta',
    position: 'Frontend Engineer',
    status: 'Draft',
    matchScore: 76,
    createdAt: '1 day ago'
  },
  {
    id: '3',
    company: 'Apple',
    position: 'iOS Developer',
    status: 'Generated',
    matchScore: 88,
    createdAt: '3 days ago'
  },
  {
    id: '4',
    company: 'Microsoft',
    position: 'Full Stack Engineer',
    status: 'Downloaded',
    matchScore: 81,
    createdAt: '1 week ago'
  },
  {
    id: '5',
    company: 'Amazon',
    position: 'Backend Developer',
    status: 'Generated',
    matchScore: 85,
    createdAt: '2 weeks ago'
  },
  {
    id: '6',
    company: 'Netflix',
    position: 'Platform Engineer',
    status: 'Draft',
    matchScore: 72,
    createdAt: '3 weeks ago'
  }
]

const filteredResumes = computed(() => {
  return resumes.filter((resume) => {
    const matchesSearch =
      searchQuery.value === '' ||
      resume.company.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      resume.position.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesStatus = statusFilter.value === '' || resume.status === statusFilter.value

    return matchesSearch && matchesStatus
  })
})

function getStatusColor(status: string) {
  const colors: Record<string, 'primary' | 'warning' | 'secondary'> = {
    Generated: 'primary',
    Draft: 'warning',
    Downloaded: 'secondary'
  }
  return colors[status] || 'neutral'
}

function getScoreColor(score: number) {
  if (score >= 85) return 'bg-emerald-500'
  if (score >= 70) return 'bg-amber-500'
  return 'bg-red-500'
}
</script>
