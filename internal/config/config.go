package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

const Filename = ".portabilitycheck.json"

var supportedPlatforms = map[string]bool{
	"linux":   true,
	"macos":   true,
	"windows": true,
}

type Config struct {
	IgnoreRules     []string `json:"ignore_rules,omitempty"`
	IgnorePaths     []string `json:"ignore_paths,omitempty"`
	TargetPlatforms []string `json:"target_platforms,omitempty"`
}

func Load(root string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(root, Filename))
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", Filename, err)
	}
	targets, err := NormalizeTargets(cfg.TargetPlatforms)
	if err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", Filename, err)
	}
	cfg.TargetPlatforms = targets
	return cfg, nil
}

func ParseTargets(value string) ([]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	return NormalizeTargets(strings.Split(value, ","))
}

func NormalizeTargets(values []string) ([]string, error) {
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}
		if !supportedPlatforms[value] {
			return nil, fmt.Errorf("unsupported target platform %q; expected linux, macos, or windows", value)
		}
		seen[value] = true
	}
	out := make([]string, 0, len(seen))
	for value := range seen {
		out = append(out, value)
	}
	sort.Strings(out)
	return out, nil
}

func (c Config) Filter(findings []model.Finding) []model.Finding {
	out := findings[:0]
	for _, finding := range findings {
		if contains(c.IgnoreRules, finding.RuleID) || matchesAny(c.IgnorePaths, finding.Path) {
			continue
		}
		if !appliesToTargets(finding, c.TargetPlatforms) {
			continue
		}
		out = append(out, finding)
	}
	return out
}

func appliesToTargets(finding model.Finding, targets []string) bool {
	if len(targets) == 0 || len(finding.Platforms) == 0 {
		return true
	}

	hasOSTag := false
	for _, platform := range finding.Platforms {
		platform = strings.ToLower(platform)
		if !supportedPlatforms[platform] {
			continue
		}
		hasOSTag = true
		for _, target := range targets {
			if platform == target {
				return true
			}
		}
	}
	// Some findings use Platforms for architecture labels (for example
	// amd64/arm64). OS target selection must not hide those findings.
	return !hasOSTag
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func matchesAny(patterns []string, value string) bool {
	if value == "" {
		return false
	}
	for _, pattern := range patterns {
		matched, err := path.Match(pattern, value)
		if err == nil && matched {
			return true
		}
	}
	return false
}
