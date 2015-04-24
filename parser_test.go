package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPN(t *testing.T) {
	ts := NewTokenStreamFromLexer(Lex("212 + 341"))
	rpn := NewTokenStreamInRPN(ts)

	assert.Equal(t, tokenTypeNumber, rpn.Next().GetType())
	assert.Equal(t, tokenTypeNumber, rpn.Next().GetType())
	assert.Equal(t, tokenTypeOperator, rpn.Next().GetType())

	ts = NewTokenStreamFromLexer(Lex("40 + (3 * 30.6)"))
	rpn = NewTokenStreamInRPN(ts)
	assert.Equal(t, "40 3 30.6 * + ", rpn.String())

	ts = NewTokenStreamFromLexer(Lex("27 * pi() + max(22, 33)"))
	rpn = NewTokenStreamInRPN(ts)
	assert.Equal(t, "27 pi() 0 * 22 33 max() 2 + ", rpn.String())

	ts = NewTokenStreamFromLexer(Lex("40 + max(4, 4+1, 88, 88) + 990"))
	rpn = NewTokenStreamInRPN(ts)
	assert.Equal(t, "40 4 4 1 + 88 88 max() 4 + 990 + ", rpn.String())

	ts = NewTokenStreamFromLexer(Lex("max(min(33, 413), 300, 102)"))
	rpn = NewTokenStreamInRPN(ts)
	assert.Equal(t, "33 413 min() 2 300 102 max() 3 ", rpn.String())
}

func TestRPNEvaluation(t *testing.T) {
	ts := NewTokenStreamFromLexer(Lex("212 + 341"))
	rpn := NewTokenStreamInRPN(ts)

	result := EvaluateRPN(rpn)

	assert.Equal(t, 553, NumberFromToken(result))
}

func TestForLoopParsing(t *testing.T) {
	fd := ParseForDeclaration(NewTokenStreamFromLexer(Lex("for 1 -> 20, x, 1")))
	assert.Equal(t, 1.0, fd.start)
	assert.Equal(t, 20, fd.end)
	assert.Equal(t, "x", fd.counter.name)
	assert.Equal(t, 1.0, fd.step)
}
