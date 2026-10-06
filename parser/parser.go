package parser

import (
	"fmt"
	"slices"
	"strconv"
	"unicode"

	. "github.com/leandrolopesm/eunuch/core"
)

type Parser struct {
	pos   int
	input []rune

	file string
}

func (p *Parser) currLinePos() (int, int) {
	nl := 1
	lastPos := 0

	for pos, char := range p.input {
		if char == '\n' {
			if pos+1 == len(p.input) || pos >= p.pos {
				break
			}

			lastPos = pos
			nl++
		}
	}

	return nl, p.pos - lastPos
}

func (p *Parser) currPos() Position {
	line, char := p.currLinePos()
	return Position{
		Line: line,
		Char: char,
		File: p.file,
	}
}

func (p *Parser) err(msg string, args ...any) error {
	lastNL, char := p.currLinePos()

	return fmt.Errorf("%s:%d:%d: %s", p.file, lastNL, char, fmt.Sprintf(msg, args...))
}

func (p *Parser) next() rune {
	if p.pos >= len(p.input) {
		return rune(0)
	}

	return p.input[p.pos]
}

func (p *Parser) rewind(howMuch int) {
	p.pos -= howMuch
}

func (p *Parser) StartsWith(prefix string) bool {
	trueStr := string(p.input[p.pos:])
	if len(trueStr) < len(prefix) {
		return false
	}

	for i := range prefix {
		if prefix[i] != trueStr[i] {
			return false
		}
	}

	return true
}

func (p *Parser) expect(what string) {
	if p.StartsWith(what) {
		p.pos += len(what)
	} else {
		panic(fmt.Sprintf("Parser:expect failed, expected %s, got %s", what, string(p.input[p.pos:])))
	}
}

func (p *Parser) EOF() bool {
	return p.pos >= len(p.input)
}

func (p *Parser) consume() rune {
	p.pos += 1
	return p.input[p.pos-1]
}

func (p *Parser) consumeWhile(condition func(rune) bool) string {
	var buffer []rune

	for !p.EOF() && condition(p.next()) {
		buffer = append(buffer, p.consume())
	}

	return string(buffer)
}

func (p *Parser) consumeWhitespace() {
	p.consumeWhile(unicode.IsSpace)
}

func (p *Parser) parseName() string {
	return p.consumeWhile(func(c rune) bool { return unicode.IsLetter(c) || unicode.IsDigit(c) })
}

// ============================================================
// SCHEME PARSER
// ============================================================
// Very rough EBNF Grammar
// ============================================================
// File 	  := (Expression)*
//      	  | ';' (Any)* '\n';
// Expression := Scheme
//			  | '#' Tag
//			  | Numeric
//			  | String
//			  | Symbol
//			  | Scheme;
//
// Scheme	  := '(' IDENT (Expression)* ')';
// Tag	  	  := '\\' CHAR
// 			  |  't' | 'f';
//
// Numeric	  := Float | Int;
// Float	  := '.' NUMBER | NUMBER '.' NUMBER | NUMBER '.';
// Int		  := NUMBER;
// ============================================================

func ParseFile(fileName, fileContent string) ([]Unit, error) {
	parser := Parser{
		pos:   0,
		input: []rune(fileContent),
		file:  fileName,
	}

	var file []Unit
	for !parser.EOF() {
		parser.consumeWhitespace()

		if val, err := parser.parseExpression(); err != nil {
			return []Unit{}, err
		} else {
			file = append(file, val)
		}
	}

	return file, nil
}

func (p *Parser) parseExpression() (Unit, error) {
tryAgain:
	switch p.next() {
	case ';':
		p.consumeWhile(func(a rune) bool { return a != '\n' })
		p.expect("\n")
		goto tryAgain

	case '#':
		return p.parseLiteral()

	case '.':
		return p.parseFloat()

	case '(':
		return p.parseScheme()

	case '"':
		return p.parseString()

	case '\'':
		return p.parseQuote()

	default:
		switch {
		case unicode.IsDigit(p.next()):
			return p.parseInt()

		case unicode.IsSpace(p.next()):
			p.consumeWhitespace()
			goto tryAgain

		default:
			return p.parseSymbol()
		}
	}
}

