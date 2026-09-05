package model

import "testing"

func TestFindingIdentityIgnoresLineMovement(t *testing.T) {
	first := Finding{RuleID: "paths.absolute", Title: "Absolute path", Description: "example", Severity: SeverityWarning, Path: `src\config.go`, Line: 4, Platforms: []string{"windows", "linux"}}
	second := first
	second.Path = "src/config.go"
	second.Line = 20
	second.Platforms = []string{"LINUX", "WINDOWS"}

	if FindingIdentity(first) != FindingIdentity(second) {
		t.Fatal("semantic identity changed after line movement or path normalization")
	}
	if FindingFingerprint(first) == FindingFingerprint(second) {
		t.Fatal("exact fingerprint should change when the source line changes")
	}
}

func TestAttachFingerprintsPreservesExistingValue(t *testing.T) {
	findings := []Finding{{RuleID: "x", Fingerprint: "external"}, {RuleID: "y", Title: "Y"}}
	AttachFingerprints(findings)
	if findings[0].Fingerprint != "external" {
		t.Fatalf("existing fingerprint overwritten: %q", findings[0].Fingerprint)
	}
	if findings[1].Fingerprint == "" {
		t.Fatal("missing generated fingerprint")
	}
}
