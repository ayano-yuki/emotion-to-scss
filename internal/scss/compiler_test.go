package scss

import (
	"strings"
	"testing"
)

func TestCompileNestedSelector(t *testing.T) {
	got, err := Compile(`.buttonStyle {
  color: red;

  &:hover {
    color: blue;
  }
}`)
	if err != nil {
		t.Fatalf("Compile returned error: %v", err)
	}
	if !strings.Contains(got, ".buttonStyle:hover") {
		t.Fatalf("expected nested selector, got:\n%s", got)
	}
	if !strings.Contains(got, "color: blue;") {
		t.Fatalf("expected nested declaration, got:\n%s", got)
	}
}

func TestCompileAtRule(t *testing.T) {
	got, err := Compile(`.buttonStyle {
  @media (min-width: 640px) {
    color: red;
  }
}`)
	if err != nil {
		t.Fatalf("Compile returned error: %v", err)
	}
	if !strings.Contains(got, "@media (min-width: 640px)") {
		t.Fatalf("expected media rule, got:\n%s", got)
	}
	if !strings.Contains(got, ".buttonStyle") {
		t.Fatalf("expected selector, got:\n%s", got)
	}
}

func TestCompileInvalidDeclaration(t *testing.T) {
	_, err := Compile(".x { color red; }")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompileDeclarationWithoutTrailingSemicolon(t *testing.T) {
	got, err := Compile(".x { color: red }")
	if err != nil {
		t.Fatalf("Compile returned error: %v", err)
	}
	if !strings.Contains(got, "color: red;") {
		t.Fatalf("expected declaration, got:\n%s", got)
	}
}
