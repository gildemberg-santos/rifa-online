import { defineStore } from "pinia"
import { computed } from "vue"
import { api } from "../utils/api"
import router from "../router"
import {
  accessToken,
  refreshToken,
  currentUser,
  setSession,
  clearSession,
  type SessionUser,
} from "../utils/session"

interface AuthResponse {
  user: SessionUser
  accessToken: string
  refreshToken: string
}

export const useAuthStore = defineStore("auth", () => {
  // Refs vindos do módulo de sessão: mantidos em sincronia com o api.ts
  // (inclusive em renovações silenciosas de token).
  const user = currentUser
  const isAuthenticated = computed(() => !!accessToken.value)

  async function register(name: string, email: string, password: string) {
    const res = await api.post<AuthResponse>("/auth/register", { name, email, password })
    setSession(res)
  }

  async function login(email: string, password: string) {
    const res = await api.post<AuthResponse>("/auth/login", { email, password })
    setSession(res)
  }

  async function refresh() {
    if (!refreshToken.value) throw new Error("No refresh token")
    const res = await api.post<AuthResponse>("/auth/refresh", { refreshToken: refreshToken.value })
    setSession(res)
  }

  function setTokens(access: string, refreshTok: string) {
    setSession({ accessToken: access, refreshToken: refreshTok })
  }

  function logout() {
    clearSession()
    router.push("/login")
  }

  return { user, accessToken, refreshToken, isAuthenticated, setTokens, register, login, refresh, logout }
})
