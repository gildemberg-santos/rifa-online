<script setup lang="ts">
import { ref, onMounted } from "vue"
import { api } from "../utils/api"
import RaffleCard from "../components/RaffleCard.vue"

interface Raffle {
  id: string
  title: string
  description: string
  ticketPrice: number
  maxNumbers: number
  drawDate: string
  imageUrl?: string
  status: string
}

const raffles = ref<Raffle[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    raffles.value = await api.get<Raffle[]>("/raffles")
  } catch {
    console.error("Failed to load raffles")
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <section class="relative overflow-hidden bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 text-white">
      <div class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmYiIGZpbGwtb3BhY2l0eT0iMC4wNSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMiIvPjwvZz48L2c+PC9zdmc+')] opacity-40"></div>
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 md:py-28 relative">
        <div class="text-center animate-fade-in">
          <div class="inline-flex items-center gap-2 px-3 py-1 bg-white/15 backdrop-blur rounded-full text-sm text-indigo-100 mb-6">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
            </svg>
            Teste grátis por 7 dias &mdash; sem compromisso
          </div>
          <h1 class="text-4xl md:text-6xl font-extrabold tracking-tight mb-6 leading-tight">
            A plataforma mais simples<br>
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-amber-300 to-yellow-200">para criar e gerenciar rifas online</span>
          </h1>
          <p class="text-lg md:text-xl text-indigo-100 max-w-2xl mx-auto mb-10 leading-relaxed">
            Crie rifas em segundos, venda números com pagamento via PIX, realize sorteios automáticos e receba tudo direto na sua conta InfinitePay.
          </p>
          <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
            <router-link
              to="/register"
              class="inline-flex items-center px-8 py-3.5 bg-white text-indigo-700 font-bold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all text-base"
            >
              Começar Grátis
              <svg class="w-5 h-5 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </router-link>
            <a
              href="#features"
              class="inline-flex items-center px-8 py-3.5 bg-white/10 backdrop-blur text-white font-semibold rounded-xl hover:bg-white/20 transition-all text-base border border-white/20"
            >
              Saiba Mais
            </a>
          </div>
        </div>
      </div>
      <div class="absolute bottom-0 left-0 right-0 h-16 bg-gradient-to-t from-gray-50 to-transparent"></div>
    </section>

    <section id="features" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="text-center mb-14">
        <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">Tudo que você precisa para organizar rifas</h2>
        <p class="text-gray-500 max-w-xl mx-auto">Ferramentas completas para criar, divulgar e sortear sua rifa com segurança.</p>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div class="bg-white rounded-2xl p-6 border border-gray-200 hover:shadow-lg hover:-translate-y-0.5 transition-all">
          <div class="w-12 h-12 bg-indigo-100 rounded-xl flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Criação Rápida</h3>
          <p class="text-gray-500 text-sm leading-relaxed">Cadastre sua rifa em segundos: título, descrição, valor do número, quantidade e data do sorteio. Sem burocracia.</p>
        </div>
        <div class="bg-white rounded-2xl p-6 border border-gray-200 hover:shadow-lg hover:-translate-y-0.5 transition-all">
          <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Pagamento via PIX</h3>
          <p class="text-gray-500 text-sm leading-relaxed">Integração com InfinitePay. Seus participantes pagam via PIX e o dinheiro cai direto na sua conta.</p>
        </div>
        <div class="bg-white rounded-2xl p-6 border border-gray-200 hover:shadow-lg hover:-translate-y-0.5 transition-all">
          <div class="w-12 h-12 bg-amber-100 rounded-xl flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Sorteio Automático</h3>
          <p class="text-gray-500 text-sm leading-relaxed">Na data marcada, o sistema sorteia automaticamente um vencedor entre os números pagos. Resultado divulgado na hora.</p>
        </div>
        <div class="bg-white rounded-2xl p-6 border border-gray-200 hover:shadow-lg hover:-translate-y-0.5 transition-all">
          <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Dashboard Completo</h3>
          <p class="text-gray-500 text-sm leading-relaxed">Acompanhe vendas em tempo real, veja receita, percentual de ingressos vendidos e estatísticas detalhadas.</p>
        </div>
        <div class="bg-white rounded-2xl p-6 border border-gray-200 hover:shadow-lg hover:-translate-y-0.5 transition-all">
          <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Segurança e Privacidade</h3>
          <p class="text-gray-500 text-sm leading-relaxed">Dados criptografados, pagamento processado pela InfinitePay e proteção de informações dos participantes.</p>
        </div>
        <div class="bg-white rounded-2xl p-6 border border-gray-200 hover:shadow-lg hover:-translate-y-0.5 transition-all">
          <div class="w-12 h-12 bg-rose-100 rounded-xl flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-rose-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Reserva com TTL</h3>
          <p class="text-gray-500 text-sm leading-relaxed">Números ficam reservados por 10 minutos durante o checkout. Se o pagamento não for concluído, voltam automaticamente.</p>
        </div>
      </div>
    </section>

    <section class="bg-gradient-to-br from-indigo-50 via-white to-purple-50 py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">Como funciona</h2>
          <p class="text-gray-500 max-w-xl mx-auto">Em poucos passos você cria sua rifa e começa a vender.</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div class="text-center">
            <div class="w-14 h-14 mx-auto bg-indigo-600 rounded-2xl flex items-center justify-center text-white text-xl font-bold shadow-lg mb-4">1</div>
            <h3 class="font-bold text-gray-900 mb-2">Crie sua conta</h3>
            <p class="text-gray-500 text-sm">Cadastre-se gratuitamente e ative seu trial de 7 dias.</p>
          </div>
          <div class="text-center">
            <div class="w-14 h-14 mx-auto bg-indigo-600 rounded-2xl flex items-center justify-center text-white text-xl font-bold shadow-lg mb-4">2</div>
            <h3 class="font-bold text-gray-900 mb-2">Configure a rifa</h3>
            <p class="text-gray-500 text-sm">Defina título, valor do número, quantidade e data do sorteio.</p>
          </div>
          <div class="text-center">
            <div class="w-14 h-14 mx-auto bg-indigo-600 rounded-2xl flex items-center justify-center text-white text-xl font-bold shadow-lg mb-4">3</div>
            <h3 class="font-bold text-gray-900 mb-2">Divulgue e venda</h3>
            <p class="text-gray-500 text-sm">Compartilhe o link e receba pagamentos via PIX automaticamente.</p>
          </div>
          <div class="text-center">
            <div class="w-14 h-14 mx-auto bg-indigo-600 rounded-2xl flex items-center justify-center text-white text-xl font-bold shadow-lg mb-4">4</div>
            <h3 class="font-bold text-gray-900 mb-2">Sorteie o vencedor</h3>
            <p class="text-gray-500 text-sm">Na data marcada, o sistema sorteia e divulga o resultado.</p>
          </div>
        </div>
      </div>
    </section>

    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 rounded-3xl p-8 md:p-12 text-white text-center relative overflow-hidden">
        <div class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmYiIGZpbGwtb3BhY2l0eT0iMC4wNSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMiIvPjwvZz48L2c+PC9zdmc+')] opacity-40"></div>
        <div class="relative">
          <h2 class="text-3xl md:text-4xl font-bold mb-4">Preço justo, sem surpresas</h2>
          <p class="text-indigo-100 text-lg mb-10 max-w-xl mx-auto">Teste grátis por 7 dias, depois apenas R$10 por mês.</p>
          <div class="max-w-sm mx-auto bg-white/10 backdrop-blur rounded-2xl p-8 border border-white/20">
            <p class="text-sm text-indigo-200 uppercase tracking-wider font-medium mb-2">Assinatura Mensal</p>
            <p class="text-5xl font-extrabold mb-2">R$10</p>
            <p class="text-indigo-200 text-sm mb-6">por mês, cancele quando quiser</p>
            <ul class="text-left space-y-3 mb-8">
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-green-300 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-sm">Rifas ilimitadas</span>
              </li>
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-green-300 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-sm">Até 1.000 números por rifa</span>
              </li>
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-green-300 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-sm">Pagamento via PIX (InfinitePay)</span>
              </li>
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-green-300 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-sm">Sorteio automático</span>
              </li>
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-green-300 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-sm">Dashboard com estatísticas</span>
              </li>
            </ul>
            <router-link
              to="/register"
              class="block w-full py-3 bg-white text-indigo-700 font-bold rounded-xl hover:bg-indigo-50 transition-all text-center shadow-lg"
            >
              Começar Teste Grátis
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <section id="raffles-section" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Rifas Ativas</h2>
          <p class="text-gray-500 text-sm mt-1">Participe escolhendo seus números da sorte</p>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-16">
        <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
      </div>

      <div v-else-if="raffles.length === 0" class="text-center py-16 bg-white rounded-2xl border border-gray-200">
        <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-gray-500 text-lg">Nenhuma rifa ativa no momento.</p>
        <p class="text-gray-400 text-sm mt-1">Volte mais tarde para conferir as novidades.</p>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <RaffleCard
          v-for="(raffle, idx) in raffles"
          :key="raffle.id"
          :raffle="raffle"
          :style="{ animationDelay: `${idx * 0.08}s` }"
          class="animate-slide-up opacity-0"
        />
      </div>
    </section>

    <section class="bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 py-16">
      <div class="max-w-4xl mx-auto px-4 text-center text-white">
        <h2 class="text-3xl md:text-4xl font-bold mb-4">Pronto para começar?</h2>
        <p class="text-indigo-100 text-lg mb-8 max-w-lg mx-auto">Crie sua primeira rifa agora mesmo. Teste grátis por 7 dias, sem compromisso.</p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
          <router-link
            to="/register"
            class="inline-flex items-center px-8 py-3.5 bg-white text-indigo-700 font-bold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all text-base"
          >
            Criar Conta Grátis
            <svg class="w-5 h-5 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </router-link>
          <router-link
            to="/login"
            class="inline-flex items-center px-8 py-3.5 bg-white/10 backdrop-blur text-white font-semibold rounded-xl hover:bg-white/20 transition-all text-base border border-white/20"
          >
            Já tenho conta
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.animate-slide-up {
  animation: slideUp 0.4s ease-out forwards;
}
@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
