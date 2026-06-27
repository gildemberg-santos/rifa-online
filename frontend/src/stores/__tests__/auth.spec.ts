import { describe, it, expect, beforeEach, vi } from "vitest"
import { setActivePinia, createPinia } from "pinia"
import { useAuthStore } from "../auth"

const { mockPost } = vi.hoisted(() => ({ mockPost: vi.fn() }))

vi.mock("../../utils/api", () => ({
  api: { post: mockPost },
}))

vi.mock("../../router", () => ({
  default: { push: vi.fn() },
}))

vi.mock("../../utils/analytics", () => ({
  sendEvent: vi.fn(),
}))

describe("AuthStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
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

  it("clears the session on logout", async () => {
    mockPost.mockResolvedValue({})
    const store = useAuthStore()
    store.setTokens("access")
    store.logout()
    await vi.waitFor(() => {
      expect(store.isAuthenticated).toBe(false)
      expect(store.accessToken).toBeNull()
    })
    expect(localStorage.getItem("accessToken")).toBeNull()
  })

  it("hydrates the access token from localStorage on boot", async () => {
    localStorage.setItem("accessToken", "stored-token")
    vi.resetModules()
    const { useAuthStore: freshUseAuthStore } = await import("../auth")
    setActivePinia(createPinia())
    const store = freshUseAuthStore()
    expect(store.accessToken).toBe("stored-token")
    expect(store.isAuthenticated).toBe(true)
  })

  it("register calls api.post with correct data", async () => {
    mockPost.mockResolvedValueOnce({})
    const store = useAuthStore()
    await store.register("John", "john@test.com", "123456")
    expect(mockPost).toHaveBeenCalledWith("/auth/register", {
      name: "John",
      email: "john@test.com",
      password: "123456",
    })
  })

  it("verifyEmail calls api.post and sets session", async () => {
    mockPost.mockResolvedValueOnce({
      accessToken: "token123",
      user: { id: "1", name: "John", email: "john@test.com" },
    })
    const store = useAuthStore()
    await store.verifyEmail("john@test.com", "123456")
    expect(mockPost).toHaveBeenCalledWith("/auth/verify-email", {
      email: "john@test.com",
      code: "123456",
    })
    expect(store.accessToken).toBe("token123")
    expect(store.isAuthenticated).toBe(true)
  })

  it("resendCode calls api.post with email", async () => {
    mockPost.mockResolvedValueOnce({})
    const store = useAuthStore()
    await store.resendCode("john@test.com")
    expect(mockPost).toHaveBeenCalledWith("/auth/resend-code", {
      email: "john@test.com",
    })
  })

  it("login calls api.post and sets session", async () => {
    mockPost.mockResolvedValueOnce({
      accessToken: "token456",
      user: { id: "1", name: "John", email: "john@test.com" },
    })
    const store = useAuthStore()
    await store.login("john@test.com", "123456")
    expect(mockPost).toHaveBeenCalledWith("/auth/login", {
      email: "john@test.com",
      password: "123456",
    })
    expect(store.accessToken).toBe("token456")
  })
})
