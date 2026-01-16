// Exemplo avancado com timeout customizado e tratamento de erros
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/raphaeltorquat0/iptuapi-go"
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

	// Consulta com tratamento de erros
	fmt.Println("=== Consulta com Tratamento de Erros ===")
	resultado, err := client.ConsultaEndereco("Avenida Paulista", "1000")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("SQL: %s\n", resultado.DadosIPTU.SQL)
	fmt.Printf("Valor Venal: R$ %.2f\n", resultado.DadosIPTU.ValorVenal)

	// Consulta por SQL
	fmt.Println("\n=== Consulta por SQL ===")
	sqlResult, err := client.ConsultaSQL(resultado.DadosIPTU.SQL)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("SQL: %s\n", sqlResult.SQL)
	fmt.Printf("Ano: %d\n", sqlResult.Ano)
	fmt.Printf("Logradouro: %s, %s\n", sqlResult.Logradouro, sqlResult.Numero)
	fmt.Printf("Bairro: %s\n", sqlResult.Bairro)
	fmt.Printf("Area Terreno: %.2f m²\n", sqlResult.AreaTerreno)
	fmt.Printf("Area Construida: %.2f m²\n", sqlResult.AreaConstruida)
	fmt.Printf("Valor Venal Total: R$ %.2f\n", sqlResult.ValorVenal)
	fmt.Printf("IPTU Valor: R$ %.2f\n", sqlResult.IPTUValor)
}

func handleError(err error) {
	if iptuapi.IsAuthError(err) {
		fmt.Println("Erro: API Key invalida ou expirada")
	} else if iptuapi.IsNotFound(err) {
		fmt.Println("Erro: Imovel nao encontrado")
	} else if iptuapi.IsRateLimit(err) {
		fmt.Println("Erro: Rate limit excedido. Aguarde antes de tentar novamente.")
	} else if apiErr, ok := err.(*iptuapi.APIError); ok {
		fmt.Printf("Erro da API (status %d): %s\n", apiErr.StatusCode, apiErr.Message)
	} else {
		fmt.Printf("Erro: %v\n", err)
	}
}
