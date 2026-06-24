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
const loading = ref(false)
const error = ref("")

async function submit() {
  if (!buyerName.value || !buyerEmail.value) return
  loading.value = true
  error.value = ""

  try {
    const result = await api.post<{ checkoutUrl: string }>(`/raffles/${raffleId}/checkout`, {
      numbers: numbers.value,
      buyerName: buyerName.value,
      buyerEmail: buyerEmail.value,
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

        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-indigo-600 text-white py-3 rounded-lg font-medium hover:bg-indigo-700 disabled:opacity-50"
        >
          {{ loading ? "Processando..." : "Ir para Pagamento" }}
        </button>
      </form>
    </div>
  </div>
</template>
