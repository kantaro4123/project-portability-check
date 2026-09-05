# Changelog

All notable changes to this project will be documented here.

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
- Linux, macOS, and Windows CI.
