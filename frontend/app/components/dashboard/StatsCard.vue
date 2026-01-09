<template>
  <UCard :ui="{ body: 'p-4' }">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-sm font-medium text-zinc-400">
          {{ title }}
        </p>
        <p class="mt-1 text-2xl font-bold text-zinc-100">
          {{ formattedValue }}
        </p>
        <p v-if="change" class="mt-1 flex items-center gap-1 text-xs">
          <UIcon
            :name="changeIcon"
            :class="['h-3 w-3', changeColor]"
          />
          <span :class="changeColor">{{ change }}%</span>
          <span class="text-zinc-500">vs last month</span>
        </p>
      </div>
      <div
        :class="[
          'flex h-10 w-10 items-center justify-center rounded-lg',
          iconBackground
        ]"
      >
        <UIcon :name="icon" :class="['h-5 w-5', iconColor]" />
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
interface Props {
  title: string
  value: number | string
  icon: string
  color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error'
  change?: number
  suffix?: string
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary',
  change: undefined,
  suffix: ''
})

const formattedValue = computed(() => {
  if (typeof props.value === 'number') {
    return props.value.toLocaleString() + props.suffix
  }
  return props.value + props.suffix
})

const iconBackground = computed(() => {
  const backgrounds: Record<string, string> = {
    primary: 'bg-emerald-500/10',
    secondary: 'bg-violet-500/10',
    success: 'bg-emerald-500/10',
    warning: 'bg-amber-500/10',
    error: 'bg-red-500/10'
  }
  return backgrounds[props.color]
})

const iconColor = computed(() => {
  const colors: Record<string, string> = {
    primary: 'text-emerald-500',
    secondary: 'text-violet-500',
    success: 'text-emerald-500',
    warning: 'text-amber-500',
    error: 'text-red-500'
  }
  return colors[props.color]
})

const changeIcon = computed(() => {
  if (!props.change) return ''
  return props.change > 0 ? 'i-lucide-trending-up' : 'i-lucide-trending-down'
})

const changeColor = computed(() => {
  if (!props.change) return ''
  return props.change > 0 ? 'text-emerald-500' : 'text-red-500'
})
</script>
