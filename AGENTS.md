# AGENTS.md

## Project Overview
Go 1.22 CLI and Docker image for converting EPUB/FB2 ebooks to PDF via Calibre’s `ebook-convert`. Docker is the recommended path; local runs require Calibre installed on PATH.

## Repository Structure
- `main.go` / `main_test.go` – CLI entrypoint and unit tests for file discovery, flag parsing, and command execution.
- `Makefile` – build/test/lint helpers and Docker image target (`ebook-pdf`).
- `Dockerfile` – multi-stage build bundling the Go binary and Calibre runtime.
- `epub.sh`, `fb2.sh`, `run.sh` – Docker wrappers for format-specific or sample conversions.
- `books/` – sample input directory (mounted by `run.sh`).
- `AGENTS.md` – AI-facing instructions (this file).

## Setup & Tooling
- Go: 1.22 (see `go.mod`).
- Docker: required for the supported workflow; Calibre lives in the image.
- Local Calibre: only needed if running the Go binary directly (not preferred).
- Cache: `GOCACHE` is pinned in Make targets to avoid polluting the global cache.

## Build, Run & Test
### Build
```bash
make build          # builds ./converter binary
make docker         # builds the ebook-pdf image
```

### Run (Docker-first)
```bash
# Generic
make docker
# Mount an absolute books dir and convert recursively
# Replace /abs/books with your path
ID="$(id -u):$(id -g)"
docker run --rm -u "$ID" -v /abs/books:/data ebook-pdf /data -o /data -r
```

### Run locally (not preferred; needs Calibre on PATH)
```bash
go build -o converter .
./converter /path/to/books -o /path/to/output -r --ext epub
```

### Tests
```bash
make test           # or: go test ./...
make coverage       # coverage summary
make lint           # installs golangci-lint v2.7.1, runs lint + tests
```

## Coding Conventions
- Use only the Go standard library; avoid new dependencies unless justified.
- Keep CLI flags and defaults stable: recursion flag, overwrite flag, `--ext` filter, and default 36pt margins passed to `ebook-convert`.
- Format Go with `gofmt` (required before commit).
- Shell scripts should remain POSIX-sh compatible and preserve the Docker-first UX.

## Workflows
- Branches: use short, descriptive names (e.g., `feat/...`, `fix/...`, `chore/...`).
- Commits: concise summaries; conventional commit prefixes are welcome but not enforced.
- PRs: include a short description of behavior changes and note any impact on Docker or CLI usage. Run `make test` (and `make lint` when touching Go) before raising a PR.

## Security, Secrets & Data
- Do not add secrets, tokens, or private file paths. No secrets are expected in this repo.
- Use test/fake data in examples; ensure you have rights to process any real ebooks.
- Do not commit local Calibre configs or caches; `books/` is sample-only.

## AI-Specific Instructions
- Prefer minimal, localized changes that keep parity between Docker scripts and CLI behavior.
- When adding behavior, include unit tests in `main_test.go` and update README examples if user-facing behavior changes.
- Avoid altering CI/packaging patterns (Makefile, Dockerfile) unless explicitly required by a task.
- Keep commands literal—use Make targets where available rather than re-creating them in scripts.

## Appendix
Example Docker invocation (recursive, keep formats as-is):
```bash
ID="$(id -u):$(id -g)"
docker run --rm -u "$ID" -v /abs/books:/data ebook-pdf /data -o /data -r --ext fb2
```
