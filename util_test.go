package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the stack struct
func TestStack(t *testing.T) {
	s := newStack()
	ts := newTokenStream("2+70.2*29").combine()

	for _, tok := range ts.tokens {
		s.push(tok)
	}

	assert.Equal(t, 29, s.pop().integer())
	assert.Equal(t, opTypeMultiplication, s.pop().opType())
	assert.Equal(t, 70.2, s.pop().float())
	assert.Equal(t, opTypeAddition, s.pop().opType())
	assert.Equal(t, 2, s.pop().integer())

	assert.Equal(t, tokenNull, s.pop())
}

// Test the queue struct
func TestQueue(t *testing.T) {
	q := newQueue()
	ts := newTokenStream("\"string\" + 200 + 23.8").combine()

	for _, tok := range ts.tokens {
		q.enqueue(tok)
	}

	assert.Equal(t, "string", q.dequeue().string())
	assert.Equal(t, opTypeAddition, q.dequeue().opType())
	assert.Equal(t, 200, q.dequeue().integer())
	assert.Equal(t, opTypeAddition, q.dequeue().opType())
	assert.Equal(t, 23.8, q.dequeue().float())

	assert.Equal(t, tokenNull, q.dequeue())
}
