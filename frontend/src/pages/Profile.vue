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
  <div class="max-w-lg mx-auto px-4 py-12">
    <div v-if="fetching" class="text-center py-12 text-gray-500">Carregando...</div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">Meu Perfil</h1>

      <form @submit.prevent="submit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Nome</label>
          <input
            v-model="name"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            v-model="email"
            type="email"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Nova senha <span class="text-gray-400 font-normal">(deixe em branco para manter)</span>
          </label>
          <input
            v-model="password"
            type="password"
            placeholder="Nova senha"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <p
          v-if="message"
          class="text-sm"
          :class="isError ? 'text-red-600' : 'text-green-600'"
        >
          {{ message }}
        </p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-indigo-600 text-white py-2 rounded-lg hover:bg-indigo-700 disabled:opacity-50 font-medium"
        >
          {{ loading ? "Salvando..." : "Salvar" }}
        </button>
      </form>
    </div>
  </div>
</template>