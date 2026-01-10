<template>
  <div class="space-y-6">
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
          {{ resume.company }} - {{ resume.position }}
        </h1>
        <p class="mt-1 text-zinc-400">Created {{ resume.createdAt }}</p>
      </div>
      <div class="flex items-center gap-2">
        <UBadge
          :color="getStatusColor(resume.status)"
          variant="subtle"
          size="lg"
        >
          {{ resume.status }}
        </UBadge>
        <UButton
          color="neutral"
          variant="outline"
          icon="i-lucide-download"
        >
          Download PDF
        </UButton>
        <UButton
          color="primary"
          icon="i-lucide-refresh-cw"
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
            <p class="text-3xl font-bold text-emerald-400">{{ resume.matchScore }}%</p>
          </div>
        </div>
        <div class="text-right">
          <p class="text-sm text-zinc-400">{{ resume.bulletCount }} bullets selected</p>
          <p class="text-sm text-zinc-400">{{ resume.skillsMatched }} skills matched</p>
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
            >
              Add Bullet
            </UButton>
          </div>
        </template>

        <div class="space-y-4">
          <div
            v-for="experience in resume.experiences"
            :key="experience.id"
            class="group"
          >
            <div class="mb-2 flex items-center justify-between">
              <span class="text-sm font-medium text-zinc-400">{{ experience.company }}</span>
              <span class="text-xs text-zinc-500">{{ experience.period }}</span>
            </div>

            <div
              v-for="bullet in experience.bullets"
              :key="bullet.id"
              class="relative mb-3 rounded-lg border border-zinc-800 bg-zinc-900 p-4 transition-colors hover:border-zinc-700"
            >
              <div class="flex items-start gap-3">
                <UIcon
                  name="i-lucide-grip-vertical"
                  class="mt-1 h-4 w-4 cursor-grab text-zinc-600"
                />
                <div class="flex-1">
                  <p class="text-sm text-zinc-300">{{ bullet.text }}</p>
                  <div class="mt-2 flex flex-wrap gap-2">
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
                    icon="i-lucide-trash-2"
                    size="xs"
                  />
                </div>
              </div>
            </div>
          </div>
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
              />
              <span class="text-sm text-zinc-400">100%</span>
              <UButton
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-zoom-in"
              />
            </div>
          </div>
        </template>

        <!-- Mock PDF Preview -->
        <div class="aspect-[8.5/11] overflow-hidden rounded-lg border border-zinc-800 bg-white p-8">
          <div class="space-y-4">
            <div class="border-b-2 border-zinc-900 pb-4">
              <h3 class="text-xl font-bold text-zinc-900">John Developer</h3>
              <p class="text-sm text-zinc-600">
                john@example.com • +1 (555) 123-4567 • San Francisco, CA
              </p>
              <p class="text-sm text-zinc-600">
                linkedin.com/in/johndeveloper • github.com/johndeveloper
              </p>
            </div>

            <div>
              <h4 class="mb-2 font-semibold text-zinc-900">Summary</h4>
              <p class="text-xs text-zinc-700">
                Senior Software Engineer with 7+ years of experience building scalable web
                applications. Expert in React, TypeScript, and Node.js with a passion for clean code
                and user experience.
              </p>
            </div>

            <div>
              <h4 class="mb-2 font-semibold text-zinc-900">Experience</h4>
              <div class="space-y-3">
                <div
                  v-for="exp in resume.experiences"
                  :key="exp.id"
                >
                  <div class="flex items-baseline justify-between">
                    <span class="text-sm font-medium text-zinc-900"
                      >{{ exp.role }} @ {{ exp.company }}</span
                    >
                    <span class="text-xs text-zinc-500">{{ exp.period }}</span>
                  </div>
                  <ul class="ml-4 mt-1 list-disc space-y-1">
                    <li
                      v-for="bullet in exp.bullets"
                      :key="bullet.id"
                      class="text-xs text-zinc-700"
                    >
                      {{ bullet.text }}
                    </li>
                  </ul>
                </div>
              </div>
            </div>

            <div>
              <h4 class="mb-2 font-semibold text-zinc-900">Skills</h4>
              <p class="text-xs text-zinc-700">
                React, TypeScript, Next.js, Node.js, PostgreSQL, GraphQL, AWS, Docker, Kubernetes,
                CI/CD
              </p>
            </div>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard'
})

// Mock resume data.
const resume = reactive({
  id: '1',
  company: 'Google',
  position: 'Senior Software Engineer',
  status: 'Generated',
  matchScore: 92,
  bulletCount: 8,
  skillsMatched: 12,
  createdAt: '2 hours ago',
  experiences: [
    {
      id: '1',
      company: 'TechCorp Inc.',
      role: 'Senior Frontend Engineer',
      period: '2021 - Present',
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
        }
      ]
    },
    {
      id: '2',
      company: 'StartupXYZ',
      role: 'Full Stack Developer',
      period: '2019 - 2021',
      bullets: [
        {
          id: '3',
          text: 'Built and deployed microservices architecture on AWS, reducing infrastructure costs by 35%.',
          tags: ['AWS', 'Microservices', 'Cost Optimization']
        },
        {
          id: '4',
          text: 'Implemented comprehensive testing strategy (unit, integration, e2e) achieving 95% code coverage.',
          tags: ['Testing', 'Quality', 'CI/CD']
        }
      ]
    }
  ]
})

function getStatusColor(status: string) {
  const colors: Record<string, 'primary' | 'warning' | 'secondary'> = {
    Generated: 'primary',
    Draft: 'warning',
    Downloaded: 'secondary'
  }
  return colors[status] || 'neutral'
}
</script>
