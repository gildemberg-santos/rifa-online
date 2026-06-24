<script setup lang="ts">
import { ref, computed } from "vue"
import { useRoute } from "vue-router"
import { api } from "../utils/api"

const route = useRoute()

const raffleId = route.params.id as string
const numbers = computed(() =>
  (route.query.numbers as string || "").split(",").map(Number).filter(Boolean),
)

const buyerName = ref("")
const buyerPhone = ref("")
const loading = ref(false)
const error = ref("")

function formatPhone(value: string): string {
  const digits = value.replace(/\D/g, "").slice(0, 11)
  if (digits.length <= 2) return `(${digits}`
  if (digits.length <= 7) return `(${digits.slice(0, 2)}) ${digits.slice(2)}`
  return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`
}

function onPhoneInput(e: Event) {
  const input = e.target as HTMLInputElement
  const cursor = input.selectionStart ?? 0
  const prevLen = buyerPhone.value.length
  buyerPhone.value = formatPhone(input.value)
  const newLen = buyerPhone.value.length
  if (input.setSelectionRange) {
    input.setSelectionRange(cursor + (newLen - prevLen), cursor + (newLen - prevLen))
  }
}

function isValidPhone(phone: string): boolean {
  return phone.replace(/\D/g, "").length === 11
}

async function submit() {
  if (!buyerName.value || !isValidPhone(buyerPhone.value)) return
  loading.value = true
  error.value = ""

  try {
    const result = await api.post<{ checkoutUrl: string }>(`/raffles/${raffleId}/checkout`, {
      numbers: numbers.value,
      buyerName: buyerName.value,
      buyerEmail: "",
      buyerPhone: buyerPhone.value.replace(/\D/g, ""),
    })
    window.location.href = result.checkoutUrl
  } catch (e: any) {
    error.value = e.message || "Erro ao criar checkout"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 animate-fade-in">
    <div class="text-center mb-8">
      <div class="inline-flex items-center justify-center w-14 h-14 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl shadow-lg mb-4">
        <svg class="w-7 h-7 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-gray-900">Finalizar Compra</h1>
      <p class="text-gray-500 text-sm mt-1">Preencha seus dados para continuar</p>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
      <div class="mb-6 p-4 bg-indigo-50 rounded-xl border border-indigo-100">
        <p class="text-sm text-gray-600">Números selecionados:</p>
        <p class="text-lg font-bold text-indigo-700 mt-1">{{ numbers.join(", ") }}</p>
      </div>

      <form @submit.prevent="submit" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">Nome completo</label>
          <input
            v-model="buyerName"
            type="text"
            required
            placeholder="Seu nome"
            class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">Telefone</label>
          <input
            :value="buyerPhone"
            @input="onPhoneInput"
            type="tel"
            inputmode="numeric"
            required
            placeholder="(11) 99999-9999"
            class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
          />
          <p v-if="buyerPhone && !isValidPhone(buyerPhone)" class="text-xs text-red-500 mt-1.5">
            Telefone inválido. Informe DDD + número com 11 dígitos.
          </p>
        </div>

        <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading || (buyerPhone.length > 0 && !isValidPhone(buyerPhone))"
          class="w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md hover:shadow-lg transition-all"
        >
          <span v-if="loading" class="inline-flex items-center gap-2">
            <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            Processando...
          </span>
          <span v-else>Ir para Pagamento</span>
        </button>
      </form>
    </div>
  </div>
</template>