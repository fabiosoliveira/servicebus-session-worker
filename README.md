# 🚀 Azure Service Bus Session Worker

> **Processamento assíncrono de alta performance com sessões do Azure Service Bus**

[![Go Version](https://img.shields.io/badge/Go-1.24.4-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Azure Service Bus](https://img.shields.io/badge/Azure-Service%20Bus-0078D4?style=for-the-badge&logo=microsoft-azure)](https://azure.microsoft.com/services/service-bus/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

## 📋 Visão Geral

Este projeto implementa um **worker distribuído** para processamento de mensagens com sessões do Azure Service Bus, demonstrando conceitos avançados de programação concorrente em Go e integração com serviços de mensageria empresariais.

### ✨ Características Principais

- 🔄 **Processamento Concorrente**: Controle de workers através de semáforos
- 🎯 **Session-Based Processing**: Garante ordem e afinidade de processamento
- 🛡️ **Graceful Shutdown**: Encerramento seguro com tratamento de sinais
- ⚡ **Alta Performance**: Pool de workers configurável
- 🔒 **Error Handling**: Tratamento robusto de erros e timeouts
- 📊 **Observabilidade**: Logs estruturados para monitoramento

## 🏗️ Arquitetura

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Sender App    │───▶│  Azure Service   │◀───│  Consumer App   │
│                 │    │      Bus         │    │                 │
│ • Envia msgs    │    │                  │    │ • Pool Workers  │
│ • 100 sessões   │    │ • Session Queue  │    │ • Graceful Stop │
│ • Batch send    │    │ • FIFO per       │    │ • Auto-scaling  │
└─────────────────┘    │   session        │    └─────────────────┘
                       └──────────────────┘
```

## 🚀 Como Executar

### Pré-requisitos

```bash
# Go 1.24.4+
go version

# Azure Service Bus configurado
az servicebus namespace create --name <namespace> --resource-group <rg>
```

### Configuração

```bash
# Variáveis de ambiente necessárias
export SERVICEBUS_CONNECTION_STRING="Endpoint=sb://..."
export SERVICEBUS_QUEUE_NAME="session-queue"
export WORKER_COUNT=50  # Opcional (padrão: 100)
```

### Executando o Producer

```bash
# Envia 1000 mensagens (10 por sessão × 100 sessões)
go run cmd/sender/main.go
```

### Executando o Consumer

```bash
# Inicia o pool de workers
go run cmd/consumer/main.go
```

## 🔧 Funcionalidades Técnicas

### Gerenciamento de Concorrência

```go
// Semáforo para controlar workers simultâneos
sem := make(chan struct{}, envs.WorkerCount)

// Padrão de aquisição/liberação segura
sem <- struct{}{}        // Adquire
defer func() { <-sem }() // Libera
```

### Session Management

- **Session Affinity**: Mensagens da mesma sessão sempre processadas em sequência
- **Auto-discovery**: Workers descobrem sessões disponíveis automaticamente
- **Idle Detection**: Sessões inativas são liberadas após 3 tentativas vazias

### Graceful Shutdown

```go
// Captura sinais do sistema
signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

// Cancela contexto para propagação limpa
ctx, cancel := context.WithCancel(context.Background())
```

## 📈 Características de Performance

| Métrica | Valor |
|---------|-------|
| **Workers Simultâneos** | Configurável (padrão: 100) |
| **Timeout por Receção** | 20 segundos |
| **Retry em Timeout** | 5 segundos |
| **Idle Threshold** | 3 tentativas vazias |

## 🏢 Casos de Uso Empresariais

### E-commerce
- Processamento de pedidos por cliente (sessão = customer_id)
- Garantia de ordem de operações por usuário

### Sistemas Financeiros  
- Transações bancárias sequenciais por conta
- Auditoria e compliance garantidos

### IoT & Telemetria
- Dados de sensores agrupados por device_id
- Processamento temporal ordenado

## 🛠️ Stack Tecnológica

- **Linguagem**: Go 1.24.4
- **Cloud**: Azure Service Bus
- **SDK**: azure-sdk-for-go v1.9.0
- **Padrões**: Worker Pool, Producer-Consumer, Graceful Shutdown

## 📊 Estrutura do Projeto

```
servicebus-session-worker/
├── cmd/
│   ├── consumer/     # Worker consumer app
│   └── sender/       # Message producer app
├── internal/
│   └── envs/         # Environment configuration
├── go.mod
└── README.md
```

## 🎯 Demonstra Competências

- **Programação Concorrente**: Goroutines, channels, semáforos
- **Cloud Native**: Integração Azure, configuração por env vars
- **Reliability**: Error handling, timeouts, graceful shutdown
- **Observabilidade**: Logging estruturado, métricas implícitas
- **Clean Architecture**: Separação de concerns, modularidade

## 🚀 Possíveis Evoluções

- [ ] Métricas com Prometheus
- [ ] Health checks endpoints  
- [ ] Circuit breaker pattern
- [ ] Dead letter queue handling
- [ ] Distributed tracing
- [ ] Auto-scaling baseado em queue depth

---

<div align="center">

**Desenvolvido com ❤️ usando Go e Azure Service Bus**

*Demonstrando expertise em sistemas distribuídos e messaging patterns*

</div>