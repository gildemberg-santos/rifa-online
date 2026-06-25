import { describe, it, expect, beforeEach, vi } from "vitest"
import { setActivePinia, createPinia } from "pinia"
import { usePaymentStore } from "../payment"

const mockGet = vi.hoisted(() => vi.fn())
vi.mock("../../utils/api", () => ({
  api: {
    get: mockGet,
  },
}))

describe("PaymentStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it("starts with empty payments and loading=false", () => {
    const store = usePaymentStore()
    expect(store.payments).toEqual([])
    expect(store.loading).toBe(false)
  })

  it("fetchMyPayments should update payments", async () => {
    const fakePayments = [
      { id: "1", raffleId: "r1", buyerName: "John", buyerEmail: "john@test.com", amount: 100, status: "paid" },
    ]
    mockGet.mockResolvedValue(fakePayments)

    const store = usePaymentStore()
    await store.fetchMyPayments("john@test.com")

    expect(mockGet).toHaveBeenCalledWith("/payments/my?email=john%40test.com")
    expect(store.payments).toEqual(fakePayments)
    expect(store.loading).toBe(false)
  })
})
