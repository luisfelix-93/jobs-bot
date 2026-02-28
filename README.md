# Jobs Bot v2 🤖

Personal Job Search Assistant driven by Multi-Profile Configuration, AI Analysis, and Smart Deduplication.

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
  - JSearch (RapidAPI)
  - Findwork.dev
  - Jobicy
  - WeWorkRemotely
  - LinkedIn (RSS)

## Configuration

### 1. Environment Variables (`.env`)

```ini
# Trello
TRELLO_API_KEY="your_key"
TRELLO_API_TOKEN="your_token"

# MongoDB (Atlas or Local)
MONGO_URI="mongodb+srv://user:pass@cluster.mongodb.net/?retryWrites=true&w=majority"

# AI & APIs (Optional but recommended)
DEEPSEEK_API_KEY="sk-..."
JSEARCH_API_KEY="rapidapi_key"
FINDWORK_API_KEY="findwork_token"

# Email (SMTP)
SMTP_HOST="smtp.gmail.com"
SMTP_PORT=587
SMTP_USER="your_email@gmail.com"
SMTP_PASSWORD="app_password"
EMAIL_TO="recipient@email.com"

# Global Settings
JOB_LIMIT=50
```

### 2. Profiles (`profiles.yaml`)

Define your search profiles in the root directory:

```yaml
profiles:
  - name: "SRE-Platform"
    resume_path: "curriculos/RESUME_SRE.md"
    positive_keywords: ["Kubernetes", "Go", "AWS"]
    negative_keywords: ["Java", "Junior"]
    trello_list_id: "list_id_for_sre"
    sources:
      jsearch_query: "SRE Remote"
      findwork_search: "devops"

  - name: "DotNet-Backend"
    resume_path: "curriculos/RESUME_DOTNET.md"
    # ...
```

## How to Run

### Local Development

1. Start local MongoDB (optional, if not using Atlas):
   ```powershell
   docker compose up -d
   ```
2. Run the bot:
   ```powershell
   go run cmd/bot/main.go
   ```

### GitHub Actions

The workflow `.github/workflows/schedule.yml` runs automatically Mon-Fri at 09:00 UTC.
Ensure you add the following **Secrets** in your GitHub Repository settings:
- `MONGO_URI`
- `TRELLO_API_KEY`, `TRELLO_API_TOKEN`
- `DEEPSEEK_API_KEY`, `JSEARCH_API_KEY`, `FINDWORK_API_KEY`
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `EMAIL_TO`

## Architecture

- **Domain**: Core logic (`Job`, `ProfileStats`, `AIAnalysis`).
- **Application**: `JobService` orchestrates fetching -> filtering -> analyzing -> notifying.
- **Infrastructure**: Implementations for MongoDB, Trello, Email, DeepSeek, and Job APIs.
- **Config**: Loads `profiles.yaml` and environment variables.

### Workflow

```mermaid
graph TD
    Start([Cron / Manual Start]) --> Config[Load Config & profiles.yaml]
    Config --> ProfileLoop{For each Profile}
    
    ProfileLoop -->|Processing Profile| Fetch[Fetch Jobs from Sources<br/>JSearch, Findwork, LinkedIn, etc.]
    Fetch --> Dedup[Deduplication<br/>Check MongoDB]
    
    Dedup --> AILoop{For each new Job}
    AILoop -->|Score Job| AI[AI Analysis via DeepSeek<br/>Compare Job vs Resume]
    AI --> Trello[Create Trello Card<br/>Add Labels & Summary]
    Trello --> AILoop
    
    AILoop -->|Jobs Processed| ProfileLoop
    ProfileLoop -->|All Profiles Done| Email[Send Consolidated<br/>HTML Email Summary]
    Email --> Finish([End])
```