package detectors

import (
	"context"
	"path"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type Lockfiles struct{}

func (Lockfiles) ID() string { return "deps.lockfile" }

func (Lockfiles) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	files := normalizedFileSet(project.Files)
	var findings []model.Finding

	for _, rel := range project.Files {
		lower := strings.ToLower(path.Clean(rel))
		base := path.Base(lower)
		dir := path.Dir(lower)

		switch base {
		case "package.json":
			if !hasInAncestors(files, dir, "package-lock.json", "npm-shrinkwrap.json", "yarn.lock", "pnpm-lock.yaml", "bun.lock", "bun.lockb") {
				findings = append(findings, model.Finding{
					RuleID:      "deps.node-lockfile",
					Title:       "JavaScript dependencies are not locked",
					Description: "package.json exists without a recognized lockfile in this package or a parent workspace, so installs can resolve different dependency versions.",
					Severity:    model.SeverityWarning,
					Path:        rel,
					Suggestion:  "Commit the lockfile produced by the package manager used by this package or workspace.",
				})
			}
		case "cargo.toml":
			if !hasInAncestors(files, dir, "cargo.lock") {
				findings = append(findings, model.Finding{
					RuleID:      "deps.cargo-lockfile",
					Title:       "Cargo dependencies are not locked",
					Description: "Cargo.toml exists without Cargo.lock in this crate or a parent workspace.",
					Severity:    model.SeverityInfo,
					Path:        rel,
					Suggestion:  "For applications and CLI tools, consider committing Cargo.lock for reproducible builds.",
				})
			}
		}
	}

	return findings, nil
}
