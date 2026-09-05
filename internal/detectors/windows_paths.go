package detectors

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type WindowsPaths struct{}

func (WindowsPaths) ID() string { return "fs.windows-paths" }

func (WindowsPaths) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		for _, part := range strings.Split(rel, "/") {
			if strings.ContainsAny(part, `<>:"|?*`) {
				findings = append(findings, model.Finding{RuleID: "fs.windows-forbidden-char", Title: "Windows-forbidden path character", Description: "A path component contains a character that Windows does not permit in normal filenames.", Severity: model.SeverityError, Path: rel, Platforms: []string{"windows"}, Suggestion: "Rename the path using characters supported by Windows filesystems."})
				break
			}
		}
		if utf8.RuneCountInString(rel) >= 240 {
			findings = append(findings, model.Finding{RuleID: "fs.windows-long-path", Title: "Very long project path", Description: "This repository-relative path is close to common Windows path-length limits once a checkout directory is added.", Severity: model.SeverityWarning, Path: rel, Platforms: []string{"windows"}, Suggestion: "Shorten deeply nested directories or filenames, especially for generated dependency trees."})
		}
	}
	return findings, nil
}
