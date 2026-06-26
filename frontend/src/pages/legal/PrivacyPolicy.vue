<script setup lang="ts">
import LegalLayout from "../../components/LegalLayout.vue"
</script>

<template>
  <LegalLayout title="Política de Privacidade" updated-at="26 de junho de 2026">
    <p>
      Esta Política de Privacidade descreve como a <strong>Rifa Online</strong> ("Plataforma"),
      operada sob a marca <strong>Rifa Online</strong> ("nós"), trata dados pessoais, em conformidade
      com a Lei nº 13.709/2018 (Lei Geral de Proteção de Dados — LGPD) e o Marco Civil da Internet
      (Lei nº 12.965/2014).
    </p>

    <h2>1. Papéis no tratamento de dados</h2>
    <p>
      Para fins da LGPD, atuamos em dois papéis distintos:
    </p>
    <ul>
      <li><strong>Controladores</strong> dos dados de cadastro dos Organizadores (conta, assinatura,
        uso da Plataforma);</li>
      <li><strong>Operadores</strong> dos dados dos Participantes coletados em nome do Organizador
        (nome e telefone do comprador), os quais são <strong>controlados pelo Organizador</strong> da
        rifa. O Organizador é responsável por possuir base legal e informar seus Participantes.</li>
    </ul>

    <h2>2. Dados que coletamos</h2>
    <div class="table-wrap"><table>
      <thead>
        <tr><th>Categoria</th><th>Dados</th><th>Origem</th></tr>
      </thead>
      <tbody>
        <tr><td>Cadastro do Organizador</td><td>Nome, e-mail, senha (criptografada), telefone, handle InfinitePay</td><td>Fornecidos pelo Usuário</td></tr>
        <tr><td>Assinatura</td><td>Status, datas, identificadores de pagamento</td><td>Plataforma / InfinitePay</td></tr>
        <tr><td>Participantes</td><td>Nome e telefone do comprador, números adquiridos</td><td>Fornecidos no checkout</td></tr>
        <tr><td>Técnicos</td><td>Endereço IP, data/hora de acesso, registros de uso, cookies</td><td>Coleta automática</td></tr>
      </tbody>
    </table></div>
    <p>
      <strong>Não coletamos nem armazenamos dados de cartão de crédito.</strong> O pagamento é processado
      diretamente pela InfinitePay.
    </p>

    <h2>3. Finalidades e bases legais</h2>
    <div class="table-wrap"><table>
      <thead>
        <tr><th>Finalidade</th><th>Base legal (art. 7º LGPD)</th></tr>
      </thead>
      <tbody>
        <tr><td>Criar e manter a conta e prestar o serviço</td><td>Execução de contrato (II)</td></tr>
        <tr><td>Processar assinaturas e pagamentos</td><td>Execução de contrato / obrigação legal</td></tr>
        <tr><td>Identificar o comprador e seus números</td><td>Execução de contrato (entre Organizador e Participante)</td></tr>
        <tr><td>Segurança, prevenção a fraudes e registros de acesso</td><td>Cumprimento de obrigação legal (Marco Civil) / legítimo interesse</td></tr>
        <tr><td>Comunicações sobre o serviço</td><td>Execução de contrato / legítimo interesse</td></tr>
        <tr><td>Comunicações de marketing</td><td>Consentimento (I)</td></tr>
      </tbody>
    </table></div>

    <h2>4. Compartilhamento de dados</h2>
    <p>Podemos compartilhar dados com:</p>
    <ul>
      <li><strong>InfinitePay</strong> — para processar pagamentos e assinaturas;</li>
      <li><strong>Organizadores</strong> — que recebem os dados dos Participantes de suas rifas;</li>
      <li><strong>Provedores de infraestrutura</strong> (hospedagem, e-mail, analytics) que atuam como
        operadores sob contrato;</li>
      <li><strong>Autoridades</strong> — quando exigido por lei, ordem judicial ou requisição legítima.</li>
    </ul>
    <p>Não vendemos dados pessoais.</p>

    <h2>5. Armazenamento e retenção</h2>
    <p>
      Os dados são mantidos pelo tempo necessário às finalidades acima e às obrigações legais. Registros
      de acesso a aplicações são guardados por, no mínimo, 6 meses (art. 15 do Marco Civil). Dados de
      cadastro são retidos enquanto a conta existir e por prazos legais de guarda após o encerramento.
    </p>

    <h2>6. Segurança</h2>
    <p>
      Adotamos medidas técnicas e organizacionais razoáveis para proteger os dados, como criptografia de
      senhas, controle de acesso e tráfego por conexão segura (HTTPS). Nenhum sistema é totalmente imune;
      em caso de incidente de segurança relevante, comunicaremos os titulares e a ANPD conforme a lei.
    </p>

    <h2>7. Direitos do titular</h2>
    <p>Nos termos do art. 18 da LGPD, você pode solicitar:</p>
    <ul>
      <li>confirmação da existência de tratamento e acesso aos dados;</li>
      <li>correção de dados incompletos, inexatos ou desatualizados;</li>
      <li>anonimização, bloqueio ou eliminação de dados desnecessários ou tratados em desconformidade;</li>
      <li>portabilidade e informação sobre compartilhamentos;</li>
      <li>revogação do consentimento e eliminação dos dados tratados com base nele.</li>
    </ul>
    <p>
      Para dados de Participantes, as solicitações podem ser direcionadas ao Organizador da rifa
      (controlador). Demais pedidos podem ser feitos pelos contatos abaixo.
    </p>

    <h2>8. Cookies</h2>
    <p>
      Utilizamos cookies e tecnologias semelhantes conforme nossa
      <router-link to="/politica-de-cookies">Política de Cookies</router-link>.
    </p>

    <h2>9. Transferência internacional</h2>
    <p>
      Caso provedores estejam localizados fora do Brasil, a transferência observará as hipóteses e
      salvaguardas dos arts. 33 a 36 da LGPD.
    </p>

    <h2>10. Encarregado (DPO) e contato</h2>
    <p>
      O Encarregado pelo Tratamento de Dados (DPO) pode ser contatado pela nossa
      <router-link to="/contato">página de contato</router-link>.
    </p>
    <p>
      Você também pode apresentar reclamação à Autoridade Nacional de Proteção de Dados (ANPD).
    </p>

    <h2>11. Alterações</h2>
    <p>
      Esta Política pode ser atualizada periodicamente; a versão vigente é sempre a publicada nesta página,
      com a data de atualização.
    </p>
  </LegalLayout>
</template>
