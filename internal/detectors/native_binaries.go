package detectors

import (
	"context"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"

	"github.com/kantaro4123/project-portability-check/internal/analyzer"
	"github.com/kantaro4123/project-portability-check/internal/model"
)

type NativeBinaries struct{}

func (NativeBinaries) ID() string { return "binary.native" }

func (NativeBinaries) Detect(_ context.Context, project analyzer.Project) ([]model.Finding, error) {
	var findings []model.Finding
	for _, rel := range project.Files {
		full := filepath.Join(project.Root, filepath.FromSlash(rel))
		f, err := os.Open(full)
		if err != nil {
			continue
		}
		header := make([]byte, 64)
		n, _ := io.ReadFull(f, header)
		_ = f.Close()
		header = header[:n]
		kind, platforms := nativeKind(header)
		if kind == "" {
			continue
		}
		findings = append(findings, model.Finding{
			RuleID:      "binary.native",
			Title:       "Tracked platform-specific native binary",
			Description: "Detected a " + kind + " executable/object. Committed native binaries commonly fail on other operating systems or CPU architectures.",
			Severity:    model.SeverityWarning,
			Path:        rel,
			Platforms:   platforms,
			Suggestion:  "Build native artifacts per target platform or publish them as release artifacts instead of relying on one checked-in binary.",
		})
	}
	return findings, nil
}

func nativeKind(data []byte) (string, []string) {
	if len(data) >= 4 && data[0] == 0x7f && data[1] == 'E' && data[2] == 'L' && data[3] == 'F' {
		return "ELF", []string{"macos", "windows"}
	}
	if len(data) >= 4 {
		magic := binary.BigEndian.Uint32(data[:4])
		switch magic {
		case 0xfeedface, 0xfeedfacf, 0xcafebabe, 0xcefaedfe, 0xcffaedfe, 0xbebafeca:
			return "Mach-O", []string{"linux", "windows"}
		}
	}
	if len(data) >= 2 && data[0] == 'M' && data[1] == 'Z' {
		return "PE/COFF", []string{"macos", "linux"}
	}
	return "", nil
}
