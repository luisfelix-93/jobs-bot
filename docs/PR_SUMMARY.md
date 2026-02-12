# 🚀 Refactor: Multi-Profile Support, AI Analysis & MongoDB Integration

## 📋 Summary
This pull request introduces a major architectural overhaul (**Jobs Bot v2**) to support **multiple job search profiles**, **intelligent deduplication** via MongoDB, and **AI-powered resume analysis** using DeepSeek. It also establishes a robust CI/CD pipeline and expands job sources.

## 🛠️ Technical Details

### 1. Core Architecture Changes
- **Multi-Profile Configuration**: Replaced single-resume logic with `profiles.yaml`. The bot now iterates through defined profiles (e.g., "SRE", ".NET Backend"), each with its own resume, keywords, and sources.
- **MongoDB Integration**: Added `internal/infrastructure/mongodb` to persist processed jobs, enabling efficient deduplication (retention policy applied) and preventing repeated notifications.
- **AI Analysis**: Integrated **DeepSeek API** (`internal/infrastructure/deepseek`) to analyze job descriptions against resumes, providing match scores and qualitative feedback.

### 2. New Data Sources
- **JSearch (RapidAPI)**: Added support for JSearch API.
- **Findwork.dev**: Added support for Findwork.dev API.
- **Email Notifications**: Implemented `internal/infrastructure/email` to send consolidated daily summaries of findings.

### 3. CI/CD & DevOps
- **Workflows Added**: 
    - `ci-pr.yml`: Automates PR labeling and tagging (`v0.0.x-rc`).
    - `ci-release.yml`: Automates GitHub Releases and version tagging upon merge to main.
- **Secrets Management**: Updated `job.yaml` to include new supported secrets (`MONGODB_URI`, `DEEPSEEK_API_KEY`, etc.).

### 4. Configuration & Documentation
- **Updated `config` package**: Added structures for `ProfileConfig` and `Sources` to parse the new YAML configuration.
- **Revamped README**: Updated to reflect v2 features, configuration steps, and architecture.
- **Gitignore**: Updated to exclude local config/agent files and track documentation.

## ⚠️ Breaking Changes
- **Configuration Format**: Applications strictly rely on `profiles.yaml` now. Legacy single-resume environment variables (like `JobicyRssURL` in root `.env` for single run) are superseded by profile-specific configs.
- **Dependency**: A MongoDB connection (local or Atlas) is now required for deduplication.

## 🧪 Testing
- Validated loading of multiple profiles.
- Verified MongoDB connection and document insertion.
- Tested DeepSeek AI integration for mock job descriptions.
