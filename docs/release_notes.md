# Release Notes — v2.3.0

**Data:** 2026-06-30
**Branch:** `feature/11-support-ats-provider`

---

## ✨ Nova Feature: Suporte a Provedores ATS (Greenhouse + Catálogo)

### Problema
O bot só buscava vagas de quadros de empregos genéricos (JSearch, Himalayas, etc.), perdendo oportunidades diretas de empresas que utilizam sistemas ATS (Applicant Tracking Systems). Muitas empresas de alto valor (Stripe, Mercury, Ramp, etc.) publicam vagas exclusivamente em suas boards ATS, e essas vagas não eram capturadas.

### Solução
Integração nativa com APIs públicas de ATS via sistema de catálogo YAML. Começamos com suporte ao **Greenhouse** (o ATS mais popular entre fintechs), com arquitetura extensível para Lever, Ashby e outros.

### Como Funciona

```
profiles.yaml → ATS Collections/Companies → catálogo YAML → Resolução → Fetch concorrente → domain.Job
```

1. **Catálogo YAML**: Arquivos em `catalog/` mapeiam empresas por provedor ATS com seus board tokens
2. **Coleções**: Agrupamentos temáticos de empresas (ex: `fintech`, `fintech-payments`)
3. **Fetch Concorrente**: Cada empresa é buscada em paralelo via goroutines
4. **Resiliência**: Falha em uma empresa não afeta as demais (erro é logado e ignorado)

### Empresas Catalogadas (Greenhouse)

| Empresa | Token | País | Categoria |
|---------|-------|------|-----------|
| Stripe | `stripe` | US | Fintech, Payments |
| Plaid | `plaid` | US | Fintech, Open Banking |
| Brex | `brex` | US | Fintech, Corporate Cards |
| Mercury | `mercury` | US | Fintech, Banking |
| Ramp | `ramp` | US | Fintech, Expense Management |
| Alloy | `alloy` | US | Fintech, Identity |
| Modern Treasury | `moderntreasury` | US | Fintech, Payments |
| Unit | `unit` | US | Fintech, Banking-as-a-Service |
| Increase | `increase` | US | Fintech, Payments API |
| Check | `check` | US | Fintech, Payroll |
| Pinwheel | `pinwheel` | US | Fintech, Payroll Connectivity |
| Coast | `coast` | US | Fintech, Fleet Payments |
| Mesh | `mesh` | US | Fintech, Crypto |
| Lithic | `lithic` | US | Fintech, Card Issuing |
| Adyen | `adyen` | NL | Fintech, Payments |

### Como Ativar

```yaml
# profiles.yaml — adicione no sources do perfil desejado
sources:
  ats:
    collections:
      - fintech        # Busca vagas de todas as empresas fintech
    companies:
      - stripe         # Empresa específica (não precisa estar em coleção)
```

### Arquitetura Extensível

Para adicionar um novo provedor ATS (ex: Lever), basta:
1. Implementar a interface `AtsClient` em `internal/infrastructure/providers/ats/lever/`
2. Registrar o cliente no `repository.go`
3. Criar o arquivo de catálogo `catalog/lever.yaml`

Veja `docs/ATS-SUPPORT-GUIDE.md` para o passo a passo completo.

### Como testar

```bash
go test ./internal/infrastructure/providers/... -v
go build ./...
```
