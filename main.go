package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

var defaultExts = map[string]struct{}{
	".fb2":  {},
	".epub": {},
}

type extList struct {
	values []string
}

func (e *extList) String() string {
	if len(e.values) == 0 {
		return ""
	}
	return strings.Join(e.values, ",")
}

func (e *extList) Set(value string) error {
	if value == "" {
		return errors.New("extension cannot be empty")
	}

	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized != "fb2" && normalized != "epub" {
		return fmt.Errorf("invalid extension %q (allowed: epub, fb2)", value)
	}

	e.values = append(e.values, normalized)
	return nil
}

func main() {
	var (
		outputDir string
		recursive bool
		overwrite bool
		extFlag   extList
	)

	flag.StringVar(&outputDir, "o", "", "Directory to save PDFs (default: same as input)")
	flag.StringVar(&outputDir, "output-dir", "", "Directory to save PDFs (default: same as input)")
	flag.BoolVar(&recursive, "r", false, "Recursively search for files in directories")
	flag.BoolVar(&recursive, "recursive", false, "Recursively search for files in directories")
	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite existing PDF files")
	flag.Var(&extFlag, "ext", "Limit to these extensions (can be used multiple times). Default: fb2 and epub.")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: converter [options] <input>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := checkEbookConvert(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	inputPath, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to resolve input path: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "ERROR: Input path does not exist: %s\n", inputPath)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "ERROR: unable to read input path: %v\n", err)
		os.Exit(1)
	}

	allowedExts := resolveExtensions(extFlag)
	allowedList := formatExtList(allowedExts)

	var resolvedOutput string
	if outputDir != "" {
		resolvedOutput, err = filepath.Abs(outputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to resolve output directory: %v\n", err)
			os.Exit(1)
		}
	} else if info.Mode().IsRegular() {
		resolvedOutput = filepath.Dir(inputPath)
	} else {
		resolvedOutput = inputPath
	}

	if info.Mode().IsRegular() {
		ext := strings.ToLower(filepath.Ext(inputPath))
		if _, ok := allowedExts[ext]; !ok {
			fmt.Printf("No matching extension for file: %s (allowed: %s)\n", inputPath, strings.Join(allowedList, ", "))
			return
		}
		convertFile(inputPath, resolvedOutput, overwrite)
		return
	}

	files, err := collectFiles(inputPath, recursive, allowedExts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Printf("No files with extensions %s found in %s.\n", strings.Join(allowedList, ", "), inputPath)
		return
	}

	for _, f := range files {
		convertFile(f, resolvedOutput, overwrite)
	}
}

func resolveExtensions(extFlag extList) map[string]struct{} {
	if len(extFlag.values) == 0 {
		return defaultExts
	}

	res := make(map[string]struct{}, len(extFlag.values))
	for _, v := range extFlag.values {
		res["."+strings.ToLower(v)] = struct{}{}
	}
	return res
}

func formatExtList(exts map[string]struct{}) []string {
	result := make([]string, 0, len(exts))
	for ext := range exts {
		result = append(result, ext)
	}
	sort.Strings(result)
	return result
}

func collectFiles(root string, recursive bool, allowed map[string]struct{}) ([]string, error) {
	var files []string

	if !recursive {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, fmt.Errorf("unable to read directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			fullPath := filepath.Join(root, entry.Name())
			if matchExtension(fullPath, allowed) {
				files = append(files, fullPath)
			}
		}
	} else {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if matchExtension(path, allowed) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("unable to walk directory: %w", err)
		}
	}

	sort.Strings(files)
	return files, nil
}

func matchExtension(path string, allowed map[string]struct{}) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := allowed[ext]
	return ok
}

func convertFile(inputPath, outputDir string, overwrite bool) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to create output directory %s: %v\n", outputDir, err)
		return
	}

	stem := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outputPDF := filepath.Join(outputDir, stem+".pdf")

	if _, err := os.Stat(outputPDF); err == nil && !overwrite {
		fmt.Printf("[SKIP] Already exists: %s\n", outputPDF)
		return
	}

	cmd := exec.Command(
		"ebook-convert",
		inputPath,
		outputPDF,
		"--pdf-page-margin-top", "36",
		"--pdf-page-margin-bottom", "36",
		"--pdf-page-margin-left", "36",
		"--pdf-page-margin-right", "36",
	)

	fmt.Printf("[CONVERT] %s \u2192 %s\n", inputPath, outputPDF)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to convert %s: %v\n", inputPath, err)
	}
}

func checkEbookConvert() error {
	cmd := exec.Command("ebook-convert", "--version")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return errors.New("ERROR: 'ebook-convert' not found. Install Calibre.")
	}
	return nil
}
