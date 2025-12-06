#!/usr/bin/env python3
import argparse
import subprocess
from pathlib import Path
import sys

SUPPORTED_EXT = {".epub", ".fb2"}

def check_ebook_convert():
    try:
        subprocess.run(["ebook-convert", "--version"],
                       stdout=subprocess.DEVNULL,
                       stderr=subprocess.DEVNULL,
                       check=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("ERROR: 'ebook-convert' not found. Install Calibre.", file=sys.stderr)
        sys.exit(1)

def convert_file(input_path: Path, output_dir: Path, overwrite: bool = False):
    ext = input_path.suffix.lower()
    if ext not in SUPPORTED_EXT:
        print(f"Skipping unsupported file type: {input_path}")
        return

    output_dir.mkdir(parents=True, exist_ok=True)
    output_pdf = output_dir / (input_path.stem + ".pdf")

    if output_pdf.exists() and not overwrite:
        print(f"Already exists, skipping: {output_pdf}")
        return

    cmd = [
        "ebook-convert",
        str(input_path),
        str(output_pdf),
        "--pdf-page-margin-top", "36",
        "--pdf-page-margin-bottom", "36",
        "--pdf-page-margin-left", "36",
        "--pdf-page-margin-right", "36",
    ]

    print(f"Converting: {input_path} → {output_pdf}")
    try:
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError as e:
        print(f"Failed to convert {input_path}: {e}", file=sys.stderr)

def main():
    parser = argparse.ArgumentParser(
        description="Convert EPUB or FB2 ebooks to PDF using Calibre."
    )
    parser.add_argument("input", help="EPUB/FB2 file or directory")
    parser.add_argument("-o", "--output-dir", default=None,
                        help="Directory to save PDFs")
    parser.add_argument("-r", "--recursive", action="true",
                        help="Search directories recursively")
    parser.add_argument("--overwrite", action="store_true",
                        help="Overwrite existing PDFs")

    args = parser.parse_args()
    check_ebook_convert()

    input_path = Path(args.input).expanduser().resolve()

    if not input_path.exists():
        print("ERROR: Input path does not exist.", file=sys.stderr)
        sys.exit(1)

    # Determine output directory
    if args.output_dir:
        output_dir = Path(args.output_dir).expanduser().resolve()
    else:
        output_dir = input_path.parent if input_path.is_file() else input_path

    if input_path.is_file():
        convert_file(input_path, output_dir, args.overwrite)
    else:
        pattern = "**/*" if args.recursive else "*"
        files = [p for p in input_path.glob(pattern) if p.suffix.lower() in SUPPORTED_EXT]

        if not files:
            print("No EPUB/FB2 files found.")
            sys.exit(0)

        for f in files:
            convert_file(f, output_dir, args.overwrite)

if __name__ == "__main__":
    main()
