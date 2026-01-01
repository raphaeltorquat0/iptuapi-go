# IPTU API - Go SDK

SDK oficial Go para integração com a IPTU API. Acesso a dados de IPTU e ITBI de São Paulo, Belo Horizonte e Recife.

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/raphaeltorquat0/iptuapi-go.svg)](https://pkg.go.dev/github.com/raphaeltorquat0/iptuapi-go)
[![License](https://img.shields.io/badge/license-Proprietary-red)](LICENSE)

## Instalação

```bash
go get github.com/raphaeltorquat0/iptuapi-go
```

## Requisitos

- Go 1.21+

## Uso Rápido

```go
package main

import (
    "fmt"
    "log"

    "github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
    client := iptuapi.NewClient("sua_api_key")

    // Consulta por endereço
    resultado, err := client.ConsultaEndereco("Avenida Paulista", "1000", "sp")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("SQL: %s\n", resultado.Data[0].SQL)
    fmt.Printf("Valor Venal: R$ %.2f\n", resultado.Data[0].ValorVenal)
}
```

## Configuração

### Cliente Básico

```go
client := iptuapi.NewClient("sua_api_key")
```

### Configuração Avançada

```go
import "time"

client := iptuapi.NewClient("sua_api_key",
    iptuapi.WithBaseURL("https://iptuapi.com.br/api/v1"),
    iptuapi.WithTimeout(60 * time.Second),
    iptuapi.WithRetry(iptuapi.RetryConfig{
        MaxRetries:      5,
        InitialDelay:    1 * time.Second,
        MaxDelay:        30 * time.Second,
        BackoffFactor:   2.0,
        RetryableStatus: []int{429, 500, 502, 503, 504},
    }),
)
```

---

## Endpoints de Consulta IPTU

### Consulta por Endereço

Busca dados de IPTU por logradouro e número. Disponível em **todos os planos**.

```go
// Consulta básica
resultado, err := client.ConsultaEndereco("Avenida Paulista", "1000", "sp")
if err != nil {
    log.Fatal(err)
}

for _, imovel := range resultado.Data {
    fmt.Printf("SQL: %s\n", imovel.SQL)
    fmt.Printf("Bairro: %s\n", imovel.Bairro)
    fmt.Printf("Valor Venal: R$ %.2f\n", imovel.ValorVenal)
}
```

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| logradouro | string | Sim | Nome da rua/avenida |
| numero | string | Sim | Número do imóvel |
| cidade | string | Não | Código da cidade (sp, bh, recife). Default: sp |

**Resposta:**

```json
{
  "success": true,
  "data": [
    {
      "sql": "008.045.0123-4",
      "logradouro": "AV PAULISTA",
      "numero": "1000",
      "bairro": "BELA VISTA",
      "cep": "01310-100",
      "area_terreno": 500.0,
      "area_construida": 1200.0,
      "valor_venal": 2500000.0,
      "valor_venal_terreno": 1500000.0,
      "valor_venal_construcao": 1000000.0,
      "ano_construcao": 1985,
      "uso": "Comercial",
      "padrao": "Alto"
    }
  ],
  "metadata": {
    "total": 1,
    "cidade": "sp",
    "ano_referencia": 2024
  }
}
```

---

### Consulta por SQL/Índice Cadastral

Busca por identificador único do imóvel. Disponível a partir do plano **Starter**.

```go
// São Paulo - número SQL
resultado, err := client.ConsultaSQL("008.045.0123-4", "sp")

// Belo Horizonte - índice cadastral
resultado, err := client.ConsultaSQL("007028 005 0086", "bh")

// Recife - sequencial
resultado, err := client.ConsultaSQL("123456", "recife")
```

---

### Consulta por CEP

Busca todos os imóveis de um CEP. Disponível em **todos os planos**.

```go
resultado, err := client.ConsultaCEP("01310-100", "sp")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Encontrados: %d imóveis\n", len(resultado.Data))
```

---

### Consulta Zoneamento

Retorna informações de zoneamento por coordenadas. Disponível em **todos os planos**.

```go
resultado, err := client.ConsultaZoneamento(-23.5505, -46.6333)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Zona: %s\n", resultado.Data.Zona)
fmt.Printf("Uso Permitido: %s\n", resultado.Data.UsoPermitido)
```

---

## Endpoints de Valuation

### Estimativa de Valor de Mercado

Calcula valor estimado de mercado. Disponível a partir do plano **Pro**.

```go
avaliacao, err := client.ValuationEstimate(iptuapi.ValuationParams{
    AreaTerreno:    250.0,
    AreaConstruida: 180.0,
    Bairro:         "Pinheiros",
    Cidade:         "sp",
    Zona:           "ZM",
    TipoUso:        "Residencial",
    TipoPadrao:     "Medio",
    AnoConstrucao:  2010,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Valor Estimado: R$ %.2f\n", avaliacao.ValorEstimado)
fmt.Printf("Confiança: %.1f%%\n", avaliacao.Confianca*100)
fmt.Printf("Faixa: R$ %.2f - R$ %.2f\n", avaliacao.ValorMinimo, avaliacao.ValorMaximo)
```

**Resposta:**

```json
{
  "success": true,
  "data": {
    "valor_estimado": 1250000.0,
    "valor_minimo": 1100000.0,
    "valor_maximo": 1400000.0,
    "confianca": 0.85,
    "valor_m2": 6944.44,
    "metodologia": "regressao_comparativa",
    "data_referencia": "2024-01-15"
  }
}
```

---

### Buscar Comparáveis

Retorna imóveis similares para comparação. Disponível a partir do plano **Pro**.

```go
comparaveis, err := client.ValuationComparables(iptuapi.ComparablesParams{
    Bairro:   "Pinheiros",
    AreaMin:  150.0,
    AreaMax:  250.0,
    Cidade:   "sp",
    Limit:    10,
})
if err != nil {
    log.Fatal(err)
}

for _, comp := range comparaveis.Data {
    fmt.Printf("SQL: %s, Área: %.0fm², Valor: R$ %.2f\n",
        comp.SQL, comp.AreaConstruida, comp.ValorVenal)
}
```

---

## Endpoints de ITBI

### Status da Transação ITBI

Consulta status de uma transação ITBI. Disponível em **todos os planos**.

```go
status, err := client.ITBIStatus("ITBI-2024-123456", "sp")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Data.Status)
fmt.Printf("Valor ITBI: R$ %.2f\n", status.Data.ValorITBI)
```

**Resposta:**

```json
{
  "success": true,
  "data": {
    "protocolo": "ITBI-2024-123456",
    "status": "aprovado",
    "data_solicitacao": "2024-01-10",
    "data_aprovacao": "2024-01-12",
    "valor_transacao": 500000.0,
    "valor_venal_referencia": 480000.0,
    "base_calculo": 500000.0,
    "aliquota": 0.03,
    "valor_itbi": 15000.0
  }
}
```

---

### Cálculo de ITBI

Calcula valor do ITBI para uma transação. Disponível em **todos os planos**.

```go
calculo, err := client.ITBICalcular(iptuapi.ITBICalculoParams{
    SQL:            "008.045.0123-4",
    ValorTransacao: 500000.0,
    Cidade:         "sp",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Valor ITBI: R$ %.2f\n", calculo.Data.ValorITBI)
fmt.Printf("Alíquota: %.1f%%\n", calculo.Data.Aliquota*100)
```

**Resposta:**

```json
{
  "success": true,
  "data": {
    "sql": "008.045.0123-4",
    "valor_transacao": 500000.0,
    "valor_venal_referencia": 480000.0,
    "base_calculo": 500000.0,
    "aliquota": 0.03,
    "valor_itbi": 15000.0,
    "isencao_aplicavel": false,
    "fundamentacao_legal": "Lei Municipal 11.154/1991"
  }
}
```

---

### Histórico de Transações ITBI

Retorna histórico de transações de um imóvel. Disponível a partir do plano **Starter**.

```go
historico, err := client.ITBIHistorico("008.045.0123-4", "sp")
if err != nil {
    log.Fatal(err)
}

for _, tx := range historico.Data {
    fmt.Printf("%s - R$ %.2f (%s)\n",
        tx.DataTransacao, tx.ValorTransacao, tx.TipoTransacao)
}
```

**Resposta:**

```json
{
  "success": true,
  "data": [
    {
      "protocolo": "ITBI-2024-123456",
      "data_transacao": "2024-01-15",
      "tipo_transacao": "compra_venda",
      "valor_transacao": 500000.0,
      "valor_itbi": 15000.0
    },
    {
      "protocolo": "ITBI-2020-098765",
      "data_transacao": "2020-06-20",
      "tipo_transacao": "compra_venda",
      "valor_transacao": 380000.0,
      "valor_itbi": 11400.0
    }
  ],
  "metadata": {
    "total_transacoes": 2,
    "sql": "008.045.0123-4"
  }
}
```

---

### Alíquotas ITBI

Retorna alíquotas vigentes por cidade. Disponível em **todos os planos**.

```go
aliquotas, err := client.ITBIAliquotas("sp")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Alíquota padrão: %.1f%%\n", aliquotas.Data.AliquotaPadrao*100)
```

**Resposta:**

```json
{
  "success": true,
  "data": {
    "cidade": "sp",
    "aliquota_padrao": 0.03,
    "aliquota_financiamento_sfh": 0.005,
    "valor_minimo_isencao": 0,
    "base_legal": "Lei Municipal 11.154/1991",
    "vigencia": "2024-01-01"
  }
}
```

---

### Isenções ITBI

Verifica isenções aplicáveis. Disponível em **todos os planos**.

```go
isencoes, err := client.ITBIIsencoes("sp")
if err != nil {
    log.Fatal(err)
}

for _, isencao := range isencoes.Data {
    fmt.Printf("- %s: %s\n", isencao.Tipo, isencao.Descricao)
}
```

---

### Guia ITBI

Gera guia de pagamento do ITBI. Disponível a partir do plano **Starter**.

```go
guia, err := client.ITBIGuia(iptuapi.ITBIGuiaParams{
    SQL:            "008.045.0123-4",
    ValorTransacao: 500000.0,
    Comprador: iptuapi.Pessoa{
        Nome:      "João da Silva",
        Documento: "123.456.789-00",
        Email:     "joao@email.com",
    },
    Vendedor: iptuapi.Pessoa{
        Nome:      "Maria Santos",
        Documento: "987.654.321-00",
    },
    Cidade: "sp",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Protocolo: %s\n", guia.Data.Protocolo)
fmt.Printf("Código de barras: %s\n", guia.Data.CodigoBarras)
fmt.Printf("Vencimento: %s\n", guia.Data.DataVencimento)
```

**Resposta:**

```json
{
  "success": true,
  "data": {
    "protocolo": "ITBI-2024-789012",
    "codigo_barras": "23793.38128 60000.000003 00000.000400 1 84340000015000",
    "linha_digitavel": "23793381286000000000300000000400184340000015000",
    "data_emissao": "2024-01-15",
    "data_vencimento": "2024-02-14",
    "valor_itbi": 15000.0
  }
}
```

---

### Validar Guia ITBI

Valida autenticidade de uma guia. Disponível em **todos os planos**.

```go
validacao, err := client.ITBIValidarGuia("ITBI-2024-789012", "sp")
if err != nil {
    log.Fatal(err)
}

if validacao.Data.Valido {
    fmt.Println("Guia válida!")
    if validacao.Data.Pago {
        fmt.Printf("Pago em: %s\n", validacao.Data.DataPagamento)
    }
}
```

---

### Simular ITBI

Simula cálculo sem gerar guia. Disponível em **todos os planos**.

```go
simulacao, err := client.ITBISimular(iptuapi.ITBISimularParams{
    ValorTransacao:    500000.0,
    Cidade:            "sp",
    TipoFinanciamento: "sfh",
    ValorFinanciado:   400000.0,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Valor ITBI Total: R$ %.2f\n", simulacao.Data.ValorITBITotal)
fmt.Printf("  - Parte financiada (SFH): R$ %.2f\n", simulacao.Data.ValorITBIFinanciado)
fmt.Printf("  - Parte não financiada: R$ %.2f\n", simulacao.Data.ValorITBINaoFinanciado)
```

**Resposta:**

```json
{
  "success": true,
  "data": {
    "valor_transacao": 500000.0,
    "valor_financiado": 400000.0,
    "valor_nao_financiado": 100000.0,
    "aliquota_sfh": 0.005,
    "aliquota_padrao": 0.03,
    "valor_itbi_financiado": 2000.0,
    "valor_itbi_nao_financiado": 3000.0,
    "valor_itbi_total": 5000.0,
    "economia_sfh": 10000.0
  }
}
```

---

## Tratamento de Erros

```go
import (
    "errors"
    "github.com/raphaeltorquat0/iptuapi-go"
)

resultado, err := client.ConsultaEndereco("Rua Teste", "100", "sp")
if err != nil {
    var apiErr *iptuapi.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.StatusCode {
        case 401:
            fmt.Println("API Key inválida")
        case 403:
            fmt.Printf("Plano não autorizado. Requer: %s\n", apiErr.RequiredPlan)
        case 404:
            fmt.Println("Imóvel não encontrado")
        case 429:
            fmt.Printf("Rate limit excedido. Retry em %ds\n", apiErr.RetryAfter)
        case 422:
            fmt.Printf("Parâmetros inválidos: %v\n", apiErr.ValidationErrors)
        default:
            fmt.Printf("Erro %d: %s\n", apiErr.StatusCode, apiErr.Message)
        }
        fmt.Printf("Request ID: %s\n", apiErr.RequestID)
    } else {
        // Erro de rede ou outro
        fmt.Printf("Erro: %v\n", err)
    }
}
```

### Funções Auxiliares de Erro

```go
// Verificar tipo de erro
if iptuapi.IsNotFound(err) {
    fmt.Println("Recurso não encontrado")
}

if iptuapi.IsRateLimit(err) {
    // Aguardar e tentar novamente
    time.Sleep(time.Duration(iptuapi.GetRetryAfter(err)) * time.Second)
}

if iptuapi.IsAuthError(err) {
    fmt.Println("Problema de autenticação")
}

if iptuapi.IsRetryable(err) {
    fmt.Println("Erro temporário, pode tentar novamente")
}
```

---

## Rate Limiting

```go
// Verificar rate limit após requisição
info := client.GetRateLimitInfo()
if info != nil {
    fmt.Printf("Limite: %d\n", info.Limit)
    fmt.Printf("Restantes: %d\n", info.Remaining)
    fmt.Printf("Reset em: %s\n", info.ResetAt.Format(time.RFC3339))
}

// ID da última requisição (útil para suporte)
fmt.Printf("Request ID: %s\n", client.GetLastRequestID())
```

### Limites por Plano

| Plano | Requisições/mês | Requisições/minuto |
|-------|-----------------|-------------------|
| Free | 100 | 10 |
| Starter | 5.000 | 60 |
| Pro | 50.000 | 300 |
| Enterprise | Ilimitado | 1.000 |

---

## Cidades Suportadas

| Código | Cidade | Identificador | Registros |
|--------|--------|---------------|-----------|
| sp | São Paulo | Número SQL | 4.5M+ |
| bh | Belo Horizonte | Índice Cadastral | 800K+ |
| recife | Recife | Sequencial | 400K+ |

---

## Exemplo Completo

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
    // Configuração
    client := iptuapi.NewClient(os.Getenv("IPTU_API_KEY"),
        iptuapi.WithTimeout(30*time.Second),
        iptuapi.WithRetry(iptuapi.RetryConfig{
            MaxRetries: 3,
        }),
    )

    // Lista de endereços para consultar
    enderecos := []struct {
        Logradouro string
        Numero     string
    }{
        {"Avenida Paulista", "1000"},
        {"Rua Augusta", "500"},
        {"Avenida Faria Lima", "3000"},
    }

    for _, end := range enderecos {
        resultado, err := client.ConsultaEndereco(end.Logradouro, end.Numero, "sp")
        if err != nil {
            var apiErr *iptuapi.APIError
            if errors.As(err, &apiErr) {
                if apiErr.StatusCode == 429 {
                    fmt.Printf("Rate limit. Aguardando %ds...\n", apiErr.RetryAfter)
                    time.Sleep(time.Duration(apiErr.RetryAfter) * time.Second)
                    continue
                }
            }
            fmt.Printf("Erro ao consultar %s: %v\n", end.Logradouro, err)
            continue
        }

        for _, imovel := range resultado.Data {
            fmt.Printf("SQL: %s, Valor Venal: R$ %.2f\n",
                imovel.SQL,
                imovel.ValorVenal,
            )
        }

        // Verificar rate limit
        info := client.GetRateLimitInfo()
        if info != nil && info.Remaining < 10 {
            fmt.Printf("Atenção: Apenas %d requisições restantes\n", info.Remaining)
        }
    }

    // Exemplo ITBI
    fmt.Println("\n--- Simulação ITBI ---")
    simulacao, err := client.ITBISimular(iptuapi.ITBISimularParams{
        ValorTransacao:    800000.0,
        Cidade:            "sp",
        TipoFinanciamento: "sfh",
        ValorFinanciado:   600000.0,
    })
    if err != nil {
        log.Printf("Erro na simulação ITBI: %v", err)
    } else {
        fmt.Printf("Valor ITBI: R$ %.2f\n", simulacao.Data.ValorITBITotal)
        fmt.Printf("Economia com SFH: R$ %.2f\n", simulacao.Data.EconomiaSFH)
    }
}
```

---

## Testes

```bash
# Rodar testes
go test ./...

# Com coverage
go test -cover ./...

# Verbose
go test -v ./...
```

---

## Licença

Copyright (c) 2025-2026 IPTU API. Todos os direitos reservados.

Este software é propriedade exclusiva da IPTU API. O uso está sujeito aos termos de serviço disponíveis em https://iptuapi.com.br/termos

---

## Links

- [Documentação](https://iptuapi.com.br/docs)
- [API Reference](https://iptuapi.com.br/docs/api)
- [Portal do Desenvolvedor](https://iptuapi.com.br/dashboard)
- [Suporte](mailto:suporte@iptuapi.com.br)
