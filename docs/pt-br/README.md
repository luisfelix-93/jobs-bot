# Jobs Bot

O Jobs Bot é um bot que busca novas vagas de emprego no LinkedIn, WeWorkRemotely e Jobicy através de um feed RSS, filtra as vagas com base em palavras-chave, analisa a compatibilidade com o seu currículo e envia as melhores para um quadro do Trello.

## Arquitetura

O projeto é dividido em três camadas principais:

- **`cmd`**: A camada de entrada da aplicação, onde a inicialização e a injeção de dependência são feitas.
- **`internal`**: A camada principal da aplicação, dividida em:
    - **`application`**: A camada de serviço, que orquestra a lógica de negócios.
    - **`domain`**: A camada de domínio, que contém as entidades e a lógica de negócios principal.
    - **`infrastructure`**: A camada de infraestrutura, que contém as implementações de repositórios e serviços externos.
- **`config`**: A camada de configuração, que carrega as configurações da aplicação a partir de variáveis de ambiente.

## Análise de Currículo

Uma das funcionalidades principais do bot é a capacidade de analisar a descrição de uma vaga e compará-la com o conteúdo de um arquivo de currículo em formato `.txt`. O bot calcula uma porcentagem de compatibilidade com base nas palavras-chave e informa quais foram encontradas e quais estão faltando.

O resultado da análise é adicionado ao card do Trello, facilitando a decisão de se candidatar ou não para a vaga. O título do card conterá a fonte da vaga e o nome da vaga, e a descrição terá os detalhes da análise.

## Configuração

O bot é configurado através de variáveis de ambiente, que podem ser definidas em um arquivo `.env` na raiz do projeto.

| Variável | Descrição |
| --- | --- |
| `LINKEDIN_RSS_URL` | A URL do feed RSS do LinkedIn com os filtros de busca de vagas. |
| `WEWORKREMOTELY_RSS_URL` | A URL do feed RSS do WeWorkRemotely com os filtros de busca de vagas. |
| `JOBICY_RSS_URL` | A URL do feed RSS do Jobicy com os filtros de busca de vagas. |
| `TRELLO_API_KEY` | A chave da API do Trello. |
| `TRELLO_API_TOKEN` | O token da API do Trello. |
| `TRELLO_LIST_ID` | O ID da lista do Trello onde os cards de vagas serão criados. |
| `POSITIVE_KEYWORDS` | Uma lista de palavras-chave separadas por vírgula para filtrar as vagas. |
| `NEGATIVE_KEYWORDS` | Uma lista de palavras-chave separadas por vírgula para excluir as vagas. |
| `JOB_LIMIT` | O número máximo de vagas a serem enviadas para o Trello. |
| `RESUME_FILE_PATH` | O caminho para o arquivo de currículo em formato `.txt`. |

## Palavras-Chave

As `POSITIVE_KEYWORDS` são usadas tanto para a filtragem inicial de vagas quanto para a análise de compatibilidade do currículo. O bot verifica quais dessas palavras-chave estão presentes na descrição da vaga e no seu currículo para calcular a pontuação.

## Como executar

1. Clone o repositório:
```bash
git clone https://github.com/luisfelix-93/jobs-bot.git
```
2. Crie um arquivo de currículo em formato `.txt` na raiz do projeto (ou em outro local) e adicione o conteúdo do seu currículo a ele.

3. Crie um arquivo `.env` na raiz do projeto e adicione as variáveis de ambiente, conforme a seção de [Configuração](#configuração). Certifique-se de que `RESUME_FILE_PATH` aponte para o seu arquivo de currículo.

4. Instale as dependências:
```bash
go mod tidy
```
5. Execute o bot:
```bash
go run cmd/bot/main.go
```
