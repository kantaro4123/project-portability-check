package detectors

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type ShellPortability struct{}

func (ShellPortability) ID() string { return "shell.portability" }

var shellRules = []struct {
	id, title, description, suggestion string
	platforms                          []string
	re                                 *regexp.Regexp
}{
	{"shell.grep-p", "GNU grep -P usage", "grep -P is not supported by the BSD grep shipped with macOS.", "Use portable grep/awk, or document GNU grep as a dependency.", []string{"macos"}, regexp.MustCompile(`(?m)(^|[;&|]\s*)grep\s+[^\n]*-P`)},
	{"shell.sed-i", "Non-portable sed -i usage", "GNU and BSD sed use different syntax for in-place editing.", "Avoid sed -i or branch on the platform with an explicit backup suffix.", []string{"macos", "linux"}, regexp.MustCompile(`(?m)(^|[;&|]\s*)sed\s+[^\n]*-i(?:\s|$)`)},
	{"shell.readlink-f", "GNU readlink -f usage", "readlink -f is not available in the default macOS readlink.", "Use a language runtime for path resolution or provide a portable fallback.", []string{"macos"}, regexp.MustCompile(`(?m)\breadlink\s+-f\b`)},
	{"shell.date-d", "GNU date -d usage", "date -d is not supported by BSD date on macOS.", "Use a portable runtime for date parsing or handle BSD date separately.", []string{"macos"}, regexp.MustCompile(`(?m)\bdate\s+[^\n]*-d(?:\s|$)`)},
	{"shell.xargs-r", "GNU xargs -r usage", "xargs -r is not supported by the BSD xargs shipped with macOS.", "Avoid -r and guard empty input explicitly.", []string{"macos"}, regexp.MustCompile(`(?m)\bxargs\s+[^\n]*-r(?:\s|$)`)},
}

func (ShellPortability) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		ext := strings.ToLower(filepath.Ext(rel))
		if ext != ".sh" && ext != ".bash" && ext != ".zsh" && ext != "" {
			continue
		}
		data, ok := readText(project.Root, rel)
		if !ok {
			continue
		}
		for _, rule := range shellRules {
			for _, loc := range rule.re.FindAllIndex(data, -1) {
				findings = append(findings, model.Finding{RuleID: rule.id, Title: rule.title, Description: rule.description, Severity: model.SeverityWarning, Path: rel, Line: lineNumber(data, loc[0]), Platforms: rule.platforms, Suggestion: rule.suggestion})
			}
		}
	}
	return findings, nil
}
