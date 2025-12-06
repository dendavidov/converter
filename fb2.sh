#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 /path/to/books_dir"
  exit 1
fi

BOOKS_DIR="$(realpath "$1")"

sudo docker run --rm \
  -u "$(id -u):$(id -g)" \
  -v "$BOOKS_DIR:/data" \
  ebook-pdf \
  /data -o /data -r --ext fb2
