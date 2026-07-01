## Resumo do Pull Request

Correção do catálogo de coleções ATS no perfil `SRE-Platform` do `profiles.yaml`. A coleção `remote-ai` não existia em `catalog/collections.yaml`, o que impedia o correto funcionamento do provedor ATS. A correção substitui por coleções fintech válidas e adiciona a empresa Stripe para monitoramento direto.

### Problema

O perfil `SRE-Platform` referenciava a coleção `remote-ai` no bloco `ats.collections`, porém esta coleção não estava definida em `catalog/collections.yaml`. Isso fazia com que o sistema de catálogo não resolvesse nenhuma empresa, resultando em zero vagas ATS sendo buscadas para o perfil.

### Solução

- Substituição da coleção inexistente `remote-ai` por coleções fintech reais definidas no catálogo:
  - `fintech` — todas as 15 empresas fintech
  - `fintech-payments` — empresas de pagamento (Stripe, Adyen, Modern Treasury, Increase, Lithic)
  - `fintech-banking` — empresas de banking (Mercury, Unit, Plaid)
  - `fintech-startups` — startups fintech (Mercury, Ramp, Alloy, Modern Treasury, Unit, Increase, Check, Pinwheel, Coast, Lithic)
- Adição de `stripe` como empresa específica no `ats.companies`

### Impacto

O perfil `SRE-Platform` agora busca vagas de **15 empresas fintech** via Greenhouse, cobrindo os segmentos de pagamentos, banking, startups e o ecossistema fintech completo.

### Arquivos Modificados

| Arquivo | Alteração |
|---------|-----------|
| `profiles.yaml` | Correção das coleções ATS no perfil SRE-Platform |

### Como Testar

```bash
go test ./internal/infrastructure/providers/... -v
go run cmd/bot/main.go
```

Verificar logs: as coleções fintech devem ser resolvidas para os board tokens correspondentes e o fetch de vagas ATS deve ocorrer sem erros.
