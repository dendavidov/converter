import sys
from pathlib import Path

import ebook_to_pdf


def test_convert_file_skips_existing_pdf(monkeypatch, tmp_path, capsys):
    input_file = tmp_path / "book.epub"
    input_file.write_text("content")
    output_dir = tmp_path / "output"
    output_dir.mkdir()
    (output_dir / "book.pdf").write_text("existing pdf")

    calls = []

    def fake_run(*args, **kwargs):
        calls.append(args)

    monkeypatch.setattr(ebook_to_pdf.subprocess, "run", fake_run)

    ebook_to_pdf.convert_file(input_file, output_dir)

    assert calls == []
    assert "SKIP" in capsys.readouterr().out


def test_convert_file_invokes_ebook_convert(monkeypatch, tmp_path):
    input_file = tmp_path / "story.fb2"
    input_file.write_text("dummy")
    output_dir = tmp_path / "pdfs"

    run_calls = []

    def fake_run(cmd, check):
        run_calls.append(cmd)

    monkeypatch.setattr(ebook_to_pdf.subprocess, "run", fake_run)

    ebook_to_pdf.convert_file(input_file, output_dir, overwrite=True)

    expected_output = output_dir / "story.pdf"
    expected_cmd = [
        "ebook-convert",
        str(input_file),
        str(expected_output),
        "--pdf-page-margin-top",
        "36",
        "--pdf-page-margin-bottom",
        "36",
        "--pdf-page-margin-left",
        "36",
        "--pdf-page-margin-right",
        "36",
    ]

    assert run_calls == [expected_cmd]
    assert expected_output.parent.exists()


def test_main_filters_extensions(monkeypatch, tmp_path):
    epub = tmp_path / "keep.epub"
    fb2 = tmp_path / "skip.fb2"
    mobi = tmp_path / "ignore.mobi"
    for path in (epub, fb2, mobi):
        path.write_text("content")

    convert_calls = []

    def fake_convert_file(input_path: Path, output_dir: Path, overwrite: bool = False):
        convert_calls.append((input_path, output_dir, overwrite))

    monkeypatch.setattr(ebook_to_pdf, "convert_file", fake_convert_file)
    monkeypatch.setattr(ebook_to_pdf, "check_ebook_convert", lambda: None)
    monkeypatch.setattr(sys, "argv", ["ebook_to_pdf.py", str(tmp_path), "--ext", "epub"])

    ebook_to_pdf.main()

    assert convert_calls == [(epub, tmp_path, False)]
