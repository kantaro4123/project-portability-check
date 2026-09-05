package detectors

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const maxTextBytes = 2 << 20

func readText(root, rel string) ([]byte, bool) {
	path := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxTextBytes {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil || bytes.IndexByte(data, 0) >= 0 || !utf8.Valid(data) {
		return nil, false
	}
	return data, true
}

func lineNumber(data []byte, offset int) int {
	if offset < 0 {
		return 0
	}
	return bytes.Count(data[:offset], []byte{'\n'}) + 1
}

func normalizedFileSet(files []string) map[string]bool {
	set := make(map[string]bool, len(files))
	for _, rel := range files {
		set[strings.ToLower(path.Clean(rel))] = true
	}
	return set
}

func hasInAncestors(files map[string]bool, dir string, names ...string) bool {
	dir = strings.ToLower(path.Clean(dir))
	for {
		for _, name := range names {
			candidate := strings.ToLower(name)
			if dir != "." && dir != "" {
				candidate = path.Join(dir, candidate)
			}
			if files[candidate] {
				return true
			}
		}
		if dir == "." || dir == "" || dir == "/" {
			return false
		}
		parent := path.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
}
