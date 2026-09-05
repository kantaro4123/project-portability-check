package detectors

import (
	"context"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type Lockfiles struct{}

func (Lockfiles) ID() string { return "deps.lockfile" }

func (Lockfiles) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	files := make(map[string]bool, len(project.Files))
	for _, rel := range project.Files {
		files[strings.ToLower(rel)] = true
	}
	var findings []model.Finding
	if files["package.json"] && !any(files, "package-lock.json", "npm-shrinkwrap.json", "yarn.lock", "pnpm-lock.yaml", "bun.lock", "bun.lockb") {
		findings = append(findings, model.Finding{RuleID: "deps.node-lockfile", Title: "JavaScript dependencies are not locked", Description: "package.json exists without a recognized lockfile, so installs can resolve different dependency versions.", Severity: model.SeverityWarning, Path: "package.json", Suggestion: "Commit the lockfile produced by the package manager used by this project."})
	}
	if files["cargo.toml"] && !files["cargo.lock"] {
		findings = append(findings, model.Finding{RuleID: "deps.cargo-lockfile", Title: "Cargo dependencies are not locked", Description: "Cargo.toml exists without Cargo.lock.", Severity: model.SeverityInfo, Path: "Cargo.toml", Suggestion: "For applications and CLI tools, consider committing Cargo.lock for reproducible builds."})
	}
	return findings, nil
}
