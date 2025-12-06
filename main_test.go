package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestResolveExtensions_Defaults(t *testing.T) {
	got := resolveExtensions(extList{})
	if !reflect.DeepEqual(got, defaultExts) {
		t.Fatalf("resolveExtensions default = %v, want %v", got, defaultExts)
	}
}

func TestResolveExtensions_Custom(t *testing.T) {
	var exts extList
	if err := exts.Set("epub"); err != nil {
		t.Fatalf("Set epub: %v", err)
	}
	if err := exts.Set("fb2"); err != nil {
		t.Fatalf("Set fb2: %v", err)
	}

	got := resolveExtensions(exts)
	expected := map[string]struct{}{".epub": {}, ".fb2": {}}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("resolveExtensions custom = %v, want %v", got, expected)
	}
}

func TestFormatExtListSorts(t *testing.T) {
	exts := map[string]struct{}{".fb2": {}, ".epub": {}}
	got := formatExtList(exts)
	want := []string{".epub", ".fb2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("formatExtList = %v, want %v", got, want)
	}
}

func TestMatchExtension(t *testing.T) {
	allowed := map[string]struct{}{".epub": {}}
	if !matchExtension("book.epub", allowed) {
		t.Fatalf("expected book.epub to match")
	}
	if matchExtension("book.fb2", allowed) {
		t.Fatalf("did not expect book.fb2 to match")
	}
}

