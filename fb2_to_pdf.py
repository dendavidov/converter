#!/usr/bin/env python3
import argparse
import subprocess
from pathlib import Path
import sys

def check_ebook_convert():
    try:
        subprocess.run(["ebook-convert", "--version"],
                       stdout=subprocess.DEVNULL,
                       stderr=subprocess.DEVNULL,
                       check=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("ERROR: 'ebook-convert' (Calibre) not found in PATH.", file=sys.stderr)
        print("Install Calibre and make sure 'ebook-convert' is available.", file=sys.stderr)
        sys.exit(1)

def convert_file(input_path: Path, output_dir: Path, overwrite: bool = False):
    if input_path.suffix.lower() != ".fb2":
        print(f"Skipping non-fb2 file: {input_path}")
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
        description="Convert FB2 ebooks to PDF using Calibre's ebook-convert."
    )
    parser.add_argument(
        "input",
        help="FB2 file or directory containing .fb2 files"
    )
    parser.add_argument(
        "-o", "--output-dir",
        help="Directory to save PDFs (default: same as input or current dir)",
        default=None,
    )
    parser.add_argument(
        "-r", "--recursive",
        action="store_true",
        help="Recursively search for .fb2 files in the input directory",
    )
    parser.add_argument(
        "--overwrite",
        action="store_true",
        help="Overwrite existing PDF files",
    )

    args = parser.parse_args()
    check_ebook_convert()

    input_path = Path(args.input).expanduser().resolve()

    if not input_path.exists():
        print(f"ERROR: Input path does not exist: {input_path}", file=sys.stderr)
        sys.exit(1)

    if args.output_dir:
        output_dir = Path(args.output_dir).expanduser().resolve()
    else:
        output_dir = input_path.parent if input_path.is_file() else input_path

    if input_path.is_file():
        convert_file(input_path, output_dir, overwrite=args.overwrite)
    else:
        if args.recursive:
            fb2_files = list(input_path.rglob("*.fb2"))
        else:
            fb2_files = list(input_path.glob("*.fb2"))

        if not fb2_files:
            print("No .fb2 files found.", file=sys.stderr)
            sys.exit(1)

        for fb2 in fb2_files:
            convert_file(fb2, output_dir, overwrite=args.overwrite)

if __name__ == "__main__":
    main()
