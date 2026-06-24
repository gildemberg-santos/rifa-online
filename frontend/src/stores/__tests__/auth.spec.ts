import { describe, it, expect, beforeEach } from "vitest"
import { setActivePinia, createPinia } from "pinia"
import { useAuthStore } from "../auth"

describe("AuthStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it("starts as not authenticated", () => {
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
  })

  it("sets tokens and marks as authenticated", () => {
    const store = useAuthStore()
    store.setTokens("access123", "refresh456")
    expect(store.accessToken).toBe("access123")
    expect(store.refreshToken).toBe("refresh456")
    expect(store.isAuthenticated).toBe(true)
    expect(localStorage.getItem("accessToken")).toBe("access123")
  })

  it("clears tokens on logout", () => {
    const store = useAuthStore()
    store.setTokens("access", "refresh")
    store.logout()
    expect(store.isAuthenticated).toBe(false)
    expect(store.accessToken).toBeNull()
    expect(localStorage.getItem("accessToken")).toBeNull()
  })

  it("reads token from localStorage on init", () => {
    localStorage.setItem("accessToken", "stored-token")
    localStorage.setItem("refreshToken", "stored-refresh")
    const store = useAuthStore()
    expect(store.accessToken).toBe("stored-token")
    expect(store.isAuthenticated).toBe(true)
  })
})
