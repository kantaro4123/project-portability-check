# Architecture

The project is intentionally dependency-free at runtime. The standard library is enough to walk a project tree, inspect files, parse the small set of structured manifests used by built-in rules, and serialize machine-readable output.

## Pipeline

1. `internal/project` discovers regular files and symbolic links while excluding VCS metadata and common dependency environments.
2. `internal/config` loads `.portabilitycheck.json` and normalizes optional operating-system targets.
3. `internal/analyzer` runs independent `Detector` implementations concurrently. Results are collected and sorted deterministically before stable finding fingerprints are attached.
4. `internal/detectors` contains focused portability rules. Detectors do not mutate the analyzed project. Workspace-aware rules can inherit runtime pins and lockfiles from parent directories.
5. Configuration suppressions and target-platform filtering are applied after detection so detector behavior stays reusable and deterministic.
6. If `--baseline` is active, `internal/baseline` suppresses matching findings from an earlier JSON report. Matching tolerates line movement and human-readable copy changes while preserving occurrence counts so new duplicate problems remain visible.
7. `internal/report` computes the score and renders text, JSON, or SARIF. SARIF includes public rule descriptors and partial fingerprints.
8. `internal/cli` owns flags, baseline resolution, exit codes, and output selection.

## Detector contract

A detector has an internal group ID and receives only the project root plus normalized relative file paths. Findings contain a stable public rule ID, severity, affected path/line when known, affected platforms or architectures, and a suggested remediation.

Public finding IDs are cataloged separately from detector group IDs. `--list-rules`, JSON, SARIF, and `ignore_rules` all use the public finding IDs.

Detectors should prefer high-confidence static evidence. If a condition is commonly intentional, it should generally be `info` instead of `warning` or `error`.

## Determinism and concurrency

Detector groups are read-only and independent, so the analyzer can run them concurrently. It never relies on goroutine completion order: results are collected, then sorted by severity, path, line, and rule ID. This keeps text, JSON, SARIF, and baseline behavior stable across machines while reducing wall-clock time on larger repositories.

The Linux CI job runs the Go race detector over the full test suite to guard this concurrent pipeline against shared-state regressions.

## Finding and baseline identity

Machine output has two related identities:

- an exact finding fingerprint includes the rule, normalized path, line, severity, display copy, and normalized platform metadata and is exposed in JSON/SARIF;
- baseline identity intentionally uses stable semantic fields — rule ID, normalized path, and severity — so line movement, remediation wording, or description improvements do not reintroduce a known finding.

Baseline matching is count-aware rather than set-based. If one instance existed in the baseline and a second equivalent instance is introduced, only one is suppressed.

## Target filtering

Target selection applies only to findings explicitly tagged with operating-system names (`linux`, `macos`, `windows`). Architecture-only findings such as `amd64`/`arm64` remain visible under an OS filter.

## Safety

Text scanning uses `Lstat` and size limits. Symlink targets are inspected as metadata but are not followed by text or binary scanners. The analyzer never executes project code, package-manager scripts, build systems, or discovered binaries.

## Portability of the analyzer

The analyzer itself is tested on Linux, macOS, and Windows. The CLI uses Go's platform-independent path APIs internally and normalizes report paths to slash-separated repository-relative strings. CI also exercises strict self-analysis, target filtering, baselines, the composite GitHub Action, and machine-readable output.
