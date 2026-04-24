package himalayas

import (
	"encoding/json"
	"testing"
)

func TestLocationList_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantName string
	}{
		{
			name:     "array of objects (normal case)",
			input:    `[{"alpha2":"US","name":"United States","slug":"united-states"}]`,
			wantLen:  1,
			wantName: "United States",
		},
		{
			name:     "array of strings (actual API shape causing the bug)",
			input:    `["USA", "Canada"]`,
			wantLen:  2,
			wantName: "USA",
		},
		{
			name:     "plain string (inconsistent API)",
			input:    `"Worldwide"`,
			wantLen:  1,
			wantName: "Worldwide",
		},
		{
			name:    "null (no restriction)",
			input:   `null`,
			wantLen: 0,
		},
		{
			name:    "empty array",
			input:   `[]`,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ll LocationList
			if err := json.Unmarshal([]byte(tt.input), &ll); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(ll) != tt.wantLen {
				t.Errorf("got len=%d, want %d", len(ll), tt.wantLen)
			}
			if tt.wantLen > 0 && ll[0].Name != tt.wantName {
				t.Errorf("got name=%q, want %q", ll[0].Name, tt.wantName)
			}
		})
	}
}

func TestStringList_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantVal string
	}{
		{
			name:    "null",
			input:   `null`,
			wantLen: 0,
		},
		{
			name:    "empty array",
			input:   `[]`,
			wantLen: 0,
		},
		{
			name:    "array of strings (normal case)",
			input:   `["UTC+0", "UTC-5"]`,
			wantLen: 2,
			wantVal: "UTC+0",
		},
		{
			name:    "plain string",
			input:   `"UTC+0"`,
			wantLen: 1,
			wantVal: "UTC+0",
		},
		{
			name:    "number (actual API shape causing the bug)",
			input:   `42`,
			wantLen: 1,
			wantVal: "42",
		},
		{
			name:    "array of numbers",
			input:   `[1, 2, 3]`,
			wantLen: 3,
			wantVal: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sl StringList
			if err := json.Unmarshal([]byte(tt.input), &sl); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(sl) != tt.wantLen {
				t.Errorf("got len=%d, want %d", len(sl), tt.wantLen)
			}
			if tt.wantLen > 0 && sl[0] != tt.wantVal {
				t.Errorf("got value=%q, want %q", sl[0], tt.wantVal)
			}
		})
	}
}

// TestJob_FullPayloadWithMixedTypes verifies a full Job decode with all
// the inconsistent types the API has been observed sending in production.
func TestJob_FullPayloadWithMixedTypes(t *testing.T) {
	payload := `{
		"title": "Senior SRE",
		"companyName": "Acme",
		"locationRestrictions": ["USA", "Canada"],
		"timezoneRestrictions": 42,
		"seniority": "Senior",
		"categories": null,
		"parentCategories": [],
		"applicationLink": "https://example.com",
		"guid": "abc-123"
	}`

	var j Job
	if err := json.Unmarshal([]byte(payload), &j); err != nil {
		t.Fatalf("should not error on mixed types: %v", err)
	}
	if len(j.LocationRestrictions) != 2 {
		t.Errorf("locationRestrictions: got len=%d, want 2", len(j.LocationRestrictions))
	}
	if len(j.TimezoneRestrictions) != 1 || j.TimezoneRestrictions[0] != "42" {
		t.Errorf("timezoneRestrictions: got %v, want [42]", j.TimezoneRestrictions)
	}
	if len(j.Seniority) != 1 || j.Seniority[0] != "Senior" {
		t.Errorf("seniority: got %v, want [Senior]", j.Seniority)
	}
}
