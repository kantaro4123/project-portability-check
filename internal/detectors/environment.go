package detectors

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type EnvironmentVariables struct{}

func (EnvironmentVariables) ID() string { return "env.documentation" }

var envReferencePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bos\.Getenv\(["']([A-Z][A-Z0-9_]*)["']\)`),
	regexp.MustCompile(`\bprocess\.env\.([A-Z][A-Z0-9_]*)\b`),
	regexp.MustCompile(`\bENV\[["']([A-Z][A-Z0-9_]*)["']\]`),
}

func (EnvironmentVariables) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	files := normalizedFileSet(project.Files)
	if hasInAncestors(files, ".", ".env.example", ".env.sample", ".env.template", "env.example") {
		return nil, nil
	}
	vars := map[string]bool{}
	for _, rel := range project.Files {
		data, ok := readText(project.Root, rel)
		if !ok {
			continue
		}
		for _, re := range envReferencePatterns {
			for _, match := range re.FindAllSubmatch(data, -1) {
				if len(match) > 1 {
					vars[string(match[1])] = true
				}
			}
		}
	}
	if len(vars) == 0 {
		return nil, nil
	}
	names := make([]string, 0, len(vars))
	for name := range vars {
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) > 6 {
		names = append(names[:6], "…")
	}
	return []model.Finding{{
		RuleID:      "env.no-example",
		Title:       "Environment variables have no example file",
		Description: "Code references environment variables (" + strings.Join(names, ", ") + ") but no .env.example-style file was found.",
		Severity:    model.SeverityInfo,
		Suggestion:  "Document required environment variables or provide a safe example environment file.",
	}}, nil
}
