<template>
  <div class="flex flex-col items-center justify-center py-12 text-center">
    <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-zinc-800">
      <UIcon :name="icon" class="h-8 w-8 text-zinc-500" />
    </div>
    <h3 class="mb-2 text-lg font-medium text-zinc-100">
      {{ title }}
    </h3>
    <p class="mb-6 max-w-sm text-sm text-zinc-400">
      {{ description }}
    </p>
    <UButton v-if="actionLabel" :to="actionTo" :color="actionColor" @click="$emit('action')">
      <UIcon v-if="actionIcon" :name="actionIcon" class="mr-1 h-4 w-4" />
      {{ actionLabel }}
    </UButton>
  </div>
</template>

<script setup lang="ts">
interface Props {
  icon: string
  title: string
  description: string
  actionLabel?: string
  actionTo?: string
  actionIcon?: string
  actionColor?: 'primary' | 'secondary'
}

withDefaults(defineProps<Props>(), {
  actionLabel: undefined,
  actionTo: undefined,
  actionIcon: undefined,
  actionColor: 'primary'
})

defineEmits<{
  action: []
}>()
</script>
