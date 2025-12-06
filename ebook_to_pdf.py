#!/usr/bin/env python3
import argparse
import subprocess
from pathlib import Path
import sys

DEFAULT_EXTS = {".fb2", ".epub"}


def check_ebook_convert():
    try:
        subprocess.run(
            ["ebook-convert", "--version"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            check=True,
        )
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("ERROR: 'ebook-convert' not found. Install Calibre.", file=sys.stderr)
        sys.exit(1)


def convert_file(input_path: Path, output_dir: Path, overwrite: bool = False):
    output_dir.mkdir(parents=True, exist_ok=True)
    output_pdf = output_dir / (input_path.stem + ".pdf")

    if output_pdf.exists() and not overwrite:
        print(f"[SKIP] Already exists: {output_pdf}")
        return

    cmd = [
        "ebook-convert",
        str(input_path),
        str(output_pdf),
        "--pdf-page-margin-top",
        "36",
        "--pdf-page-margin-bottom",
        "36",
        "--pdf-page-margin-left",
        "36",
        "--pdf-page-margin-right",
        "36",
    ]

    print(f"[CONVERT] {input_path} → {output_pdf}")
    try:
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError as e:
        print(f"[ERROR] Failed to convert {input_path}: {e}", file=sys.stderr)


def main():
    parser = argparse.ArgumentParser(
        description="Convert ebooks (EPUB/FB2) to PDF using Calibre."
    )
    parser.add_argument("input", help="Input file or directory")
    parser.add_argument(
        "-o",
        "--output-dir",
        default=None,
        help="Directory to save PDFs (default: same as input)",
    )
    parser.add_argument(
        "-r",
        "--recursive",
        action="store_true",
        help="Recursively search for files in directories",
    )
    parser.add_argument(
        "--overwrite",
        action="store_true",
        help="Overwrite existing PDF files",
    )
    parser.add_argument(
        "--ext",
        action="append",
        choices=["fb2", "epub"],
        help="Limit to these extensions (can be used multiple times). "
             "Default: fb2 and epub.",
    )

    args = parser.parse_args()
    check_ebook_convert()

    input_path = Path(args.input).expanduser().resolve()

    if not input_path.exists():
        print(f"ERROR: Input path does not exist: {input_path}", file=sys.stderr)
        sys.exit(1)

    # Compute allowed extensions
    if args.ext:
        exts = {"." + e.lower() for e in args.ext}
    else:
        exts = set(DEFAULT_EXTS)

    # Determine output directory
    if args.output_dir:
        output_dir = Path(args.output_dir).expanduser().resolve()
    else:
        output_dir = input_path.parent if input_path.is_file() else input_path

    if input_path.is_file():
        if input_path.suffix.lower() not in exts:
            print(
                f"No matching extension for file: {input_path} "
                f"(allowed: {', '.join(exts)})"
            )
            sys.exit(0)
        convert_file(input_path, output_dir, overwrite=args.overwrite)
    else:
        pattern = "**/*" if args.recursive else "*"
        files = [
            p
            for p in input_path.glob(pattern)
            if p.is_file() and p.suffix.lower() in exts
        ]

        if not files:
            print(
                f"No files with extensions {', '.join(exts)} found in {input_path}."
            )
            sys.exit(0)

        for f in sorted(files):
            convert_file(f, output_dir, overwrite=args.overwrite)


if __name__ == "__main__":
    main()
