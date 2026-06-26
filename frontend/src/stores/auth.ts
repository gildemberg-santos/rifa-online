import { defineStore } from "pinia"
import { computed } from "vue"
import { api } from "../utils/api"
import router from "../router"
import { accessToken, currentUser, setSession, clearSession, type SessionUser } from "../utils/session"

interface AuthResponse {
  user: SessionUser
  accessToken: string
}

export const useAuthStore = defineStore("auth", () => {
  // Refs vindos do módulo de sessão: mantidos em sincronia com o api.ts
  // (inclusive em renovações silenciosas de token). O refresh token fica
  // apenas no cookie HttpOnly do backend.
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
    const res = await api.post<AuthResponse>("/auth/refresh")
    setSession(res)
  }

  function setTokens(access: string) {
    setSession({ accessToken: access })
  }

  function logout() {
    clearSession()
    // Invalida o cookie de refresh no servidor (fire-and-forget).
    api.post("/auth/logout").catch(() => {})
    router.push("/login")
  }

  return { user, accessToken, isAuthenticated, setTokens, register, login, refresh, logout }
})
