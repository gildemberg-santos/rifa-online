import { defineStore } from "pinia"
import { ref, computed } from "vue"
import { api } from "../utils/api"
import router from "../router"

interface User {
  id: string
  name: string
  email: string
}

interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
}

export const useAuthStore = defineStore("auth", () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem("accessToken"))
  const refreshToken = ref<string | null>(localStorage.getItem("refreshToken"))

  const isAuthenticated = computed(() => !!accessToken.value)

  async function register(name: string, email: string, password: string) {
    const res = await api.post<AuthResponse>("/auth/register", { name, email, password })
    setTokens(res.accessToken, res.refreshToken)
    user.value = res.user
  }

  async function login(email: string, password: string) {
    const res = await api.post<AuthResponse>("/auth/login", { email, password })
    setTokens(res.accessToken, res.refreshToken)
    user.value = res.user
  }

  async function refresh() {
    if (!refreshToken.value) throw new Error("No refresh token")
    const res = await api.post<AuthResponse>("/auth/refresh", { refreshToken: refreshToken.value })
    setTokens(res.accessToken, res.refreshToken)
    user.value = res.user
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem("accessToken", access)
    localStorage.setItem("refreshToken", refresh)
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem("accessToken")
    localStorage.removeItem("refreshToken")
    router.push("/login")
  }

  return { user, accessToken, refreshToken, isAuthenticated, setTokens, register, login, refresh, logout }
})
