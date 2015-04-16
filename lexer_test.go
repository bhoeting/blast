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
}

func TestLexerItemType(t *testing.T) {
	lexer := NewLexer("200 200.89 -.99 && xx || \"derp\" (),")
	lexer.Lex()

	assertItemType(t, itemTypeNum, lexer.NextItem())
	assertItemType(t, itemTypeNum, lexer.NextItem())
	assertItemType(t, itemTypeNum, lexer.NextItem())
	assertItemType(t, itemTypeOperator, lexer.NextItem())
	assertItemType(t, itemTypeIdentifier, lexer.NextItem())
	assertItemType(t, itemTypeOperator, lexer.NextItem())
	assertItemType(t, itemTypeString, lexer.NextItem())
	assertItemType(t, itemTypeOpenParen, lexer.NextItem())
	assertItemType(t, itemTypeCloseParen, lexer.NextItem())
	assertItemType(t, itemTypeComma, lexer.NextItem())
}

func TestReservedItems(t *testing.T) {
	lexer := NewLexer("return if else end function XX")
	lexer.Lex()

	assertItemType(t, itemTypeReturn, lexer.NextItem())
	assertItemType(t, itemTypeIf, lexer.NextItem())
	assertItemType(t, itemTypeElse, lexer.NextItem())
	assertItemType(t, itemTypeEnd, lexer.NextItem())
	assertItemType(t, itemTypeFunction, lexer.NextItem())
	assertItemType(t, itemTypeIdentifier, lexer.NextItem())
}

func assertItemType(t *testing.T, expected itemType, actualItem *Item) {
	actual := actualItem.typ

	if actual != expected {
		t.Errorf("Item types not equal\n\t expected: \t%v\n\t   actual:\t%v <%v>",
			expected, actual, actualItem)
	}
}
