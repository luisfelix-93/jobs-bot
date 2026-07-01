## Resumo do Pull Request

Suporte nativo a provedores de ATS (Applicant Tracking Systems) com integração inicial ao Greenhouse. O bot agora consegue buscar vagas diretamente das boards públicas do Greenhouse usando um sistema de catálogo baseado em YAML, permitindo monitorar empresas específicas (ex: Stripe, Mercury, Ramp) e coleções temáticas (ex: `fintech`, `remote-ai`) via configuração em `profiles.yaml`.

### Principais Alterações

#### 1. Arquitetura de Provedores ATS (`internal/infrastructure/providers/ats/`)

Novo pacote com arquitetura extensível para múltiplos provedores:

- **`client.go`** — Interface `AtsClient` com método `FetchJobs(boardToken string) ([]domain.Job, error)`, permitindo adicionar novos provedores sem alterar o orquestrador
- **`repository.go`** — Orquestrador concurrente que implementa `domain.JobRepository`, fazendo fetch paralelo para todas as empresas configuradas com tratamento resiliente de erros (falha em uma empresa não afeta as demais)
- **`catalog.go`** — Sistema de catálogo que carrega dinamicamente todos os arquivos YAML do diretório `/catalog`, resolve coleções em listas de empresas e trata duplicatas entre coleções+empresas
- **`greenhouse/greenhouse.go`** — Cliente nativo da API pública Greenhouse (`boards-api.greenhouse.io/v1/boards/{token}/jobs?content=true`), com suporte a autenticação opcional via Bearer token

#### 2. Sistema de Catálogo (`catalog/`)

Catálogo estruturado em YAML para fácil gerenciamento de empresas:

- **`catalog/collections.yaml`** — Coleções temáticas que agrupam empresas por categoria (ex: `fintech`, `fintech-payments`, `fintech-banking`, `fintech-startups`)
- **`catalog/greenhouse.yaml`** — 15 empresas fintech catalogadas com board tokens, nomes amigáveis, país, suporte a remoto, categorias e URL da página de carreiras:

| Empresa | Board Token | Categoria |
|---------|-------------|-----------|
| Stripe | `stripe` | Fintech, Payments |
| Plaid | `plaid` | Fintech, Open Banking |
| Brex | `brex` | Fintech, Corporate Cards |
| Mercury | `mercury` | Fintech, Banking |
| Ramp | `ramp` | Fintech, Expense Management |
| Alloy | `alloy` | Fintech, Identity |
| Modern Treasury | `moderntreasury` | Fintech, Payments |
| Unit | `unit` | Fintech, Banking-as-a-Service |
| Increase | `increase` | Fintech, Payments API |
| Check | `check` | Fintech, Payroll |
| Pinwheel | `pinwheel` | Fintech, Payroll Connectivity |
| Coast | `coast` | Fintech, Fleet Payments |
| Mesh | `mesh` | Fintech, Crypto |
| Lithic | `lithic` | Fintech, Card Issuing |
| Adyen | `adyen` | NL, Payments |

- **`catalog/lever.yaml`** e **`catalog/ashby.yaml`** — Placeholders para suporte futuro

#### 3. Configuração (`config/config.go`)

- Nova struct `AtsConfig` com campos `Collections []string` e `Companies []string`
- Novo campo `Ats AtsConfig` no struct `Sources` (configurável por perfil em `profiles.yaml`)
- Novas variáveis de ambiente: `GREENHOUSE_API_KEY`, `LEVER_API_KEY`, `ASHBY_API_KEY`

#### 4. Integração no `cmd/bot/main.go`

- `buildRepos` agora aceita `*config.Config` como parâmetro
- Se o perfil tiver `ats.collections` ou `ats.companies` configurados, instancia `ats.NewRepository`

#### 5. Guia de Suporte (`docs/ATS-SUPPORT-GUIDE.md`)

Documentação detalhada cobrindo:
- Arquitetura do sistema ATS (diagrama Mermaid)
- Gerenciamento do catálogo de empresas
- Passo a passo para adicionar novas empresas ao catálogo
- Guia para adicionar suporte a novos provedores ATS (Lever, Ashby, etc.)
- Extensibilidade via plugin (interface `AtsClient` + YAML + switch-case)

### Exemplo de Uso

```yaml
# profiles.yaml
sources:
  ats:
    collections:
      - fintech
    companies:
      - stripe
```

```bash
# .env
GREENHOUSE_API_KEY=optional_bearer_token
```

### Arquivos Modificados

| Arquivo | Alteração |
|---------|-----------|
| `internal/infrastructure/providers/ats/client.go` | **Novo** — Interface AtsClient |
| `internal/infrastructure/providers/ats/repository.go` | **Novo** — Orquestrador concurrente |
| `internal/infrastructure/providers/ats/catalog.go` | **Novo** — Loader/Resolver de catálogo |
| `internal/infrastructure/providers/ats/catalog_test.go` | **Novo** — Testes do catálogo |
| `internal/infrastructure/providers/ats/repository_test.go` | **Novo** — Testes do repositório (falha resiliente) |
| `internal/infrastructure/providers/ats/greenhouse/greenhouse.go` | **Novo** — Cliente Greenhouse |
| `internal/infrastructure/providers/ats/greenhouse/greenhouse_test.go` | **Novo** — Testes do cliente Greenhouse |
| `config/config.go` | +3 env vars, +1 struct (AtsConfig), +1 field (Ats) |
| `cmd/bot/main.go` | Integração do ATS repository |
| `catalog/collections.yaml` | **Novo** — Coleções de empresas |
| `catalog/greenhouse.yaml` | **Novo** — 15 empresas fintech |
| `catalog/lever.yaml` | **Novo** — Placeholder |
| `catalog/ashby.yaml` | **Novo** — Placeholder |
| `profiles.yaml` | Exemplo de configuração ATS |
| `docs/ATS-SUPPORT-GUIDE.md` | **Novo** — Guia completo de suporte ATS |

### Como Testar

```bash
go test ./internal/infrastructure/providers/... -v
go build ./...
```
