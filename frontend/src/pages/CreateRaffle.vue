<script setup lang="ts">
import { ref, computed } from "vue"
import { useRouter } from "vue-router"
import { api } from "../utils/api"
import { sendEvent } from "../utils/analytics"

const router = useRouter()

const title = ref("")
const description = ref("")
const ticketPrice = ref(0)
const maxNumbers = ref(100)
const drawDate = ref("")
const error = ref("")
const loading = ref(false)

const totalPotential = computed(() => {
  return (ticketPrice.value * maxNumbers.value * 100) / 100
})

const isFutureDate = computed(() => {
  if (!drawDate.value) return false
  return new Date(drawDate.value) > new Date()
})

const minDate = computed(() => {
  const now = new Date()
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset())
  return now.toISOString().slice(0, 16)
})

async function submit() {
  if (!isFutureDate.value) {
    error.value = "A data do sorteio deve ser futura"
    return
  }
  loading.value = true
  error.value = ""

  try {
    await api.post("/raffles", {
      title: title.value,
      description: description.value,
      ticketPrice: Math.round(ticketPrice.value * 100),
      maxNumbers: maxNumbers.value,
      drawDate: new Date(drawDate.value).toISOString(),
    })
    sendEvent("raffle_created", { title: title.value, max_numbers: maxNumbers.value })
    router.push("/dashboard")
  } catch (e: any) {
    const msg = e.message || ""
    if (msg.includes("subscription is not active")) {
      error.value = "⚠ Sua assinatura não está ativa. Você precisa de um plano mensal para criar rifas."
    } else {
      error.value = msg
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8 animate-fade-in">
    <div class="flex items-center gap-3 mb-8">
      <button @click="router.back()" class="p-2 rounded-xl hover:bg-gray-100 transition-colors">
        <svg class="w-5 h-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
      </button>
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Criar Nova Rifa</h1>
        <p class="text-gray-500 text-sm">Preencha os dados para criar uma nova rifa</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-5 gap-8">
      <div class="lg:col-span-3">
        <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
          <form @submit.prevent="submit" class="space-y-6">
            <div>
              <label class="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
                Título
              </label>
              <input
                v-model="title"
                type="text"
                required
                maxlength="100"
                placeholder="Ex: Rifa de Final de Ano"
                class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
              />
              <p class="text-xs text-gray-400 mt-1 text-right">{{ title.length }}/100</p>
            </div>

            <div>
              <label class="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h7" />
                </svg>
                Descrição
              </label>
              <textarea
                v-model="description"
                rows="3"
                maxlength="500"
                placeholder="Descreva o prêmio, as regras e informações importantes..."
                class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow resize-none"
              ></textarea>
              <p class="text-xs text-gray-400 mt-1 text-right">{{ description.length }}/500</p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                  <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Valor do número (R$)
                </label>
                <div class="relative">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 font-medium">R$</span>
                  <input
                    v-model.number="ticketPrice"
                    type="number"
                    step="0.01"
                    min="0.01"
                    required
                    placeholder="10,00"
                    class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
                  />
                </div>
              </div>
              <div>
                <label class="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                  <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M7 3l-3 3m4 0l-3-3m6 0l-3 3m4 0l-3-3m6 0l-3 3m4 0l-3-3" />
                  </svg>
                  Quantidade
                </label>
                <input
                  v-model.number="maxNumbers"
                  type="number"
                  min="1"
                  max="1000"
                  required
                  class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
                />
              </div>
            </div>

            <div>
              <label class="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                Data do sorteio
              </label>
              <input
                v-model="drawDate"
                type="datetime-local"
                :min="minDate"
                required
                class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
              />
              <p v-if="drawDate && !isFutureDate" class="text-xs text-red-500 mt-1">A data do sorteio deve ser no futuro</p>
            </div>

            <p v-if="error && !error.includes('assinatura')" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>
            <div v-if="error && error.includes('assinatura')" class="bg-amber-50 border border-amber-200 rounded-xl p-4">
              <p class="text-sm text-amber-800">{{ error }}</p>
              <router-link to="/subscription" class="mt-2 inline-block text-sm font-semibold text-indigo-600 hover:text-indigo-800 underline">
                Ver planos de assinatura →
              </router-link>
            </div>

            <div class="flex gap-3">
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
                  Criando...
                </span>
                <span v-else>Criar Rifa</span>
              </button>
            </div>
          </form>
        </div>
      </div>

      <div class="lg:col-span-2">
        <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6 lg:sticky lg:top-24">
          <h3 class="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">Prévia</h3>
          <div class="bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 rounded-xl h-32 flex items-center justify-center mb-4">
            <span class="text-5xl text-white font-extrabold drop-shadow-lg">{{ title ? title.charAt(0).toUpperCase() : "?" }}</span>
          </div>
          <h4 class="font-semibold text-gray-900 truncate">{{ title || "Título da Rifa" }}</h4>
          <p class="text-sm text-gray-500 mt-1 line-clamp-2">{{ description || "Descrição da rifa..." }}</p>
          <div class="flex justify-between items-center mt-4 pt-4 border-t border-gray-100">
            <div>
              <span class="text-xs text-gray-400">Preço</span>
              <p class="text-lg font-bold text-indigo-600">
                R$ {{ ticketPrice > 0 ? ticketPrice.toFixed(2) : "0,00" }}
              </p>
            </div>
            <div class="text-right">
              <span class="text-xs text-gray-400">Números</span>
              <p class="text-sm font-medium text-gray-700">{{ maxNumbers || 0 }}</p>
            </div>
          </div>
          <div class="mt-4 p-3 bg-indigo-50 rounded-xl">
            <span class="text-xs text-gray-500">Receita potencial total</span>
            <p class="text-lg font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600">
              R$ {{ totalPotential > 0 ? totalPotential.toFixed(2) : "0,00" }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>