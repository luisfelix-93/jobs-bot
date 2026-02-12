# Enhancement Specs — Jobs Bot v2

> Documento de especificações das melhorias planejadas para o Jobs Bot.
> Baseado no brainstorm finalizado em 10/02/2026.

---

## 1. Multi-Perfil (profiles.yaml)

### Objetivo
Suportar múltiplos currículos com keywords, fontes de vagas e Trello boards independentes.

### Especificação

- Migrar configuração de perfis do `.env` para um arquivo `profiles.yaml` na raiz do projeto.
- Cada perfil contém:
  - **name** — Identificador único do perfil (ex: `SRE-Platform`)
  - **resume_path** — Caminho para o arquivo de currículo em `curriculos/`
  - **positive_keywords** — Lista de keywords positivas para filtro e análise
  - **negative_keywords** — Lista de keywords negativas para exclusão
  - **trello_list_id** — ID da lista Trello do perfil (cada perfil tem board separado)
  - **sources** — Configuração de queries por API (JSearch query, Findwork search/location)

### Perfis Iniciais

| Perfil | Currículo | Foco |
|--------|-----------|------|
| **SRE-Platform** | `curriculos/RESUME_LUIS_FELIX.md` | Go, Kubernetes, SRE, DevOps, Platform, Infrastructure |
| **DotNet-Backend** | `curriculos/Luis_Felipe_Felix_Filho_Resume.md` | .NET, C#, Backend, ASP.NET, Azure, React |

### Exemplo `profiles.yaml`

```yaml
profiles:
  - name: "SRE-Platform"
    resume_path: "curriculos/RESUME_LUIS_FELIX.md"
    positive_keywords: ["Go", "Golang", "Platform", "SRE", "Infrastructure", "DevOps", "Backend", "Remote", "Kubernetes"]
    negative_keywords: ["Java", "Frontend", "Manager", "Estágio", "Júnior"]
    trello_list_id: "690b91b6e6ff9a72c3a2e052"
    sources:
      jsearch_query: "SRE remote Europe"
      findwork_search: "golang devops"
      findwork_location: "remote"

  - name: "DotNet-Backend"
    resume_path: "curriculos/Luis_Felipe_Felix_Filho_Resume.md"
    positive_keywords: [".NET", "C#", "Backend", "ASP.NET", "Azure", "React", "SQL Server", "Remote"]
    negative_keywords: ["Java", "Frontend only", "Manager", "Estágio", "Júnior"]
    trello_list_id: "OUTRO_LIST_ID"
    sources:
      jsearch_query: ".NET developer remote Europe"
      findwork_search: "dotnet backend"
      findwork_location: "remote"
```

### Impacto no Código

| Arquivo | Mudança |
|---------|---------|
| `config/config.go` | Novo `ProfileConfig` struct, parser YAML, manter env vars globais (Trello keys, Mongo URI, etc.) |
| `cmd/bot/main.go` | Loop sobre perfis, instanciar `JobService` por perfil |
| `go.mod` | Nova dependência: `gopkg.in/yaml.v3` |

---

## 2. Novas APIs de Vagas

### 2.1 JSearch (RapidAPI / OpenWebNinja)

Agregador que puxa vagas do Google Jobs, LinkedIn, Indeed, Glassdoor, ZipRecruiter, Monster.

| Propriedade | Valor |
|------------|-------|
| **Endpoint** | `GET https://jsearch.p.rapidapi.com/search` |
| **Autenticação** | Header `X-RapidAPI-Key` + `X-RapidAPI-Host: jsearch.p.rapidapi.com` |
| **Free Tier** | 200 requests/mês |
| **Rate Limit** | 1000 req/hora |

**Parâmetros principais:** `query`, `page`, `num_pages`, `date_posted` (`today`/`week`/`month`), `remote_jobs_only`, `employment_types`.

**Mapeamento para `domain.Job`:**

