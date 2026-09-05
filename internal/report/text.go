package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func WriteText(w io.Writer, report model.Report) error {
	if _, err := fmt.Fprintf(w, "project-portability-check v%s\nProject: %s\nScanned: %d file(s)\n", report.Version, report.Root, report.Summary.FilesScanned); err != nil {
		return err
	}
	if len(report.TargetPlatforms) > 0 {
		if _, err := fmt.Fprintf(w, "Targets: %s\n", strings.Join(report.TargetPlatforms, ", ")); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if len(report.Findings) == 0 {
		if _, err := fmt.Fprintln(w, "No portability findings."); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintln(w, "Findings"); err != nil {
			return err
		}
		for _, f := range report.Findings {
			marker := "·"
			switch f.Severity {
			case model.SeverityError:
				marker = "✗"
			case model.SeverityWarning:
				marker = "!"
			}
			location := f.Path
			if f.Line > 0 {
				location = fmt.Sprintf("%s:%d", f.Path, f.Line)
			}
			if location != "" {
				location = " (" + location + ")"
			}
			if _, err := fmt.Fprintf(w, "  %s %s%s [%s]\n    %s\n", marker, f.Title, location, f.RuleID, f.Description); err != nil {
				return err
			}
			if len(f.Platforms) > 0 {
				if _, err := fmt.Fprintf(w, "    Affects: %s\n", strings.Join(f.Platforms, ", ")); err != nil {
					return err
				}
			}
			if f.Suggestion != "" {
				if _, err := fmt.Fprintf(w, "    Fix: %s\n", f.Suggestion); err != nil {
					return err
				}
			}
		}
	}
	_, err := fmt.Fprintf(w, "\nPortability Score: %d/100\nErrors: %d  Warnings: %d  Info: %d\n", report.Summary.Score, report.Summary.Errors, report.Summary.Warnings, report.Summary.Info)
	return err
}
