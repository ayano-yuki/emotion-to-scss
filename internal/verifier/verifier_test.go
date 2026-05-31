package verifier

import (
	"testing"

	"emotion-to-scss/internal/domain"
)

func TestVerifyDynamicStyle(t *testing.T) {
	result, err := Verify(domain.Style{
		Name:      "boxStyle",
		ClassName: "boxStyle",
		CSS:       "color: __emotion_to_scss_dynamic_1__;",
		Expressions: []domain.Expression{{
			Source:      "color",
			Placeholder: "__emotion_to_scss_dynamic_1__",
		}},
	})
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if !result.Equivalent {
		t.Fatalf("expected equivalent, got reason=%q", result.Reason)
	}
}
