<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { api } from "../utils/api"
import { sendEvent } from "../utils/analytics"

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

interface ContactMessage {
  id: string
  name: string
  contact: string
  message: string
  createdAt: string
}

const stats = ref<AdminStats | null>(null)
const users = ref<AdminUser[]>([])
const raffles = ref<AdminRaffle[]>([])
const contactMessages = ref<ContactMessage[]>([])
const activeTab = ref<"users" | "raffles" | "contact">("users")
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
    const [s, u, r, c] = await Promise.all([
      api.get<AdminStats>("/admin/stats"),
      api.get<AdminUser[]>("/admin/users"),
      api.get<AdminRaffle[]>("/admin/raffles"),
      api.get<ContactMessage[]>("/admin/contact-messages"),
    ])
    stats.value = s
    users.value = u
    raffles.value = r
    contactMessages.value = c
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
    sendEvent("admin_subscription_updated", { target_id: userId, action })
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
      sendEvent("admin_raffle_drawn", { raffle_id: raffleId })
    } else {
      await api.patch(`/raffles/${raffleId}/cancel`)
      sendEvent("admin_raffle_cancelled", { raffle_id: raffleId })
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
      <div v-if="stats" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 md:gap-4 mb-6">
        <div class="bg-white rounded-xl border border-gray-200 p-3 md:p-4 hover:shadow-md transition-shadow cursor-pointer" @click="activeTab = 'users'">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-indigo-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"/></svg>
            </div>
            <div class="min-w-0">
              <p class="text-xs text-gray-500 uppercase font-medium truncate">Usuários</p>
              <p class="text-lg font-bold text-gray-900">{{ stats.totalUsers }}</p>
              <p class="text-[11px] text-green-600">{{ stats.activeUsers }} ativos</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-3 md:p-4 hover:shadow-md transition-shadow cursor-pointer" @click="activeTab = 'users'">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-amber-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase font-medium">Trial</p>
              <p class="text-lg font-bold text-amber-600">{{ stats.trialUsers }}</p>
              <p class="text-[11px] text-gray-500">{{ stats.pastDueUsers }} vencidos</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-3 md:p-4 hover:shadow-md transition-shadow cursor-pointer" @click="activeTab = 'raffles'">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-purple-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase font-medium">Rifas</p>
              <p class="text-lg font-bold text-gray-900">{{ stats.totalRaffles }}</p>
              <p class="text-[11px] text-indigo-600">{{ stats.activeRaffles }} ativas</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-3 md:p-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-green-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase font-medium">Ingressos</p>
              <p class="text-lg font-bold text-gray-900">{{ stats.totalPaidTickets }}</p>
              <p class="text-[11px] text-gray-500">pagos</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-3 md:p-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-emerald-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase font-medium">Receita</p>
              <p class="text-lg font-bold text-gray-900">R$ {{ (stats.totalRevenue / 100).toFixed(2) }}</p>
              <p class="text-[11px] text-gray-500">total</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-3 md:p-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-rose-100 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-rose-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 uppercase font-medium">Inativos</p>
              <p class="text-lg font-bold text-gray-900">{{ stats.totalUsers - stats.activeUsers }}</p>
              <p class="text-[11px] text-gray-500">sem assinatura</p>
            </div>
          </div>
        </div>
      </div>

      <div class="border-b border-gray-200 mb-6">
        <nav class="flex gap-1" role="tablist">
          <button @click="activeTab = 'users'" role="tab"
            class="px-5 py-3 text-sm font-medium rounded-t-lg transition-all"
            :class="activeTab === 'users'
              ? 'text-indigo-700 bg-white border border-b-0 border-gray-200 -mb-px shadow-sm'
              : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'"
          >
            <div class="flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"/></svg>
              Usuários
              <span class="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded-full">{{ users.length }}</span>
            </div>
          </button>
          <button @click="activeTab = 'raffles'" role="tab"
            class="px-5 py-3 text-sm font-medium rounded-t-lg transition-all"
            :class="activeTab === 'raffles'
              ? 'text-indigo-700 bg-white border border-b-0 border-gray-200 -mb-px shadow-sm'
              : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'"
          >
            <div class="flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/></svg>
              Rifas
              <span class="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded-full">{{ raffles.length }}</span>
            </div>
          </button>
          <button @click="activeTab = 'contact'" role="tab"
            class="px-5 py-3 text-sm font-medium rounded-t-lg transition-all"
            :class="activeTab === 'contact'
              ? 'text-indigo-700 bg-white border border-b-0 border-gray-200 -mb-px shadow-sm'
              : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'"
          >
            <div class="flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"/></svg>
              Mensagens
              <span v-if="contactMessages.length" class="text-xs bg-indigo-100 text-indigo-600 px-1.5 py-0.5 rounded-full">{{ contactMessages.length }}</span>
            </div>
          </button>
        </nav>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200">
        <div v-if="activeTab === 'users'">
          <div class="p-4 border-b border-gray-100 flex flex-col sm:flex-row gap-3">
            <input v-model="searchQuery" type="text" placeholder="Buscar por nome, email ou telefone..."
              class="w-full sm:max-w-xs px-4 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none" />
            <div class="flex gap-2 flex-wrap">
              <button @click="searchQuery = ''; userSort = 'createdAt'; userSortDir = 'desc'"
                class="px-3 py-2 text-xs font-medium text-gray-500 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors"
              >Limpar filtros</button>
            </div>
            <div class="sm:ml-auto flex items-center gap-2 text-xs text-gray-500">
              <span>Ordenar:</span>
              <select v-model="userSort" class="text-xs border border-gray-200 rounded-lg px-2 py-1.5 outline-none">
                <option value="createdAt">Registro</option>
                <option value="name">Nome</option>
                <option value="email">Email</option>
                <option value="subscriptionStatus">Assinatura</option>
              </select>
              <button @click="userSortDir = userSortDir === 'asc' ? 'desc' : 'asc'" class="p-1 hover:bg-gray-100 rounded" :title="userSortDir === 'asc' ? 'Crescente' : 'Decrescente'">
                <svg class="w-4 h-4 text-gray-400" :class="userSortDir === 'desc' ? 'rotate-180' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
            </div>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead class="bg-gray-50">
                <tr>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Nome</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Email</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Telefone</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Tipo</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Assinatura</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Expira</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Registro</th>
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
                    <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold" :class="subColor(u.subscriptionStatus)">
                      {{ subLabel(u.subscriptionStatus) }}
                      <span v-if="u.subscriptionIsTrial" class="text-amber-500">(teste)</span>
                    </span>
                  </td>
                  <td class="px-5 py-3 text-gray-500 text-xs">{{ u.subscriptionExpiresAt ? formatDate(u.subscriptionExpiresAt) : '-' }}</td>
                  <td class="px-5 py-3 text-gray-500 text-xs">{{ formatDate(u.createdAt) }}</td>
                  <td class="px-5 py-3">
                    <div class="flex gap-1 flex-wrap">
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

        <div v-if="activeTab === 'raffles'">
          <div class="p-4 border-b border-gray-100 flex flex-col sm:flex-row gap-3">
            <div class="sm:ml-auto flex items-center gap-2 text-xs text-gray-500">
              <span>Ordenar:</span>
              <select v-model="raffleSort" class="text-xs border border-gray-200 rounded-lg px-2 py-1.5 outline-none">
                <option value="createdAt">Data</option>
                <option value="title">Título</option>
                <option value="organizerName">Criador</option>
                <option value="status">Status</option>
                <option value="soldTickets">Vendas</option>
              </select>
              <button @click="raffleSortDir = raffleSortDir === 'asc' ? 'desc' : 'asc'" class="p-1 hover:bg-gray-100 rounded" :title="raffleSortDir === 'asc' ? 'Crescente' : 'Decrescente'">
                <svg class="w-4 h-4 text-gray-400" :class="raffleSortDir === 'desc' ? 'rotate-180' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
            </div>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead class="bg-gray-50">
                <tr>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Título</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Criador</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Status</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Valor</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Vendas</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Receita</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Data</th>
                  <th class="text-left px-5 py-3 font-semibold text-gray-600">Ações</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-for="r in sortedRaffles" :key="r.id" class="hover:bg-gray-50">
                  <td class="px-5 py-3">
                    <router-link :to="`/raffles/${r.id}`" class="font-medium text-indigo-600 hover:text-indigo-800 hover:underline max-w-[200px] truncate block" :title="r.title">
                      {{ r.title }}
                    </router-link>
                  </td>
                  <td class="px-5 py-3 text-gray-600 whitespace-nowrap">{{ r.organizerName }}</td>
                  <td class="px-5 py-3">
                    <span class="inline-flex px-2.5 py-1 rounded-full text-xs font-semibold" :class="statusColor(r.status)">
                      {{ statusLabel(r.status) }}
                    </span>
                  </td>
                  <td class="px-5 py-3 text-gray-700 whitespace-nowrap">R$ {{ (r.ticketPrice / 100).toFixed(2) }}</td>
                  <td class="px-5 py-3 min-w-[140px]">
                    <div class="flex items-center gap-2">
                      <div class="flex-1 min-w-[60px] bg-gray-100 rounded-full h-2">
                        <div class="h-2 rounded-full transition-all" :class="r.status === 'ACTIVE' ? 'bg-indigo-500' : r.status === 'DRAWN' ? 'bg-blue-500' : 'bg-gray-300'"
                          :style="{ width: salesProgress(r) + '%' }"></div>
                      </div>
                      <span class="text-xs text-gray-600 whitespace-nowrap">{{ r.paidTickets }}/{{ r.maxNumbers }}</span>
                    </div>
                  </td>
                  <td class="px-5 py-3 text-gray-700 whitespace-nowrap">R$ {{ (r.revenue / 100).toFixed(2) }}</td>
                  <td class="px-5 py-3 text-gray-500 text-xs whitespace-nowrap">{{ formatDate(r.createdAt) }}</td>
                  <td class="px-5 py-3">
                    <div class="flex gap-1 flex-wrap" v-if="r.status === 'ACTIVE'">
                      <router-link :to="`/dashboard/raffles/${r.id}/edit?redirect=/admin`"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg bg-indigo-50 text-indigo-600 hover:bg-indigo-100 transition-colors"
                      >Editar</router-link>
                      <button @click="confirmAction = { id: r.id, action: 'draw', label: 'realizar sorteio' }" :disabled="actionLoading === r.id || r.paidTickets === 0"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 disabled:opacity-40 transition-colors"
                        :title="r.paidTickets === 0 ? 'Nenhum ingresso pago' : ''"
                      >Sortear</button>
                      <button @click="confirmAction = { id: r.id, action: 'cancel', label: 'cancelar' }" :disabled="actionLoading === r.id"
                        class="px-2.5 py-1 text-xs font-medium rounded-lg bg-red-50 text-red-600 hover:bg-red-100 disabled:opacity-50 transition-colors"
                      >Cancelar</button>
                    </div>
                    <div v-else-if="r.status === 'DRAWN'" class="flex items-center gap-2 text-xs text-blue-600">
                      <span class="font-medium">Ganhador: {{ r.winnerNumber ?? '-' }}</span>
                      <router-link :to="`/raffles/${r.id}/result`" class="underline hover:text-blue-800">Ver resultado</router-link>
                    </div>
                    <span v-else class="text-xs text-gray-400">-</span>
                  </td>
                </tr>
              </tbody>
            </table>
            <p v-if="raffles.length === 0" class="text-center py-8 text-gray-500">Nenhuma rifa encontrada</p>
          </div>
        </div>

        <div v-if="activeTab === 'contact'" class="divide-y divide-gray-100">
          <div v-for="m in contactMessages" :key="m.id" class="p-5 hover:bg-gray-50 transition-colors">
            <div class="flex items-start justify-between gap-4 mb-2">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center text-white text-sm font-bold shrink-0">
                  {{ (m.name || 'A').charAt(0).toUpperCase() }}
                </div>
                <div>
                  <p class="font-semibold text-gray-900">{{ m.name }}</p>
                  <p class="text-sm text-indigo-600">{{ m.contact }}</p>
                </div>
              </div>
              <span class="text-xs text-gray-400 shrink-0">{{ new Date(m.createdAt).toLocaleString("pt-BR") }}</span>
            </div>
            <div class="ml-12">
              <p class="text-sm text-gray-700 whitespace-pre-wrap leading-relaxed">{{ m.message }}</p>
            </div>
          </div>
          <p v-if="contactMessages.length === 0" class="text-center py-12 text-gray-400">Nenhuma mensagem recebida</p>
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
