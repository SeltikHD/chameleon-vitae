<template>
  <NuxtLink
    :to="to"
    class="flex items-center gap-2"
  >
    <img
      src="~/assets/icons/logo.svg"
      alt="Chameleon Vitae Logo"
      class="-scale-x-100"
      :class="logoSizeClass"
    />
    <span
      v-if="showText"
      class="font-bold text-zinc-100"
      :class="textSizeClass"
    >
      {{ text }}
    </span>
  </NuxtLink>
</template>

<script setup lang="ts">
interface Props {
  to?: string
  showText?: boolean
  text?: string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  to: '/',
  showText: true,
  text: 'Chameleon Vitae',
  size: 'md'
})

const logoSizeClass = computed(() => {
  const sizes = {
    sm: 'h-6 w-6',
    md: 'h-8 w-8',
    lg: 'h-10 w-10'
  }
  return sizes[props.size]
})

const textSizeClass = computed(() => {
  const sizes = {
    sm: 'text-base',
    md: 'text-lg',
    lg: 'text-xl'
  }
  return sizes[props.size]
})
</script>
