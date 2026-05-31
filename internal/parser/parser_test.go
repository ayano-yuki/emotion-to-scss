package parser

import "testing"

func TestParseStaticCSSAndStyledTemplate(t *testing.T) {
	source := `import styled from "@emotion/styled"

const buttonStyle = styled.button` + "`" + `
  color: red;

  &:hover {
    color: blue;
  }
` + "`" + `
`

	styles, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(styles) != 1 {
		t.Fatalf("expected 1 style, got %d", len(styles))
	}
	if styles[0].Name != "buttonStyle" {
		t.Fatalf("unexpected style name: %s", styles[0].Name)
	}
	if styles[0].Line != 3 {
		t.Fatalf("unexpected line: %d", styles[0].Line)
	}
}

func TestParseDynamicInterpolation(t *testing.T) {
	source := `const boxStyle = css` + "`" + `
  color: ${color};
  width: ${({ size }) => size}px;
` + "`"

	styles, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(styles) != 1 {
		t.Fatalf("expected 1 style, got %d", len(styles))
	}
	if got := len(styles[0].Expressions); got != 2 {
		t.Fatalf("expected 2 expressions, got %d", got)
	}
	if styles[0].Expressions[0].Source != "color" {
		t.Fatalf("unexpected expression: %q", styles[0].Expressions[0].Source)
	}
	if styles[0].Expressions[1].Placeholder != "__emotion_to_scss_dynamic_2__" {
		t.Fatalf("unexpected placeholder: %q", styles[0].Expressions[1].Placeholder)
	}
}

func TestParseStandaloneDynamicInterpolation(t *testing.T) {
	source := "const boxStyle = css`${baseStyle}`"

	styles, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(styles) != 1 {
		t.Fatalf("expected 1 style, got %d", len(styles))
	}
	want := "--emotion-to-scss-dynamic-1: __emotion_to_scss_dynamic_1__"
	if styles[0].CSS != want {
		t.Fatalf("unexpected CSS: %q", styles[0].CSS)
	}
}

func TestParseUnterminatedTemplate(t *testing.T) {
	_, err := Parse("const badStyle = css`color: red;")
	if err == nil {
		t.Fatal("expected parse error")
	}
}
