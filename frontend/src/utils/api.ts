import router from "../router"
import { accessToken, refreshToken, setSession, clearSession } from "./session"

const API_URL = import.meta.env.VITE_API_URL || ""
const BASE_URL = API_URL ? `${API_URL}/api/v1` : "/api/v1"

class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message)
  }
}

function isAuthEndpoint(path: string): boolean {
  return path.startsWith("/auth/")
}

function redirectToLogin() {
  clearSession()
  const current = router.currentRoute.value
  if (current.name === "login") return
  router.push({ name: "login", query: { redirect: current.fullPath } })
}

// Renovação de token com "single-flight": chamadas concorrentes que recebem 401
// compartilham a mesma tentativa de refresh, evitando múltiplas chamadas a /auth/refresh.
let refreshing: Promise<boolean> | null = null

function refreshAccessToken(): Promise<boolean> {
  if (refreshing) return refreshing

  refreshing = (async () => {
    if (!refreshToken.value) return false
    try {
      const res = await fetch(`${BASE_URL}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refreshToken: refreshToken.value }),
      })
      if (!res.ok) return false
      const data = await res.json()
      if (!data?.accessToken) return false
      setSession({ accessToken: data.accessToken, refreshToken: data.refreshToken, user: data.user })
      return true
    } catch {
      return false
    }
  })()

  refreshing.finally(() => {
    refreshing = null
  })
  return refreshing
}

async function request<T>(method: string, path: string, body?: unknown, retried = false): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  }

  if (accessToken.value) {
    headers["Authorization"] = `Bearer ${accessToken.value}`
  }

  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  // Token expirado/inválido: tenta renovar uma vez; se não conseguir, vai para o login.
  if (res.status === 401 && !isAuthEndpoint(path)) {
    if (!retried && (await refreshAccessToken())) {
      return request<T>(method, path, body, true)
    }
    redirectToLogin()
    throw new ApiError(401, "Sessão expirada. Faça login novamente.")
  }

  if (!res.ok) {
    throw new ApiError(res.status, `Request failed (${res.status})`)
  }

  const data = await res.json()

  return data as T
}

export const api = {
  get: <T>(path: string) => request<T>("GET", path),
  post: <T>(path: string, body?: unknown) => request<T>("POST", path, body),
  put: <T>(path: string, body?: unknown) => request<T>("PUT", path, body),
  patch: <T>(path: string, body?: unknown) => request<T>("PATCH", path, body),
  delete: <T>(path: string) => request<T>("DELETE", path),
}

export { ApiError }
