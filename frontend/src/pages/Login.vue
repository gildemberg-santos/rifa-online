<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "../stores/auth"

const auth = useAuthStore()
const router = useRouter()

const email = ref("")
const password = ref("")
const error = ref("")
const loading = ref(false)

async function submit() {
  loading.value = true
  error.value = ""
  try {
    await auth.login(email.value, password.value)
    router.push("/dashboard")
  } catch (e: any) {
    error.value = e.message || "Erro ao fazer login"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-md mx-auto px-4 py-12">
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <h1 class="text-2xl font-bold text-gray-900 mb-6 text-center">Entrar</h1>

      <form @submit.prevent="submit" class="space-y-4">
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
          <label class="block text-sm font-medium text-gray-700 mb-1">Senha</label>
          <input
            v-model="password"
            type="password"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-indigo-600 text-white py-2 rounded-lg hover:bg-indigo-700 disabled:opacity-50"
        >
          {{ loading ? "Entrando..." : "Entrar" }}
        </button>
      </form>

      <p class="text-sm text-center mt-4 text-gray-500">
        Não tem conta?
        <router-link to="/register" class="text-indigo-600 hover:underline">Cadastre-se</router-link>
      </p>
    </div>
  </div>
</template>
