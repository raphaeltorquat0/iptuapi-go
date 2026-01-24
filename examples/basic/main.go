// Exemplo básico de uso do SDK
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	iptuapi "github.com/raphaeltorquat0/iptuapi-go"
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

	// Consulta por endereço
	fmt.Println("=== Consulta por Endereço ===")
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
	fmt.Printf("CEP: %s\n", resultado.CEP)
	fmt.Printf("Área Terreno: %.2f m²\n", resultado.AreaTerreno)
	fmt.Printf("Área Construída: %.2f m²\n", resultado.AreaConstruida)
	fmt.Printf("Tipo Uso: %s\n", resultado.TipoUso)

	// Dados de valor
	fmt.Println("\n=== Valores ===")
	fmt.Printf("Valor Venal Total: R$ %.2f\n", resultado.ValorVenalTotal)
	fmt.Printf("Valor Venal Terreno: R$ %.2f\n", resultado.ValorVenalTerreno)
	fmt.Printf("Valor Venal Construção: R$ %.2f\n", resultado.ValorVenalConstrucao)
	fmt.Printf("IPTU: R$ %.2f\n", resultado.IPTUValor)

	// Exemplo IPTU Tools - Cidades
	fmt.Println("\n=== IPTU Tools - Cidades ===")
	cidades, err := client.IPTUToolsCidades(ctx)
	if err != nil {
		log.Printf("Erro ao buscar cidades: %v", err)
	} else {
		fmt.Printf("Total de cidades: %d\n", cidades.Total)
		for _, c := range cidades.Cidades {
			fmt.Printf("  - %s (%s)\n", c.Nome, c.Codigo)
		}
	}
}
