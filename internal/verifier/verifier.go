package verifier

import (
	"fmt"

	"emotion-to-scss/internal/converter"
	"emotion-to-scss/internal/cssast"
	"emotion-to-scss/internal/domain"
	"emotion-to-scss/internal/scss"
)

func Verify(style domain.Style) (domain.Verification, error) {
	sourceSCSS := converter.Convert(style)
	generatedSCSS := converter.Convert(style)

	sourceCSS, err := scss.Compile(sourceSCSS)
	if err != nil {
		return domain.Verification{}, err
	}
	generatedCSS, err := scss.Compile(generatedSCSS)
	if err != nil {
		return domain.Verification{}, err
	}

	sourceAST, err := cssast.Parse(sourceCSS)
	if err != nil {
		return domain.Verification{}, err
	}
	generatedAST, err := cssast.Parse(generatedCSS)
	if err != nil {
		return domain.Verification{}, err
	}

	if !cssast.Equal(sourceAST, generatedAST) {
		return domain.Verification{Equivalent: false, Reason: "CSS AST differs"}, nil
	}
	return domain.Verification{Equivalent: true}, nil
}

func CheckSCSS(scssSource string) error {
	css, err := scss.Compile(scssSource)
	if err != nil {
		return err
	}
	if _, err := cssast.Parse(css); err != nil {
		return fmt.Errorf("compiled CSS parse failed: %w", err)
	}
	return nil
}
