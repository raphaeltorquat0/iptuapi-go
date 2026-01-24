// Exemplo de consulta IPTU em múltiplas cidades
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	iptuapi "github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
	apiKey := os.Getenv("IPTU_API_KEY")
	if apiKey == "" {
		log.Fatal("IPTU_API_KEY environment variable is required")
	}

	client := iptuapi.NewClient(apiKey)
	ctx := context.Background()

	// Lista de cidades disponíveis
	fmt.Println("=== Cidades Disponíveis ===")
	cidades, err := client.IPTUToolsCidades(ctx)
	if err != nil {
		log.Printf("Erro ao listar cidades: %v", err)
	} else {
		for _, c := range cidades.Cidades {
			fmt.Printf("  - %s (%s)\n", c.Nome, c.Codigo)
		}
	}

	// Consulta em São Paulo
	fmt.Println("\n=== São Paulo ===")
	spResult, err := client.ConsultaEndereco(ctx, &iptuapi.ConsultaEnderecoParams{
		Logradouro: "Avenida Paulista",
		Numero:     "1000",
		Cidade:     iptuapi.CidadeSaoPaulo,
	})
	if err != nil {
		log.Printf("Erro SP: %v", err)
	} else {
		fmt.Printf("  %s, %s - %s\n", spResult.Logradouro, spResult.Numero, spResult.Bairro)
		fmt.Printf("  Valor Venal: R$ %.2f\n", spResult.ValorVenalTotal)
	}

	// Consulta em Belo Horizonte
	fmt.Println("\n=== Belo Horizonte ===")
	bhResult, err := client.ConsultaEndereco(ctx, &iptuapi.ConsultaEnderecoParams{
		Logradouro: "Avenida Afonso Pena",
		Numero:     "1000",
		Cidade:     iptuapi.CidadeBeloHorizonte,
	})
	if err != nil {
		log.Printf("Erro BH: %v", err)
	} else {
		fmt.Printf("  %s, %s - %s\n", bhResult.Logradouro, bhResult.Numero, bhResult.Bairro)
		fmt.Printf("  Valor Venal: R$ %.2f\n", bhResult.ValorVenalTotal)
	}

	// Consulta em Rio de Janeiro
	fmt.Println("\n=== Rio de Janeiro ===")
	rjResult, err := client.ConsultaEndereco(ctx, &iptuapi.ConsultaEnderecoParams{
		Logradouro: "Avenida Atlântica",
		Numero:     "1000",
		Cidade:     iptuapi.CidadeRioDeJaneiro,
	})
	if err != nil {
		log.Printf("Erro RJ: %v", err)
	} else {
		fmt.Printf("  %s, %s - %s\n", rjResult.Logradouro, rjResult.Numero, rjResult.Bairro)
		fmt.Printf("  Valor Venal: R$ %.2f\n", rjResult.ValorVenalTotal)
	}

	// Calendário IPTU de cada cidade
	fmt.Println("\n=== Calendários IPTU 2026 ===")
	cidadesCodigos := []iptuapi.Cidade{
		iptuapi.CidadeSaoPaulo,
		iptuapi.CidadeBeloHorizonte,
		iptuapi.CidadeRioDeJaneiro,
	}

	for _, cidade := range cidadesCodigos {
		cal, err := client.IPTUToolsCalendario(ctx, cidade)
		if err != nil {
			log.Printf("Erro calendário %s: %v", cidade, err)
			continue
		}
		fmt.Printf("  %s: %d parcelas, %.0f%% desconto à vista\n",
			cal.Cidade, cal.ParcelasMax, cal.DescontoVistaPercentual)
	}
}
