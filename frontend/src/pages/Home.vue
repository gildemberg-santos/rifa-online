<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"
import RaffleCard from "../components/RaffleCard.vue"

interface Raffle {
  id: string
  title: string
  description: string
  ticketPrice: number
  maxNumbers: number
  drawDate: string
  imageUrl?: string
  status: string
}

const raffles = ref<Raffle[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    raffles.value = await api.get<Raffle[]>("/raffles")
  } catch (e) {
    console.error("Failed to load raffles", e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="text-center mb-10">
      <h1 class="text-4xl font-bold text-gray-900 mb-2">Rifa Online</h1>
      <p class="text-lg text-gray-600">
        Participe de rifas e concorra a prêmios incríveis!
      </p>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">
      Carregando...
    </div>

    <div
      v-else-if="raffles.length === 0"
      class="text-center py-12 text-gray-500"
    >
      Nenhuma rifa ativa no momento.
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      <RaffleCard
        v-for="raffle in raffles"
        :key="raffle.id"
        :raffle="raffle"
      />
    </div>
  </div>
</template>
