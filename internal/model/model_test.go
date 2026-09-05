package model

import "testing"

func TestFindingIdentityIgnoresLineAndMessageMovement(t *testing.T) {
	first := Finding{RuleID: "paths.absolute", Title: "Absolute path", Description: "old wording", Severity: SeverityWarning, Path: `src\config.go`, Line: 4, Platforms: []string{"windows", "linux"}}
	second := first
	second.Path = "src/config.go"
	second.Line = 20
	second.Title = "Machine-specific absolute path"
	second.Description = "new wording"
	second.Platforms = []string{"LINUX", "MACOS", "WINDOWS"}

	if FindingIdentity(first) != FindingIdentity(second) {
		t.Fatal("baseline identity changed after line, copy, platform, or path normalization changes")
	}
	if FindingFingerprint(first) == FindingFingerprint(second) {
		t.Fatal("exact fingerprint should change when emitted finding metadata changes")
	}
}

func TestFindingIdentityChangesWhenSeverityChanges(t *testing.T) {
	first := Finding{RuleID: "x", Path: "config.txt", Severity: SeverityWarning}
	second := first
	second.Severity = SeverityError
	if FindingIdentity(first) == FindingIdentity(second) {
		t.Fatal("severity upgrades should re-surface a baseline finding")
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
