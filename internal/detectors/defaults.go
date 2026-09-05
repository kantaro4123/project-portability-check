package detectors

import "github.com/kantaro4123/project-portability-check/internal/analyzer"

// Default returns the built-in detector set in a stable order.
func Default() []analyzer.Detector {
	return []analyzer.Detector{
		AbsolutePaths{},
		WindowsNames{},
		WindowsPaths{},
		CaseCollisions{},
		Symlinks{},
		LineEndings{},
		ExecutableScripts{},
		ShellPortability{},
		PackageScripts{},
		RuntimePins{},
		Lockfiles{},
		EnvironmentVariables{},
		NativeBinaries{},
		DockerPlatform{},
		GitAttributes{},
		TextEncoding{},
		CIMatrix{},
	}
}
