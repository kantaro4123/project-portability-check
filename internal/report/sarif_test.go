package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func TestSARIFOutput(t *testing.T) {
	input := model.Report{Version: "0.2.0", Findings: []model.Finding{{RuleID: "paths.absolute", Title: "Absolute path", Description: "example", Severity: model.SeverityWarning, Path: "config.txt", Line: 3, Suggestion: "Use a relative path."}}}
	var buf bytes.Buffer
	if err := WriteSARIF(&buf, input); err != nil {
		t.Fatal(err)
	}
	var doc sarifDocument
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Version != "2.1.0" || len(doc.Runs) != 1 {
		t.Fatalf("unexpected SARIF: %s", buf.String())
	}
	driver := doc.Runs[0].Tool.Driver
	if len(driver.Rules) != 1 || driver.Rules[0].ID != "paths.absolute" || driver.Rules[0].Help.Text == "" {
		t.Fatalf("missing rule metadata: %+v", driver.Rules)
	}
	result := doc.Runs[0].Results[0]
	if result.PartialFingerprints["projectPortabilityCheck/v1"] == "" {
		t.Fatalf("missing fingerprint: %+v", result)
	}
}

func TestFindingFingerprintIsDeterministic(t *testing.T) {
	finding := model.Finding{RuleID: "x", Title: "title", Description: "desc", Path: `src\\file.go`, Line: 7}
	first := model.FindingFingerprint(finding)
	second := model.FindingFingerprint(finding)
	if first == "" || first != second {
		t.Fatalf("fingerprint is not deterministic: %q %q", first, second)
	}
}
