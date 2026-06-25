import { defineStore } from "pinia"
import { ref, computed } from "vue"
import { api } from "../utils/api"
import router from "../router"

interface User {
  id: string
  name: string
  email: string
  phone?: string
  role: string
  infinitePayHandle?: string
  subscriptionStatus: string
  subscriptionIsTrial: boolean
  hasSubscriptionBefore?: boolean
}

interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
}

function loadUser(): User | null {
  const raw = localStorage.getItem("user")
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

export const useAuthStore = defineStore("auth", () => {
  const user = ref<User | null>(loadUser())
  const accessToken = ref<string | null>(localStorage.getItem("accessToken"))
  const refreshToken = ref<string | null>(localStorage.getItem("refreshToken"))

  const isAuthenticated = computed(() => !!accessToken.value)

  async function register(name: string, email: string, password: string) {
    const res = await api.post<AuthResponse>("/auth/register", { name, email, password })
    setAuth(res)
  }

  async function login(email: string, password: string) {
    const res = await api.post<AuthResponse>("/auth/login", { email, password })
    setAuth(res)
  }

  async function refresh() {
    if (!refreshToken.value) throw new Error("No refresh token")
    const res = await api.post<AuthResponse>("/auth/refresh", { refreshToken: refreshToken.value })
    setAuth(res)
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem("accessToken", access)
    localStorage.setItem("refreshToken", refresh)
  }

  function setAuth(data: AuthResponse) {
    setTokens(data.accessToken, data.refreshToken)
    user.value = data.user
    localStorage.setItem("user", JSON.stringify(data.user))
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem("accessToken")
    localStorage.removeItem("refreshToken")
    localStorage.removeItem("user")
    router.push("/login")
  }

  return { user, accessToken, refreshToken, isAuthenticated, setTokens, register, login, refresh, logout }
})
