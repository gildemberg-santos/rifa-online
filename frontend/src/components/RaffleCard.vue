<script setup lang="ts">
import { ref } from "vue"

const props = defineProps<{
  raffle: {
    id: string
    title: string
    description: string
    ticketPrice: number
    maxNumbers: number
    drawDate: string
    imageUrl?: string
    status: string
  }
}>()

const copied = ref(false)

function share(e: Event) {
  e.stopPropagation()
  e.preventDefault()
  const url = `${window.location.origin}/raffles/${props.raffle.id}`
  if (navigator.share) {
    navigator.share({ title: props.raffle.title, url })
  } else {
    navigator.clipboard.writeText(url)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}
</script>

<template>
  <router-link
    :to="`/raffles/${raffle.id}`"
    class="group block bg-white rounded-2xl shadow-sm border border-gray-200/80 overflow-hidden hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
  >
    <div v-if="raffle.imageUrl" class="h-36 sm:h-44 bg-gray-100 overflow-hidden">
      <img
        :src="raffle.imageUrl"
        :alt="raffle.title"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
      />
    </div>
    <div v-else class="h-36 sm:h-44 bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 flex items-center justify-center relative overflow-hidden">
      <div class="absolute inset-0 bg-white/10"></div>
      <span class="text-5xl text-white font-extrabold drop-shadow-lg">{{ raffle.title.charAt(0).toUpperCase() }}</span>
    </div>

    <div class="p-5">
      <div class="flex items-start justify-between gap-2">
        <h3 class="text-lg font-semibold text-gray-900 truncate group-hover:text-indigo-600 transition-colors">
          {{ raffle.title }}
        </h3>
        <div class="flex items-center gap-1 shrink-0">
          <button
            @click="share"
            class="p-1.5 rounded-lg hover:bg-gray-100 text-gray-400 hover:text-indigo-600 transition-colors"
            title="Compartilhar"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
            </svg>
          </button>
          <span
            class="px-2 py-0.5 text-xs font-semibold rounded-full"
            :class="raffle.status === 'ACTIVE' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'"
          >
            {{ raffle.status === "ACTIVE" ? "Ativa" : raffle.status }}
          </span>
        </div>
      </div>
      <p class="text-sm text-gray-500 mt-1.5 line-clamp-2 leading-relaxed">
        {{ raffle.description }}
      </p>

      <div class="flex justify-between items-center mt-4 pt-4 border-t border-gray-100">
        <div>
          <span class="text-xs text-gray-400">Preço</span>
          <p class="text-lg font-bold text-indigo-600">
            R$ {{ (raffle.ticketPrice / 100).toFixed(2) }}
          </p>
        </div>
        <div v-if="copied" class="text-xs font-medium text-green-600 animate-fade-in">
          Link copiado!
        </div>
        <div v-else class="text-right">
          <span class="text-xs text-gray-400">Números</span>
          <p class="text-sm font-medium text-gray-700">{{ raffle.maxNumbers }}</p>
        </div>
      </div>
    </div>
  </router-link>
</template>