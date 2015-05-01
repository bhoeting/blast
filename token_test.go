package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStream(t *testing.T) {
	ts := NewNodeStreamFromLexer(Lex("200.98 + 300"))
	assert.Equal(t, nodeTypeNumber, ts.Next().GetType())
	assert.Equal(t, nodeTypeOperator, ts.Next().GetType())
	assert.Equal(t, nodeTypeNumber, ts.Next().GetType())

	ts = NewNodeStreamFromLexer(Lex("true false \"derpsause\" + 300 == 41 && x <= y(220)"))
	assert.Equal(t, nodeTypeBoolean, ts.Next().GetType())
	assert.Equal(t, nodeTypeBoolean, ts.Next().GetType())
	assert.Equal(t, nodeTypeString, ts.Next().GetType())
	assert.Equal(t, nodeTypeOperator, ts.Next().GetType())
	assert.Equal(t, nodeTypeNumber, ts.Next().GetType())
	assert.Equal(t, nodeTypeOperator, ts.Next().GetType())
	assert.Equal(t, nodeTypeNumber, ts.Next().GetType())
	assert.Equal(t, nodeTypeOperator, ts.Next().GetType())
	assert.Equal(t, nodeTypeVariable, ts.Next().GetType())
	assert.Equal(t, nodeTypeOperator, ts.Next().GetType())
	assert.Equal(t, nodeTypeFuncCall, ts.Next().GetType())

	ts = NewNodeStreamFromLexer(Lex("return n1 + n2"))
	assert.Equal(t, nodeTypeReserved, ts.Next().GetType())
	assert.Equal(t, nodeTypeVariable, ts.Next().GetType())
	assert.Equal(t, nodeTypeOperator, ts.Next().GetType())
	assert.Equal(t, nodeTypeVariable, ts.Next().GetType())
}

func TestTokenConversions(t *testing.T) {
	ts := NewNodeStreamFromLexer(Lex("200.9 false true \"derp\""))
	num := ts.Next()
	fls := ts.Next()
	tru := ts.Next()
	str := ts.Next()

	// Numbers
	assert.Equal(t, 200.9, Float64FromNode(num))
	assert.Equal(t, 0.0, Float64FromNode(fls))
	assert.Equal(t, 1.0, Float64FromNode(tru))

	// Bools
	assert.Equal(t, true, BooleanFromNode(tru))
	assert.Equal(t, false, BooleanFromNode(fls))

	// Strings
	assert.Equal(t, "200.9", StringFromNode(num))
	assert.Equal(t, "false", StringFromNode(fls))
	assert.Equal(t, "true", StringFromNode(tru))
	assert.Equal(t, "derp", StringFromNode(str))
}
