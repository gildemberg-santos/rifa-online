<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"

interface Raffle {
  id: string
  title: string
  ticketPrice: number
  maxNumbers: number
  status: string
  drawDate: string
}
const raffles = ref<Raffle[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    raffles.value = await api.get<Raffle[]>("/raffles/my")
  } catch (e) {
    console.error("Failed to load my raffles", e)
  } finally {
    loading.value = false
  }
})

async function drawRaffle(id: string) {
  if (!confirm("Realizar sorteio?")) return
  try {
    await api.post(`/raffles/${id}/draw`)
    const idx = raffles.value.findIndex((r) => r.id === id)
    if (idx !== -1) raffles.value[idx].status = "DRAWN"
  } catch (e: any) {
    alert(e.message || "Erro ao realizar sorteio")
  }
}
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-2xl font-bold">Meu Dashboard</h1>
      <router-link
        to="/dashboard/raffles/new"
        class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 text-sm"
      >
        Criar Rifa
      </router-link>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">
      Carregando...
    </div>

    <div v-else-if="raffles.length === 0" class="text-center py-12 text-gray-500">
      Você ainda não criou nenhuma rifa.
      <br />
      <router-link to="/dashboard/raffles/new" class="text-indigo-600 hover:underline mt-2 inline-block">
        Criar primeira rifa
      </router-link>
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50">
          <tr>
            <th class="text-left px-4 py-3 font-medium text-gray-600">Título</th>
            <th class="text-left px-4 py-3 font-medium text-gray-600">Status</th>
            <th class="text-left px-4 py-3 font-medium text-gray-600">Valor</th>
            <th class="text-left px-4 py-3 font-medium text-gray-600">Números</th>
            <th class="text-left px-4 py-3 font-medium text-gray-600">Ações</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="raffle in raffles" :key="raffle.id">
            <td class="px-4 py-3 font-medium">{{ raffle.title }}</td>
            <td class="px-4 py-3">
              <span
                class="px-2 py-1 rounded-full text-xs font-medium"
                :class="{
                  'bg-green-100 text-green-800': raffle.status === 'ACTIVE',
                  'bg-red-100 text-red-800': raffle.status === 'CANCELLED',
                  'bg-blue-100 text-blue-800': raffle.status === 'DRAWN',
                }"
              >
                {{ raffle.status }}
              </span>
            </td>
            <td class="px-4 py-3">R$ {{ (raffle.ticketPrice / 100).toFixed(2) }}</td>
            <td class="px-4 py-3">{{ raffle.maxNumbers }}</td>
            <td class="px-4 py-3 flex gap-2">
              <router-link
                :to="`/dashboard/raffles/${raffle.id}/edit`"
                v-if="raffle.status === 'ACTIVE'"
                class="text-indigo-600 hover:underline text-xs"
              >
                Editar
              </router-link>
              <button
                v-if="raffle.status === 'ACTIVE'"
                @click="drawRaffle(raffle.id)"
                class="text-green-600 hover:underline text-xs"
              >
                Sortear
              </button>
              <router-link
                :to="`/raffles/${raffle.id}/result`"
                v-if="raffle.status === 'DRAWN'"
                class="text-blue-600 hover:underline text-xs"
              >
                Resultado
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
