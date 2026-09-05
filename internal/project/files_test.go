package project

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListFilesSkipsDependencyEnvironments(t *testing.T) {
	root := t.TempDir()
	write := func(rel string) {
		t.Helper()
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("src/main.go")
	write("node_modules/pkg/index.js")
	write(".venv/lib/example.py")
	write("venv/lib/other.py")
	write(".git/config")

	files, err := ListFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"src/main.go"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("files=%v, want %v", files, want)
	}
}
