<script setup lang="ts">
import { ref, computed } from "vue"
import { useRoute } from "vue-router"
import { useAuthStore } from "../stores/auth"

const auth = useAuthStore()
const route = useRoute()
const mobileOpen = ref(false)

const isActive = (path: string) => route.path.startsWith(path)

const navLinks = computed(() => {
  if (!auth.isAuthenticated) return []
  return [
    { path: "/dashboard", label: "Dashboard", icon: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" },
    { path: "/profile", label: "Perfil", icon: "M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" },
    { path: "/subscription", label: "Assinatura", icon: "M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" },
    { path: "/my-purchases", label: "Minhas Compras", icon: "M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" },
  ]
})

const userInitial = computed(() => (auth.user?.name || "U").charAt(0).toUpperCase())
</script>

<template>
  <nav class="sticky top-0 z-50 backdrop-blur-xl bg-white/90 border-b border-gray-200/60 shadow-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">
        <router-link to="/" class="flex items-center gap-2.5 group shrink-0">
          <div class="w-9 h-9 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-xl flex items-center justify-center shadow-md group-hover:shadow-lg transition-shadow">
            <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span class="text-lg font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600">Rifa Online</span>
        </router-link>

        <button class="sm:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors" @click="mobileOpen = !mobileOpen" :aria-label="mobileOpen ? 'Fechar menu' : 'Abrir menu'">
          <svg class="w-6 h-6 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path v-if="!mobileOpen" stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <div class="hidden sm:flex items-center gap-1">
          <template v-if="auth.isAuthenticated">
            <router-link
              v-for="link in navLinks"
              :key="link.path"
              :to="link.path"
              class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-lg transition-all"
              :class="isActive(link.path)
                ? 'text-indigo-700 bg-indigo-50 shadow-sm'
                : 'text-gray-600 hover:text-indigo-600 hover:bg-indigo-50/60'"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" :d="link.icon" />
              </svg>
              {{ link.label }}
            </router-link>

            <router-link
              v-if="auth.user?.role === 'ADMIN'"
              to="/admin"
              class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-lg transition-all"
              :class="isActive('/admin')
                ? 'text-amber-700 bg-amber-50 shadow-sm'
                : 'text-amber-600 hover:text-amber-700 hover:bg-amber-50/60'"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z" />
              </svg>
              Admin
            </router-link>

            <router-link
              to="/ajuda"
              class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-lg transition-all"
              :class="isActive('/ajuda')
                ? 'text-indigo-700 bg-indigo-50 shadow-sm'
                : 'text-gray-600 hover:text-indigo-600 hover:bg-indigo-50/60'"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Ajuda
            </router-link>

            <div class="w-px h-6 bg-gray-200 mx-2"></div>

            <div class="relative group">
              <button class="flex items-center gap-2 px-3 py-2 text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-indigo-50/60 rounded-lg transition-all">
                <div class="w-7 h-7 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center text-white text-xs font-bold shadow-sm">
                  {{ userInitial }}
                </div>
                <span class="hidden lg:inline max-w-24 truncate">{{ auth.user?.name }}</span>
                <svg class="w-4 h-4 text-gray-400 group-hover:text-indigo-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div class="absolute right-0 mt-1 w-48 bg-white rounded-xl shadow-lg border border-gray-200 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 translate-y-1 group-hover:translate-y-0 z-50">
                <div class="p-2 space-y-0.5">
                  <router-link to="/profile" class="flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors" @click="mobileOpen = false">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                    Meu Perfil
                  </router-link>
                  <router-link to="/subscription" class="flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                    </svg>
                    Assinatura
                  </router-link>
                  <div class="border-t border-gray-100 my-1"></div>
                  <button @click="auth.logout" class="flex items-center gap-2 w-full px-3 py-2 text-sm text-red-600 hover:bg-red-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    Sair
                  </button>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <router-link
              to="/ajuda"
              class="flex items-center gap-1.5 px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50/60 rounded-lg transition-all"
              :class="isActive('/ajuda') ? 'text-indigo-700 bg-indigo-50 shadow-sm' : ''"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Ajuda
            </router-link>
            <router-link to="/login" class="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-indigo-600 to-purple-600 rounded-lg hover:from-indigo-700 hover:to-purple-700 shadow-sm hover:shadow-md transition-all">
              Entrar
            </router-link>
            <router-link to="/register" class="px-4 py-2 text-sm font-medium text-indigo-600 bg-indigo-50 rounded-lg hover:bg-indigo-100 transition-all">
              Criar Conta
            </router-link>
          </template>
        </div>
      </div>

      <div
        v-if="mobileOpen"
        class="sm:hidden border-t border-gray-100 overflow-hidden animate-fade-in"
      >
        <template v-if="auth.isAuthenticated">
          <div class="flex items-center gap-3 px-3 py-4 border-b border-gray-100">
            <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center text-white text-base font-bold shadow-sm">
              {{ userInitial }}
            </div>
            <div class="min-w-0">
              <p class="text-sm font-semibold text-gray-900 truncate">{{ auth.user?.name }}</p>
              <p class="text-xs text-gray-500 truncate">{{ auth.user?.email }}</p>
            </div>
          </div>
          <div class="py-2 space-y-0.5">
            <router-link
              v-for="link in navLinks"
              :key="link.path"
              :to="link.path"
              class="flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-lg mx-2 transition-all"
              :class="isActive(link.path)
                ? 'text-indigo-700 bg-indigo-50'
                : 'text-gray-700 hover:text-indigo-600 hover:bg-indigo-50/60'"
              @click="mobileOpen = false"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" :d="link.icon" />
              </svg>
              {{ link.label }}
            </router-link>
            <router-link
              v-if="auth.user?.role === 'ADMIN'"
              to="/admin"
              class="flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-lg mx-2 transition-all"
              :class="isActive('/admin') ? 'text-amber-700 bg-amber-50' : 'text-amber-600 hover:bg-amber-50/60'"
              @click="mobileOpen = false"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z" />
              </svg>
              Admin
            </router-link>
            <router-link
              to="/ajuda"
              class="flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-lg mx-2 transition-all"
              :class="isActive('/ajuda') ? 'text-indigo-700 bg-indigo-50' : 'text-gray-700 hover:text-indigo-600 hover:bg-indigo-50/60'"
              @click="mobileOpen = false"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Ajuda
            </router-link>
          </div>
          <div class="border-t border-gray-100 py-2">
            <button @click="auth.logout" class="flex items-center gap-3 w-full px-4 py-2.5 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg mx-2 transition-colors" style="width: calc(100% - 16px)">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Sair
            </button>
          </div>
        </template>
        <template v-else>
          <div class="py-3 space-y-1">
            <router-link to="/ajuda" class="flex items-center gap-3 px-4 py-2.5 text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-indigo-50/60 rounded-lg mx-2 transition-colors" @click="mobileOpen = false">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Ajuda
            </router-link>
          </div>
          <div class="border-t border-gray-100 py-3 px-4 flex gap-2">
            <router-link to="/login" class="flex-1 text-center px-4 py-2.5 text-sm font-medium text-white bg-gradient-to-r from-indigo-600 to-purple-600 rounded-lg hover:from-indigo-700 hover:to-purple-700 shadow-sm transition-all" @click="mobileOpen = false">
              Entrar
            </router-link>
            <router-link to="/register" class="flex-1 text-center px-4 py-2.5 text-sm font-medium text-indigo-600 bg-indigo-50 rounded-lg hover:bg-indigo-100 transition-all" @click="mobileOpen = false">
              Criar Conta
            </router-link>
          </div>
        </template>
      </div>
    </div>
  </nav>
</template>
