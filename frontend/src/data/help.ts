// Central de Ajuda — conteúdo dos tutoriais.
// Conteúdo estático e confiável (renderizado via v-html nas páginas de ajuda).

export interface HelpCategory {
  id: string
  title: string
  description: string
  /** SVG path "d" para o ícone (heroicons outline). */
  icon: string
}

export interface HelpArticle {
  slug: string
  category: string
  title: string
  summary: string
  /** Corpo em HTML simples (h2/h3/p/ul/ol/strong). */
  body: string
}

export const helpCategories: HelpCategory[] = [
  {
    id: "primeiros-passos",
    title: "Primeiros passos",
    description: "Entenda o que é a plataforma e como ela funciona.",
    icon: "M13 10V3L4 14h7v7l9-11h-7z",
  },
  {
    id: "organizador",
    title: "Para organizadores",
    description: "Crie e gerencie suas rifas, acompanhe vendas e faça sorteios.",
    icon: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
  },
  {
    id: "participante",
    title: "Para participantes",
    description: "Compre números, acompanhe suas compras e veja resultados.",
    icon: "M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z",
  },
  {
    id: "pagamentos",
    title: "Pagamentos e assinatura",
    description: "Como funcionam os pagamentos via InfinitePay e a assinatura.",
    icon: "M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z",
  },
  {
    id: "administracao",
    title: "Administração",
    description: "Recursos do painel administrativo da plataforma.",
    icon: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z",
  },
  {
    id: "conta-suporte",
    title: "Conta, privacidade e suporte",
    description: "Gerencie sua conta, conheça seus direitos e fale conosco.",
    icon: "M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z",
  },
]

