# Contributing to SBOMForge

SBOMForge is maintained as a personal project and does not accept pull requests from external contributors. However, bug reports and feature suggestions are very welcome through GitHub Issues.

## Developer Certificate of Origin (DCO)

All commits must include a `Signed-off-by` line per the [Developer Certificate of Origin v1.1](https://developercertificate.org). Use `git commit -s` to add this automatically.

## Coding standards

- Code formatted with `gofmt`
- Linting enforced with `golangci-lint` (runs in CI)
- Standard Go idioms per [Effective Go](https://go.dev/doc/effective_go)

## Contribution criteria

Changes must satisfy all of the following:

- Pass `golangci-lint run ./...`
- Pass `go test -race ./...`
- New features include tests with >80% statement coverage for `internal/sbom` and `internal/release`
- No secrets or credentials committed (Gitleaks validation required)

## Local setup

```bash
# Run tests
go test -race ./...

# Run linter
golangci-lint run ./...

# Build the binary
go build -o sbomforge ./cmd/

# Simulate CI locally (requires Docker + act)
act -j test
```

## Security

Report vulnerabilities through [GitHub Security Advisories](https://github.com/Richonn/SBOMForge/security/advisories/new), not public issues.
