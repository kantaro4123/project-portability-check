package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func TestVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"--version"}, &stdout, &stderr)
	if code != 0 || stdout.String() != "project-portability-check "+Version+"\n" {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestJSONCleanProject(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"--json", root}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%q output=%q", code, stderr.String(), stdout.String())
	}
	var result model.Report
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Version != Version || result.Summary.FilesScanned != 1 {
		t.Fatalf("unexpected report: %+v", result)
	}
}

func TestStrictFailsOnWarning(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "config.txt"), []byte("/Users/test/work\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"--strict", root}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d, want 1; stderr=%q output=%q", code, stderr.String(), stdout.String())
	}
}

func TestOutputModesAreExclusive(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run(context.Background(), []string{"--json", "--sarif"}, &stdout, &stderr); code != 2 {
		t.Fatalf("code=%d, want 2", code)
	}
}
