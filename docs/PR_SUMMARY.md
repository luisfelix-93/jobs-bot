## Resumo do Pull Request

Este Pull Request implementa a integração com a **API pública do Himalayas Remote Jobs** como nova fonte de busca de vagas para os dois perfis configurados no projeto (`SRE-Platform` e `DotNet-Backend`). A API é gratuita, não exige autenticação e entrega resultados estruturados em JSON focados exclusivamente em vagas remotas.

---

### Principais Alterações

#### 1. Novo pacote de infraestrutura: `internal/infrastructure/himalayas`

Três arquivos foram criados seguindo o padrão modular já estabelecido no projeto:

- **`models.go`**: Define as structs `SearchResponse`, `Job`, `Location`, `LocationList` e `StringList` que espelham o schema JSON da API do Himalayas. Inclui unmarshalers customizados para lidar com o contrato inconsistente da API:
  - **`LocationList`**: o campo `locationRestrictions` pode chegar como `null`, array de objetos `[{alpha2, name, slug}]`, array de strings `["USA"]` ou string pura `"Worldwide"`.
  - **`StringList`**: campos `seniority`, `timezoneRestrictions`, `categories` e `parentCategories` podem chegar como `null`, array de strings, string pura ou número escalar. Todos os shapes são absorvidos sem panic.

- **`client.go`**: Implementa o `APIClient` com método `Search(ctx, SearchParams)`. Realiza requisições `GET` ao endpoint `https://himalayas.app/jobs/api/search` com query parameters (`q`, `employment_type`, `country`, `worldwide`, `page`, `limit`). Trata explicitamente:
  - `HTTP 429 Too Many Requests` — retorna erro descritivo sem panic
  - `HTTP 400 Bad Request` — retorna corpo da resposta para diagnóstico
  - Outros status inesperados — retorna código + corpo

- **`repository.go`**: Implementa `domain.JobRepository` via `FetchJobs()`. Realiza **duas buscas sequenciais** para garantir relevância geográfica para o Brasil:
  1. `country=BR` — vagas onde o Brasil é explicitamente listado como país aceito
  2. `worldwide=true` — vagas sem restrição geográfica (abertas para qualquer país)

  Os resultados são **merged e deduplicados por GUID**. Vagas restritas a países específicos como India são automaticamente excluídas por não aparecerem em nenhuma das duas queries. Cada busca pagina até `maxPages = 5` ou até consumir o `totalCount`.

- **`models_test.go`**: Testes unitários cobrindo todos os shapes de `LocationList` e `StringList`, além de um teste de integração com payload completo (`TestJob_FullPayloadWithMixedTypes`) simulando os tipos inconsistentes observados em produção.

#### 2. Extensão do contrato de configuração: `config/config.go`

- Adicionado campo `HimalayasQuery string \`yaml:"himalayas_query"\`` ao struct `Sources`.
- Nenhuma variável de ambiente adicional é necessária (a API não requer autenticação).

#### 3. Configuração de perfis: `profiles.yaml`

Ambos os perfis receberam o novo campo `himalayas_query`:

| Perfil | Query configurada |
|---|---|
| `SRE-Platform` | `"golang devops kubernetes sre platform engineer"` |
| `DotNet-Backend` | `"dotnet c# backend asp.net azure"` |

#### 4. Registro da fonte: `cmd/bot/main.go`

- Import do pacote `himalayas` adicionado.
- Bloco de registro da fonte adicionado em `buildRepos()`: se `sources.HimalayasQuery != ""`, o repositório é instanciado e adicionado à lista de fontes concorrentes. **Não requer API Key**, diferentemente das demais fontes opcionais.

---

### Bugs Corrigidos (pós-integração)

| Campo | Erro observado | Correção |
|---|---|---|
| `locationRestrictions` | `cannot unmarshal string into []Location` | Tipo `LocationList` com `UnmarshalJSON` customizado |
| `locationRestrictions` | `cannot unmarshal string into Location` (array de strings) | Branch adicional para `[]string → []Location` |
| `timezoneRestrictions` | `cannot unmarshal number into string` | Tipo `StringList` absorve number, string, array ou null |
| `seniority`, `categories`, `parentCategories` | Potencial falha silenciosa com tipos errados | Migrados para `StringList` preventivamente |

---

### Arquivos Criados

| Arquivo | Tipo |
|---|---|
| `internal/infrastructure/himalayas/models.go` | Novo |
| `internal/infrastructure/himalayas/client.go` | Novo |
| `internal/infrastructure/himalayas/repository.go` | Novo |
| `internal/infrastructure/himalayas/models_test.go` | Novo |

### Arquivos Modificados

| Arquivo | Alteração |
|---|---|
| `config/config.go` | +1 campo em `Sources` |
| `profiles.yaml` | +1 campo `himalayas_query` em cada perfil |
| `cmd/bot/main.go` | +1 import, +bloco em `buildRepos()` |
