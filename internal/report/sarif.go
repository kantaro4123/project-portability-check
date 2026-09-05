package report

import (
	"encoding/json"
	"io"

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
	Name           string `json:"name"`
	Version        string `json:"version"`
	InformationURI string `json:"informationUri"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations,omitempty"`
}

type sarifMessage struct{ Text string `json:"text"` }
type sarifLocation struct{ PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"` }
type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}
type sarifArtifactLocation struct{ URI string `json:"uri"` }
type sarifRegion struct{ StartLine int `json:"startLine"` }

func WriteSARIF(w io.Writer, report model.Report) error {
	results := make([]sarifResult, 0, len(report.Findings))
	for _, finding := range report.Findings {
		result := sarifResult{RuleID: finding.RuleID, Level: sarifLevel(finding.Severity), Message: sarifMessage{Text: finding.Title + ": " + finding.Description}}
		if finding.Path != "" {
			physical := sarifPhysicalLocation{ArtifactLocation: sarifArtifactLocation{URI: finding.Path}}
			if finding.Line > 0 {
				physical.Region = &sarifRegion{StartLine: finding.Line}
			}
			result.Locations = []sarifLocation{{PhysicalLocation: physical}}
		}
		results = append(results, result)
	}
	doc := sarifDocument{Version: "2.1.0", Schema: "https://json.schemastore.org/sarif-2.1.0.json", Runs: []sarifRun{{Tool: sarifTool{Driver: sarifDriver{Name: "project-portability-check", Version: report.Version, InformationURI: "https://github.com/kantaro4123/project-portability-check"}}, Results: results}}}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
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
