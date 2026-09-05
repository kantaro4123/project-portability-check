package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func TestBaselineIgnoresLineMovementButKeepsNewDuplicate(t *testing.T) {
	root := t.TempDir()
	baselinePath := filepath.Join(root, "baseline.json")
	report := model.Report{Findings: []model.Finding{{
		RuleID:      "paths.absolute",
		Title:       "Machine-specific absolute path",
		Description: "Found a macOS user path that may fail on another machine.",
		Severity:    model.SeverityWarning,
		Path:        "config.txt",
		Line:        4,
		Platforms:   []string{"linux", "windows"},
	}}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(baselinePath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	b, err := Load(baselinePath)
	if err != nil {
		t.Fatal(err)
	}
	current := []model.Finding{
		{RuleID: "paths.absolute", Title: "Machine-specific absolute path", Description: "Found a macOS user path that may fail on another machine.", Severity: model.SeverityWarning, Path: "config.txt", Line: 9, Platforms: []string{"linux", "windows"}},
		{RuleID: "paths.absolute", Title: "Machine-specific absolute path", Description: "Found a macOS user path that may fail on another machine.", Severity: model.SeverityWarning, Path: "config.txt", Line: 14, Platforms: []string{"linux", "windows"}},
	}
	got, suppressed := b.Filter(current)
	if suppressed != 1 || len(got) != 1 || got[0].Line != 14 {
		t.Fatalf("suppressed=%d findings=%+v", suppressed, got)
	}
}
