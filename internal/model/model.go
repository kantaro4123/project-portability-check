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

// FindingFingerprint identifies one exact emitted finding. It intentionally
// includes message metadata and the source line for precise machine output.
func FindingFingerprint(f Finding) string {
	platforms := append([]string(nil), f.Platforms...)
	for i := range platforms {
		platforms[i] = strings.ToLower(platforms[i])
	}
	sort.Strings(platforms)
	return hashParts(
		f.RuleID,
		normalizeFindingPath(f.Path),
		strconv.Itoa(f.Line),
		f.Title,
		f.Description,
		string(f.Severity),
		strings.Join(platforms, ","),
	)
}

// FindingIdentity is the deliberately conservative identity used by baselines.
// Stable rule ID, normalized path, and severity survive line movement and copy
// edits while count-aware matching still surfaces additional occurrences.
func FindingIdentity(f Finding) string {
	return hashParts(f.RuleID, normalizeFindingPath(f.Path), string(f.Severity))
}

func normalizeFindingPath(value string) string {
	return strings.ReplaceAll(value, "\\", "/")
}

func hashParts(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
