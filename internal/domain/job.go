package domain


type Job struct {
	Title string
	Link  string
	GUID  string
}
type JobRepository interface {
	FetchJobs() ([]Job, error)
}

type NotificationService interface {
	Notify(job Job) error
}