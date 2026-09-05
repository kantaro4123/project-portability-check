package detectors

import (
	"context"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
)

func TestPackageScriptUnixCommand(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "package.json", `{"scripts":{"clean":"rm -rf dist"}}`)
	findings, err := (PackageScripts{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{"package.json"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "package.script-unix" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestNodeRuntimeAndLockfileWarnings(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "package.json", `{}`)
	project := analyzer.Project{Root: root, Files: []string{"package.json"}}
	runtimeFindings, err := (RuntimePins{}).Detect(context.Background(), project)
	if err != nil {
		t.Fatal(err)
	}
	lockFindings, err := (Lockfiles{}).Detect(context.Background(), project)
	if err != nil {
		t.Fatal(err)
	}
	if len(runtimeFindings) != 1 || len(lockFindings) != 1 {
		t.Fatalf("runtime=%+v lock=%+v", runtimeFindings, lockFindings)
	}
}

func TestDockerFixedPlatform(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "Dockerfile", "FROM --platform=linux/amd64 alpine:latest\n")
	findings, err := (DockerPlatform{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{"Dockerfile"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "docker.fixed-platform" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestFullCIMatrixHasNoFinding(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, ".github/workflows/ci.yml", "runs-on: ubuntu-latest\n# macos-latest windows-latest\n")
	findings, err := (CIMatrix{}).Detect(context.Background(), analyzer.Project{Root: root, Files: []string{".github/workflows/ci.yml"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}
