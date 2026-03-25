## Resumo do Pull Request

Este Pull Request resolve um problema de tipagem com o campo de identificação das vagas retornadas pela API do Findwork. A estrutura de dados interna e o processo de parse foram ajustados para acomodar o formato correto do respectivo ID.

### Principais Alterações:

#### 1. Correção da tipagem do ID do Findwork

- O tipo do campo `ID` na estrutura `findworkJob` foi alterado de `int` para `string` (`internal/infrastructure/findwork/repository.go`).
- A atribuição ao campo `GUID` no mapeamento de resposta modelo de domínio agora utiliza o valor da string de forma direta (`item.ID`), substituindo a antiga conversão `strconv.Itoa(item.ID)`.
- O import da biblioteca `strconv` foi removido do repositório, otimizando as dependências do pacote.
