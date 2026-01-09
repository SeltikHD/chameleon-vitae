<template>
  <div class="flex min-h-screen bg-zinc-950">
    <!-- Sidebar -->
    <aside
      :class="[
        'fixed left-0 top-0 z-40 h-screen border-r border-zinc-800 bg-zinc-900 transition-all duration-300',
        isSidebarOpen ? 'w-60' : 'w-16'
      ]"
    >
      <!-- Logo -->
      <div class="flex h-16 items-center border-b border-zinc-800 px-4">
        <NuxtLink to="/dashboard" class="flex items-center gap-2">
          <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-500">
            <UIcon name="i-lucide-file-text" class="h-5 w-5 text-zinc-950" />
          </div>
          <span
            :class="[
              'text-lg font-bold text-zinc-100 transition-opacity duration-300',
              isSidebarOpen ? 'opacity-100' : 'opacity-0 w-0 overflow-hidden'
            ]"
          >
            Chameleon
          </span>
        </NuxtLink>
      </div>

      <!-- Navigation -->
      <nav class="flex flex-col gap-1 p-2">
        <DashboardNavItem
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          :icon="item.icon"
          :label="item.label"
          :collapsed="!isSidebarOpen"
        />
      </nav>

      <!-- Bottom Actions -->
      <div class="absolute bottom-0 left-0 right-0 border-t border-zinc-800 p-2">
        <DashboardNavItem
          to="/dashboard/profile"
          icon="i-lucide-settings"
          label="Settings"
          :collapsed="!isSidebarOpen"
        />
        <button
          class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-zinc-400 transition hover:bg-zinc-800 hover:text-zinc-100"
          @click="handleSignOut"
        >
          <UIcon name="i-lucide-log-out" class="h-5 w-5 shrink-0" />
          <span
            :class="[
              'text-sm transition-opacity duration-300',
              isSidebarOpen ? 'opacity-100' : 'opacity-0 w-0 overflow-hidden'
            ]"
          >
            Sign Out
          </span>
        </button>
      </div>

      <!-- Collapse Toggle -->
      <button
        class="absolute -right-3 top-20 flex h-6 w-6 items-center justify-center rounded-full border border-zinc-700 bg-zinc-800 text-zinc-400 transition hover:text-zinc-100"
        @click="isSidebarOpen = !isSidebarOpen"
      >
        <UIcon
          :name="isSidebarOpen ? 'i-lucide-chevron-left' : 'i-lucide-chevron-right'"
          class="h-4 w-4"
        />
      </button>
    </aside>

    <!-- Main Content -->
    <div
      :class="[
        'flex flex-1 flex-col transition-all duration-300',
        isSidebarOpen ? 'ml-60' : 'ml-16'
      ]"
    >
      <!-- Top Header -->
      <header class="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-zinc-800 bg-zinc-950/80 px-6 backdrop-blur-sm">
        <div class="flex items-center gap-4">
          <h1 class="text-xl font-semibold text-zinc-100">
            {{ pageTitle }}
          </h1>
        </div>

        <!-- User Menu -->
        <div class="flex items-center gap-4">
          <UButton to="/dashboard/resumes/new" color="primary" size="sm">
            <UIcon name="i-lucide-plus" class="mr-1 h-4 w-4" />
            New Resume
          </UButton>

          <UDropdownMenu
            :items="userMenuItems"
            :ui="{ content: 'w-48' }"
          >
            <UButton variant="ghost" color="neutral" size="sm" class="gap-2">
              <UAvatar :src="mockUser.avatar" :alt="mockUser.name" size="xs" />
              <span class="hidden text-sm md:inline">{{ mockUser.name }}</span>
              <UIcon name="i-lucide-chevron-down" class="h-4 w-4" />
            </UButton>
          </UDropdownMenu>
        </div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 p-6">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const isSidebarOpen = ref(true)

// Mock user data.
const mockUser = {
  name: 'John Developer',
  email: 'john@example.com',
  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=john'
}

// Navigation items.
const navItems = [
  { to: '/dashboard', icon: 'i-lucide-layout-dashboard', label: 'Dashboard' },
  { to: '/dashboard/resumes', icon: 'i-lucide-file-text', label: 'Resumes' },
  { to: '/dashboard/experiences', icon: 'i-lucide-briefcase', label: 'Experiences' },
  { to: '/dashboard/skills', icon: 'i-lucide-tags', label: 'Skills' },
  { to: '/dashboard/profile', icon: 'i-lucide-user', label: 'Profile' }
]

// User menu items.
const userMenuItems = [
  [
    { label: 'Profile', icon: 'i-lucide-user', to: '/dashboard/profile' },
    { label: 'Settings', icon: 'i-lucide-settings', to: '/dashboard/profile' }
  ],
  [
    { label: 'Sign Out', icon: 'i-lucide-log-out', click: () => handleSignOut() }
  ]
]

// Page title based on route.
const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/dashboard': 'Dashboard',
    '/dashboard/resumes': 'My Resumes',
    '/dashboard/resumes/new': 'Create Resume',
    '/dashboard/experiences': 'Experiences',
    '/dashboard/skills': 'Skills',
    '/dashboard/profile': 'Profile'
  }

  // Handle dynamic routes.
  if (route.path.startsWith('/dashboard/resumes/') && route.params.id) {
    return 'Resume Workbench'
  }

  return titles[route.path] || 'Dashboard'
})

function handleSignOut() {
  // Mock sign out - will be replaced with real auth.
  navigateTo('/login')
}
</script>
