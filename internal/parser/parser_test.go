package parser

import (
	"strings"
	"testing"
)

func TestParseEmotionTemplates(t *testing.T) {
	source := `import { css } from "@emotion/react"
import styled from "@emotion/styled"

export const buttonStyle = css` + "`" + `
  color: red;

  &:hover {
    color: blue;
  }
` + "`" + `

export const Button = styled.button<ButtonProps>` + "`" + `
  display: inline-flex;
` + "`"

	styles, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(styles) != 2 {
		t.Fatalf("expected 2 styles, got %d", len(styles))
	}
	if styles[0].Name != "buttonStyle" || styles[1].Name != "Button" {
		t.Fatalf("unexpected style names: %s, %s", styles[0].Name, styles[1].Name)
	}
}

func TestToComparableSCSS(t *testing.T) {
	styles, err := Parse("const buttonStyle = css`color: red;`")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	got := ToComparableSCSS(styles)
	if !strings.Contains(got, ".buttonStyle") {
		t.Fatalf("missing class selector: %s", got)
	}
	if !strings.Contains(got, "color: red;") {
		t.Fatalf("missing declaration: %s", got)
	}
}

func TestParseUnterminatedTemplate(t *testing.T) {
	_, err := Parse("const badStyle = css`color: red;")
	if err == nil {
		t.Fatal("expected parse error")
	}
}
