import sys
from pathlib import Path

import fb2_to_pdf


def test_convert_file_skips_non_fb2(monkeypatch, tmp_path, capsys):
    input_file = tmp_path / "note.txt"
    input_file.write_text("content")

    calls = []
    monkeypatch.setattr(fb2_to_pdf.subprocess, "run", lambda *args, **kwargs: calls.append(args))

    fb2_to_pdf.convert_file(input_file, tmp_path)

    assert calls == []
    assert "Skipping non-fb2 file" in capsys.readouterr().out


def test_main_converts_recursive(monkeypatch, tmp_path):
    nested = tmp_path / "nested"
    nested.mkdir()
    fb2_file = nested / "story.fb2"
    fb2_file.write_text("content")

    convert_calls = []

    def fake_convert_file(input_path: Path, output_dir: Path, overwrite: bool = False):
        convert_calls.append((input_path, output_dir, overwrite))

    monkeypatch.setattr(fb2_to_pdf, "convert_file", fake_convert_file)
    monkeypatch.setattr(fb2_to_pdf, "check_ebook_convert", lambda: None)
    monkeypatch.setattr(sys, "argv", ["fb2_to_pdf.py", str(tmp_path), "-r"])

    fb2_to_pdf.main()

    assert convert_calls == [(fb2_file, tmp_path, False)]
