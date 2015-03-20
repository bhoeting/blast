package blast

var precedenceMap = map[opType]int{
	opTypeAddition:       1,
	opTypeDivision:       2,
	opTypeSubtraction:    1,
	opTypeMultiplication: 2,
}

type elType int

const (
	elTypeConst = iota
	elTypeExp
)

// element is a struct that can either
// be an "expression", in which case it
// will store two more elements and an
// operation. It can also be a constant,
// in which case it stores one token
type element struct {
	e1 *element
	e2 *element
	op *token
	c  *token
	t  elType
}

// newElementTree creates an abstract syntax tree
// of elements
func newElementTree(ts *tokenStream) *element {
	return buildElementTree(ts.combine().toReversePolishNotation().tokens)
}

// buildElementTree recursively builds an abstract
// syntax tree of elements. The token slice argument
// must be in reverse polish notation
func buildElementTree(tokens []*token) *element {
	tree := new(element)

	if len(tokens) == 1 {
		tree.t = elTypeConst
		tree.c = tokens[0]
		return tree
	}

	for i := len(tokens) - 1; i >= 0; i-- {
		if isOp(tokens[i]) {
			tree.op = tokens[i]
			tree.t = elTypeExp
			tree.e1 = buildElementTree(tokens[0:1])
			tree.e2 = buildElementTree(tokens[1:i])
			return tree
		}
	}

	return tree
}

// string returns a string representation
// of an element tree
func (e *element) string() string {
	str := buildElementTreeString("", e)
	return str
}

// buildElementTreeString recursively creates a string
// representing an element tree
func buildElementTreeString(str string, e *element) string {
	if e.t == elTypeConst {
		str += e.c.string()
	} else {
		str = buildElementTreeString(str, e.e1)
		str += e.op.string()
		str = buildElementTreeString(str, e.e2)
	}

	return str
}

// toReversePolishNotation converts a group of components
// in infix notation to reverse polish notation. This is
// an implementation of the "Shunting yard" algorithm
// (wikipedia.org/wiki/Shunting-yard_algorithm)
func (input *tokenStream) toReversePolishNotation() *tokenStream {
	ops := newStack()
	output := newQueue()

	for _, token := range input.tokens {
		if isValue(token) {
			output.enqueue(token)
		}

		if isOp(token) {
			for tos := ops.top(); isOp(tos); tos = ops.top() {
				if p := compPrecedence(token, tos); p < 0 {
					output.enqueue(ops.pop())
				} else {
					break
				}
			}

			ops.push(token)
		}

		if isParen(token) {
			if token.parenType() == parenTypeOpen {
				ops.push(token)
			} else {
				for tos := ops.top(); !isOpenParen(tos); tos = ops.top() {
					output.enqueue(ops.pop())
				}

				if isOpenParen(token) {
					ops.pop()
				}
			}
		}
	}

	for ops.top() != tokenNull {
		if isOp(ops.top()) {
			output.enqueue(ops.pop())
		} else {
			ops.pop()
		}
	}

	input.tokens = output.items
	return input
}

// isValue determines if the token
// is a value, which can be a literal
// value or a variable that will
// be converted into a value
func isValue(token *token) bool {
	return token.t == tokenTypeString ||
		token.t == tokenTypeFloat ||
		token.t == tokenTypeInt ||
		token.t == tokenTypeVar
}

// isOp determines if the token
// is an operator
func isOp(token *token) bool {
	return token.t == tokenTypeOp
}

// isParen determines if the token is a paren
func isParen(token *token) bool {
	return token.t == tokenTypeParen
}

func isOpenParen(token *token) bool {
	return isParen(token) && token.parenType() == parenTypeOpen
}

// compPrecedence compares the precedence of two
// operator tokens. If the first token has a higher
// precedence, a positive number is returned. If they're
// equal, 0 is returned. If the second token has a higher
// precedence, a negative number is returned.
func compPrecedence(token1 *token, token2 *token) int {
	return precedenceMap[token1.opType()] - precedenceMap[token2.opType()]
}
