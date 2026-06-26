import { ref } from "vue"

// Estado de sessão reativo e único, compartilhado entre o cliente HTTP (api.ts)
// e a store de autenticação. Mantém os refs e o localStorage sempre em sincronia,
// para que uma renovação silenciosa de token reflita na UI (navbar, guards, etc.).
//
// O refresh token NÃO é armazenado no cliente: ele vive apenas num cookie HttpOnly
// gerenciado pelo backend (protegido contra XSS). Aqui guardamos só o access token
// (curta duração) e os dados do usuário.

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
export const currentUser = ref<SessionUser | null>(loadUser())

export function setSession(data: { accessToken: string; user?: SessionUser }) {
  accessToken.value = data.accessToken
  localStorage.setItem("accessToken", data.accessToken)
  if (data.user) {
    currentUser.value = data.user
    localStorage.setItem("user", JSON.stringify(data.user))
  }
}

export function clearSession() {
  accessToken.value = null
  currentUser.value = null
  localStorage.removeItem("accessToken")
  localStorage.removeItem("refreshToken") // limpeza de versões antigas
  localStorage.removeItem("user")
}
