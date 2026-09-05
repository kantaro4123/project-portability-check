# project-portability-check

[![CI](https://github.com/kantaro4123/project-portability-check/actions/workflows/ci.yml/badge.svg)](https://github.com/kantaro4123/project-portability-check/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**Find the reasons a project works on one machine but breaks on another.**

`project-portability-check` is a static analyzer for cross-platform portability hazards. It scans a project without executing its code and reports filesystem, path, shell, runtime, dependency, text-encoding, architecture, Docker, and CI assumptions that can break across Windows, macOS, Linux, or CPU architectures.

```text
$ project-portability-check --target windows,linux .
project-portability-check v0.2.0
Project: /path/to/project
Scanned: 84 file(s)
Targets: linux, windows

Findings
  ! Machine-specific absolute path (src/config.ts:18) [paths.absolute]
    Found a macOS user path that may fail on another machine.
    Affects: linux, windows
    Fix: Use a relative path, environment variable, or configurable project root.

  ! Unix-specific package script (package.json) [package.script-unix]
    package.json script clean uses rm, which is not portable to the default Windows shell.
    Affects: windows

Portability Score: 90/100
Errors: 0  Warnings: 2  Info: 0
```

## Why use it

- **Static and safe** — it never executes project files, package scripts, build tools, or discovered binaries.
- **Cross-platform by design** — checks filesystem, shell, runtime, dependency, architecture, Docker, and CI assumptions.
- **Monorepo-aware** — nested Node.js, Python, Go, and Cargo projects can inherit runtime pins and lockfiles from parent workspaces.
- **Incremental adoption** — JSON baselines let an existing project fail CI only on newly introduced portability problems.
- **CI-friendly** — stable rule IDs, exit codes, JSON, richer SARIF, deterministic fingerprints, and a reusable GitHub Action.
- **Fast on large repositories** — dependency environments are pruned and independent detector groups run concurrently while output remains deterministic.

## Install

With Go:

```bash
go install github.com/kantaro4123/project-portability-check/cmd/project-portability-check@latest
```

Or build from source:

```bash
git clone https://github.com/kantaro4123/project-portability-check.git
cd project-portability-check
go build ./cmd/project-portability-check
```

The analyzer has no runtime dependencies beyond the compiled binary.

## GitHub Action

After a tagged public release, the same analyzer can run directly in GitHub Actions:

```yaml
- uses: actions/checkout@v7
- uses: kantaro4123/project-portability-check@v0.2.0
  with:
    path: .
    strict: "true"
    format: json
    target: linux,macos,windows
    output: artifacts/portability.json
```

For gradual adoption, commit or restore an earlier JSON report and pass it as a baseline:

```yaml
- uses: kantaro4123/project-portability-check@v0.2.0
  with:
    path: .
    strict: "true"
    baseline: .portability-baseline.json
```

Action inputs:

| Input | Default | Meaning |
| --- | --- | --- |
| `path` | `.` | Project path relative to the workspace |
| `strict` | `true` | Fail on warnings as well as errors |
| `format` | `text` | `text`, `json`, or `sarif` |
| `target` | empty | Optional comma-separated `linux`, `macos`, `windows` targets |
| `baseline` | empty | Previous JSON report relative to the analyzed project root |
| `output` | empty | Optional report path relative to the workspace |

The composite action builds the checker from its own source and then analyzes the caller's workspace.

## What it checks

### Filesystems and paths

- macOS `/Users/...`, Linux `/home/...`, and Windows `C:\Users\...` absolute paths
- Windows reserved names such as `CON`, `NUL`, `COM1`, and `LPT1`
- Windows-forbidden filename characters and risky long paths
- files that collide on case-insensitive filesystems
- symbolic links, including absolute and project-external targets
- shebang scripts that lack executable permission on Unix-like systems

### Shell and text

- mixed LF/CRLF line endings
- non-UTF-8 source text and UTF-8 BOMs
- GNU/BSD incompatibilities including `grep -P`, `sed -i`, `readlink -f`, `date -d`, and `xargs -r`
- Unix-only commands and environment syntax in `package.json` scripts
- missing `.gitattributes` guidance for repositories with platform-sensitive scripts

### Runtimes and dependencies

- unpinned Node.js and Python development runtimes, including nested monorepo packages
- missing Go language version directives across multiple `go.mod` files
- missing JavaScript lockfiles and informational Cargo lockfile guidance, with parent-workspace inheritance
- code that references environment variables without an example environment file

### Build and delivery

- checked-in ELF, Mach-O, and PE/COFF native binaries
- Dockerfiles fixed to one CPU platform such as `linux/amd64`
- GitHub Actions configurations that do not exercise Linux, macOS, and Windows; YAML comments do not count as real coverage

See [the full rule reference](docs/rules.md).

## Usage

```bash
project-portability-check [options] [path]
```

The path defaults to the current directory.

```bash
# Human-readable report
project-portability-check .

# Machine-readable report with stable finding fingerprints
project-portability-check --json .

# SARIF 2.1.0 with rule metadata and partial fingerprints
project-portability-check --sarif .

# Make warnings fail CI too
project-portability-check --strict .

# Only keep findings relevant to selected operating-system targets
project-portability-check --target windows,linux .

# Suppress findings that already existed in a previous JSON report
project-portability-check --strict --baseline .portability-baseline.json .

# Show stable finding rule IDs accepted by ignore_rules
project-portability-check --list-rules

# Version
project-portability-check --version
```

## Baselines

Baselines make it practical to introduce strict portability checks to a repository that already has known issues.

Create a JSON report:

```bash
project-portability-check --json . > .portability-baseline.json
```

Then use it on later runs:

```bash
project-portability-check --strict --baseline .portability-baseline.json .
```

Known findings are matched by semantic identity rather than line number, so adding unrelated lines does not make an old issue look new. Matching uses counts, so introducing an additional duplicate of an existing problem is still reported. If the baseline is stored inside the analyzed project, it is automatically excluded from that scan.

## Configuration

Create `.portabilitycheck.json` in the project root to suppress intentional findings or select target operating systems:

```json
{
  "$schema": "https://raw.githubusercontent.com/kantaro4123/project-portability-check/main/schemas/portabilitycheck.schema.json",
  "target_platforms": [
    "linux",
    "windows"
  ],
  "ignore_rules": [
    "deps.cargo-lockfile"
  ],
  "ignore_paths": [
    "vendor/*",
    "testdata/*"
  ]
}
```

`--target` overrides `target_platforms` for a command. OS targeting only filters findings explicitly tagged with an operating system; architecture findings such as `arm64`/`amd64` remain visible. `ignore_paths` uses slash-separated glob patterns. The repository includes the published schema at [`schemas/portabilitycheck.schema.json`](schemas/portabilitycheck.schema.json).

## Machine-readable output

JSON findings include an exact deterministic `fingerprint`. SARIF 2.1.0 output includes rule descriptors, remediation help, normalized locations, and `partialFingerprints` for code-scanning integrations. Public finding IDs are stable and can be listed with `--list-rules`.

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | No portability error was found |
| `1` | At least one error was found, or a warning was found with `--strict` |
| `2` | The analyzer could not run or the CLI arguments were invalid |

Informational findings never fail the command.

## Safety and privacy

The analyzer does **not** execute project files, package scripts, build tools, or discovered binaries. Text and binary scanners do not follow symlinks. All analysis is local and no project content is uploaded by the CLI.

## Scoring

The portability score starts at 100. Errors carry a larger penalty than warnings; informational findings do not reduce the score. The score is a prioritization aid, not a guarantee that software will behave identically on every platform.

## Development

```bash
gofmt -w .
go vet ./...
go test ./...
go run ./cmd/project-portability-check --strict --target linux,macos,windows .
```

CI runs formatting, vet, tests, strict self-analysis, baseline and target smoke tests, the local composite action, builds, and output smoke tests on Linux, macOS, and Windows.

See [Architecture](docs/architecture.md), [Rule reference](docs/rules.md), and [Contributing](CONTRIBUTING.md).

## License

MIT
