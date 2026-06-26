<script setup lang="ts">
import { onMounted } from "vue"
import { sendEvent } from "../utils/analytics"

onMounted(() => {
  sendEvent("payment_pending")
})
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 text-center animate-fade-in">
    <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8 sm:p-10">
      <div class="w-16 h-16 sm:w-20 sm:h-20 mx-auto bg-amber-100 rounded-full flex items-center justify-center mb-6">
        <svg class="w-8 h-8 sm:w-10 sm:h-10 text-amber-600 animate-pulse" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h1 class="text-xl sm:text-2xl font-bold text-amber-600 mb-2">Aguardando Pagamento</h1>
      <p class="text-gray-500 mb-2">Sua compra está sendo processada.</p>
      <p class="text-sm text-gray-400 mb-8">Assim que o pagamento for confirmado, você receberá uma confirmação.</p>
      <div class="flex justify-center gap-3 mb-6">
        <span class="w-2 h-2 bg-amber-400 rounded-full animate-ping" style="animation-delay: 0s"></span>
        <span class="w-2 h-2 bg-amber-400 rounded-full animate-ping" style="animation-delay: 0.3s"></span>
        <span class="w-2 h-2 bg-amber-400 rounded-full animate-ping" style="animation-delay: 0.6s"></span>
      </div>
      <router-link
        to="/"
        class="inline-flex items-center px-5 sm:px-6 py-2.5 sm:py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md transition-all text-sm sm:text-base"
      >
        Voltar para Home
      </router-link>
    </div>
  </div>
</template>
