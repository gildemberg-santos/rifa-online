<script setup lang="ts">
import { ref } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useAuthStore } from "../stores/auth"
import { useNotification } from "../composables/useNotification"

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const notify = useNotification()

const email = ref("")
const password = ref("")
const error = ref("")
const loading = ref(false)

const verifying = ref(false)
const code = ref("")

async function submit() {
  loading.value = true
  error.value = ""
  try {
    await auth.login(email.value, password.value)
    const redirect = route.query.redirect as string | undefined
    router.push(redirect || "/dashboard")
  } catch (e: any) {
    const msg = e.message || ""
    if (msg.includes("email not verified")) {
      verifying.value = true
      notify.show("Email nao verificado. Informe o codigo enviado no cadastro.", "info")
    } else {
      error.value = msg
    }
  } finally {
    loading.value = false
  }
}

async function verifyAndLogin() {
  loading.value = true
  error.value = ""
  try {
    await auth.verifyEmail(email.value, code.value)
    notify.show("Email verificado com sucesso!", "success")
    router.push("/dashboard")
  } catch (e: any) {
    error.value = e.message || "Codigo invalido"
  } finally {
    loading.value = false
  }
}

async function resendCode() {
  try {
    await auth.resendCode(email.value)
    notify.show("Codigo reenviado", "info")
  } catch {}
}
</script>

<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-md animate-scale-in">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-14 h-14 sm:w-16 sm:h-16 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-xl sm:rounded-2xl shadow-lg mb-4">
          <svg class="w-7 h-7 sm:w-8 sm:h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-gray-900">Bem-vindo</h1>
        <p class="text-gray-500 text-sm mt-1">Faça login para acessar o painel</p>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
        <!-- Login form -->
        <form v-if="!verifying" @submit.prevent="submit" class="space-y-5">
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
              maxlength="72"
              placeholder="••••••••"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md hover:shadow-lg transition-all"
          >
            <span v-if="loading" class="inline-flex items-center gap-2">
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              Entrando...
            </span>
            <span v-else>Entrar</span>
          </button>
        </form>

        <!-- Verification code -->
        <div v-else class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Codigo de verificacao</label>
            <p class="text-xs text-gray-400 mb-2">Enviamos um codigo para <strong>{{ email }}</strong> no cadastro</p>
            <input
              v-model="code"
              type="text"
              maxlength="6"
              placeholder="000000"
              class="w-full text-center text-2xl tracking-[0.5em] px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
            />
          </div>

          <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

          <button
            @click="verifyAndLogin"
            :disabled="loading || code.length !== 6"
            class="w-full py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md hover:shadow-lg transition-all"
          >
            <span v-if="loading" class="inline-flex items-center gap-2">
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              Verificando...
            </span>
            <span v-else>Verificar e Entrar</span>
          </button>

          <button
            @click="resendCode"
            class="w-full py-2.5 text-gray-600 font-medium bg-gray-100 hover:bg-gray-200 rounded-xl transition-colors"
          >
            Reenviar codigo
          </button>
        </div>

        <p class="text-center text-sm text-gray-500 mt-6">
          Não tem conta?
          <router-link to="/register" class="text-indigo-600 hover:text-indigo-700 font-medium">Criar conta</router-link>
        </p>
      </div>
    </div>
  </div>
</template>