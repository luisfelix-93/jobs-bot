package himalayas

import (
	"encoding/json"
	"fmt"
)

// SearchResponse is the top-level JSON response from the Himalayas Jobs API.
type SearchResponse struct {
	UpdatedAt  int64 `json:"updatedAt"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
	TotalCount int   `json:"totalCount"`
	Jobs       []Job `json:"jobs"`
}

// Job represents a single job listing returned by the Himalayas API.
// All []string fields use StringList to survive the API's inconsistent typing.
type Job struct {
	Title                string       `json:"title"`
	Excerpt              string       `json:"excerpt"`
	CompanyName          string       `json:"companyName"`
	CompanySlug          string       `json:"companySlug"`
	CompanyLogo          string       `json:"companyLogo"`
	EmploymentType       string       `json:"employmentType"`
	MinSalary            *float64     `json:"minSalary"`
	MaxSalary            *float64     `json:"maxSalary"`
	Seniority            StringList   `json:"seniority"`
	Currency             string       `json:"currency"`
	LocationRestrictions LocationList `json:"locationRestrictions"`
	TimezoneRestrictions StringList   `json:"timezoneRestrictions"`
	Categories           StringList   `json:"categories"`
	ParentCategories     StringList   `json:"parentCategories"`
	Description          string       `json:"description"`
	PubDate              int64        `json:"pubDate"`
	ExpiryDate           int64        `json:"expiryDate"`
	ApplicationLink      string       `json:"applicationLink"`
	GUID                 string       `json:"guid"`
}

// Location represents a country restriction on a job listing.
type Location struct {
	Alpha2 string `json:"alpha2"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

// StringList is a []string that tolerates the Himalayas API's inconsistent typing.
// A field declared as StringList survives: null, a plain string, a number,
// an array of strings, and an array of numbers — all without panicking.
type StringList []string

// UnmarshalJSON absorbs all shapes the API may send for a string-array field.
func (s *StringList) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*s = StringList{}
		return nil
	}

	if data[0] == '[' {
		// Happy path: array of strings.
		var ss []string
		if err := json.Unmarshal(data, &ss); err == nil {
			*s = ss
			return nil
		}
		// Fallback: array of raw values (numbers, booleans…) — stringify each.
		var raws []json.RawMessage
		if err := json.Unmarshal(data, &raws); err != nil {
			return err
		}
		result := make(StringList, 0, len(raws))
		for _, r := range raws {
			result = append(result, string(r))
		}
		*s = result
		return nil
	}

	// Plain string.
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StringList{str}
		return nil
	}

	// Any other scalar (number, bool) — convert raw bytes to string.
	*s = StringList{fmt.Sprintf("%s", data)}
	return nil
}

// LocationList is a []Location that tolerates the Himalayas API returning
// locationRestrictions as an array of objects, array of strings, or a plain string.
type LocationList []Location

// UnmarshalJSON handles the inconsistent API contract for locationRestrictions.
func (l *LocationList) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*l = LocationList{}
		return nil
	}

	if data[0] == '[' {
		// Try array of Location objects first.
		var locs []Location
		if err := json.Unmarshal(data, &locs); err == nil {
			*l = locs
			return nil
		}
		// Fallback: array of plain strings e.g. ["USA", "Canada"].
		var names []string
		if err := json.Unmarshal(data, &names); err != nil {
			return err
		}
		result := make(LocationList, 0, len(names))
		for _, n := range names {
			result = append(result, Location{Name: n})
		}
		*l = result
		return nil
	}

	// Plain string — wrap into a single Location.
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*l = LocationList{{Name: str}}
	return nil
}
