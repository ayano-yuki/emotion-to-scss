package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunCheckEquivalentFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "button.ts")
	scss := filepath.Join(dir, "button.scss")

	if err := os.WriteFile(input, []byte("const buttonStyle = css`color: red;`"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scss, []byte(".buttonStyle { color: red; }"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code, err := Run([]string{"check", input}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if code != ExitSuccess {
		t.Fatalf("unexpected exit code: %d stderr=%s", code, stderr.String())
	}
}

func TestRunCheckDifferentFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "button.ts")
	scss := filepath.Join(dir, "button.scss")

	if err := os.WriteFile(input, []byte("const buttonStyle = css`color: red;`"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scss, []byte(".buttonStyle { color: blue; }"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, err := Run([]string{"check", input}, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Run returned infrastructure error: %v", err)
	}
	if code != ExitVerificationFailed {
		t.Fatalf("unexpected exit code: %d", code)
	}
}

func TestRunInvalidArgs(t *testing.T) {
	code, err := Run([]string{"verify", "button.ts"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error")
	}
	if code != ExitInvalidArgs {
		t.Fatalf("unexpected exit code: %d", code)
	}
}