| JSearch Field | domain.Job Field |
|--------------|-----------------|
| `job_id` | `GUID` |
| `job_title` | `Title` |
| `job_apply_link` | `Link` |
| `job_description` | `FullDescription` |
| `job_city` + `job_country` | `Location` |
| (hardcoded) | `SourceFeed = "JSearch"` |

**Arquivo:** `internal/infrastructure/jsearch/repository.go`

### 2.2 Findwork.dev (OpenPublicAPIs)

API dev-focused para vagas de TI.

| Propriedade | Valor |
|------------|-------|
| **Endpoint** | `GET https://findwork.dev/api/jobs/` |
| **Autenticação** | Header `Authorization: Token API_KEY` |
| **Free Tier** | Ilimitado |
| **Rate Limit** | 60 req/min |

**Parâmetros principais:** `search`, `location`, `remote`, `full_time`, `page`.

**Mapeamento para `domain.Job`:**

| Findwork Field | domain.Job Field |
|---------------|-----------------|
| `id` | `GUID` (como string) |
| `role` | `Title` |
| `url` | `Link` |
| `text` | `FullDescription` |
| `location` | `Location` |
| (hardcoded) | `SourceFeed = "Findwork"` |

**Arquivo:** `internal/infrastructure/findwork/repository.go`

### Env Vars Novas

```env
JSEARCH_API_KEY=your_rapidapi_key
FINDWORK_API_KEY=your_findwork_token
```

---

## 3. MongoDB Atlas — Persistência e Deduplicação

### Objetivo
Armazenar todas as vagas processadas para evitar notificações duplicadas entre execuções.

### Infraestrutura
- **Serviço:** MongoDB Atlas (free tier, 512MB)
- **Ambiente:** Separado do bot (cloud Atlas)
- **Go Driver:** `go.mongodb.org/mongo-driver`

### Schema — Collection `processed_jobs`

```json
{
  "_id": "ObjectId",
  "guid": "jsearch-abc123",
  "source": "JSearch",
  "profile": "SRE-Platform",
  "title": "Senior SRE Engineer",
  "link": "https://...",
  "location": "Remote, Europe",
  "description": "...",
  "keyword_analysis": {
    "match_percentage": 72.5,
    "found_keywords": ["Go", "Kubernetes"],
    "missing_keywords": ["Terraform"]
  },
  "ai_analysis": {
    "score": 85,
    "strengths": ["5+ anos Go", "experiência K8s"],
    "gaps": ["sem Terraform"],
    "recommendation": "apply",
    "summary": "Forte candidato, falta IaC"
  },
  "notified": true,
  "notified_at": "ISODate",
  "created_at": "ISODate",
  "ttl_expire_at": "ISODate (created_at + 90 dias)"
}
```

### Indexes

| Index | Tipo | Propósito |
|-------|------|-----------|
| `{guid, profile}` | Compound Unique | Deduplicação — mesma vaga pode ser relevante para perfis diferentes |
| `ttl_expire_at` | TTL (expireAfterSeconds: 0) | Auto-limpeza após 90 dias |
| `source` | Regular | Queries por fonte |

### Interface no Domínio

```go
// internal/domain/job.go
type JobStore interface {
    Exists(guid, profile string) (bool, error)
    Save(job ProcessedJob) error
}
```

### Fluxo de Dedup

1. Vaga chega do fetch → gera GUID: `"{source}-{id}"` (ex: `"jsearch-abc123"`)
2. Consulta `db.processed_jobs.findOne({guid, profile})`
3. **Se existe** → SKIP (já notificado para este perfil)
4. **Se não existe** → processa → analisa → salva → notifica

### Env Var

```env
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/jobs-bot?retryWrites=true&w=majority
```

**Arquivo:** `internal/infrastructure/mongodb/repository.go`

---

## 4. DeepSeek AI — Análise Semântica

### Objetivo
Substituir o matching básico por keywords por uma análise semântica usando IA, mantendo o sistema antigo como fallback.

