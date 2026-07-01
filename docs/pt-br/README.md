# Jobs Bot

Um sistema inteligente de automação de busca de empregos que monitora múltiplos sites de vagas, analisa oportunidades contra seu currículo usando IA, filtra baseado em suas preferências e organiza os resultados no Trello com resumos diários por email.

## Funcionalidades

- **Suporte a Múltiplos Perfis**: Configure múltiplos perfis de busca (ex: "SRE", "Backend .NET") via `profiles.yaml`.
- **Desduplicação Inteligente**: Usa MongoDB Atlas para rastrear vagas processadas e evitar duplicatas (retenção de 90 dias).
- **Pipeline de Normalização**: Estrutura dados brutos de múltiplas fontes em um modelo padronizado com senioridade, modalidade de trabalho, tipo de contratação, skills técnicas, faixa salarial e localização normalizada.
- **Suporte a Provedores ATS**: Busca vagas diretamente de sistemas ATS como Greenhouse, com catálogo YAML de empresas e coleções temáticas (fintech, etc.).
- **Análise por IA (DeepSeek)**: Analisa descrições de vagas contra seu currículo, fornecendo:
  - Pontuação de Compatibilidade (0-100)
  - Pontos Fortes e Lacunas
  - Recomendação (Candidatar/Revisar/Pular)
  - *Fallback para correspondência de palavras-chave se a IA estiver indisponível.*
- **Resumo Diário por Email**: Envia um email HTML consolidado com estatísticas, badges de senioridade/modalidade/salário e principais recomendações para todos os perfis.
- **Integração com Trello**: Cria cards ricos com tags de senioridade, modalidade, empresa, score da IA e seção de dados normalizados.
- **Múltiplas Fontes**:
  - Himalayas *(gratuito, sem necessidade de API Key)*
  - JSearch (RapidAPI)
  - Findwork.dev
  - Jobicy
  - WeWorkRemotely
  - LinkedIn (RSS)
  - TheirStack
  - **Greenhouse (ATS)** — Stripe, Mercury, Ramp, e mais via catálogo YAML

## Arquitetura

O projeto é dividido em três camadas principais:

- **`cmd`**: O ponto de entrada da aplicação, onde a inicialização e a injeção de dependência são feitas.
- **`internal`**: A camada principal da aplicação, dividida em:
  - **`application`**: A camada de serviço, que orquestra a lógica de negócios (`JobService`).
  - **`domain`**: A camada de domínio, que contém as entidades (`Job`, `ProcessedJob`, `AIAnalysis`, `ResumeAnalysis`, `ProfileStats`), a lógica de negócios (`JobFilter`, `ResumeAnalyzer`) e o pipeline de normalização (`normalization/`).
  - **`infrastructure`**: A camada de infraestrutura, que contém implementações para serviços externos:
    - **Fontes de Vagas**: `himalayas`, `jobicy`, `weworkremotely`, `linkedin`, `jsearch`, `findwork`, `theirstack`, `providers/ats` (Greenhouse, Lever, Ashby)
    - **IA**: `deepseek`
    - **Notificações**: `trello`, `email`
    - **Persistência**: `mongodb`
- **`config`**: A camada de configuração, que carrega as configurações do `profiles.yaml` e variáveis de ambiente.

### Entidades de Domínio

| Entidade | Descrição |
|----------|-----------|
| `Job` | Dados brutos da vaga de qualquer fonte (título, empresa, descrição, URL, etc.), agora com campos normalizados (senioridade, modalidade, skills, salário, etc.) |
| `ProcessedJob` | Vaga armazenada no MongoDB com resultados da análise, TTL e dados normalizados |
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
                               ▼
