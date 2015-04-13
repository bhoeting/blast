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

	test := func(lexemeString string, lexemeType lexemeType) {
		for _, lexeme := range newLexemeStream(lexemeString).lexemes {
			assert.Equal(t, lexemeType, lexeme.t)
		}
	}

	test(lets, lexemeTypeCharacter)
	test(nums, lexemeTypeNum)
	test(parens, lexemeTypeParen)
	test(ops, lexemeTypeOperator)
	test(e, lexemeTypeEOL)
	test(d, lexemeTypeDecimalDigit)
	test(s, lexemeTypeSpace)
	test(q, lexemeTypeQuote)
	test(c, lexemeTypeComma)
}

func TestLexemeString(t *testing.T) {
	ls := newLexemeStream("200+AA+\"derp^\"")

	// 200
	assert.Equal(t, lexemeTypeNum, ls.get(0).t)
	assert.Equal(t, lexemeTypeNum, ls.get(1).t)
	assert.Equal(t, lexemeTypeNum, ls.get(2).t)

	// +
	assert.Equal(t, lexemeTypeOperator, ls.get(3).t)

	// AA
	assert.Equal(t, lexemeTypeCharacter, ls.get(4).t)
	assert.Equal(t, lexemeTypeCharacter, ls.get(5).t)

	// +
	assert.Equal(t, lexemeTypeOperator, ls.get(6).t)
	assert.Equal(t, lexemeTypeQuote, ls.get(7).t)

	// "derp^"
	assert.Equal(t, lexemeTypeCharacter, ls.get(8).t)
	assert.Equal(t, lexemeTypeCharacter, ls.get(9).t)
	assert.Equal(t, lexemeTypeCharacter, ls.get(10).t)
	assert.Equal(t, lexemeTypeCharacter, ls.get(11).t)
	assert.Equal(t, lexemeTypeCharacter, ls.get(12).t)
	assert.Equal(t, lexemeTypeQuote, ls.get(13).t)
}

/*func TestUnknownLexeme(t *testing.T) {
	ls := newLexemeStream("\"^%$$!@#$\"")

}*/
