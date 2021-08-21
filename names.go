//go:generate stringer -type=Kind
package names

import (
	"unicode"
	"unicode/utf8"
)

type Kind int

const (
	EOF Kind = iota
	Invalid
	Word
	Symbol
)

type Flags int

const (
	Uppercase Flags = iota + 1
)

type Token struct {
	Kind  Kind
	Value string
	Flags Flags
}

func (tkn Token) Uppercase() bool {
	return tkn.Flags&Uppercase == Uppercase
}

type Tokeniser struct {
	Input      string
	Characters string
	pos        int
	width      int
}

func (tok *Tokeniser) Token() Token {
	start := tok.pos
	r, allowed := tok.next()
	if !allowed {
		return Token{Kind: Invalid, Value: string(r)}
	}
	switch {
	case r == eof:
		return Token{Kind: EOF}
	case unicode.IsLetter(r):
		return tok.readWord(r, start)
	default:
		return Token{Kind: Symbol, Value: string(r)}
	}
}

func (tok *Tokeniser) allowed(r rune) bool {
	for _, ch := range tok.Characters {
		if ch == r {
			return true
		}
	}
	return false
}

const eof rune = 0

func (tok *Tokeniser) next() (r rune, allowed bool) {
	if tok.pos >= len(tok.Input) {
		tok.width = 0
		return eof, true
	}
	r, tok.width = utf8.DecodeRuneInString(tok.Input[tok.pos:])
	tok.pos += tok.width
	return r, tok.allowed(r)
}

func (tok *Tokeniser) peek() (rune, bool) {
	pos := tok.pos
	width := tok.width
	r, allowed := tok.next()
	tok.pos = pos
	tok.width = width
	return r, allowed
}

func (err *Tokeniser) invalid(value string) Token {
	return Token{Kind: Invalid, Value: value}
}

func flags(uppercase bool) Flags {
	var flags Flags
	if uppercase {
		flags |= Uppercase
	}
	return flags
}

func (tok *Tokeniser) readWord(r rune, start int) Token {
	uppercaseFlag := true
	count := 1

	peek, allowed := tok.peek()
	if !unicode.IsLetter(peek) {
		return Token{Kind: Word, Value: string(r), Flags: flags(uppercaseFlag)}
	}
	if !allowed {
		return tok.invalid(string(r))
	}

	for {
		r, allowed := tok.next()
		if !allowed {
			return tok.invalid(string(r))
		}
		switch {
		case unicode.IsLetter(r):
			uppercaseRune := unicode.IsUpper(r)
			if uppercaseFlag && !uppercaseRune {
				uppercaseFlag = false
			}
			peek, allowed := tok.peek()
			if !allowed {
				return tok.invalid(string(peek))
			}
			uppercasePeek := unicode.IsUpper(peek)
			if peek != eof && uppercaseFlag && count > 2 && !uppercasePeek {
				tok.pos -= tok.width
				return Token{Kind: Word, Value: tok.Input[start:tok.pos], Flags: flags(uppercaseFlag)}
			} else if !uppercaseFlag && uppercasePeek {
				return Token{Kind: Word, Value: tok.Input[start:tok.pos], Flags: flags(uppercaseFlag)}
			}
			count++
		default:
			tok.pos -= tok.width
			return Token{Kind: Word, Value: tok.Input[start:tok.pos], Flags: flags(uppercaseFlag)}
		}
	}
}
