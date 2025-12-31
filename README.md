# IPTU API - Go SDK

SDK oficial para integração com a IPTU API - Dados de IPTU de São Paulo e Belo Horizonte.

## Instalação

```bash
go get github.com/raphaeltorquat0/iptuapi-go
```

## Cidades Suportadas

| Cidade | Constante | Identificador |
|--------|-----------|---------------|
| São Paulo | `iptuapi.CidadeSaoPaulo` | Número SQL |
| Belo Horizonte | `iptuapi.CidadeBeloHorizonte` | Índice Cadastral |

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

    // Consulta por endereço (São Paulo - endpoint legado)
    resultado, err := client.ConsultaEndereco("Avenida Paulista", "1000")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("SQL: %s, Bairro: %s\n", resultado.Data.SQLBase, resultado.Data.Bairro)

    // Consulta por SQL (Starter+)
    dados, err := client.ConsultaSQL("100-01-001-001")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Valor Venal: R$ %.2f\n", dados.ValorVenal)
}
```

## Consulta Multi-Cidade (Novo!)

```go
package main

import (
    "fmt"
    "log"

    "github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
    client := iptuapi.NewClient("sua_api_key")

    // São Paulo - busca por endereço
    resultadosSP, err := client.ConsultaIPTU(iptuapi.CidadeSaoPaulo, "Avenida Paulista", &iptuapi.ConsultaIPTUOptions{
        Numero: intPtr(1000),
        Ano:    2024,
        Limit:  20,
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, imovel := range resultadosSP {
        fmt.Printf("SQL: %s, Valor Venal: R$ %.2f\n", imovel.SQL, imovel.ValorVenal)
    }

    // Belo Horizonte - busca por endereço
    resultadosBH, err := client.ConsultaIPTU(iptuapi.CidadeBeloHorizonte, "Afonso Pena", nil)
    if err != nil {
        log.Fatal(err)
    }
    for _, imovel := range resultadosBH {
        fmt.Printf("Índice: %s, Valor Venal: R$ %.2f\n", imovel.SQL, imovel.ValorVenal)
    }

    // Busca por identificador único
    // São Paulo (SQL)
    dadosSP, err := client.ConsultaIPTUSQL(iptuapi.CidadeSaoPaulo, "00904801381", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Encontrados: %d registros\n", len(dadosSP))

    // Belo Horizonte (Índice Cadastral)
    dadosBH, err := client.ConsultaIPTUSQL(iptuapi.CidadeBeloHorizonte, "007028 005 0086", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Encontrados: %d registros\n", len(dadosBH))
}

func intPtr(i int) *int {
    return &i
}
```

## Avaliação de Mercado (Pro+)

```go
avaliacao, err := client.ValuationEstimate(iptuapi.ValuationParams{
    AreaTerreno:    250,
    AreaConstruida: 180,
    Bairro:         "Pinheiros",
    Zona:           "ZM",
    TipoUso:        "Residencial",
    TipoPadrao:     "Médio",
    AnoConstrucao:  2010,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Valor Estimado: R$ %.2f\n", avaliacao.ValorEstimado)
```

## Tratamento de Erros

```go
resultado, err := client.ConsultaIPTU(iptuapi.CidadeSaoPaulo, "Rua Inexistente", nil)
if err != nil {
    if iptuapi.IsNotFound(err) {
        fmt.Println("Imóvel não encontrado")
    } else if iptuapi.IsRateLimit(err) {
        fmt.Println("Limite de requisições excedido")
    } else if iptuapi.IsAuthError(err) {
        fmt.Println("API Key inválida")
    } else {
        log.Fatal(err)
    }
}
```

## Opções do Cliente

```go
import "time"

// Com URL customizada e timeout
client := iptuapi.NewClient("sua_api_key",
    iptuapi.WithBaseURL("https://custom.api.com"),
    iptuapi.WithTimeout(60 * time.Second),
)
```

## Documentação

Acesse a documentação completa em [iptuapi.com.br/docs](https://iptuapi.com.br/docs)
