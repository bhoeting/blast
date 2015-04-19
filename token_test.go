package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStream(t *testing.T) {
	ts := NewTokenStreamFromLexer(Lex("200.98 + 300"))
	assert.Equal(t, tokenTypeNumber, ts.Next().GetType())
	assert.Equal(t, tokenTypeOperator, ts.Next().GetType())
	assert.Equal(t, tokenTypeNumber, ts.Next().GetType())

	ts = NewTokenStreamFromLexer(Lex("true false \"derpsause\" + 300 == 41 && x <= y(220)"))
	assert.Equal(t, tokenTypeBoolean, ts.Next().GetType())
	assert.Equal(t, tokenTypeBoolean, ts.Next().GetType())
	assert.Equal(t, tokenTypeString, ts.Next().GetType())
	assert.Equal(t, tokenTypeOperator, ts.Next().GetType())
	assert.Equal(t, tokenTypeNumber, ts.Next().GetType())
	assert.Equal(t, tokenTypeOperator, ts.Next().GetType())
	assert.Equal(t, tokenTypeNumber, ts.Next().GetType())
	assert.Equal(t, tokenTypeOperator, ts.Next().GetType())
	assert.Equal(t, tokenTypeVariable, ts.Next().GetType())
	assert.Equal(t, tokenTypeOperator, ts.Next().GetType())
	assert.Equal(t, tokenTypeFuncCall, ts.Next().GetType())
}

func TestTokenConversions(t *testing.T) {
	ts := NewTokenStreamFromLexer(Lex("200.9 false true \"derp\""))
	num := ts.Next()
	fls := ts.Next()
	tru := ts.Next()
	str := ts.Next()

	// Numbers
	assert.Equal(t, 200.9, NumberFromToken(num))
	assert.Equal(t, 0.0, NumberFromToken(fls))
	assert.Equal(t, 1.0, NumberFromToken(tru))

	// Bools
	assert.Equal(t, true, BooleanFromToken(tru))
	assert.Equal(t, false, BooleanFromToken(fls))

	// Strings
	assert.Equal(t, "200.9", StringFromToken(num))
	assert.Equal(t, "false", StringFromToken(fls))
	assert.Equal(t, "true", StringFromToken(tru))
	assert.Equal(t, "derp", StringFromToken(str))
}
