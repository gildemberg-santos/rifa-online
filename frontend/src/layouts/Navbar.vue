<script setup lang="ts">
import { useAuthStore } from "../stores/auth"

const auth = useAuthStore()
</script>

<template>
  <nav class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">
        <router-link to="/" class="text-xl font-bold text-indigo-600">
          Rifa Online
        </router-link>

        <div class="flex items-center gap-4">
          <router-link to="/" class="text-sm text-gray-600 hover:text-gray-900">
            Rifas
          </router-link>

          <template v-if="auth.isAuthenticated">
            <router-link to="/dashboard" class="text-sm text-gray-600 hover:text-gray-900">
              Dashboard
            </router-link>
            <button
              @click="auth.logout"
              class="text-sm text-red-600 hover:text-red-800"
            >
              Sair
            </button>
          </template>

          <template v-else>
            <router-link
              to="/login"
              class="text-sm text-gray-600 hover:text-gray-900"
            >
              Entrar
            </router-link>

          </template>
        </div>
      </div>
    </div>
  </nav>
</template>
