package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func TestSARIFOutput(t *testing.T) {
	input := model.Report{Version: "0.1.0", Findings: []model.Finding{{RuleID: "paths.absolute", Title: "Absolute path", Description: "example", Severity: model.SeverityWarning, Path: "config.txt", Line: 3}}}
	var buf bytes.Buffer
	if err := WriteSARIF(&buf, input); err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if doc["version"] != "2.1.0" {
		t.Fatalf("unexpected SARIF: %s", buf.String())
	}
}
