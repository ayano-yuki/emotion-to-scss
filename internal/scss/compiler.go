package scss

import (
	"fmt"
	"strings"
)

type Declaration struct {
	Property string
	Value    string
}

type Rule struct {
	AtRules      []string
	Selector     string
	Declarations []Declaration
}

func Compile(input string) (string, error) {
	p := parser{input: stripComments(input)}
	rules, err := p.parseStylesheet(nil, nil)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	for _, rule := range rules {
		writeRule(&out, rule, 0)
	}
	return out.String(), nil
}

func writeRule(out *strings.Builder, rule Rule, depth int) {
	indent := strings.Repeat("  ", depth)
	if len(rule.AtRules) > 0 {
		nested := rule
		nested.AtRules = nested.AtRules[1:]
		out.WriteString(indent)
		out.WriteString(rule.AtRules[0])
		out.WriteString(" {\n")
		writeRule(out, nested, depth+1)
		out.WriteString(indent)
		out.WriteString("}\n")
		return
	}

	out.WriteString(indent)
	out.WriteString(rule.Selector)
	out.WriteString(" {\n")
	for _, decl := range rule.Declarations {
		out.WriteString(indent)
		out.WriteString("  ")
		out.WriteString(decl.Property)
		out.WriteString(": ")
		out.WriteString(decl.Value)
		out.WriteString(";\n")
	}
	out.WriteString(indent)
	out.WriteString("}\n")
}

type parser struct {
	input string
	pos   int
}

func (p *parser) parseStylesheet(parents []string, atRules []string) ([]Rule, error) {
	var rules []Rule
	for {
		p.skipSpace()
		if p.eof() || p.peek() == '}' {
			return rules, nil
		}
		head, delimiter := p.readUntilTopLevel("{};")
		head = strings.TrimSpace(head)
		if head == "" {
			if delimiter == ';' {
				p.pos++
				continue
			}
			return nil, fmt.Errorf("empty rule")
		}
		if delimiter != '{' {
			return nil, fmt.Errorf("expected block for %q", head)
		}
		p.pos++

		if strings.HasPrefix(head, "@") {
			children, err := p.parseBlock(parents, append(atRules, head))
			if err != nil {
				return nil, err
			}
			rules = append(rules, children...)
			continue
		}

		selectors := combineSelectors(parents, splitSelectors(head))
		children, err := p.parseBlock(selectors, atRules)
		if err != nil {
			return nil, err
		}
		rules = append(rules, children...)
	}
}

func (p *parser) parseBlock(selectors []string, atRules []string) ([]Rule, error) {
	var declarations []Declaration
	var nested []Rule

	for {
		p.skipSpace()
		if p.eof() {
			return nil, fmt.Errorf("unterminated block")
		}
		if p.peek() == '}' {
			p.pos++
			break
		}

		head, delimiter := p.readUntilTopLevel("{};")
		head = strings.TrimSpace(head)
		if head == "" {
			if delimiter == ';' {
				p.pos++
				continue
			}
			return nil, fmt.Errorf("empty statement")
		}

		switch delimiter {
		case ';':
			p.pos++
			decl, err := parseDeclaration(head)
			if err != nil {
				return nil, err
			}
			declarations = append(declarations, decl)
		case '}':
			decl, err := parseDeclaration(head)
			if err != nil {
				return nil, err
			}
			declarations = append(declarations, decl)
			p.pos++
			return buildRules(selectors, atRules, declarations, nested)
		case '{':
			p.pos++
			if strings.HasPrefix(head, "@") {
				children, err := p.parseBlock(selectors, append(atRules, head))
				if err != nil {
					return nil, err
				}
				nested = append(nested, children...)
				continue
			}
			childSelectors := combineSelectors(selectors, splitSelectors(head))
			children, err := p.parseBlock(childSelectors, atRules)
			if err != nil {
				return nil, err
			}
			nested = append(nested, children...)
		default:
			return nil, fmt.Errorf("unexpected end of statement")
		}
	}

	return buildRules(selectors, atRules, declarations, nested)
}

func buildRules(selectors []string, atRules []string, declarations []Declaration, nested []Rule) ([]Rule, error) {
	var rules []Rule
	if len(declarations) > 0 {
		if len(selectors) == 0 {
			return nil, fmt.Errorf("declarations require a selector")
		}
		for _, selector := range selectors {
			rules = append(rules, Rule{
				AtRules:      append([]string(nil), atRules...),
				Selector:     selector,
				Declarations: append([]Declaration(nil), declarations...),
			})
		}
	}
	rules = append(rules, nested...)
	return rules, nil
}

func parseDeclaration(statement string) (Declaration, error) {
	idx := strings.Index(statement, ":")
	if idx <= 0 {
		return Declaration{}, fmt.Errorf("invalid declaration %q", statement)
	}
	property := strings.TrimSpace(statement[:idx])
	value := strings.TrimSpace(statement[idx+1:])
	if property == "" || value == "" {
		return Declaration{}, fmt.Errorf("invalid declaration %q", statement)
	}
	return Declaration{Property: property, Value: value}, nil
}

func combineSelectors(parents, children []string) []string {
	if len(parents) == 0 {
		return children
	}
	var out []string
	for _, parent := range parents {
		for _, child := range children {
			if strings.Contains(child, "&") {
				out = append(out, strings.ReplaceAll(child, "&", parent))
			} else {
				out = append(out, parent+" "+child)
			}
		}
	}
	return out
}

func splitSelectors(selector string) []string {
	parts := strings.Split(selector, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func (p *parser) readUntilTopLevel(delimiters string) (string, byte) {
	start := p.pos
	var quote byte
	parenDepth := 0
	bracketDepth := 0
	for !p.eof() {
		c := p.peek()
		if quote != 0 {
			if c == '\\' {
				p.pos += 2
				continue
			}
			if c == quote {
				quote = 0
			}
			p.pos++
			continue
		}
		switch c {
		case '\'', '"':
			quote = c
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		default:
			if parenDepth == 0 && bracketDepth == 0 && strings.ContainsRune(delimiters, rune(c)) {
				return p.input[start:p.pos], c
			}
		}
		p.pos++
	}
	return p.input[start:p.pos], 0
}

func (p *parser) skipSpace() {
	for !p.eof() {
		switch p.peek() {
		case ' ', '\t', '\r', '\n':
			p.pos++
		default:
			return
		}
	}
}

func (p *parser) eof() bool {
	return p.pos >= len(p.input)
}

func (p *parser) peek() byte {
	return p.input[p.pos]
}

func stripComments(input string) string {
	var out strings.Builder
	for i := 0; i < len(input); {
		if i+1 < len(input) && input[i] == '/' && input[i+1] == '*' {
			i += 2
			for i+1 < len(input) && !(input[i] == '*' && input[i+1] == '/') {
				if input[i] == '\n' {
					out.WriteByte('\n')
				}
				i++
			}
			if i+1 < len(input) {
				i += 2
			}
			continue
		}
		out.WriteByte(input[i])
		i++
	}
	return out.String()
}
