package detectors

import (
	"context"
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type PackageScripts struct{}

func (PackageScripts) ID() string { return "package.scripts" }

var unixScriptPatterns = []struct {
	name string
	re   *regexp.Regexp
}{
	{"rm", regexp.MustCompile(`(^|[;&|]\s*)rm\s+-`)},
	{"cp", regexp.MustCompile(`(^|[;&|]\s*)cp\s+`)},
	{"mv", regexp.MustCompile(`(^|[;&|]\s*)mv\s+`)},
	{"mkdir -p", regexp.MustCompile(`(^|[;&|]\s*)mkdir\s+-p\b`)},
	{"export", regexp.MustCompile(`(^|[;&|]\s*)export\s+[A-Za-z_]`)},
	{"Unix environment assignment", regexp.MustCompile(`(^|[;&|]\s*)[A-Za-z_][A-Za-z0-9_]*=[^ ]+\s+[^;&|]+`)},
	{"/dev/null", regexp.MustCompile(`/dev/null`)},
}

func (PackageScripts) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	data, ok := readText(project.Root, "package.json")
	if !ok {
		return nil, nil
	}
	var manifest struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, nil
	}
	names := make([]string, 0, len(manifest.Scripts))
	for name := range manifest.Scripts {
		names = append(names, name)
	}
	sort.Strings(names)
	var findings []model.Finding
	for _, name := range names {
		command := manifest.Scripts[name]
		for _, pattern := range unixScriptPatterns {
			if !pattern.re.MatchString(command) {
				continue
			}
			findings = append(findings, model.Finding{
				RuleID:      "package.script-unix",
				Title:       "Unix-specific package script",
				Description: "package.json script " + name + " uses " + pattern.name + ", which is not portable to the default Windows shell.",
				Severity:    model.SeverityWarning,
				Path:        "package.json",
				Platforms:   []string{"windows"},
				Suggestion:  "Use a cross-platform Node.js helper or a portable package such as rimraf/cross-env where appropriate.",
			})
			break
		}
	}
	_ = strings.Builder{}
	return findings, nil
}
