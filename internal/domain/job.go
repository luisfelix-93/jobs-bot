package domain

import "time"

type Job struct {
	Title           string
	Link            string
	GUID            string
	SourceFeed      string
	Location        string
	FullDescription string
}

type JobRepository interface {
	FetchJobs() ([]Job, error)
}

type NotificationService interface {
	Notify(job Job, analysis ResumeAnalysis, aiAnalysis *AIAnalysis) error
}

type AIAnalysis struct {
	Score          int      `json:"score"`
	Strengths      []string `json:"strengths"`
	Gaps           []string `json:"gaps"`
	Recommendation string   `json:"recommendation"`
	Summary        string   `json:"summary"`
	Source         string   // "deepseek" or "keyword_fallback"
}

type AIAnalyzer interface {
	Analyze(resumeContent, jobDescription string) (*AIAnalysis, error)
}

type ProcessedJob struct {
	GUID            string
	Source          string
	Profile         string
	Title           string
	Link            string
	Location        string
	Description     string
	KeywordAnalysis ResumeAnalysis
	AIAnalysis      *AIAnalysis
	Notified        bool
	NotifiedAt      time.Time
	CreatedAt       time.Time
	TTLExpireAt     time.Time
}

type JobStore interface {
	Exists(guid, profile string) (bool, error)
	Save(job ProcessedJob) error
	Close() error
}

type ProfileStats struct {
	ProfileName    string
	TotalFound     int
	TotalFiltered  int
	TotalNotified  int
	TotalSkipped   int // duplicatas
	BelowThreshold int
	TopJobs        []ProcessedJob // top 3-5 jobs
}
