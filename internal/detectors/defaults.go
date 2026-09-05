package detectors

import "github.com/kantaro4123/project-portability-check/internal/analyzer"

// Default returns the built-in detector set in a stable order.
func Default() []analyzer.Detector {
	return []analyzer.Detector{
		AbsolutePaths{},
		WindowsNames{},
		CaseCollisions{},
		Symlinks{},
		LineEndings{},
		ExecutableScripts{},
		ShellPortability{},
		RuntimePins{},
		Lockfiles{},
		EnvironmentVariables{},
		NativeBinaries{},
		GitAttributes{},
		TextEncoding{},
	}
}
