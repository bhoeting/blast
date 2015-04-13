package blast

var precedenceMap = map[opType]int{
	opTypeAssignment:     0,
	opTypeAddition:       1,
	opTypeDivision:       2,
	opTypeSubtraction:    1,
	opTypeMultiplication: 2,
}

// parse parses the tokens in a tokenStream
// into one token
func (ts *tokenStream) parse() *token {
	newTS := *ts

	switch ts.get(0).t {
	case tokenTypeReturn, tokenTypeIf:
		newTS.tokens = newTS.tokens[1:]
	}

	return ts.toRPN().evaluate()
}

// evaluate evaluates a token stream into one token
func (ts *tokenStream) evaluate() *token {
	s := newStack()

	if ts.isRPN == false {
		ts.toRPN()
	}

	for _, token := range ts.tokens {
		if isValue(token) {
			s.push(token)
			continue
		}

		if isOp(token) {
			t1, t2 := s.pop(), s.pop()
			s.push(evaluateTokens(t2, t1, token))
			continue
		}

		if token.t == tokenTypeFunctionCall {
			result, nParams := callFunction(token.str(), s.toTokenStream().tokens)

			for nParams > 0 {
				s.pop()
				nParams--
			}

			s.push(result)
		}
	}

	return evaluateToken(s.pop())
}

// toRPN converts a group of tokens in infix notation
// to reverse polish notation. This is an implementation
// of the "Shunting yard" algorithm
// (wikipedia.org/wiki/Shunting-yard_algorithm)
func (input *tokenStream) toRPN() *tokenStream {
	input.isRPN = true
	ops := newStack()
	output := newQueue()

	for _, token := range input.tokens {
		if isValue(token) {
			output.enqueue(token)
			continue
		}

		if token.t == tokenTypeFunctionCall {
			ops.push(token)
			continue
		}

		// If the token is an argument seperator
		if token.t == tokenTypeComma {
			// Until the token at the top of the stack is a left parenthesis,
			// pop operators off the stack onto the ouput queue
			for tos := ops.top(); !isOpenParen(tos); tos = ops.top() {
				output.enqueue(ops.pop())
			}

			continue
		}

		// If the token is an operator
		if isOp(token) {
			// While there is an operator at the top of the stack
			for tos := ops.top(); isOp(tos); tos = ops.top() {
				// If the precendence is less than operator
				// at the top of the stack, pop & queue
				if p := compPrecedence(token, tos); p < 0 {
					output.enqueue(ops.pop())
				} else {
					break
				}
			}

			ops.push(token)
			continue
		}

		// If the token is a parenthesis
		if isParen(token) {
			// If right paren, push to stack
			if token.parenType() == parenTypeOpen {
				ops.push(token)
			} else {
				// If left paren, pop operators to the the stack until
				// a left paren is detected
				for tos := ops.top(); !isOpenParen(tos); tos = ops.top() {
					output.enqueue(ops.pop())
				}

				// Pop the left paren
				ops.pop()

				// If the top of the stack contains a function
				// call, pop & push
				if ops.top().t == tokenTypeFunctionCall {
					output.enqueue(ops.pop())
				}
			}

			continue
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
		token.t == tokenTypeVar ||
		token.t == tokenTypeUnkown
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

// isEvaluatable determines if a token
// can be evaulated
func isEvaluatable(token *token) bool {
	return token.t != tokenTypeComma &&
		token.t != tokenTypeFunctionCall &&
		token.t != tokenTypeReturn &&
		token.t != tokenTypeParen &&
		token.t != tokenTypeOperator &&
		token.t != tokenTypeIf &&
		token.t != tokenTypeEnd
}
