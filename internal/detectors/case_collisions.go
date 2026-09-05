package detectors

import (
	"context"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type CaseCollisions struct{}

func (CaseCollisions) ID() string { return "fs.case-collision" }

func (CaseCollisions) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	seen := make(map[string]string)
	var findings []model.Finding
	for _, rel := range project.Files {
		key := strings.ToLower(rel)
		if previous, ok := seen[key]; ok && previous != rel {
			findings = append(findings, model.Finding{
				RuleID:      "fs.case-collision",
				Title:       "Case-insensitive path collision",
				Description: previous + " and " + rel + " differ only by letter case.",
				Severity:    model.SeverityError,
				Path:        rel,
				Platforms:   []string{"windows", "macos"},
				Suggestion:  "Rename one path so the names differ beyond capitalization.",
			})
			continue
		}
		seen[key] = rel
	}
	return findings, nil
}
