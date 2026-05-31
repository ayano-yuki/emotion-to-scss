package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"emotion-to-scss/internal/converter"
	"emotion-to-scss/internal/domain"
	"emotion-to-scss/internal/parser"
	"emotion-to-scss/internal/report"
	"emotion-to-scss/internal/verifier"
)

const (
	ExitSuccess            = 0
	ExitVerificationFailed = 1
	ExitInvalidArgs        = 2
	ExitFileError          = 3
	ExitParseError         = 4
)

type options struct {
	outDir            string
	dryRun            bool
	overwrite         bool
	failOnUnsupported bool
	reportPath        string
}

type inputFile struct {
	Path string
	Rel  string
}

func Run(args []string, stdout, stderr io.Writer) (int, error) {
	if len(args) == 0 {
		return ExitInvalidArgs, usageError("missing command")
	}

	command := args[0]
	if command != "convert" && command != "verify" && command != "check" {
		return ExitInvalidArgs, usageError("unknown command %q", command)
	}

	input, opts, err := parseArgs(command, args[1:])
	if err != nil {
		return ExitInvalidArgs, err
	}

	files, err := collectInput(input)
	if err != nil {
		return ExitFileError, err
	}

	result, failed, err := runCommand(command, input, files, opts, stdout)
	if err != nil {
		var parseErr *parser.ParseError
		if errors.As(err, &parseErr) {
			writeReportIfRequested(opts.reportPath, result)
			return ExitParseError, err
		}
		writeReportIfRequested(opts.reportPath, result)
		return ExitFileError, err
	}

	if opts.reportPath != "" {
		if err := report.Write(opts.reportPath, result); err != nil {
			return ExitFileError, err
		}
	}

	if failed {
		return ExitVerificationFailed, nil
	}
	return ExitSuccess, nil
}

func parseArgs(command string, args []string) (string, options, error) {
	opts := options{}
	var input string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--out":
			i++
			if i >= len(args) {
				return "", options{}, usageError("--out requires a value")
			}
			opts.outDir = args[i]
		case "--dry-run":
			opts.dryRun = true
		case "--overwrite":
			opts.overwrite = true
		case "--fail-on-unsupported":
			opts.failOnUnsupported = true
		case "--report":
			i++
			if i >= len(args) {
				return "", options{}, usageError("--report requires a value")
			}
			opts.reportPath = args[i]
		default:
			if strings.HasPrefix(arg, "-") {
				return "", options{}, usageError("unknown option %q", arg)
			}
			if input != "" {
				return "", options{}, usageError("expected exactly one input")
			}
			input = arg
		}
	}

	if input == "" {
		return "", options{}, usageError("expected exactly one input")
	}
	if command != "check" && opts.outDir == "" {
		return "", options{}, usageError("--out is required")
	}
	return input, opts, nil
}

func runCommand(command, input string, files []inputFile, opts options, stdout io.Writer) (report.Report, bool, error) {
	result := report.Report{}
	failed := false

	for _, file := range files {
		data, err := os.ReadFile(file.Path)
		if err != nil {
			return result, false, err
		}
		styles, err := parser.Parse(string(data))
		if err != nil {
			return result, false, err
		}

		fileReport := report.FileReport{Path: file.Path}
		for _, style := range styles {
			styleReport := report.StyleReport{
				Name:         style.Name,
				Status:       "converted",
				Line:         style.Line,
				ClassName:    style.ClassName,
				Verification: "skipped",
			}

			if command == "verify" {
				verification, err := verifier.Verify(style)
				if err != nil {
					styleReport.Status = "failed"
					styleReport.Verification = "failed"
					styleReport.Reason = err.Error()
					failed = true
				} else if !verification.Equivalent {
					styleReport.Status = "failed"
					styleReport.Verification = "failed"
					styleReport.Reason = verification.Reason
					failed = true
				} else {
					styleReport.Verification = "passed"
				}
			}

			if command == "check" {
				if err := verifier.CheckSCSS(converter.Convert(style)); err != nil {
					styleReport.Status = "failed"
					styleReport.Reason = err.Error()
					failed = true
				} else {
					styleReport.Status = "checked"
				}
			}

			fileReport.Styles = append(fileReport.Styles, styleReport)
		}

		if command == "convert" || command == "verify" {
			if err := writeConverted(input, file, styles, opts); err != nil {
				return result, false, err
			}
			if opts.dryRun {
				fmt.Fprintln(stdout, converter.ConvertAll(styles))
			}
		}

		result.Files = append(result.Files, fileReport)
	}

	fillSummary(&result)
	return result, failed, nil
}

func writeConverted(input string, file inputFile, styles []domain.Style, opts options) error {
	if opts.dryRun {
		return nil
	}

	outputPath := outputPath(input, file, opts.outDir)
	if !opts.overwrite {
		if _, err := os.Stat(outputPath); err == nil {
			return fmt.Errorf("output file already exists: %s", outputPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(converter.ConvertAll(styles)), 0o644)
}

func outputPath(input string, file inputFile, outDir string) string {
	info, err := os.Stat(input)
	if err == nil && info.IsDir() {
		rel := strings.TrimSuffix(file.Rel, filepath.Ext(file.Rel)) + ".scss"
		return filepath.Join(outDir, rel)
	}
	base := strings.TrimSuffix(filepath.Base(file.Path), filepath.Ext(file.Path)) + ".scss"
	return filepath.Join(outDir, base)
}

func fillSummary(r *report.Report) {
	r.Summary.Files = len(r.Files)
	for _, file := range r.Files {
		r.Summary.Styles += len(file.Styles)
		for _, style := range file.Styles {
			switch style.Status {
			case "converted", "checked":
				r.Summary.Converted++
			case "failed":
				r.Summary.Failed++
			}
			if style.Verification == "passed" {
				r.Summary.Verified++
			}
		}
	}
}

func collectInput(input string) ([]inputFile, error) {
	info, err := os.Stat(input)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if !isSupported(input) {
			return nil, fmt.Errorf("unsupported extension: %s", input)
		}
		return []inputFile{{Path: input, Rel: filepath.Base(input)}}, nil
	}

	var files []inputFile
	err = filepath.WalkDir(input, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !isSupported(path) {
			return nil
		}
		rel, err := filepath.Rel(input, path)
		if err != nil {
			return err
		}
		files = append(files, inputFile{Path: path, Rel: rel})
		return nil
	})
	return files, err
}

func isSupported(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ts", ".tsx", ".js", ".jsx":
		return true
	default:
		return false
	}
}

func writeReportIfRequested(path string, r report.Report) {
	if path == "" {
		return
	}
	_ = report.Write(path, r)
}

func usageError(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
