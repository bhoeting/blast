package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexerMovement(t *testing.T) {
	lexer := NewLexer("200+\n")

	assert.Equal(t, '2', lexer.Next())
	assert.Equal(t, '0', lexer.Next())
	assert.Equal(t, '0', lexer.Next())

	lexer.Backup()

	assert.Equal(t, '0', lexer.Next())
	assert.Equal(t, '+', lexer.Peek())
	assert.Equal(t, '+', lexer.Next())
	assert.Equal(t, '\n', lexer.Next())
	assert.Equal(t, eof, lexer.Next())
}

func TestLexerItemText(t *testing.T) {
	lexer := NewLexer("201+901")
	lexer.Lex()
	assert.Equal(t, "201", lexer.NextItem().text)
	assert.Equal(t, "+", lexer.NextItem().text)
	assert.Equal(t, "901", lexer.NextItem().text)

	lexer = NewLexer("x2,x1,91,-99.9,-.78 .87")
	lexer.Lex()
	assert.Equal(t, "x2", lexer.NextItem().text)
	assert.Equal(t, ",", lexer.NextItem().text)
	assert.Equal(t, "x1", lexer.NextItem().text)
	assert.Equal(t, ",", lexer.NextItem().text)
	assert.Equal(t, "91", lexer.NextItem().text)
	assert.Equal(t, ",", lexer.NextItem().text)
	assert.Equal(t, "-99.9", lexer.NextItem().text)
	assert.Equal(t, ",", lexer.NextItem().text)
	assert.Equal(t, "-.78", lexer.NextItem().text)
	assert.Equal(t, ".87", lexer.NextItem().text)

	lexer = NewLexer("300 - -30")
	lexer.Lex()
	assert.Equal(t, "300", lexer.NextItem().text)
	assert.Equal(t, "-", lexer.NextItem().text)
	assert.Equal(t, "-30", lexer.NextItem().text)

	lexer = NewLexer("314.9 + \"str()_=182ing\"")
	lexer.Lex()
	assert.Equal(t, "314.9", lexer.NextItem().text)
	assert.Equal(t, "+", lexer.NextItem().text)
	assert.Equal(t, "str()_=182ing", lexer.NextItem().text)

	lexer = NewLexer("if x == 1 + 100")
	lexer.Lex()
	assert.Equal(t, "if", lexer.NextItem().text)
	assert.Equal(t, "x", lexer.NextItem().text)
	assert.Equal(t, "==", lexer.NextItem().text)
	assert.Equal(t, "1", lexer.NextItem().text)
	assert.Equal(t, "+", lexer.NextItem().text)
	assert.Equal(t, "100", lexer.NextItem().text)

	lexer = NewLexer("((100 + 13) * 78)")
	lexer.Lex()
	assert.Equal(t, "(", lexer.NextItem().text)
	assert.Equal(t, "(", lexer.NextItem().text)
	assert.Equal(t, "100", lexer.NextItem().text)
	assert.Equal(t, "+", lexer.NextItem().text)
	assert.Equal(t, "13", lexer.NextItem().text)
	assert.Equal(t, ")", lexer.NextItem().text)
	assert.Equal(t, "*", lexer.NextItem().text)
	assert.Equal(t, "78", lexer.NextItem().text)
	assert.Equal(t, ")", lexer.NextItem().text)

	lexer = NewLexer("true false")
	lexer.Lex()
	assert.Equal(t, "true", lexer.NextItem().text)
	assert.Equal(t, "false", lexer.NextItem().text)

	lexer = NewLexer("for 1 -> 20, i, 2")
	lexer.Lex()
	assert.Equal(t, "for", lexer.NextItem().text)
	assert.Equal(t, "1", lexer.NextItem().text)
	assert.Equal(t, "->", lexer.NextItem().text)
	assert.Equal(t, "20", lexer.NextItem().text)
	assert.Equal(t, ",", lexer.NextItem().text)
	assert.Equal(t, "i", lexer.NextItem().text)
	assert.Equal(t, ",", lexer.NextItem().text)
	assert.Equal(t, "2", lexer.NextItem().text)
}

func TestLexerItemType(t *testing.T) {
	lexer := NewLexer("200 200.89 -.99 && xx || \"derp\" (),")
	lexer.Lex()

	assertItemType(t, tokenTypeNum, lexer.NextItem())
	assertItemType(t, tokenTypeNum, lexer.NextItem())
	assertItemType(t, tokenTypeNum, lexer.NextItem())
	assertItemType(t, tokenTypeOperator, lexer.NextItem())
	assertItemType(t, tokenTypeIdentifier, lexer.NextItem())
	assertItemType(t, tokenTypeOperator, lexer.NextItem())
	assertItemType(t, tokenTypeString, lexer.NextItem())
	assertItemType(t, tokenTypeOpenParen, lexer.NextItem())
	assertItemType(t, tokenTypeCloseParen, lexer.NextItem())
	assertItemType(t, tokenTypeComma, lexer.NextItem())

	lexer = NewLexer("200.9 false true \"derp\"")
	lexer.Lex()
	lexer.NextItem()

	assertItemType(t, tokenTypeBool, lexer.NextItem())
	assertItemType(t, tokenTypeBool, lexer.NextItem())

	lexer = NewLexer("return n1 + n2")
	lexer.Lex()

	assertItemType(t, tokenTypeReturn, lexer.NextItem())
	assertItemType(t, tokenTypeIdentifier, lexer.NextItem())
	assertItemType(t, tokenTypeOperator, lexer.NextItem())
	assertItemType(t, tokenTypeIdentifier, lexer.NextItem())
	assert.Equal(t, false, lexer.HasNextItem())
}

func TestReservedItems(t *testing.T) {
	lexer := NewLexer("for return if else end function XX")
	lexer.Lex()

	assertItemType(t, tokenTypeEnd, lexer.NextItem())
	assertItemType(t, tokenTypeReturn, lexer.NextItem())
	assertItemType(t, tokenTypeIf, lexer.NextItem())
	assertItemType(t, tokenTypeElse, lexer.NextItem())
	assertItemType(t, itemTypeEnd, lexer.NextItem())
	assertItemType(t, tokenTypeFunction, lexer.NextItem())
	assertItemType(t, tokenTypeIdentifier, lexer.NextItem())
}

func assertItemType(t *testing.T, expected itemType, actualItem *Token) {
	actual := actualItem.typ

	if actual != expected {
		t.Errorf("Item types not equal\n\t expected: \t%v\n\t   actual:\t%v <%v>",
			expected, actual, actualItem)
	}
}
