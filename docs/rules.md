# Rule reference

`project-portability-check --list-rules` prints the detector groups enabled by the current binary. Individual findings use stable rule IDs so they can be suppressed in `.portabilitycheck.json` and consumed from JSON or SARIF.

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

- `runtime.node-unpinned` — Node.js project without a recognized runtime pin.
- `runtime.python-unpinned` — Python project without a recognized development interpreter pin.
- `runtime.go-unpinned` — `go.mod` without a `go` directive.
- `deps.node-lockfile` — JavaScript project without a recognized lockfile.
- `deps.cargo-lockfile` — Cargo project without `Cargo.lock` (informational because libraries often omit it intentionally).
- `env.no-example` — environment-dependent code without an example environment file.

## Build and delivery

- `binary.native` — checked-in ELF, Mach-O, or PE/COFF artifacts.
- `docker.fixed-platform` — Dockerfile pinned to a single CPU platform.
- `git.no-gitattributes` — platform-sensitive scripts without a repository line-ending policy.
- `ci.platform-coverage` — GitHub Actions configuration that does not reference all three major desktop OS families.

## Severity

`error` means the project contains a concrete cross-platform blocker. `warning` means behavior is likely to vary between common environments. `info` is a reproducibility or maintenance recommendation that is often intentional.
