package cssast

import (
	"fmt"
	"reflect"
	"strings"
)

type Stylesheet struct {
	Nodes []Node
}

type Node struct {
	AtRule *AtRule
	Rule   *Rule
}

type AtRule struct {
	Prelude  string
	Children []Node
}

type Rule struct {
	Selector     string
	Declarations []Declaration
}

type Declaration struct {
	Property string
	Value    string
}

func Parse(input string) (Stylesheet, error) {
	p := parser{input: stripComments(input)}
	nodes, err := p.parseNodes(false)
	if err != nil {
		return Stylesheet{}, err
	}
	return Stylesheet{Nodes: nodes}, nil
}

func Equal(a, b Stylesheet) bool {
	return reflect.DeepEqual(a, b)
}

type parser struct {
	input string
	pos   int
}

func (p *parser) parseNodes(stopOnBrace bool) ([]Node, error) {
	var nodes []Node
	for {
		p.skipSpace()
		if p.eof() {
			if stopOnBrace {
				return nil, fmt.Errorf("unterminated block")
			}
			return nodes, nil
		}
		if p.peek() == '}' {
			if !stopOnBrace {
				return nil, fmt.Errorf("unexpected closing brace")
			}
			p.pos++
			return nodes, nil
		}

		head, delimiter := p.readUntilTopLevel("{}")
		head = normalizeSpace(head)
		if head == "" || delimiter != '{' {
			return nil, fmt.Errorf("expected CSS block")
		}
		p.pos++

		if strings.HasPrefix(head, "@") {
			children, err := p.parseNodes(true)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, Node{AtRule: &AtRule{Prelude: head, Children: children}})
			continue
		}

		declarations, err := p.parseDeclarations()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, Node{Rule: &Rule{Selector: head, Declarations: declarations}})
	}
}

func (p *parser) parseDeclarations() ([]Declaration, error) {
	var declarations []Declaration
	for {
		p.skipSpace()
		if p.eof() {
			return nil, fmt.Errorf("unterminated rule")
		}
		if p.peek() == '}' {
			p.pos++
			return declarations, nil
		}

		statement, delimiter := p.readUntilTopLevel(";}")
		statement = strings.TrimSpace(statement)
		if statement == "" {
			if delimiter == ';' {
				p.pos++
				continue
			}
			if delimiter == '}' {
				p.pos++
				return declarations, nil
			}
		}
		decl, err := parseDeclaration(statement)
		if err != nil {
			return nil, err
		}
		declarations = append(declarations, decl)
		if delimiter == ';' {
			p.pos++
			continue
		}
		if delimiter == '}' {
			p.pos++
			return declarations, nil
		}
		return nil, fmt.Errorf("unterminated declaration")
	}
}

func parseDeclaration(statement string) (Declaration, error) {
	idx := strings.Index(statement, ":")
	if idx <= 0 {
		return Declaration{}, fmt.Errorf("invalid declaration %q", statement)
	}
	property := normalizeSpace(statement[:idx])
	value := normalizeSpace(statement[idx+1:])
	if property == "" || value == "" {
		return Declaration{}, fmt.Errorf("invalid declaration %q", statement)
	}
	return Declaration{Property: property, Value: value}, nil
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

func normalizeSpace(input string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(input)), " ")
}

func stripComments(input string) string {
	var out strings.Builder
	for i := 0; i < len(input); {
		if i+1 < len(input) && input[i] == '/' && input[i+1] == '*' {
			i += 2
			for i+1 < len(input) && !(input[i] == '*' && input[i+1] == '/') {
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
