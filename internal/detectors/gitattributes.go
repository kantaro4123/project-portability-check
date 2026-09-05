package detectors

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type GitAttributes struct{}

func (GitAttributes) ID() string { return "git.attributes" }

func (GitAttributes) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	hasAttributes := false
	hasScripts := false
	for _, rel := range project.Files {
		if strings.EqualFold(rel, ".gitattributes") {
			hasAttributes = true
		}
		switch strings.ToLower(filepath.Ext(rel)) {
		case ".sh", ".bash", ".zsh", ".ps1", ".cmd", ".bat":
			hasScripts = true
		}
	}
	if hasAttributes || !hasScripts {
		return nil, nil
	}
	return []model.Finding{{
		RuleID:      "git.no-gitattributes",
		Title:       "Cross-platform scripts without .gitattributes",
		Description: "The repository contains platform-sensitive scripts but does not define Git text/line-ending attributes.",
		Severity:    model.SeverityInfo,
		Suggestion:  "Consider a .gitattributes file that normalizes text and preserves required script line endings.",
	}}, nil
}
