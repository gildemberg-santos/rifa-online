<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { api } from "../utils/api"
import { useAuthStore } from "../stores/auth"
import NumberGrid from "../components/NumberGrid.vue"

const RESERVATION_TTL = 30 * 60

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
  reservedAt?: string
  reservationExpiresIn?: number
}

interface RaffleDetail {
  raffle: Raffle
  tickets: Ticket[]
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const detail = ref<RaffleDetail | null>(null)
const loading = ref(true)
const selectedNumbers = ref<number[]>([])
const countdown = ref(0)
const copied = ref(false)
let countdownTimer: ReturnType<typeof setInterval> | null = null

function shareRaffle() {
  const url = window.location.href
  if (navigator.share) {
    navigator.share({ title: detail.value?.raffle.title, url })
  } else {
    navigator.clipboard.writeText(url)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}

onMounted(async () => {
  try {
    const status = await api.get<{ subscriptionStatus: string }>("/subscription/status")
    if (status.subscriptionStatus !== "ACTIVE") {
      router.push({ name: "subscription" })
      return
    }

    detail.value = await api.get<RaffleDetail>(`/raffles/${route.params.id}`)
  } catch (e) {
    console.error("Failed to load raffle detail", e)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
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

const reservedTickets = computed(() => {
  if (!detail.value) return []
  return detail.value.tickets.filter(t => t.status === "RESERVED").map(t => {
    let expiresIn = t.reservationExpiresIn ?? 0
    if (t.reservedAt) {
      const elapsed = (Date.now() - new Date(t.reservedAt).getTime()) / 1000
      expiresIn = Math.max(0, RESERVATION_TTL - elapsed)
    }
    return { ...t, expiresIn }
  })
})

countdownTimer = setInterval(() => {
  countdown.value++
}, 1000)

function formatExpiresIn(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, "0")}`
}

function formatReservedAt(iso: string): string {
  return new Date(iso).toLocaleTimeString("pt-BR")
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else-if="detail">
      <div class="bg-gradient-to-br from-indigo-600 to-purple-700 rounded-2xl p-6 md:p-8 text-white mb-8 shadow-lg">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div>
            <h1 class="text-3xl font-extrabold">{{ detail.raffle.title }}</h1>
            <p class="text-indigo-100 mt-2">{{ detail.raffle.description }}</p>
          </div>
          <div class="flex flex-wrap gap-3">
            <div class="bg-white/15 backdrop-blur rounded-xl px-4 py-3 text-center">
              <p class="text-xs text-indigo-200 uppercase tracking-wider">Preço</p>
              <p class="text-xl font-bold">R$ {{ (detail.raffle.ticketPrice / 100).toFixed(2) }}</p>
            </div>
            <div class="bg-white/15 backdrop-blur rounded-xl px-4 py-3 text-center">
              <p class="text-xs text-indigo-200 uppercase tracking-wider">Sorteio</p>
              <p class="text-sm font-semibold">{{ new Date(detail.raffle.drawDate).toLocaleDateString("pt-BR") }}</p>
            </div>
            <div class="bg-white/15 backdrop-blur rounded-xl px-4 py-3 text-center">
              <p class="text-xs text-indigo-200 uppercase tracking-wider">Status</p>
              <p class="text-sm font-semibold" :class="detail.raffle.status === 'ACTIVE' ? 'text-green-300' : 'text-red-300'">
                {{ detail.raffle.status === "ACTIVE" ? "Ativa" : detail.raffle.status }}
              </p>
            </div>
            <router-link
              v-if="authStore.isAuthenticated && detail.raffle.status === 'ACTIVE'"
              :to="`/dashboard/raffles/${detail.raffle.id}/edit`"
              class="self-center bg-white/20 hover:bg-white/30 backdrop-blur rounded-xl px-4 py-3 text-center transition-colors"
            >
              <p class="text-xs text-indigo-200 uppercase tracking-wider">Ações</p>
              <p class="text-sm font-semibold text-white">Editar</p>
            </router-link>
            <button
              @click="shareRaffle"
              class="self-center bg-white/20 hover:bg-white/30 backdrop-blur rounded-xl px-4 py-3 text-center transition-colors cursor-pointer"
            >
              <p class="text-xs text-indigo-200 uppercase tracking-wider">Compartilhar</p>
              <p class="text-sm font-semibold text-white flex items-center justify-center gap-1">
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                </svg>
                {{ copied ? "Copiado!" : "Compartilhar" }}
              </p>
            </button>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">Escolha seus números</h2>

        <NumberGrid
          :tickets="detail.tickets"
          :selected-numbers="selectedNumbers"
          @toggle="toggleNumber"
        />

        <div v-if="selectedNumbers.length > 0" class="mt-6 p-5 bg-gradient-to-r from-indigo-50 to-purple-50 rounded-xl border border-indigo-100 animate-fade-in">
          <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <div>
              <span class="text-sm text-gray-600">{{ selectedNumbers.length }} número(s) selecionado(s)</span>
              <p class="text-2xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600">
                R$ {{ totalPrice.toFixed(2) }}
              </p>
            </div>
            <router-link
              :to="`/raffles/${detail.raffle.id}/checkout?numbers=${selectedNumbers.join(',')}`"
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md hover:shadow-lg transition-all"
            >
              Comprar
              <svg class="w-5 h-5 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </router-link>
          </div>
        </div>
      </div>

      <div v-if="authStore.isAuthenticated && reservedTickets.length > 0" class="mt-6 bg-white rounded-2xl shadow-sm border border-gray-200 p-6 animate-fade-in">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">Reservas ativas ({{ reservedTickets.length }})</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200">
                <th class="text-left px-4 py-3 font-medium text-gray-600">Número</th>
                <th class="text-left px-4 py-3 font-medium text-gray-600">Reservado às</th>
                <th class="text-left px-4 py-3 font-medium text-gray-600">Expira em</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="t in reservedTickets" :key="t.number" class="hover:bg-gray-50 transition-colors">
                <td class="px-4 py-3 font-semibold text-gray-900">{{ t.number }}</td>
                <td class="px-4 py-3 text-gray-600">{{ formatReservedAt(t.reservedAt!) }}</td>
                <td class="px-4 py-3" :class="t.expiresIn <= 60 ? 'text-red-600 font-semibold' : 'text-gray-600'">
                  <span class="inline-flex items-center gap-1">
                    <svg v-if="t.expiresIn <= 60" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ formatExpiresIn(t.expiresIn) }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>