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
const buyerEmail = ref("")
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
  if (!buyerName.value || !buyerEmail.value || !isValidPhone(buyerPhone.value)) return
  loading.value = true
  error.value = ""

  try {
    const result = await api.post<{ checkoutUrl: string }>(`/raffles/${raffleId}/checkout`, {
      numbers: numbers.value,
      buyerName: buyerName.value,
      buyerEmail: buyerEmail.value,
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
  <div class="max-w-lg mx-auto px-4 py-12">
    <h1 class="text-2xl font-bold text-gray-900 mb-6">Finalizar Compra</h1>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
      <div class="mb-4">
        <p class="text-sm text-gray-500">Números selecionados:</p>
        <p class="text-lg font-medium">{{ numbers.join(", ") }}</p>
      </div>

      <form @submit.prevent="submit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Nome</label>
          <input
            v-model="buyerName"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            v-model="buyerEmail"
            type="email"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Telefone</label>
          <input
            :value="buyerPhone"
            @input="onPhoneInput"
            type="tel"
            inputmode="numeric"
            required
            placeholder="(11) 99999-9999"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
          <p v-if="buyerPhone && !isValidPhone(buyerPhone)" class="text-xs text-red-500 mt-1">
            Telefone inválido. Informe DDD + número com 11 dígitos.
          </p>
        </div>

        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading || (buyerPhone.length > 0 && !isValidPhone(buyerPhone))"
          class="w-full bg-indigo-600 text-white py-3 rounded-lg font-medium hover:bg-indigo-700 disabled:opacity-50"
        >
          {{ loading ? "Processando..." : "Ir para Pagamento" }}
        </button>
      </form>
    </div>
  </div>
</template>
