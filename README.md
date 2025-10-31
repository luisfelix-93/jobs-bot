# Jobs Bot

O Jobs Bot é um bot que busca novas vagas de emprego no LinkedIn através de um feed RSS, filtra as vagas com base em palavras-chave e envia as melhores para um quadro do Trello.

## Arquitetura

O projeto é dividido em três camadas principais:

- **`cmd`**: A camada de entrada da aplicação, onde a inicialização e a injeção de dependência são feitas.
- **`internal`**: A camada principal da aplicação, dividida em:
    - **`application`**: A camada de serviço, que orquestra a lógica de negócios.
    - **`domain`**: A camada de domínio, que contém as entidades e a lógica de negócios principal.
    - **`infrastructure`**: A camada de infraestrutura, que contém as implementações de repositórios e serviços externos.
- **`config`**: A camada de configuração, que carrega as configurações da aplicação a partir de variáveis de ambiente.

## Configuração

O bot é configurado através de variáveis de ambiente, que podem ser definidas em um arquivo `.env` na raiz do projeto.

| Variável               | Descrição                                                                 |
| ---------------------- | ------------------------------------------------------------------------- |
| `LINKEDIN_RSS_URL`     | A URL do feed RSS do LinkedIn com os filtros de busca de vagas.           |
| `TRELLO_API_KEY`       | A chave da API do Trello.                                                 |
| `TRELLO_API_TOKEN`     | O token da API do Trello.                                                 |
| `TRELLO_LIST_ID`       | O ID da lista do Trello onde os cards de vagas serão criados.             |
| `POSITIVE_KEYWORDS`    | Uma lista de palavras-chave separadas por vírgula para filtrar as vagas. |
| `NEGATIVE_KEYWORDS`    | Uma lista de palavras-chave separadas por vírgula para excluir as vagas.  |
| `JOB_LIMIT`            | O número máximo de vagas a serem enviadas para o Trello.                  |

## Como executar

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/jobs-bot.git
```
2. Crie um arquivo `.env` na raiz do projeto e adicione as variáveis de ambiente, conforme a seção de [Configuração](#configuração).

3. Instale as dependências:
```bash
go mod tidy
```
4. Execute o bot:
```bash
go run cmd/bot/main.go
```
