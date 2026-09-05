package report

import (
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func TestSummarize(t *testing.T) {
	findings := []model.Finding{
		{Severity: model.SeverityError},
		{Severity: model.SeverityWarning},
		{Severity: model.SeverityInfo},
	}
	got := Summarize(12, findings)
	if got.FilesScanned != 12 || got.Errors != 1 || got.Warnings != 1 || got.Info != 1 || got.Score != 80 {
		t.Fatalf("unexpected summary: %+v", got)
	}
}

func TestSummarizeDoesNotGoBelowZero(t *testing.T) {
	findings := make([]model.Finding, 10)
	for i := range findings {
		findings[i].Severity = model.SeverityError
	}
	if got := Summarize(1, findings).Score; got != 0 {
		t.Fatalf("score = %d, want 0", got)
	}
}
