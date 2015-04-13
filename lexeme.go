package blast

import (
	"fmt"
	"log"
)

type lexemeType int

const (
	lexemeTypeEOL = iota
	lexemeTypeNum
	lexemeTypeTab
	lexemeTypeComma
	lexemeTypeSpace
	lexemeTypeParen
	lexemeTypeQuote
	lexemeTypeCharacter
	lexemeTypeOperator
	lexemeTypeDecimalDigit
	lexemeTypeUnidentified
)

var lexemeIdentifiers = map[lexemeType]string{
	lexemeTypeEOL:          "\n",
	lexemeTypeNum:          "0123456789",
	lexemeTypeTab:          "\t",
	lexemeTypeParen:        "()",
	lexemeTypeComma:        ",",
	lexemeTypeQuote:        "\"",
	lexemeTypeSpace:        " ",
	lexemeTypeCharacter:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	lexemeTypeOperator:     "+-<>*/=&|",
	lexemeTypeDecimalDigit: ".",
}

// lexemeIdentificationError is an error
// that occurs when a character cannot
// be identified as any type of lexeme
type lexemeIdentificationError struct {
	pos  int
	text rune
}

// Error returns a string with information
// about a lexemeIdentificationError
func (l *lexemeIdentificationError) Error() string {
	return fmt.Sprintf(
		"Could not identify item \"%v\" at position %v", string(l.text), l.pos,
	)
}

// lexeme is a struct that
// stores a rune, position,
// and is a type of lexeme
type lexeme struct {
	pos  int
	text rune
	t    lexemeType
}

// lexemeStream is a struct
// that stores a slice of
// lexemes
type lexemeStream struct {
	lexemes []*lexeme
	size    int
	start   int
	end     int
}

// newLexeme returns a new lexeme
func newLexeme(rLexeme rune, pos int, t lexemeType) *lexeme {
	lexeme := new(lexeme)
	lexeme.pos = pos
	lexeme.text = rLexeme
	lexeme.t = t
	return lexeme
}

// newLexemeStream returns a new lexemeStream
func newLexemeStream(strLexemes string) *lexemeStream {
	ls := new(lexemeStream)
	isInQuotes := false
	var lexType lexemeType
	var err error

	for pos, rLexeme := range strLexemes {
		lexType, err = getLexemeType(rLexeme, pos)
		if lexType == lexemeTypeQuote {
			isInQuotes = !isInQuotes
		}

		// If an unknown character
		// appears in quotes, we
		// can ignore it
		if err != nil {
			if isInQuotes {
				ls.push(newLexeme(rLexeme, pos, lexemeTypeCharacter))
			} else {
				log.Fatal(err.Error())
			}
		} else {
			ls.push(newLexeme(rLexeme, pos, lexType))
		}
	}

	return ls
}

// getLexemeType returns a lexeme type and an error
// if the lexeme did not match any type
func getLexemeType(rLexeme rune, pos int) (lexemeType, error) {
	for lType, identifiers := range lexemeIdentifiers {
		for _, identifier := range identifiers {
			if identifier == rLexeme {
				return lType, nil
			}
		}
	}

	return lexemeTypeUnidentified, &lexemeIdentificationError{pos, rLexeme}
}

// get returns the lexeme at the specified index
func (ls *lexemeStream) get(index int) *lexeme {
	return ls.lexemes[index]
}

func (ls *lexemeStream) pop() *lexeme {
	lex := ls.lexemes[ls.size-1]
	ls.lexemes = ls.lexemes[:ls.size]
	ls.size--
	return lex
}

func (ls *lexemeStream) top() *lexeme {
	if ls.size < 1 {
		return new(lexeme)
	}

	return ls.lexemes[ls.size-1]
}

func (ls *lexemeStream) push(lex *lexeme) {
	ls.lexemes = append(ls.lexemes, lex)
	ls.size++
}

func (ls *lexemeStream) clear() {
	ls.lexemes = nil
	ls.size = 0
}

func (lex *lexeme) string() string {
	return fmt.Sprintf("%v", string(lex.text))
}

func (ls *lexemeStream) string() string {
	str := ""

	for _, lex := range ls.lexemes {
		str += lex.string()
	}

	return str
}
