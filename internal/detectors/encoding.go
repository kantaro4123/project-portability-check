package detectors

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type TextEncoding struct{}

func (TextEncoding) ID() string { return "text.encoding" }

var likelyTextExt = map[string]bool{
	".go": true, ".py": true, ".js": true, ".ts": true, ".tsx": true, ".jsx": true,
	".json": true, ".yaml": true, ".yml": true, ".toml": true, ".xml": true,
	".md": true, ".txt": true, ".sh": true, ".bash": true, ".zsh": true, ".ps1": true,
	".css": true, ".html": true, ".java": true, ".rs": true, ".c": true, ".h": true,
}

func (TextEncoding) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		if !likelyTextExt[strings.ToLower(filepath.Ext(rel))] {
			continue
		}
		full := filepath.Join(project.Root, filepath.FromSlash(rel))
		info, err := os.Lstat(full)
		if err != nil || !info.Mode().IsRegular() || info.Size() > maxTextBytes {
			continue
		}
		data, err := os.ReadFile(full)
		if err != nil || bytes.IndexByte(data, 0) >= 0 {
			continue
		}
		if !utf8.Valid(data) {
			findings = append(findings, model.Finding{RuleID: "text.non-utf8", Title: "Non-UTF-8 source text", Description: "A likely text/source file is not valid UTF-8, which can be decoded differently across tools and locales.", Severity: model.SeverityWarning, Path: rel, Suggestion: "Convert the file to UTF-8 or document the required encoding."})
			continue
		}
		if bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
			findings = append(findings, model.Finding{RuleID: "text.utf8-bom", Title: "UTF-8 BOM present", Description: "Some Unix tools and shebang handling can be sensitive to a UTF-8 BOM.", Severity: model.SeverityInfo, Path: rel, Platforms: []string{"macos", "linux"}, Suggestion: "Prefer UTF-8 without BOM unless a consuming tool requires it."})
		}
	}
	return findings, nil
}
