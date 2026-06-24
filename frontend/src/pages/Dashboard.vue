<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"

interface Raffle {
  id: string
  title: string
  ticketPrice: number
  maxNumbers: number
  status: string
  drawDate: string
}

interface DashboardStats {
  totalRaffles: number
  activeRaffles: number
  drawnRaffles: number
  cancelledRaffles: number
  totalSoldTickets: number
  totalRevenue: number
  totalReservedTickets: number
  totalMaxNumbers: number
  totalAvailableTickets: number
}

const raffles = ref<Raffle[]>([])
const stats = ref<DashboardStats | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const [r, s] = await Promise.all([
      api.get<Raffle[]>("/raffles/my"),
      api.get<DashboardStats>("/dashboard/stats"),
    ])
    raffles.value = r
    stats.value = s
  } catch (e) {
    console.error("Failed to load dashboard data", e)
  } finally {
    loading.value = false
  }
})

async function deleteRaffle(id: string) {
  if (!confirm("Tem certeza que deseja excluir esta rifa? Esta ação não pode ser desfeita.")) return
  try {
    await api.delete(`/raffles/${id}`)
    raffles.value = raffles.value.filter((r) => r.id !== id)
  } catch (e: any) {
    alert(e.message || "Erro ao excluir rifa")
  }
}

async function drawRaffle(id: string) {
  if (!confirm("Realizar sorteio?")) return
  try {
    await api.post(`/raffles/${id}/draw`)
    const idx = raffles.value.findIndex((r) => r.id === id)
    if (idx !== -1) raffles.value[idx].status = "DRAWN"
    if (stats.value) {
      stats.value.activeRaffles--
      stats.value.drawnRaffles++
    }
  } catch (e: any) {
    alert(e.message || "Erro ao realizar sorteio")
  }
}

function statusBadge(status: string) {
  switch (status) {
    case "ACTIVE": return "bg-green-100 text-green-700"
    case "DRAWN": return "bg-blue-100 text-blue-700"
    case "CANCELLED": return "bg-red-100 text-red-700"
    default: return "bg-gray-100 text-gray-600"
  }
}

function statusLabel(status: string) {
  switch (status) {
    case "ACTIVE": return "Ativa"
    case "DRAWN": return "Sorteada"
    case "CANCELLED": return "Cancelada"
    default: return status
  }
}
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8 animate-fade-in">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-gray-500 text-sm mt-1">Gerencie suas rifas</p>
      </div>
      <router-link
        to="/dashboard/raffles/new"
        class="inline-flex items-center px-4 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white text-sm font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md transition-all"
      >
        <svg class="w-4 h-4 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
        </svg>
        Criar Rifa
      </router-link>
    </div>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else>
      <div v-if="stats" class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-white rounded-xl border border-gray-200 p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Total Rifas</p>
              <p class="text-xl font-bold text-gray-900">{{ stats.totalRaffles }}</p>
            </div>
          </div>
          <div class="flex gap-2 mt-3 text-xs text-gray-500">
            <span class="inline-flex items-center gap-1">
              <span class="w-1.5 h-1.5 rounded-full bg-green-500"></span>
              {{ stats.activeRaffles }} ativas
            </span>
            <span class="inline-flex items-center gap-1">
              <span class="w-1.5 h-1.5 rounded-full bg-blue-500"></span>
              {{ stats.drawnRaffles }} sorteadas
            </span>
            <span v-if="stats.cancelledRaffles" class="inline-flex items-center gap-1">
              <span class="w-1.5 h-1.5 rounded-full bg-red-500"></span>
              {{ stats.cancelledRaffles }} canceladas
            </span>
          </div>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Ingressos Vendidos</p>
              <p class="text-xl font-bold text-gray-900">{{ stats.totalSoldTickets }}</p>
            </div>
          </div>
          <p class="mt-3 text-xs text-gray-500">
            {{ stats.totalReservedTickets }} reservados · {{ stats.totalAvailableTickets }} disponíveis
          </p>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-amber-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Receita Total</p>
              <p class="text-xl font-bold text-gray-900">R$ {{ (stats.totalRevenue / 100).toFixed(2) }}</p>
            </div>
          </div>
          <p class="mt-3 text-xs text-gray-500">
            {{ stats.totalMaxNumbers }} números no total
          </p>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Vendidos %</p>
              <p class="text-xl font-bold text-gray-900">
                {{ stats.totalMaxNumbers > 0 ? ((stats.totalSoldTickets / stats.totalMaxNumbers) * 100).toFixed(1) : 0 }}%
              </p>
            </div>
          </div>
          <p class="mt-3 text-xs text-gray-500">
            {{ stats.totalSoldTickets }} de {{ stats.totalMaxNumbers }} números
          </p>
        </div>
      </div>

      <div v-if="raffles.length === 0" class="text-center py-16 bg-white rounded-2xl border border-gray-200">
        <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        <p class="text-gray-500 text-lg">Você ainda não criou nenhuma rifa.</p>
        <router-link to="/dashboard/raffles/new" class="mt-4 inline-flex items-center text-indigo-600 hover:text-indigo-700 font-medium">
          <svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          Criar primeira rifa
        </router-link>
      </div>

      <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50">
              <tr>
                <th class="text-left px-5 py-3.5 font-semibold text-gray-600">Título</th>
                <th class="text-left px-5 py-3.5 font-semibold text-gray-600">Status</th>
                <th class="text-left px-5 py-3.5 font-semibold text-gray-600">Valor</th>
                <th class="text-left px-5 py-3.5 font-semibold text-gray-600">Números</th>
                <th class="text-left px-5 py-3.5 font-semibold text-gray-600">Ações</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="raffle in raffles" :key="raffle.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-5 py-4 font-medium text-gray-900">{{ raffle.title }}</td>
                <td class="px-5 py-4">
                  <span class="inline-flex px-2.5 py-1 rounded-full text-xs font-semibold" :class="statusBadge(raffle.status)">
                    {{ statusLabel(raffle.status) }}
                  </span>
                </td>
                <td class="px-5 py-4 text-gray-700">R$ {{ (raffle.ticketPrice / 100).toFixed(2) }}</td>
                <td class="px-5 py-4 text-gray-700">{{ raffle.maxNumbers }}</td>
                <td class="px-5 py-4">
                  <div class="flex gap-2">
                    <router-link
                      v-if="raffle.status === 'ACTIVE'"
                      :to="`/dashboard/raffles/${raffle.id}/edit`"
                      class="inline-flex items-center px-3 py-1.5 text-xs font-medium text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors"
                    >
                      Editar
                    </router-link>
                    <button
                      v-if="raffle.status === 'ACTIVE'"
                      @click="drawRaffle(raffle.id)"
                      class="inline-flex items-center px-3 py-1.5 text-xs font-medium text-green-600 bg-green-50 hover:bg-green-100 rounded-lg transition-colors"
                    >
                      Sortear
                    </button>
                    <router-link
                      v-if="raffle.status === 'DRAWN'"
                      :to="`/raffles/${raffle.id}/result`"
                      class="inline-flex items-center px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors"
                    >
                      Resultado
                    </router-link>
                    <button
                      @click="deleteRaffle(raffle.id)"
                      class="inline-flex items-center px-3 py-1.5 text-xs font-medium text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors"
                    >
                      Excluir
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>