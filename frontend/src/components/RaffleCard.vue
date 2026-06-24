<script setup lang="ts">
defineProps<{
  raffle: {
    id: string
    title: string
    description: string
    ticketPrice: number
    maxNumbers: number
    drawDate: string
    imageUrl?: string
    status: string
  }
}>()
</script>

<template>
  <router-link
    :to="`/raffles/${raffle.id}`"
    class="block bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
  >
    <div v-if="raffle.imageUrl" class="h-40 bg-gray-200">
      <img
        :src="raffle.imageUrl"
        :alt="raffle.title"
        class="w-full h-full object-cover"
      />
    </div>
    <div v-else class="h-40 bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center">
      <span class="text-4xl text-white font-bold">{{ raffle.title.charAt(0) }}</span>
    </div>

    <div class="p-4">
      <h3 class="text-lg font-semibold text-gray-900 truncate">
        {{ raffle.title }}
      </h3>
      <p class="text-sm text-gray-500 mt-1 line-clamp-2">
        {{ raffle.description }}
      </p>

      <div class="flex justify-between items-center mt-4 text-sm">
        <span class="text-indigo-600 font-medium">
          R$ {{ (raffle.ticketPrice / 100).toFixed(2) }}
        </span>
        <span class="text-gray-500">
          {{ raffle.maxNumbers }} números
        </span>
      </div>
    </div>
  </router-link>
</template>