┌──────────────────────────────────────────────────────────────┐
│                      JobService                               │
│                                                               │
│  1. Buscar vagas (paralelo)                                   │
│     ├── APIs tradicionais (JSearch, Himalayas, etc.)          │
│     └── Provedores ATS (Greenhouse + catálogo YAML)           │
│  2. NORMALIZAR (pipeline com 7 normalizers)                   │
│  3. Filtrar & ranquear por palavras                           │
│  4. Desduplicar (MongoDB)                                     │
│  5. Análise IA / fallback palavras                            │
│  6. Armazenar com TTL de 90 dias                              │
└──────────────────────────────────────┬───────────────────────┘
                                       │
          ┌────────────────────────────┼────────────────────────┐
          │                            │                        │
          ▼                            ▼                        ▼
   ┌─────────────┐             ┌─────────────┐          ┌─────────────┐
   │ Fontes de   │             │   DeepSeek  │          │   Email     │
   │ Vagas (8)   │             │    IA       │          │  Resumo     │
   │ concurrentes│             └─────────────┘          │  c/ badges  │
   └──────┬──────┘                                      └─────────────┘
          │                                                    │
          ▼                                                    ▼
   ┌──────────────────────────────────────┐          ┌─────────────┐
   │           JobService                 │          │   Trello    │
   │                                      │          │  c/ tags    │
   │  - Buscar vagas (paralelo)           │          │  e dados    │
   │  - Normalizar (Pipeline)             │          └─────────────┘
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
| `GREENHOUSE_API_KEY` | Token Bearer para API do Greenhouse (opcional) | Opcional |
| `LEVER_API_KEY` | Chave da API do Lever (suporte futuro) | Opcional |
| `ASHBY_API_KEY` | Chave da API do Ashby (suporte futuro) | Opcional |
| `SMTP_HOST` | Host do servidor SMTP | Para email |
| `SMTP_PORT` | Porta do servidor SMTP | Para email |
| `SMTP_USER` | Usuário SMTP | Para email |
| `SMTP_PASSWORD` | Senha SMTP ou senha de app | Para email |
| `EMAIL_TO` | Endereço de email do destinatário | Para email |
| `JOB_LIMIT` | Máx de vagas por perfil (padrão: 10) | Não |

> **Nota:** A fonte do Himalayas não requer chave de API. A API Greenhouse é pública — a chave é opcional para autenticação avançada.

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
      ats:
        collections:
          - fintech
          - fintech-payments
          - fintech-banking
          - fintech-startups
        companies:
          - stripe

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
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `EMAIL_TO`

> A fonte do Himalayas não requer secret — configure apenas o `himalayas_query` no `profiles.yaml`.

## Pipeline de Normalização

Entre a coleta e a filtragem, um pipeline com 7 normalizers transforma dados brutos em um modelo padronizado:

| Normalizer | O que faz | Exemplo |
|-----------|-----------|---------|
| **Senioridade** | Extrai senioridade do título | `"Senior Go Engineer"` → `Senior` |
| **Modalidade** | Detecta Remote/Hybrid/On-site | `"Remote"` na localização → `Remote` |
| **Contratação** | Normaliza tipo de contratação | `"CLT"` → `FullTime` |
| **Título** | Remove prefixos/sufixos do título | `"Google - Dev (Remote)"` → `"Dev"` |
| **Skills** | Extrai skills técnicas da descrição | `"Go, Kubernetes, AWS"` |
| **Salário** | Parseia faixa salarial | `"$120k-$150k"` → `USD 120000-150000` |
| **Localização** | Padroniza nomes de países | `"USA"` → `"United States"` |

## Suporte a Provedores ATS

O bot busca vagas diretamente de sistemas ATS como o Greenhouse, utilizando um catálogo YAML de empresas.

### Como Configurar

```yaml
sources:
  ats:
    collections:
      - fintech          # Todas as empresas da coleção
    companies:
      - stripe           # Empresa específica
```

### Empresas Disponíveis (Greenhouse)

Stripe, Plaid, Brex, Mercury, Ramp, Alloy, Modern Treasury, Unit, Increase, Check, Pinwheel, Coast, Mesh, Lithic, Adyen — e mais podem ser adicionadas editando `catalog/greenhouse.yaml`.

### Arquitetura

```
catalog/
├── collections.yaml    # Coleções temáticas (fintech, fintech-payments, etc.)
├── greenhouse.yaml     # 15 empresas + board tokens
├── lever.yaml          # Placeholder para suporte futuro
└── ashby.yaml          # Placeholder para suporte futuro
```

Veja `docs/ATS-SUPPORT-GUIDE.md` para o guia completo de gerenciamento do catálogo e adição de novos provedores.

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
