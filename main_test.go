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

func createTestFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}
