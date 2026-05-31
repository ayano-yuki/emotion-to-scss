package parser

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"emotion-to-scss/internal/domain"
)

var styleStartPatterns = []*regexp.Regexp{
	regexp.MustCompile("\\b(?:export\\s+)?(?:const|let|var)\\s+([A-Za-z_$][A-Za-z0-9_$]*)\\s*(?::[^=]+)?=\\s*(?:css(?:<[^`]+>)?|styled(?:\\.[A-Za-z_$][A-Za-z0-9_$]*)?(?:<[^`]+>)?)\\s*`"),
	regexp.MustCompile("(?m)(?:^|[,{]\\s*)([A-Za-z_$][A-Za-z0-9_$]*)\\s*:\\s*css(?:<[^`]+>)?\\s*`"),
}

type ParseError struct {
	Line int
	Err  error
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("parse error at line %d: %v", e.Line, e.Err)
	}
	return fmt.Sprintf("parse error: %v", e.Err)
}

func Parse(source string) ([]domain.Style, error) {
	matches := findStyleStarts(source)
	styles := make([]domain.Style, 0, len(matches))

	for _, match := range matches {
		name := source[match.nameStart:match.nameEnd]
		css, _, err := scanTemplate(source, match.templateStart)
		if err != nil {
			return nil, &ParseError{Line: lineNumber(source, match.nameStart), Err: err}
		}
		styles = append(styles, domain.Style{
			Name:      name,
			ClassName: name,
			CSS:       css,
			Line:      lineNumber(source, match.nameStart),
		})
	}

	return styles, nil
}

func ToComparableSCSS(styles []domain.Style) string {
	var out strings.Builder
	for i, style := range styles {
		if i > 0 {
			out.WriteString("\n")
		}
		out.WriteString(".")
		out.WriteString(style.ClassName)
		out.WriteString(" {\n")
		body := strings.Trim(style.CSS, "\r\n")
		if strings.TrimSpace(body) != "" {
			for _, line := range strings.Split(body, "\n") {
				out.WriteString("  ")
				out.WriteString(strings.TrimRight(line, "\r"))
				out.WriteString("\n")
			}
		}
		out.WriteString("}\n")
	}
	return out.String()
}

type styleStart struct {
	start         int
	nameStart     int
	nameEnd       int
	templateStart int
}

func findStyleStarts(source string) []styleStart {
	var starts []styleStart
	for _, pattern := range styleStartPatterns {
		for _, match := range pattern.FindAllStringSubmatchIndex(source, -1) {
			starts = append(starts, styleStart{
				start:         match[0],
				nameStart:     match[2],
				nameEnd:       match[3],
				templateStart: match[1] - 1,
			})
		}
	}
	sort.Slice(starts, func(i, j int) bool {
		return starts[i].start < starts[j].start
	})
	return starts
}

func scanTemplate(source string, start int) (string, int, error) {
	if start < 0 || start >= len(source) || source[start] != '`' {
		return "", 0, fmt.Errorf("expected template literal")
	}

	var out strings.Builder
	dynamicIndex := 0
	for i := start + 1; i < len(source); {
		switch source[i] {
		case '\\':
			if i+1 >= len(source) {
				out.WriteByte(source[i])
				i++
				continue
			}
			out.WriteByte(source[i])
			out.WriteByte(source[i+1])
			i += 2
		case '`':
			return out.String(), i + 1, nil
		case '$':
			if i+1 < len(source) && source[i+1] == '{' {
				_, next, err := scanExpression(source, i+2)
				if err != nil {
					return "", 0, err
				}
				dynamicIndex++
				placeholder := fmt.Sprintf("__emotion_to_scss_dynamic_%d__", dynamicIndex)
				if interpolationStartsStatement(out.String()) {
					out.WriteString("--emotion-to-scss-dynamic-")
					out.WriteString(fmt.Sprint(dynamicIndex))
					out.WriteString(": ")
					out.WriteString(placeholder)
					out.WriteString(";")
				} else {
					out.WriteString(placeholder)
				}
				i = next
				continue
			}
			out.WriteByte(source[i])
			i++
		default:
			out.WriteByte(source[i])
			i++
		}
	}

	return "", 0, fmt.Errorf("unterminated template literal")
}

func scanExpression(source string, start int) (string, int, error) {
	depth := 1
	var quote byte
	templateDepth := 0

	for i := start; i < len(source); i++ {
		c := source[i]
		if quote != 0 {
			if c == '\\' {
				i++
				continue
			}
			if quote == '`' && c == '$' && i+1 < len(source) && source[i+1] == '{' {
				templateDepth++
				i++
				continue
			}
			if c == quote && templateDepth == 0 {
				quote = 0
				continue
			}
			if quote == '`' && c == '}' && templateDepth > 0 {
				templateDepth--
			}
			continue
		}

		switch c {
		case '\'', '"', '`':
			quote = c
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return source[start:i], i + 1, nil
			}
		}
	}

	return "", 0, fmt.Errorf("unterminated interpolation")
}

func interpolationStartsStatement(css string) bool {
	for i := len(css) - 1; i >= 0; i-- {
		switch css[i] {
		case ' ', '\t', '\r', '\n':
			continue
		case '{', '}', ';':
			return true
		default:
			return false
		}
	}
	return true
}

func lineNumber(source string, offset int) int {
	line := 1
	for i := 0; i < offset && i < len(source); i++ {
		if source[i] == '\n' {
			line++
		}
	}
	return line
}
