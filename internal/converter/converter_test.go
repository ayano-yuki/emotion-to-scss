package converter

import (
	"strings"
	"testing"

	"emotion-to-scss/internal/domain"
)

func TestConvertWrapsCSSWithClassName(t *testing.T) {
	got := Convert(domain.Style{
		Name:      "buttonStyle",
		ClassName: "buttonStyle",
		CSS:       "\ncolor: red;\n",
	})

	if !strings.Contains(got, ".buttonStyle {") {
		t.Fatalf("missing class wrapper: %s", got)
	}
	if !strings.Contains(got, "  color: red;") {
		t.Fatalf("missing declaration: %s", got)
	}
}
