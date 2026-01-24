// Exemplo de uso do endpoint de Valuation (Pro+)
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

	// Estimativa de valor de mercado com parâmetros manuais
	fmt.Println("=== Valuation Estimate ===")
	avaliacao, err := client.ValuationEstimate(ctx, &iptuapi.ValuationParams{
		AreaTerreno:    250,
		AreaConstruida: 180,
		Bairro:         "Pinheiros",
		Zona:           "ZM",
		TipoUso:        "Residencial",
		TipoPadrao:     "Medio",
		AnoConstrucao:  2010,
		Cidade:         iptuapi.CidadeSaoPaulo,
	})
	if err != nil {
		if apiErr, ok := err.(*iptuapi.APIError); ok && apiErr.StatusCode == 403 {
			fmt.Println("Este endpoint requer plano Pro ou superior")
			return
		}
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("Valor Estimado: R$ %.2f\n", avaliacao.ValorEstimado)
	fmt.Printf("Valor Mínimo:   R$ %.2f\n", avaliacao.ValorMinimo)
	fmt.Printf("Valor Máximo:   R$ %.2f\n", avaliacao.ValorMaximo)
	fmt.Printf("Confiança:      %.1f%%\n", avaliacao.Confianca*100)

	// Buscar imóveis comparáveis
	fmt.Println("\n=== Imóveis Comparáveis ===")
	comparaveis, err := client.ValuationComparables(ctx, "Pinheiros", 200, 300, iptuapi.CidadeSaoPaulo, 5)
	if err != nil {
		if apiErr, ok := err.(*iptuapi.APIError); ok && apiErr.StatusCode == 403 {
			fmt.Println("Este endpoint requer plano Pro ou superior")
			return
		}
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("Encontrados %d imóveis comparáveis:\n", len(comparaveis))
	for i, comp := range comparaveis {
		fmt.Printf("  %d. %s, %s - %.2f m² - R$ %.2f\n",
			i+1, comp.Logradouro, comp.Numero, comp.AreaTerreno, comp.ValorVenalTotal)
	}
}
