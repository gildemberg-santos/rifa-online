<script setup lang="ts">
defineProps<{
  size?: "sm" | "md" | "lg"
}>()
</script>

<template>
  <div class="flex justify-center items-center">
    <div
      class="border-indigo-600 border-b-transparent rounded-full animate-spin"
      :class="{
        'w-5 h-5 border-2': size === 'sm',
        'w-8 h-8 border-3': size === 'md' || !size,
        'w-12 h-12 border-4': size === 'lg',
      }"
    ></div>
  </div>
</template>
