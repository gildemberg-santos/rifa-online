<template>
  <footer class="bg-gray-900 text-gray-400 mt-auto">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-8">
        <div>
          <div class="flex items-center gap-2 text-lg font-bold text-white mb-3">
            <svg class="w-6 h-6 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Rifa Online
          </div>
          <p class="text-sm leading-relaxed">Plataforma para criação e gerenciamento de rifas online.</p>
        </div>
        <div>
          <h4 class="text-white font-semibold mb-3 text-sm uppercase tracking-wider">Links</h4>
          <ul class="space-y-2 text-sm">
            <li><router-link to="/" class="hover:text-white transition-colors">Início</router-link></li>
            <li><router-link to="/ajuda" class="hover:text-white transition-colors">Central de Ajuda</router-link></li>
            <li><router-link to="/login" class="hover:text-white transition-colors">Entrar</router-link></li>
          </ul>
        </div>
        <div>
          <h4 class="text-white font-semibold mb-3 text-sm uppercase tracking-wider">Legal</h4>
          <ul class="space-y-2 text-sm">
            <li><router-link to="/termos-de-uso" class="hover:text-white transition-colors">Termos de Uso</router-link></li>
            <li><router-link to="/politica-de-privacidade" class="hover:text-white transition-colors">Política de Privacidade</router-link></li>
            <li><router-link to="/politica-de-cookies" class="hover:text-white transition-colors">Política de Cookies</router-link></li>
            <li><router-link to="/termo-do-organizador" class="hover:text-white transition-colors">Termo do Organizador</router-link></li>
          </ul>
        </div>
        <div>
          <h4 class="text-white font-semibold mb-3 text-sm uppercase tracking-wider">Atendimento</h4>
          <ul class="space-y-2 text-sm">
            <li><router-link to="/contato" class="hover:text-white transition-colors">Fale conosco</router-link></li>
          </ul>
        </div>
      </div>
      <div class="border-t border-gray-800 mt-8 pt-8 text-center text-sm">
        &copy; {{ new Date().getFullYear() }} Rifa Online. Todos os direitos reservados.
      </div>
    </div>
  </footer>
</template>
