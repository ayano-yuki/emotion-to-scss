package verifier

import (
	"fmt"
	"os"
	"path/filepath"

	"emotion-to-scss/internal/cssast"
	"emotion-to-scss/internal/domain"
	"emotion-to-scss/internal/parser"
	"emotion-to-scss/internal/scss"
)

func VerifyFile(inputPath string) (domain.VerificationResult, error) {
	result := domain.VerificationResult{
		InputPath: inputPath,
		SCSSPath:  matchingSCSSPath(inputPath),
	}

	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return result, err
	}
	scssBytes, err := os.ReadFile(result.SCSSPath)
	if err != nil {
		return result, err
	}

	styles, err := parser.Parse(string(inputBytes))
	if err != nil {
		return result, err
	}
	if len(styles) == 0 {
		return result, fmt.Errorf("no Emotion CSS templates found")
	}

	emotionCSS, err := scss.Compile(parser.ToComparableSCSS(styles))
	if err != nil {
		return result, fmt.Errorf("compile Emotion CSS: %w", err)
	}
	scssCSS, err := scss.Compile(string(scssBytes))
	if err != nil {
		return result, fmt.Errorf("compile SCSS: %w", err)
	}

	emotionAST, err := cssast.Parse(emotionCSS)
	if err != nil {
		return result, fmt.Errorf("parse Emotion CSS: %w", err)
	}
	scssAST, err := cssast.Parse(scssCSS)
	if err != nil {
		return result, fmt.Errorf("parse SCSS CSS: %w", err)
	}
	if !cssast.Equal(emotionAST, scssAST) {
		result.Reason = "normalized CSS AST differs"
		return result, nil
	}

	result.OK = true
	return result, nil
}

func matchingSCSSPath(inputPath string) string {
	ext := filepath.Ext(inputPath)
	return inputPath[:len(inputPath)-len(ext)] + ".scss"
}
