# Exemplos - Go SDK

## Executando os Exemplos

```bash
# Configurar API Key
export IPTU_API_KEY="sua_api_key"

# Executar exemplo basico
go run ./examples/basic/

# Executar exemplo avancado
go run ./examples/advanced/

# Executar exemplo de valuation (requer plano Pro+)
go run ./examples/valuation/

# Executar exemplo de IPTU Tools
go run ./examples/iptu-tools/
```

## Exemplos Disponiveis

| Exemplo | Descricao | Plano |
|---------|-----------|-------|
| `basic/` | Consulta simples por endereco | Free |
| `advanced/` | Retry, timeout, tratamento de erros | Free |
| `valuation/` | Estimativa de valor de mercado | Pro+ |
| `iptu-tools/` | Calendario, simulador, isencao | Free |
