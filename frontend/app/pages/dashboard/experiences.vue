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
        @click="showAddModal = true"
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
        title="Skills Covered"
        :value="String(uniqueSkills.length)"
        icon="i-lucide-tags"
        color="primary"
      />
    </div>

    <!-- Timeline View -->
    <div class="relative">
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
              name="i-lucide-briefcase"
              class="h-5 w-5 text-emerald-400"
            />
          </div>

          <UCard>
            <template #header>
              <div class="flex items-start justify-between">
                <div>
                  <h3 class="font-semibold text-zinc-100">
                    {{ experience.role }}
                  </h3>
                  <p class="text-sm text-zinc-400">{{ experience.company }}</p>
                </div>
                <div class="text-right">
                  <UBadge
                    color="neutral"
                    variant="subtle"
                  >
                    {{ experience.period }}
                  </UBadge>
                  <p class="mt-1 text-xs text-zinc-500">{{ experience.bullets.length }} bullets</p>
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
                    <p class="text-sm text-zinc-300">{{ bullet.text }}</p>
                    <div class="mt-2 flex flex-wrap gap-1">
                      <UBadge
                        v-for="tag in bullet.tags"
                        :key="tag"
                        color="primary"
                        variant="subtle"
                        size="xs"
                      >
                        {{ tag }}
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
                    />
                    <UButton
                      color="neutral"
                      variant="ghost"
                      icon="i-lucide-copy"
                      size="xs"
                    />
                    <UButton
                      color="neutral"
                      variant="ghost"
                      icon="i-lucide-trash-2"
                      size="xs"
                    />
                  </div>
                </div>
              </div>

              <UButton
                color="neutral"
                variant="dashed"
                block
                size="sm"
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
                  <span>Used in {{ experience.usedInResumes }} resumes</span>
                </div>
                <div class="flex items-center gap-2">
                  <UButton
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    icon="i-lucide-pencil"
                  >
                    Edit
                  </UButton>
                  <UButton
                    color="error"
                    variant="ghost"
                    size="xs"
                    icon="i-lucide-trash-2"
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

    <!-- Empty State -->
    <SharedEmptyState
      v-if="experiences.length === 0"
      icon="i-lucide-briefcase"
      title="No experiences yet"
      description="Add your work experiences to build your bullet library."
      action-label="Add Experience"
      @action="showAddModal = true"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard'
})

const showAddModal = ref(false)

// Mock experiences data.
const experiences = [
  {
    id: '1',
    role: 'Senior Frontend Engineer',
    company: 'TechCorp Inc.',
    period: '2021 - Present',
    usedInResumes: 8,
    bullets: [
      {
        id: '1',
        text: 'Led development of a React-based design system used by 50+ engineers, reducing UI development time by 40%.',
        tags: ['React', 'Design System', 'Leadership']
      },
      {
        id: '2',
        text: 'Architected and implemented real-time collaboration features using WebSockets, serving 100K+ daily active users.',
        tags: ['WebSockets', 'Architecture', 'Scale']
      },
      {
        id: '3',
        text: 'Mentored team of 5 junior developers, conducting code reviews and pair programming sessions.',
        tags: ['Mentoring', 'Code Review', 'Leadership']
      }
    ]
  },
  {
    id: '2',
    role: 'Full Stack Developer',
    company: 'StartupXYZ',
    period: '2019 - 2021',
    usedInResumes: 6,
    bullets: [
      {
        id: '4',
        text: 'Built and deployed microservices architecture on AWS, reducing infrastructure costs by 35%.',
        tags: ['AWS', 'Microservices', 'Cost Optimization']
      },
      {
        id: '5',
        text: 'Implemented comprehensive testing strategy (unit, integration, e2e) achieving 95% code coverage.',
        tags: ['Testing', 'Quality', 'CI/CD']
      }
    ]
  },
  {
    id: '3',
    role: 'Junior Developer',
    company: 'WebAgency Co.',
    period: '2017 - 2019',
    usedInResumes: 3,
    bullets: [
      {
        id: '6',
        text: 'Developed responsive web applications for 20+ clients using React and Node.js.',
        tags: ['React', 'Node.js', 'Client Work']
      },
      {
        id: '7',
        text: 'Improved site performance by 60% through code optimization and lazy loading implementation.',
        tags: ['Performance', 'Optimization']
      }
    ]
  }
]

const totalBullets = computed(() => experiences.reduce((acc, exp) => acc + exp.bullets.length, 0))

const uniqueSkills = computed(() => {
  const skills = new Set<string>()
  experiences.forEach((exp) => {
    exp.bullets.forEach((bullet) => {
      bullet.tags.forEach((tag) => skills.add(tag))
    })
  })
  return Array.from(skills)
})
</script>
