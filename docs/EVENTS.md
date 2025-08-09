# Sistema de Eventos - Master Construtora

## Visão Geral

O Master Construtora implementa um sistema de eventos interno baseado no padrão **Event-Driven Architecture**. Este sistema permite a comunicação assíncrona e desacoplada entre os diferentes módulos (bounded contexts) da aplicação.

## Arquitetura de Eventos

### Componentes Principais

#### 1. EventBus (`internal/platform/bus/eventbus.go`)
- **Publisher/Subscriber pattern**
- **Comunicação in-memory**
- **Execução assíncrona via goroutines**
- **Thread-safe com sync.RWMutex**

#### 2. Event Definitions (`internal/events/`)
- **Constantes para nomes de eventos**
- **Structs para payloads de eventos**
- **Tipagem forte para garantir consistência**

#### 3. Event Handlers
- **Handlers específicos por módulo**
- **Implementação de lógica de negócio reativa**
- **Error handling e logging**

## Implementação do EventBus

### Interface Principal

```go
type EventBus interface {
    Publicar(eventoNome string, dados interface{})
    Subscrever(eventoNome string, handler func(dados interface{}))
}
```

### Estrutura Interna

```go
type eventBus struct {
    handlers map[string][]func(interface{})
    mu       sync.RWMutex
    logger   *slog.Logger
}
```

### Métodos Principais

#### Subscrever (Subscribe)
Registra um handler para um evento específico:

```go
func (eb *eventBus) Subscrever(eventoNome string, handler func(dados interface{})) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    if eb.handlers[eventoNome] == nil {
        eb.handlers[eventoNome] = make([]func(interface{}), 0)
    }
    
    eb.handlers[eventoNome] = append(eb.handlers[eventoNome], handler)
    eb.logger.Info("handler subscrito ao evento", "evento", eventoNome)
}
```

#### Publicar (Publish)
Publica um evento para todos os handlers subscritos:

```go
func (eb *eventBus) Publicar(eventoNome string, dados interface{}) {
    eb.mu.RLock()
    handlers := eb.handlers[eventoNome]
    eb.mu.RUnlock()
    
    if len(handlers) == 0 {
        eb.logger.Warn("nenhum handler encontrado para evento", "evento", eventoNome)
        return
    }
    
    eb.logger.Info("publicando evento", "evento", eventoNome, "handlers", len(handlers))
    
    // Executa handlers em paralelo
    for _, handler := range handlers {
        go func(h func(interface{})) {
            defer func() {
                if r := recover(); r != nil {
                    eb.logger.Error("panic no handler de evento", 
                        "evento", eventoNome, 
                        "panic", r)
                }
            }()
            h(dados)
        }(handler)
    }
}
```

## Eventos Implementados

### 1. OrcamentoStatusAtualizado

**Evento**: `orcamento:status_atualizado`  
**Módulo Origem**: Suprimentos  
**Módulos Destino**: Obras

#### Payload

```go
type OrcamentoStatusAtualizadoPayload struct {
    OrcamentoID string
    EtapaID     string
    NovoStatus  string
    Valor       float64
}
```

#### Quando é Disparado
- Orçamento é aprovado
- Orçamento é rejeitado
- Orçamento é marcado como pago
- Status do orçamento é alterado

#### Exemplo de Uso

```go
// No serviço de Suprimentos
func (s *Service) AtualizarStatusOrcamento(ctx context.Context, orcamentoID string, input dto.AtualizarStatusOrcamentoInput) error {
    // 1. Atualiza o orçamento no banco
    err := s.orcamentoRepo.AtualizarStatus(ctx, orcamentoID, input.Status)
    if err != nil {
        return err
    }
    
    // 2. Busca dados do orçamento
    orcamento, err := s.orcamentoRepo.BuscarPorID(ctx, orcamentoID)
    if err != nil {
        return err
    }
    
    // 3. Publica evento
    payload := events.OrcamentoStatusAtualizadoPayload{
        OrcamentoID: orcamento.ID,
        EtapaID:     orcamento.EtapaID,
        NovoStatus:  input.Status,
        Valor:       orcamento.ValorTotal,
    }
    
    s.eventBus.Publicar(events.OrcamentoStatusAtualizado, payload)
    
    return nil
}
```

### 2. ApontamentoPago (Planejado)

