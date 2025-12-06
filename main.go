// Package main provides a CLI for converting EPUB/FB2 files to PDF via Calibre.
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

var ebookConvertBin = "ebook-convert"

type extList struct {
	values []string
}

type config struct {
	inputPath   string
	outputDir   string
	recursive   bool
	overwrite   bool
	allowedExts map[string]struct{}
	allowedList []string
	isFile      bool
}

type flagOptions struct {
	outputDir string
	recursive bool
	overwrite bool
	extFlag   extList
	input     string
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
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseConfig(args)
	if err != nil {
		return err
	}

	if err := checkEbookConvert(); err != nil {
		return err
	}

	if cfg.isFile {
		ext := strings.ToLower(filepath.Ext(cfg.inputPath))
		if _, ok := cfg.allowedExts[ext]; !ok {
			fmt.Printf("No matching extension for file: %s (allowed: %s)\n", cfg.inputPath, strings.Join(cfg.allowedList, ", "))
			return nil
		}
		convertFile(cfg.inputPath, cfg.outputDir, cfg.overwrite)
		return nil
	}

	files, err := collectFiles(cfg.inputPath, cfg.recursive, cfg.allowedExts)
	if err != nil {
		return fmt.Errorf("failed to collect files: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("No files with extensions %s found in %s.\n", strings.Join(cfg.allowedList, ", "), cfg.inputPath)
		return nil
	}

	for _, f := range files {
		convertFile(f, cfg.outputDir, cfg.overwrite)
	}

	return nil
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
	info, err := os.Stat(outputDir)
	switch {
	case err == nil && !info.IsDir():
		fmt.Fprintf(os.Stderr, "[ERROR] Output path is not a directory: %s\n", outputDir)
		return
	case err != nil && !errors.Is(err, os.ErrNotExist):
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to inspect output directory %s: %v\n", outputDir, err)
		return
	}

	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to create output directory %s: %v\n", outputDir, err)
		return
	}

	stem := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outputPDF := filepath.Join(outputDir, stem+".pdf")

	if _, err := os.Stat(outputPDF); err == nil && !overwrite {
		fmt.Printf("[SKIP] Already exists: %s\n", outputPDF)
		return
	}

	// #nosec G204 - inputs are intentionally user-provided for conversion
	cmd := exec.Command(
		ebookConvertBin,
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
	cmd := exec.Command(ebookConvertBin, "--version")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return errors.New("error: 'ebook-convert' not found. install Calibre")
	}
	return nil
}

func parseConfig(args []string) (config, error) {
	opts, err := parseFlags(args)
	if err != nil {
		return config{}, err
	}

	inputPath, err := filepath.Abs(opts.input)
	if err != nil {
		return config{}, fmt.Errorf("unable to resolve input path: %w", err)
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config{}, fmt.Errorf("input path does not exist: %s", inputPath)
		}
		return config{}, fmt.Errorf("unable to read input path: %w", err)
	}

	allowedExts := resolveExtensions(opts.extFlag)
	allowedList := formatExtList(allowedExts)

	resolvedOutput, err := resolveOutputDir(opts.outputDir, info, inputPath)
	if err != nil {
		return config{}, err
	}

	return config{
		inputPath:   inputPath,
		outputDir:   resolvedOutput,
		recursive:   opts.recursive,
		overwrite:   opts.overwrite,
		allowedExts: allowedExts,
		allowedList: allowedList,
		isFile:      info.Mode().IsRegular(),
	}, nil
}

func parseFlags(args []string) (flagOptions, error) {
	var opts flagOptions

	fs := flag.NewFlagSet("converter", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.outputDir, "o", "", "Directory to save PDFs (default: same as input)")
	fs.StringVar(&opts.outputDir, "output-dir", "", "Directory to save PDFs (default: same as input)")
	fs.BoolVar(&opts.recursive, "r", false, "Recursively search for files in directories")
	fs.BoolVar(&opts.recursive, "recursive", false, "Recursively search for files in directories")
	fs.BoolVar(&opts.overwrite, "overwrite", false, "Overwrite existing PDF files")
	fs.Var(&opts.extFlag, "ext", "Limit to these extensions (can be used multiple times). Default: fb2 and epub.")

	if err := fs.Parse(args); err != nil {
		return flagOptions{}, err
	}

	if fs.NArg() < 1 {
		return flagOptions{}, errors.New("input path is required")
	}

	opts.input = fs.Arg(0)
	return opts, nil
}

func resolveOutputDir(outputFlag string, info os.FileInfo, inputPath string) (string, error) {
	if outputFlag != "" {
		resolved, err := filepath.Abs(outputFlag)
		if err != nil {
			return "", fmt.Errorf("unable to resolve output directory: %w", err)
		}

		existing, err := os.Stat(resolved)
		switch {
		case err == nil && !existing.IsDir():
			return "", fmt.Errorf("output path is not a directory: %s", resolved)
		case err != nil && !errors.Is(err, os.ErrNotExist):
			return "", fmt.Errorf("unable to check output directory: %w", err)
		}

		return resolved, nil
	}

	if info.Mode().IsRegular() {
		return filepath.Dir(inputPath), nil
	}

	return inputPath, nil
}
