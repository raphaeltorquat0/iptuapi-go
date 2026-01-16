# Exemplos - Go SDK

## Executando os Exemplos

```bash
# Configurar API Key
export IPTU_API_KEY="sua_api_key"

# Executar exemplo basico
go run ./examples/basic/

# Executar exemplo avancado (tratamento de erros)
go run ./examples/advanced/

# Executar exemplo de valuation (requer plano Pro+)
go run ./examples/valuation/

# Executar exemplo multi-cidade
go run ./examples/multi-city/
```

## Exemplos Disponiveis

| Exemplo | Descricao | Plano |
|---------|-----------|-------|
| `basic/` | Consulta simples por endereco | Free |
| `advanced/` | Timeout customizado, tratamento de erros | Free |
| `valuation/` | Estimativa de valor (AVM + ITBI) | Pro+ |
| `multi-city/` | Consultas em SP, BH e Recife | Free |

## Cidades Suportadas

- `iptuapi.CidadeSaoPaulo` - Sao Paulo (sp)
- `iptuapi.CidadeBeloHorizonte` - Belo Horizonte (bh)
- `iptuapi.CidadeRecife` - Recife (recife) - inclui coordenadas
- `iptuapi.CidadePortoAlegre` - Porto Alegre (poa)
- `iptuapi.CidadeFortaleza` - Fortaleza (fortaleza)
- `iptuapi.CidadeCuritiba` - Curitiba (curitiba)
- `iptuapi.CidadeRioDeJaneiro` - Rio de Janeiro (rj)
- `iptuapi.CidadeBrasilia` - Brasilia (brasilia)
