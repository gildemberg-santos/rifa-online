import { describe, it, expect, beforeEach, vi } from "vitest"
import { api, ApiError } from "../api"

const BASE = "http://localhost:8081/api/v1"

describe("API Utility", () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it("api.get makes GET request with correct URL", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: "test" }),
    })
    globalThis.fetch = mockFetch

    const result = await api.get("/test")

    expect(mockFetch).toHaveBeenCalledWith(
      `${BASE}/test`,
      expect.objectContaining({ method: "GET" }),
    )
    expect(result).toEqual({ data: "test" })
  })

  it("api.post makes POST request with body", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ id: 1 }),
    })
    globalThis.fetch = mockFetch

    const result = await api.post("/test", { name: "test" })

    expect(mockFetch).toHaveBeenCalledWith(
      `${BASE}/test`,
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ name: "test" }),
      }),
    )
    expect(result).toEqual({ id: 1 })
  })

  it("includes Authorization header when token is in localStorage", async () => {
    localStorage.setItem("accessToken", "my-token")

    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({}),
    })
    globalThis.fetch = mockFetch

    await api.get("/secure")

    expect(mockFetch).toHaveBeenCalledWith(
      `${BASE}/secure`,
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: "Bearer my-token",
        }),
      }),
    )
  })

  it("ApiError class works correctly", () => {
    const error = new ApiError(404, "Not found")

    expect(error).toBeInstanceOf(Error)
    expect(error.status).toBe(404)
    expect(error.message).toBe("Not found")
  })
})
