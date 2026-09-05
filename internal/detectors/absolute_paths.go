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
	name string
	re   *regexp.Regexp
}{
	{"macOS user path", regexp.MustCompile(`/Users/[A-Za-z0-9._-]+/`)},
	{"Linux home path", regexp.MustCompile(`/home/[A-Za-z0-9._-]+/`)},
	{"Windows user path", regexp.MustCompile(`[A-Za-z]:\\Users\\[^\\\r\n]+\\`)},
}

var desktopPlatforms = []string{"linux", "macos", "windows"}

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
					Description: "Found a " + pattern.name + " that may fail on another machine, including another machine running the same operating system.",
					Severity:    model.SeverityWarning,
					Path:        rel,
					Line:        lineNumber(data, loc[0]),
					Platforms:   append([]string(nil), desktopPlatforms...),
					Suggestion:  "Use a relative path, environment variable, or configurable project root.",
				})
			}
		}
	}
	return findings, nil
}
