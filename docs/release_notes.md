# Job Search Bot v2.0 - Multi-Profile & AI Intelligence 🤖

We are excited to announce **v2.0** of the Job Search Bot! This release transforms the tool from a simple RSS scraper into a robust, intelligent job search assistant capable of managing multiple career profiles simultaneously.

## ✨ Key Features

### 1. Multi-Profile Support
Stop running multiple instances! You can now define multiple search profiles in `profiles.yaml`. 
- **Example**: Search for "Go Backend Developer" and "SRE Platform Engineer" at the same time.
- Each profile has its own resume, keywords, and source configurations.

### 2. Smart Deduplication with MongoDB
- We integrated **MongoDB** to track every job processed.
- **Benefit**: No more duplicate notifications. If a job was seen in the last 90 days, it won't bother you again.
- **Storage**: Supports both local MongoDB (Docker) and MongoDB Atlas (Cloud).

### 3. AI-Powered Analysis (DeepSeek) 🧠
- **DeepSeek Integration**: The bot now uses the DeepSeek API to analyze job descriptions against your resume.
- **Smart Scoring**: Get a compatibility score (0-100) and a summary of "Strengths" vs. "Gaps".
- **Falback**: If AI is unavailable, we gracefully fall back to keyword matching.

### 4. Consolidated Email Summaries 📧
- Receive a single daily email with the top job opportunities across all your profiles.
- HTML format with quick links to apply or view details.

### 5. Expanded Job Sources
- **JSearch (RapidAPI)**: Access to a broader range of job listings.
- **Findwork.dev**: Specific source for tech-focused roles.

## ⚠️ Breaking Changes & Upgrades

- **Configuration Overhaul**: The `profiles.yaml` file is now mandatory. Please migrate your settings from `.env` to this file.
- **New Dependencies**: 
  - **MongoDB**: A connection string is required in `.env` (`MONGO_URI`).
  - **AI Keys**: `DEEPSEEK_API_KEY` (Optional but recommended).
- **Environment Variables**: Updated list of required keys in `.env`.

## 🚀 Getting Started

1. **Update your `.env`**: Add `MONGO_URI`, `TRELLO_API_KEY`, etc.
2. **Create `profiles.yaml`**: Define your search profiles.
3. **Run**: 
   ```bash
   docker compose up -d
   go run cmd/bot/main.go
   ```

**Full Changelog**: https://github.com/luisfelix-93/jobs-bot/compare/v1.0.0...v2.0.0
