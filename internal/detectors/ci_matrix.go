package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type CIMatrix struct{}

func (CIMatrix) ID() string { return "ci.platform-coverage" }

var ciOSPatterns = map[string]*regexp.Regexp{
	"linux":   regexp.MustCompile(`(?i)ubuntu-|linux`),
	"macos":   regexp.MustCompile(`(?i)macos-`),
	"windows": regexp.MustCompile(`(?i)windows-`),
}

func (CIMatrix) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	seenWorkflow := false
	covered := map[string]bool{}
	for _, rel := range project.Files {
		lower := strings.ToLower(rel)
		if !strings.HasPrefix(lower, ".github/workflows/") || !(strings.HasSuffix(lower, ".yml") || strings.HasSuffix(lower, ".yaml")) {
			continue
		}
		seenWorkflow = true
		data, ok := readText(project.Root, rel)
		if !ok {
			continue
		}
		data = stripYAMLComments(data)
		for osName, re := range ciOSPatterns {
			if re.Match(data) {
				covered[osName] = true
			}
		}
	}
	if !seenWorkflow || len(covered) >= 3 {
		return nil, nil
	}
	missing := make([]string, 0, 3-len(covered))
	for _, osName := range []string{"linux", "macos", "windows"} {
		if !covered[osName] {
			missing = append(missing, osName)
		}
	}
	return []model.Finding{{RuleID: "ci.platform-coverage", Title: "CI does not cover all major desktop platforms", Description: "GitHub Actions workflows were found, but no runner reference was detected for: " + strings.Join(missing, ", ") + ".", Severity: model.SeverityInfo, Platforms: missing, Suggestion: "Add representative CI jobs where cross-platform support is a project goal."}}, nil
}

func stripYAMLComments(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		lines[i] = stripYAMLComment(line)
	}
	return []byte(strings.Join(lines, "\n"))
}

func stripYAMLComment(line string) string {
	inSingle := false
	inDouble := false
	escaped := false
	for i := 0; i < len(line); i++ {
		ch := line[i]
		if inDouble && ch == '\\' && !escaped {
			escaped = true
			continue
		}
		if ch == '"' && !inSingle && !escaped {
			inDouble = !inDouble
		} else if ch == '\'' && !inDouble {
			inSingle = !inSingle
		} else if ch == '#' && !inSingle && !inDouble && (i == 0 || line[i-1] == ' ' || line[i-1] == '\t') {
			return line[:i]
		}
		escaped = false
	}
	return line
}
