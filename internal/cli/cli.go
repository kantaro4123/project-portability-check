package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/baseline"
	"github.com/kantaro4123/project-portability-check/internal/config"
	"github.com/kantaro4123/project-portability-check/internal/detectors"
	"github.com/kantaro4123/project-portability-check/internal/model"
	projectfiles "github.com/kantaro4123/project-portability-check/internal/project"
	"github.com/kantaro4123/project-portability-check/internal/report"
)

const Version = "0.2.0"

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("project-portability-check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOutput := fs.Bool("json", false, "emit JSON")
	sarifOutput := fs.Bool("sarif", false, "emit SARIF 2.1.0")
	strict := fs.Bool("strict", false, "fail on warnings as well as errors")
	target := fs.String("target", "", "comma-separated target platforms: linux, macos, windows")
	baselineFile := fs.String("baseline", "", "suppress findings already present in a previous JSON report")
	listRules := fs.Bool("list-rules", false, "list built-in finding rule IDs")
	showVersion := fs.Bool("version", false, "print version")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: project-portability-check [options] [path]")
		fmt.Fprintln(stderr, "\nAnalyze a project for cross-platform portability hazards.")
		fmt.Fprintln(stderr, "\nOptions:")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		fmt.Fprintf(stdout, "project-portability-check %s\n", Version)
		return 0
	}
	if *listRules {
		for _, rule := range detectors.Rules() {
			fmt.Fprintln(stdout, rule.ID)
		}
		return 0
	}
	if *jsonOutput && *sarifOutput {
		fmt.Fprintln(stderr, "error: --json and --sarif are mutually exclusive")
		return 2
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(stderr, "error: expected at most one project path")
		return 2
	}
	root := "."
	if fs.NArg() == 1 {
		root = fs.Arg(0)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(stderr, "error: resolve project path: %v\n", err)
		return 2
	}
	info, err := os.Stat(absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "error: open project: %v\n", err)
		return 2
	}
	if !info.IsDir() {
		fmt.Fprintln(stderr, "error: project path must be a directory")
		return 2
	}
	cfg, err := config.Load(absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "error: load configuration: %v\n", err)
		return 2
	}
	if *target != "" {
		targets, targetErr := config.ParseTargets(*target)
		if targetErr != nil {
			fmt.Fprintf(stderr, "error: %v\n", targetErr)
			return 2
		}
		cfg.TargetPlatforms = targets
	}

	resolvedBaseline := ""
	if *baselineFile != "" {
		resolvedBaseline = *baselineFile
		if !filepath.IsAbs(resolvedBaseline) {
			resolvedBaseline = filepath.Join(absRoot, filepath.FromSlash(resolvedBaseline))
		}
	}

	files, err := projectfiles.ListFiles(absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "error: scan project: %v\n", err)
		return 2
	}
	if resolvedBaseline != "" {
		files = excludeProjectFile(files, absRoot, resolvedBaseline)
	}

	engine := analyzer.New(detectors.Default()...)
	findings, err := engine.Analyze(ctx, analyzer.Project{Root: absRoot, Files: files})
	if err != nil {
		fmt.Fprintf(stderr, "error: analyze project: %v\n", err)
		return 2
	}
	findings = cfg.Filter(findings)

	baselineSuppressed := 0
	if resolvedBaseline != "" {
		previous, baselineErr := baseline.Load(resolvedBaseline)
		if baselineErr != nil {
			fmt.Fprintf(stderr, "error: load baseline: %v\n", baselineErr)
			return 2
		}
		findings, baselineSuppressed = previous.Filter(findings)
	}

	result := model.Report{
		Version:            Version,
		Root:               absRoot,
		TargetPlatforms:    cfg.TargetPlatforms,
		BaselineSuppressed: baselineSuppressed,
		Findings:           findings,
		Summary:            report.Summarize(len(files), findings),
	}
	if *jsonOutput {
		err = report.WriteJSON(stdout, result)
	} else if *sarifOutput {
		err = report.WriteSARIF(stdout, result)
	} else {
		err = report.WriteText(stdout, result)
	}
	if err != nil && !errors.Is(err, io.ErrClosedPipe) {
		fmt.Fprintf(stderr, "error: write report: %v\n", err)
		return 2
	}
	if result.Summary.Errors > 0 || (*strict && result.Summary.Warnings > 0) {
		return 1
	}
	return 0
}

func excludeProjectFile(files []string, root, filename string) []string {
	rel, err := filepath.Rel(root, filename)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return files
	}
	rel = filepath.ToSlash(rel)
	out := files[:0]
	for _, file := range files {
		if file != rel {
			out = append(out, file)
		}
	}
	return out
}