### Infraestrutura

| Propriedade | Valor |
|------------|-------|
| **Endpoint** | `https://api.deepseek.com/chat/completions` |
| **Modelo** | `deepseek-chat` |
| **Custo** | ~$0.14/M tokens input, $0.28/M output |
| **Go Client** | `github.com/cohesion-org/deepseek-go` |

### Prompt

```
Você é um analista de vagas de emprego.
Compare o currículo abaixo com a descrição da vaga e avalie a compatibilidade.

CURRÍCULO:
{resume_content}

VAGA:
{job_description}

Retorne APENAS um JSON válido com:
{
  "score": 0-100,
  "strengths": ["competência que o candidato tem e a vaga requer"],
  "gaps": ["competência que a vaga requer e o candidato não tem"],
  "recommendation": "apply" | "review" | "skip",
  "summary": "análise em 2-3 frases"
}
```

### Threshold

- **Score ≥ 50** → Salva no MongoDB + Notifica via Trello + inclui no email
- **Score < 50** → Salva no MongoDB (para histórico) mas **não** notifica

### Estratégia de Fallback

```
1. TRY: DeepSeek AI analysis
2. CATCH (timeout, rate limit, API down):
   → Usa ResumeAnalyzer atual (keyword matching)
   → Mapeia resultado para o mesmo formato AIAnalysis
   → Marca source como "keyword_fallback"
3. Salva resultado no MongoDB (independente da fonte)
```

### Env Var

```env
DEEPSEEK_API_KEY=your_deepseek_api_key
```

**Arquivos:**
- `internal/infrastructure/deepseek/analyzer.go` — Client + prompt
- `internal/domain/resume_analyzer.go` — Interface `AIAnalyzer` + fallback

---

## 5. Email de Resumo Diário

### Objetivo
Enviar um email ao final de cada execução com o resumo das vagas selecionadas, scores e recomendações.

### Infraestrutura

| Propriedade | Valor |
|------------|-------|
| **Método** | SMTP direto via `net/smtp` (Go stdlib) |
| **Provedor** | Gmail com App Password |
| **Frequência** | 1x por execução (1x/dia) |

### Conteúdo do Email

Para cada perfil, uma seção com tabela:

```
📬 Jobs Bot — Resumo Diário (10/02/2026)

━━ Perfil: SRE-Platform ━━━━━━━━━━━━━━━
| Vaga                   | Fonte    | AI Score | Recomendação |
|------------------------|----------|----------|--------------|
| Senior SRE Engineer    | JSearch  | 85       | ✅ Apply     |
| Platform Engineer      | Findwork | 72       | 🔍 Review    |

━━ Perfil: DotNet-Backend ━━━━━━━━━━━━━
| Vaga                   | Fonte    | AI Score | Recomendação |
|------------------------|----------|----------|--------------|
| .NET Backend Developer | Findwork | 91       | ✅ Apply     |

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total: 3 vagas notificadas | 15 filtradas | 5 duplicadas
```

### Env Vars

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu.email@gmail.com
SMTP_PASSWORD=app_password_aqui
EMAIL_TO=seu.email@gmail.com
```

**Arquivo:** `internal/infrastructure/email/notification_service.go`

---

## 6. GitHub Actions — Atualização do Workflow

### Secrets Novos

| Secret | Descrição |
|--------|-----------|
| `JSEARCH_API_KEY` | RapidAPI key para JSearch |
| `FINDWORK_API_KEY` | Token da Findwork.dev |
| `MONGO_URI` | Connection string do MongoDB Atlas |
| `DEEPSEEK_API_KEY` | API key do DeepSeek |
| `SMTP_PASSWORD` | App password do Gmail |

### Mudanças no Workflow

- Passar todas as novas env vars como secrets
- Manter schedule `cron` 1x/dia
- Commit do `profiles.yaml` no repositório (não é secret)
