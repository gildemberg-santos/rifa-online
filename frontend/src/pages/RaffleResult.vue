<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { api } from "../utils/api"

interface Raffle {
  id: string
  title: string
  winnerNumber: number | null
  status: string
}

interface RaffleDetail {
  raffle: Raffle
}

const route = useRoute()
const detail = ref<RaffleDetail | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    detail.value = await api.get<RaffleDetail>(`/raffles/${route.params.id}`)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 text-center animate-fade-in">
    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="detail" class="bg-white rounded-2xl shadow-sm border border-gray-200 p-10">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">{{ detail.raffle.title }}</h1>

      <div v-if="detail.raffle.status === 'DRAWN' && detail.raffle.winnerNumber">
        <div class="w-28 h-28 sm:w-36 sm:h-36 mx-auto bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 rounded-2xl sm:rounded-3xl flex items-center justify-center shadow-xl mb-4 animate-bounce-in">
          <span class="text-3xl sm:text-5xl font-extrabold text-white drop-shadow-lg">{{ detail.raffle.winnerNumber }}</span>
        </div>
        <p class="text-base sm:text-lg text-gray-600">Número vencedor!</p>
      </div>

      <div v-else class="py-8">
        <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-gray-500">O sorteio ainda não foi realizado.</p>
      </div>

      <div class="mt-8">
        <router-link
          to="/"
          class="inline-flex items-center text-indigo-600 hover:text-indigo-700 font-medium"
        >
          <svg class="w-4 h-4 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Voltar para Home
        </router-link>
      </div>
    </div>
  </div>
</template>