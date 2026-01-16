// Exemplo basico de uso do SDK
package main

import (
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

	// Consulta por endereco
	fmt.Println("=== Consulta por Endereco ===")
	resultado, err := client.ConsultaEndereco("Avenida Paulista", "1000")
	if err != nil {
		log.Fatalf("Erro na consulta: %v", err)
	}

	fmt.Printf("SQL Base: %s\n", resultado.Data.SQLBase)
	fmt.Printf("Logradouro: %s, %s\n", resultado.Data.Logradouro, resultado.Data.Numero)
	fmt.Printf("Bairro: %s\n", resultado.Data.Bairro)
	fmt.Printf("CEP: %s\n", resultado.Data.CEP)
	fmt.Printf("Area Terreno: %.2f m²\n", resultado.Data.AreaTerreno)
	fmt.Printf("Tipo Uso: %s\n", resultado.Data.TipoUso)

	// Dados IPTU detalhados
	fmt.Println("\n=== Dados IPTU ===")
	fmt.Printf("SQL: %s\n", resultado.DadosIPTU.SQL)
	fmt.Printf("Ano Referencia: %d\n", resultado.DadosIPTU.AnoReferencia)
	fmt.Printf("Area Construida: %.2f m²\n", resultado.DadosIPTU.AreaConstruida)
	fmt.Printf("Valor Venal: R$ %.2f\n", resultado.DadosIPTU.ValorVenal)
	fmt.Printf("Valor Terreno: R$ %.2f\n", resultado.DadosIPTU.ValorTerreno)
	fmt.Printf("Valor Construcao: R$ %.2f\n", resultado.DadosIPTU.ValorConstrucao)
}
