package detectors

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
)

func TestAbsolutePathDetection(t *testing.T) {
	root := t.TempDir()
	userPath := "cache=/" + "Users" + "/example/Library/cache\n"
	writeTestFile(t, root, "config.txt", userPath)
	findings, err := (AbsolutePaths{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{"config.txt"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Line != 1 {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestMixedLineEndings(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "mixed.txt", "one\r\ntwo\n")
	findings, err := (LineEndings{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{"mixed.txt"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "text.mixed-line-endings" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestShellGNUConstructs(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "build.sh", "#!/bin/sh\nreadlink -f ./thing\ndate -d tomorrow\n")
	findings, err := (ShellPortability{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{"build.sh"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 2 {
		t.Fatalf("got %d findings, want 2: %+v", len(findings), findings)
	}
}

func writeTestFile(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
