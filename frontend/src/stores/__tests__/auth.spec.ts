import { describe, it, expect, beforeEach, vi } from "vitest"
import { setActivePinia, createPinia } from "pinia"
import { useAuthStore } from "../auth"

// O refresh token vive apenas num cookie HttpOnly gerenciado pelo backend;
// o cliente só guarda o access token e os dados do usuário.

describe("AuthStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it("starts as not authenticated", () => {
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
  })

  it("sets the access token and marks as authenticated", () => {
    const store = useAuthStore()
    store.setTokens("access123")
    expect(store.accessToken).toBe("access123")
    expect(store.isAuthenticated).toBe(true)
    expect(localStorage.getItem("accessToken")).toBe("access123")
  })

  it("clears the session on logout", () => {
    const store = useAuthStore()
    store.setTokens("access")
    store.logout()
    expect(store.isAuthenticated).toBe(false)
    expect(store.accessToken).toBeNull()
    expect(localStorage.getItem("accessToken")).toBeNull()
  })

  it("hydrates the access token from localStorage on boot", async () => {
    // O módulo de sessão hidrata na carga; resetModules força reimportar
    // após gravar o token, simulando o boot da aplicação.
    localStorage.setItem("accessToken", "stored-token")
    vi.resetModules()
    const { useAuthStore: freshUseAuthStore } = await import("../auth")
    setActivePinia(createPinia())
    const store = freshUseAuthStore()
    expect(store.accessToken).toBe("stored-token")
    expect(store.isAuthenticated).toBe(true)
  })
})
