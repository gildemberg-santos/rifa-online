<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"

interface AdminStats {
  totalUsers: number
  activeUsers: number
  totalRaffles: number
  activeRaffles: number
  totalPaidTickets: number
  totalRevenue: number
  trialUsers: number
  pastDueUsers: number
}

interface AdminUser {
  id: string
  name: string
  email: string
  role: string
  subscriptionStatus: string
  subscriptionIsTrial: boolean
  createdAt: string
}

interface AdminRaffle {
  id: string
  title: string
  organizerId: string
  status: string
  ticketPrice: number
  maxNumbers: number
  createdAt: string
}

const stats = ref<AdminStats | null>(null)
const users = ref<AdminUser[]>([])
const raffles = ref<AdminRaffle[]>([])
const activeTab = ref<"stats" | "users" | "raffles">("stats")
const loading = ref(true)

onMounted(async () => {
  loading.value = true
  try {
    const [s, u, r] = await Promise.all([
      api.get<AdminStats>("/admin/stats"),
      api.get<AdminUser[]>("/admin/users"),
      api.get<AdminRaffle[]>("/admin/raffles"),
    ])
    stats.value = s
    users.value = u
    raffles.value = r
  } catch {
    console.error("Failed to load admin data")
  } finally {
    loading.value = false
  }
})

function formatDate(d: string): string {
  return new Date(d).toLocaleDateString("pt-BR")
}

function statusLabel(s: string): string {
  switch (s) {
    case "ACTIVE": return "Ativa"
    case "DRAWN": return "Sorteada"
    case "CANCELLED": return "Cancelada"
    default: return s
  }
}

function subLabel(s: string): string {
  switch (s) {
    case "ACTIVE": return "Ativa"
    case "INACTIVE": return "Inativa"
    case "PAST_DUE": return "Vencida"
    case "CANCELLED": return "Cancelada"
    default: return s
  }
}
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8 animate-fade-in">
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Administração</h1>
      <p class="text-gray-500 text-sm mt-1">Visão geral da plataforma</p>
    </div>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else>
      <div v-if="stats" class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Usuários</p>
          <p class="text-xl font-bold text-gray-900 mt-1">{{ stats.totalUsers }}</p>
          <p class="text-xs text-green-600 mt-1">{{ stats.activeUsers }} ativos</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Teste</p>
          <p class="text-xl font-bold text-amber-600 mt-1">{{ stats.trialUsers }}</p>
          <p class="text-xs text-gray-500 mt-1">em período de teste</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Rifas</p>
          <p class="text-xl font-bold text-gray-900 mt-1">{{ stats.totalRaffles }}</p>
          <p class="text-xs text-indigo-600 mt-1">{{ stats.activeRaffles }} ativas</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Receita</p>
          <p class="text-xl font-bold text-gray-900 mt-1">R$ {{ (stats.totalRevenue / 100).toFixed(2) }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ stats.totalPaidTickets }} ingressos</p>
        </div>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="flex border-b border-gray-200">
          <button
            @click="activeTab = 'stats'"
            class="px-5 py-3 text-sm font-medium"
            :class="activeTab === 'stats' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
          >Resumo</button>
          <button
            @click="activeTab = 'users'"
            class="px-5 py-3 text-sm font-medium"
            :class="activeTab === 'users' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
          >Usuários ({{ users.length }})</button>
          <button
            @click="activeTab = 'raffles'"
            class="px-5 py-3 text-sm font-medium"
            :class="activeTab === 'raffles' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
          >Rifas ({{ raffles.length }})</button>
        </div>

        <div v-if="activeTab === 'users'" class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50">
              <tr>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Nome</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Email</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Assinatura</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Cadastro</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="u in users" :key="u.id" class="hover:bg-gray-50">
                <td class="px-5 py-3 font-medium text-gray-900">{{ u.name }}</td>
                <td class="px-5 py-3 text-gray-600">{{ u.email }}</td>
                <td class="px-5 py-3">
                  <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold"
                    :class="u.subscriptionStatus === 'ACTIVE' ? 'bg-green-100 text-green-700' : u.subscriptionStatus === 'PAST_DUE' ? 'bg-red-100 text-red-700' : 'bg-gray-100 text-gray-600'"
                  >
                    {{ subLabel(u.subscriptionStatus) }}
                    <span v-if="u.subscriptionIsTrial" class="text-amber-500">(teste)</span>
                  </span>
                </td>
                <td class="px-5 py-3 text-gray-500">{{ formatDate(u.createdAt) }}</td>
              </tr>
            </tbody>
          </table>
          <p v-if="users.length === 0" class="text-center py-8 text-gray-500">Nenhum usuário encontrado</p>
        </div>

        <div v-if="activeTab === 'raffles'" class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50">
              <tr>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Título</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Status</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Valor</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Números</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Criação</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="r in raffles" :key="r.id" class="hover:bg-gray-50">
                <td class="px-5 py-3 font-medium text-gray-900">{{ r.title }}</td>
                <td class="px-5 py-3">
                  <span class="inline-flex px-2.5 py-1 rounded-full text-xs font-semibold"
                    :class="r.status === 'ACTIVE' ? 'bg-green-100 text-green-700' : r.status === 'DRAWN' ? 'bg-blue-100 text-blue-700' : 'bg-red-100 text-red-700'"
                  >{{ statusLabel(r.status) }}</span>
                </td>
                <td class="px-5 py-3 text-gray-700">R$ {{ (r.ticketPrice / 100).toFixed(2) }}</td>
                <td class="px-5 py-3 text-gray-700">{{ r.maxNumbers }}</td>
                <td class="px-5 py-3 text-gray-500">{{ formatDate(r.createdAt) }}</td>
              </tr>
            </tbody>
          </table>
          <p v-if="raffles.length === 0" class="text-center py-8 text-gray-500">Nenhuma rifa encontrada</p>
        </div>
      </div>
    </template>
  </div>
</template>