**Evento**: `apontamento:pago`  
**Módulo Origem**: Pessoal  
**Módulos Destino**: Financeiro

#### Payload

```go
type ApontamentoPagoPayload struct {
    ApontamentoID   string
    FuncionarioID   string
    ObraID          string
    ValorPago       float64
    PeriodoInicio   time.Time
    PeriodoFim      time.Time
    ContaBancariaID string
}
```

#### Quando é Disparado
- Apontamento é marcado como "Pago"
- Pagamento de apontamento é processado

## Event Handlers

### ObrasEventHandler

Localizado em: `internal/service/obras/events/handler.go`

#### HandleOrcamentoStatusAtualizado

```go
type ObrasEventHandler struct {
    logger *slog.Logger
}

func (h *ObrasEventHandler) HandleOrcamentoStatusAtualizado(dados interface{}) {
    payload, ok := dados.(events.OrcamentoStatusAtualizadoPayload)
    if !ok {
        h.logger.Error("tipo de payload inválido para evento OrcamentoStatusAtualizado")
        return
    }
    
    h.logger.Info("processando atualização de orçamento",
        "orcamentoId", payload.OrcamentoID,
        "etapaId", payload.EtapaID,
        "novoStatus", payload.NovoStatus,
        "valor", payload.Valor)
    
    // Lógica específica baseada no status
    switch payload.NovoStatus {
    case "Aprovado":
        h.logger.Info("orçamento aprovado - atualizando métricas financeiras da obra")
        // TODO: Atualizar métricas de custo da obra
        
    case "Rejeitado":
        h.logger.Info("orçamento rejeitado - sem ação necessária")
        
    case "Pago":
        h.logger.Info("orçamento pago - atualizando custos realizados")
        // TODO: Atualizar custos realizados da obra
    }
}
```

### Configuração no main.go

```go
// Criação do event handler
obrasEventHandler := obras_events.NovoObrasEventHandler(logger)

// Subscrição ao evento
eventBus.Subscrever(events.OrcamentoStatusAtualizado, obrasEventHandler.HandleOrcamentoStatusAtualizado)
```

## Padrões de Uso

### 1. Publicação de Eventos

#### No Service Layer
```go
func (s *Service) ExecutarOperacao(ctx context.Context, input InputDTO) error {
    // 1. Executar operação principal
    resultado, err := s.repository.ExecutarOperacao(ctx, input)
    if err != nil {
        return err
    }
    
    // 2. Publicar evento apenas se operação foi bem-sucedida
    payload := EventPayload{
        ID:    resultado.ID,
        Dados: resultado.DadosRelevantes,
    }
    
    s.eventBus.Publicar("operacao:executada", payload)
    
    return nil
}
```

#### Tratamento de Erros
```go
func (s *Service) OperacaoComEventos(ctx context.Context) error {
    // Usar transação para garantir consistência
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback() // Rollback se não commitado
    
    // Executar operações no banco
    err = s.repository.ExecutarComTransacao(ctx, tx, dados)
    if err != nil {
        return err
    }
    
    // Commit primeiro
    err = tx.Commit()
    if err != nil {
        return err
    }
    
    // Publicar eventos apenas após commit bem-sucedido
    s.eventBus.Publicar("evento:executado", payload)
    
    return nil
}
```

### 2. Implementação de Handlers

#### Estrutura Padrão
```go
type ModuloEventHandler struct {
    service ServiceInterface
    logger  *slog.Logger
}

func NovoModuloEventHandler(service ServiceInterface, logger *slog.Logger) *ModuloEventHandler {
    return &ModuloEventHandler{
        service: service,
        logger:  logger.With("component", "ModuloEventHandler"),
    }
}

func (h *ModuloEventHandler) HandleEvento(dados interface{}) {
    // 1. Validar tipo do payload
    payload, ok := dados.(ExpectedPayloadType)
    if !ok {
        h.logger.Error("tipo de payload inválido", "expected", "ExpectedPayloadType")
        return
    }
    
    // 2. Log do início do processamento
    h.logger.Info("processando evento", "id", payload.ID)
    
    // 3. Executar lógica de negócio
    err := h.service.ProcessarEvento(context.Background(), payload)
    if err != nil {
        h.logger.Error("erro ao processar evento", "error", err, "id", payload.ID)
        return
    }
    
    // 4. Log de sucesso
    h.logger.Info("evento processado com sucesso", "id", payload.ID)
}
```

