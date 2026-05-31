package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"emotion-to-scss/internal/cssast"
	"emotion-to-scss/internal/domain"
	"emotion-to-scss/internal/parser"
	"emotion-to-scss/internal/scss"
)

func VerifyFile(inputPath string) (domain.VerificationResult, error) {
	return VerifyFileWithOptions(inputPath, Options{})
}

type Options struct {
	WriteAST bool
}

func VerifyFileWithOptions(inputPath string, opts Options) (domain.VerificationResult, error) {
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
	if opts.WriteAST {
		if err := writeASTFiles(inputPath, emotionAST, scssAST); err != nil {
			return result, err
		}
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

func writeASTFiles(inputPath string, emotionAST cssast.Stylesheet, scssAST cssast.Stylesheet) error {
	base := inputPath[:len(inputPath)-len(filepath.Ext(inputPath))]
	if err := writeJSON(base+".emotion.ast.json", emotionAST); err != nil {
		return err
	}
	return writeJSON(base+".scss.ast.json", scssAST)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
