<script setup lang="ts">
defineProps<{
  title: string
  updatedAt: string
}>()
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-12 animate-fade-in">
    <div class="mb-8">
      <router-link to="/" class="text-sm text-indigo-600 hover:text-indigo-700 font-medium inline-flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        Voltar ao início
      </router-link>
      <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 mt-4">{{ title }}</h1>
      <p class="text-sm text-gray-500 mt-2">Última atualização: {{ updatedAt }}</p>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-5 sm:p-10 legal-prose">
      <slot />
    </div>

    <div class="mt-8 flex flex-wrap gap-x-6 gap-y-2 text-sm text-gray-500">
      <router-link to="/termos-de-uso" class="hover:text-indigo-600">Termos de Uso</router-link>
      <router-link to="/politica-de-privacidade" class="hover:text-indigo-600">Política de Privacidade</router-link>
      <router-link to="/politica-de-cookies" class="hover:text-indigo-600">Política de Cookies</router-link>
      <router-link to="/termo-do-organizador" class="hover:text-indigo-600">Termo do Organizador</router-link>
      <router-link to="/contato" class="hover:text-indigo-600">Contato</router-link>
    </div>
  </div>
</template>

<style scoped>
.legal-prose :deep(h2) {
  font-size: 1.125rem;
  font-weight: 700;
  color: #111827;
  margin-top: 2rem;
  margin-bottom: 0.75rem;
}
.legal-prose :deep(h2:first-child) {
  margin-top: 0;
}
.legal-prose :deep(h3) {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin-top: 1.25rem;
  margin-bottom: 0.5rem;
}
.legal-prose :deep(p) {
  color: #374151;
  line-height: 1.7;
  margin-bottom: 0.875rem;
}
.legal-prose :deep(ul) {
  list-style: disc;
  padding-left: 1.5rem;
  margin-bottom: 0.875rem;
  color: #374151;
}
.legal-prose :deep(li) {
  line-height: 1.7;
  margin-bottom: 0.375rem;
}
.legal-prose :deep(strong) {
  color: #111827;
  font-weight: 600;
}
.legal-prose :deep(a) {
  color: #4f46e5;
  text-decoration: underline;
}
.legal-prose :deep(.table-wrap) {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  margin-bottom: 1rem;
}
.legal-prose :deep(table) {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
  min-width: 480px;
}
.legal-prose :deep(th),
.legal-prose :deep(td) {
  border: 1px solid #e5e7eb;
  padding: 0.5rem 0.75rem;
  text-align: left;
  vertical-align: top;
}
.legal-prose :deep(th) {
  background: #f9fafb;
  font-weight: 600;
  color: #111827;
}
.legal-prose :deep(.callout) {
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 0.75rem;
  padding: 1rem 1.25rem;
  margin: 1.25rem 0;
  color: #78350f;
}
</style>