export const helpArticles: HelpArticle[] = [
  // ---------------- PRIMEIROS PASSOS ----------------
  {
    slug: "o-que-e-a-rifa-online",
    category: "primeiros-passos",
    title: "O que é a Rifa Online",
    summary: "Visão geral da plataforma e para quem ela serve.",
    body: `
      <h2>O que é</h2>
      <p>A <strong>Rifa Online</strong> é uma plataforma que permite criar, divulgar e gerenciar rifas pela internet, com pagamento via PIX ou cartão. Ela é uma <strong>ferramenta tecnológica</strong>: quem cria e conduz a rifa é o organizador.</p>
      <h2>Para quem é</h2>
      <ul>
        <li><strong>Organizadores</strong> — assinam o serviço para criar rifas, vender números e realizar sorteios.</li>
        <li><strong>Participantes</strong> — compram números de uma rifa, sem precisar criar conta.</li>
        <li><strong>Administração</strong> — opera a plataforma (usuários, assinaturas, rifas e mensagens).</li>
      </ul>
      <h2>Por onde começar</h2>
      <ul>
        <li>Vai criar rifas? Veja <a href="/ajuda/criar-conta-e-trial">Criar conta e período de teste</a>.</li>
        <li>Vai participar de uma rifa? Veja <a href="/ajuda/comprar-numeros">Como comprar números</a>.</li>
      </ul>
    `,
  },
  {
    slug: "como-a-plataforma-funciona",
    category: "primeiros-passos",
    title: "Como a plataforma funciona",
    summary: "O fluxo completo, do cadastro ao sorteio.",
    body: `
      <h2>Fluxo do organizador</h2>
      <ol>
        <li>Cria a conta e entra no período de teste (trial) de 7 dias.</li>
        <li>Assina o plano para continuar usando após o trial.</li>
        <li>Configura o recebimento (handle InfinitePay) no perfil.</li>
        <li>Cria a rifa (título, valor, quantidade de números, data do sorteio).</li>
        <li>Divulga o link e acompanha as vendas pelo Dashboard.</li>
        <li>Na data marcada, realiza o sorteio e divulga o resultado.</li>
      </ol>
      <h2>Fluxo do participante</h2>
      <ol>
        <li>Abre a rifa e escolhe um ou mais números disponíveis.</li>
        <li>Informa nome e telefone e vai para o pagamento (PIX ou cartão).</li>
        <li>Após pagar, os números ficam registrados em seu nome.</li>
        <li>Acompanha em "Minhas compras" pelo telefone e vê o resultado no dia do sorteio.</li>
      </ol>
      <div class="tip">Dica: o telefone informado na compra é o que identifica o participante. Use sempre o mesmo número para encontrar todas as suas compras.</div>
    `,
  },

  // ---------------- ORGANIZADOR ----------------
  {
    slug: "criar-conta-e-trial",
    category: "organizador",
    title: "Criar conta e período de teste (trial)",
    summary: "Cadastre-se e use gratuitamente por 7 dias.",
    body: `
      <h2>Passo a passo</h2>
      <ol>
        <li>Acesse <a href="/register">Criar conta</a>.</li>
        <li>Informe nome, e-mail e senha (mínimo de 6 caracteres).</li>
        <li>Marque o aceite dos Termos de Uso e da Política de Privacidade.</li>
        <li>Clique em <strong>Criar conta</strong>. Você já entra com 7 dias de teste.</li>
      </ol>
      <h2>O que o trial libera</h2>
      <p>Durante os 7 dias você pode criar e gerenciar rifas normalmente. Quando o teste expira, a criação e a gestão de rifas ficam bloqueadas até a assinatura ser ativada.</p>
      <div class="tip">Você pode assinar a qualquer momento durante o trial — veja <a href="/ajuda/assinatura-mensal">Assinatura mensal</a>.</div>
    `,
  },
  {
    slug: "configurar-recebimento",
    category: "organizador",
    title: "Configurar recebimento (handle InfinitePay)",
    summary: "Defina para onde vão os pagamentos das suas rifas.",
    body: `
      <h2>Por que é necessário</h2>
      <p>Os pagamentos dos participantes vão <strong>direto para a sua conta InfinitePay</strong>. Sem o handle configurado, as compras não têm para onde ir.</p>
      <h2>Passo a passo</h2>
      <ol>
        <li>Acesse <a href="/profile">Perfil</a>.</li>
        <li>No campo <strong>handle InfinitePay</strong>, informe sua tag (ex.: <em>$suatag</em>).</li>
        <li>Salve as alterações.</li>
      </ol>
      <h2>Onde conseguir o handle</h2>
      <p>O handle é o identificador da sua conta na InfinitePay (começa com <strong>$</strong>). Ele é criado no aplicativo da InfinitePay. A Rifa Online não cria nem administra essa conta.</p>
    `,
  },
  {
    slug: "criar-rifa",
    category: "organizador",
    title: "Como criar uma rifa",
    summary: "Configure título, valor, quantidade de números e data do sorteio.",
    body: `
      <h2>Passo a passo</h2>
      <ol>
        <li>No <a href="/dashboard">Dashboard</a>, clique em <strong>Nova rifa</strong>.</li>
        <li>Preencha:
          <ul>
            <li><strong>Título</strong> e <strong>descrição</strong> — o que está sendo rifado.</li>
            <li><strong>Valor do número</strong> — preço de cada número.</li>
            <li><strong>Quantidade de números</strong> — total disponível.</li>
            <li><strong>Data do sorteio</strong>.</li>
            <li><strong>Imagem</strong> (opcional) — foto do prêmio.</li>
          </ul>
        </li>
        <li>Confira os dados e clique em <strong>Criar</strong>.</li>
      </ol>
      <h2>Depois de criar</h2>
      <p>A rifa fica pública na lista e você recebe um link para divulgar. Acompanhe as vendas pelo Dashboard.</p>
      <div class="tip">Importante: o organizador é o responsável pela legalidade da rifa e pela entrega do prêmio. Veja o <a href="/termo-do-organizador">Termo do Organizador</a>.</div>
    `,
  },
  {
    slug: "gerenciar-rifa",
    category: "organizador",
    title: "Editar, cancelar ou excluir uma rifa",
    summary: "Ajuste os dados ou encerre uma rifa.",
    body: `
      <h2>Editar</h2>
      <ol>
        <li>No <a href="/dashboard">Dashboard</a>, localize a rifa e clique em <strong>Editar</strong>.</li>
        <li>Altere os campos permitidos e salve.</li>
      </ol>
      <p><strong>Atenção:</strong> a edição fica limitada (ou bloqueada) quando a rifa já tem números vendidos, para proteger quem comprou.</p>
      <h2>Cancelar</h2>
      <p>Cancelar interrompe as vendas e marca a rifa como encerrada. Lembre-se de tratar eventuais reembolsos com os participantes — os pagamentos passam pela sua conta InfinitePay.</p>
      <h2>Excluir</h2>
      <p>A exclusão remove a rifa. Use com cuidado: rifas com vendas devem ser tratadas com transparência com os participantes.</p>
    `,
  },
  {
    slug: "acompanhar-vendas",
    category: "organizador",
    title: "Acompanhar vendas no Dashboard",
    summary: "Veja faturamento, números vendidos e estatísticas.",
    body: `
      <h2>O que você vê</h2>
      <ul>
        <li><strong>Resumo</strong> — total de rifas, vendas e faturamento.</li>
        <li><strong>Por rifa</strong> — números disponíveis, reservados e pagos.</li>
        <li><strong>Gráficos</strong> — evolução das vendas.</li>
      </ul>
      <h2>Como acessar</h2>
      <ol>
        <li>Entre em <a href="/dashboard">Dashboard</a>.</li>
        <li>Clique em uma rifa para ver estatísticas detalhadas.</li>
      </ol>
      <div class="tip">Números <strong>reservados</strong> são os que estão em processo de pagamento. Se o pagamento não se concretiza, eles voltam a ficar disponíveis automaticamente.</div>
    `,
  },
  {
    slug: "realizar-sorteio",
    category: "organizador",
    title: "Realizar o sorteio e divulgar o resultado",
    summary: "Sorteie o número vencedor e publique o resultado.",
    body: `
      <h2>Passo a passo</h2>
      <ol>
        <li>No <a href="/dashboard">Dashboard</a>, abra a rifa que chegou à data do sorteio.</li>
        <li>Clique em <strong>Sortear</strong>. O sistema escolhe um número de forma aleatória entre os pagos.</li>
        <li>O resultado é registrado e fica disponível na página pública de resultado da rifa.</li>
      </ol>
      <h2>Divulgação</h2>
      <p>Compartilhe o link de resultado com os participantes. A entrega do prêmio ao ganhador é responsabilidade do organizador.</p>
      <div class="tip">O sorteio é manual: ele acontece quando você clica em "Sortear". Garanta que a data anunciada seja respeitada.</div>
    `,
  },

  // ---------------- PARTICIPANTE ----------------
  {
    slug: "comprar-numeros",
    category: "participante",
    title: "Como comprar números de uma rifa",
    summary: "Escolha os números e finalize o pagamento — sem precisar de conta.",
    body: `
      <h2>Passo a passo</h2>
      <ol>
        <li>Abra a rifa desejada na lista pública.</li>
        <li>Selecione um ou mais números disponíveis.</li>
        <li>Clique em comprar e informe <strong>nome</strong> e <strong>telefone</strong>.</li>
        <li>Marque o aceite e clique em <strong>Ir para pagamento</strong>.</li>
        <li>Pague via PIX ou cartão na tela da InfinitePay.</li>
      </ol>
      <h2>Depois de pagar</h2>
      <p>O pagamento é confirmado automaticamente e seus números ficam registrados pelo telefone informado. Você não precisa criar conta.</p>
      <div class="tip">Guarde o telefone usado: é com ele que você acompanha suas compras. Veja <a href="/ajuda/minhas-compras">Minhas compras</a>.</div>
    `,
  },
  {
    slug: "minhas-compras",
    category: "participante",
    title: "Acompanhar Minhas Compras pelo telefone",
    summary: "Encontre todos os números que você comprou.",
    body: `
      <h2>Passo a passo</h2>
      <ol>
        <li>Acesse <a href="/my-purchases">Minhas compras</a>.</li>
        <li>Informe o <strong>telefone</strong> usado na compra.</li>
        <li>Veja a lista de rifas e números registrados nesse telefone.</li>
      </ol>
      <h2>Não encontrou sua compra?</h2>
      <ul>
        <li>Confirme se usou o <strong>mesmo telefone</strong> do pagamento.</li>
        <li>Se acabou de pagar, aguarde alguns instantes para a confirmação.</li>
        <li>Em caso de dúvida, fale com o organizador da rifa ou use a <a href="/contato">página de contato</a>.</li>
      </ul>
    `,
  },
  {
    slug: "ver-resultado",
    category: "participante",
    title: "Ver o resultado do sorteio",
    summary: "Descubra o número sorteado e se você ganhou.",
    body: `
      <h2>Como acompanhar</h2>
      <ol>
        <li>Abra a página da rifa.</li>
        <li>Após a data do sorteio, o <strong>número vencedor</strong> aparece na página de resultado.</li>
        <li>Compare com os números que você comprou em <a href="/my-purchases">Minhas compras</a>.</li>
      </ol>
      <h2>Ganhou?</h2>
      <p>A entrega do prêmio é feita pelo <strong>organizador</strong> da rifa. Combine com ele a forma de recebimento.</p>
    `,
  },

  // ---------------- PAGAMENTOS E ASSINATURA ----------------
  {
    slug: "assinatura-mensal",
    category: "pagamentos",
    title: "Assinatura mensal: como assinar e gerenciar",
    summary: "Mantenha sua conta de organizador ativa.",
    body: `
      <h2>Como assinar</h2>
      <ol>
        <li>Acesse <a href="/subscription">Assinatura</a>.</li>
        <li>Clique em assinar e conclua o pagamento via InfinitePay.</li>
        <li>Pronto: sua conta passa a ficar ativa e libera a criação de rifas.</li>
      </ol>
      <h2>Status da assinatura</h2>
      <ul>
        <li><strong>Trial</strong> — período de teste de 7 dias.</li>
        <li><strong>Ativa</strong> — assinatura em dia.</li>
        <li><strong>Pendente / vencida</strong> — pagamento em atraso.</li>
        <li><strong>Cancelada / inativa</strong> — sem acesso às funções de organizador.</li>
      </ul>
      <h2>Cancelamento</h2>
      <p>O cancelamento interrompe as renovações futuras. Você continua com acesso até o fim do período já pago.</p>
    `,
  },
  {
    slug: "pagamentos-infinitepay",
    category: "pagamentos",
    title: "Pagamentos e segurança (InfinitePay)",
    summary: "Como o dinheiro circula e por que seus dados ficam protegidos.",
    body: `
      <h2>Quem processa</h2>
      <p>Todos os pagamentos — assinaturas e compras de números — são processados pela <strong>InfinitePay</strong>. A Rifa Online não recebe nem guarda o dinheiro das rifas.</p>
      <h2>Para onde vai o dinheiro</h2>
      <ul>
        <li><strong>Compras de números</strong> → direto para a conta InfinitePay do organizador (handle configurado).</li>
        <li><strong>Assinatura</strong> → cobrança mensal da plataforma via InfinitePay.</li>
      </ul>
      <h2>Segurança</h2>
      <p>A Rifa Online <strong>não armazena dados de cartão</strong>. Esses dados são tratados diretamente na tela segura da InfinitePay.</p>
      <h2>Formas de pagamento</h2>
      <p>PIX (confirmação rápida) e cartão. A confirmação é automática após o pagamento.</p>
    `,
  },
  {
    slug: "integracao-pagamento-infinitepay",
    category: "pagamentos",
    title: "Integração de pagamento (InfinitePay) — referência",
    summary: "Cada componente da InfinitePay que o sistema usa e como integrar o seu recebimento.",
    body: `
      <h2>Nosso papel</h2>
      <p>A Rifa Online usa o <strong>Checkout da InfinitePay</strong> como meio de pagamento. Nós geramos os links de pagamento e confirmamos as transações; <strong>a InfinitePay processa o dinheiro</strong> e o repassa direto ao organizador. A plataforma <strong>não armazena dados de cartão</strong>.</p>

      <h2>Componentes que utilizamos</h2>
      <table>
        <thead><tr><th>Componente</th><th>Para que serve</th><th>Endpoint</th></tr></thead>
        <tbody>
          <tr><td><strong>Link de checkout</strong></td><td>Criar a página de pagamento de uma compra ou assinatura</td><td><code>POST /links</code></td></tr>
          <tr><td><strong>Verificação de pagamento</strong></td><td>Confirmar se um pedido foi pago (fallback do webhook)</td><td><code>POST /payment_check</code></td></tr>
          <tr><td><strong>Webhook</strong></td><td>Receber a notificação automática quando o pagamento é aprovado</td><td>seu endpoint <code>webhook_url</code></td></tr>
          <tr><td><strong>InfiniteTag (handle)</strong></td><td>Identifica a conta que <em>recebe</em> o pagamento</td><td>campo <code>handle</code></td></tr>
          <tr><td><strong>Redirect URL</strong></td><td>Para onde o cliente volta após pagar</td><td>campo <code>redirect_url</code></td></tr>
        </tbody>
      </table>
      <p>A base da API utilizada é <code>https://api.checkout.infinitepay.io</code>.</p>

      <h2>Ciclo de vida de um pagamento</h2>
      <ol>
        <li>O participante seleciona números e informa nome e telefone.</li>
        <li>O sistema cria o checkout (<code>POST /links</code>) com o <code>handle</code> do organizador, os <code>items</code>, um <code>order_nsu</code> único, o <code>redirect_url</code> e o <code>webhook_url</code>.</li>
        <li>A InfinitePay devolve a <code>url</code> e o participante é levado à página de pagamento (PIX ou cartão).</li>
        <li>Ao aprovar, a InfinitePay chama o nosso <code>webhook_url</code> com o <code>order_nsu</code>; marcamos os números como pagos (de forma idempotente).</li>
        <li>Como reforço, o sistema também consulta <code>POST /payment_check</code> para confirmar o status.</li>
        <li>O participante retorna pelo <code>redirect_url</code> à tela de sucesso.</li>
      </ol>

      <h2>Exemplo — criação de checkout</h2>
      <pre><code>POST https://api.checkout.infinitepay.io/links
{
  "handle": "suatag",
  "order_nsu": "rifa_123_abc",
  "redirect_url": "https://seusite.com/payment/success",
  "webhook_url": "https://seusite.com/api/v1/webhooks/infinitepay",
  "items": [
    { "quantity": 2, "price": 500, "description": "Rifa X - numeros 7, 22" }
  ],
  "customer": { "name": "Maria", "phone_number": "11999999999" }
}</code></pre>
      <p>Os valores em <code>price</code> são em <strong>centavos</strong> (ex.: <code>500</code> = R$ 5,00).</p>

      <h2>Exemplo — corpo do webhook</h2>
      <pre><code>{
  "invoice_slug": "abc123",
  "order_nsu": "rifa_123_abc",
  "amount": 1000,
  "paid_amount": 1000,
  "installments": 1,
  "capture_method": "pix",
  "transaction_nsu": "xyz789",
  "receipt_url": "https://...",
  "items": [ ... ]
}</code></pre>
      <table>
        <thead><tr><th>Campo</th><th>Descrição</th></tr></thead>
        <tbody>
          <tr><td><code>order_nsu</code></td><td>O identificador que você gerou — liga o pagamento à compra.</td></tr>
          <tr><td><code>invoice_slug</code></td><td>Identificador da fatura na InfinitePay.</td></tr>
          <tr><td><code>amount</code> / <code>paid_amount</code></td><td>Valor previsto / efetivamente pago (centavos).</td></tr>
          <tr><td><code>installments</code></td><td>Número de parcelas.</td></tr>
          <tr><td><code>capture_method</code></td><td>Forma de pagamento (ex.: <code>pix</code>, <code>credit_card</code>).</td></tr>
          <tr><td><code>transaction_nsu</code></td><td>Identificador da transação.</td></tr>
          <tr><td><code>receipt_url</code></td><td>Link do comprovante.</td></tr>
        </tbody>
      </table>
      <div class="tip">Boas práticas de webhook: responda rápido (idealmente em menos de 1 segundo) e trate a notificação de forma <strong>idempotente</strong> — o mesmo <code>order_nsu</code> pode chegar mais de uma vez.</div>

      <h2>Como o cliente integra o recebimento</h2>
      <ol>
        <li>Crie uma conta na <strong>InfinitePay</strong> pelo aplicativo dela.</li>
        <li>Pegue o seu <strong>InfiniteTag</strong> (seu nome de usuário na InfinitePay).</li>
        <li>Informe esse handle no <a href="/profile">Perfil</a> da Rifa Online (veja <a href="/ajuda/configurar-recebimento">Configurar recebimento</a>).</li>
        <li>A partir daí, toda compra das suas rifas cai direto na sua conta InfinitePay.</li>
        <li>Faça um teste com um valor pequeno (PIX) e confirme em <a href="/my-purchases">Minhas compras</a> e no <a href="/dashboard">Dashboard</a>.</li>
      </ol>
      <p>Observação: na API, o handle é usado <strong>sem o símbolo <code>$</code></strong> no início.</p>

      <h2>Configuração técnica (auto-hospedagem)</h2>
      <p>Para quem opera a própria instância, a integração é configurada por variáveis de ambiente:</p>
      <ul>
        <li><code>INFINITEPAY_BASE_URL</code> — base da API (<code>https://api.checkout.infinitepay.io</code>).</li>
        <li><code>INFINITEPAY_HANDLE</code> — handle da plataforma (usado nas assinaturas).</li>
        <li>O <code>webhook_url</code> é montado automaticamente como <code>FRONTEND_URL + /api/v1/webhooks/infinitepay</code>.</li>
      </ul>

      <h2>Referências oficiais</h2>
      <ul>
        <li><a href="https://www.infinitepay.io/checkout-documentacao" target="_blank" rel="noopener">Documentação do Checkout InfinitePay</a></li>
        <li><a href="https://www.infinitepay.io/desenvolvedores" target="_blank" rel="noopener">Integração para desenvolvedores</a></li>
        <li><a href="https://ajuda.infinitepay.io/pt-BR/articles/10766888-como-usar-o-checkout-da-infinitepay" target="_blank" rel="noopener">Central de Ajuda — Como usar o Checkout</a></li>
      </ul>
    `,
  },

  // ---------------- ADMINISTRAÇÃO ----------------
  {
    slug: "painel-admin",
    category: "administracao",
    title: "Painel administrativo",
    summary: "Gerencie usuários, assinaturas, rifas e veja estatísticas.",
    body: `
      <h2>Acesso</h2>
      <p>Disponível apenas para contas com perfil <strong>Admin</strong>, em <a href="/admin">Admin</a>.</p>
      <h2>Abas</h2>
      <ul>
        <li><strong>Resumo</strong> — estatísticas globais da plataforma.</li>
        <li><strong>Usuários</strong> — busca, detalhes e alteração de status de assinatura.</li>
        <li><strong>Rifas</strong> — visão geral, com ações de sortear/cancelar.</li>
        <li><strong>Mensagens</strong> — mensagens recebidas pela página de contato.</li>
      </ul>
      <h2>Alterar assinatura de um usuário</h2>
      <ol>
        <li>Abra a aba <strong>Usuários</strong>.</li>
        <li>Use os botões de status (ativar, cancelar, vencer, zerar) na linha do usuário.</li>
      </ol>
    `,
  },
  {
    slug: "mensagens-de-contato",
    category: "administracao",
    title: "Mensagens de contato recebidas",
    summary: "Leia o que os usuários enviam pela página de contato.",
    body: `
      <h2>Onde ver</h2>
      <ol>
        <li>Acesse <a href="/admin">Admin</a>.</li>
        <li>Abra a aba <strong>Mensagens</strong>.</li>
        <li>Cada item mostra nome, contato (se informado), mensagem e data.</li>
      </ol>
      <h2>Como responder</h2>
      <p>Se a pessoa informou e-mail ou telefone, responda por esse meio. Solicitações de privacidade (LGPD) ficam registradas com data para acompanhamento.</p>
    `,
  },

  // ---------------- CONTA, PRIVACIDADE E SUPORTE ----------------
  {
    slug: "gerenciar-perfil",
    category: "conta-suporte",
    title: "Gerenciar perfil e conta",
    summary: "Atualize seus dados e o recebimento.",
    body: `
      <h2>O que você pode ajustar</h2>
      <ul>
        <li>Nome e telefone.</li>
        <li>Handle InfinitePay (recebimento das rifas).</li>
      </ul>
      <h2>Passo a passo</h2>
      <ol>
        <li>Acesse <a href="/profile">Perfil</a>.</li>
        <li>Edite os campos desejados.</li>
        <li>Salve as alterações.</li>
      </ol>
    `,
  },
  {
    slug: "privacidade-lgpd",
    category: "conta-suporte",
    title: "Privacidade e seus direitos (LGPD)",
    summary: "Quais dados são tratados e como exercer seus direitos.",
    body: `
      <h2>Dados tratados</h2>
      <ul>
        <li><strong>Organizadores</strong> — dados de cadastro e da assinatura.</li>
        <li><strong>Participantes</strong> — nome e telefone, para identificar a compra.</li>
      </ul>
      <h2>Seus direitos</h2>
      <p>Você pode solicitar acesso, correção ou exclusão dos seus dados, entre outros direitos previstos na LGPD. Detalhes na <a href="/politica-de-privacidade">Política de Privacidade</a>.</p>
      <h2>Como solicitar</h2>
      <p>Use a <a href="/contato">página de contato</a>. Pedidos sobre dados de participantes também podem ser direcionados ao organizador da rifa.</p>
      <h2>Documentos</h2>
      <ul>
        <li><a href="/termos-de-uso">Termos de Uso</a></li>
        <li><a href="/politica-de-privacidade">Política de Privacidade</a></li>
        <li><a href="/politica-de-cookies">Política de Cookies</a></li>
        <li><a href="/termo-do-organizador">Termo do Organizador</a></li>
      </ul>
    `,
  },
  {
    slug: "suporte-e-contato",
    category: "conta-suporte",
    title: "Suporte e contato",
    summary: "Como falar com a Rifa Online.",
    body: `
      <h2>Canal de atendimento</h2>
      <p>Use a <a href="/contato">página de contato</a> para enviar dúvidas, problemas ou solicitações. Se quiser retorno, informe um e-mail ou telefone.</p>
      <h2>Antes de enviar</h2>
      <ul>
        <li>Procure sua dúvida aqui na Central de Ajuda.</li>
        <li>Para dúvidas sobre uma rifa específica (prêmio, sorteio, entrega), fale também com o organizador.</li>
      </ul>
    `,
  },
]

export function getArticle(slug: string): HelpArticle | undefined {
  return helpArticles.find((a) => a.slug === slug)
}

export function articlesByCategory(categoryId: string): HelpArticle[] {
  return helpArticles.filter((a) => a.category === categoryId)
}
