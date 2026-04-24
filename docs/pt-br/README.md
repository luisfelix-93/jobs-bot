# Jobs Bot

Um sistema inteligente de automação de busca de empregos que monitora múltiplos sites de vagas, analisa oportunidades contra seu currículo usando IA, filtra baseado em suas preferências e organiza os resultados no Trello com resumos diários por email.

## Funcionalidades

- **Suporte a Múltiplos Perfis**: Configure múltiplos perfis de busca (ex: "SRE", "Backend .NET") via `profiles.yaml`.
- **Desduplicação Inteligente**: Usa MongoDB Atlas para rastrear vagas processadas e evitar duplicatas (retenção de 90 dias).
- **Análise por IA (DeepSeek)**: Analisa descrições de vagas contra seu currículo, fornecendo:
  - Pontuação de Compatibilidade (0-100)
  - Pontos Fortes e Lacunas
  - Recomendação (Candidatar/Revisar/Pular)
  - *Fallback para correspondência de palavras-chave se a IA estiver indisponível.*
- **Resumo Diário por Email**: Envia um email HTML consolidado com estatísticas e principais recomendações para todos os perfis.
- **Integração com Trello**: Cria cards ricos com resumos de IA e tags.
- **Múltiplas Fontes**:
  - Himalayas *(gratuito, sem necessidade de API Key)*
  - JSearch (RapidAPI)
  - Findwork.dev
  - Jobicy
  - WeWorkRemotely
  - LinkedIn (RSS)
  - TheirStack

## Arquitetura

O projeto é dividido em três camadas principais:

- **`cmd`**: O ponto de entrada da aplicação, onde a inicialização e a injeção de dependência são feitas.
- **`internal`**: A camada principal da aplicação, dividida em:
  - **`application`**: A camada de serviço, que orquestra a lógica de negócios (`JobService`).
  - **`domain`**: A camada de domínio, que contém as entidades (`Job`, `ProcessedJob`, `AIAnalysis`, `ResumeAnalysis`, `ProfileStats`) e a lógica de negócios (`JobFilter`, `ResumeAnalyzer`).
  - **`infrastructure`**: A camada de infraestrutura, que contém implementações para serviços externos:
    - **Fontes de Vagas**: `himalayas`, `jobicy`, `weworkremotely`, `linkedin`, `jsearch`, `findwork`, `theirstack`
    - **IA**: `deepseek`
    - **Notificações**: `trello`, `email`
    - **Persistência**: `mongodb`
- **`config`**: A camada de configuração, que carrega as configurações do `profiles.yaml` e variáveis de ambiente.

### Entidades de Domínio

| Entidade | Descrição |
|----------|-----------|
| `Job` | Dados brutos da vaga de qualquer fonte (título, empresa, descrição, URL, etc.) |
| `ProcessedJob` | Vaga armazenada no MongoDB com resultados da análise e TTL |
| `AIAnalysis` | Avaliação gerada pelo DeepSeek (pontuação, pontos fortes, lacunas, recomendação) |
| `ResumeAnalysis` | Resultados da correspondência por palavras-chave |
| `ProfileStats` | Estatísticas de processamento por perfil |

### Fluxo de Trabalho

```
┌─────────────────────────────────────────────────────────────┐
│                    config/profiles.yaml + .env              │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                     cmd/bot/main.go                          │
│  - Carregar perfis  - Inicializar MongoDB                    │
│  - Inicializar DeepSeek (opcional)  - Construir repositórios │
└──────────────────────────────┬──────────────────────────────┘
                               │
     ┌──────────────────────────┼──────────────────────────┐
     │                          │                          │
     ▼                          ▼                          ▼
┌─────────────┐          ┌─────────────┐           ┌─────────────┐
│ Fontes de   │          │   DeepSeek  │           │   Email     │
│ Vagas (6)   │          │    IA       │           │  Resumo     │
│ concurrentes│          └─────────────┘           └─────────────┘
└──────┬──────┘                                                  │
       │                                                  ▼
       ▼                                          ┌─────────────┐
┌──────────────────────────────────────┐          │   Trello    │
│           JobService                 │          │   Serviço   │
│                                      │          └─────────────┘
│  - Buscar vagas (paralelo)           │
│  - Filtrar & ranquear por palavras   │
│  - Desduplicar (MongoDB)             │
│  - Análise IA / fallback palavras    │
│  - Armazenar com TTL de 90 dias     │
└──────────────────────────────────────┘
```

## Configuração

