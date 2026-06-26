import { ref } from "vue"

// Estado de sessão reativo e único, compartilhado entre o cliente HTTP (api.ts)
// e a store de autenticação. Mantém os refs e o localStorage sempre em sincronia,
// para que uma renovação silenciosa de token reflita na UI (navbar, guards, etc.).

export interface SessionUser {
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

function loadUser(): SessionUser | null {
  const raw = localStorage.getItem("user")
  if (!raw) return null
  try {
    return JSON.parse(raw) as SessionUser
  } catch {
    return null
  }
}

export const accessToken = ref<string | null>(localStorage.getItem("accessToken"))
export const refreshToken = ref<string | null>(localStorage.getItem("refreshToken"))
export const currentUser = ref<SessionUser | null>(loadUser())

export function setSession(data: {
  accessToken: string
  refreshToken?: string
  user?: SessionUser
}) {
  accessToken.value = data.accessToken
  localStorage.setItem("accessToken", data.accessToken)
  if (data.refreshToken) {
    refreshToken.value = data.refreshToken
    localStorage.setItem("refreshToken", data.refreshToken)
  }
  if (data.user) {
    currentUser.value = data.user
    localStorage.setItem("user", JSON.stringify(data.user))
  }
}

export function clearSession() {
  accessToken.value = null
  refreshToken.value = null
  currentUser.value = null
  localStorage.removeItem("accessToken")
  localStorage.removeItem("refreshToken")
  localStorage.removeItem("user")
}
