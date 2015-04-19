package blast

import "log"

var opPrecedenceMap = map[opType]int{
	opTypeAddition:             2,
	opTypeSubtraction:          2,
	opTypeMultiplication:       3,
	opTypeDivision:             3,
	opTypeExponent:             4,
	opTypeLessThan:             1,
	opTypeLessThanOrEqualTo:    1,
	opTypeEqualTo:              1,
	opTypeNotEqualTo:           1,
	opTypeGreaterThan:          1,
	opTypeGreaterThanOrEqualTo: 1,
	opTypeAnd:                  0,
	opTypeOr:                   0,
	opTypeAssignment:           -1,
}

func EvaluateRPN(ts *TokenStream) Token {
	tokens := NewTokenStream()

	for ts.HasNext() {
		token := ts.Next()
		switch token.GetType() {
		case tokenTypeNumber, tokenTypeBoolean,
			tokenTypeVariable, tokenTypeString:
			tokens.Push(token)
		case tokenTypeOperator:
			t1, t2 := tokens.Pop(), tokens.Pop()
			tokens.Push(EvaluateTokens(t2, t1, token))
		case tokenTypeFuncCall:
			argCount := ts.Next().(ArgCount)
			args := NewTokenStream()
			for argCount > 0 {
				args.Push(tokens.Pop())
				argCount--
			}
			tokens.Push(EvalulateFunctionCall(token, args))
		}
	}

	return EvaluateToken(tokens.Pop())
}

func NewTokenStreamInRPN(ts *TokenStream) *TokenStream {
	funcArgCounts := make(map[int]int, 0)
	currFuncId := -1
	ops, output := NewTokenStream(), NewTokenStream()

	for ts.HasNext() {
		token := ts.Next()
		switch token.GetType() {
		case tokenTypeNumber, tokenTypeBoolean,
			tokenTypeVariable, tokenTypeString:
			output.Push(token)
		case tokenTypeFuncCall:
			currFuncId++
			funcArgCounts[currFuncId] = 1
			ops.Push(token)
		case tokenTypeComma:
			funcArgCounts[currFuncId]++
			for top := ops.Top(); !isLeftParen(top); top = ops.Top() {
				output.Push(ops.Pop())
			}
		case tokenTypeOperator:
			for top := ops.Top(); top.GetType() == tokenTypeOperator; top = ops.Top() {
				if shouldPopOperator(top, token) {
					output.Push(ops.Pop())
				} else {
					break
				}
			}

			ops.Push(token)
		}

		switch pType := getParenType(token); pType {
		case parenTypeOpen:
			ops.Push(token)
			if getParenType(ts.Peek()) == parenTypeClose {
				funcArgCounts[currFuncId] = 0
			}
		case parenTypeClose:
			for top := ops.Top(); !isLeftParen(top); top = ops.Top() {
				output.Push(ops.Pop())
			}

			ops.Pop()
			if ops.Top().GetType() == tokenTypeFuncCall {
				output.Push(ops.Pop())
				output.Push(NewArgCount(funcArgCounts[currFuncId]))
				currFuncId--
			}
		}
	}

	for ops.Length() > 0 {
		output.Push(ops.Pop())
	}

	return output
}

func EvaluateTokens(t1 Token, t2 Token, tokOp Token) Token {
	switch tokOp.(*Operator).typ {
	case opTypeAddition:
		return AddTokens(t1, t2)
	case opTypeSubtraction:
		return SubtractTokens(t1, t2)
	case opTypeMultiplication:
		return MultiplyTokens(t1, t2)
	case opTypeDivision:
		return MultiplyTokens(t1, t2)
	case opTypeExponent:
		return RaiseTokens(t1, t2)
	case opTypeAssignment:
		return AssignTokens(t1, t2)
	case opTypeGreaterThan,
		opTypeLessThan,
		opTypeLessThanOrEqualTo,
		opTypeGreaterThanOrEqualTo,
		opTypeNotEqualTo,
		opTypeEqualTo:
		return CompareTokens(t1, t2, tokOp)
	}

	log.Fatalf("Could not %v on %v and %v", tokOp, t1, t2)
	return &tokenNil{}
}

func EvaluateToken(t1 Token) Token {
	return &tokenNil{}
}

func EvalulateFunctionCall(funcCall Token, args *TokenStream) Token {
	return &tokenNil{}
}

func isLeftParen(token Token) bool {
	if paren, ok := token.(*Paren); ok {
		return paren.typ == parenTypeOpen
	}

	return false
}

func getParenType(token Token) parenType {
	if paren, ok := token.(*Paren); ok {
		return paren.typ
	}

	return parenTypeNil
}

func shouldPopOperator(topToken Token, opToken Token) bool {
	var ok bool
	var topOp, op *Operator

	if topOp, ok = topToken.(*Operator); !ok {
		log.Fatal("topToken is not an operator")
	}

	if op, ok = opToken.(*Operator); !ok {
		log.Fatal("topToken is not an operator")
	}

	if op.typ == opTypeExponent {
		return opPrecedenceMap[op.typ] < opPrecedenceMap[topOp.typ]
	} else {
		return opPrecedenceMap[op.typ] <= opPrecedenceMap[topOp.typ]
	}

	return false
}
