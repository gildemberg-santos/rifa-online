<template>
  <div class="max-w-lg mx-auto px-4 py-12 text-center">
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <div class="text-6xl mb-4 animate-pulse">&#8987;</div>
      <h1 class="text-2xl font-bold text-yellow-600 mb-2">Aguardando Pagamento</h1>
      <p class="text-gray-600 mb-6">
        Sua compra está sendo processada. Assim que o pagamento for confirmado, você receberá uma confirmação.
      </p>
      <router-link
        to="/"
        class="inline-block bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700"
      >
        Voltar para Home
      </router-link>
    </div>
  </div>
</template>
