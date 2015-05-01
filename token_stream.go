package blast

// NodeStream is a Node slice
// wrapper with stack and
// queue funcitonality
type NodeStream struct {
	pos   int
	size  int
	nodes []Node
}

// Push adds a Node to the NodeStream
func (ns *NodeStream) Push(t Node) Node {
	ns.nodes = append(ns.nodes, t)
	ns.size++
	return t
}

// Pop removes and returns the Node
// at the top of the NodeStream
func (ns *NodeStream) Pop() Node {
	ns.size--
	n := ns.nodes[ns.size]
	ns.nodes = ns.nodes[:ns.size]
	return n
}

// Top returns the Node at the
// top of the NodeStream
func (ns *NodeStream) Top() Node {
	if ns.size-1 < 0 {
		return &nodeNil{}
	}

	n := ns.nodes[ns.size-1]
	return n
}

// RemoveLast removes the first Node
// from a NodeStream
func (ns *NodeStream) RemoveLast() Node {
	ns.size--
	n := ns.nodes[0]
	ns.nodes = ns.nodes[1:]
	return n
}

// Length returns the length of a
// NodeStream
func (ns *NodeStream) Length() int {
	return ns.size
}

// HasNext determines if there is
// another node to return from
// the NodeStream
func (ns *NodeStream) HasNext() bool {
	return ns.pos < ns.size
}

// Next returns the next node
// from the NodeStream
func (ns *NodeStream) Next() Node {
	if ns.pos >= ns.size {
		return &nodeNil{}
	}

	node := ns.nodes[ns.pos]
	ns.pos++
	return node
}

// Backup decremenns the position
// in the NodeStream
func (ns *NodeStream) Backup() *NodeStream {
	ns.pos--
	return ns
}

// Peek returns the next Node without
// incrementing the position
func (ns *NodeStream) Peek() Node {
	n := ns.Next()
	ns.Backup()
	return n
}

// Evaluate converns the NodeStream
// to RPN notation and evaluates it
func (ns *NodeStream) Evaluate() Node {
	rpn := NewNodeStreamInRPN(ns)
	return EvaluateRPN(rpn)
}

// String returns a string representation
// of the NodeStream
func (ns *NodeStream) String() string {
	str := ""

	for _, node := range ns.nodes {
		str += node.String() + " "
	}

	return str
}

// NewNodeStream returns a new NodeStream
func NewNodeStream() *NodeStream {
	ns := new(NodeStream)
	return ns
}

// Chop returns a new NodeStream from
// containing all the Nodes starting at the
// position of the current NodeStream
func (ns *NodeStream) Chop() *NodeStream {
	newns := NewNodeStream()

	for ns.HasNext() {
		newns.Push(ns.Next())
	}

	return newns
}

// Reverse reverses the order
// of the Nodes in the NodeStream
func (ns *NodeStream) Reverse() {
	reversed := NewNodeStream()

	for len(ns.nodes) != 0 {
		node := ns.Pop()
		reversed.Push(node)
	}

	ns.nodes = reversed.nodes
	ns.size = reversed.size
}

// Reset sens the position to
// zero so Nodes can be returned
// from it again with Next()
func (ns *NodeStream) Reset() {
	ns.pos = 0
}

// NewNodeStreamFromLexer returns a new NodeStream
// from a lexer that already has a slice of Tokens
func NewNodeStreamFromLexer(l *Lexer) *NodeStream {
	ns := NewNodeStream()

	for l.HasNextItem() {
		switch item := l.NextItem(); item.typ {
		case tokenTypeNum:
			ns.Push(NewNumber(item.text))
		case tokenTypeBool:
			ns.Push(NewBoolean(item.text))
		case tokenTypeString:
			ns.Push(NewString(item.text))
		case tokenTypeOperator:
			ns.Push(NewOperator(item.text))
		case tokenTypeOpenParen, tokenTypeCloseParen:
			ns.Push(NewParen(item.text))
		case tokenTypeComma:
			ns.Push(NewComma())
		case tokenTypeIdentifier:
			if l.HasNextItem() && l.PeekItem().typ == tokenTypeOpenParen {
				ns.Push(NewFunctionCall(item.text))
			} else {
				ns.Push(NewVariable(item.text))
			}
		default:
			ns.Push(NewReserved())
		}
	}

	l.tokenPos = 0
	return ns
}
