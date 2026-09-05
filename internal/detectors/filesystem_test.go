package detectors

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
)

func TestCaseCollision(t *testing.T) {
	findings, err := (CaseCollisions{}).Detect(context.Background(), analyzer.Project{Files: []string{"README.md", "readme.md"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "fs.case-collision" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestWindowsReservedName(t *testing.T) {
	findings, err := (WindowsNames{}).Detect(context.Background(), analyzer.Project{Files: []string{"docs/CON.txt"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "fs.windows-reserved" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestExternalSymlink(t *testing.T) {
	if os.Getenv("GOOS") == "windows" {
		t.Skip("symlink creation may require Windows developer mode")
	}
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "external")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	findings, err := (Symlinks{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{"external"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Severity != "error" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}
