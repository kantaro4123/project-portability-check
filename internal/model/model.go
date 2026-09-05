package model

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
)

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
	Fingerprint string   `json:"fingerprint,omitempty"`
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
	Version            string    `json:"version"`
	Root               string    `json:"root"`
	TargetPlatforms    []string  `json:"target_platforms,omitempty"`
	BaselineSuppressed int       `json:"baseline_suppressed,omitempty"`
	Summary            Summary   `json:"summary"`
	Findings           []Finding `json:"findings"`
}

// AttachFingerprints populates stable exact fingerprints for machine output.
func AttachFingerprints(findings []Finding) {
	for i := range findings {
		if findings[i].Fingerprint == "" {
			findings[i].Fingerprint = FindingFingerprint(findings[i])
		}
	}
}

// FindingFingerprint identifies one finding at one location.
func FindingFingerprint(f Finding) string {
	return hashFinding(f, true)
}

// FindingIdentity identifies the semantic finding while intentionally ignoring
// line numbers. It is used for baseline matching so harmless line movement does
// not make an existing issue look new.
func FindingIdentity(f Finding) string {
	return hashFinding(f, false)
}

func hashFinding(f Finding, includeLine bool) string {
	platforms := append([]string(nil), f.Platforms...)
	for i := range platforms {
		platforms[i] = strings.ToLower(platforms[i])
	}
	sort.Strings(platforms)
	parts := []string{
		f.RuleID,
		strings.ReplaceAll(f.Path, "\\", "/"),
		f.Title,
		f.Description,
		string(f.Severity),
		strings.Join(platforms, ","),
	}
	if includeLine {
		parts = append(parts, strconv.Itoa(f.Line))
	}
	h := sha256.New()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
