<script setup lang="ts">
defineProps<{
  tickets: { number: number; status: string }[]
  selectedNumbers: number[]
}>()

const emit = defineEmits<{
  toggle: [number: number]
}>()

function statusClass(status: string, number: number, selected: number[]): string {
  if (selected.includes(number)) return "bg-indigo-600 text-white ring-2 ring-indigo-400"
  switch (status) {
    case "AVAILABLE":
      return "bg-green-100 text-green-800 hover:bg-green-200 cursor-pointer"
    case "RESERVED":
      return "bg-yellow-100 text-yellow-800"
    case "PAID":
      return "bg-blue-100 text-blue-800"
    default:
      return "bg-gray-100 text-gray-400"
  }
}

function isClickable(status: string): boolean {
  return status === "AVAILABLE"
}
</script>

<template>
  <div class="grid grid-cols-5 sm:grid-cols-8 md:grid-cols-10 gap-2">
    <button
      v-for="ticket in tickets"
      :key="ticket.number"
      :disabled="!isClickable(ticket.status)"
      :class="statusClass(ticket.status, ticket.number, selectedNumbers)"
      class="aspect-square rounded-lg text-sm font-medium flex items-center justify-center transition-colors"
      @click="emit('toggle', ticket.number)"
    >
      {{ ticket.number }}
    </button>
  </div>
</template>
