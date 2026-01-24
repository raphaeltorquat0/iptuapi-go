// Exemplo avançado com timeout customizado e tratamento de erros
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	iptuapi "github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
	apiKey := os.Getenv("IPTU_API_KEY")
	if apiKey == "" {
		log.Fatal("IPTU_API_KEY environment variable is required")
	}

	// Cliente com timeout customizado
	client := iptuapi.NewClient(apiKey,
		iptuapi.WithTimeout(60*time.Second),
	)

	ctx := context.Background()

	// Consulta com tratamento de erros
	fmt.Println("=== Consulta com Tratamento de Erros ===")
	resultado, err := client.ConsultaEndereco(ctx, &iptuapi.ConsultaEnderecoParams{
		Logradouro: "Avenida Paulista",
		Numero:     "1000",
		Cidade:     iptuapi.CidadeSaoPaulo,
	})
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("SQL: %s\n", resultado.SQL)
	fmt.Printf("Valor Venal Total: R$ %.2f\n", resultado.ValorVenalTotal)

	// Consulta por SQL
	if resultado.SQL != "" {
		fmt.Println("\n=== Consulta por SQL ===")
		sqlResult, err := client.ConsultaSQL(ctx, resultado.SQL, iptuapi.CidadeSaoPaulo)
		if err != nil {
			handleError(err)
			return
		}

		fmt.Printf("SQL: %s\n", sqlResult.SQL)
		fmt.Printf("Ano: %d\n", sqlResult.Ano)
		fmt.Printf("Logradouro: %s, %s\n", sqlResult.Logradouro, sqlResult.Numero)
		fmt.Printf("Bairro: %s\n", sqlResult.Bairro)
		fmt.Printf("Área Terreno: %.2f m²\n", sqlResult.AreaTerreno)
		fmt.Printf("Área Construída: %.2f m²\n", sqlResult.AreaConstruida)
		fmt.Printf("Valor Venal Total: R$ %.2f\n", sqlResult.ValorVenalTotal)
		fmt.Printf("IPTU Valor: R$ %.2f\n", sqlResult.IPTUValor)
	}
}

func handleError(err error) {
	if iptuapi.IsAuthError(err) {
		fmt.Println("Erro: API Key inválida ou expirada")
	} else if iptuapi.IsNotFound(err) {
		fmt.Println("Erro: Imóvel não encontrado")
	} else if iptuapi.IsRateLimit(err) {
		fmt.Println("Erro: Rate limit excedido. Aguarde antes de tentar novamente.")
	} else if apiErr, ok := err.(*iptuapi.APIError); ok {
		fmt.Printf("Erro da API (status %d): %s\n", apiErr.StatusCode, apiErr.Message)
	} else {
		fmt.Printf("Erro: %v\n", err)
	}
}
