<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { api } from "../utils/api"
import { useAuthStore } from "../stores/auth"
import { sendEvent } from "../utils/analytics"

const auth = useAuthStore()
const loading = ref(false)
const creatingCheckout = ref(false)
const subscriptionStatus = ref("")
const subscriptionExpiresAt = ref("")
const subscriptionIsTrial = ref(false)
const hasSubscriptionBefore = ref(false)
const error = ref("")
const trialActivated = ref(false)

onMounted(async () => {
  await loadStatus()
})

async function loadStatus() {
  loading.value = true
  try {
    const data = await api.get<{ subscriptionStatus: string; subscriptionExpiresAt?: string; subscriptionIsTrial?: boolean; hasSubscriptionBefore?: boolean }>("/subscription/status")
    subscriptionStatus.value = data.subscriptionStatus
    subscriptionExpiresAt.value = data.subscriptionExpiresAt || ""
    subscriptionIsTrial.value = data.subscriptionIsTrial ?? false
    hasSubscriptionBefore.value = data.hasSubscriptionBefore ?? false
  } catch {
    subscriptionStatus.value = auth.user?.subscriptionStatus || "INACTIVE"
  } finally {
    loading.value = false
  }
}

async function createCheckout() {
  creatingCheckout.value = true
  error.value = ""
  trialActivated.value = false

  const checkoutPath = import.meta.env.DEV ? "/subscription/dev-checkout" : "/subscription/checkout"

  try {
    const result = await api.post<{ checkoutUrl?: string; isTrial: boolean }>(checkoutPath)

    if (result.isTrial) {
      sendEvent("subscription_trial_started")
      trialActivated.value = true
      subscriptionStatus.value = "ACTIVE"
      subscriptionIsTrial.value = true
      hasSubscriptionBefore.value = true
      const expires = new Date()
      expires.setDate(expires.getDate() + 7)
      subscriptionExpiresAt.value = expires.toISOString()
      if (auth.user) {
        auth.user.subscriptionStatus = "ACTIVE"
        auth.user.subscriptionIsTrial = true
      }
    } else if (result.checkoutUrl) {
      sendEvent("subscription_checkout_started")
      window.location.href = result.checkoutUrl
    }
  } catch (e: any) {
    const msg = e?.message || ""
    if (msg.includes("já possui uma assinatura ativa")) {
      error.value = "Sua assinatura já está ativa. Renove quando estiver próxima do vencimento."
    } else {
      error.value = "Erro ao criar checkout. Tente novamente."
    }
  } finally {
    creatingCheckout.value = false
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ""
  return new Date(dateStr).toLocaleDateString("pt-BR")
}

const isExpired = computed(() => {
  if (!subscriptionExpiresAt.value) return false
  return new Date(subscriptionExpiresAt.value) < new Date()
})

const showTrialInfo = computed(() =>
  !hasSubscriptionBefore.value && (subscriptionStatus.value !== "ACTIVE" || isExpired.value),
)

const canCreateCheckout = computed(() =>
  subscriptionStatus.value !== "ACTIVE" || !hasSubscriptionBefore.value || isExpired.value,
)
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 animate-fade-in">
    <div class="text-center mb-8">
      <div class="inline-flex items-center justify-center w-12 h-12 sm:w-14 sm:h-14 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl sm:rounded-2xl shadow-lg mb-4">
        <svg class="w-6 h-6 sm:w-7 sm:h-7 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
        </svg>
      </div>
      <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Assinatura</h1>
      <p class="text-gray-500 text-sm mt-1">Gerencie sua assinatura mensal</p>
    </div>

    <div v-if="loading" class="flex justify-center py-8">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else>
      <div v-if="trialActivated" class="bg-white rounded-2xl shadow-sm border border-green-200 p-8 text-center">
        <div class="w-16 h-16 mx-auto bg-green-100 rounded-full flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-gray-900 mb-2">Teste Gratuito Ativado!</h2>
        <p class="text-gray-600 mb-2">Você tem 7 dias grátis para testar a plataforma.</p>
        <p class="text-sm text-gray-500 mb-6">Válido até {{ formatDate(subscriptionExpiresAt) }}</p>
        <router-link
          to="/dashboard"
          class="inline-flex items-center px-6 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md transition-all"
        >
          Ir para o Dashboard
        </router-link>
      </div>

      <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
        <div v-if="subscriptionIsTrial" class="mb-6 p-4 bg-amber-50 border border-amber-200 rounded-xl">
          <div class="flex items-start gap-3">
            <svg class="w-5 h-5 text-amber-600 mt-0.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <p class="text-sm font-semibold text-amber-800">Você está no período de teste do sistema</p>
              <p class="text-xs text-amber-700 mt-0.5">Aproveite para testar todas as funcionalidades.</p>
            </div>
          </div>
        </div>

        <div class="mb-6 p-4 rounded-xl border" :class="subscriptionStatus === 'ACTIVE' && !isExpired ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
          <p class="text-sm font-medium" :class="subscriptionStatus === 'ACTIVE' && !isExpired ? 'text-green-700' : 'text-red-700'">
            <template v-if="isExpired">
              Sua assinatura expirou
            </template>
            <template v-else-if="subscriptionStatus === 'ACTIVE'">
              Sua assinatura está ativa
            </template>
            <template v-else>
              Sua assinatura está inativa
            </template>
          </p>
          <p v-if="subscriptionExpiresAt" class="text-xs mt-1" :class="subscriptionStatus === 'ACTIVE' && !isExpired ? 'text-green-600' : 'text-red-600'">
            Válida até {{ formatDate(subscriptionExpiresAt) }}
          </p>
        </div>

        <div v-if="showTrialInfo" class="mb-6 p-4 bg-amber-50 border border-amber-200 rounded-xl">
          <p class="text-sm font-medium text-amber-800">Teste gratuito de 7 dias</p>
          <p class="text-xs text-amber-700 mt-1">Sem cobrança imediata. Após o período, assine por R$ 10,00/mês.</p>
        </div>

        <div class="mb-6">
          <div class="flex items-center justify-between py-3 border-b border-gray-100">
            <span class="text-gray-600">Plano</span>
            <span class="font-semibold text-gray-900">Assinatura Mensal</span>
          </div>
          <div class="flex items-center justify-between py-3 border-b border-gray-100">
            <span class="text-gray-600">Valor</span>
            <span class="font-semibold text-gray-900">R$ 10,00 /mês</span>
          </div>
          <div class="flex items-center justify-between py-3">
            <span class="text-gray-600">Recebimento</span>
            <img src="../assets/infinitepay-logo.svg" alt="InfinitePay" class="h-5" />
          </div>
        </div>

        <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2 mb-4">{{ error }}</p>

        <button
          v-if="canCreateCheckout"
          @click="createCheckout"
          :disabled="creatingCheckout"
          class="w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md hover:shadow-lg transition-all"
        >
          <span v-if="creatingCheckout" class="inline-flex items-center gap-2">
            <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            {{ showTrialInfo ? 'Ativando teste...' : 'Processando...' }}
          </span>
          <span v-else>
            <template v-if="showTrialInfo">Iniciar Teste Gratuito - 7 Dias</template>
            <template v-else>Assinar Agora - R$ 10,00</template>
          </span>
        </button>

        <div v-if="subscriptionStatus === 'ACTIVE' && !showTrialInfo && !isExpired" class="text-center">
          <p class="text-sm text-green-600 font-medium">Sua assinatura está vigente</p>
          <p v-if="subscriptionExpiresAt" class="text-xs text-gray-500 mt-1">Válida até {{ formatDate(subscriptionExpiresAt) }}</p>
        </div>
      </div>
    </template>
  </div>
</template>
