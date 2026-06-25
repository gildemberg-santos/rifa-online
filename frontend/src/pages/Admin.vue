<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
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
  phone: string
  role: string
  subscriptionStatus: string
  subscriptionIsTrial: boolean
  subscriptionExpiresAt?: string
  createdAt: string
}

interface AdminRaffle {
  id: string
  title: string
  organizerId: string
  organizerName: string
  status: string
  ticketPrice: number
  maxNumbers: number
  soldTickets: number
  paidTickets: number
  revenue: number
  drawDate: string
  winnerNumber?: number
  createdAt: string
}

const stats = ref<AdminStats | null>(null)
const users = ref<AdminUser[]>([])
const raffles = ref<AdminRaffle[]>([])
const activeTab = ref<"stats" | "users" | "raffles">("stats")
const loading = ref(true)
const actionLoading = ref<string | null>(null)
const selectedUser = ref<AdminUser | null>(null)
const showUserModal = ref(false)
const userRaffles = ref<any[]>([])
const userPayments = ref<any[]>([])
const userTickets = ref<any[]>([])
const searchQuery = ref("")
const userSort = ref<"name" | "email" | "createdAt" | "subscriptionStatus">("createdAt")
const userSortDir = ref<"asc" | "desc">("desc")
const raffleSort = ref<"title" | "organizerName" | "status" | "createdAt" | "soldTickets">("createdAt")
const raffleSortDir = ref<"asc" | "desc">("desc")
const confirmAction = ref<{ id: string; action: string; label: string } | null>(null)

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

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value
  const q = searchQuery.value.toLowerCase()
  return users.value.filter(u =>
    u.name.toLowerCase().includes(q) ||
    u.email.toLowerCase().includes(q) ||
    (u.phone && u.phone.includes(q))
  )
})

const sortedUsers = computed(() => {
  const arr = [...filteredUsers.value]
  arr.sort((a, b) => {
    let cmp = 0
    switch (userSort.value) {
      case "name": cmp = a.name.localeCompare(b.name); break
      case "email": cmp = a.email.localeCompare(b.email); break
      case "createdAt": cmp = new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(); break
      case "subscriptionStatus": cmp = a.subscriptionStatus.localeCompare(b.subscriptionStatus); break
    }
    return userSortDir.value === "asc" ? cmp : -cmp
  })
  return arr
})

const sortedRaffles = computed(() => {
  const arr = [...raffles.value]
  arr.sort((a, b) => {
    let cmp = 0
    switch (raffleSort.value) {
      case "title": cmp = a.title.localeCompare(b.title); break
      case "organizerName": cmp = a.organizerName.localeCompare(b.organizerName); break
      case "status": cmp = a.status.localeCompare(b.status); break
      case "createdAt": cmp = new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(); break
      case "soldTickets": cmp = a.soldTickets - b.soldTickets; break
    }
    return raffleSortDir.value === "asc" ? cmp : -cmp
  })
  return arr
})

