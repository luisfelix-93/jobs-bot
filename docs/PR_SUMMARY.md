## Resumo do Pull Request

Este Pull Request engloba as correções descritas no `bug-fixing.specs.md` abordando as Issues 01 e 03. Estes commits restauraram a estabilidade do fluxo de captação RSS e a higienização de carga proveniente da Inteligência Artificial.

### Principais Alterações:

#### 1. Validação de HTTP StatusCode (Issue #01)

- **Correção da ingestão de RSS:** Os arquivos `internal/infrastructure/jobicy/rss_repository.go` e `internal/infrastructure/weworkremotely/rss_repository.go` agora avaliam a viabilidade da resposta HTTP `(resp.StatusCode)` antes de efetuar o body parsing JSON/XML.
- **Prevenção de Crashes:** Evitamos falhas e pânicos da aplicação causados pelo parsing de documentos inválidos que mascaram erros crônicos de rede ou negações de serviço (como do Cloudflare) nas APIs parceiras.

#### 2. Extração segura da resposta LLM DeepSeek (Issue #03)

- **Sanitização de String:** Respostas da API DeepSeek agora passam por uma rotina de deleção de blocos formados como "markdown code chunks" via `strings.TrimPrefix` e `strings.TrimSuffix` no repositório `internal/infrastructure/deepseek/analyzer.go`.
- **Prevenção de Fallbacks falsos:** A correção garante que o Unmarshal nativo capture dados puramente JSON sem quebrar a deserialização da estrutura final por falhas de formatação não computacionais vindas da IA.
