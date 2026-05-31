package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"emotion-to-scss/internal/verifier"
)

const (
	ExitSuccess            = 0
	ExitVerificationFailed = 1
	ExitInvalidArgs        = 2
	ExitFileError          = 3
	ExitParseError         = 4
)

func Run(args []string, stdout, stderr io.Writer) (int, error) {
	input, opts, err := parseArgs(args)
	if err != nil {
		return ExitInvalidArgs, err
	}

	files, err := collectInputs(input)
	if err != nil {
		return ExitFileError, err
	}
	if len(files) == 0 {
		return ExitInvalidArgs, fmt.Errorf("no supported input files found")
	}

	failed := false
	for _, file := range files {
		result, err := verifier.VerifyFileWithOptions(file, verifier.Options{
			WriteAST: opts.debugAST,
		})
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return ExitFileError, err
			}
			fmt.Fprintf(stderr, "FAIL %s: %v\n", file, err)
			failed = true
			continue
		}
		if !result.OK {
			fmt.Fprintf(stderr, "FAIL %s: %s\n", file, result.Reason)
			failed = true
			continue
		}
		fmt.Fprintf(stdout, "OK %s\n", file)
	}

	if failed {
		return ExitVerificationFailed, nil
	}
	return ExitSuccess, nil
}

type options struct {
	debugAST bool
}

func parseArgs(args []string) (string, options, error) {
	if len(args) < 2 || args[0] != "check" {
		return "", options{}, fmt.Errorf("usage: emotion-to-scss check <input> [--debug-ast]")
	}

	input := ""
	opts := options{}
	for _, arg := range args[1:] {
		switch arg {
		case "--debug-ast":
			opts.debugAST = true
		default:
			if strings.HasPrefix(arg, "-") {
				return "", options{}, fmt.Errorf("unknown option %q", arg)
			}
			if input != "" {
				return "", options{}, fmt.Errorf("expected exactly one input")
			}
			input = arg
		}
	}
	if input == "" {
		return "", options{}, fmt.Errorf("expected input")
	}
	return input, opts, nil
}

func collectInputs(input string) ([]string, error) {
	info, err := os.Stat(input)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if !isSupported(input) {
			return nil, fmt.Errorf("unsupported extension: %s", input)
		}
		return []string{input}, nil
	}

	var files []string
	err = filepath.WalkDir(input, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if isSupported(path) {
			files = append(files, path)
		}
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
