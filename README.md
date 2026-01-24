# IPTU API - Go SDK

SDK oficial Go para integracao com a IPTU API. Acesso a dados de IPTU de Sao Paulo, Belo Horizonte e Recife.

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/iptuapi/iptuapi-go.svg)](https://pkg.go.dev/github.com/iptuapi/iptuapi-go)
[![License](https://img.shields.io/badge/license-Proprietary-red)](LICENSE)

## Instalacao

```bash
go get github.com/iptuapi/iptuapi-go
```

## Uso Rapido

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/iptuapi/iptuapi-go"
)

func main() {
    client := iptuapi.NewClient("sua_api_key")

    ctx := context.Background()

    // Consulta por endereco
    resultado, err := client.ConsultaEndereco(ctx, "Avenida Paulista", "1000", "sp")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("SQL: %s, Bairro: %s\n", resultado.SQL, resultado.Bairro)
}
```

## Configuracao

### Cliente Basico

```go
client := iptuapi.NewClient("sua_api_key")
```

### Configuracao Avancada

```go
import (
    "log/slog"
    "time"

    "github.com/iptuapi/iptuapi-go"
)

// Configuracao de retry
retryConfig := iptuapi.RetryConfig{
    MaxRetries:       5,
    InitialDelay:     time.Second,
    MaxDelay:         30 * time.Second,
    BackoffFactor:    2.0,
    RetryableStatus:  []int{429, 500, 502, 503, 504},
}

// Configuracao do cliente
config := iptuapi.ClientConfig{
    BaseURL:     "https://iptuapi.com.br/api/v1",
    Timeout:     60 * time.Second,
    RetryConfig: retryConfig,
    Logger:      slog.Default(), // Logger compativel com slog
}

client := iptuapi.NewClientWithConfig("sua_api_key", config)
```

### Logging Customizado

```go
import "log/slog"

// Usar slog logger customizado
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client := iptuapi.NewClient("sua_api_key",
    iptuapi.WithLogger(logger),
)
```

## Endpoints da API

### Consultas (Todos os Planos)

```go
ctx := context.Background()

// Consulta por endereco
resultado, err := client.ConsultaEndereco(ctx, "Avenida Paulista", "1000", "sp")

// Consulta por CEP
resultado, err := client.ConsultaCEP(ctx, "01310-100", "sp")

// Consulta por coordenadas (zoneamento)
resultado, err := client.ConsultaZoneamento(ctx, -23.5505, -46.6333)
```

### Consultas Avancadas (Starter+)

```go
// Consulta por numero SQL
resultado, err := client.ConsultaSQL(ctx, "100-01-001-001", "sp")

// Historico de valores IPTU
historico, err := client.DadosIPTUHistorico(ctx, "100-01-001-001", "sp")

// Consulta CNPJ
empresa, err := client.DadosCNPJ(ctx, "12345678000100")

// Correcao monetaria IPCA
corrigido, err := client.DadosIPCACorrigir(ctx, 100000.0, "2020-01", "2024-01")
```

### Valuation (Pro+)

```go
// Estimativa de valor de mercado
params := iptuapi.ValuationParams{
    AreaTerreno:    250,
    AreaConstruida: 180,
    Bairro:         "Pinheiros",
    Zona:           "ZM",
    TipoUso:        "Residencial",
    TipoPadrao:     "Medio",
    AnoConstrucao:  2010,
}
avaliacao, err := client.ValuationEstimate(ctx, params)
fmt.Printf("Valor estimado: R$ %.2f\n", avaliacao.ValorEstimado)

// Buscar comparaveis
comparaveis, err := client.ValuationComparables(ctx, iptuapi.ComparablesParams{
    Bairro:  "Pinheiros",
    AreaMin: 150,
    AreaMax: 250,
    Cidade:  "sp",
    Limit:   10,
})
```

### Batch Operations (Enterprise)

```go
// Valuation em lote (ate 100 imoveis)
imoveis := []iptuapi.ValuationParams{
    {AreaTerreno: 250, AreaConstruida: 180, Bairro: "Pinheiros"},
    {AreaTerreno: 300, AreaConstruida: 200, Bairro: "Moema"},
}
resultados, err := client.ValuationBatch(ctx, imoveis)
```

## Context e Cancelamento

```go
// Com timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resultado, err := client.ConsultaEndereco(ctx, "Avenida Paulista", "1000", "sp")
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Requisicao cancelada por timeout")
    }
}

// Com cancelamento manual
ctx, cancel := context.WithCancel(context.Background())

go func() {
    time.Sleep(5 * time.Second)
    cancel() // Cancela a requisicao
}()

resultado, err := client.ConsultaEndereco(ctx, "Avenida Paulista", "1000", "sp")
```

## Tratamento de Erros

```go
import "github.com/iptuapi/iptuapi-go"

