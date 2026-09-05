package detectors

import (
	"bytes"
	"context"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type LineEndings struct{}

func (LineEndings) ID() string { return "text.line-endings" }

func (LineEndings) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		data, ok := readText(project.Root, rel)
		if !ok || len(data) == 0 {
			continue
		}
		crlf := bytes.Count(data, []byte("\r\n"))
		lf := bytes.Count(data, []byte("\n")) - crlf
		bareCR := bytes.Count(data, []byte("\r")) - crlf
		if (crlf > 0 && lf > 0) || bareCR > 0 {
			findings = append(findings, model.Finding{
				RuleID:      "text.mixed-line-endings",
				Title:       "Mixed line endings",
				Description: "This text file mixes newline conventions, which can cause noisy diffs and script failures across platforms.",
				Severity:    model.SeverityWarning,
				Path:        rel,
				Platforms:   []string{"windows", "macos", "linux"},
				Suggestion:  "Normalize line endings and consider a .gitattributes policy such as '* text=auto'.",
			})
		}
		if hasCRLFShebang(data) {
			findings = append(findings, model.Finding{
				RuleID:      "text.shell-crlf",
				Title:       "Shell shebang uses CRLF",
				Description: "The carriage return becomes part of the interpreter token on common Unix execution paths and can make the script fail before it starts.",
				Severity:    model.SeverityError,
				Path:        rel,
				Line:        1,
				Platforms:   []string{"linux", "macos"},
				Suggestion:  "Store executable scripts with LF line endings and enforce that policy with .gitattributes.",
			})
		}
	}
	return findings, nil
}

func hasCRLFShebang(data []byte) bool {
	if !bytes.HasPrefix(data, []byte("#!")) {
		return false
	}
	newline := bytes.IndexByte(data, '\n')
	return newline > 0 && data[newline-1] == '\r'
}
