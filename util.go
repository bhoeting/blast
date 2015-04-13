package blast

import "errors"

var (
	errQueueIsEmpty = errors.New("Queue is empty")
	errStackIsEmpty = errors.New("Stack is empty")
	tokenNull       = &token{nil, 0, 0, tokenTypeNull}
)

// stack is a basic
// stack implementation
type stack struct {
	items map[int]*token
	count int
}

//  is a basic
// stack implementation
type queue struct {
	items []*token
	count int
}

// newStack returns a new stack
func newStack() *stack {
	s := new(stack)
	s.items = make(map[int]*token, 0)
	s.count = 0
	return s
}

// push adds an item to the stack and
// returns the new count
func (s *stack) push(item *token) int {
	s.items[s.count] = item
	s.count++

	return s.count
}

// pop removes and returns the last
// item from the stack
func (s *stack) pop() *token {
	if s.count > 0 {
		i := s.items[s.count-1]
		delete(s.items, s.count)
		s.count--
		return i
	}

	return tokenNull
}

// top returns the next item
// that will be popped without
// removing the item
func (s *stack) top() *token {
	if s.count > 0 {
		return s.items[s.count-1]
	}

	return tokenNull
}

// bottom returns the first
// item in the stack
func (s *stack) bottom() *token {
	if s.count > 0 {
		return s.items[0]
	}

	return tokenNull
}

// toTokenStream returns a tokenStream
// from the stack
func (s *stack) toTokenStream() *tokenStream {
	ts := new(tokenStream)

	for i := 0; i < s.count; i++ {
		ts.add(s.items[i])
	}

	return ts
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
func (q *queue) enqueue(item *token) int {
	q.items = append(q.items, item)
	q.count++
	return q.count
}

// dequeue removes and returns the first
// item in the queue
func (q *queue) dequeue() *token {
	if q.count > 0 {
		i := q.items[0]
		q.items = q.items[1:]
		q.count--
		return i
	}

	return tokenNull
}

func (q *queue) each(f func(item *token)) {
	for _, tok := range q.items {
		f(tok)
	}
}

// size returns the queue's stize
func (q *queue) size() int {
	return q.count
}
