package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReversePolishNotation(t *testing.T) {
	actual := newTokenStream(newLexemeStream("40 + (3 * 30.6)")).toRPN().string()
	expected := "40,3,30.6,*,+,"
	assert.Equal(t, expected, actual)
}

// func TestEvaluationOfTokenStream(t *testing.T) {
// 	actual := newTokenStream(newLexemeStream("40 + 3 * 0.5")).toRPN().string().data
// 	expected := 41.5
// 	assert.Equal(t, expected, actual)
// }
