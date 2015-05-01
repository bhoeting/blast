package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPN(t *testing.T) {
	ts := NewNodeStreamFromLexer(Lex("212 + 341"))
	rpn := NewNodeStreamInRPN(ts)

	assert.Equal(t, nodeTypeNumber, rpn.Next().GetType())
	assert.Equal(t, nodeTypeNumber, rpn.Next().GetType())
	assert.Equal(t, nodeTypeOperator, rpn.Next().GetType())

	ts = NewNodeStreamFromLexer(Lex("40 + (3 * 30.6)"))
	rpn = NewNodeStreamInRPN(ts)
	assert.Equal(t, "40 3 30.6 * + ", rpn.String())

	ts = NewNodeStreamFromLexer(Lex("27 * pi() + max(22, 33)"))
	rpn = NewNodeStreamInRPN(ts)
	assert.Equal(t, "27 pi() 0 * 22 33 max() 2 + ", rpn.String())

	ts = NewNodeStreamFromLexer(Lex("40 + max(4, 4+1, 88, 88) + 990"))
	rpn = NewNodeStreamInRPN(ts)
	assert.Equal(t, "40 4 4 1 + 88 88 max() 4 + 990 + ", rpn.String())

	ts = NewNodeStreamFromLexer(Lex("max(min(33, 413), 300, 102)"))
	rpn = NewNodeStreamInRPN(ts)
	assert.Equal(t, "33 413 min() 2 300 102 max() 3 ", rpn.String())
}

func TestRPNEvaluation(t *testing.T) {
	ts := NewNodeStreamFromLexer(Lex("212 + 341"))
	rpn := NewNodeStreamInRPN(ts)

	result := EvaluateRPN(rpn)

	assert.Equal(t, 553, Float64FromNode(result))
}

func TestForLoopParsing(t *testing.T) {
	fd := ParseForDeclaration(NewNodeStreamFromLexer(Lex("for 1 -> 20, x, 1")))
	assert.Equal(t, 1.0, fd.start)
	assert.Equal(t, 20, fd.end)
	assert.Equal(t, "x", fd.counter.name)
	assert.Equal(t, 1.0, fd.step)
}
