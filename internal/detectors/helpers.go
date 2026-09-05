package detectors

import (
	"bytes"
	"os"
	"path/filepath"
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