#### Error Handling
```go
func (h *Handler) HandleEventoComRecovery(dados interface{}) {
    defer func() {
        if r := recover(); r != nil {
            h.logger.Error("panic no handler de evento", 
                "panic", r,
                "evento", "nome-do-evento")
        }
    }()
    
    // Processamento normal do evento
    h.processarEvento(dados)
}
```

## Casos de Uso Implementados

### 1. Atualização de Métricas Financeiras

**Fluxo:**
1. Orçamento é aprovado no módulo Suprimentos
2. Evento `OrcamentoStatusAtualizado` é publicado
3. Módulo Obras recebe o evento
4. Métricas financeiras da obra são atualizadas
5. Dashboard reflete novos valores automaticamente

### 2. Sincronização de Status

**Fluxo:**
1. Status de entidade muda em um módulo
2. Evento é publicado com novos dados
3. Módulos interessados atualizam suas views/cache
4. Consistência eventual é mantida

## Benefícios do Sistema de Eventos

### 1. Desacoplamento
- Módulos não precisam conhecer uns aos outros diretamente
- Redução de dependências circulares
- Facilita testes unitários

### 2. Extensibilidade
- Novos handlers podem ser adicionados facilmente
- Funcionalidades podem ser adicionadas sem modificar código existente
- Suporte a múltiplos consumidores por evento

### 3. Escalabilidade
- Processamento assíncrono não bloqueia operações principais
- Handlers podem ser otimizados independentemente
- Base para futuras implementações distribuídas

### 4. Auditoria
- Todos os eventos são logados
- Facilita debugging e troubleshooting
- Rastreabilidade de mudanças no sistema

## Limitações Atuais

### 1. In-Memory Only
- Eventos são perdidos se aplicação reiniciar
- Não há garantia de entrega
- Sem persistência de eventos

### 2. Sem Retry Logic
- Falhas nos handlers são apenas logadas
- Não há tentativas automáticas de reprocessamento
- Eventos podem ser perdidos em caso de erro

### 3. Ordem de Processamento
- Não há garantia de ordem entre eventos
- Handlers executam em paralelo
- Pode causar race conditions em casos específicos

## Evoluções Futuras

### 1. Event Store
Implementar persistência de eventos:
```go
type EventStore interface {
    Save(event Event) error
    GetEvents(aggregateID string) ([]Event, error)
    GetEventsSince(timestamp time.Time) ([]Event, error)
}
```

### 2. Message Queue
Integração com sistemas externos:
- Redis Streams
- Apache Kafka
- RabbitMQ

### 3. Retry e Dead Letter Queue
```go
type RetryPolicy struct {
    MaxRetries int
    BackoffStrategy string
    DeadLetterQueue string
}
```

### 4. Event Sourcing
Usar eventos como fonte de verdade:
```go
type Aggregate interface {
    Apply(event Event)
    GetUncommittedEvents() []Event
    MarkEventsAsCommitted()
}
```

## Debugging e Troubleshooting

### Logs de Eventos

Os logs incluem informações detalhadas:

```json
{
  "level": "info",
  "time": "2024-02-20T14:30:00Z",
  "msg": "publicando evento",
  "evento": "orcamento:status_atualizado",
  "handlers": 1,
  "component": "EventBus"
}
```

### Verificação de Handlers

Para verificar quais handlers estão registrados:

```go
func (eb *eventBus) GetHandlers() map[string]int {
    eb.mu.RLock()
    defer eb.mu.RUnlock()
    
    result := make(map[string]int)
    for eventName, handlers := range eb.handlers {
        result[eventName] = len(handlers)
    }
    return result
}
```

### Métricas Úteis

- Número de eventos publicados por tipo
- Tempo de processamento dos handlers
- Taxa de erro nos handlers
- Handlers registrados por evento

## Exemplo de Teste

```go
func TestEventBus(t *testing.T) {
    // Arrange
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    eventBus := bus.NovoEventBus(logger)
    
    var receivedData interface{}
    handler := func(data interface{}) {
        receivedData = data
    }
    
    // Act
    eventBus.Subscrever("test:event", handler)
    eventBus.Publicar("test:event", "test data")
    
    // Wait for async processing
    time.Sleep(100 * time.Millisecond)
    
    // Assert
    assert.Equal(t, "test data", receivedData)
}
```