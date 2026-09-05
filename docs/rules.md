# Rule reference

`project-portability-check --list-rules` prints the stable finding rule IDs enabled by the current binary. These are the same IDs emitted in JSON and SARIF and accepted by `.portabilitycheck.json` `ignore_rules`.

## Filesystem and paths

- `paths.absolute` — machine-specific macOS, Linux, and Windows home paths.
- `fs.windows-reserved` — Windows device names such as `CON`, `NUL`, `COM1`, and `LPT1`.
- `fs.windows-trailing` — path components ending in a space or dot.
- `fs.windows-forbidden-char` — characters such as `?`, `*`, `|`, or `:` that Windows rejects.
- `fs.windows-long-path` — repository-relative paths likely to become problematic after checkout.
- `fs.case-collision` — paths that differ only by letter case.
- `fs.symlink` — symlinks, with stronger severity for absolute or project-external targets.
- `fs.script-not-executable` — shebang scripts without Unix executable permission.

## Text and shell

- `text.mixed-line-endings` — files containing multiple newline conventions.
- `text.non-utf8` — likely source/text files that are not valid UTF-8.
- `text.utf8-bom` — UTF-8 BOMs that can surprise Unix tooling.
- `shell.grep-p`, `shell.sed-i`, `shell.readlink-f`, `shell.date-d`, `shell.xargs-r` — GNU/BSD command-line incompatibilities.
- `package.script-unix` — npm-compatible scripts that rely on Unix shell commands or syntax.

## Dependencies and runtimes

- `runtime.node-unpinned` — Node.js package without a recognized runtime pin in the package or a parent workspace.
- `runtime.python-unpinned` — Python project without a recognized development interpreter pin in the project or a parent workspace.
- `runtime.go-unpinned` — any `go.mod` without a `go` directive.
- `deps.node-lockfile` — JavaScript package without a recognized lockfile in the package or a parent workspace.
- `deps.cargo-lockfile` — Cargo project without `Cargo.lock` in the crate or parent workspace (informational because libraries often omit it intentionally).
- `env.no-example` — environment-dependent code without an example environment file.

## Build and delivery

- `binary.native` — checked-in ELF, Mach-O, or PE/COFF artifacts.
- `docker.fixed-platform` — Dockerfile pinned to a single CPU platform.
- `git.no-gitattributes` — platform-sensitive scripts without a repository line-ending policy.
- `ci.platform-coverage` — GitHub Actions configuration that does not reference all three major desktop OS families. Runner names that only occur in YAML comments do not count as coverage.

## Target filtering

`--target` and the `target_platforms` configuration field filter findings that are explicitly tagged with operating systems. Findings that only carry architecture labels such as `amd64` or `arm64` remain visible because OS selection does not remove CPU-architecture risks.

## Severity

`error` means the project contains a concrete cross-platform blocker. `warning` means behavior is likely to vary between common environments. `info` is a reproducibility or maintenance recommendation that is often intentional.
