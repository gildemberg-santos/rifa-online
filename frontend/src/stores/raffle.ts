import { defineStore } from "pinia"
import { ref } from "vue"
import { api } from "../utils/api"

interface Raffle {
  id: string
  title: string
  description: string
  ticketPrice: number
  maxNumbers: number
  drawDate: string
  imageUrl?: string
  status: string
}

export const useRaffleStore = defineStore("raffle", () => {
  const activeRaffles = ref<Raffle[]>([])
  const loading = ref(false)

  async function fetchActive() {
    loading.value = true
    try {
      activeRaffles.value = await api.get<Raffle[]>("/raffles")
    } finally {
      loading.value = false
    }
  }

  return { activeRaffles, loading, fetchActive }
})
