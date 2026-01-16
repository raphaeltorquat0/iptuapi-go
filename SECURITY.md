# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 2.x.x   | :white_check_mark: |
| 1.x.x   | :x:                |

## Reporting a Vulnerability

A seguranca dos nossos usuarios e prioridade. Se voce descobrir uma vulnerabilidade de seguranca neste SDK, por favor reporte de forma responsavel.

### Como Reportar

1. **NAO** abra uma issue publica no GitHub
2. Envie um email para **security@iptuapi.com.br** com:
   - Descricao detalhada da vulnerabilidade
   - Passos para reproduzir o problema
   - Impacto potencial
   - Sugestao de correcao (se tiver)

### O Que Esperar

- **Confirmacao**: Responderemos em ate 48 horas confirmando o recebimento
- **Avaliacao**: Avaliaremos a vulnerabilidade em ate 7 dias
- **Correcao**: Vulnerabilidades criticas serao corrigidas em ate 30 dias
- **Credito**: Voce sera creditado no changelog (se desejar)

### Escopo

Este policy cobre:
- Codigo fonte do SDK
- Dependencias diretas
- Configuracoes de seguranca

Fora do escopo:
- A API em si (reporte para security@iptuapi.com.br separadamente)
- Aplicacoes de terceiros que usam o SDK

## Boas Praticas de Seguranca

### Proteja sua API Key

```go
// NUNCA faca isso
client := iptuapi.NewClient("sk_live_abc123") // API key hardcoded

// Faca isso
client := iptuapi.NewClient(os.Getenv("IPTU_API_KEY"))
```

### Use HTTPS

O SDK sempre usa HTTPS por padrao. Nunca desabilite a verificacao SSL em producao.

### Mantenha Atualizado

Sempre use a versao mais recente do SDK para ter as ultimas correcoes de seguranca.

```bash
go get -u github.com/raphaeltorquat0/iptuapi-go
```

## Vulnerabilidades Conhecidas

Nenhuma vulnerabilidade conhecida na versao atual.

Consulte o [Go Vulnerability Database](https://pkg.go.dev/vuln/) para verificar dependencias.
