# PRD — Rifa Online

## 1. Visão Geral

Sistema web para gerenciamento de rifas online com pagamento via AbacatePay. Permite que organizadores criem rifas, vendam números e gerenciem sorteios, e que participantes comprem números e acompanhem resultados.

## 2. Objetivos

- Facilitar a criação e gestão de rifas por qualquer pessoa
- Processar pagamentos de forma segura via PIX/cartão (AbacatePay)
- Automatizar a marcação de números pagos
- Sortear números aleatórios e divulgar resultados
- Dashboard administrativo com relatórios básicos

## 3. Personas

| Persona | Descrição |
|---------|-----------|
| **Organizador** | Cria e gerencia rifas, visualiza vendas, realiza sorteios |
| **Participante** | Compra números, acompanha rifas, vê resultados |

## 4. Funcionalidades (MVP)

### 4.1 Módulo de Autenticação

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-01 | Cadastro de organizador (email + senha) | Alta |
| F-02 | Login com JWT | Alta |
| F-03 | Recuperação de senha | Média |

### 4.2 Módulo de Rifas

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-04 | Criar rifa (título, descrição, valor do número, data do sorteio, quantidade de números, imagem) | Alta |
| F-05 | Listar rifas públicas (com filtro por status) | Alta |
| F-06 | Visualizar detalhes de uma rifa (números disponíveis/ocupados) | Alta |
| F-07 | Editar rifa (antes do sorteio) | Média |
| F-08 | Cancelar rifa | Média |

### 4.3 Módulo de Compra (Participante)

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-09 | Selecionar números disponíveis em uma rifa | Alta |
| F-10 | Checkout via AbacatePay (PIX / Cartão) | Alta |
| F-11 | Webhook para confirmação automática de pagamento | Alta |
| F-12 | Visualizar meus números comprados | Alta |

### 4.4 Módulo de Sorteio

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-13 | Sortear vencedor (automático na data agendada ou manual) | Alta |
| F-14 | Registrar resultado do sorteio | Alta |
| F-15 | Notificar vencedor (por email) | Média |

### 4.5 Dashboard

| ID | Funcionalidade | Prioridade |
|----|---------------|------------|
| F-16 | Painel do organizador com vendas, rifas ativas, faturamento | Alta |
| F-17 | Lista de participantes por rifa | Média |

## 5. Fluxo do Usuário (Participante)

1. Acessa lista de rifas públicas
2. Escolhe uma rifa e visualiza os números
3. Seleciona um ou mais números
4. Informa nome e email
5. É redirecionado ao checkout AbacatePay
6. Paga via PIX ou cartão
7. É redirecionado de volta ao site com confirmação
8. Recebe confirmação por email (quando implementado)
9. Acompanha o sorteio na data marcada

## 6. Fluxo do Organizador

1. Cadastra-se e faz login
2. Cria rifa (define dados, valor, data)
3. Acompanha vendas no dashboard
4. Sorteia o vencedor (manual ou automático)
5. Divulga o resultado

## 7. Critérios de Sucesso

- Uma compra é confirmada em < 30s após o pagamento via webhook
- Organizador consegue criar rifa em < 5 passos
- Sistema lida com 1000 participantes simultâneos

## 8. Não-Escopo (MVP)

- App mobile nativo
- Split de pagamento entre organizadores
- Rifas recorrentes
- Marketplace de rifas
- Chat entre organizador e participantes

## 9. Glossário

| Termo | Definição |
|-------|-----------|
| Rifa | Sorteio onde participantes compram números |
| Número | Unidade vendável de uma rifa |
| Checkout | Página de pagamento hospedada pela AbacatePay |
| Webhook | Notificação HTTP enviada pela AbacatePay sobre eventos de pagamento |
