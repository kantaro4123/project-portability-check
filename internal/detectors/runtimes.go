package detectors

import (
	"context"
	"path"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type RuntimePins struct{}

func (RuntimePins) ID() string { return "runtime.pin" }

func (RuntimePins) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	files := make(map[string]bool, len(project.Files))
	for _, rel := range project.Files {
		files[strings.ToLower(rel)] = true
	}
	var findings []model.Finding
	if files["package.json"] && !any(files, ".nvmrc", ".node-version", ".tool-versions", "mise.toml") {
		findings = append(findings, model.Finding{RuleID: "runtime.node-unpinned", Title: "Node.js version is not pinned", Description: "Different Node.js versions can change dependency resolution and runtime behavior.", Severity: model.SeverityWarning, Path: "package.json", Suggestion: "Add .nvmrc, .node-version, mise.toml, or another explicit runtime pin."})
	}
	pythonProject := files["pyproject.toml"] || files["requirements.txt"] || files["setup.py"] || files["setup.cfg"]
	if pythonProject && !any(files, ".python-version", ".tool-versions", "mise.toml") {
		findings = append(findings, model.Finding{RuleID: "runtime.python-unpinned", Title: "Python version is not pinned", Description: "A broad Python requirement can hide interpreter-specific behavior.", Severity: model.SeverityWarning, Suggestion: "Pin a development interpreter with .python-version, mise, or your environment manager."})
	}
	for rel := range files {
		if path.Base(rel) == "go.mod" {
			data, ok := readText(project.Root, rel)
			if ok && !strings.Contains(string(data), "\ngo ") && !strings.HasPrefix(string(data), "go ") {
				findings = append(findings, model.Finding{RuleID: "runtime.go-unpinned", Title: "Go language version is not declared", Description: "go.mod does not declare a Go language version.", Severity: model.SeverityWarning, Path: rel, Suggestion: "Add a go directive to go.mod."})
			}
			break
		}
	}
	return findings, nil
}

func any(files map[string]bool, names ...string) bool {
	for _, name := range names {
		if files[strings.ToLower(name)] {
			return true
		}
	}
	return false
}
