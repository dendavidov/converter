# Agent Guidelines

This repository contains a Go 1.22 CLI and Docker tooling for converting EPUB/FB2 ebooks to PDF via Calibre. Calibre is expected to be available inside the container image; local execution requires `ebook-convert` on `PATH`.

## Scope
- These instructions apply to the entire repository.

## Development Practices
- Prefer Go 1.22 tooling. Run `gofmt` on any Go files you modify.
- Keep the CLI behavior stable: recursive search, extension filtering (EPUB/FB2), overwrite flag, and `ebook-convert` invocation with the existing margin defaults.
- Avoid introducing external dependencies unless necessary; the project currently uses only the Go standard library.
- When updating shell scripts, keep them POSIX-sh compatible and preserve the current Docker-first workflow.

## Testing & Linting
- Before committing, run `make test` (or `go test ./...`) to ensure the CLI logic remains sound.
- When touching Go code, prefer `make lint` locally if available; it installs and runs `golangci-lint` and the test suite.

## Documentation & Messages
- Align README/usage comments with the Docker-centered workflow (Calibre inside the image, optional local build).
- Commit messages should summarize the change succinctly (e.g., "Add AGENTS instructions for repo").
