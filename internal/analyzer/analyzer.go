package analyzer

import (
	"context"
	"sort"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

// Project contains the normalized project inputs available to detectors.
type Project struct {
	Root  string
	Files []string
}

// Detector inspects a project and returns zero or more portability findings.
type Detector interface {
	ID() string
	Detect(context.Context, Project) ([]model.Finding, error)
}

// Analyzer runs a deterministic detector pipeline.
type Analyzer struct {
	detectors []Detector
}

func New(detectors ...Detector) *Analyzer {
	return &Analyzer{detectors: append([]Detector(nil), detectors...)}
}

func (a *Analyzer) Analyze(ctx context.Context, project Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, detector := range a.detectors {
		items, err := detector.Detect(ctx, project)
		if err != nil {
			return nil, err
		}
		findings = append(findings, items...)
	}
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].Severity != findings[j].Severity {
			return severityRank(findings[i].Severity) > severityRank(findings[j].Severity)
		}
		if findings[i].Path != findings[j].Path {
			return findings[i].Path < findings[j].Path
		}
		return findings[i].RuleID < findings[j].RuleID
	})
	return findings, nil
}

func severityRank(s model.Severity) int {
	switch s {
	case model.SeverityError:
		return 3
	case model.SeverityWarning:
		return 2
	default:
		return 1
	}
}
