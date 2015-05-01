package blast

// NodeStream is a Node slice
// wrapper with stack and
// queue funcitonality
type TokenStream struct {
	pos    int
	size   int
	tokens []Token
}

// Push adds a Node to the NodeStream
func (ts *TokenStream) Push(t Token) Token {
	ts.tokens = append(ts.tokens, t)
	ts.size++
	return t
}

// Pop removes and returns the Node
// at the top of the NodeStream
func (ts *TokenStream) Pop() Token {
	ts.size--
	t := ts.tokens[ts.size]
	ts.tokens = ts.tokens[:ts.size]
	return t
}

// Top returns the Node at the
// top of the NodeStream
func (ts *TokenStream) Top() Token {
	if ts.size-1 < 0 {
		return &tokenNil{}
	}

	t := ts.tokens[ts.size-1]
	return t
}

// RemoveLast removes the first Node
// from a NodeStream
func (ts *TokenStream) RemoveLast() Token {
	ts.size--
	t := ts.tokens[0]
	ts.tokens = ts.tokens[1:]
	return t
}

// Length returns the length of a
// NodeStream
func (ts *TokenStream) Length() int {
	return ts.size
}

// HasNext determines if there is
// another node to return from
// the NodeStream
func (ts *TokenStream) HasNext() bool {
	return ts.pos < ts.size
}

// Next returns the next node
// from the NodeStream
func (ts *TokenStream) Next() Token {
	if ts.pos >= ts.size {
		return &tokenNil{}
	}

	token := ts.tokens[ts.pos]
	ts.pos++
	return token
}

// Backup decrements the position
// in the NodeStream
func (ts *TokenStream) Backup() *TokenStream {
	ts.pos--
	return ts
}

// Peek returns the next Node without
// incrementing the position
func (ts *TokenStream) Peek() Token {
	t := ts.Next()
	ts.Backup()
	return t
}

// Evaluate converts the NodeStream
// to RPN notation and evaluates it
func (ts *TokenStream) Evaluate() Token {
	rpn := NewTokenStreamInRPN(ts)
	return EvaluateRPN(rpn)
}

// String returns a string representation
// of the NodeStream
func (ts *TokenStream) String() string {
	str := ""

	for _, token := range ts.tokens {
		str += token.String() + " "
	}

	return str
}

// NewNodeStream returns a new NodeStream
func NewTokenStream() *TokenStream {
	ts := new(TokenStream)
	return ts
}

// Chop returns a new NodeStream from
// containing all the Nodes starting at the
// position of the current NodeStream
func (ts *TokenStream) Chop() *TokenStream {
	newTS := NewTokenStream()

	for ts.HasNext() {
		newTS.Push(ts.Next())
	}

	return newTS
}

// Reverse reverses the order
// of the Nodes in the NodeStream
func (ts *TokenStream) Reverse() {
	reversed := NewTokenStream()

	for len(ts.tokens) != 0 {
		token := ts.Pop()
		reversed.Push(token)
	}

	ts.tokens = reversed.tokens
	ts.size = reversed.size
}

// Reset sets the position to
// zero so Nodes can be returned
// from it again with Next()
func (ts *TokenStream) Reset() {
	ts.pos = 0
}

// NewNodeStreamFromLexer returns a new NodeStream
// from a lexer that already has a slice of Tokens
func NewTokenStreamFromLexer(l *Lexer) *TokenStream {
	ts := NewTokenStream()

	for l.HasNextItem() {
		switch item := l.NextItem(); item.typ {
		case itemTypeNum:
			ts.Push(NewNumber(item.text))
		case itemTypeBool:
			ts.Push(NewBoolean(item.text))
		case itemTypeString:
			ts.Push(NewString(item.text))
		case itemTypeOperator:
			ts.Push(NewOperator(item.text))
		case itemTypeOpenParen, itemTypeCloseParen:
			ts.Push(NewParen(item.text))
		case itemTypeComma:
			ts.Push(NewComma())
		case itemTypeIdentifier:
			if l.HasNextItem() && l.PeekItem().typ == itemTypeOpenParen {
				ts.Push(NewFunctionCall(item.text))
			} else {
				ts.Push(NewVariable(item.text))
			}
		default:
			ts.Push(NewReserved())
		}
	}

	l.itemPos = 0
	return ts
}
