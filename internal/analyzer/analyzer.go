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
	if len(a.detectors) == 0 {
		return nil, nil
	}

	type detectorResult struct {
		index int
		items []model.Finding
		err   error
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan detectorResult, len(a.detectors))
	for index, detector := range a.detectors {
		index, detector := index, detector
		go func() {
			items, err := detector.Detect(ctx, project)
			results <- detectorResult{index: index, items: items, err: err}
		}()
	}

	batches := make([][]model.Finding, len(a.detectors))
	errs := make([]error, len(a.detectors))
	for range a.detectors {
		result := <-results
		batches[result.index] = result.items
		errs[result.index] = result.err
		if result.err != nil {
			cancel()
		}
	}
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	var findings []model.Finding
	for _, batch := range batches {
		findings = append(findings, batch...)
	}
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].Severity != findings[j].Severity {
			return severityRank(findings[i].Severity) > severityRank(findings[j].Severity)
		}
		if findings[i].Path != findings[j].Path {
			return findings[i].Path < findings[j].Path
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
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
