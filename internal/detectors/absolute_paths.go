package detectors

import (
	"context"
	"regexp"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type AbsolutePaths struct{}

func (AbsolutePaths) ID() string { return "paths.absolute" }

var absolutePathPatterns = []struct {
	name      string
	platforms []string
	re        *regexp.Regexp
}{
	{"macOS user path", []string{"linux", "windows"}, regexp.MustCompile(`/Users/[A-Za-z0-9._-]+/`)},
	{"Linux home path", []string{"macos", "windows"}, regexp.MustCompile(`/home/[A-Za-z0-9._-]+/`)},
	{"Windows user path", []string{"macos", "linux"}, regexp.MustCompile(`[A-Za-z]:\\Users\\[^\\\r\n]+\\`)},
}

func (AbsolutePaths) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		data, ok := readText(project.Root, rel)
		if !ok {
			continue
		}
		for _, pattern := range absolutePathPatterns {
			for _, loc := range pattern.re.FindAllIndex(data, -1) {
				findings = append(findings, model.Finding{
					RuleID:      "paths.absolute",
					Title:       "Machine-specific absolute path",
					Description: "Found a " + pattern.name + " that may fail on another machine.",
					Severity:    model.SeverityWarning,
					Path:        rel,
					Line:        lineNumber(data, loc[0]),
					Platforms:   pattern.platforms,
					Suggestion:  "Use a relative path, environment variable, or configurable project root.",
				})
			}
		}
	}
	return findings, nil
}
