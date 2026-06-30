# Release Notes — v2.2.0

**Data:** 2026-06-30
**Branch:** `feature/10-normalization-jobs`

---

## ✨ Nova Feature: Pipeline de Normalização de Vagas

### Problema
As vagas coletadas de diferentes fontes (Himalayas, JSearch, Findwork, TheirStack, etc.) chegavam em formatos inconsistentes — umas com empresa, outras sem; umas com salário, outras com senioridade no título; cada uma com sua própria nomenclatura para o mesmo conceito. Isso dificultava a comparação, filtragem e exibição unificada dos dados.

### Solução
Pipeline de normalização com 7 normalizers que transformam dados brutos em um modelo padronizado e enriquecido:

### Novos Campos Normalizados

Cada vaga agora contém até 9 campos estruturados:

| Campo | Exemplo | Fonte |
|-------|---------|-------|
| Empresa | `Google`, `AWS` | Provider ou extraído |
| Senioridade | `Senior`, `Mid`, `Junior`, `Lead` | Título da vaga |
| Modalidade | `Remote`, `Hybrid`, `On-site` | Localização + Título + Descrição |
| Contratação | `FullTime`, `Contract`, `PartTime` | Título + Descrição |
| Skills | `Go, Kubernetes, AWS` | Descrição da vaga |
| Salário | `USD 120k-150k` | Título + Descrição |
| Localização | `United States`, `Brazil`, `Remote` | Normalizada |
| Título Limpo | `Software Engineer` (sem `(Remote)`) | Título original |

### O que muda na prática

**Antes:**
- Vagas apareciam como `"Google - Software Engineer (Remote) [Hybrid]"`
- Senioridade e salário enterrados na descrição
- Skills eram ignoradas na notificação
- Trello mostrava apenas título + fonte
- Email mostrava apenas título + score

**Depois:**
- Título é limpo para `"Software Engineer"` (com senoridade, modalidade e empresa em campos separados)
- Badges coloridos no email: Senioridade (azul), Modalidade (verde), Salário (amarelo)
- Top-3 skills exibidas no email
- Card do Trello com tags: `[AI: 85] [JSearch] [Google] [Senior] [Remote]`
- Seção "Informações Normalizadas" com todos os campos no body da card
- Salários parseados de formatos variados: `$120k`, `USD 80,000`, `€100k`, `BRL 10.000`
- Dados persistidos no MongoDB para consulta futura

### Fluxo Atualizado

```
1. Buscar vagas de todas as APIs
2. Normalizar (preencher lacunas via análise de texto)
   ├── Extrair senioridade do título
   ├── Detectar modalidade (Remote/Hybrid/On-site)
   ├── Detectar tipo de contratação
   ├── Extrair skills da descrição
   ├── Parsear faixa salarial
   ├── Padronizar localização
   └── Limpar título
3. Filtrar e ordenar por relevância
4. Desduplicar (MongoDB)
5. Analisar com IA / fallback
6. Notificar (Trello + Email com dados enriquecidos)
```

### Arquivos criados

| Arquivo | Descrição |
|---------|-----------|
| `internal/domain/normalization/normalizer.go` | Interface + Pipeline orquestrador |
| `internal/domain/normalization/seniority.go` | Normalizer de senioridade |
| `internal/domain/normalization/work_mode.go` | Normalizer de modalidade |
| `internal/domain/normalization/employment_type.go` | Normalizer de tipo de contratação |
| `internal/domain/normalization/title.go` | Normalizer de título |
| `internal/domain/normalization/skills.go` | Extrator de skills técnicas |
| `internal/domain/normalization/salary.go` | Parser de faixa salarial |
| `internal/domain/normalization/location.go` | Normalizer de localização |
| `internal/domain/normalization/normalizer_test.go` | Testes unitários (8 testes, todos passando) |

### Como testar

```bash
go test ./internal/domain/normalization/... -v
go build ./...
```
