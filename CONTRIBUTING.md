# Contributing

Contributions are welcome, especially new high-confidence portability rules, false-positive reductions, and cross-platform regression cases.

## Development

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
go run ./cmd/project-portability-check --strict --target linux,macos,windows .
```

The repository CI runs on Linux, macOS, and Windows. Race detection runs on Linux for the concurrent analyzer pipeline.

## Adding a rule

1. Add a focused detector under `internal/detectors`, or extend an existing detector when the new rule shares the same evidence scan.
2. Give findings stable rule IDs; changing an existing rule ID is a breaking machine-output change.
3. Prefer evidence that can be established without executing untrusted project code.
4. Add a regression test with the smallest realistic positive fixture and an important non-match when false positives are plausible.
5. Register new detector groups in `internal/detectors/defaults.go` when needed.
6. Add every public finding ID to `internal/detectors/rules.go` so `--list-rules`, JSON/SARIF consumers, and suppressions stay in sync.
7. Document the rule in `docs/rules.md` and mention user-visible additions in `CHANGELOG.md`.
8. Run the strict self-check so the repository remains a clean dogfood target.

Use `error` only for concrete incompatibilities, `warning` for strong portability risks, and `info` for recommendations that are often intentional.

## Pull requests

Keep changes focused and explain the operating systems, filesystems, shells, runtimes, or architectures affected. New rules should include examples of both a positive match and an important non-match when false positives are plausible. Machine-output compatibility matters: stable rule IDs and baseline identities should not change for cosmetic wording updates.
