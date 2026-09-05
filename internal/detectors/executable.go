package detectors

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type ExecutableScripts struct{}

func (ExecutableScripts) ID() string { return "fs.executable" }

func (ExecutableScripts) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		data, ok := readText(project.Root, rel)
		if !ok || !bytes.HasPrefix(data, []byte("#!")) {
			continue
		}
		info, err := os.Stat(filepath.Join(project.Root, filepath.FromSlash(rel)))
		if err != nil {
			continue
		}
		if info.Mode().Perm()&0111 == 0 {
			findings = append(findings, model.Finding{
				RuleID:      "fs.script-not-executable",
				Title:       "Script lacks executable permission",
				Description: "A shebang script is not executable, so './script' can fail on Unix-like systems.",
				Severity:    model.SeverityWarning,
				Path:        rel,
				Platforms:   []string{"macos", "linux"},
				Suggestion:  "Mark the script executable in Git or document that it must be invoked through its interpreter.",
			})
		}
	}
	return findings, nil
}
