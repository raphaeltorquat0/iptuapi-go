// Exemplo de consulta IPTU em multiplas cidades
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

	// Consulta em Sao Paulo
	fmt.Println("=== Sao Paulo ===")
	spResults, err := client.ConsultaIPTU(iptuapi.CidadeSaoPaulo, "Paulista", &iptuapi.ConsultaIPTUOptions{
		Limit: 3,
	})
	if err != nil {
		log.Printf("Erro SP: %v", err)
	} else {
		for _, r := range spResults {
			fmt.Printf("  %s, %v - %s - R$ %.2f\n", r.Logradouro, r.Numero, derefString(r.Bairro), r.ValorVenal)
		}
	}

	// Consulta em Belo Horizonte
	fmt.Println("\n=== Belo Horizonte ===")
	bhResults, err := client.ConsultaIPTU(iptuapi.CidadeBeloHorizonte, "Afonso Pena", &iptuapi.ConsultaIPTUOptions{
		Limit: 3,
	})
	if err != nil {
		log.Printf("Erro BH: %v", err)
	} else {
		for _, r := range bhResults {
			fmt.Printf("  %s, %v - %s - R$ %.2f\n", r.Logradouro, r.Numero, derefString(r.Bairro), r.ValorVenal)
		}
	}

	// Consulta em Recife (inclui coordenadas)
	fmt.Println("\n=== Recife ===")
	reResults, err := client.ConsultaIPTU(iptuapi.CidadeRecife, "Boa Viagem", &iptuapi.ConsultaIPTUOptions{
		Limit: 3,
	})
	if err != nil {
		log.Printf("Erro Recife: %v", err)
	} else {
		for _, r := range reResults {
			coords := ""
			if r.Latitude != nil && r.Longitude != nil {
				coords = fmt.Sprintf(" (%.6f, %.6f)", *r.Latitude, *r.Longitude)
			}
			fmt.Printf("  %s, %v - R$ %.2f%s\n", r.Logradouro, r.Numero, r.ValorVenal, coords)
		}
	}

	// Consulta por SQL/Identificador
	fmt.Println("\n=== Consulta por SQL (SP) ===")
	sqlResults, err := client.ConsultaIPTUSQL(iptuapi.CidadeSaoPaulo, "00904801381", nil)
	if err != nil {
		log.Printf("Erro: %v", err)
	} else {
		for _, r := range sqlResults {
			fmt.Printf("  SQL: %s\n", r.SQL)
			fmt.Printf("  Endereco: %s, %v\n", r.Logradouro, r.Numero)
			fmt.Printf("  Valor Venal: R$ %.2f\n", r.ValorVenal)
			fmt.Printf("  Area Terreno: %.2f m²\n", derefFloat(r.AreaTerreno))
			fmt.Printf("  Area Construida: %.2f m²\n", derefFloat(r.AreaConstruida))
		}
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
