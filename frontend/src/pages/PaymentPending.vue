<template>
  <div class="max-w-lg mx-auto px-4 py-12 text-center animate-scale-in">
    <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-10">
      <div class="w-20 h-20 mx-auto bg-amber-100 rounded-full flex items-center justify-center mb-6">
        <svg class="w-10 h-10 text-amber-600 animate-pulse" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-amber-600 mb-2">Aguardando Pagamento</h1>
      <p class="text-gray-500 mb-2">Sua compra está sendo processada.</p>
      <p class="text-sm text-gray-400 mb-8">Assim que o pagamento for confirmado, você receberá uma confirmação.</p>
      <router-link
        to="/"
        class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 shadow-md transition-all"
      >
        Voltar para Home
      </router-link>
    </div>
  </div>
</template>