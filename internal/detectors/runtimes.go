package detectors

import (
	"context"
	"path"
	"sort"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type RuntimePins struct{}

func (RuntimePins) ID() string { return "runtime.pin" }

func (RuntimePins) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	files := normalizedFileSet(project.Files)
	var findings []model.Finding
	pythonProjects := map[string]string{}

	for _, rel := range project.Files {
		lower := strings.ToLower(path.Clean(rel))
		base := path.Base(lower)
		dir := path.Dir(lower)

		switch base {
		case "package.json":
			if !hasInAncestors(files, dir, ".nvmrc", ".node-version", ".tool-versions", "mise.toml") {
				findings = append(findings, model.Finding{
					RuleID:      "runtime.node-unpinned",
					Title:       "Node.js version is not pinned",
					Description: "Different Node.js versions can change dependency resolution and runtime behavior.",
					Severity:    model.SeverityWarning,
					Path:        rel,
					Suggestion:  "Add .nvmrc, .node-version, mise.toml, or another explicit runtime pin in this package or a parent workspace.",
				})
			}
		case "pyproject.toml", "requirements.txt", "setup.py", "setup.cfg", "pipfile":
			if _, exists := pythonProjects[dir]; !exists {
				pythonProjects[dir] = rel
			}
		case "go.mod":
			data, ok := readText(project.Root, rel)
			if ok && !hasGoDirective(data) {
				findings = append(findings, model.Finding{
					RuleID:      "runtime.go-unpinned",
					Title:       "Go language version is not declared",
					Description: "go.mod does not declare a Go language version.",
					Severity:    model.SeverityWarning,
					Path:        rel,
					Suggestion:  "Add a go directive to go.mod.",
				})
			}
		}
	}

	dirs := make([]string, 0, len(pythonProjects))
	for dir := range pythonProjects {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	for _, dir := range dirs {
		if hasInAncestors(files, dir, ".python-version", ".tool-versions", "mise.toml") {
			continue
		}
		findings = append(findings, model.Finding{
			RuleID:      "runtime.python-unpinned",
			Title:       "Python version is not pinned",
			Description: "A broad Python requirement can hide interpreter-specific behavior.",
			Severity:    model.SeverityWarning,
			Path:        pythonProjects[dir],
			Suggestion:  "Pin a development interpreter in this project or a parent workspace with .python-version, mise, or your environment manager.",
		})
	}

	return findings, nil
}

func hasGoDirective(data []byte) bool {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") && len(strings.Fields(line)) >= 2 {
			return true
		}
	}
	return false
}
