// Exemplo de uso do endpoint de Valuation (Pro+)
package main

import (
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

	// Estimativa de valor de mercado com parametros manuais
	fmt.Println("=== Valuation Estimate ===")
	avaliacao, err := client.ValuationEstimate(iptuapi.ValuationParams{
		AreaTerreno:    250,
		AreaConstruida: 180,
		Bairro:         "Pinheiros",
		Zona:           "ZM",
		TipoUso:        "Residencial",
		TipoPadrao:     "Medio",
		AnoConstrucao:  2010,
	})
	if err != nil {
		if apiErr, ok := err.(*iptuapi.APIError); ok && apiErr.StatusCode == 403 {
			fmt.Println("Este endpoint requer plano Pro ou superior")
			return
		}
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("Valor Estimado: R$ %.2f\n", avaliacao.ValorEstimado)
	fmt.Printf("Valor Minimo:   R$ %.2f\n", avaliacao.ValorMinimo)
	fmt.Printf("Valor Maximo:   R$ %.2f\n", avaliacao.ValorMaximo)
	fmt.Printf("Valor por m²:   R$ %.2f\n", avaliacao.ValorM2)
	fmt.Printf("Confianca:      %.1f%%\n", avaliacao.Confianca*100)
	fmt.Printf("Modelo Versao:  %s\n", avaliacao.ModeloVersao)

	// Avaliacao completa por SQL (combina AVM + ITBI)
	fmt.Println("\n=== Valuation Evaluate (por SQL) ===")
	evaluation, err := client.ValuationEvaluate(iptuapi.EvaluateParams{
		SQL:    "00904801381",
		Cidade: "sp",
	})
	if err != nil {
		if apiErr, ok := err.(*iptuapi.APIError); ok && apiErr.StatusCode == 403 {
			fmt.Println("Este endpoint requer plano Pro ou superior")
			return
		}
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("\nValor Final Estimado: R$ %.2f\n", evaluation.ValorFinal.Estimado)
	fmt.Printf("Valor Final Minimo:   R$ %.2f\n", evaluation.ValorFinal.Minimo)
	fmt.Printf("Valor Final Maximo:   R$ %.2f\n", evaluation.ValorFinal.Maximo)
	fmt.Printf("Metodo:               %s\n", evaluation.ValorFinal.Metodo)
	fmt.Printf("Confianca:            %.1f%%\n", evaluation.ValorFinal.Confianca*100)

	if evaluation.AvaliacaoAvm != nil {
		fmt.Println("\n--- AVM (Machine Learning) ---")
		fmt.Printf("  Valor: R$ %.2f\n", evaluation.AvaliacaoAvm.ValorEstimado)
		fmt.Printf("  Confianca: %.1f%%\n", evaluation.AvaliacaoAvm.Confianca*100)
	}

	if evaluation.AvaliacaoItbi != nil {
		fmt.Println("\n--- ITBI (Transacoes Reais) ---")
		fmt.Printf("  Valor: R$ %.2f\n", evaluation.AvaliacaoItbi.ValorEstimado)
		fmt.Printf("  Transacoes: %d\n", evaluation.AvaliacaoItbi.TotalTransacoes)
		fmt.Printf("  Periodo: %s\n", evaluation.AvaliacaoItbi.Periodo)
	}
}
