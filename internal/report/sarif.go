package report

import (
	"encoding/json"
	"io"
	"sort"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

type sarifDocument struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string                     `json:"name"`
	Version        string                     `json:"version"`
	InformationURI string                     `json:"informationUri"`
	Rules          []sarifReportingDescriptor `json:"rules,omitempty"`
}

type sarifReportingDescriptor struct {
	ID               string       `json:"id"`
	Name             string       `json:"name,omitempty"`
	ShortDescription sarifMessage `json:"shortDescription"`
	Help             sarifMessage `json:"help"`
}

type sarifResult struct {
	RuleID              string            `json:"ruleId"`
	Level               string            `json:"level"`
	Message             sarifMessage      `json:"message"`
	Locations           []sarifLocation   `json:"locations,omitempty"`
	PartialFingerprints map[string]string `json:"partialFingerprints,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}
type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}
type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}
type sarifArtifactLocation struct {
	URI string `json:"uri"`
}
type sarifRegion struct {
	StartLine int `json:"startLine"`
}

func WriteSARIF(w io.Writer, report model.Report) error {
	results := make([]sarifResult, 0, len(report.Findings))
	for _, finding := range report.Findings {
		fingerprint := finding.Fingerprint
		if fingerprint == "" {
			fingerprint = model.FindingFingerprint(finding)
		}
		result := sarifResult{
			RuleID:  finding.RuleID,
			Level:   sarifLevel(finding.Severity),
			Message: sarifMessage{Text: finding.Title + ": " + finding.Description},
			PartialFingerprints: map[string]string{
				"projectPortabilityCheck/v1": fingerprint,
			},
		}
		if finding.Path != "" {
			physical := sarifPhysicalLocation{ArtifactLocation: sarifArtifactLocation{URI: strings.ReplaceAll(finding.Path, "\\", "/")}}
			if finding.Line > 0 {
				physical.Region = &sarifRegion{StartLine: finding.Line}
			}
			result.Locations = []sarifLocation{{PhysicalLocation: physical}}
		}
		results = append(results, result)
	}

	driver := sarifDriver{
		Name:           "project-portability-check",
		Version:        report.Version,
		InformationURI: "https://github.com/kantaro4123/project-portability-check",
		Rules:          sarifRules(report.Findings),
	}
	doc := sarifDocument{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs:    []sarifRun{{Tool: sarifTool{Driver: driver}, Results: results}},
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func sarifRules(findings []model.Finding) []sarifReportingDescriptor {
	byID := map[string]model.Finding{}
	for _, finding := range findings {
		if _, exists := byID[finding.RuleID]; !exists {
			byID[finding.RuleID] = finding
		}
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	rules := make([]sarifReportingDescriptor, 0, len(ids))
	for _, id := range ids {
		finding := byID[id]
		help := finding.Description
		if finding.Suggestion != "" {
			help += " Fix: " + finding.Suggestion
		}
		rules = append(rules, sarifReportingDescriptor{
			ID:               id,
			Name:             id,
			ShortDescription: sarifMessage{Text: finding.Title},
			Help:             sarifMessage{Text: help},
		})
	}
	return rules
}

func sarifLevel(severity model.Severity) string {
	switch severity {
	case model.SeverityError:
		return "error"
	case model.SeverityWarning:
		return "warning"
	default:
		return "note"
	}
}
