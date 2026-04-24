# Jobs Bot

An intelligent job search automation system that monitors multiple job boards, analyzes opportunities against your resume using AI, filters based on your preferences, and organizes results in Trello with daily email summaries.

## Features

- **Multi-Profile Support**: Configure multiple search profiles (e.g., "SRE", "Backend .NET") via `profiles.yaml`.
- **Intelligent Deduplication**: Uses MongoDB Atlas to track processed jobs and prevent duplicates (90-day retention).
- **AI Analysis (DeepSeek)**: Analyzes job descriptions against your resume, providing:
  - Match Score (0-100)
  - Strengths & Gaps interpretation
  - Recommendation (Apply/Review/Skip)
  - *Fallback to Keyword Matching if AI is unavailable.*
- **Daily Email Summary**: Sends a consolidated HTML email with stats and top recommendations for all profiles.
- **Trello Integration**: Creates rich cards with AI summaries and tags.
- **Multiple Sources**:
  - Himalayas *(free, no API key required)*
  - JSearch (RapidAPI)
  - Findwork.dev
  - Jobicy
  - WeWorkRemotely
  - LinkedIn (RSS)
  - TheirStack

## Architecture

The project is divided into three main layers:

- **`cmd`**: The application's entry point, where initialization and dependency injection are done.
- **`internal`**: The main application layer, divided into:
  - **`application`**: The service layer, which orchestrates the business logic (`JobService`).
  - **`domain`**: The domain layer, which contains the entities (`Job`, `ProcessedJob`, `AIAnalysis`, `ResumeAnalysis`, `ProfileStats`) and business logic (`JobFilter`, `ResumeAnalyzer`).
  - **`infrastructure`**: The infrastructure layer, which contains implementations for external services:
    - **Job Sources**: `himalayas`, `jobicy`, `weworkremotely`, `linkedin`, `jsearch`, `findwork`, `theirstack`
    - **AI**: `deepseek`
    - **Notifications**: `trello`, `email`
    - **Persistence**: `mongodb`
- **`config`**: The configuration layer, which loads settings from `profiles.yaml` and environment variables.

### Domain Entities

| Entity | Description |
|--------|-------------|
| `Job` | Raw job data from any source (title, company, description, URL, etc.) |
| `ProcessedJob` | Job stored in MongoDB with analysis results and TTL |
| `AIAnalysis` | DeepSeek-generated evaluation (score, strengths, gaps, recommendation) |
| `ResumeAnalysis` | Keyword-based matching results |
| `ProfileStats` | Per-profile processing statistics |

### Workflow

```
┌─────────────────────────────────────────────────────────────┐
│                    config/profiles.yaml + .env              │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                     cmd/bot/main.go                          │
│  - Load profiles  - Initialize MongoDB                       │
│  - Initialize DeepSeek (optional)  - Build repositories     │
└──────────────────────────────┬──────────────────────────────┘
                               │
     ┌──────────────────────────┼──────────────────────────┐
     │                          │                          │
     ▼                          ▼                          ▼
┌─────────────┐          ┌─────────────┐           ┌─────────────┐
│ Job Sources │          │   DeepSeek  │           │   Email     │
│ (6 concurrent)         │    AI       │           │  Summary    │
└──────┬──────┘          └─────────────┘           └─────────────┘
       │                                                  │
       ▼                                                  ▼
┌──────────────────────────────────────┐         ┌─────────────┐
│           JobService                 │         │   Trello    │
│                                      │         │   Service   │
│  - Fetch jobs (parallel)             │         └─────────────┘
│  - Filter & rank by keywords         │                                  
│  - Deduplicate (MongoDB)             │
│  - AI analysis / keyword fallback    │
│  - Store with 90-day TTL             │
└──────────────────────────────────────┘
```

## Configuration

### 1. Environment Variables (`.env`)

| Variable | Description | Required |
|----------|-------------|----------|
| `TRELLO_API_KEY` | Trello API key | Yes |
| `TRELLO_API_TOKEN` | Trello API token | Yes |
| `MONGO_URI` | MongoDB connection string | Yes |
| `DEEPSEEK_API_KEY` | DeepSeek API key for AI analysis | Recommended |
| `JSEARCH_API_KEY` | RapidAPI key for JSearch | Optional |
| `FINDWORK_API_KEY` | findwork.dev API key | Optional |
| `THEIRSTACK_API_KEY` | TheirStack API key | Optional |

> **Note:** The Himalayas source requires no API key and is enabled simply by setting `himalayas_query` in `profiles.yaml`.
| `SMTP_HOST` | SMTP server host | For email |
| `SMTP_PORT` | SMTP server port | For email |
| `SMTP_USER` | SMTP username | For email |
| `SMTP_PASSWORD` | SMTP password or app password | For email |
| `EMAIL_TO` | Recipient email address | For email |
| `JOB_LIMIT` | Max jobs per profile (default: 10) | No |

### 2. Profiles (`profiles.yaml`)

Define your search profiles in the root directory:

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
    trello_list_id: "your_trello_list_id"
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
    trello_list_id: "your_other_list_id"
    sources:
      jsearch_query: ".NET Backend Remote"
      himalayas_query: "dotnet c# backend asp.net azure"
```

### 3. Resume Files

Place your resume files in the `curriculos/` directory. Files should be in markdown format (`.md`).

## How to Run

### Local Development

1. Start local MongoDB (if not using Atlas):
   ```bash
   docker compose up -d
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the bot:
   ```bash
   go run cmd/bot/main.go
   ```

### GitHub Actions

The workflow `.github/workflows/schedule.yml` runs automatically Mon-Fri at 09:00 UTC.

Add the following **Secrets** in your GitHub Repository settings:
- `MONGO_URI`
- `TRELLO_API_KEY`, `TRELLO_API_TOKEN`
- `DEEPSEEK_API_KEY`, `JSEARCH_API_KEY`, `FINDWORK_API_KEY`, `THEIRSTACK_API_KEY`

> The Himalayas source requires no secret — just configure `himalayas_query` in `profiles.yaml`.
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `EMAIL_TO`

## Resume Analysis

The bot analyzes job descriptions against your resume using two methods:

1. **AI Analysis (DeepSeek)**: Compares job description with resume content to generate:
   - Match Score (0-100)
   - Identified strengths
   - Skill gaps
   - Recommendation (Apply/Review/Skip)

2. **Keyword Fallback**: If AI is unavailable, the bot performs keyword matching to calculate a compatibility percentage and identify found/missing keywords.

Jobs with AI score >= 50 are sent to Trello; lower-scoring jobs are saved but not notified.

## Statistics Tracked

For each profile, the system tracks:
- Total jobs found across all sources
- Jobs remaining after filtering
- Jobs notified (score >= 50)
- Jobs below threshold (saved but not notified)