package detectors

import (
	"context"
	"regexp"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type DockerPlatform struct{}

func (DockerPlatform) ID() string { return "docker.platform" }

var dockerPlatformRE = regexp.MustCompile(`(?mi)^\s*FROM\s+--platform=(linux/(?:amd64|arm64|386|arm/v\d+))\s+`)

func (DockerPlatform) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		base := strings.ToLower(rel)
		if !strings.HasSuffix(base, "dockerfile") && !strings.Contains(base, "dockerfile.") {
			continue
		}
		data, ok := readText(project.Root, rel)
		if !ok {
			continue
		}
		for _, match := range dockerPlatformRE.FindAllSubmatchIndex(data, -1) {
			platform := string(data[match[2]:match[3]])
			findings = append(findings, model.Finding{RuleID: "docker.fixed-platform", Title: "Docker build pins one CPU platform", Description: "Dockerfile fixes the base image to " + platform + ", which can force emulation or fail on another architecture.", Severity: model.SeverityWarning, Path: rel, Line: lineNumber(data, match[0]), Platforms: []string{"arm64", "amd64"}, Suggestion: "Use build-time TARGETPLATFORM/BUILDPLATFORM values or publish a documented single-architecture constraint."})
		}
	}
	return findings, nil
}
