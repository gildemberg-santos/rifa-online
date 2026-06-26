<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "../stores/auth"

const auth = useAuthStore()
const router = useRouter()

const name = ref("")
const email = ref("")
const password = ref("")
const acceptedTerms = ref(false)
const error = ref("")
const loading = ref(false)

async function submit() {
  if (!acceptedTerms.value) {
    error.value = "É necessário aceitar os Termos de Uso e a Política de Privacidade."
    return
  }
  loading.value = true
  error.value = ""
  try {
    await auth.register(name.value, email.value, password.value)
    router.push("/dashboard")
  } catch (e: any) {
    error.value = e.message || "Erro ao criar conta"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-md animate-scale-in">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-2xl shadow-lg mb-4">
          <svg class="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-gray-900">Criar Conta</h1>
        <p class="text-gray-500 text-sm mt-1">Cadastre-se para criar rifas</p>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
        <form @submit.prevent="submit" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Nome</label>
            <input
              v-model="name"
              type="text"
              required
              maxlength="100"
              placeholder="Seu nome"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
            <input
              v-model="email"
              type="email"
              required
              maxlength="255"
              placeholder="seu@email.com"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Senha</label>
            <input
              v-model="password"
              type="password"
              required
              minlength="6"
              maxlength="72"
              placeholder="Mínimo 6 caracteres"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <label class="flex items-start gap-2 text-sm text-gray-600">
            <input
              v-model="acceptedTerms"
              type="checkbox"
              class="mt-0.5 h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
            />
            <span>
              Li e aceito os
              <router-link to="/termos-de-uso" target="_blank" class="text-indigo-600 hover:text-indigo-700 font-medium">Termos de Uso</router-link>,
              o
              <router-link to="/termo-do-organizador" target="_blank" class="text-indigo-600 hover:text-indigo-700 font-medium">Termo do Organizador</router-link>
              e a
              <router-link to="/politica-de-privacidade" target="_blank" class="text-indigo-600 hover:text-indigo-700 font-medium">Política de Privacidade</router-link>.
            </span>
          </label>

          <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

          <button
            type="submit"
            :disabled="loading || !acceptedTerms"
            class="w-full py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md hover:shadow-lg transition-all"
          >
            <span v-if="loading" class="inline-flex items-center gap-2">
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              Criando...
            </span>
            <span v-else>Criar Conta</span>
          </button>
        </form>

        <p class="text-center text-sm text-gray-500 mt-6">
          Já tem conta?
          <router-link to="/login" class="text-indigo-600 hover:text-indigo-700 font-medium">Fazer login</router-link>
        </p>
      </div>
    </div>
  </div>
</template>
