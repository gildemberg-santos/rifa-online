<script setup lang="ts">
defineProps<{
  tickets: { number: number; status: string }[]
  selectedNumbers: number[]
}>()

const emit = defineEmits<{
  toggle: [number: number]
}>()

function statusClass(status: string, number: number, selected: number[]): string {
  if (selected.includes(number))
    return "bg-gradient-to-br from-indigo-600 to-purple-600 text-white shadow-md shadow-indigo-300 ring-2 ring-indigo-400 scale-105 z-10"
  switch (status) {
    case "AVAILABLE":
      return "bg-white text-gray-700 border border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 hover:shadow-sm cursor-pointer"
    case "RESERVED":
      return "bg-amber-50 text-amber-700 border border-amber-200 cursor-not-allowed"
    case "PAID":
      return "bg-emerald-50 text-emerald-700 border border-emerald-200 cursor-not-allowed"
    default:
      return "bg-gray-50 text-gray-300 border border-gray-100 cursor-not-allowed"
  }
}

function statusLabel(status: string): string | null {
  if (status === "PAID") return "✓"
  return null
}
</script>

<template>
  <div class="grid grid-cols-5 sm:grid-cols-8 md:grid-cols-10 gap-2">
    <button
      v-for="ticket in tickets"
      :key="ticket.number"
      :disabled="ticket.status !== 'AVAILABLE'"
      :class="statusClass(ticket.status, ticket.number, selectedNumbers)"
      class="relative aspect-square rounded-xl text-sm font-semibold flex items-center justify-center transition-all duration-150"
      @click="emit('toggle', ticket.number)"
    >
      <span v-if="statusLabel(ticket.status)" class="absolute top-0.5 right-0.5 text-[10px]">{{ statusLabel(ticket.status) }}</span>
      {{ ticket.number }}
    </button>
  </div>
</template>