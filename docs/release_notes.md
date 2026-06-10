# Release Notes — v2.1.0

**Data:** 2026-06-10
**Branch:** `feature/001_fix_job_processing_limit`

---

## 🐛 Bug Fix: Correção do processamento de vagas

### Problema
O bot buscava ~100 vagas nas APIs mas **cortava para 50 antes** de verificar se já existiam no banco de dados. Se essas 50 vagas já tivessem sido processadas em execuções anteriores, o bot simplesmente as ignorava — e as outras 50 vagas (potencialmente novas) nunca eram avaliadas.

**Impacto em produção:** Perda significativa de aproveitamento de vagas, com o bot repetidamente processando o mesmo conjunto de vagas e ignorando oportunidades novas.

### Solução
- O filtro de vagas (`FilterAndRankJobs`) agora retorna **todas** as vagas elegíveis, sem truncamento.
- O limite de processamento é aplicado **apenas sobre vagas novas** (que não existem no banco de dados).
- A análise de IA (DeepSeek) **só é executada em vagas novas**, economizando tokens e custos.

### Fluxo corrigido
```
1. Buscar vagas de todas as APIs
2. Filtrar e ordenar por relevância (SEM limite)
3. Para cada vaga:
   → Verificar se já existe no banco de dados
   → Se existir: pular (não gasta tokens)
   → Se for nova: analisar com IA, notificar, salvar
   → Parar quando atingir o limite de vagas NOVAS
```

### Arquivos modificados
| Arquivo | Alteração |
|---------|-----------|
| `internal/domain/filter.go` | Removido parâmetro `limit` de `FilterAndRankJobs` |
| `internal/application/job_service.go` | Loop reestruturado para dedup-first + limit em novas |
| `internal/domain/filter_test.go` | **Novo** — 5 testes unitários |
| `internal/application/job_service_test.go` | **Novo** — 5 testes unitários |

### Como testar
```bash
go test ./internal/domain/ -v
go test ./internal/application/ -v
go build ./...
```
