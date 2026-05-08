# Architectural Decisions

This document explains the key technical choices made in SBOMForge and the reasoning behind them.

---

## Why Syft for SBOM generation?

[Syft](https://github.com/anchore/syft) is the de facto standard for source-based SBOM generation in the open source ecosystem. Compared to alternatives:

| Tool | Reason not chosen |
|---|---|
| `cyclonedx-gomod` | Go-only, no multi-language support |
| `spdx-tools` | Java-based, heavy runtime dependency |
| `trivy sbom` | Security scanner first, SBOM second — less accurate for source scanning |
| `tern` | Container-focused, poor source code support |

Syft supports all three major SBOM formats (SPDX, CycloneDX, Syft JSON), detects ecosystems automatically (Go, Node, Python, Java, Rust, etc.), and is actively maintained by Anchore. It is invoked as a binary via `exec.CommandContext`, keeping the integration simple and decoupled from any Go API changes.

---

## Why Cosign keyless signing?

Traditional signing requires managing a private key — generating it, storing it securely (usually as a repository secret), rotating it, and distributing the public key for verification. This is error-prone and creates operational overhead.

[Cosign keyless signing](https://docs.sigstore.dev/cosign/signing/overview/) via Sigstore OIDC eliminates this entirely:

- GitHub Actions provides a native OIDC token that proves the identity of the workflow
- Cosign exchanges this token for a short-lived certificate from Sigstore's Fulcio CA
- The signature is recorded in Sigstore's Rekor transparency log
- Verification requires no pre-shared key — just the SBOM and the `.bundle` file

This is a better security posture than key-based signing: there is no long-lived secret to leak, and every signature is publicly auditable via the transparency log.

---

## Why Go?

The DevSecOps tooling ecosystem is largely written in Go (Syft, Cosign, Trivy, kubectl, Terraform, etc.). Choosing Go means:

- **Static binary**: the compiled output has zero runtime dependencies, which is ideal for a minimal Docker image
- **Cross-compilation**: trivial to build for `linux/amd64` and `linux/arm64` from any machine
- **Standard library**: `os/exec`, `encoding/json`, `net/http` cover most needs without external dependencies
- **Performance**: fast startup time matters in CI environments where every second counts

The only external dependencies are `google/go-github` (typed GitHub API client) and `golang.org/x/oauth2` (token authentication) — both well-maintained and widely used.

---

## Why a multi-stage Docker build?

GitHub Actions Docker actions run the image directly — there is no separate build step on the runner. The image must therefore contain everything needed to run: the compiled binary, Syft, and Cosign.

A multi-stage build keeps the final image minimal:

- **Stage 1 (`golang:1.22-alpine`)**: compiles the Go binary with `CGO_ENABLED=0` for a fully static executable
- **Stage 2 (`alpine:3.21`)**: installs Syft and Cosign, then copies only the compiled binary from stage 1

This approach avoids shipping the Go toolchain (~300MB) in the final image. The `go mod download` step is placed before `COPY . .` to maximize Docker layer caching: dependency downloads are only re-run when `go.mod` or `go.sum` change, not on every code change.

---

## Why a separate token input instead of using `GITHUB_TOKEN` directly?

`GITHUB_TOKEN` is always available in GitHub Actions and could be read directly from the environment. However, requiring it as an explicit input (`github-token`) has two advantages:

1. **Clarity**: the user explicitly grants the action access to their token, making the permission grant visible in the workflow file
2. **Flexibility**: a user can provide a PAT with broader permissions if needed (e.g., to upload to a release in a different repository)

This mirrors the convention used by most official GitHub Actions (`actions/checkout`, `actions/create-release`, etc.).
