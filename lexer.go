package mapq

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

const (
	eof = iota
	ident
	illegal
	and
	or
	xor
	leftParenthesis
	rightParenthesis
)

type Token int

func (t Token) String() string {
	return tokens[t]
}

var tokens = []string{
	eof:              "EOF",
	ident:            "IDENT",
	illegal:          "ILLEGAL",
	and:              "and",
	or:               "or",
	xor:              "xor",
	leftParenthesis:  "(",
	rightParenthesis: ")",
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func CreateLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, eof, ""
			}

			panic(err)
		}

		l.pos.column++

		if r == '\n' {
			l.nextLine()
		}

		if unicode.IsSpace(r) {
			continue
		}

		if r == '(' {
			return l.pos, leftParenthesis, ""
		}

		if r == ')' {
			return l.pos, rightParenthesis, ""
		}

		if r == 'a' {
			startPos := l.pos
			l.backup()
			t, s := l.lexKeyword(and, "and")

			return startPos, t, s
		}

		if r == 'o' {
			startPos := l.pos
			l.backup()
			t, s := l.lexKeyword(or, "or")

			return startPos, t, s
		}

		if r == 'x' {
			startPos := l.pos
			l.backup()
			t, s := l.lexKeyword(xor, "xor")

			return startPos, t, s
		}

		if r == '"' {
			startPos := l.pos
			l.backup()
			t, s := l.lexString('"')

			return startPos, t, s
		}

		if r == '\'' {
			startPos := l.pos
			l.backup()
			t, s := l.lexString('\'')

			return startPos, t, s
		}

		if unicode.IsDigit(r) {
			startPos := l.pos
			l.backup()
			t, s := l.lexNumber()

			return startPos, t, s
		}

	}
}

func (l *Lexer) nextLine() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) lexKeyword(token Token, keyword string) (Token, string) {
	var lit string
	for range len(keyword) {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return illegal, lit
			}
			panic(err)
		}

		l.pos.column++
		lit = lit + string(r)
	}

	if lit != keyword {
		return illegal, lit
	}

	return token, ""
}

func (l *Lexer) lexString(starter rune) (Token, string) {
	var lit string
	var escape bool
	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return illegal, lit
			}
			panic(err)
		}

		l.pos.column++

		if r == '\\' {
			escape = true
			continue
		}

		if r == starter && !escape {
			return ident, fmt.Sprintf(`"%s"`, lit)
		}

		lit = lit + string(r)
		escape = false
	}
}

func (l *Lexer) lexNumber() (Token, string) {
	var lit string
	var dotSeen bool
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return ident, lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else if r == '.' && !dotSeen {
			lit = lit + "."
		} else {
			// scanned something not in the integer
			l.backup()
			return ident, lit
		}
	}
}
