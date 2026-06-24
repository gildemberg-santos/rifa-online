<script setup lang="ts">
defineProps<{
  type?: "info" | "success" | "warning" | "error"
  message: string
}>()
</script>

<template>
  <div
    class="px-4 py-3 rounded-lg text-sm"
    :class="{
      'bg-blue-50 text-blue-700 border border-blue-200': type === 'info' || !type,
      'bg-green-50 text-green-700 border border-green-200': type === 'success',
      'bg-yellow-50 text-yellow-700 border border-yellow-200': type === 'warning',
      'bg-red-50 text-red-700 border border-red-200': type === 'error',
    }"
  >
    {{ message }}
  </div>
</template>
