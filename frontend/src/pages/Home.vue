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
  <div>
    <section class="relative overflow-hidden bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 text-white">
      <div class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmYiIGZpbGwtb3BhY2l0eT0iMC4wNSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMiIvPjwvZz48L2c+PC9zdmc+')] opacity-40"></div>
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 md:py-28 relative">
        <div class="text-center animate-fade-in">
          <h1 class="text-4xl md:text-6xl font-extrabold tracking-tight mb-4">
            Rifa Online
          </h1>
          <p class="text-lg md:text-xl text-indigo-100 max-w-2xl mx-auto">
            Participe de rifas e concorra a prêmios incríveis! Escolha seus números e garanta sua participação.
          </p>
          <div class="mt-8 flex justify-center gap-4">
            <a href="#raffles-section" class="inline-flex items-center px-6 py-3 bg-white text-indigo-700 font-semibold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all">
              Ver Rifas
              <svg class="w-5 h-5 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </a>
          </div>
        </div>
      </div>
      <div class="absolute bottom-0 left-0 right-0 h-16 bg-gradient-to-t from-gray-50 to-transparent"></div>
    </section>

    <section id="raffles-section" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Rifas Ativas</h2>
          <p class="text-gray-500 text-sm mt-1">Participe escolhendo seus números da sorte</p>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-16">
        <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
      </div>

      <div v-else-if="raffles.length === 0" class="text-center py-16 bg-white rounded-2xl border border-gray-200">
        <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-gray-500 text-lg">Nenhuma rifa ativa no momento.</p>
        <p class="text-gray-400 text-sm mt-1">Volte mais tarde para conferir as novidades.</p>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <RaffleCard
          v-for="(raffle, idx) in raffles"
          :key="raffle.id"
          :raffle="raffle"
          :style="{ animationDelay: `${idx * 0.08}s` }"
          class="animate-slide-up opacity-0"
        />
      </div>
    </section>
  </div>
</template>

<style scoped>
.animate-slide-up {
  animation: slideUp 0.4s ease-out forwards;
}
@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>