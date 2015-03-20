package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReversePolishNotation(t *testing.T) {
	actual := newTokenStream("40 + (3 * 30.6)").combine().toReversePolishNotation().string()
	expected := "40330.6*+"

	assert.Equal(t, expected, actual)
}

func TestElementTree(t *testing.T) {
	tree := newElementTree(newTokenStream("40+(3*30.6))"))
	tree.string()
	t.Log(tree.string())
}
