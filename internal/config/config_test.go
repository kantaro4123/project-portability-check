package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func TestLoadAndFilter(t *testing.T) {
	root := t.TempDir()
	data := `{"ignore_rules":["paths.absolute"],"ignore_paths":["vendor/*"]}`
	if err := os.WriteFile(filepath.Join(root, Filename), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	findings := []model.Finding{
		{RuleID: "paths.absolute", Path: "src/main.go"},
		{RuleID: "other", Path: "vendor/generated.go"},
		{RuleID: "other", Path: "src/keep.go"},
	}
	got := cfg.Filter(findings)
	if len(got) != 1 || got[0].Path != "src/keep.go" {
		t.Fatalf("unexpected filtered findings: %+v", got)
	}
}

func TestMissingConfigIsAllowed(t *testing.T) {
	cfg, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.IgnoreRules) != 0 || len(cfg.IgnorePaths) != 0 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
