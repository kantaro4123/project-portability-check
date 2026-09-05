package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/config"
	"github.com/kantaro4123/project-portability-check/internal/detectors"
	"github.com/kantaro4123/project-portability-check/internal/model"
	projectfiles "github.com/kantaro4123/project-portability-check/internal/project"
	"github.com/kantaro4123/project-portability-check/internal/report"
)

const Version = "0.1.0"

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("project-portability-check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOutput := fs.Bool("json", false, "emit JSON")
	sarifOutput := fs.Bool("sarif", false, "emit SARIF 2.1.0")
	strict := fs.Bool("strict", false, "fail on warnings as well as errors")
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
	files, err := projectfiles.ListFiles(absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "error: scan project: %v\n", err)
		return 2
	}
	cfg, err := config.Load(absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "error: load configuration: %v\n", err)
		return 2
	}
	engine := analyzer.New(detectors.Default()...)
	findings, err := engine.Analyze(ctx, analyzer.Project{Root: absRoot, Files: files})
	if err != nil {
		fmt.Fprintf(stderr, "error: analyze project: %v\n", err)
		return 2
	}
	findings = cfg.Filter(findings)
	result := model.Report{Version: Version, Root: absRoot, Findings: findings, Summary: report.Summarize(len(files), findings)}
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