### 1. Variáveis de Ambiente (`.env`)

| Variável | Descrição | Necessário |
|----------|-----------|-----------|
| `TRELLO_API_KEY` | Chave da API do Trello | Sim |
| `TRELLO_API_TOKEN` | Token da API do Trello | Sim |
| `MONGO_URI` | String de conexão do MongoDB | Sim |
| `DEEPSEEK_API_KEY` | Chave da API do DeepSeek para análise de IA | Recomendado |
| `JSEARCH_API_KEY` | Chave do RapidAPI para JSearch | Opcional |
| `FINDWORK_API_KEY` | Chave da API do findwork.dev | Opcional |
| `THEIRSTACK_API_KEY` | Chave da API do TheirStack | Opcional |

> **Nota:** A fonte do Himalayas não requer chave de API. Basta configurar `himalayas_query` no `profiles.yaml`.
| `SMTP_HOST` | Host do servidor SMTP | Para email |
| `SMTP_PORT` | Porta do servidor SMTP | Para email |
| `SMTP_USER` | Usuário SMTP | Para email |
| `SMTP_PASSWORD` | Senha SMTP ou senha de app | Para email |
| `EMAIL_TO` | Endereço de email do destinatário | Para email |
| `JOB_LIMIT` | Máx de vagas por perfil (padrão: 10) | Não |

### 2. Perfis (`profiles.yaml`)

Defina seus perfis de busca no diretório raiz:

```yaml
profiles:
  - name: "SRE-Platform"
    resume_path: "curriculos/RESUME_SRE.md"
    positive_keywords:
      - "Kubernetes"
      - "Go"
      - "AWS"
      - "Terraform"
    negative_keywords:
      - "Java"
      - "Junior"
    trello_list_id: "seu_id_da_lista_trello"
    sources:
      jsearch_query: "SRE Remote"
      findwork_search: "devops"
      theirstack_url: "https://api.theirstack.com/v1/jobs/search"
      himalayas_query: "golang devops kubernetes sre platform engineer"

  - name: "DotNet-Backend"
    resume_path: "curriculos/RESUME_DOTNET.md"
    positive_keywords:
      - ".NET"
      - "C#"
      - "SQL Server"
    negative_keywords:
      - "Java"
      - "Junior"
    trello_list_id: "id_de_outra_lista_trello"
    sources:
      jsearch_query: ".NET Backend Remote"
      himalayas_query: "dotnet c# backend asp.net azure"
```

### 3. Arquivos de Currículo

Coloque seus arquivos de currículo no diretório `curriculos/`. Os arquivos devem estar em formato markdown (`.md`).

## Como Executar

### Desenvolvimento Local

1. Inicie o MongoDB local (se não estiver usando Atlas):
   ```bash
   docker compose up -d
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

3. Execute o bot:
   ```bash
   go run cmd/bot/main.go
   ```

### GitHub Actions

O workflow `.github/workflows/schedule.yml` executa automaticamente de Seg-Sex às 09:00 UTC.

Adicione os seguintes **Secrets** nas configurações do seu Repositório GitHub:
- `MONGO_URI`
- `TRELLO_API_KEY`, `TRELLO_API_TOKEN`
- `DEEPSEEK_API_KEY`, `JSEARCH_API_KEY`, `FINDWORK_API_KEY`, `THEIRSTACK_API_KEY`

> A fonte do Himalayas não requer secret — configure apenas o `himalayas_query` no `profiles.yaml`.
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `EMAIL_TO`

## Análise de Currículo

O bot analisa descrições de vagas contra seu currículo usando dois métodos:

1. **Análise por IA (DeepSeek)**: Compara a descrição da vaga com o conteúdo do currículo para gerar:
   - Pontuação de Compatibilidade (0-100)
   - Pontos fortes identificados
   - Lacunas de habilidades
   - Recomendação (Candidatar/Revisar/Pular)

2. **Fallback por Palavras-Chave**: Se a IA estiver indisponível, o bot realiza correspondência de palavras-chave para calcular uma porcentagem de compatibilidade e identificar palavras-chave encontradas/faltantes.

Vagas com pontuação de IA >= 50 são enviadas ao Trello; vagas com pontuação menor são salvas mas não são notificadas.

## Estatísticas Rastreadas

Para cada perfil, o sistema rastreia:
- Total de vagas encontradas em todas as fontes
- Vagas restantes após filtragem
- Vagas notificadas (pontuação >= 50)
- Vagas abaixo do limiar (salvas mas não notificadas)