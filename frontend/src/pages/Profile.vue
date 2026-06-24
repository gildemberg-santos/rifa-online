<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"

interface User {
  id: string
  name: string
  email: string
}

const user = ref<User | null>(null)
const name = ref("")
const email = ref("")
const password = ref("")
const loading = ref(false)
const fetching = ref(true)
const message = ref("")
const isError = ref(false)

onMounted(async () => {
  try {
    user.value = await api.get<User>("/me")
    name.value = user.value.name
    email.value = user.value.email
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
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 animate-fade-in">
    <div v-if="fetching" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
      <div class="flex items-center gap-4 mb-6">
        <div class="w-14 h-14 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center text-white text-xl font-bold shadow-md">
          {{ (user?.name || "U").charAt(0).toUpperCase() }}
        </div>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Meu Perfil</h1>
          <p class="text-sm text-gray-500">Gerencie seus dados</p>
        </div>
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
  </div>
</template>