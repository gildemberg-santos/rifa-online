<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { api } from "../utils/api"

const route = useRoute()
const router = useRouter()
const raffleId = route.params.id as string

const title = ref("")
const description = ref("")
const ticketPrice = ref(0)
const drawDate = ref("")
const error = ref("")
const loading = ref(false)
const fetching = ref(true)

onMounted(async () => {
  try {
    const detail = await api.get<any>(`/raffles/${raffleId}`)
    title.value = detail.raffle.title
    description.value = detail.raffle.description
    ticketPrice.value = detail.raffle.ticketPrice / 100
    drawDate.value = detail.raffle.drawDate
  } catch (e) {
    error.value = "Erro ao carregar rifa"
  } finally {
    fetching.value = false
  }
})

async function submit() {
  loading.value = true
  error.value = ""

  try {
    await api.put(`/raffles/${raffleId}`, {
      title: title.value,
      description: description.value,
      ticketPrice: Math.round(ticketPrice.value * 100),
      maxNumbers: 0,
      drawDate: new Date(drawDate.value).toISOString(),
    })
    router.push("/dashboard")
  } catch (e: any) {
    error.value = e.message || "Erro ao atualizar rifa"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-lg mx-auto px-4 py-12">
    <div v-if="fetching" class="text-center py-12 text-gray-500">Carregando...</div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">Editar Rifa</h1>

      <form @submit.prevent="submit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Título</label>
          <input
            v-model="title"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Descrição</label>
          <textarea
            v-model="description"
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          ></textarea>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Valor do número (R$)</label>
          <input
            v-model.number="ticketPrice"
            type="number"
            step="0.01"
            min="0.01"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Data do sorteio</label>
          <input
            v-model="drawDate"
            type="datetime-local"
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
          {{ loading ? "Salvando..." : "Salvar" }}
        </button>
      </form>
    </div>
  </div>
</template>
