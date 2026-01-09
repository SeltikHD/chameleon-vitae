<template>
  <div class="space-y-8">
    <!-- Welcome Header -->
    <div>
      <h1 class="text-2xl font-bold text-zinc-100">
        Welcome back, {{ mockUser.name.split(' ')[0] }}!
      </h1>
      <p class="mt-1 text-zinc-400">
        Here's what's happening with your resumes this week.
      </p>
    </div>

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
        <h2 class="text-lg font-semibold text-zinc-100">
          Recent Resumes
        </h2>
        <UButton to="/dashboard/resumes" color="neutral" variant="ghost" size="sm">
          View all
          <UIcon name="i-lucide-arrow-right" class="ml-1 h-4 w-4" />
        </UButton>
      </div>

      <UCard>
        <UTable :data="recentResumes" :columns="resumeColumns">
          <template #status-cell="{ row }">
            <UBadge :color="getStatusColor(row.original.status)" variant="subtle" size="sm">
              {{ row.original.status }}
            </UBadge>
          </template>

          <template #matchScore-cell="{ row }">
            <div class="flex items-center gap-2">
              <div class="h-2 w-16 overflow-hidden rounded-full bg-zinc-800">
                <div
                  class="h-full rounded-full bg-emerald-500"
                  :style="{ width: `${row.original.matchScore}%` }"
                />
              </div>
              <span class="text-sm text-zinc-400">{{ row.original.matchScore }}%</span>
            </div>
          </template>

          <template #actions-cell="{ row }">
            <UDropdownMenu
              :items="[[
                { label: 'View', icon: 'i-lucide-eye', click: () => navigateTo(`/dashboard/resumes/${row.original.id}`) },
                { label: 'Download', icon: 'i-lucide-download' },
                { label: 'Duplicate', icon: 'i-lucide-copy' }
              ], [
                { label: 'Delete', icon: 'i-lucide-trash-2', color: 'error' }
              ]]"
            >
              <UButton color="neutral" variant="ghost" icon="i-lucide-more-horizontal" size="xs" />
            </UDropdownMenu>
          </template>
        </UTable>
      </UCard>
    </div>

    <!-- Quick Actions -->
    <div>
      <h2 class="mb-4 text-lg font-semibold text-zinc-100">
        Quick Actions
      </h2>

      <div class="grid gap-4 sm:grid-cols-3">
        <UCard
          class="cursor-pointer transition-colors hover:border-emerald-500/40"
          @click="navigateTo('/dashboard/resumes/new')"
        >
          <div class="flex items-center gap-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-emerald-500/10">
              <UIcon name="i-lucide-plus" class="h-6 w-6 text-emerald-400" />
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
              <UIcon name="i-lucide-briefcase" class="h-6 w-6 text-violet-400" />
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
              <UIcon name="i-lucide-tags" class="h-6 w-6 text-sky-400" />
            </div>
            <div>
              <h3 class="font-medium text-zinc-100">Manage Skills</h3>
              <p class="text-sm text-zinc-400">Update your skill matrix</p>
            </div>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard'
})

// Mock user data.
const mockUser = {
  name: 'John Developer',
  email: 'john@example.com'
}

// Stats data.
const stats = [
  {
    title: 'Total Resumes',
    value: '12',
    icon: 'i-lucide-file-text',
    color: 'primary' as const,
    change: 3,
    changeType: 'increase' as const
  },
  {
    title: 'Experience Bullets',
    value: '47',
    icon: 'i-lucide-layers',
    color: 'secondary' as const,
    change: 8,
    changeType: 'increase' as const
  },
  {
    title: 'Avg. Match Score',
    value: '78%',
    icon: 'i-lucide-target',
    color: 'primary' as const,
    change: 5,
    changeType: 'increase' as const
  },
  {
    title: 'Downloads',
    value: '24',
    icon: 'i-lucide-download',
    color: 'primary' as const,
    change: 12,
    changeType: 'increase' as const
  }
]

// Resume columns.
const resumeColumns = [
  { accessorKey: 'company', header: 'Company' },
  { accessorKey: 'position', header: 'Position' },
  { accessorKey: 'status', header: 'Status' },
  { accessorKey: 'matchScore', header: 'Match Score' },
  { accessorKey: 'createdAt', header: 'Created' },
  { accessorKey: 'actions', header: '' }
]

// Mock recent resumes.
const recentResumes = [
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
  }
]

function getStatusColor(status: string) {
  const colors: Record<string, 'primary' | 'warning' | 'secondary'> = {
    Generated: 'primary',
    Draft: 'warning',
    Downloaded: 'secondary'
  }
  return colors[status] || 'neutral'
}
</script>
