# Converter

Scripts and Docker tooling for converting EPUB and FB2 ebooks to PDF using Calibre’s `ebook-convert`.

## Features
- Convert individual files or entire directories to PDF.
- Optional recursive search and overwrite control.
- Docker image that bundles Calibre for isolated use.
- Simple helper scripts for EPUB-only or FB2-only batches.

## Requirements
- Calibre installed locally (`ebook-convert` must be in `PATH`) for direct script use.
- Docker (if you prefer the containerized workflow).
- Python 3.8+ for running the helper scripts directly.

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

## Quick Start (Python)
Convert a file:
```
python ebook_to_pdf.py /path/to/book.epub
```
Convert a directory (non-recursive by default):
```
python ebook_to_pdf.py /path/to/books_dir -o /path/to/output
```
Convert recursively and overwrite existing PDFs:
```
python ebook_to_pdf.py /path/to/books_dir -r --overwrite
```
Limit to extensions (can be provided multiple times):
```
python ebook_to_pdf.py /path/to/books_dir --ext epub --ext fb2
```

## Script Reference
- `ebook_to_pdf.py` – primary CLI; supports recursion, overwrite flag, and extension filtering.
- `epub_to_pdf.py` – legacy variant for EPUB/FB2.
- `fb2_to_pdf.py` – legacy FB2-only variant.
- `epub.sh`, `fb2.sh`, `run.sh` – Docker helpers that call the built image with sensible defaults.

## Notes on Content
- Real ebooks should **not** be committed to version control. Keep your own collection outside the repo or add `books/` to `.gitignore` and supply sample placeholders if needed.
- Ensure you have the right to convert and store the files you process.

## Troubleshooting
- If you see `ERROR: 'ebook-convert' not found`, install Calibre and ensure `ebook-convert` is on `PATH`, or use the Docker workflow.
- PDF margins default to 36 points on all sides; adjust flags in the scripts if you need different page settings.

## Development
- Python formatting: follow default `black`/`ruff` conventions if you add checks.
- Docker: the image uses `python:3.12-slim` and installs Calibre with minimal dependencies.

## License
- MIT License (`LICENSE`). Provided as-is, without warranty; use at your own risk.
