package theirstack

// SearchJobsRequest represents the request body for the TheirStack Jobs Search API
type SearchJobsRequest struct {
	Page                    int      `json:"page"`
	Limit                   int      `json:"limit"`
	OrderBy                 []Order  `json:"order_by,omitempty"`
	JobTitleOr              []string `json:"job_title_or,omitempty"`
	JobDescriptionPatternOr []string `json:"job_description_pattern_or,omitempty"`
	JobCountryCodeOr        []string `json:"job_country_code_or,omitempty"` // For worldwide we might leave it empty
	Remote                  *bool    `json:"remote,omitempty"`
	PostedAtGte             string   `json:"posted_at_gte,omitempty"`
	PostedAtMaxAgeDays      *int     `json:"posted_at_max_age_days,omitempty"`
}

type Order struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}

// SearchJobsResponse represents the response from TheirStack API
type SearchJobsResponse struct {
	Data    []TheirStackJob `json:"data"`
	HasMore bool            `json:"has_more"`
	Total   int             `json:"total"`
}

type TheirStackJob struct {
	ID          int64  `json:"id"`
	JobTitle    string `json:"job_title"`
	Company     string `json:"company"`
	URL         string `json:"url"`
	DatePosted  string `json:"date_posted"`
	Location    string `json:"location"`
	CountryCode string `json:"country_code"`
	Description string `json:"description"`
}

// Ensure the models cover what is needed to map to domain.Job
