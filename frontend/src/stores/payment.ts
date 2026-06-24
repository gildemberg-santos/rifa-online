import { defineStore } from "pinia"
import { ref } from "vue"
import { api } from "../utils/api"

interface Payment {
  id: string
  raffleId: string
  buyerName: string
  buyerEmail: string
  amount: number
  status: string
}

export const usePaymentStore = defineStore("payment", () => {
  const payments = ref<Payment[]>([])
  const loading = ref(false)

  async function fetchMyPayments(email: string) {
    loading.value = true
    try {
      payments.value = await api.get<Payment[]>(`/payments/my?email=${encodeURIComponent(email)}`)
    } finally {
      loading.value = false
    }
  }

  return { payments, loading, fetchMyPayments }
})
