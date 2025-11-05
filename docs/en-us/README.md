# Jobs Bot

The Jobs Bot is a bot that searches for new job openings on LinkedIn, WeWorkRemotely, and Jobicy through an RSS feed, filters the jobs based on keywords, analyzes compatibility with your resume, and sends the best ones to a Trello board.

## Architecture

The project is divided into three main layers:

- **`cmd`**: The application's entry layer, where initialization and dependency injection are done.
- **`internal`**: The main application layer, divided into:
    - **`application`**: The service layer, which orchestrates the business logic.
    - **`domain`**: The domain layer, which contains the entities and the main business logic.
    - **`infrastructure`**: The infrastructure layer, which contains the implementations of repositories and external services.
- **`config`**: The configuration layer, which loads the application's settings from environment variables.

## Resume Analysis

One of the main features of the bot is the ability to analyze a job description and compare it with the content of a resume file in `.txt` format. The bot calculates a compatibility percentage based on keywords and reports which ones were found and which are missing.

The analysis result is added to the Trello card, making it easier to decide whether to apply for the job or not. The card title will contain the job source and the job title, and the description will have the analysis details.

## Configuration

The bot is configured through environment variables, which can be defined in an `.env` file in the project's root.

| Variable | Description |
| --- | --- |
| `LINKEDIN_RSS_URL` | The URL of the LinkedIn RSS feed with the job search filters. |
| `WEWORKREMOTELY_RSS_URL` | The URL of the WeWorkRemotely RSS feed with the job search filters. |
| `JOBICY_RSS_URL` | The URL of the Jobicy RSS feed with the job search filters. |
| `TRELLO_API_KEY` | The Trello API key. |
| `TRELLO_API_TOKEN` | The Trello API token. |
| `TRELLO_LIST_ID` | The ID of the Trello list where the job cards will be created. |
| `POSITIVE_KEYWORDS` | A comma-separated list of keywords to filter the jobs. |
| `NEGATIVE_KEYWORDS` | A comma-separated list of keywords to exclude the jobs. |
| `JOB_LIMIT` | The maximum number of jobs to be sent to Trello. |
| `RESUME_FILE_PATH` | The path to the resume file in `.txt` format. |

## Keywords

The `POSITIVE_KEYWORDS` are used for both the initial job filtering and the resume compatibility analysis. The bot checks which of these keywords are present in the job description and in your resume to calculate the score.

## How to run

1. Clone the repository:
```bash
git clone https://github.com/luisfelix-93/jobs-bot.git
```
2. Create a resume file in `.txt` format in the project's root (or in another location) and add your resume's content to it.

3. Create an `.env` file in the project's root and add the environment variables, as per the [Configuration](#configuration) section. Make sure that `RESUME_FILE_PATH` points to your resume file.

4. Install the dependencies:
```bash
go mod tidy
```
5. Run the bot:
```bash
go run cmd/bot/main.go
```