func TestCollectFilesNonRecursive(t *testing.T) {
	root := t.TempDir()
	allowed := map[string]struct{}{".epub": {}, ".fb2": {}}

	files := []string{"one.epub", "two.fb2", "skip.txt", filepath.Join("nested", "deep.epub")}
	for _, name := range files {
		full := filepath.Join(root, name)
		if err := createTestFile(full); err != nil {
			t.Fatalf("createTestFile %s: %v", full, err)
		}
	}

	got, err := collectFiles(root, false, allowed)
	if err != nil {
		t.Fatalf("collectFiles returned error: %v", err)
	}
	want := []string{
		filepath.Join(root, "one.epub"),
		filepath.Join(root, "two.fb2"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectFiles non-recursive = %v, want %v", got, want)
	}
}

func TestCollectFilesRecursive(t *testing.T) {
	root := t.TempDir()
	allowed := map[string]struct{}{".epub": {}, ".fb2": {}}

	files := []string{"one.epub", "two.fb2", "skip.txt", filepath.Join("nested", "deep.epub")}
	for _, name := range files {
		full := filepath.Join(root, name)
		if err := createTestFile(full); err != nil {
			t.Fatalf("createTestFile %s: %v", full, err)
		}
	}

	got, err := collectFiles(root, true, allowed)
	if err != nil {
		t.Fatalf("collectFiles returned error: %v", err)
	}
	want := []string{
		filepath.Join(root, "nested", "deep.epub"),
		filepath.Join(root, "one.epub"),
		filepath.Join(root, "two.fb2"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectFiles recursive = %v, want %v", got, want)
	}
}

func TestParseFlagsRequiresInput(t *testing.T) {
	if _, err := parseFlags([]string{}); err == nil {
		t.Fatalf("expected error for missing input")
	}
}

func TestResolveOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := createNamedFile(t, tmpDir, "file.epub")
	fileInfo := mustStatFile(t, filePath)
	dirInfo := mustStatDir(t, t.TempDir())
	existingFile := createNamedFile(t, t.TempDir(), "output.txt")

	out, err := resolveOutputDir("", fileInfo, "/tmp/dir/file.epub")
	if err != nil {
		t.Fatalf("resolveOutputDir file: %v", err)
	}
	if out != "/tmp/dir" {
		t.Fatalf("resolveOutputDir file out = %s, want /tmp/dir", out)
	}

	dirPath := "/tmp/dir"
	out, err = resolveOutputDir("", dirInfo, dirPath)
	if err != nil {
		t.Fatalf("resolveOutputDir dir: %v", err)
	}
	if out != dirPath {
		t.Fatalf("resolveOutputDir dir out = %s, want %s", out, dirPath)
	}

	customOut := "/custom/out"
	out, err = resolveOutputDir(customOut, dirInfo, dirPath)
	if err != nil {
		t.Fatalf("resolveOutputDir custom: %v", err)
	}
	if out != customOut {
		t.Fatalf("resolveOutputDir custom out = %s, want %s", out, customOut)
	}

	if _, err := resolveOutputDir(existingFile, dirInfo, dirPath); err == nil {
		t.Fatalf("expected error for existing non-directory output")
	}
}

func TestCheckEbookConvert(t *testing.T) {
	mockEbookConvert(t)
	if err := checkEbookConvert(); err != nil {
		t.Fatalf("checkEbookConvert returned error: %v", err)
	}
}

func TestConvertFileCreatesOutputAndSkipsExisting(t *testing.T) {
	mockEbookConvert(t)
	input := createNamedFile(t, t.TempDir(), "book.epub")
	outDir := t.TempDir()

	convertFile(input, outDir, false)
	output := filepath.Join(outDir, "book.pdf")
	info1 := mustStatFile(t, output)

	// Second call should skip and preserve mod time.
	convertFile(input, outDir, false)
	info2 := mustStatFile(t, output)
	if !info1.ModTime().Equal(info2.ModTime()) {
		t.Fatalf("expected modtime unchanged on skip")
	}
}

func TestRunConvertsDirectory(t *testing.T) {
	mockEbookConvert(t)
	root := t.TempDir()
	createNamedFile(t, root, "one.epub")
	createNamedFile(t, root, "two.fb2")
	createNamedFile(t, root, "skip.txt")

	if err := run([]string{root}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "one.pdf")); err != nil {
		t.Fatalf("expected one.pdf: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "two.pdf")); err != nil {
		t.Fatalf("expected two.pdf: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "skip.pdf")); !os.IsNotExist(err) {
		t.Fatalf("unexpected skip.pdf created")
	}
}

func TestRunSkipsNonMatchingFile(t *testing.T) {
	mockEbookConvert(t)
	root := t.TempDir()
	file := createNamedFile(t, root, "note.txt")

	if err := run([]string{file}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "note.pdf")); !os.IsNotExist(err) {
		t.Fatalf("unexpected pdf created for non-matching file")
	}
}

func TestRunErrorsOnNonDirectoryOutput(t *testing.T) {
	outputFile := createNamedFile(t, t.TempDir(), "out.txt")
	inputDir := t.TempDir()
	createNamedFile(t, inputDir, "book.epub")

	if err := run([]string{"-o", outputFile, inputDir}); err == nil {
		t.Fatalf("expected error for output path that is not a directory")
	}
}

func TestParseFlagsCollectsExts(t *testing.T) {
	opts, err := parseFlags([]string{"-ext", "epub", "-ext", "fb2", "file"})
	if err != nil {
		t.Fatalf("parseFlags error: %v", err)
	}
	if opts.input != "file" {
		t.Fatalf("expected input 'file', got %s", opts.input)
	}
	if got := opts.extFlag.String(); got != "epub,fb2" {
		t.Fatalf("expected ext list, got %s", got)
	}
}

func createNamedFile(t *testing.T, dir, name string) string {
	t.Helper()
	full := filepath.Join(dir, name)
	if err := createTestFile(full); err != nil {
		t.Fatalf("createNamedFile %s: %v", full, err)
	}
	return full
}

func createTestFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	// #nosec G304 - test helper writes to a controlled temporary path
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

func mockEbookConvert(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	bin := filepath.Join(tmp, "ebook-convert")
	script := "#!/usr/bin/env bash\nif [ \"$1\" = \"--version\" ]; then exit 0; fi\nOUT=\"$2\"\nmkdir -p \"$(dirname \"$OUT\")\"\ntouch \"$OUT\"\n"
	// Create an executable mock; gosec warning suppressed because this is a test helper.
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil { //nolint:gosec // test helper uses execute bit intentionally
		t.Fatalf("write mock ebook-convert: %v", err)
	}
	t.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))
	prev := ebookConvertBin
	ebookConvertBin = bin
	t.Cleanup(func() {
		ebookConvertBin = prev
	})
}

func mustStatFile(t *testing.T, path string) os.FileInfo {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat file %s: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file at %s", path)
	}
	return info
}

func mustStatDir(t *testing.T, dir string) os.FileInfo {
	t.Helper()
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat dir %s: %v", dir, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected dir at %s", dir)
	}
	return info
}
