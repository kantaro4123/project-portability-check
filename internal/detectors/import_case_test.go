package detectors

import (
	"context"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
)

func TestImportCaseMismatchIsReported(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "src/main.ts", `import { value } from "./Utils";\n`)
	writeTestFile(t, root, "src/utils.ts", "export const value = 1;\n")
	project := analyzer.Project{Root: root, Files: []string{"src/main.ts", "src/utils.ts"}}

	findings, err := (ImportCase{}).Detect(context.Background(), project)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "imports.case-mismatch" || findings[0].Line != 1 {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestImportCaseExactPathIsClean(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "src/main.ts", `import { value } from "./utils";\n`)
	writeTestFile(t, root, "src/utils.ts", "export const value = 1;\n")
	project := analyzer.Project{Root: root, Files: []string{"src/main.ts", "src/utils.ts"}}

	findings, err := (ImportCase{}).Detect(context.Background(), project)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestImportCaseResolvesIndexFiles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "src/main.ts", `import thing from "./Widget";\n`)
	writeTestFile(t, root, "src/widget/index.ts", "export default 1;\n")
	project := analyzer.Project{Root: root, Files: []string{"src/main.ts", "src/widget/index.ts"}}

	findings, err := (ImportCase{}).Detect(context.Background(), project)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].RuleID != "imports.case-mismatch" {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func TestImportCaseIgnoresPackageImports(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "src/main.ts", `import React from "React";\n`)
	project := analyzer.Project{Root: root, Files: []string{"src/main.ts"}}

	findings, err := (ImportCase{}).Detect(context.Background(), project)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}
