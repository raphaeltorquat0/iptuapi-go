// Exemplo de uso do endpoint de Valuation (Pro+)
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
	apiKey := os.Getenv("IPTU_API_KEY")
	if apiKey == "" {
		log.Fatal("IPTU_API_KEY environment variable is required")
	}

	client := iptuapi.NewClient(apiKey)
	ctx := context.Background()

	// Estimativa de valor de mercado
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
		if iptuapi.IsForbidden(err) {
			fmt.Println("Este endpoint requer plano Pro ou superior")
			return
		}
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("Valor Estimado: R$ %.2f\n", avaliacao.ValorEstimado)
	fmt.Printf("Valor Minimo:   R$ %.2f\n", avaliacao.ValorMinimo)
	fmt.Printf("Valor Maximo:   R$ %.2f\n", avaliacao.ValorMaximo)
	fmt.Printf("Confianca:      %.1f%%\n", avaliacao.Confianca*100)
	fmt.Printf("Comparaveis:    %d\n", avaliacao.ComparaveisUtilizados)

	// Buscar comparaveis
	fmt.Println("\n=== Comparaveis ===")
	comparaveis, err := client.ValuationComparables(ctx, "Pinheiros", 150, 250, iptuapi.CidadeSaoPaulo, 5)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	for i, c := range comparaveis {
		fmt.Printf("%d. %s, %s - %.0f m² - R$ %.2f (%.0fm)\n",
			i+1, c.Logradouro, c.Numero, c.AreaConstruida, c.ValorVenalTotal, c.DistanciaMetros)
	}
}
