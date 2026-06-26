# PRD — Rifa Online

## 1. Visão Geral

Sistema web (SaaS) para gerenciamento de rifas online com pagamento via **InfinitePay**. Organizadores assinam o serviço para criar rifas, vender números e gerenciar sorteios; participantes compram números (sem necessidade de conta) e acompanham resultados. Há um painel administrativo para operação da plataforma.

## 2. Modelo de Negócio (SaaS)

- Organizadores pagam **assinatura mensal de R$ 10,00** (via InfinitePay) para usar a plataforma.
- **Trial de 7 dias** ao criar a conta — durante o trial o organizador já pode criar e gerenciar rifas.
- Sem assinatura ativa (ou trial expirado), a criação/gestão de rifas é bloqueada.
- **Admin** ignora a verificação de assinatura (bypass).

## 3. Objetivos

- Facilitar a criação e gestão de rifas por qualquer organizador assinante
- Processar pagamentos de forma segura via PIX/cartão (InfinitePay)
- Automatizar a marcação de números pagos via webhook + confirmação
- Sortear números aleatórios e divulgar resultados
- Dashboard com relatórios e estatísticas
- Monetizar via assinatura mensal recorrente

## 4. Personas

| Persona | Descrição |
|---------|-----------|
| **Organizador** | Assina o serviço, cria e gerencia rifas, visualiza vendas, realiza sorteios |
| **Participante** | Compra números (sem conta), acompanha rifas, vê resultados e suas compras |
| **Admin** | Opera a plataforma: gerencia usuários, assinaturas, rifas e vê estatísticas globais |

## 5. Funcionalidades

### 5.1 Autenticação

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-01 | Cadastro de organizador (nome, email, senha) com trial de 7 dias | Alta |
| F-02 | Login com JWT (access + refresh) | Alta |
| F-03 | Refresh de token | Alta |
| F-04 | Perfil do usuário (nome, telefone, handle InfinitePay) | Alta |

### 5.2 Assinatura (SaaS)

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-05 | Checkout de assinatura via InfinitePay (R$10/mês) | Alta |
| F-06 | Status da assinatura (ACTIVE/INACTIVE/PAST_DUE/CANCELLED, trial, expiração) | Alta |
| F-07 | Botão de assinar exibido durante o trial | Média |
| F-08 | Bloqueio de criação de rifa sem assinatura/trial | Alta |

### 5.3 Rifas

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-09 | Criar rifa (título, descrição, valor, data, quantidade, imagem) | Alta |
| F-10 | Listar rifas públicas | Alta |
| F-11 | Visualizar detalhes (números disponíveis/reservados/pagos) | Alta |
| F-12 | Editar rifa (bloqueado se já houver vendas) | Média |
| F-13 | Cancelar / excluir rifa | Média |
| F-14 | Compartilhar rifa | Baixa |
| F-15 | Estatísticas da rifa | Média |

### 5.4 Compra (Participante)

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-16 | Selecionar números disponíveis | Alta |
| F-17 | Checkout via InfinitePay (PIX / Cartão) — informa nome e **telefone** | Alta |
| F-18 | Webhook + confirmação automática de pagamento | Alta |
| F-19 | "Minhas compras" — buscar tickets pelo telefone | Alta |

### 5.5 Sorteio

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-20 | Sortear vencedor (manual) | Alta |
| F-21 | Registrar e divulgar resultado | Alta |

### 5.6 Dashboard & Admin

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-22 | Dashboard do organizador (vendas, rifas, faturamento, gráficos) | Alta |
| F-23 | Painel admin: usuários, busca, rifas, estatísticas globais | Alta |
| F-24 | Admin alterar assinatura de um usuário | Média |

## 6. Fluxo do Participante

1. Acessa a lista de rifas públicas
2. Escolhe uma rifa e visualiza os números
3. Seleciona um ou mais números
4. Informa **nome e telefone** (o telefone identifica o comprador)
5. É redirecionado ao checkout InfinitePay
6. Paga via PIX ou cartão
7. Volta ao site; pagamento confirmado por webhook/polling
8. Acompanha "Minhas compras" pelo telefone e o sorteio na data marcada

## 7. Fluxo do Organizador

1. Cadastra-se (entra em trial de 7 dias) e faz login
2. (Durante/após o trial) assina o plano via InfinitePay
3. Configura seu **handle InfinitePay** no perfil (para receber os pagamentos das rifas)
4. Cria rifa (dados, valor, data)
5. Acompanha vendas no dashboard
6. Sorteia o vencedor e divulga o resultado

## 8. Critérios de Sucesso

- Pagamento confirmado em < 30s após o pagamento (webhook + confirmação)
- Organizador cria rifa em < 5 passos
- Sistema lida com 1000 participantes simultâneos

## 9. Não-Escopo

- App mobile nativo
- Split de pagamento entre organizadores
- Rifas recorrentes
- Marketplace de rifas
- Chat organizador ↔ participante
- Sorteio automático agendado (apenas manual no momento)

## 10. Glossário

| Termo | Definição |
|-------|-----------|
| Rifa | Sorteio onde participantes compram números |
| Número / Ticket | Unidade vendável de uma rifa |
| Handle InfinitePay | Identificador (`$tag`) do organizador que recebe os pagamentos |
| Checkout | Página de pagamento hospedada pela InfinitePay |
| Webhook | Notificação HTTP enviada pela InfinitePay sobre eventos de pagamento |
| Trial | Período de 7 dias de uso gratuito após o cadastro |
