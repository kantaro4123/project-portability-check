package baseline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

// Baseline stores semantic finding counts from a previous JSON report. Counts
// are used instead of a set so a newly introduced duplicate finding in the same
// file is still surfaced.
type Baseline struct {
	counts map[string]int
}

func Load(filename string) (Baseline, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Baseline{}, err
	}
	var report model.Report
	if err := json.Unmarshal(data, &report); err != nil {
		return Baseline{}, fmt.Errorf("parse baseline JSON: %w", err)
	}
	counts := make(map[string]int, len(report.Findings))
	for _, finding := range report.Findings {
		counts[model.FindingIdentity(finding)]++
	}
	return Baseline{counts: counts}, nil
}

func (b Baseline) Filter(findings []model.Finding) ([]model.Finding, int) {
	remaining := make(map[string]int, len(b.counts))
	for key, count := range b.counts {
		remaining[key] = count
	}

	out := make([]model.Finding, 0, len(findings))
	suppressed := 0
	for _, finding := range findings {
		key := model.FindingIdentity(finding)
		if remaining[key] > 0 {
			remaining[key]--
			suppressed++
			continue
		}
		out = append(out, finding)
	}
	return out, suppressed
}
