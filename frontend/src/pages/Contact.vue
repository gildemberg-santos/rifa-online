<script setup lang="ts">
import { ref } from "vue"
import { api } from "../utils/api"

const name = ref("")
const contact = ref("")
const message = ref("")
const loading = ref(false)
const error = ref("")
const sent = ref(false)

async function submit() {
  if (!name.value.trim() || message.value.trim().length < 10) {
    error.value = "Informe seu nome e uma mensagem com pelo menos 10 caracteres."
    return
  }
  loading.value = true
  error.value = ""
  try {
    await api.post("/contact", {
      name: name.value.trim(),
      contact: contact.value.trim(),
      message: message.value.trim(),
    })
    sent.value = true
  } catch (e: any) {
    error.value = e.message || "Não foi possível enviar a mensagem. Tente novamente."
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12 animate-fade-in">
    <div class="text-center mb-8">
      <div class="inline-flex items-center justify-center w-14 h-14 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl shadow-lg mb-4">
        <svg class="w-7 h-7 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-gray-900">Fale com a Rifa Online</h1>
      <p class="text-gray-500 text-sm mt-1">Envie sua mensagem, dúvida ou solicitação de privacidade (LGPD).</p>
    </div>

    <div v-if="sent" class="bg-white rounded-2xl shadow-sm border border-gray-200 p-10 text-center">
      <div class="w-16 h-16 mx-auto bg-green-100 rounded-full flex items-center justify-center mb-5">
        <svg class="w-8 h-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <h2 class="text-xl font-bold text-gray-900 mb-2">Mensagem enviada!</h2>
      <p class="text-gray-500">Recebemos sua mensagem e responderemos pelo contato informado, quando aplicável.</p>
    </div>

    <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8">
      <form @submit.prevent="submit" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">Nome</label>
          <input
            v-model="name"
            type="text"
            required
            maxlength="150"
            placeholder="Seu nome"
            class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            E-mail ou telefone para retorno <span class="text-gray-400 font-normal">(opcional)</span>
          </label>
          <input
            v-model="contact"
            type="text"
            maxlength="150"
            placeholder="Como podemos te responder"
            class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">Mensagem</label>
          <textarea
            v-model="message"
            required
            rows="5"
            maxlength="2000"
            placeholder="Escreva sua mensagem..."
            class="w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow resize-y"
          ></textarea>
        </div>

        <p v-if="error" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 shadow-md hover:shadow-lg transition-all"
        >
          <span v-if="loading" class="inline-flex items-center gap-2">
            <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            Enviando...
          </span>
          <span v-else>Enviar mensagem</span>
        </button>
      </form>
    </div>
  </div>
</template>
