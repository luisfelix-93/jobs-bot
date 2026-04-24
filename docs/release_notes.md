# Release Notes

## [v0.0.11] — 2026-04-24

### ✨ Nova Fonte de Vagas: Himalayas Remote Jobs

O Jobs Bot agora busca vagas diretamente na **[Himalayas](https://himalayas.app)**, uma das maiores plataformas de empregos 100% remotos do mundo.

#### O que muda para você

- **Mais vagas encontradas por execução**: O bot agora consulta uma fonte adicional focada exclusivamente em trabalho remoto, aumentando o volume de oportunidades analisadas diariamente.
- **Zero configuração extra**: A API do Himalayas é pública e gratuita — nenhuma chave de API, cadastro ou variável de ambiente nova é necessária.
- **Busca personalizada por perfil**: Cada perfil de busca (`SRE-Platform`, `DotNet-Backend`) utiliza uma query diferente, alinhada com as tecnologias e palavras-chave que você já configurou.
- **Paginação automática**: O bot recupera múltiplas páginas de resultados em cada execução, garantindo que nenhuma vaga relevante seja perdida.

#### Fontes de vagas agora disponíveis

| Fonte | Tipo | Autenticação |
|---|---|---|
| Himalayas ✨ **novo** | JSON API (remoto) | Não requer |
| TheirStack | JSON API | API Key |
| Findwork.dev | JSON API | API Key |
| JSearch (RapidAPI) | JSON API | API Key |
| Jobicy | RSS | Não requer |
| WeWorkRemotely | RSS | Não requer |
| LinkedIn | RSS | Não requer |

#### Como configurar a query de busca

Edite seu `profiles.yaml` e adicione o campo `himalayas_query` na seção `sources` de cada perfil:

```yaml
sources:
  himalayas_query: "golang devops kubernetes sre"
```

Se o campo estiver ausente ou vazio, a fonte é ignorada silenciosamente.

---

## [v0.0.10] — Anterior

### 🐛 Correção

- **API de Vagas (Findwork):** Corrigido erro de conversão de tipo no campo `ID` da integração Findwork. O identificador interno foi ajustado de `int` para `string`, garantindo a correta ingestão de vagas no sistema.
