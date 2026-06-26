<script setup lang="ts">
import LegalLayout from "../../components/LegalLayout.vue"
</script>

<template>
  <LegalLayout title="Termos de Uso" updated-at="26 de junho de 2026">
    <p>
      Estes Termos de Uso ("Termos") regem o acesso e a utilização da plataforma
      <strong>Rifa Online</strong> ("Rifa Online", "Plataforma", "nós" ou "nosso").
    </p>
    <p>
      Ao se cadastrar, acessar ou utilizar a Plataforma, você ("Usuário") declara que leu,
      compreendeu e concorda integralmente com estes Termos e com a
      <router-link to="/politica-de-privacidade">Política de Privacidade</router-link>.
      Caso não concorde, não utilize a Plataforma.
    </p>

    <div class="callout">
      <strong>Aviso importante sobre rifas e sorteios.</strong> A distribuição gratuita de prêmios
      mediante sorteio, vale-brinde, concurso ou operação assemelhada é atividade regulada no Brasil
      (Lei nº 5.768/1971, Decreto nº 70.951/1972 e normas da Secretaria de Prêmios e Apostas/Ministério
      da Fazenda, além da legislação de loterias). A Plataforma é apenas uma <strong>ferramenta
      tecnológica</strong>. A obtenção de autorizações, licenças e o cumprimento de toda a legislação
      aplicável são de <strong>responsabilidade exclusiva do Organizador</strong> da rifa.
    </div>

    <h2>1. Definições</h2>
    <ul>
      <li><strong>Plataforma:</strong> o sistema web Rifa Online e seus serviços associados.</li>
      <li><strong>Organizador:</strong> usuário cadastrado que assina o serviço para criar e gerir rifas.</li>
      <li><strong>Participante:</strong> pessoa que adquire números de uma rifa criada por um Organizador.</li>
      <li><strong>InfinitePay:</strong> instituição de pagamento terceira responsável pelo processamento financeiro.</li>
      <li><strong>Assinatura:</strong> plano mensal pago pelo Organizador para uso da Plataforma.</li>
    </ul>

    <h2>2. Natureza do serviço (papel da Plataforma)</h2>
    <p>
      A Plataforma é um software como serviço (SaaS) que disponibiliza ferramentas para que
      Organizadores criem, divulguem e administrem rifas. <strong>A Rifa Online não é promotora,
      organizadora, patrocinadora ou realizadora das rifas</strong> publicadas pelos Usuários.
    </p>
    <p>A Rifa Online <strong>não</strong>:</p>
    <ul>
      <li>vende números de rifa em nome próprio;</li>
      <li>recebe, retém, custodia ou administra os valores pagos pelos Participantes — os pagamentos
        são processados diretamente pela InfinitePay e destinados ao Organizador (handle InfinitePay);</li>
      <li>realiza os sorteios nem garante a sua lisura;</li>
      <li>entrega prêmios ou garante a sua existência, qualidade ou entrega;</li>
      <li>obtém autorizações regulatórias em nome do Organizador.</li>
    </ul>
    <p>
      A relação jurídica decorrente da rifa estabelece-se <strong>exclusivamente entre Organizador e
      Participante</strong>. A Rifa Online atua apenas como provedora de aplicação, nos termos do art. 19
      do Marco Civil da Internet (Lei nº 12.965/2014).
    </p>

    <h2>3. Cadastro e conta</h2>
    <ul>
      <li>O cadastro é destinado a maiores de 18 anos, plenamente capazes.</li>
      <li>O Usuário compromete-se a fornecer informações verdadeiras, completas e atualizadas.</li>
      <li>As credenciais de acesso são pessoais e intransferíveis; o Usuário é responsável por toda
        atividade realizada em sua conta.</li>
      <li>A Rifa Online pode recusar, suspender ou cancelar contas que violem estes Termos ou a legislação.</li>
    </ul>

    <h2>4. Assinatura, pagamento e trial</h2>
    <ul>
      <li>O uso das funções de Organizador exige assinatura mensal, conforme valor informado na Plataforma.</li>
      <li>É oferecido período de teste gratuito ("trial") por tempo determinado, após o qual a criação e
        gestão de rifas é bloqueada sem assinatura ativa.</li>
      <li>Os pagamentos da assinatura são processados pela InfinitePay; a Rifa Online não armazena dados
        de cartão.</li>
      <li>A assinatura é renovada automaticamente até cancelamento pelo Usuário. O cancelamento
        interrompe renovações futuras, sem reembolso de períodos já iniciados, salvo exigência legal.</li>
    </ul>

    <h2>5. Obrigações e responsabilidade do Organizador</h2>
    <p>Ao publicar uma rifa, o Organizador declara e garante que:</p>
    <ul>
      <li>possui todas as autorizações, licenças e registros legais exigíveis para a realização da rifa
        ou promoção, e é o único responsável por obtê-los e mantê-los;</li>
      <li>cumprirá integralmente a legislação aplicável, incluindo normas tributárias, consumeristas,
        de proteção de dados e de sorteios/promoções;</li>
      <li>é o legítimo proprietário ou está autorizado a oferecer o prêmio anunciado, responsabilizando-se
        por sua existência, qualidade e entrega ao vencedor;</li>
      <li>prestará informações verdadeiras sobre a rifa, o prêmio, as regras e a data do sorteio;</li>
      <li>tratará os dados pessoais dos Participantes de forma lícita, na condição de controlador desses
        dados, conforme a <router-link to="/politica-de-privacidade">Política de Privacidade</router-link>;</li>
      <li>não utilizará a Plataforma para fins ilícitos, fraudulentos, de lavagem de dinheiro ou que
        violem direitos de terceiros.</li>
    </ul>
    <p>
      O detalhamento dessas responsabilidades consta do
      <router-link to="/termo-do-organizador">Termo de Responsabilidade do Organizador</router-link>,
      que integra estes Termos.
    </p>

    <h2>6. Condutas proibidas</h2>
    <p>É vedado a qualquer Usuário:</p>
    <ul>
      <li>criar rifas ou promoções ilegais, sem autorização exigida, ou com prêmios proibidos;</li>
      <li>publicar conteúdo falso, enganoso, difamatório, discriminatório ou que viole direitos de terceiros;</li>
      <li>utilizar a Plataforma para golpes, estelionato, lavagem de dinheiro ou financiamento ilícito;</li>
      <li>tentar burlar, sobrecarregar, copiar, sondar ou comprometer a segurança da Plataforma;</li>
      <li>coletar dados de outros usuários sem base legal.</li>
    </ul>

    <h2>7. Pagamentos por terceiros (InfinitePay)</h2>
    <p>
      Todo o processamento financeiro é realizado pela <strong>InfinitePay</strong>, instituição de
      pagamento independente, sujeita aos seus próprios termos e políticas. A Rifa Online não tem ingerência
      sobre aprovação, recusa, estorno (chargeback), liquidação ou repasse de valores, que se regem pela
      relação entre Participante, Organizador e InfinitePay.
    </p>

    <h2>8. Isenção de garantias</h2>
    <p>
      A Plataforma é fornecida "no estado em que se encontra" e "conforme disponível". A Rifa Online não
      garante que o serviço será ininterrupto, livre de erros ou de que atenderá a finalidades específicas.
      A Rifa Online não se responsabiliza por condutas, conteúdos, rifas, sorteios, prêmios ou pagamentos de
      Organizadores e Participantes.
    </p>

    <h2>9. Limitação de responsabilidade</h2>
    <p>
      Na máxima extensão permitida pela lei, a Rifa Online não responde por danos indiretos, lucros cessantes,
      perda de dados ou prejuízos decorrentes de: (i) rifas, sorteios e prêmios; (ii) atos de Organizadores
      ou Participantes; (iii) falhas de terceiros, incluindo a InfinitePay; (iv) indisponibilidade temporária
      do serviço. Caso reconhecida alguma responsabilidade da Rifa Online, esta limitar-se-á ao valor das
      assinaturas pagas pelo Usuário nos 12 meses anteriores ao evento.
    </p>

    <h2>10. Indenização</h2>
    <p>
      O Usuário concorda em <strong>indenizar e isentar</strong> a Rifa Online, seus sócios, administradores e
      colaboradores de quaisquer reclamações, perdas, danos, multas, custos e honorários advocatícios
      decorrentes do uso da Plataforma, da violação destes Termos ou da legislação, especialmente em razão
      de rifas, sorteios e promoções criadas pelo Usuário.
    </p>

    <h2>11. Propriedade intelectual</h2>
    <p>
      A marca, o software, o design e os demais elementos da Plataforma pertencem à Rifa Online e são protegidos
      por lei. É concedida ao Usuário licença limitada, não exclusiva e revogável de uso, vedada qualquer
      reprodução ou exploração não autorizada.
    </p>

    <h2>12. Suspensão e encerramento</h2>
    <p>
      A Rifa Online pode suspender ou encerrar o acesso de Usuários que violem estes Termos ou a legislação,
      bem como remover conteúdos e rifas irregulares, a qualquer tempo e sem aviso prévio quando necessário
      para cumprir a lei ou proteger terceiros.
    </p>

    <h2>13. Alterações dos Termos</h2>
    <p>
      Estes Termos podem ser atualizados a qualquer momento. A versão vigente será sempre publicada nesta
      página, com a respectiva data. O uso continuado após alterações implica concordância.
    </p>

    <h2>14. Lei aplicável e foro</h2>
    <p>
      Estes Termos são regidos pelas leis da República Federativa do Brasil. Fica eleito o foro do
      domicílio do Usuário ou do consumidor para dirimir eventuais controvérsias.
    </p>

    <h2>15. Contato</h2>
    <p>
      Dúvidas sobre estes Termos podem ser encaminhadas pela nossa
      <router-link to="/contato">página de contato</router-link>.
    </p>
  </LegalLayout>
</template>
