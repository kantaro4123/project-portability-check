# Changelog

All notable changes to this project will be documented here.

## [0.2.0] - 2026-09-05

### Added

- Concurrent detector execution with deterministic final ordering for faster scans on larger repositories.
- Monorepo-aware runtime pin and lockfile inheritance for nested Node.js, Python, Go, and Cargo projects.
- `--target` and `target_platforms` support for focusing findings on selected operating systems.
- JSON baselines with count-aware matching so existing findings can be suppressed while newly introduced duplicates remain visible.
- Stable exact finding fingerprints in JSON and richer SARIF rule metadata and partial fingerprints.
- GitHub Action inputs for target platforms, baselines, and report file output.
- A stable public rule catalog used by `--list-rules` and configuration suppressions.

### Fixed

- CI platform coverage no longer counts runner names that appear only in YAML comments.
- OS target filtering no longer hides architecture-only findings such as `amd64`/`arm64` risks.
- Runtime analysis now checks every `go.mod` instead of stopping after the first module.
- Project-internal baseline files are excluded from their own analysis.

### Changed

- Version advanced to `0.2.0`.
- CI now exercises baseline suppression, target filtering, GitHub Action output files, and strict self-analysis on Linux, macOS, and Windows.
- `--list-rules` now returns actual stable finding IDs rather than internal detector group IDs.

## [0.1.0] - 2026-09-05

### Added

- Cross-platform project scanner with a dependency-free Go implementation.
- Portability score and structured findings with stable rule IDs.
- Filesystem checks for Windows-reserved names, forbidden characters, long paths, case collisions, symlinks, and executable scripts.
- Text checks for mixed line endings, UTF-8 issues, and machine-specific absolute paths.
- Shell checks for common GNU/BSD incompatibilities.
- JavaScript package-script checks for Unix-only commands and environment syntax.
- Runtime and dependency reproducibility checks for Node.js, Python, Go, Cargo, and lockfiles.
- Environment-variable documentation hints.
- Native binary and fixed Docker platform detection.
- GitHub Actions platform-coverage and `.gitattributes` guidance.
- Human-readable, JSON, and SARIF 2.1.0 output.
- `.portabilitycheck.json` suppression configuration with a JSON Schema.
- Reusable composite GitHub Action for CI integration.
- Dependency-directory pruning for faster scans of real-world projects.
- Linux, macOS, and Windows CI with strict self-dogfood checks.
