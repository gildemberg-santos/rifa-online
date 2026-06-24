<template>
  <div class="max-w-lg mx-auto px-4 py-12 text-center animate-scale-in">
    <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-10">
      <div class="w-20 h-20 mx-auto bg-green-100 rounded-full flex items-center justify-center mb-6">
        <svg class="w-10 h-10 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-green-600 mb-2">Pagamento Confirmado!</h1>
      <p class="text-gray-500 mb-8">Seus números foram registrados com sucesso.</p>
      <router-link
        to="/"
        class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md transition-all"
      >
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
        </svg>
        Voltar para Home
      </router-link>
    </div>
  </div>
</template>