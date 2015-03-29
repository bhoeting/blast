package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStream(t *testing.T) {
	ls := newLexemeStream("20+4*12.99+\"string\"")
	ts := newTokenStream(ls)

	assert.Equal(t, tokenTypeInt, ts.get(0).t)
	assert.Equal(t, tokenTypeOperator, ts.get(1).t)
	assert.Equal(t, tokenTypeInt, ts.get(2).t)
	assert.Equal(t, tokenTypeOperator, ts.get(3).t)
	assert.Equal(t, tokenTypeFloat, ts.get(4).t)
	assert.Equal(t, tokenTypeOperator, ts.get(5).t)
	assert.Equal(t, tokenTypeString, ts.get(6).t)

	ls = newLexemeStream("\"derp\"==5")
	ts = newTokenStream(ls)

	assert.Equal(t, tokenTypeString, ts.get(0).t)
	assert.Equal(t, tokenTypeOperator, ts.get(1).t)
	assert.Equal(t, tokenTypeInt, ts.get(2).t)

	ls = newLexemeStream("if x <= 22")
	ts = newTokenStream(ls)

	assert.Equal(t, tokenTypeIf, ts.get(0).t)
	assert.Equal(t, tokenTypeVar, ts.get(1).t)
	assert.Equal(t, tokenTypeOperator, ts.get(2).t)
	assert.Equal(t, tokenTypeInt, ts.get(3).t)

	ls = newLexemeStream("50 + (4 * 30.6)")
	ts = newTokenStream(ls)

	t.Log(ts.string())

	assert.Equal(t, tokenTypeInt, ts.get(0).t)
	assert.Equal(t, tokenTypeOperator, ts.get(1).t)
	assert.Equal(t, tokenTypeParen, ts.get(2).t)
	assert.Equal(t, tokenTypeInt, ts.get(3).t)
	assert.Equal(t, tokenTypeOperator, ts.get(4).t)
	assert.Equal(t, tokenTypeFloat, ts.get(5).t)
	assert.Equal(t, tokenTypeParen, ts.get(6).t)
}
