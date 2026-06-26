<script setup lang="ts">
import LegalLayout from "../../components/LegalLayout.vue"
</script>

<template>
  <LegalLayout title="Política de Cookies" updated-at="26 de junho de 2026">
    <p>
      Esta Política explica como a <strong>Rifa Online</strong> utiliza cookies e tecnologias semelhantes,
      complementando a <router-link to="/politica-de-privacidade">Política de Privacidade</router-link>.
    </p>

    <h2>1. O que são cookies</h2>
    <p>
      Cookies são pequenos arquivos armazenados no seu dispositivo quando você acessa um site. Também
      utilizamos tecnologias equivalentes, como o <em>localStorage</em>, para guardar informações no navegador.
    </p>

    <h2>2. Cookies e armazenamento que utilizamos</h2>
    <table>
      <thead>
        <tr><th>Tipo</th><th>Finalidade</th><th>Exemplo</th></tr>
      </thead>
      <tbody>
        <tr><td>Estritamente necessários</td><td>Autenticação e segurança da sessão</td><td>Token de acesso (localStorage)</td></tr>
        <tr><td>Funcionais</td><td>Lembrar preferências do usuário</td><td>Dados de perfil em cache</td></tr>
        <tr><td>Analíticos (se ativados)</td><td>Entender o uso da Plataforma de forma agregada</td><td>Ferramentas de analytics</td></tr>
      </tbody>
    </table>
    <p>
      Os cookies estritamente necessários são essenciais ao funcionamento e não podem ser desativados sem
      comprometer o uso da Plataforma.
    </p>

    <h2>3. Cookies de terceiros</h2>
    <p>
      O processamento de pagamentos pela <strong>InfinitePay</strong> pode envolver cookies próprios,
      sujeitos às políticas daquela empresa, sobre os quais não temos controle.
    </p>

    <h2>4. Como gerenciar cookies</h2>
    <p>
      Você pode bloquear ou apagar cookies nas configurações do seu navegador. A desativação de cookies
      necessários pode impedir o login e o uso de funcionalidades essenciais.
    </p>

    <h2>5. Alterações</h2>
    <p>
      Esta Política pode ser atualizada periodicamente; a versão vigente é a publicada nesta página.
    </p>

    <h2>6. Contato</h2>
    <p>
      Dúvidas podem ser enviadas pela nossa <router-link to="/contato">página de contato</router-link>.
    </p>
  </LegalLayout>
</template>
