# Security Policy

## Supported versions

| Version | Status |
|---|---|
| v1.x | Active — security fixes, bug fixes, feature updates |

When v2.x launches, v1.x enters a 90-day maintenance window for critical security issues only before reaching end-of-life.

## Vulnerability management

All **CRITICAL** and **HIGH** severity vulnerabilities in dependencies and container base images must be resolved before merging. Medium and lower severity issues follow best-effort remediation in subsequent release cycles.

Trivy scans every pull request and commit to the Docker image, blocking any release that contains unresolved critical or high-severity findings.

## Credential security

SBOMForge uses keyless signing via Sigstore OIDC — no signing keys are stored anywhere. The `github-token` input is processed as an environment variable and never logged. Gitleaks runs on every commit to prevent accidental credential exposure.

## Dependency management

Dependabot is configured for weekly automated updates (Go modules + GitHub Actions). All dependencies must maintain MIT license compatibility.

## Reporting a vulnerability

Please report vulnerabilities **privately** through [GitHub Security Advisories](https://github.com/Richonn/SBOMForge/security/advisories/new) rather than opening a public issue.

Expected response time: within 7 days.
