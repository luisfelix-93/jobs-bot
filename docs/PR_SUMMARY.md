## Resumo do Pull Request

Corrige o bug em produção onde o bot processava apenas as primeiras 50 vagas (truncadas antes da verificação de duplicidade), causando perda de vagas novas que ficavam além do limite. Agora o limite é aplicado **apenas sobre vagas novas**, e a análise de IA só é executada em vagas que não existem no banco de dados, economizando tokens.

### Principais Alterações:

#### 1. Remoção do `limit` na filtragem (`internal/domain/filter.go`)

- Removido o parâmetro `limit` da função `FilterAndRankJobs`
- A função agora retorna **todas** as vagas elegíveis, ordenadas por score
- A responsabilidade de limitar a quantidade foi movida para a camada de aplicação

#### 2. Reestruturação do loop de processamento (`internal/application/job_service.go`)

- O loop agora percorre **todas** as vagas rankeadas
- A checagem de duplicidade (`s.store.Exists`) ocorre **antes** de qualquer análise
- A análise de IA (`analyzeWithAI`) só é chamada para vagas **novas** (não salvas no banco)
- O contador de limite (`s.limit`) só incrementa para vagas novas processadas
- O loop para quando `s.limit` vagas novas foram processadas

#### 3. Testes unitários (`internal/domain/filter_test.go`)

- Teste de retorno de todas as vagas válidas (sem truncamento)
- Teste de 100 vagas sem truncamento
- Teste com input vazio
- Teste de exclusão total por keyword negativa
- Teste de ordenação por score

#### 4. Testes unitários (`internal/application/job_service_test.go`)

- Teste de skip de duplicatas e descoberta de vagas novas
- Teste de limite aplicado apenas a vagas novas (cenário com 10 duplicatas + 10 novas, limit=5)
- Teste que garante que AI **não** é chamada para duplicatas (economia de tokens)
- Teste que confirma AI chamada apenas para vagas novas
- Teste de limite 0 (ilimitado)
