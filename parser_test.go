package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReversePolishNotation(t *testing.T) {
	actual := newTokenStream("40 + (3 * 30.6)").combine().toRPN().string()
	expected := "40,3,30.6,*,+,"

	assert.Equal(t, expected, actual)
}

func TestEvaluationOfTokenStream(t *testing.T) {
	actual := newTokenStream("40 + 3 * 0.5").combine().toRPN().evaluate().data
	expected := 41.5

	assert.Equal(t, expected, actual)
}
