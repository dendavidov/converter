# Converter

Go CLI and Docker tooling for converting EPUB and FB2 ebooks to PDF using Calibre’s `ebook-convert`.

## Features
- Convert individual files or entire directories to PDF.
- Optional recursive search and overwrite control.
- Docker image that bundles Calibre for isolated use.
- Simple helper scripts for EPUB-only or FB2-only batches.

## Requirements
- Docker (recommended workflow; no local Calibre install needed).
- Go 1.22+ only if you want to build the binary yourself or adjust the Docker image.

## Quick Start (Docker)
1) Build the image:
```
docker build -t ebook-pdf .
```
2) Convert a directory of ebooks to PDF in place:
```
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -v /absolute/path/to/books:/data \
  ebook-pdf \
  /data -o /data -r
```
3) Limit to a specific format by passing `--ext`:
```
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -v /absolute/path/to/books:/data \
  ebook-pdf \
  /data -o /data -r --ext epub
```

Helper wrappers:
- `./epub.sh /path/to/books_dir` – convert only EPUB files.
- `./fb2.sh /path/to/books_dir` – convert only FB2 files.
- `./run.sh` – example that mounts the bundled `books/` directory; adjust paths as needed.

## Why Docker-only?
Calibre pulls in a large dependency chain and is not recommended to install on the host. The provided image bundles Calibre plus the Go converter binary so your machine stays clean and the attack surface stays contained.

## Optional: Go CLI (with local Calibre)
If you explicitly want to run outside the container, you must install Calibre so `ebook-convert` is on `PATH`, then:
```
go build -o converter .
./converter /path/to/books_dir -r --ext epub
```
This local path is not recommended; prefer the containerized workflow above.

## CLI Reference
- `converter` – primary CLI; supports recursion, overwrite flag, and extension filtering.
- `epub.sh`, `fb2.sh`, `run.sh` – Docker helpers that call the built image with sensible defaults.

## Notes on Content
- Ensure you have the right to convert and store the files you process.

## Troubleshooting
- If you see `ERROR: 'ebook-convert' not found`, install Calibre and ensure `ebook-convert` is on `PATH`, or use the Docker workflow.
- PDF margins default to 36 points on all sides; adjust flags in the scripts if you need different page settings.

## Development
- Go formatting: `gofmt -w main.go`
- Docker: the image uses a multi-stage build and installs Calibre in the runtime stage.

## License
- MIT License (`LICENSE`). Provided as-is, without warranty; use at your own risk.
