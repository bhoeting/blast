package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReversePolishNotation(t *testing.T) {
	actual := newTokenStream(newLexemeStream("40 + (3 * 30.6)")).toRPN().string()
	expected := "40`3`30.6`*`+`"
	assert.Equal(t, expected, actual)

	actual = newTokenStream(newLexemeStream("40 + max(4, 4+1)")).toRPN().string()
	expected = "40`4`4`1`+`max()`+`"
	assert.Equal(t, expected, actual)

	actual = newTokenStream(newLexemeStream("abs(199)")).toRPN().string()
	expected = "199`abs()`"
	assert.Equal(t, expected, actual)
}

func TestEvaluationOfTokenStream(t *testing.T) {
	Init()

	actual := newTokenStream(newLexemeStream("40 + 3 * 0.5")).evaluate().data
	expected := 41.5
	assert.Equal(t, expected, actual)

	actual = newTokenStream(newLexemeStream("2 >= 3")).evaluate().data
	assert.Equal(t, false, actual)

	actual = newTokenStream(newLexemeStream("40 + max(4, 4+1)")).evaluate().data
	assert.Equal(t, 45, actual)

	actual = newTokenStream(newLexemeStream("max(199,189+1)")).evaluate().data
	assert.Equal(t, 199, actual)
}
