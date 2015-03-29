package blast

var precedenceMap = map[opType]int{
	opTypeAssignment:     0,
	opTypeAddition:       1,
	opTypeDivision:       2,
	opTypeSubtraction:    1,
	opTypeMultiplication: 2,
}

// evaluate evaluates a token stream into one token
func (input *tokenStream) evaluate() *token {
	s := newStack()

	for _, t := range input.tokens {
		if isValue(t) {
			s.push(t)
			continue
		}

		if isOp(t) {
			t1, t2 := s.pop(), s.pop()
			s.push(evaluateTokens(t2, t1, t))
		}
	}

	return evaluateToken(s.pop())
}

// toRPN converts a group of tokens
// in infix notation to reverse polish notation. This is
// an implementation of the "Shunting yard" algorithm
// (wikipedia.org/wiki/Shunting-yard_algorithm)
func (input *tokenStream) toRPN() *tokenStream {
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
	return token.t == tokenTypeOperator
}

// isParen determines if the token is a paren
func isParen(token *token) bool {
	return token.t == tokenTypeParen
}

// isOpenParen determines if the token is an open paren
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
