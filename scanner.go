package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Token int

const (
	ILLEGAL Token = iota
	SPACE
	STRING
	OPERATOR
	TOKEN
)

type Scanner struct {
	last Token
	curr Token
	s    *bufio.Scanner
}

func classOf(r rune) Token {
	if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
		return SPACE
	}
	if r == '=' || r == '<' || r == '>' || r == '!' {
		return OPERATOR
	}
	if '0' <= r && r <= '9' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'A') {
		return TOKEN
	}
	if r == '\'' {
		return STRING
	}
	return ILLEGAL
}

func (s *Scanner) splitToken(data []byte, atEOF bool) (int, []byte, error) {
	bpos := 0
	b := data
	if s.last == ILLEGAL {
		s.curr = SPACE
	} else {
		s.curr = s.last
	}
	for {
		r, i := utf8.DecodeRune(b)
		if i == 0 {
			break
		}
		clazz := classOf(r)
		if s.last == ILLEGAL {
			s.last = clazz
		} else if clazz != s.last {
			s.last = clazz
			break
		}
		bpos += i
		b = b[i:]
	}
	var err error
	if atEOF {
		err = io.EOF
	}
	return bpos, data[:bpos], err
}

func NewScanner(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	scan := &Scanner{
		s:    s,
		curr: ILLEGAL,
		last: ILLEGAL,
	}
	s.Split(scan.splitToken)
	return scan
}

func (s *Scanner) Text() string {
	return s.s.Text()
}

func (s *Scanner) Token() Token {
	return s.curr
}

func (s *Scanner) Scan() bool {
	return s.s.Scan()
}

func main() {
	s := `
	select * from foo where id = 3
	`
	scan := NewScanner(strings.NewReader(s))
	for scan.Scan() {
		fmt.Printf("%v: %q\n", scan.Token(), scan.Text())
	}
}
