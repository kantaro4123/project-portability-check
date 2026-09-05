package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

const Filename = ".portabilitycheck.json"

type Config struct {
	IgnoreRules []string `json:"ignore_rules,omitempty"`
	IgnorePaths []string `json:"ignore_paths,omitempty"`
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
	return cfg, nil
}

func (c Config) Filter(findings []model.Finding) []model.Finding {
	out := findings[:0]
	for _, finding := range findings {
		if contains(c.IgnoreRules, finding.RuleID) || matchesAny(c.IgnorePaths, finding.Path) {
			continue
		}
		out = append(out, finding)
	}
	return out
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
