<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRoute } from "vue-router"
import { api } from "../utils/api"
import NumberGrid from "../components/NumberGrid.vue"

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

interface Ticket {
  id: string
  number: number
  status: string
}

interface RaffleDetail {
  raffle: Raffle
  tickets: Ticket[]
}

const route = useRoute()
const detail = ref<RaffleDetail | null>(null)
const loading = ref(true)
const selectedNumbers = ref<number[]>([])

onMounted(async () => {
  try {
    detail.value = await api.get<RaffleDetail>(`/raffles/${route.params.id}`)
  } catch (e) {
    console.error("Failed to load raffle detail", e)
  } finally {
    loading.value = false
  }
})

function toggleNumber(num: number) {
  const idx = selectedNumbers.value.indexOf(num)
  if (idx === -1) {
    selectedNumbers.value.push(num)
  } else {
    selectedNumbers.value.splice(idx, 1)
  }
}

const totalPrice = computed(() => {
  if (!detail.value) return 0
  return (detail.value.raffle.ticketPrice * selectedNumbers.value.length) / 100
})

</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div v-if="loading" class="text-center py-12 text-gray-500">
      Carregando...
    </div>

    <template v-else-if="detail">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">{{ detail.raffle.title }}</h1>
        <p class="text-gray-600 mt-2">{{ detail.raffle.description }}</p>

        <div class="flex gap-6 mt-4 text-sm">
          <span class="text-indigo-600 font-medium">
            R$ {{ (detail.raffle.ticketPrice / 100).toFixed(2) }} por número
          </span>
          <span class="text-gray-500">
            Sorteio: {{ new Date(detail.raffle.drawDate).toLocaleDateString("pt-BR") }}
          </span>
          <span
            class="font-medium"
            :class="detail.raffle.status === 'ACTIVE' ? 'text-green-600' : 'text-red-600'"
          >
            {{ detail.raffle.status === "ACTIVE" ? "Ativa" : detail.raffle.status }}
          </span>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold mb-4">Escolha seus números</h2>

        <NumberGrid
          :tickets="detail.tickets"
          :selected-numbers="selectedNumbers"
          @toggle="toggleNumber"
        />

        <div v-if="selectedNumbers.length > 0" class="mt-6 p-4 bg-gray-50 rounded-lg flex justify-between items-center">
          <div>
            <span class="text-sm text-gray-600">
              {{ selectedNumbers.length }} número(s) selecionado(s)
            </span>
            <span class="ml-4 text-lg font-bold text-indigo-600">
              R$ {{ totalPrice.toFixed(2) }}
            </span>
          </div>

          <router-link
            :to="`/raffles/${detail.raffle.id}/checkout?numbers=${selectedNumbers.join(',')}`"
            class="bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 font-medium"
          >
            Comprar
          </router-link>
        </div>
      </div>
    </template>
  </div>
</template>
