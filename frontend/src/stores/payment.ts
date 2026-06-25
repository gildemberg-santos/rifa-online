import { defineStore } from "pinia"
import { ref } from "vue"
import { api } from "../utils/api"

interface Payment {
  id: string
  raffleId: string
  buyerName: string
  buyerPhone: string
  amount: number
  status: string
}

export const usePaymentStore = defineStore("payment", () => {
  const payments = ref<Payment[]>([])
  const loading = ref(false)

  async function fetchMyPayments(phone: string) {
    loading.value = true
    try {
      payments.value = await api.get<Payment[]>(`/payments/my?phone=${encodeURIComponent(phone)}`)
    } finally {
      loading.value = false
    }
  }

  return { payments, loading, fetchMyPayments }
})