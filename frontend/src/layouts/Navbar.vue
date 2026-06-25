<script setup lang="ts">
import { ref } from "vue"
import { useAuthStore } from "../stores/auth"

const auth = useAuthStore()
const mobileOpen = ref(false)
</script>

<template>
  <nav class="sticky top-0 z-50 backdrop-blur-lg bg-white/80 border-b border-gray-200/60">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">
        <router-link to="/" class="flex items-center gap-2 text-xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600">
          <svg class="w-7 h-7 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Rifa Online
        </router-link>

        <button class="sm:hidden p-2 rounded-lg hover:bg-gray-100" @click="mobileOpen = !mobileOpen">
          <svg class="w-6 h-6 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path v-if="!mobileOpen" stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <div class="hidden sm:flex items-center gap-1">
          <template v-if="auth.isAuthenticated">
            <router-link to="/dashboard" class="px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors">
              Dashboard
            </router-link>
            <router-link to="/profile" class="px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors">
              Perfil
            </router-link>
            <router-link to="/subscription" class="px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors">
              Assinatura
            </router-link>
            <router-link to="/my-purchases" class="px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors">
              Minhas Compras
            </router-link>
            <router-link v-if="auth.user?.role === 'ADMIN'" to="/admin" class="px-3 py-2 text-sm font-medium text-amber-600 hover:bg-amber-50 rounded-lg transition-colors">
              Admin
            </router-link>
            <div class="w-px h-6 bg-gray-200 mx-2"></div>
            <button @click="auth.logout" class="px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors">
              Sair
            </button>
          </template>
          <template v-else>
            <router-link to="/login" class="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-indigo-600 to-purple-600 rounded-lg hover:from-indigo-700 hover:to-purple-700 shadow-sm transition-all">
              Entrar
            </router-link>
          </template>
        </div>
      </div>

      <div v-if="mobileOpen" class="sm:hidden pb-4 border-t border-gray-100 pt-3 space-y-1 animate-fade-in">
        <template v-if="auth.isAuthenticated">
          <router-link to="/dashboard" class="block px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg" @click="mobileOpen = false">
            Dashboard
          </router-link>
          <router-link to="/profile" class="block px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg" @click="mobileOpen = false">
            Perfil
          </router-link>
          <router-link to="/subscription" class="block px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg" @click="mobileOpen = false">
            Assinatura
          </router-link>
          <router-link to="/my-purchases" class="block px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg" @click="mobileOpen = false">
            Minhas Compras
          </router-link>
          <router-link v-if="auth.user?.role === 'ADMIN'" to="/admin" class="block px-3 py-2 text-sm font-medium text-amber-600 hover:bg-amber-50 rounded-lg" @click="mobileOpen = false">
            Admin
          </router-link>
          <button @click="auth.logout" class="block w-full text-left px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg">
            Sair
          </button>
        </template>
        <template v-else>
          <router-link to="/login" class="block px-3 py-2 text-sm font-medium text-indigo-600 hover:bg-indigo-50 rounded-lg" @click="mobileOpen = false">
            Entrar
          </router-link>
        </template>
      </div>
    </div>
  </nav>
</template>