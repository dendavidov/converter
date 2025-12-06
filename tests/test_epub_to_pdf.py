import subprocess
from pathlib import Path

import epub_to_pdf


def test_convert_file_skips_unknown_extension(monkeypatch, tmp_path, capsys):
    input_file = tmp_path / "note.txt"
    input_file.write_text("content")

    calls = []
    monkeypatch.setattr(epub_to_pdf.subprocess, "run", lambda *args, **kwargs: calls.append(args))

    epub_to_pdf.convert_file(input_file, tmp_path)

    assert calls == []
    assert "Skipping unsupported file type" in capsys.readouterr().out


def test_convert_file_builds_expected_command(monkeypatch, tmp_path):
    input_file = tmp_path / "book.epub"
    input_file.write_text("content")
    output_dir = tmp_path / "out"

    run_calls = []

    def fake_run(cmd, check):
        run_calls.append(cmd)

    monkeypatch.setattr(epub_to_pdf.subprocess, "run", fake_run)

    epub_to_pdf.convert_file(input_file, output_dir, overwrite=True)

    expected_output = output_dir / "book.pdf"
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


def test_check_ebook_convert_exits_on_missing_binary(monkeypatch):
    def fake_run(*args, **kwargs):
        raise FileNotFoundError

    monkeypatch.setattr(epub_to_pdf.subprocess, "run", fake_run)

    try:
        epub_to_pdf.check_ebook_convert()
    except SystemExit as exc:
        assert exc.code == 1
    else:
        raise AssertionError("SystemExit not raised")
