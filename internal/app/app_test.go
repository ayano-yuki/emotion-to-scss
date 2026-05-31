package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunConvertWritesSCSS(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "button.tsx")
	outDir := filepath.Join(dir, "out")
	if err := os.WriteFile(input, []byte("const buttonStyle = css`color: red;`"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, err := Run([]string{"convert", input, "--out", outDir}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if code != ExitSuccess {
		t.Fatalf("unexpected exit code: %d", code)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "button.scss"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), ".buttonStyle") {
		t.Fatalf("unexpected output: %s", data)
	}
}

func TestRunVerifyFailsOnInvalidStyle(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "bad.ts")
	outDir := filepath.Join(dir, "out")
	if err := os.WriteFile(input, []byte("const badStyle = css`color red;`"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, err := Run([]string{"verify", input, "--out", outDir}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Run returned unexpected infrastructure error: %v", err)
	}
	if code != ExitVerificationFailed {
		t.Fatalf("unexpected exit code: %d", code)
	}
}

func TestRunInvalidArguments(t *testing.T) {
	code, err := Run([]string{"convert"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error")
	}
	if code != ExitInvalidArgs {
		t.Fatalf("unexpected exit code: %d", code)
	}
}
