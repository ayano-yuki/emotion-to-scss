package verifier

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyFileEquivalent(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "button.tsx")
	scss := filepath.Join(dir, "button.scss")

	if err := os.WriteFile(input, []byte("const buttonStyle = css`color: red; &:hover { color: blue; }`"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scss, []byte(".buttonStyle { color: red; &:hover { color: blue; } }"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := VerifyFile(input)
	if err != nil {
		t.Fatalf("VerifyFile returned error: %v", err)
	}
	if !result.OK {
		t.Fatalf("expected equivalent, got reason=%q", result.Reason)
	}
}

func TestVerifyFileDifferent(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "button.tsx")
	scss := filepath.Join(dir, "button.scss")

	if err := os.WriteFile(input, []byte("const buttonStyle = css`color: red;`"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scss, []byte(".buttonStyle { color: blue; }"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := VerifyFile(input)
	if err != nil {
		t.Fatalf("VerifyFile returned error: %v", err)
	}
	if result.OK {
		t.Fatal("expected verification difference")
	}
}