resultado, err := client.ConsultaEndereco(ctx, "Rua Teste", "100", "sp")
if err != nil {
    var apiErr *iptuapi.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("Status: %d, Request ID: %s\n", apiErr.StatusCode, apiErr.RequestID)
        fmt.Printf("Retryable: %v\n", apiErr.IsRetryable())
    }

    // Verificar tipos especificos de erro
    if iptuapi.IsAuthError(err) {
        fmt.Println("API Key invalida")
    } else if iptuapi.IsForbidden(err) {
        var forbiddenErr *iptuapi.ForbiddenError
        if errors.As(err, &forbiddenErr) {
            fmt.Printf("Plano requerido: %s\n", forbiddenErr.RequiredPlan)
        }
    } else if iptuapi.IsNotFound(err) {
        fmt.Println("Imovel nao encontrado")
    } else if iptuapi.IsRateLimit(err) {
        var rateLimitErr *iptuapi.RateLimitError
        if errors.As(err, &rateLimitErr) {
            fmt.Printf("Retry em: %d segundos\n", rateLimitErr.RetryAfter)
        }
    } else if iptuapi.IsValidation(err) {
        var validationErr *iptuapi.ValidationError
        if errors.As(err, &validationErr) {
            for field, msgs := range validationErr.Errors {
                fmt.Printf("Campo %s: %v\n", field, msgs)
            }
        }
    } else if iptuapi.IsServerError(err) {
        fmt.Println("Erro no servidor (retryable)")
    } else if iptuapi.IsTimeout(err) {
        fmt.Println("Timeout na requisicao")
    } else if iptuapi.IsNetworkError(err) {
        fmt.Println("Erro de conexao")
    }
}
```

### Funcoes de Verificacao de Erro

```go
// Verificar se erro e de um tipo especifico
if iptuapi.IsAuthError(err) { ... }
if iptuapi.IsForbidden(err) { ... }
if iptuapi.IsNotFound(err) { ... }
if iptuapi.IsRateLimit(err) { ... }
if iptuapi.IsValidation(err) { ... }
if iptuapi.IsServerError(err) { ... }
if iptuapi.IsTimeout(err) { ... }
if iptuapi.IsNetworkError(err) { ... }
```

## Rate Limiting

```go
// Verificar rate limit apos requisicao
if rateLimit := client.GetRateLimitInfo(); rateLimit != nil {
    fmt.Printf("Limite: %d\n", rateLimit.Limit)
    fmt.Printf("Restantes: %d\n", rateLimit.Remaining)
    fmt.Printf("Reset em: %s\n", rateLimit.ResetAt().Format(time.RFC3339))
}

// ID da ultima requisicao (util para suporte)
fmt.Printf("Request ID: %s\n", client.GetLastRequestID())
```

## Tipos e Structs

```go
// Resultado de consulta de endereco
type PropertyData struct {
    SQL             string   `json:"sql"`
    Logradouro      string   `json:"logradouro"`
    Numero          string   `json:"numero"`
    Bairro          string   `json:"bairro"`
    Cidade          string   `json:"cidade"`
    CEP             string   `json:"cep"`
    AreaTerreno     float64  `json:"area_terreno"`
    AreaConstruida  float64  `json:"area_construida"`
    ValorVenal      float64  `json:"valor_venal"`
    ValorIPTU       float64  `json:"valor_iptu"`
    Latitude        float64  `json:"latitude"`
    Longitude       float64  `json:"longitude"`
}

// Parametros de valuation
type ValuationParams struct {
    AreaTerreno    float64 `json:"area_terreno"`
    AreaConstruida float64 `json:"area_construida"`
    Bairro         string  `json:"bairro"`
    Zona           string  `json:"zona,omitempty"`
    TipoUso        string  `json:"tipo_uso,omitempty"`
    TipoPadrao     string  `json:"tipo_padrao,omitempty"`
    AnoConstrucao  int     `json:"ano_construcao,omitempty"`
}

// Resultado de valuation
type ValuationResult struct {
    ValorEstimado   float64 `json:"valor_estimado"`
    ValorMinimo     float64 `json:"valor_minimo"`
    ValorMaximo     float64 `json:"valor_maximo"`
    Confianca       float64 `json:"confianca"`
    Comparables     int     `json:"comparables_count"`
}

// Rate limit info
type RateLimitInfo struct {
    Limit     int
    Remaining int
    Reset     int64
}

func (r *RateLimitInfo) ResetAt() time.Time {
    return time.Unix(r.Reset, 0)
}
```

## Testes

```bash
# Rodar testes
go test ./...

# Com verbose
go test -v ./...

# Com coverage
go test -cover ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Cidades Suportadas

| Codigo | Cidade |
|--------|--------|
| sp | Sao Paulo |
| bh | Belo Horizonte |
| recife | Recife |

## Licenca

Copyright (c) 2025-2026 IPTU API. Todos os direitos reservados.

Este software e propriedade exclusiva da IPTU API. O uso esta sujeito aos termos de servico disponiveis em https://iptuapi.com.br/termos

## Links

- [Documentacao](https://iptuapi.com.br/docs)
- [API Reference](https://iptuapi.com.br/docs/api)
- [Portal do Desenvolvedor](https://iptuapi.com.br/dashboard)
