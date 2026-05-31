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

func TestVerifyFileWritesDebugAST(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "button.tsx")
	scss := filepath.Join(dir, "button.scss")

	if err := os.WriteFile(input, []byte("const buttonStyle = css`color: red;`"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scss, []byte(".buttonStyle { color: red; }"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := VerifyFileWithOptions(input, Options{WriteAST: true})
	if err != nil {
		t.Fatalf("VerifyFileWithOptions returned error: %v", err)
	}
	if !result.OK {
		t.Fatalf("expected equivalent, got reason=%q", result.Reason)
	}

	for _, path := range []string{
		filepath.Join(dir, "button.emotion.ast.json"),
		filepath.Join(dir, "button.scss.ast.json"),
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("expected AST file %s: %v", path, err)
		}
		if len(data) == 0 {
			t.Fatalf("expected AST file content for %s", path)
		}
	}
}
