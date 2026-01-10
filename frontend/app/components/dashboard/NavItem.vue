<template>
  <NuxtLink
    :to="to"
    :class="[
      'flex items-center gap-3 rounded-lg px-3 py-2 transition',
      isActive
        ? 'bg-emerald-500/10 text-emerald-500'
        : 'text-zinc-400 hover:bg-zinc-800 hover:text-zinc-100'
    ]"
  >
    <UIcon
      :name="icon"
      class="h-5 w-5 shrink-0"
    />
    <span
      :class="[
        'text-sm font-medium transition-opacity duration-300',
        collapsed ? 'opacity-0 w-0 overflow-hidden' : 'opacity-100'
      ]"
    >
      {{ label }}
    </span>
  </NuxtLink>
</template>

<script setup lang="ts">
interface Props {
  to: string
  icon: string
  label: string
  collapsed?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  collapsed: false
})

const route = useRoute()

const isActive = computed(() => {
  // Exact match for dashboard home.
  if (props.to === '/dashboard') {
    return route.path === '/dashboard'
  }
  // Prefix match for other routes.
  return route.path.startsWith(props.to)
})
</script>
