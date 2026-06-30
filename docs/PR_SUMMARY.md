## Resumo do Pull Request

Pipeline de normalização de vagas que estrutura dados brutos de múltiplas fontes em um modelo padronizado. A normalização acontece entre a coleta e a filtragem, enriquecendo cada vaga com senioridade, modalidade de trabalho, tipo de contratação, skills, faixa salarial e localização padronizada — tudo extraído de forma determinística via regras e regex.

### Principais Alterações

#### 1. Pipeline de Normalização (`internal/domain/normalization/`)

Novo pacote com 7 normalizers e um pipeline orquestrador:

- **`normalizer.go`** — Interface `Normalizer` + struct `Pipeline` que executa normalizers em sequência
- **`seniority.go`** — Extrai senioridade do título: `"Senior"/"Sr"/"III"/"IV" → Senior`, `"Pleno"/"II"/"Mid" → Mid`, `"Junior"/"Jr"/"I" → Junior`, `"Staff" → Staff`, `"Principal" → Principal`, `"Lead"/"Tech Lead" → Lead`
- **`work_mode.go`** — Detecta modalidade analisando `Location` + `Title` + `Description`: `"remote"/"home office"/"remoto" → Remote`, `"hybrid"/"híbrido" → Hybrid`, `"on-site"/"presencial" → On-site`
- **`employment_type.go`** — Normaliza tipo de contratação: `"CLT"/"full time"/"permanent" → FullTime`, `"PJ"/"contractor"/"freelance" → Contract`, `"part time" → PartTime`
- **`title.go`** — Remove prefixos de empresa, sufixos de localização e tags de senioridade do título original, gerando `NormalizedTitle`
- **`skills.go`** — Extrai skills técnicas conhecidas (~50) do título + descrição via regex com word-boundary, com deduplicação e mapeamento de variantes (ex: `"Golang" → "Go"`)
- **`salary.go`** — Parseia faixas salariais com regex: `"$120k-$150k"`, `"USD 80,000 to 100,000"`, `"€100k"`, `"BRL 10.000 - BRL 12.000"` — extrai moeda + min + max
- **`location.go`** — Padroniza localizações: `"USA"/"US" → "United States"`, `"UK" → "United Kingdom"`, `"BR"/"Brasil" → "Brazil"`, `"Anywhere"/"Worldwide"/"Remote" → "Remote"`

#### 2. Expansão dos Models (`internal/domain/job.go`)

9 novos campos adicionados aos structs `Job` e `ProcessedJob`:

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `Company` | `string` | Nome da empresa |
| `Seniority` | `string` | Junior, Mid, Senior, Staff, Principal, Lead |
| `WorkMode` | `string` | Remote, Hybrid, On-site |
| `EmploymentType` | `string` | FullTime, Contract, PartTime |
| `Skills` | `[]string` | Skills técnicas extraídas |
| `SalaryMin` | `float64` | Salário mínimo |
| `SalaryMax` | `float64` | Salário máximo |
| `SalaryCurrency` | `string` | USD, EUR, BRL, GBP |
| `NormalizedTitle` | `string` | Título limpo |

#### 3. Atualização dos Providers

Dados nativos de cada fonte agora são mapeados:

- **Himalayas** — `CompanyName`, `Seniority`, `EmploymentType`, `MinSalary`, `MaxSalary`, `Currency`
- **JSearch** — `EmployerName` → `Company`, `JobIsRemote` → `WorkMode`
- **TheirStack** — `Company` → `Company`, `WorkMode: "Remote"` (todas as vagas)
- **Findwork** — `CompanyName` → `Company`

Providers sem dados nativos (Jobicy, LinkedIn, WeWorkRemotely) continuam funcionando — o normalizer preenche via análise de texto.

#### 4. Integração no JobService (`internal/application/job_service.go`)

- Pipeline executado entre `FetchJobs` e `FilterAndRankJobs`
- Logs de estatísticas pós-normalização: `"Seniority: 45/100, WorkMode: 80/100, Salary: 12/100"`
- Campos normalizados são persistidos no MongoDB via `ProcessedJob`

#### 5. Enriquecimento de Notificações

- **Email** (`internal/infrastructure/email/notification_service.go`): Badges coloridos de Seniority (azul), WorkMode (verde) e Salary (amarelo) no HTML; top-3 skills exibidas
- **Trello** (`internal/infrastructure/trello/notification_service.go`): Card title prefixado com tags `[AI Score] [Source] [Company] [Seniority] [WorkMode]`; seção "Informações Normalizadas" no body da card

#### 6. Persistência no MongoDB (`internal/infrastructure/mongodb/repository.go`)

- Novos campos com `bson:"...,omitempty"` para compatibilidade com documentos existentes

### Arquivos Modificados

| Arquivo | Alteração |
|---------|-----------|
| `internal/domain/job.go` | +9 campos em `Job` e `ProcessedJob` |
| `internal/domain/normalization/normalizer.go` | **Novo** — Interface + Pipeline |
| `internal/domain/normalization/seniority.go` | **Novo** — Normalizer de senioridade |
| `internal/domain/normalization/work_mode.go` | **Novo** — Normalizer de modalidade |
| `internal/domain/normalization/employment_type.go` | **Novo** — Normalizer de tipo de contratação |
| `internal/domain/normalization/title.go` | **Novo** — Normalizer de título |
| `internal/domain/normalization/skills.go` | **Novo** — Extrator de skills |
| `internal/domain/normalization/salary.go` | **Novo** — Normalizer de salário |
| `internal/domain/normalization/location.go` | **Novo** — Normalizer de localização |
| `internal/domain/normalization/normalizer_test.go` | **Novo** — 8 testes unitários |
| `internal/infrastructure/himalayas/repository.go` | Mapeamento de dados nativos |
| `internal/infrastructure/jsearch/repository.go` | Mapeamento de dados nativos |
| `internal/infrastructure/theirstack/repository.go` | Mapeamento de dados nativos |
| `internal/infrastructure/findwork/repository.go` | Mapeamento de dados nativos |
| `internal/application/job_service.go` | Integração do pipeline |
| `internal/application/job_service_test.go` | Ajuste para novo parâmetro |
| `internal/infrastructure/email/notification_service.go` | Badges de dados normalizados |
| `internal/infrastructure/trello/notification_service.go` | Tags e seção de dados normalizados |
| `internal/infrastructure/mongodb/repository.go` | Persistência dos novos campos |
| `cmd/bot/main.go` | Instanciação do pipeline |

### Como Testar

```bash
go test ./internal/domain/normalization/... -v
go build ./...
```
