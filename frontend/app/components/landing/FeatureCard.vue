<template>
  <UCard
    class="text-center transition-transform duration-300 hover:-translate-y-1"
    :class="cardClasses"
  >
    <div class="flex flex-col items-center">
      <div
        class="mb-4 flex h-14 w-14 items-center justify-center rounded-xl"
        :class="iconBgClasses"
      >
        <UIcon
          :name="icon"
          class="h-7 w-7"
          :class="iconColorClasses"
        />
      </div>

      <h3 class="text-lg font-semibold text-zinc-100">
        {{ title }}
      </h3>

      <p class="mt-2 text-sm text-zinc-400">
        {{ description }}
      </p>
    </div>
  </UCard>
</template>

<script setup lang="ts">
interface Props {
  icon: string
  title: string
  description: string
  color?: 'primary' | 'secondary'
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary'
})

const cardClasses = computed(() => ({
  'border-emerald-500/20 hover:border-emerald-500/40': props.color === 'primary',
  'border-violet-500/20 hover:border-violet-500/40': props.color === 'secondary'
}))

const iconBgClasses = computed(() => ({
  'bg-emerald-500/10': props.color === 'primary',
  'bg-violet-500/10': props.color === 'secondary'
}))

const iconColorClasses = computed(() => ({
  'text-emerald-400': props.color === 'primary',
  'text-violet-400': props.color === 'secondary'
}))
</script>
