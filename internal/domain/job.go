package domain


type Job struct {
	Title        string
	Link         string
	GUID         string
	SourceFeed   string
	Location     string
	FullDescription  string
}
type JobRepository interface {
	FetchJobs() ([]Job, error)
}

type NotificationService interface {
	Notify(job Job, analysis ResumeAnalysis) error
}