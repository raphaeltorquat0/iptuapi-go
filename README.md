# IPTU API - Go SDK

SDK oficial para integração com a IPTU API.

## Instalação

```bash
go get github.com/raphaeltorquat0/iptuapi-go
```

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
    resultado, err := client.ConsultaEndereco("Avenida Paulista", "1000")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("SQL: %s, Bairro: %s\n", resultado.SQL, resultado.Bairro)

    // Consulta por SQL (Starter+)
    dados, err := client.ConsultaSQL("100-01-001-001")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Valor Venal: R$ %.2f\n", dados.ValorVenal)

    // Avaliação de mercado (Pro+)
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
}
```

## Tratamento de Erros

```go
resultado, err := client.ConsultaEndereco("Rua Inexistente", "")
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
// Com URL customizada
client := iptuapi.NewClient("sua_api_key",
    iptuapi.WithBaseURL("https://custom.api.com"),
    iptuapi.WithTimeout(60 * time.Second),
)
```

## Documentação

Acesse a documentação completa em [iptuapi.com.br/docs](https://iptuapi.com.br/docs)
