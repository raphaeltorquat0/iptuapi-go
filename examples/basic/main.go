// Exemplo basico de uso do SDK
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
	// Obter API key do ambiente
	apiKey := os.Getenv("IPTU_API_KEY")
	if apiKey == "" {
		log.Fatal("IPTU_API_KEY environment variable is required")
	}

	// Criar cliente
	client := iptuapi.NewClient(apiKey)
	ctx := context.Background()

	// Consulta por endereco
	fmt.Println("=== Consulta por Endereco ===")
	resultado, err := client.ConsultaEndereco(ctx, &iptuapi.ConsultaEnderecoParams{
		Logradouro: "Avenida Paulista",
		Numero:     "1000",
		Cidade:     iptuapi.CidadeSaoPaulo,
	})
	if err != nil {
		log.Fatalf("Erro na consulta: %v", err)
	}

	fmt.Printf("SQL: %s\n", resultado.SQL)
	fmt.Printf("Logradouro: %s, %s\n", resultado.Logradouro, resultado.Numero)
	fmt.Printf("Bairro: %s\n", resultado.Bairro)
	fmt.Printf("Area Terreno: %.2f m²\n", resultado.AreaTerreno)
	fmt.Printf("Area Construida: %.2f m²\n", resultado.AreaConstruida)
	fmt.Printf("Valor Venal: R$ %.2f\n", resultado.ValorVenalTotal)

	// Verificar rate limit
	if client.RateLimit != nil {
		fmt.Printf("\nRate Limit: %d/%d (reset em %s)\n",
			client.RateLimit.Remaining,
			client.RateLimit.Limit,
			client.RateLimit.ResetTime.Format("15:04:05"))
	}
}
