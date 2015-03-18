package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the stack struct
func TestStack(t *testing.T) {
	s := newStack()
	s.push(3)
	s.push(4)
	s.push(88)

	assert.Equal(t, 88, s.pop())
	assert.Equal(t, 4, s.pop())
	assert.Equal(t, 3, s.pop())
	assert.Equal(t, 0, s.size())

	assert.IsType(t, errStackIsEmpty, s.pop())
}

// Test the queue struct
func TestQueue(t *testing.T) {
	q := newQueue()
	q.enqueue(3)
	q.enqueue(4)
	q.enqueue(88)

	assert.Equal(t, 3, q.dequeue())
	assert.Equal(t, 4, q.dequeue())
	assert.Equal(t, 88, q.dequeue())
	assert.Equal(t, 0, q.size())

	assert.IsType(t, errQueueIsEmpty, q.dequeue())
}
