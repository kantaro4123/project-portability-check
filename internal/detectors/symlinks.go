package detectors

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type Symlinks struct{}

func (Symlinks) ID() string { return "fs.symlink" }

func (Symlinks) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	rootAbs, err := filepath.Abs(project.Root)
	if err != nil {
		return nil, err
	}
	var findings []model.Finding
	for _, rel := range project.Files {
		full := filepath.Join(project.Root, filepath.FromSlash(rel))
		info, err := os.Lstat(full)
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			continue
		}
		target, err := os.Readlink(full)
		if err != nil {
			continue
		}
		resolved := target
		if !filepath.IsAbs(target) {
			resolved = filepath.Join(filepath.Dir(full), target)
		}
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			continue
		}
		relToRoot, err := filepath.Rel(rootAbs, resolved)
		outside := err != nil || relToRoot == ".." || strings.HasPrefix(relToRoot, ".."+string(filepath.Separator))
		severity := model.SeverityWarning
		title := "Tracked symbolic link"
		description := "Symbolic links can behave differently across filesystems and Windows developer-mode settings."
		if outside || filepath.IsAbs(target) {
			severity = model.SeverityError
			title = "Non-portable symbolic link"
			description = "The symbolic link points outside the project or uses an absolute target."
		}
		findings = append(findings, model.Finding{RuleID: "fs.symlink", Title: title, Description: description, Severity: severity, Path: rel, Platforms: []string{"windows"}, Suggestion: "Prefer project-relative files or document required symlink support."})
	}
	return findings, nil
}
