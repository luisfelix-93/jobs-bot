# Release Notes — v2.3.1

**Data:** 2026-07-01
**Branch:** `fix/error-catalog`

---

## Correção: Catálogo de Coleções ATS

### Problema
O perfil `SRE-Platform` estava configurado com a coleção `remote-ai`, que não existe no catálogo de empresas (`catalog/collections.yaml`). Como resultado, nenhuma vaga de sistemas ATS (Greenhouse) era buscada para este perfil.

### Correção
As coleções ATS do perfil `SRE-Platform` foram corrigidas para utilizar coleções reais do catálogo fintech:

| Coleção | Empresas Incluídas |
|---------|-------------------|
| `fintech` | Todas as 15 empresas (Stripe, Plaid, Brex, Mercury, Ramp, etc.) |
| `fintech-payments` | Stripe, Adyen, Modern Treasury, Increase, Lithic |
| `fintech-banking` | Mercury, Unit, Plaid |
| `fintech-startups` | Mercury, Ramp, Alloy, Modern Treasury, Unit, Increase, Check, Pinwheel, Coast, Lithic |
| Empresa específica | Stripe |

Agora o perfil `SRE-Platform` busca vagas de todas as empresas fintech disponíveis no catálogo Greenhouse.

### Arquivos Alterados
- `profiles.yaml` — correção das coleções ATS

### Como verificar
Execute o bot e observe os logs para confirmação de que as coleções fintech estão sendo resolvidas corretamente:

```bash
go run cmd/bot/main.go
```
