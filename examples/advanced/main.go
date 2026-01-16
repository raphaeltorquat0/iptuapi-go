// Exemplo avancado com retry, timeout e tratamento de erros
package main

import (
	"context"
	"errors"
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

	// Cliente com configuracao avancada
	client := iptuapi.NewClient(apiKey,
		iptuapi.WithTimeout(60*time.Second),
		iptuapi.WithRetry(&iptuapi.RetryConfig{
			MaxRetries:      5,
			InitialDelay:    time.Second,
			MaxDelay:        30 * time.Second,
			BackoffFactor:   2.0,
			RetryableStatus: []int{429, 500, 502, 503, 504},
		}),
		iptuapi.WithLogger(&iptuapi.DefaultLogger{Enabled: true}),
	)

	// Context com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Consulta com tratamento de erros
	fmt.Println("=== Consulta com Tratamento de Erros ===")
	resultado, err := client.ConsultaEndereco(ctx, &iptuapi.ConsultaEnderecoParams{
		Logradouro:        "Avenida Paulista",
		Numero:            "1000",
		Cidade:            iptuapi.CidadeSaoPaulo,
		IncluirHistorico:  true,
		IncluirZoneamento: true,
	})

	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("SQL: %s\n", resultado.SQL)
	fmt.Printf("Valor Venal: R$ %.2f\n", resultado.ValorVenalTotal)

	// Historico
	if len(resultado.Historico) > 0 {
		fmt.Println("\n=== Historico ===")
		for _, h := range resultado.Historico {
			fmt.Printf("  %d: R$ %.2f\n", h.Ano, h.ValorVenalTotal)
		}
	}

	// Zoneamento
	if resultado.Zoneamento != nil {
		fmt.Println("\n=== Zoneamento ===")
		fmt.Printf("  Zona: %s (%s)\n", resultado.Zoneamento.Zona, resultado.Zoneamento.ZonaDescricao)
		fmt.Printf("  CA Basico: %.2f\n", resultado.Zoneamento.CoeficienteAproveitamentoBasico)
		fmt.Printf("  CA Maximo: %.2f\n", resultado.Zoneamento.CoeficienteAproveitamentoMaximo)
	}
}

func handleError(err error) {
	var authErr *iptuapi.AuthenticationError
	var forbiddenErr *iptuapi.ForbiddenError
	var notFoundErr *iptuapi.NotFoundError
	var rateLimitErr *iptuapi.RateLimitError
	var validationErr *iptuapi.ValidationError
	var serverErr *iptuapi.ServerError

	switch {
	case errors.As(err, &authErr):
		fmt.Println("Erro: API Key invalida")
	case errors.As(err, &forbiddenErr):
		fmt.Printf("Erro: Plano nao autorizado. Requer: %s\n", forbiddenErr.RequiredPlan)
	case errors.As(err, &notFoundErr):
		fmt.Println("Erro: Imovel nao encontrado")
	case errors.As(err, &rateLimitErr):
		fmt.Printf("Erro: Rate limit excedido. Retry em %d segundos\n", rateLimitErr.RetryAfter)
	case errors.As(err, &validationErr):
		fmt.Println("Erro: Parametros invalidos")
		for _, e := range validationErr.Errors {
			fmt.Printf("  - %s: %s\n", e.Field, e.Message)
		}
	case errors.As(err, &serverErr):
		fmt.Printf("Erro: Servidor (status %d)\n", serverErr.StatusCode)
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("Erro: Timeout")
	case errors.Is(err, context.Canceled):
		fmt.Println("Erro: Cancelado")
	default:
		fmt.Printf("Erro: %v\n", err)
	}
}
