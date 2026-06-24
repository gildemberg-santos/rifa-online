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
  <div class="max-w-lg mx-auto px-4 py-12 text-center">
    <div v-if="loading" class="text-gray-500">Carregando...</div>

    <div v-else-if="detail" class="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <h1 class="text-2xl font-bold text-gray-900 mb-4">{{ detail.raffle.title }}</h1>

      <div v-if="detail.raffle.status === 'DRAWN' && detail.raffle.winnerNumber">
        <div class="text-6xl font-bold text-indigo-600 my-6">
          {{ detail.raffle.winnerNumber }}
        </div>
        <p class="text-lg text-gray-600">Número vencedor!</p>
      </div>

      <div v-else>
        <p class="text-gray-500">O sorteio ainda não foi realizado.</p>
      </div>

      <router-link to="/" class="inline-block mt-6 text-indigo-600 hover:underline">
        Voltar para Home
      </router-link>
    </div>
  </div>
</template>
