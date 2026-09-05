package detectors

import (
	"context"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type ImportCase struct{}

func (ImportCase) ID() string { return "imports.case" }

var relativeJSImportRE = regexp.MustCompile(`(?m)(?:\bfrom\s*|\brequire\(\s*|\bimport\(\s*|\bimport\s*)["'](\.{1,2}/[^"']+)["']`)

var jsModuleExt = map[string]bool{
	".js": true, ".jsx": true, ".ts": true, ".tsx": true, ".mjs": true, ".cjs": true, ".mts": true, ".cts": true,
}

var jsResolutionExts = []string{".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs", ".mts", ".cts"}

func (ImportCase) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	exact := make(map[string]bool, len(project.Files))
	folded := make(map[string][]string, len(project.Files))
	for _, rel := range project.Files {
		clean := path.Clean(rel)
		exact[clean] = true
		key := strings.ToLower(clean)
		folded[key] = append(folded[key], clean)
	}
	for key := range folded {
		sort.Strings(folded[key])
	}

	var findings []model.Finding
	for _, rel := range project.Files {
		if !jsModuleExt[strings.ToLower(path.Ext(rel))] {
			continue
		}
		data, ok := readText(project.Root, rel)
		if !ok {
			continue
		}
		for _, match := range relativeJSImportRE.FindAllSubmatchIndex(data, -1) {
			if len(match) < 4 || match[2] < 0 || match[3] < 0 {
				continue
			}
			specifier := string(data[match[2]:match[3]])
			if cut := strings.IndexAny(specifier, "?#"); cut >= 0 {
				specifier = specifier[:cut]
			}
			if specifier == "." || specifier == ".." || strings.HasSuffix(specifier, "/") {
				continue
			}

			requested, actual, mismatch := resolveImportCase(rel, specifier, exact, folded)
			if !mismatch {
				continue
			}
			findings = append(findings, model.Finding{
				RuleID:      "imports.case-mismatch",
				Title:       "Relative import path has the wrong letter case",
				Description: "The import resolves only on a case-insensitive filesystem: requested " + requested + " but the tracked path is " + actual + ".",
				Severity:    model.SeverityError,
				Path:        rel,
				Line:        lineNumber(data, match[2]),
				Platforms:   []string{"linux"},
				Suggestion:  "Change the import specifier so every path component matches the tracked filename exactly.",
			})
		}
	}
	return findings, nil
}

func resolveImportCase(source, specifier string, exact map[string]bool, folded map[string][]string) (string, string, bool) {
	base := path.Clean(path.Join(path.Dir(source), specifier))
	candidates := []string{base}
	if path.Ext(base) == "" {
		for _, ext := range jsResolutionExts {
			candidates = append(candidates, base+ext)
		}
		for _, ext := range jsResolutionExts {
			candidates = append(candidates, path.Join(base, "index"+ext))
		}
	}

	for _, candidate := range candidates {
		if exact[candidate] {
			return candidate, candidate, false
		}
	}
	for _, candidate := range candidates {
		matches := folded[strings.ToLower(candidate)]
		if len(matches) > 0 {
			return candidate, matches[0], candidate != matches[0]
		}
	}
	return "", "", false
}
