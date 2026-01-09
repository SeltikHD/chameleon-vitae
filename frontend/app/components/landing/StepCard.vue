<template>
  <div class="relative flex flex-col items-center text-center">
    <!-- Step Number Badge -->
    <div class="absolute -top-3 left-1/2 -translate-x-1/2">
      <span
        class="flex h-8 w-8 items-center justify-center rounded-full bg-emerald-500 text-sm font-bold text-zinc-950"
      >
        {{ number }}
      </span>
    </div>

    <UCard class="w-full pt-8">
      <div class="flex flex-col items-center">
        <!-- Icon -->
        <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-zinc-800">
          <UIcon :name="icon" class="h-6 w-6 text-emerald-400" />
        </div>

        <!-- Title -->
        <h3 class="text-lg font-semibold text-zinc-100">
          {{ title }}
        </h3>

        <!-- Description -->
        <p class="mt-2 text-sm text-zinc-400">
          {{ description }}
        </p>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
interface Props {
  number: number
  title: string
  description: string
  icon: string
}

defineProps<Props>()
</script>
