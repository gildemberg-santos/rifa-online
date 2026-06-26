<script setup lang="ts">
import { computed } from "vue"
import { useRoute } from "vue-router"
import { getArticle, helpCategories, articlesByCategory } from "../data/help"

const route = useRoute()
const article = computed(() => getArticle(route.params.slug as string))
const category = computed(() =>
  article.value ? helpCategories.find((c) => c.id === article.value!.category) : undefined,
)
const related = computed(() =>
  article.value
    ? articlesByCategory(article.value.category).filter((a) => a.slug !== article.value!.slug)
    : [],
)
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-12 animate-fade-in">
    <div class="mb-6 text-sm">
      <router-link to="/ajuda" class="text-indigo-600 hover:text-indigo-700 font-medium inline-flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        Central de Ajuda
      </router-link>
    </div>

    <template v-if="article">
      <p v-if="category" class="text-sm font-medium text-indigo-600">{{ category.title }}</p>
      <h1 class="text-3xl font-bold text-gray-900 mt-1 mb-2">{{ article.title }}</h1>
      <p class="text-gray-500 mb-8">{{ article.summary }}</p>

      <article class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6 sm:p-10 doc-prose" v-html="article.body"></article>

      <div v-if="related.length" class="mt-10">
        <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">Veja também</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <router-link
            v-for="a in related"
            :key="a.slug"
            :to="`/ajuda/${a.slug}`"
            class="block bg-white rounded-xl border border-gray-200 p-4 hover:border-indigo-300 hover:shadow-sm transition-all"
          >
            <h3 class="font-semibold text-gray-900 text-sm">{{ a.title }}</h3>
            <p class="text-xs text-gray-500 mt-0.5">{{ a.summary }}</p>
          </router-link>
        </div>
      </div>

      <div class="mt-10 text-center text-sm text-gray-500">
        Ainda com dúvida?
        <router-link to="/contato" class="text-indigo-600 hover:text-indigo-700 font-medium">Fale conosco</router-link>.
      </div>
    </template>

    <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-200 p-10 text-center">
      <h1 class="text-xl font-bold text-gray-900 mb-2">Tutorial não encontrado</h1>
      <p class="text-gray-500 mb-6">O conteúdo que você procura não existe ou foi movido.</p>
      <router-link to="/ajuda" class="inline-flex items-center px-5 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 transition-all">
        Voltar à Central de Ajuda
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.doc-prose :deep(h2) {
  font-size: 1.125rem;
  font-weight: 700;
  color: #111827;
  margin-top: 2rem;
  margin-bottom: 0.75rem;
}
.doc-prose :deep(h2:first-child) {
  margin-top: 0;
}
.doc-prose :deep(h3) {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin-top: 1.25rem;
  margin-bottom: 0.5rem;
}
.doc-prose :deep(p) {
  color: #374151;
  line-height: 1.7;
  margin-bottom: 0.875rem;
}
.doc-prose :deep(ul) {
  list-style: disc;
  padding-left: 1.5rem;
  margin-bottom: 0.875rem;
  color: #374151;
}
.doc-prose :deep(ol) {
  list-style: decimal;
  padding-left: 1.5rem;
  margin-bottom: 0.875rem;
  color: #374151;
}
.doc-prose :deep(li) {
  line-height: 1.7;
  margin-bottom: 0.375rem;
}
.doc-prose :deep(strong) {
  color: #111827;
  font-weight: 600;
}
.doc-prose :deep(a) {
  color: #4f46e5;
  text-decoration: underline;
}
.doc-prose :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.85em;
  background: #f3f4f6;
  padding: 0.1rem 0.35rem;
  border-radius: 0.35rem;
  color: #4338ca;
}
.doc-prose :deep(pre) {
  background: #0f172a;
  color: #e2e8f0;
  padding: 1rem 1.25rem;
  border-radius: 0.75rem;
  overflow-x: auto;
  margin-bottom: 1rem;
  font-size: 0.85rem;
  line-height: 1.6;
}
.doc-prose :deep(pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
}
.doc-prose :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1rem;
  font-size: 0.875rem;
}
.doc-prose :deep(th),
.doc-prose :deep(td) {
  border: 1px solid #e5e7eb;
  padding: 0.5rem 0.75rem;
  text-align: left;
  vertical-align: top;
}
.doc-prose :deep(th) {
  background: #f9fafb;
  font-weight: 600;
  color: #111827;
}
.doc-prose :deep(.tip) {
  background: #eef2ff;
  border: 1px solid #c7d2fe;
  border-radius: 0.75rem;
  padding: 0.875rem 1.125rem;
  margin: 1.25rem 0;
  color: #3730a3;
  font-size: 0.925rem;
}
</style>
