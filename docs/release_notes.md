# Job# Release Notes v2.0.1 - Bug Fixes & Stability

Esta release engloba correções críticas de estabilidade referentes à captação de RSS e ao parsing das respostas de Inteligência Artificial vindas do modelo recém integrado (DeepSeek).

## 🐛 Bug Fixes

- **Correção da ingestão de RSS (Issue #01):** Os repositórios de vagas Jobicy e WeWorkRemotely agora validam corteramente respostas de HTTP fora da casa dos `200 OK` (ex: Server Error ou Cloudflare block) antes de parsearem o XML/JSON. Isso elimina crashes abruptos e sujos em tempo de execução quando a rede parceira apresenta instabilidade.
- **Extração Segura da IA (Issue #03):** Se o DeepSeek responder encapsulando o objeto de análise dentro de formatações de Code Blocks nativos do Markdown (` ```json `), a string final é higienizada (trimmed) para extirpar sufixos e prefixos literários. Isso restaura o funcinamento do parser nativo, prevenindo o `fallback` constante incorreto.

## 🛠️ Detalhes Adicionais
* **Hash dos Commits:**
   * `20260228 - validação de status code`
   * `20260228 - correção formatação DeepSeek`
