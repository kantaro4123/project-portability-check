# Architecture

The project is intentionally dependency-free at runtime. The standard library is enough to walk a project tree, inspect files, parse the small set of structured manifests used by built-in rules, and serialize machine-readable output.

## Pipeline

1. `internal/project` discovers regular files and symbolic links while excluding VCS metadata.
2. `internal/config` loads `.portabilitycheck.json`.
3. `internal/analyzer` runs an ordered set of `Detector` implementations.
4. `internal/detectors` contains focused portability rules. Detectors do not mutate the analyzed project.
5. Configuration suppressions are applied after detection so detector behavior stays deterministic.
6. `internal/report` computes the score and renders text, JSON, or SARIF.
7. `internal/cli` owns flags, exit codes, and output selection.

## Detector contract

A detector has a stable group ID and receives only the project root plus normalized relative file paths. Findings contain a stable rule ID, severity, affected path/line when known, affected platforms, and a suggested remediation.

Detectors should prefer high-confidence static evidence. If a condition is commonly intentional, it should generally be `info` instead of `warning` or `error`.

## Safety

Text scanning uses `Lstat` and size limits. Symlink targets are inspected as metadata but are not followed by text or binary scanners. The analyzer never executes project code, package-manager scripts, build systems, or discovered binaries.

## Portability of the analyzer

The analyzer itself is tested on Linux, macOS, and Windows. The CLI uses Go's platform-independent path APIs internally and normalizes report paths to slash-separated repository-relative strings.
