# Contributing

Contributions are welcome, especially new high-confidence portability rules, false-positive reductions, and cross-platform regression cases.

## Development

```bash
go test ./...
go vet ./...
gofmt -w .
go run ./cmd/project-portability-check .
```

The repository CI runs on Linux, macOS, and Windows.

## Adding a rule

1. Add a focused detector under `internal/detectors`.
2. Give findings stable rule IDs; changing an existing rule ID is a breaking machine-output change.
3. Prefer evidence that can be established without executing untrusted project code.
4. Add a regression test with the smallest realistic fixture.
5. Register the detector in `internal/detectors/defaults.go`.
6. Document the rule in `docs/rules.md`.

Use `error` only for concrete incompatibilities, `warning` for strong portability risks, and `info` for recommendations that are often intentional.

## Pull requests

Keep changes focused and explain the operating systems or filesystems affected. New rules should include examples of both a positive match and an important non-match when false positives are plausible.
