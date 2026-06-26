<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"

interface Purchase {
  id: string
  type: "RAFFLE" | "SUBSCRIPTION"
  amount: number
  status: string
  buyerName: string
  createdAt: string
  paidAt?: string
  raffleId?: string
  raffleTitle?: string
  raffleStatus?: string
  ticketNumbers?: number[]
}

const items = ref<Purchase[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    items.value = await api.get<Purchase[]>("/me/purchases")
  } catch {
    console.error("Failed to load purchases")
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
    case "PAID": return "Pago"
    case "PENDING": return "Pendente"
    case "REFUNDED": return "Reembolsado"
    case "EXPIRED": return "Expirado"
    default: return s
  }
}

function statusColor(s: string): string {
  switch (s) {
    case "PAID": return "bg-green-100 text-green-700"
    case "PENDING": return "bg-yellow-100 text-yellow-700"
    case "REFUNDED": return "bg-red-100 text-red-700"
    case "EXPIRED": return "bg-gray-100 text-gray-600"
    default: return "bg-gray-100 text-gray-600"
  }
}

function typeIcon(type: string): string {
  return type === "RAFFLE" ? "🎫" : "⭐"
}

function typeLabel(type: string): string {
  return type === "RAFFLE" ? "Ingresso" : "Assinatura"
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8 animate-fade-in">
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Minhas Compras</h1>
      <p class="text-gray-500 text-sm mt-1">Histórico de todas as suas compras na plataforma</p>
    </div>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="items.length === 0" class="text-center py-16 bg-white rounded-2xl border border-gray-200">
      <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <p class="text-gray-500 text-lg">Nenhuma compra encontrada.</p>
      <p class="text-gray-400 text-sm mt-1">Você ainda não comprou ingressos ou assinaturas.</p>
      <router-link to="/" class="inline-flex items-center mt-6 px-5 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md transition-all text-sm">
        Ver rifas disponíveis
      </router-link>
    </div>

    <div v-else class="space-y-3 sm:space-y-4">
      <div v-for="p in items" :key="p.id" class="bg-white rounded-xl border border-gray-200 p-4 sm:p-5 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-2">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg shrink-0"
              :class="p.type === 'RAFFLE' ? 'bg-indigo-100' : 'bg-amber-100'"
            >
              {{ typeIcon(p.type) }}
            </div>
            <div class="min-w-0">
              <p class="font-semibold text-gray-900 truncate">
                <template v-if="p.type === 'RAFFLE'">{{ p.raffleTitle || p.buyerName }}</template>
                <template v-else>Assinatura Mensal</template>
              </p>
              <p class="text-xs text-gray-400">{{ typeLabel(p.type) }}</p>
            </div>
          </div>
          <span class="inline-flex self-start sm:self-auto px-2.5 py-1 rounded-full text-xs font-semibold shrink-0" :class="statusColor(p.status)">
            {{ statusLabel(p.status) }}
          </span>
        </div>

        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mt-3 pt-3 border-t border-gray-100">
          <div class="text-xs sm:text-sm text-gray-500 flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3">
            <span>{{ formatDate(p.createdAt) }}</span>
            <span v-if="p.paidAt" class="sm:ml-0">Pago em {{ formatDate(p.paidAt) }}</span>
          </div>
          <p class="font-bold text-gray-900 text-lg sm:text-base">R$ {{ (p.amount / 100).toFixed(2) }}</p>
        </div>

        <div v-if="p.type === 'RAFFLE' && p.ticketNumbers && p.ticketNumbers.length" class="mt-3">
          <div class="flex flex-wrap gap-1.5">
            <span
              v-for="n in p.ticketNumbers"
              :key="n"
              class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-indigo-50 text-indigo-700 border border-indigo-100"
            >#{{ n }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
