package blast

import "errors"

var (
	errQueueIsEmpty = errors.New("Queue is empty")
	errStackIsEmpty = errors.New("Stack is empty")
)

// stack is a basic
// stack implementation
type stack struct {
	items map[int]interface{}
	count int
}

//  is a basic
// stack implementation
type queue struct {
	items []interface{}
	count int
}

// newStack returns a new stack
func newStack() *stack {
	s := new(stack)
	s.items = make(map[int]interface{}, 0)
	s.count = 0
	return s
}

// push adds an item to the stack and
// returns the new count
func (s *stack) push(item interface{}) int {
	s.items[s.count] = item
	s.count++

	return s.count
}

// pop removes and returns the last
// item from the stack
func (s *stack) pop() interface{} {
	if s.count > 0 {
		i := s.items[s.count-1]
		delete(s.items, s.count)
		s.count--
		return i
	}

	return errStackIsEmpty
}

// size returns the stack's stize
func (s *stack) size() int {
	return s.count
}

// newQueue returns a new queue
func newQueue() *queue {
	q := new(queue)
	q.count = 0
	return q
}

// enqueue adds a new item
// to the back of the queue
func (q *queue) enqueue(item interface{}) int {
	q.items = append(q.items, item)
	q.count++
	return q.count
}

// dequeue removes and returns the first
// item in the queue
func (q *queue) dequeue() interface{} {
	if q.count > 0 {
		i := q.items[0]
		q.items = q.items[1:]
		q.count--
		return i
	}

	return errQueueIsEmpty
}

// size returns the queue's stize
func (q *queue) size() int {
	return q.count
}
