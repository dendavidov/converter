# Security Policy

## Project Overview

`converter` is a Go 1.25 CLI and Docker image that wraps Calibre's `ebook-convert` tool to batch-convert EPUB/FB2 ebooks to PDF. Because users run the binary or Docker image against their own content, we treat security and supply-chain integrity with care.

## Supported Versions

We provide security updates for the **current main branch and the latest published Docker image tag**. Older tags or forks are not maintained.

- Go toolchain requirement: **Go 1.25+**
- Docker workflow: **latest `ebook-pdf` image** built from this repository

Always pull the latest image or rebuild from `main` to receive security fixes.

## Security Scope

This security policy covers:

**In Scope:**
- Vulnerabilities in the CLI or Docker image that could lead to arbitrary code execution, privilege escalation, or data leakage
- Supply chain risks in build scripts (`Dockerfile`, `Makefile`, shell wrappers) and Go dependencies
- Malicious or unsafe defaults when invoking `ebook-convert` (e.g., unsafe path handling)
- Integrity of published Docker images and tagged releases

**Out of Scope:**
- Vulnerabilities in Calibre or `ebook-convert` itself (report upstream)
- Issues caused by third-party ebook files with malicious content
- Problems in user-provided Docker configurations or host mounts
- Performance issues without a clear security impact

## Reporting a Vulnerability

### Private Reporting

Please report suspected vulnerabilities **privately** to avoid exposing users before a fix is ready.

- **GitHub Security Advisories**: Use the "Report a vulnerability" feature on this repository for coordinated disclosure.
- **Email (fallback)**: [mailto:mail@dendavidov.com](mailto:mail@dendavidov.com)

Include the following where possible:
- Description of the vulnerability and affected components (CLI, Docker image, shell scripts)
- Steps to reproduce, including sample commands and input files if applicable
- Expected vs. actual behavior and potential impact
- Environment details (OS, Docker version, Go version)
- Whether you would like public credit

**Please DO NOT:**
- Open public GitHub issues for security-sensitive reports
- Share exploit details publicly before we confirm and release a fix

### Non-security Bugs

For non-security issues (feature requests, usability bugs), use the standard [GitHub Issues](https://github.com/dendavidov/converter/issues) tracker.

## Response Process

We aim for the following timelines after receiving a private security report:

| Stage | Target Timeline |
|-------|-----------------|
| Initial response | Within 72 hours |
| Vulnerability triage and confirmation | Within 7 business days |
| Fix development and validation | Varies by complexity |
| Patch release | See severity targets below |

Patch release targets by severity (CVSS v3.1 as reference):

| Severity | Patch Target |
|----------|--------------|
| Critical | 14 business days |
| High     | 30 business days |
| Medium   | 60 business days |
| Low      | Best effort |

Our process:
1. Acknowledge receipt of your report and establish a secure communication channel.
2. Reproduce and assess impact/severity.
3. Develop and test a fix (and tests where applicable).
4. Coordinate disclosure and publish patched releases/Docker images.
5. Credit reporters in release notes if desired.

## Disclosure Policy

We follow coordinated disclosure. Please allow an embargo period while we validate and release fixes. Security fixes will be highlighted in release notes. If warranted, we will request CVE assignment.

## Security Measures We Implement

- Minimal dependency footprint (standard library–only Go code)
- Docker-first workflow to isolate Calibre runtime and reduce host exposure
- Reproducible builds via Make targets and pinned Go version
- Protected main branch with required reviews for changes
- Sample data only; no secrets committed or expected at runtime

## Known Security Considerations

- The tool invokes Calibre's `ebook-convert` on user-provided files; ensure input files come from trusted sources.
- Host mounts in Docker (`-v /abs/books:/data`) grant the container access to those paths—limit mounts to required directories.
- Running the binary directly inherits the permissions of the invoking user; prefer Docker when possible.

## Security Resources

- [Calibre Security Notes](https://manual.calibre-ebook.com/faq.html#security)
- [Go Security Policy](https://go.dev/security/policy)
- [Docker Security Best Practices](https://docs.docker.com/security/)

## Questions?

If you have questions about this policy or need guidance on secure usage, please reach out via the GitHub Issues tracker for non-sensitive topics or the private reporting channels above.

---

**Last Updated**: 2025-12-06
**Policy Version**: 1.0
