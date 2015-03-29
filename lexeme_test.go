package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexeme(t *testing.T) {
	lets := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	nums := "1234567890"
	parens := "()"
	ops := "+-*/<>="
	e := "\n"
	d := "."
	s := " "
	q := "\""
	c := ","

	for i, r := range nums {
		assert.Equal(t, newLexeme(r, i).t, lexemeTypeNum)
	}

	for i, r := range lets {
		assert.Equal(t, newLexeme(r, i).t, lexemeTypeLetter)
	}

	for i, r := range parens {
		assert.Equal(t, newLexeme(r, i).t, lexemeTypeParen)
	}

	for i, r := range ops {
		assert.Equal(t, newLexeme(r, i).t, lexemeTypeOperator)
	}

	assert.Equal(t, newLexeme(rune(d[0]), 0).t, lexemeTypeDecimalDigit)
	assert.Equal(t, newLexeme(rune(s[0]), 0).t, lexemeTypeSpace)
	assert.Equal(t, newLexeme(rune(q[0]), 0).t, lexemeTypeQuote)
	assert.Equal(t, newLexeme(rune(c[0]), 0).t, lexemeTypeComma)
	assert.Equal(t, newLexeme(rune(e[0]), 0).t, lexemeTypeEOL)
}

func TestLexemeString(t *testing.T) {
	ls := newLexemeStream("200+AA+\"derp\"")

	// 200
	assert.Equal(t, lexemeTypeNum, ls.get(0).t)
	assert.Equal(t, lexemeTypeNum, ls.get(1).t)
	assert.Equal(t, lexemeTypeNum, ls.get(2).t)

	// +
	assert.Equal(t, lexemeTypeOperator, ls.get(3).t)

	// AA
	assert.Equal(t, lexemeTypeLetter, ls.get(4).t)
	assert.Equal(t, lexemeTypeLetter, ls.get(5).t)

	// +
	assert.Equal(t, lexemeTypeOperator, ls.get(6).t)
	assert.Equal(t, lexemeTypeQuote, ls.get(7).t)

	// "derp"
	assert.Equal(t, lexemeTypeLetter, ls.get(8).t)
	assert.Equal(t, lexemeTypeLetter, ls.get(9).t)
	assert.Equal(t, lexemeTypeLetter, ls.get(10).t)
	assert.Equal(t, lexemeTypeLetter, ls.get(11).t)
	assert.Equal(t, lexemeTypeQuote, ls.get(12).t)
}
