# Go SDK Status

**Última atualização:** 2026-01-24
**Versão:** v2.1.2
**Status:** 🟡 VERIFICAR

---

## Informações

| Item | Valor |
|------|-------|
| **Versão** | v2.1.2 |
| **Registry** | Go Modules (`go get github.com/raphaeltorquat0/iptuapi-go`) |
| **Status** | 🟡 VERIFICAR |
| **Mínimo** | Go 1.21+ |

## Instalação

```bash
go get github.com/raphaeltorquat0/iptuapi-go@latest
```

## Exemplo Rápido

```go
package main

import (
    "fmt"
    iptuapi "github.com/raphaeltorquat0/iptuapi-go"
)

func main() {
    client := iptuapi.NewClient("sua_api_key")
    cidades, _ := client.IPTUToolsCidades()
    fmt.Printf("%d cidades disponíveis\n", cidades.Total)
}
```

## Validação Automática

Este SDK é validado automaticamente:
- ✅ Instalação limpa via Go proxy
- ✅ Import do pacote
- ✅ Teste contra API real (`IPTUToolsCidades`)
- ✅ Teste autenticado (`ConsultaEndereco`)

---

*Atualizado automaticamente pelo CI/CD*
