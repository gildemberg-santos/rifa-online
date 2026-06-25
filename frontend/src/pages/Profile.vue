<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"

interface User {
  id: string
  name: string
  email: string
  infinitePayHandle?: string
  subscriptionStatus: string
}

const user = ref<User | null>(null)
const name = ref("")
const email = ref("")
const password = ref("")
const infinitePayHandle = ref("")
const loading = ref(false)
const fetching = ref(true)
const message = ref("")
const isError = ref(false)
const handleMessage = ref("")
const handleLoading = ref(false)

onMounted(async () => {
  try {
    user.value = await api.get<User>("/me")
    name.value = user.value.name
    email.value = user.value.email
    infinitePayHandle.value = user.value.infinitePayHandle || ""
  } catch {
    message.value = "Erro ao carregar perfil"
    isError.value = true
  } finally {
    fetching.value = false
  }
})

async function submit() {
  loading.value = true
  message.value = ""
  isError.value = false

  const body: Record<string, string> = {}
  if (name.value !== user.value?.name) body.name = name.value
  if (email.value !== user.value?.email) body.email = email.value
  if (password.value) body.password = password.value

  if (Object.keys(body).length === 0) {
    message.value = "Nenhuma alteração"
    loading.value = false
    return
  }

  try {
    user.value = await api.put<User>("/me", body)
    password.value = ""
    message.value = "Perfil atualizado com sucesso"
  } catch (e: any) {
    message.value = e.message || "Erro ao atualizar perfil"
    isError.value = true
  } finally {
    loading.value = false
  }
}

async function saveHandle() {
  handleLoading.value = true
  handleMessage.value = ""

  try {
    const result = await api.put<{ infinitePayHandle: string }>("/me/infinite-pay-handle", {
      infinitePayHandle: infinitePayHandle.value,
    })
    if (user.value) user.value.infinitePayHandle = result.infinitePayHandle
    handleMessage.value = "Conta InfinitePay conectada com sucesso"
  } catch {
    handleMessage.value = "Erro ao salvar conta InfinitePay"
    isError.value = true
  } finally {
    handleLoading.value = false
  }
}

function subscriptionLabel(status: string): string {
  switch (status) {
    case "ACTIVE": return "Ativa"
    case "INACTIVE": return "Inativa"
    case "PAST_DUE": return "Vencida"
    case "CANCELLED": return "Cancelada"
    default: return status
  }
}

function subscriptionColor(status: string): string {
  switch (status) {
    case "ACTIVE": return "text-green-600 bg-green-50"
    case "INACTIVE":
    case "PAST_DUE": return "text-red-600 bg-red-50"
    default: return "text-gray-600 bg-gray-50"
  }
}
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 animate-fade-in">
    <div v-if="fetching" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else class="space-y-6">
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
        <div class="flex items-center gap-4 mb-6">
          <div class="w-14 h-14 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center text-white text-xl font-bold shadow-md">
            {{ (user?.name || "U").charAt(0).toUpperCase() }}
          </div>
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Meu Perfil</h1>
            <p class="text-sm text-gray-500">Gerencie seus dados</p>
          </div>
        </div>

        <div class="mb-6 flex items-center gap-2">
          <span class="text-sm text-gray-500">Assinatura:</span>
          <span class="inline-flex px-2.5 py-1 rounded-full text-xs font-semibold" :class="subscriptionColor(user?.subscriptionStatus || 'INACTIVE')">
            {{ subscriptionLabel(user?.subscriptionStatus || 'INACTIVE') }}
          </span>
          <router-link
            v-if="user?.subscriptionStatus !== 'ACTIVE'"
            to="/subscription"
            class="text-xs text-indigo-600 hover:text-indigo-700 font-medium ml-1"
          >
            Assinar agora
          </router-link>
        </div>

        <form @submit.prevent="submit" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Nome</label>
            <input
              v-model="name"
              type="text"
              required
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
            <input
              v-model="email"
              type="email"
              required
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">
              Nova senha <span class="text-gray-400 font-normal">(opcional)</span>
            </label>
            <input
              v-model="password"
              type="password"
              placeholder="Deixe em branco para manter"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <p
            v-if="message"
            class="text-sm px-3 py-2 rounded-lg"
            :class="isError ? 'text-red-600 bg-red-50' : 'text-green-600 bg-green-50'"
          >
            {{ message }}
          </p>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md transition-all"
          >
            <span v-if="loading" class="inline-flex items-center gap-2">
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              Salvando...
            </span>
            <span v-else>Salvar</span>
          </button>
        </form>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
        <h2 class="text-lg font-bold text-gray-900 mb-4">Conta InfinitePay</h2>
        <p class="text-sm text-gray-500 mb-4">
          Conecte sua conta InfinitePay para receber os pagamentos das rifas diretamente.
        </p>

        <form @submit.prevent="saveHandle" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Handle InfinitePay</label>
            <input
              v-model="infinitePayHandle"
              type="text"
              placeholder="Seu handle (ex: usuario123)"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <p
            v-if="handleMessage"
            class="text-sm px-3 py-2 rounded-lg"
            :class="isError ? 'text-red-600 bg-red-50' : 'text-green-600 bg-green-50'"
          >
            {{ handleMessage }}
          </p>

          <button
            type="submit"
            :disabled="handleLoading"
            class="w-full py-2.5 bg-gradient-to-r from-emerald-500 to-green-600 text-white font-semibold rounded-xl hover:from-emerald-600 hover:to-green-700 disabled:opacity-50 shadow-md transition-all"
          >
            <span v-if="handleLoading" class="inline-flex items-center gap-2">
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              Salvando...
            </span>
            <span v-else>Salvar</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
