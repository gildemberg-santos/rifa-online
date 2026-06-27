import { defineStore } from "pinia"
import { computed } from "vue"
import { api } from "../utils/api"
import { sendEvent } from "../utils/analytics"
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
    await api.post("/auth/register", { name, email, password })
    sendEvent("user_registration", { method: "email", user_email: email })
  }

  async function verifyEmail(email: string, code: string) {
    const res = await api.post<AuthResponse>("/auth/verify-email", { email, code })
    setSession(res)
  }

  async function resendCode(email: string) {
    await api.post("/auth/resend-code", { email })
  }

  async function login(email: string, password: string) {
    const res = await api.post<AuthResponse>("/auth/login", { email, password })
    setSession(res)
    sendEvent("user_login", { method: "email", user_email: email })
  }

  async function refresh() {
    const res = await api.post<AuthResponse>("/auth/refresh")
    setSession(res)
  }

  function setTokens(access: string) {
    setSession({ accessToken: access })
  }

  function logout() {
    const email = user.value?.email
    clearSession()
    sendEvent("user_logout", { user_email: email })
    // Invalida o cookie de refresh no servidor (fire-and-forget).
    api.post("/auth/logout").catch(() => {})
    router.push("/login")
  }

  return { user, accessToken, isAuthenticated, setTokens, register, verifyEmail, resendCode, login, refresh, logout }
})
