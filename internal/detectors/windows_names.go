package detectors

import (
	"context"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type WindowsNames struct{}

func (WindowsNames) ID() string { return "fs.windows-names" }

var windowsReserved = map[string]bool{
	"CON": true, "PRN": true, "AUX": true, "NUL": true,
	"COM1": true, "COM2": true, "COM3": true, "COM4": true, "COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
	"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true, "LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
}

func (WindowsNames) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		for _, part := range strings.Split(rel, "/") {
			base := part
			if dot := strings.IndexByte(base, '.'); dot >= 0 {
				base = base[:dot]
			}
			upper := strings.ToUpper(base)
			if windowsReserved[upper] {
				findings = append(findings, model.Finding{RuleID: "fs.windows-reserved", Title: "Windows-reserved filename", Description: part + " cannot be created normally on Windows.", Severity: model.SeverityError, Path: rel, Platforms: []string{"windows"}, Suggestion: "Rename the file or directory to a Windows-safe name."})
				break
			}
			if strings.HasSuffix(part, " ") || strings.HasSuffix(part, ".") {
				findings = append(findings, model.Finding{RuleID: "fs.windows-trailing", Title: "Windows-unsafe trailing character", Description: "Windows strips trailing spaces and dots from path components.", Severity: model.SeverityError, Path: rel, Platforms: []string{"windows"}, Suggestion: "Remove trailing spaces or dots from the path."})
				break
			}
		}
	}
	return findings, nil
}
