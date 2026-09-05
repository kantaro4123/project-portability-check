package model

// Severity represents the impact of a portability finding.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Finding is a single portability issue detected in a project.
type Finding struct {
	RuleID      string   `json:"rule_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
	Path        string   `json:"path,omitempty"`
	Line        int      `json:"line,omitempty"`
	Platforms   []string `json:"platforms,omitempty"`
	Suggestion  string   `json:"suggestion,omitempty"`
}

// Summary is the aggregate result of an analysis.
type Summary struct {
	FilesScanned int `json:"files_scanned"`
	Errors       int `json:"errors"`
	Warnings     int `json:"warnings"`
	Info         int `json:"info"`
	Score        int `json:"score"`
}

// Report is the machine-readable analysis result.
type Report struct {
	Version         string    `json:"version"`
	Root            string    `json:"root"`
	TargetPlatforms []string  `json:"target_platforms,omitempty"`
	Summary         Summary   `json:"summary"`
	Findings        []Finding `json:"findings"`
}
