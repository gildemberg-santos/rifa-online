<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { api } from "../utils/api"

const route = useRoute()
const router = useRouter()
const raffleId = route.params.id as string

const title = ref("")
const description = ref("")
const ticketPrice = ref(0)
const drawDate = ref("")
const error = ref("")
const loading = ref(false)
const fetching = ref(true)
const stats = ref<{ totalSold: number; totalRevenue: number; percentageSold: number } | null>(null)

onMounted(async () => {
  try {
    const [detail, s] = await Promise.all([
      api.get<any>(`/raffles/${raffleId}`),
      api.get<any>(`/raffles/${raffleId}/stats`).catch(() => null),
    ])
    title.value = detail.raffle.title
    description.value = detail.raffle.description
    ticketPrice.value = detail.raffle.ticketPrice / 100
    drawDate.value = detail.raffle.drawDate
    stats.value = s
  } catch (e) {
    error.value = "Erro ao carregar rifa"
  } finally {
    fetching.value = false
  }
})

async function submit() {
  loading.value = true
  error.value = ""

  try {
    await api.put(`/raffles/${raffleId}`, {
      title: title.value,
      description: description.value,
      ticketPrice: Math.round(ticketPrice.value * 100),
      maxNumbers: 0,
      drawDate: new Date(drawDate.value).toISOString(),
    })
    router.push("/dashboard")
  } catch (e: any) {
    error.value = e.message || "Erro ao atualizar rifa"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8 animate-fade-in">
    <div v-if="fetching" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else>
      <div class="flex items-center gap-3 mb-8">
        <button @click="router.back()" class="p-2 rounded-xl hover:bg-gray-100 transition-colors">
          <svg class="w-5 h-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Editar Rifa</h1>
          <p class="text-gray-500 text-sm">Atualize os dados da rifa</p>
        </div>
      </div>

      <div v-if="stats" class="grid grid-cols-3 gap-4 mb-8">
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Ingressos Vendidos</p>
          <p class="text-xl font-bold text-gray-900 mt-1">{{ stats.totalSold }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Receita</p>
          <p class="text-xl font-bold text-gray-900 mt-1">R$ {{ (stats.totalRevenue / 100).toFixed(2) }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider font-medium">Vendidos</p>
          <p class="text-xl font-bold text-gray-900 mt-1">{{ stats.percentageSold }}%</p>
        </div>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
        <div class="flex items-center gap-2 mb-5 pb-4 border-b border-gray-100">
          <div class="w-10 h-10 bg-gradient-to-br from-amber-500 to-orange-600 rounded-xl flex items-center justify-center text-white font-bold shrink-0">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <div>
            <p class="font-semibold text-gray-900">Dados da Rifa</p>
            <p class="text-xs text-gray-500">ID: {{ raffleId.substring(0, 8) }}...</p>
          </div>
        </div>

        <form @submit.prevent="submit" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Título</label>
            <input
              v-model="title"
              type="text"
              required
              maxlength="100"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Descrição</label>
            <textarea
              v-model="description"
              rows="3"
              maxlength="500"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow resize-none"
            ></textarea>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Valor do número (R$)</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 font-medium">R$</span>
                <input
                  v-model.number="ticketPrice"
                  type="number"
                  step="0.01"
                  min="0.01"
                  required
                  class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Quantidade</label>
              <input
                :value="stats?.totalSold ?? 0"
                type="text"
                disabled
                class="w-full px-4 py-2.5 border border-gray-200 rounded-xl bg-gray-50 text-gray-500 cursor-not-allowed"
              />
              <p class="text-xs text-gray-400 mt-1">Não é possível alterar após criar</p>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Data do sorteio</label>
            <input
              v-model="drawDate"
              type="datetime-local"
              required
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

          <div class="flex gap-3 pt-2">
            <router-link
              to="/dashboard"
              class="flex-1 py-2.5 text-center text-gray-700 font-medium bg-gray-100 hover:bg-gray-200 rounded-xl transition-colors"
            >
              Cancelar
            </router-link>
            <button
              type="submit"
              :disabled="loading"
              class="flex-[2] py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md transition-all"
            >
              <span v-if="loading" class="inline-flex items-center gap-2">
                <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                Salvando...
              </span>
              <span v-else>Salvar</span>
            </button>
          </div>
        </form>
      </div>
    </template>
  </div>
</template>