function formatDate(d: string): string {
  if (!d) return ""
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

function statusColor(s: string): string {
  switch (s) {
    case "ACTIVE": return "bg-green-100 text-green-700"
    case "DRAWN": return "bg-blue-100 text-blue-700"
    case "CANCELLED": return "bg-red-100 text-red-700"
    default: return "bg-gray-100 text-gray-600"
  }
}

function subLabel(s: string): string {
  switch (s) {
    case "ACTIVE": return "Ativa"
    case "INACTIVE": return "Inativa"
    case "PAST_DUE": return "Vencida"
    case "CANCELLED": return "Cancelada"
    default: return s || "Inativa"
  }
}

function subColor(s: string): string {
  switch (s) {
    case "ACTIVE": return "bg-green-100 text-green-700"
    case "PAST_DUE": return "bg-red-100 text-red-700"
    case "CANCELLED": return "bg-gray-100 text-gray-600"
    default: return "bg-gray-100 text-gray-600"
  }
}

function roleLabel(r: string): string {
  return r === "ADMIN" ? "Admin" : "Usuário"
}

function roleColor(r: string): string {
  return r === "ADMIN" ? "bg-purple-100 text-purple-700" : "bg-gray-100 text-gray-600"
}

function salesProgress(r: AdminRaffle): number {
  if (r.maxNumbers === 0) return 0
  return Math.round((r.paidTickets / r.maxNumbers) * 100)
}

async function updateSubscription(userId: string, action: string) {
  actionLoading.value = userId
  try {
    await api.put(`/admin/users/${userId}/subscription`, { action })
    const [s, u] = await Promise.all([
      api.get<AdminStats>("/admin/stats"),
      api.get<AdminUser[]>("/admin/users"),
    ])
    stats.value = s
    users.value = u
  } catch {
    console.error("Failed to update subscription")
  } finally {
    actionLoading.value = null
  }
}

async function raffleAction(raffleId: string, action: "draw" | "cancel") {
  actionLoading.value = raffleId
  confirmAction.value = null
  try {
    if (action === "draw") {
      await api.post(`/raffles/${raffleId}/draw`)
    } else {
      await api.patch(`/raffles/${raffleId}/cancel`)
    }
    const r = await api.get<AdminRaffle[]>("/admin/raffles")
    raffles.value = r
  } catch (e: any) {
    console.error("Failed to update raffle", e?.message)
  } finally {
    actionLoading.value = null
  }
}

async function viewUserDetails(user: AdminUser) {
  selectedUser.value = user
  try {
    const data = await api.get<any>(`/admin/users/${user.id}`)
    userRaffles.value = data.raffles || []
    userPayments.value = data.payments || []
    userTickets.value = data.tickets || []
    showUserModal.value = true
  } catch {
    console.error("Failed to load user details")
  }
}

function closeModal() {
  showUserModal.value = false
  selectedUser.value = null
  userRaffles.value = []
  userPayments.value = []
  userTickets.value = []
}

function toggleSort(field: string, type: "user" | "raffle") {
  if (type === "user") {
    if (userSort.value === field) {
      userSortDir.value = userSortDir.value === "asc" ? "desc" : "asc"
    } else {
      userSort.value = field as any
      userSortDir.value = "asc"
    }
  } else {
    if (raffleSort.value === field) {
      raffleSortDir.value = raffleSortDir.value === "asc" ? "desc" : "asc"
    } else {
      raffleSort.value = field as any
      raffleSortDir.value = "asc"
    }
  }
}

function sortIcon(field: string, type: "user" | "raffle"): string {
  const current = type === "user" ? userSort.value : raffleSort.value
  const dir = type === "user" ? userSortDir.value : raffleSortDir.value
  if (current !== field) return "text-gray-300"
  return dir === "asc" ? "text-indigo-600" : "text-indigo-600 rotate-180"
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8 animate-fade-in">
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Administração</h1>
        <p class="text-gray-500 text-sm mt-1">Gerencie usuários, assinaturas e rifas</p>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else>
      <div v-if="stats" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4 mb-8">
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Usuários</p>
          <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.totalUsers }}</p>
          <p class="text-xs text-green-600 mt-1">{{ stats.activeUsers }} ativos</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Assinantes</p>
          <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.activeUsers }}</p>
          <p class="text-xs text-amber-600 mt-1">{{ stats.trialUsers }} em trial</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Vencidos</p>
          <p class="text-2xl font-bold text-red-600 mt-1">{{ stats.pastDueUsers }}</p>
          <p class="text-xs text-gray-500 mt-1">past due</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Rifas</p>
          <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.totalRaffles }}</p>
          <p class="text-xs text-indigo-600 mt-1">{{ stats.activeRaffles }} ativas</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Ingressos</p>
          <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.totalPaidTickets }}</p>
          <p class="text-xs text-gray-500 mt-1">pagos</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Receita</p>
          <p class="text-2xl font-bold text-gray-900 mt-1">R$ {{ (stats.totalRevenue / 100).toFixed(2) }}</p>
          <p class="text-xs text-gray-500 mt-1">total</p>
        </div>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="flex border-b border-gray-200">
          <button @click="activeTab = 'stats'"
            class="px-5 py-3 text-sm font-medium transition-colors"
            :class="activeTab === 'stats' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
          >Resumo</button>
          <button @click="activeTab = 'users'"
            class="px-5 py-3 text-sm font-medium transition-colors"
            :class="activeTab === 'users' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
          >Usuários ({{ users.length }})</button>
          <button @click="activeTab = 'raffles'"
            class="px-5 py-3 text-sm font-medium transition-colors"
            :class="activeTab === 'raffles' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
          >Rifas ({{ raffles.length }})</button>
        </div>

        <div v-if="activeTab === 'stats'" class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-4">Distribuição de Usuários</h3>
              <div class="space-y-3" v-if="stats">
                <div>
                  <div class="flex justify-between text-sm mb-1">
                    <span class="text-gray-600">Ativos</span>
                    <span class="font-semibold text-gray-900">{{ stats.activeUsers }}</span>
                  </div>
                  <div class="w-full bg-gray-100 rounded-full h-2">
                    <div class="bg-green-500 h-2 rounded-full" :style="{ width: stats.totalUsers > 0 ? (stats.activeUsers / stats.totalUsers * 100) + '%' : '0%' }"></div>
                  </div>
                </div>
                <div>
                  <div class="flex justify-between text-sm mb-1">
                    <span class="text-gray-600">Trial</span>
                    <span class="font-semibold text-amber-600">{{ stats.trialUsers }}</span>
                  </div>
                  <div class="w-full bg-gray-100 rounded-full h-2">
                    <div class="bg-amber-500 h-2 rounded-full" :style="{ width: stats.totalUsers > 0 ? (stats.trialUsers / stats.totalUsers * 100) + '%' : '0%' }"></div>
                  </div>
                </div>
                <div>
                  <div class="flex justify-between text-sm mb-1">
                    <span class="text-gray-600">Vencidos</span>
                    <span class="font-semibold text-red-600">{{ stats.pastDueUsers }}</span>
                  </div>
                  <div class="w-full bg-gray-100 rounded-full h-2">
                    <div class="bg-red-500 h-2 rounded-full" :style="{ width: stats.totalUsers > 0 ? (stats.pastDueUsers / stats.totalUsers * 100) + '%' : '0%' }"></div>
                  </div>
                </div>
              </div>
            </div>
            <div>
              <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-4">Distribuição de Rifas</h3>
              <div class="space-y-3" v-if="stats">
                <div>
                  <div class="flex justify-between text-sm mb-1">
                    <span class="text-gray-600">Ativas</span>
                    <span class="font-semibold text-indigo-600">{{ stats.activeRaffles }}</span>
                  </div>
                  <div class="w-full bg-gray-100 rounded-full h-2">
                    <div class="bg-indigo-500 h-2 rounded-full" :style="{ width: stats.totalRaffles > 0 ? (stats.activeRaffles / stats.totalRaffles * 100) + '%' : '0%' }"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'users'">
          <div class="p-4 border-b border-gray-100">
            <input v-model="searchQuery" type="text" placeholder="Buscar por nome, email ou telefone..."
              class="w-full max-w-md px-4 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none" />
          </div>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead class="bg-gray-50">
                <tr>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('name', 'user')">
                    Nome <span class="inline-block text-xs transition-transform" :class="sortIcon('name', 'user')">▲</span>
                  </th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('email', 'user')">
                    Email <span class="inline-block text-xs transition-transform" :class="sortIcon('email', 'user')">▲</span>
                  </th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Telefone</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Tipo</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('subscriptionStatus', 'user')">
                    Assinatura <span class="inline-block text-xs transition-transform" :class="sortIcon('subscriptionStatus', 'user')">▲</span>
                  </th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Expira</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('createdAt', 'user')">
                    Registro <span class="inline-block text-xs transition-transform" :class="sortIcon('createdAt', 'user')">▲</span>
                  </th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Ações</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-for="u in sortedUsers" :key="u.id" class="hover:bg-gray-50">
                  <td class="px-5 py-3 font-medium text-gray-900">{{ u.name }}</td>
                  <td class="px-5 py-3 text-gray-600 text-xs">{{ u.email }}</td>
                  <td class="px-5 py-3 text-gray-500 text-xs">{{ u.phone || '-' }}</td>
                  <td class="px-5 py-3">
                    <span class="inline-flex px-2 py-0.5 rounded-full text-xs font-semibold" :class="roleColor(u.role)">{{ roleLabel(u.role) }}</span>
                  </td>
                  <td class="px-5 py-3">
                    <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold" :class="subColor(u.subscriptionStatus)">
                      {{ subLabel(u.subscriptionStatus) }}
                      <span v-if="u.subscriptionIsTrial" class="text-amber-500">(teste)</span>
                    </span>
                  </td>
                  <td class="px-5 py-3 text-gray-500 text-xs">{{ u.subscriptionExpiresAt ? formatDate(u.subscriptionExpiresAt) : '-' }}</td>
                  <td class="px-5 py-3 text-gray-500 text-xs">{{ formatDate(u.createdAt) }}</td>
                  <td class="px-5 py-3">
                    <div class="flex gap-1.5 flex-wrap">
                      <button @click="updateSubscription(u.id, 'activate')" :disabled="actionLoading === u.id"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                        :class="u.subscriptionStatus === 'ACTIVE' ? 'bg-green-100 text-green-700 cursor-not-allowed' : 'bg-green-50 text-green-600 hover:bg-green-100'"
                      >Ativar</button>
                      <button @click="updateSubscription(u.id, 'cancel')" :disabled="actionLoading === u.id"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                        :class="u.subscriptionStatus === 'CANCELLED' ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-gray-50 text-gray-600 hover:bg-gray-100'"
                      >Cancelar</button>
                      <button @click="updateSubscription(u.id, 'past_due')" :disabled="actionLoading === u.id"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                        :class="u.subscriptionStatus === 'PAST_DUE' ? 'bg-red-100 text-red-400 cursor-not-allowed' : 'bg-red-50 text-red-600 hover:bg-red-100'"
                      >Vencer</button>
                      <button @click="updateSubscription(u.id, 'inactive')" :disabled="actionLoading === u.id"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg text-gray-500 bg-gray-50 hover:bg-gray-100 disabled:opacity-50 transition-colors"
                      >Zerar</button>
                      <button @click="viewUserDetails(u)"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg text-indigo-600 bg-indigo-50 hover:bg-indigo-100 transition-colors"
                      >Detalhes</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            <p v-if="sortedUsers.length === 0" class="text-center py-8 text-gray-500">Nenhum usuário encontrado</p>
          </div>
        </div>

        <div v-if="activeTab === 'raffles'" class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50">
              <tr>
                <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('title', 'raffle')">
                  Título <span class="inline-block text-xs transition-transform" :class="sortIcon('title', 'raffle')">▲</span>
                </th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('organizerName', 'raffle')">
                  Criador <span class="inline-block text-xs transition-transform" :class="sortIcon('organizerName', 'raffle')">▲</span>
                </th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('status', 'raffle')">
                  Status <span class="inline-block text-xs transition-transform" :class="sortIcon('status', 'raffle')">▲</span>
                </th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Valor</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('soldTickets', 'raffle')">
                  Vendas <span class="inline-block text-xs transition-transform" :class="sortIcon('soldTickets', 'raffle')">▲</span>
                </th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Receita</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600 cursor-pointer select-none" @click="toggleSort('createdAt', 'raffle')">
                  Data <span class="inline-block text-xs transition-transform" :class="sortIcon('createdAt', 'raffle')">▲</span>
                </th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Ações</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="r in sortedRaffles" :key="r.id" class="hover:bg-gray-50">
                <td class="px-5 py-3 font-medium text-gray-900 max-w-[200px] truncate" :title="r.title">{{ r.title }}</td>
                <td class="px-5 py-3 text-gray-600">{{ r.organizerName }}</td>
                <td class="px-5 py-3">
                  <span class="inline-flex px-2.5 py-1 rounded-full text-xs font-semibold" :class="statusColor(r.status)">
                    {{ statusLabel(r.status) }}
                  </span>
                </td>
                <td class="px-5 py-3 text-gray-700">R$ {{ (r.ticketPrice / 100).toFixed(2) }}</td>
                <td class="px-5 py-3">
                  <div class="flex items-center gap-2">
                    <div class="flex-1 min-w-[80px] bg-gray-100 rounded-full h-2">
                      <div class="h-2 rounded-full transition-all" :class="r.status === 'ACTIVE' ? 'bg-indigo-500' : r.status === 'DRAWN' ? 'bg-blue-500' : 'bg-gray-300'"
                        :style="{ width: salesProgress(r) + '%' }"></div>
                    </div>
                    <span class="text-xs text-gray-600 whitespace-nowrap">{{ r.paidTickets }}/{{ r.maxNumbers }}</span>
                  </div>
                </td>
                <td class="px-5 py-3 text-gray-700">R$ {{ (r.revenue / 100).toFixed(2) }}</td>
                <td class="px-5 py-3 text-gray-500 text-xs">{{ formatDate(r.createdAt) }}</td>
                <td class="px-5 py-3">
                  <div class="flex gap-1.5" v-if="r.status === 'ACTIVE'">
                    <button @click="confirmAction = { id: r.id, action: 'draw', label: 'realizar sorteio' }" :disabled="actionLoading === r.id"
                      class="px-2.5 py-1 text-xs font-medium rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 disabled:opacity-50 transition-colors"
                      :title="r.paidTickets === 0 ? 'Nenhum ingresso pago' : ''"
                    >Sortear</button>
                    <button @click="confirmAction = { id: r.id, action: 'cancel', label: 'cancelar' }" :disabled="actionLoading === r.id"
                      class="px-2.5 py-1 text-xs font-medium rounded-lg bg-red-50 text-red-600 hover:bg-red-100 disabled:opacity-50 transition-colors"
                    >Cancelar</button>
                  </div>
                  <span v-else-if="r.status === 'DRAWN'" class="text-xs text-blue-600 font-medium">
                    Ganhador: {{ r.winnerNumber ?? '-' }}
                  </span>
                  <span v-else class="text-xs text-gray-400">-</span>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-if="raffles.length === 0" class="text-center py-8 text-gray-500">Nenhuma rifa encontrada</p>
        </div>
      </div>
    </template>

    <div v-if="showUserModal && selectedUser" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm" @click.self="closeModal">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-2xl mx-4 max-h-[80vh] overflow-y-auto animate-scale-in">
        <div class="px-6 py-5 border-b border-gray-100 flex items-center justify-between">
          <div>
            <h2 class="text-lg font-bold text-gray-900">{{ selectedUser.name }}</h2>
            <p class="text-sm text-gray-500">{{ selectedUser.email }}</p>
            <div class="flex gap-2 mt-1">
              <span class="inline-flex px-2 py-0.5 rounded-full text-xs font-semibold" :class="roleColor(selectedUser.role)">{{ roleLabel(selectedUser.role) }}</span>
              <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold" :class="subColor(selectedUser.subscriptionStatus)">
                {{ subLabel(selectedUser.subscriptionStatus) }}
                <span v-if="selectedUser.subscriptionIsTrial" class="text-amber-500">(teste)</span>
              </span>
            </div>
          </div>
          <button @click="closeModal" class="p-2 rounded-lg hover:bg-gray-100 transition-colors">
            <svg class="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-6 space-y-6">
          <div>
            <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">Rifas ({{ userRaffles.length }})</h3>
            <div v-if="userRaffles.length === 0" class="text-sm text-gray-400">Nenhuma rifa criada</div>
            <div v-for="r in userRaffles" :key="r.id" class="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
              <div>
                <span class="text-sm font-medium text-gray-900">{{ r.title }}</span>
                <span class="text-xs text-gray-400 ml-2">{{ r.maxNumbers }} números</span>
              </div>
              <span class="text-xs px-2 py-0.5 rounded-full" :class="statusColor(r.status)">{{ statusLabel(r.status) }}</span>
            </div>
          </div>

          <div>
            <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">Pagamentos ({{ userPayments.length }})</h3>
            <div v-if="userPayments.length === 0" class="text-sm text-gray-400">Nenhum pagamento</div>
            <div v-for="p in userPayments" :key="p.id" class="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
              <div>
                <span class="text-sm text-gray-700">
                  <template v-if="p.type === 'SUBSCRIPTION'">Assinatura</template>
                  <template v-else>Rifa</template>
                  — R$ {{ (p.amount / 100).toFixed(2) }}
                </span>
                <span class="text-xs text-gray-400 ml-2">{{ new Date(p.createdAt).toLocaleString('pt-BR') }}</span>
              </div>
              <span class="text-xs px-2 py-0.5 rounded-full"
                :class="p.status === 'PAID' ? 'bg-green-100 text-green-700' : p.status === 'PENDING' ? 'bg-yellow-100 text-yellow-700' : 'bg-gray-100 text-gray-600'"
              >{{ p.status }}</span>
            </div>
          </div>

          <div v-if="userTickets.length > 0">
            <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">Ingressos ({{ userTickets.length }})</h3>
            <div v-for="t in userTickets.slice(0, 20)" :key="t.id" class="flex items-center justify-between py-1.5 border-b border-gray-50 last:border-0">
              <span class="text-sm text-gray-700">Nº {{ t.number }}</span>
              <span class="text-xs text-gray-400">{{ t.raffleTitle || '' }}</span>
            </div>
            <p v-if="userTickets.length > 20" class="text-xs text-gray-400 mt-2">...e mais {{ userTickets.length - 20 }} ingressos</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="confirmAction" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm" @click.self="confirmAction = null">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-sm mx-4 p-6 animate-scale-in">
        <h3 class="text-lg font-bold text-gray-900 mb-2">Confirmar ação</h3>
        <p class="text-sm text-gray-600 mb-6">Tem certeza que deseja {{ confirmAction.label }} esta rifa?</p>
        <div class="flex gap-3 justify-end">
          <button @click="confirmAction = null" class="px-4 py-2 text-sm font-medium text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors">Cancelar</button>
          <button @click="raffleAction(confirmAction.id, confirmAction.action as any)" :disabled="actionLoading === confirmAction.id"
            class="px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors"
            :class="confirmAction.action === 'draw' ? 'bg-blue-600 hover:bg-blue-700' : 'bg-red-600 hover:bg-red-700'"
          >Confirmar</button>
        </div>
      </div>
    </div>
  </div>
</template>
