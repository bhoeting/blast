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

	ts = NewTokenStreamFromLexer(Lex("\"derpsause\" + 300 == 41 && x <= y(220)"))
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
