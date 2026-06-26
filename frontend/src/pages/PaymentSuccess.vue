<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { api } from "../utils/api"
import { sendEvent } from "../utils/analytics"

const route = useRoute()
const status = ref<"verifying" | "confirmed" | "error">("verifying")
const message = ref("")

onMounted(async () => {
  const paymentId = (route.query.paymentId as string) || ""
  if (!paymentId) {
    status.value = "error"
    message.value = "ID do pagamento não encontrado."
    return
  }

  try {
    await api.post(`/payments/${paymentId}/confirm`)
    status.value = "confirmed"
    sendEvent("payment_success", { payment_id: paymentId })
  } catch (e: any) {
    status.value = "error"
    message.value = e.message || "Não foi possível confirmar o pagamento. Ele será processado em breve."
    sendEvent("payment_pending", { payment_id: paymentId })
  }
})
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 text-center animate-fade-in">
    <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8 sm:p-10">
      <div v-if="status === 'verifying'" class="py-8">
        <div class="w-12 h-12 mx-auto border-4 border-indigo-600 border-t-transparent rounded-full animate-spin mb-4"></div>
        <p class="text-gray-600">Verificando pagamento...</p>
      </div>

      <template v-else-if="status === 'confirmed'">
        <div class="w-16 h-16 sm:w-20 sm:h-20 mx-auto bg-green-100 rounded-full flex items-center justify-center mb-6 animate-scale-in">
          <svg class="w-8 h-8 sm:w-10 sm:h-10 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h1 class="text-xl sm:text-2xl font-bold text-green-600 mb-2">Pagamento Confirmado!</h1>
        <p class="text-gray-500 mb-8">Seus números foram registrados com sucesso.</p>
      </template>

      <template v-else>
        <div class="w-16 h-16 sm:w-20 sm:h-20 mx-auto bg-amber-100 rounded-full flex items-center justify-center mb-6">
          <svg class="w-8 h-8 sm:w-10 sm:h-10 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>
        <h1 class="text-xl sm:text-2xl font-bold text-amber-600 mb-2">Pagamento em Processamento</h1>
        <p class="text-gray-500 mb-8">{{ message }}</p>
      </template>

      <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
        <router-link
          to="/"
          class="w-full sm:w-auto inline-flex items-center justify-center px-5 sm:px-6 py-2.5 sm:py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md transition-all text-sm sm:text-base"
        >
          <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Voltar para Home
        </router-link>
        <router-link
          to="/purchases"
          class="w-full sm:w-auto inline-flex items-center justify-center px-5 sm:px-6 py-2.5 sm:py-3 bg-white border border-gray-300 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 transition-all text-sm sm:text-base"
        >
          Ver minhas compras
        </router-link>
      </div>
    </div>
  </div>
</template>