func (p *Parser) parseQuote() (Unit, error) {
	p.consume()

	if val, err := p.parseExpression(); err != nil {
		return Unit{}, err
	} else {
		return MkScheme(Scheme{
			Args:     slices.Concat([]Unit{MkSymbol("quote"), val}),
			Position: p.currPos(),
		}), nil
	}
}

func (p *Parser) parseString() (Unit, error) {
	_ = p.consume()
	start := p.pos
	str := p.consumeWhile(func(c rune) bool { return c != '"' })

	if p.EOF() {
		p.pos = start
		return MkString(""), p.err("Unclosed string starts here")
	} else {
		p.consume()
	}

	return MkString(str), nil
}

func (p *Parser) parseScheme() (Unit, error) {
	p.expect("(")
	out := Scheme{}

	out.Position = p.currPos()

	for {
		if p.next() == ')' {
			_ = p.consume()
			break
		}

		if val, err := p.parseExpression(); err != nil {
			return MkScheme(Scheme{}), err
		} else {
			out.Args = append(out.Args, val)
		}
	}

	return MkScheme(out), nil
}

func isHex(c rune) bool {
	return unicode.IsDigit(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func (p *Parser) parseInt() (Unit, error) {
	radix := 10
	if p.next() == 'x' {
		radix = 16
	}

	buffer := p.consumeWhile(isHex)
	if p.next() == '.' {
		p.rewind(len(buffer))
		return p.parseFloat()
	}

	if val, err := strconv.ParseInt(buffer, radix, 64); err != nil {
		return MkInt(0), p.err("Invalid integer literal %s", string(buffer))
	} else {
		return MkInt(val), nil
	}
}

func (p *Parser) parseFloat() (Unit, error) {
	var buffer []rune

	if p.next() == '.' {
		buffer = append(buffer, '0') // predicate the string so .123 is 0.123
	} else {
		buffer = []rune(p.consumeWhile(unicode.IsDigit))
		p.expect(".")
	}

	if !unicode.IsDigit(p.next()) {
		buffer = append(buffer, '0') // suffix the string so 123. is 123.0
	} else {
		buffer = slices.Concat(buffer, []rune(p.consumeWhile(unicode.IsDigit)))
	}

	if val, err := strconv.ParseFloat(string(buffer), 64); err != nil {
		return MkFloat(0.), p.err("Invalid float literal %s", string(buffer))
	} else {
		return MkFloat(val), nil
	}

}

func (p *Parser) parseLiteral() (Unit, error) {
	p.expect("#")

	switch p.next() {
	case '\\':
		return p.parseChar()
	case 't', 'f':
		return p.parseBool()
	}

	p.rewind(1)
	return p.parseSymbol()
}

func (p *Parser) parseSymbol() (Unit, error) {
	return MkSymbol(p.consumeWhile(isSymbolChar)), nil
}

func isSymbolChar(c rune) bool {
	return c != '(' && c != ')' && !unicode.IsSpace(c)
}

func (p *Parser) parseBool() (Unit, error) {
	val := p.consume()

	return MkBool(val == 't'), nil
}

func (p *Parser) parseChar() (Unit, error) {
	p.expect("\\")

	charName := p.parseName()

	switch charName {
	case "space":
		return MkChar(' '), nil
	case "tab":
		return MkChar('\t'), nil
	case "\\n":
		return MkChar('\n'), nil
	default:
		if len(charName) > 1 {
			return MkChar(0), p.err("Unknown character literal %v", charName)
		}

		return MkChar(([]rune(charName))[0]), nil
	}
}
