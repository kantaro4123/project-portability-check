package detectors

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
)

func TestWindowsForbiddenCharacter(t *testing.T) {
	findings, err := (WindowsPaths{}).Detect(context.Background(), analyzer.Project{Files: []string{"docs/what?.md"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "fs.windows-forbidden-char" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestWindowsLongPath(t *testing.T) {
	rel := strings.Repeat("a", 241) + ".txt"
	findings, err := (WindowsPaths{}).Detect(context.Background(), analyzer.Project{Files: []string{rel}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "fs.windows-long-path" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestELFNativeBinary(t *testing.T) {
	root := t.TempDir()
	full := filepath.Join(root, "tool")
	if err := os.WriteFile(full, []byte{0x7f, 'E', 'L', 'F', 2, 1, 1}, 0o755); err != nil {
		t.Fatal(err)
	}
	findings, err := (NativeBinaries{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{"tool"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "binary.native" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}
