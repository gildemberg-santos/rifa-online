import { describe, it, expect, beforeEach, vi } from "vitest"
import { setActivePinia, createPinia } from "pinia"
import { useRaffleStore } from "../raffle"

const mockGet = vi.hoisted(() => vi.fn())
vi.mock("../../utils/api", () => ({
  api: {
    get: mockGet,
  },
}))

describe("RaffleStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it("starts with empty raffles and loading=false", () => {
    const store = useRaffleStore()
    expect(store.activeRaffles).toEqual([])
    expect(store.loading).toBe(false)
  })

  it("fetchActive should update activeRaffles", async () => {
    const fakeRaffles = [
      { id: "1", title: "Test Raffle", description: "Desc", ticketPrice: 10, maxNumbers: 100, drawDate: "2025-01-01", status: "active" },
    ]
    mockGet.mockResolvedValue(fakeRaffles)

    const store = useRaffleStore()
    await store.fetchActive()

    expect(mockGet).toHaveBeenCalledWith("/raffles")
    expect(store.activeRaffles).toEqual(fakeRaffles)
    expect(store.loading).toBe(false)
  })
})
