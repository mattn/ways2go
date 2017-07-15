package scanner

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

// Token is type of token.
type Token int

//go:generate stringer -type Token .

const (
	// ILLEGAL is
	ILLEGAL Token = iota
	// SPACE is
	SPACE
	// STRING is
	STRING
	// COMP is
	COMP
	// TOKEN is
	TOKEN
	// QUOTE is
	QUOTE
	// OPERATOR is
	OPERATOR
	// COMMENT is
	COMMENT
	// PAREN is
	PAREN
	// COMMA is
	COMMA
)

// Scanner hold informations for bufio.Scanner.
type Scanner struct {
	last    Token
	curr    Token
	scan    *bufio.Scanner
	quote   bool
	comment bool
	pos     int
}

func (s *Scanner) classOf(r rune) Token {
	if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
		return SPACE
	}
	if r == '=' || r == '<' || r == '>' || r == '!' {
		return COMP
	}
	if r == '+' || r == '-' || r == '*' || r == '/' {
		return OPERATOR
	}
	if r == ',' {
		return COMMA
	}
	if r == '(' || r == ')' {
		return PAREN
	}
	if '0' <= r && r <= '9' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '.' {
		return TOKEN
	}
	if r == '\'' {
		return QUOTE
	}
	return ILLEGAL
}

func (s *Scanner) splitToken(data []byte, atEOF bool) (int, []byte, error) {
	bpos := 0
	b := data
	s.curr = s.last

	var clazz Token
	for {
		r, i := utf8.DecodeRune(b)
		if i == 0 {
			break
		}
		if len(b) > 2 && ((r == '/' && b[1] == '*') || (r == '*' && b[1] == '/')) {
			clazz = COMMENT
			i++
		} else {
			clazz = s.classOf(r)
			if clazz == ILLEGAL {
				return bpos, data[:bpos], fmt.Errorf("Illegal token at %v", s.pos)
			}
		}

		if s.comment {
			if bpos > 0 && clazz == COMMENT {
				s.comment = false
			} else {
				clazz = COMMENT
			}
		} else if clazz == COMMENT {
			s.comment = true
		}

		if s.quote {
			if bpos > 0 && clazz == QUOTE {
				s.quote = false
			} else {
				clazz = QUOTE
			}
		} else if clazz == QUOTE {
			s.quote = true
		}

		if clazz != s.last {
			s.last = clazz
			break
		}
		s.pos += i
		bpos += i
		b = b[i:]
	}
	var err error
	if atEOF {
		err = io.EOF
	}
	return bpos, data[:bpos], err
}

// NewScanner return new Scanner.
func NewScanner(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	scan := &Scanner{
		scan: s,
		curr: ILLEGAL,
		last: SPACE,
	}
	s.Split(scan.splitToken)
	return scan
}

// Text return text.
func (s *Scanner) Text() string {
	return s.scan.Text()
}

// Token return token.
func (s *Scanner) Token() Token {
	return s.curr
}

// Scan return true when it is possible to scan next.
func (s *Scanner) Scan() bool {
	return s.scan.Scan()
}

// Err return error while scanning.
func (s *Scanner) Err() error {
	return s.scan.Err()
}
