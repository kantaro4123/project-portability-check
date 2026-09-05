package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func TestLoadAndFilter(t *testing.T) {
	root := t.TempDir()
	data := `{"ignore_rules":["paths.absolute"],"ignore_paths":["vendor/*"],"target_platforms":["windows"]}`
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
		{RuleID: "mac-only", Path: "src/mac.go", Platforms: []string{"macos"}},
		{RuleID: "win", Path: "src/win.go", Platforms: []string{"windows"}},
		{RuleID: "general", Path: "src/general.go"},
	}
	got := cfg.Filter(findings)
	if len(got) != 2 || got[0].RuleID != "win" || got[1].RuleID != "general" {
		t.Fatalf("unexpected filtered findings: %+v", got)
	}
}

func TestParseTargetsNormalizesAndDeduplicates(t *testing.T) {
	got, err := ParseTargets("Windows, linux,WINDOWS")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "linux" || got[1] != "windows" {
		t.Fatalf("unexpected targets: %+v", got)
	}
}

func TestInvalidTargetIsRejected(t *testing.T) {
	if _, err := ParseTargets("freebsd"); err == nil {
		t.Fatal("expected unsupported target to fail")
	}
}

func TestMissingConfigIsAllowed(t *testing.T) {
	cfg, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.IgnoreRules) != 0 || len(cfg.IgnorePaths) != 0 || len(cfg.TargetPlatforms) != 0 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
