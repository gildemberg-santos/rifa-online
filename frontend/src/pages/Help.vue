<script setup lang="ts">
import { ref, computed } from "vue"
import { helpCategories, helpArticles, articlesByCategory } from "../data/help"

const query = ref("")

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return null
  return helpArticles.filter(
    (a) => a.title.toLowerCase().includes(q) || a.summary.toLowerCase().includes(q),
  )
})
</script>

<template>
  <div class="max-w-5xl mx-auto px-4 py-12 animate-fade-in">
    <div class="text-center mb-10">
      <h1 class="text-3xl font-bold text-gray-900">Central de Ajuda</h1>
      <p class="text-gray-500 mt-2">Tutoriais de todos os recursos da Rifa Online.</p>
      <div class="mt-6 max-w-xl mx-auto relative">
        <svg class="w-5 h-5 text-gray-400 absolute left-4 top-1/2 -translate-y-1/2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="query"
          type="search"
          placeholder="Buscar um tutorial..."
          class="w-full pl-11 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-shadow"
        />
      </div>
    </div>

    <!-- Resultados da busca -->
    <div v-if="filtered" class="space-y-3">
      <p class="text-sm text-gray-500">{{ filtered.length }} resultado(s) para "{{ query }}"</p>
      <router-link
        v-for="a in filtered"
        :key="a.slug"
        :to="`/ajuda/${a.slug}`"
        class="block bg-white rounded-xl border border-gray-200 p-4 hover:border-indigo-300 hover:shadow-sm transition-all"
      >
        <h3 class="font-semibold text-gray-900">{{ a.title }}</h3>
        <p class="text-sm text-gray-500 mt-0.5">{{ a.summary }}</p>
      </router-link>
      <p v-if="filtered.length === 0" class="text-center py-10 text-gray-500">
        Nada encontrado. Tente outras palavras ou
        <router-link to="/contato" class="text-indigo-600 hover:text-indigo-700">fale conosco</router-link>.
      </p>
    </div>

    <!-- Categorias -->
    <div v-else class="space-y-10">
      <section v-for="cat in helpCategories" :key="cat.id">
        <div class="flex items-start gap-3 mb-4">
          <div class="shrink-0 w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" :d="cat.icon" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-bold text-gray-900">{{ cat.title }}</h2>
            <p class="text-sm text-gray-500">{{ cat.description }}</p>
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <router-link
            v-for="a in articlesByCategory(cat.id)"
            :key="a.slug"
            :to="`/ajuda/${a.slug}`"
            class="block bg-white rounded-xl border border-gray-200 p-4 hover:border-indigo-300 hover:shadow-sm transition-all"
          >
            <h3 class="font-semibold text-gray-900">{{ a.title }}</h3>
            <p class="text-sm text-gray-500 mt-0.5">{{ a.summary }}</p>
          </router-link>
        </div>
      </section>
    </div>

    <div class="mt-12 text-center text-sm text-gray-500">
      Não encontrou o que procurava?
      <router-link to="/contato" class="text-indigo-600 hover:text-indigo-700 font-medium">Fale conosco</router-link>.
    </div>
  </div>
</template>
