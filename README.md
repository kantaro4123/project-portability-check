# project-portability-check

[![CI](https://github.com/kantaro4123/project-portability-check/actions/workflows/ci.yml/badge.svg)](https://github.com/kantaro4123/project-portability-check/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**Find the reasons a project works on one machine but breaks on another.**

`project-portability-check` is a static analyzer for cross-platform portability hazards. It scans a project without executing its code and reports filesystem, path, shell, runtime, dependency, text-encoding, architecture, Docker, and CI assumptions that can break across Windows, macOS, Linux, or CPU architectures.

```text
$ project-portability-check .
project-portability-check v0.1.0
Project: /path/to/project
Scanned: 84 file(s)

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
- uses: kantaro4123/project-portability-check@v0.1.0
  with:
    path: .
    strict: "true"
    format: text
```

The composite action builds the checker from its own source, then analyzes the caller's workspace. `format` accepts `text`, `json`, or `sarif`.

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

- unpinned Node.js and Python development runtimes
- missing Go language version directives
- missing JavaScript lockfiles and informational Cargo lockfile guidance
- code that references environment variables without an example environment file

### Build and delivery

- checked-in ELF, Mach-O, and PE/COFF native binaries
- Dockerfiles fixed to one CPU platform such as `linux/amd64`
- GitHub Actions configurations that do not exercise Linux, macOS, and Windows

See [the full rule reference](docs/rules.md).

## Usage

```bash
project-portability-check [options] [path]
```

The path defaults to the current directory.

```bash
# Human-readable report
project-portability-check .

# Machine-readable report
project-portability-check --json .

# SARIF 2.1.0 for code-scanning integrations
project-portability-check --sarif .

# Make warnings fail CI too
project-portability-check --strict .

# Show detector groups
project-portability-check --list-rules

# Version
project-portability-check --version
```

## Configuration

Create `.portabilitycheck.json` in the project root to suppress intentional findings:

```json
{
  "$schema": "https://raw.githubusercontent.com/kantaro4123/project-portability-check/main/schemas/portabilitycheck.schema.json",
  "ignore_rules": [
    "deps.cargo-lockfile"
  ],
  "ignore_paths": [
    "vendor/*",
    "testdata/*"
  ]
}
```

`ignore_paths` uses slash-separated glob patterns. The repository includes the published schema at [`schemas/portabilitycheck.schema.json`](schemas/portabilitycheck.schema.json).

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
go run ./cmd/project-portability-check --strict .
```

CI runs formatting, vet, tests, strict self-analysis, the local composite action, builds, and output smoke tests on Linux, macOS, and Windows.

See [Architecture](docs/architecture.md), [Rule reference](docs/rules.md), and [Contributing](CONTRIBUTING.md).

## License

MIT
