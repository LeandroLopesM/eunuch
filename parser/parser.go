package parser

import (
	"fmt"
	"slices"
	"strconv"
	"unicode"

	. "github.com/leandrolopesm/eunuch/core"
)

type parser struct {
	pos   int
	input []rune

	file string
}

func (p *parser) currLinePos() (int, int) {
	lastNewline := 1
	lastChar := 0

	for pos, char := range p.input {
		if char == '\n' {
			if pos+1 == len(p.input) || pos >= p.pos {
				break
			}

			lastChar = pos
			lastNewline++
		}
	}

	return lastNewline, p.pos - lastChar
}

func (p *parser) currPos() Position {
	line, char := p.currLinePos()
	return Position{
		Line: line,
		Char: char,
		File: p.file,
	}
}

func (p *parser) err(msg string, args ...any) error {
	lastNL, char := p.currLinePos()

	return fmt.Errorf("%s:%d:%d: %s", p.file, lastNL, char, fmt.Sprintf(msg, args...))
}

func (p *parser) next() rune {
	if p.pos >= len(p.input) {
		return rune(0)
	}

	return p.input[p.pos]
}

func (p *parser) rewind(howMuch int) {
	p.pos -= howMuch
}

func (p *parser) startsWith(prefix string) bool {
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

func (p *parser) expect(what string) {
	if p.startsWith(what) {
		p.pos += len(what)
	} else {
		panic(fmt.Sprintf("parser:expect failed, expected %s, got %s", what, string(p.input[p.pos:])))
	}
}

func (p *parser) eof() bool {
	return p.pos >= len(p.input)
}

func (p *parser) consume() rune {
	p.pos++
	return p.input[p.pos-1]
}

func (p *parser) consumeWhile(condition func(rune) bool) string {
	var buffer []rune

	for !p.eof() && condition(p.next()) {
		buffer = append(buffer, p.consume())
	}

	return string(buffer)
}

func (p *parser) consumeWhitespace() {
	p.consumeWhile(unicode.IsSpace)
}

func (p *parser) parseName() string {
	return p.consumeWhile(func(c rune) bool { return unicode.IsLetter(c) || unicode.IsDigit(c) })
}

/* =============================================================
 * SCHEME PARSER
 * =============================================================
 * Very rough EBNF Grammar
 * =============================================================
 * File 	  := (Expression)*
 *      	  | ';' (Any)* '\n';
 * Expression := Scheme
 *			  | '#' Tag
 *			  | Numeric
 *			  | String
 *			  | Symbol
 *			  | Scheme;
 *
 * Scheme	  := '(' IDENT (Expression)* ')';
 * Tag	  	  := '\\' CHAR
 * 			  |  't' | 'f';
 *
 * Numeric	  := Float | Int;
 * Float	  := '.' NUMBER | NUMBER '.' NUMBER | NUMBER '.';
 * Int		  := NUMBER;
 * ============================================================= */

/// ParseFile parses a file given the grammar described above
func ParseFile(fileName, fileContent string) ([]Value, error) {
	parser := parser{
		pos:   0,
		input: []rune(fileContent),
		file:  fileName,
	}

	var file []Value
	for !parser.eof() {
		parser.consumeWhitespace()

		var (
			unit Value
			err  error
		)

		if unit, err = parser.parseExpression(); err != nil {
			return []Value{}, err
		}
		file = append(file, unit)
	}

	return file, nil
}

func (p *parser) parseExpression() (Value, error) {
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

func (p *parser) parseQuote() (Value, error) {
	p.consume()

	var (
		unit Value
		err  error
	)

	if unit, err = p.parseExpression(); err != nil {
		return Value{}, err
	}

	return SchemeVal(Scheme{
			Args:     slices.Concat([]Value{SymbolVal("quote"), unit}),
			Position: p.currPos(),
	}), nil
}

func (p *parser) parseString() (Value, error) {
	_ = p.consume()
	start := p.pos
	str := p.consumeWhile(func(c rune) bool { return c != '"' })

	if p.eof() {
		p.pos = start
		return Value{}, p.err("Unclosed string starts here")
	}

	p.consume()

	return StringVal(str), nil
}

func (p *parser) parseScheme() (Value, error) {
	p.expect("(")

	start := p.pos
	out := Scheme{
		Position: p.currPos(),
		Args: []Value{},
	}

	for {
		if p.eof() {
			p.pos = start
			return Value{}, p.err("Unclosed scheme here")
		} else if p.next() == ')' {
			_ = p.consume()
			break
		}

		var (
			unit Value
			err error
		)

		if unit, err = p.parseExpression(); err != nil {
			return Value{}, err
		}

		p.consumeWhitespace()
		out.Args = append(out.Args, unit)
	}

	return SchemeVal(out), nil
}

func isHex(c rune) bool {
	return unicode.IsDigit(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func (p *parser) parseInt() (Value, error) {
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
		return Value{}, p.err("invalid integer literal %s", buffer)
	} else {
		return IntegerVal(val), nil
	}
}

func (p *parser) parseFloat() (Value, error) {
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
		return Value{}, p.err("Invalid float literal %s", string(buffer))
	} else {
		return FloatVal(val), nil
	}
}

func (p *parser) parseLiteral() (Value, error) {
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

func (p *parser) parseSymbol() (Value, error) {
	return SymbolVal(p.consumeWhile(isSymbolChar)), nil
}

func isSymbolChar(c rune) bool {
	return c != '(' && c != ')' && !unicode.IsSpace(c)
}

func (p *parser) parseBool() (Value, error) {
	val := p.consume()

	return BoolVal(val == 't'), nil
}

func (p *parser) parseChar() (Value, error) {
	p.expect("\\")

	charName := p.parseName()

	switch charName {
	case "space":
		return CharVal(' '), nil
	case "tab":
		return CharVal('\t'), nil
	case "\\n":
		return CharVal('\n'), nil
	default:
		if len(charName) > 1 {
			return Value{}, p.err("Unknown character literal %v", charName)
		}

		return CharVal(([]rune(charName))[0]), nil
	}
}
