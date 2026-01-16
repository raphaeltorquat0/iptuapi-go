// Exemplo de uso das ferramentas IPTU 2026
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

	// Listar cidades disponiveis
	fmt.Println("=== Cidades Disponiveis ===")
	cidades, err := client.IPTUToolsCidades(ctx)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	for _, c := range cidades.Cidades {
		fmt.Printf("  %s (%s) - Desconto: %s, Parcelas: %d\n",
			c.Nome, c.Codigo, c.DescontoVista, c.ParcelasMax)
	}

	// Calendario de Sao Paulo
	fmt.Println("\n=== Calendario SP 2026 ===")
	calendario, err := client.IPTUToolsCalendario(ctx, iptuapi.CidadeSaoPaulo)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("Desconto a vista: %.1f%%\n", calendario.DescontoVistaPercentual)
	fmt.Printf("Parcelas: ate %d\n", calendario.ParcelasMax)
	fmt.Printf("Proximo vencimento: %s (%d dias)\n",
		calendario.ProximoVencimento, calendario.DiasParaProximoVencimento)

	if len(calendario.Alertas) > 0 {
		fmt.Println("\nAlertas:")
		for _, a := range calendario.Alertas {
			fmt.Printf("  ⚠️  %s\n", a)
		}
	}

	// Simulador de pagamento
	fmt.Println("\n=== Simulador (IPTU R$ 2.000) ===")
	simulacao, err := client.IPTUToolsSimulador(ctx, &iptuapi.SimuladorParams{
		ValorIPTU:  2000,
		Cidade:     "sp",
		ValorVenal: 500000,
	})
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("A vista:    R$ %.2f (economia de R$ %.2f)\n",
		simulacao.ValorVista, simulacao.EconomiaVista)
	fmt.Printf("Parcelado:  %dx de R$ %.2f = R$ %.2f\n",
		simulacao.Parcelas, simulacao.ValorParcela, simulacao.ValorTotalParcelado)
	fmt.Printf("Recomendacao: %s\n", simulacao.Recomendacao)

	// Verificar isencao
	fmt.Println("\n=== Verificar Isencao ===")
	isencao, err := client.IPTUToolsIsencao(ctx, 250000, iptuapi.CidadeSaoPaulo)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("Valor venal: R$ %.2f\n", isencao.ValorVenal)
	fmt.Printf("Limite isencao: R$ %.2f\n", isencao.LimiteIsencao)
	fmt.Printf("Elegivel: %v\n", isencao.ElegivelIsencaoTotal)
	fmt.Printf("Mensagem: %s\n", isencao.Mensagem)
}
