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
  subscriptionExpiresAt?: string
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
const actionLoading = ref<string | null>(null)
const selectedUser = ref<AdminUser | null>(null)
const showUserModal = ref(false)
const userRaffles = ref<any[]>([])
const userPayments = ref<any[]>([])

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

function subLabel(s: string): string {
  switch (s) {
    case "ACTIVE": return "Ativa"
    case "INACTIVE": return "Inativa"
    case "PAST_DUE": return "Vencida"
    case "CANCELLED": return "Cancelada"
    default: return s
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

async function viewUserDetails(user: AdminUser) {
  selectedUser.value = user
  try {
    const data = await api.get<any>(`/admin/users/${user.id}`)
    userRaffles.value = data.raffles || []
    userPayments.value = data.payments || []
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
}
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8 animate-fade-in">
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Administração</h1>
      <p class="text-gray-500 text-sm mt-1">Gerencie todos os assinantes da plataforma</p>
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
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Expira</th>
                <th class="text-left px-5 py-3 font-semibold text-gray-600">Ações</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="u in users" :key="u.id" class="hover:bg-gray-50">
                <td class="px-5 py-3 font-medium text-gray-900">{{ u.name }}</td>
                <td class="px-5 py-3 text-gray-600">{{ u.email }}</td>
                <td class="px-5 py-3">
                  <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold" :class="subColor(u.subscriptionStatus)">
                    {{ subLabel(u.subscriptionStatus) }}
                    <span v-if="u.subscriptionIsTrial" class="text-amber-500">(teste)</span>
                  </span>
                </td>
                <td class="px-5 py-3 text-gray-500 text-xs">{{ u.subscriptionExpiresAt ? formatDate(u.subscriptionExpiresAt) : '-' }}</td>
                <td class="px-5 py-3">
                  <div class="flex gap-1.5 flex-wrap">
                    <button
                      @click="updateSubscription(u.id, 'activate')"
                      :disabled="actionLoading === u.id"
                      class="px-2.5 py-1 text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                      :class="u.subscriptionStatus === 'ACTIVE' ? 'bg-green-100 text-green-700 cursor-not-allowed' : 'bg-green-50 text-green-600 hover:bg-green-100'"
                    >Ativar</button>
                    <button
                      @click="updateSubscription(u.id, 'cancel')"
                      :disabled="actionLoading === u.id"
                      class="px-2.5 py-1 text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                      :class="u.subscriptionStatus === 'CANCELLED' ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-gray-50 text-gray-600 hover:bg-gray-100'"
                    >Cancelar</button>
                    <button
                      @click="updateSubscription(u.id, 'past_due')"
                      :disabled="actionLoading === u.id"
                      class="px-2.5 py-1 text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                      :class="u.subscriptionStatus === 'PAST_DUE' ? 'bg-red-100 text-red-400 cursor-not-allowed' : 'bg-red-50 text-red-600 hover:bg-red-100'"
                    >Vencer</button>
                    <button
                      @click="updateSubscription(u.id, 'inactive')"
                      :disabled="actionLoading === u.id"
                      class="px-2.5 py-1 text-xs font-medium rounded-lg text-gray-500 bg-gray-50 hover:bg-gray-100 disabled:opacity-50 transition-colors"
                    >Zerar</button>
                    <button
                      @click="viewUserDetails(u)"
                      class="px-2.5 py-1 text-xs font-medium rounded-lg text-indigo-600 bg-indigo-50 hover:bg-indigo-100 transition-colors"
                    >Detalhes</button>
                  </div>
                </td>
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

    <div v-if="showUserModal && selectedUser" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm" @click.self="closeModal">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-2xl mx-4 max-h-[80vh] overflow-y-auto animate-scale-in">
        <div class="px-6 py-5 border-b border-gray-100 flex items-center justify-between">
          <div>
            <h2 class="text-lg font-bold text-gray-900">{{ selectedUser.name }}</h2>
            <p class="text-sm text-gray-500">{{ selectedUser.email }}</p>
          </div>
          <button @click="closeModal" class="p-2 rounded-lg hover:bg-gray-100 transition-colors">
            <svg class="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-6 space-y-6">
          <div>
            <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Rifas</h3>
            <div v-if="userRaffles.length === 0" class="text-sm text-gray-400">Nenhuma rifa criada</div>
            <div v-for="r in userRaffles" :key="r.id" class="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
              <span class="text-sm font-medium text-gray-900">{{ r.title }}</span>
              <span class="text-xs px-2 py-0.5 rounded-full"
                :class="r.status === 'ACTIVE' ? 'bg-green-100 text-green-700' : r.status === 'DRAWN' ? 'bg-blue-100 text-blue-700' : 'bg-red-100 text-red-700'"
              >{{ r.status }}</span>
            </div>
          </div>

          <div>
            <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Pagamentos</h3>
            <div v-if="userPayments.length === 0" class="text-sm text-gray-400">Nenhum pagamento</div>
            <div v-for="p in userPayments" :key="p.id" class="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
              <span class="text-sm text-gray-700">
                <template v-if="p.type === 'SUBSCRIPTION'">Assinatura</template>
                <template v-else>Rifa</template>
                — R$ {{ (p.amount / 100).toFixed(2) }}
              </span>
              <span class="text-xs px-2 py-0.5 rounded-full"
                :class="p.status === 'PAID' ? 'bg-green-100 text-green-700' : p.status === 'PENDING' ? 'bg-yellow-100 text-yellow-700' : 'bg-gray-100 text-gray-600'"
              >{{ p.status }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
