package blast

import "log"

// opPrecedenceMap is used to determine
// the precedence of an operator
var opPrecedenceMap = map[opType]int{
	opTypeAddition:             2,
	opTypeSubtraction:          2,
	opTypeMultiplication:       3,
	opTypeDivision:             3,
	opTypeModulus:              3,
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

// ForDeclaration is a struct
// that holds all the parts
// of a for loop
type ForDeclaration struct {
	start   float64
	end     float64
	step    float64
	counter *Variable
}

// OneLineIf stores the components
// for a one line if test
type OneLineIf struct {
	cond      *NodeStream
	passBlock *NodeStream
	elseBlock *NodeStream
}

// EvaluateRPN evaluates an RPN expression
func EvaluateRPN(ts *NodeStream) Node {
	var node Node
	nodes := NewNodeStream()
	for ts.HasNext() {
		node = ts.Next()
		switch node.GetType() {
		// Push the value Nodes
		case nodeTypeNumber, nodeTypeBoolean,
			nodeTypeVariable, nodeTypeString:
			nodes.Push(node)
		// If an operator is detected, pop two Nodes
		// off the stack and evaluate them
		case nodeTypeOperator:
			t1, t2 := nodes.Pop(), nodes.Pop()
			nodes.Push(EvaluateNodes(t2, t1, node))
		// If a function call Node is detected, then
		// pop nodes off the stack to pass to the
		// function call until argCount is zero
		case nodeTypeFuncCall:
			argCount := ts.Next().(ArgCount)
			args := NewNodeStream()

			for argCount > 0 {
				args.Push(EvaluateNode(nodes.Pop()))
				argCount--
			}

			args.Reverse()
			t := EvalulateFunctionCall(node, args)
			nodes.Push(t)
		}
	}

	return EvaluateNode(nodes.Pop())
}

// NewNodeStreamInRPN takes a nodeStream and rearranges
// the nodes so they are in reverse polish notation
func NewNodeStreamInRPN(ts *NodeStream) *NodeStream {
	funcArgCounts := make(map[int]int, 0)
	currFuncID := -1
	ops, output := NewNodeStream(), NewNodeStream()

	for ts.HasNext() {
		node := ts.Next()
		switch node.GetType() {
		case nodeTypeNumber, nodeTypeBoolean,
			nodeTypeVariable, nodeTypeString:
			output.Push(node)
		case nodeTypeFuncCall:
			currFuncID++
			funcArgCounts[currFuncID] = 1
			ops.Push(node)
		case nodeTypeComma:
			funcArgCounts[currFuncID]++
			for top := ops.Top(); !isLeftParen(top); top = ops.Top() {
				output.Push(ops.Pop())
			}
		case nodeTypeOperator:
			for top := ops.Top(); top.GetType() == nodeTypeOperator; top = ops.Top() {
				if shouldPopOperator(top, node) {
					output.Push(ops.Pop())
				} else {
					break
				}
			}

			ops.Push(node)
		}

		switch pType := getParenType(node); pType {
		case parenTypeOpen:
			ops.Push(node)
			if getParenType(ts.Peek()) == parenTypeClose {
				funcArgCounts[currFuncID] = 0
			}
		case parenTypeClose:
			for top := ops.Top(); !isLeftParen(top); top = ops.Top() {
				output.Push(ops.Pop())
			}

			ops.Pop()
			if ops.Top().GetType() == nodeTypeFuncCall {
				output.Push(ops.Pop())
				output.Push(NewArgCount(funcArgCounts[currFuncID]))
				currFuncID--
			}
		}
	}

	for ops.Length() > 0 {
		output.Push(ops.Pop())
	}

	return output
}

// EvaluateNodes performs an operation of two Nodes
func EvaluateNodes(t1 Node, t2 Node, tokOp Node) Node {
	opType := tokOp.(*Operator).typ

	if opType != opTypeAssignment {
		t1 = EvaluateNode(t1)
	}

	t2 = EvaluateNode(t2)

	switch opType {
	case opTypeAddition:
		return AddNodes(t1, t2)
	case opTypeSubtraction:
		return SubtractNodes(t1, t2)
	case opTypeMultiplication:
		return MultiplyNodes(t1, t2)
	case opTypeDivision:
		return DivideNodes(t1, t2)
	case opTypeExponent:
		return RaiseNodes(t1, t2)
	case opTypeModulus:
		return ModNodes(t1, t2)
	case opTypeAssignment:
		return AssignNode(t1, t2)
	case opTypeGreaterThan,
		opTypeLessThan,
		opTypeLessThanOrEqualTo,
		opTypeGreaterThanOrEqualTo,
		opTypeNotEqualTo,
		opTypeEqualTo,
		opTypeAnd,
		opTypeOr:
		return CompareNodes(t1, t2, tokOp)
	}

	log.Fatalf("Could not %v on %v and %v", tokOp, t1, t2)
	return &nodeNil{}
}

// EvaluateNode returns the value of a variable Node
// or the Node if it's not a variable
func EvaluateNode(t1 Node) Node {
	if t1.GetType() == nodeTypeVariable {
		var v Node
		var err error
		if v, err = GetVar(t1.(*Variable).name); err == nil {
			return v
		}
		log.Fatal(err.Error())
	}
	return t1
}

// EvalulateFunctionCall runs the function stored in a function node
// and returns the result
func EvalulateFunctionCall(funcCall Node, args *NodeStream) Node {
	f, err := GetFunc(funcCall.(*FunctionCall).name)

	if err != nil {
		log.Fatal(err.Error())
	}

	return f.Call(args)
}

// ParseOneLineIf parses a NodeStream into a `OneLineIf` struct
func ParseOneLineIf(ns *NodeStream) *OneLineIf {
	// if x == 1 then print(x) else print(x-1)
	// if x == 2 then print(x+2)
	oli := new(OneLineIf)

	// Skip and remove the `if`
	ns.Next()
	ns = ns.Chop()

	// Extract the condition
	reserved := ns.Next()
	for ns.HasNext() && reserved.(*Reserved).value != "then" {
		reserved = ns.Next()
	}

	ns.Next()
	return oli
}

// ParseForDeclaration parses a NodeStream into a `ForDeclaration`
func ParseForDeclaration(ts *NodeStream) *ForDeclaration {
	// for 1 -> 20, counter, 2
	// for 1 -> 20, counter
	fd := new(ForDeclaration)

	// Skip the "for"
	ts.Next()

	fd.start = Float64FromNode(ts.Next())
	arrowOp := ts.Next()

	if arrowOp.GetType() != nodeTypeOperator ||
		arrowOp.(*Operator).typ != opTypeArrow {
		log.Fatal("Expected -> in for loop declaration")
	}

	fd.end = Float64FromNode(ts.Next())

	// Skip comma
	ts.Next()

	// Get counter
	if ts.HasNext() {
		next := ts.Next()
		if next.GetType() == nodeTypeNumber {
			fd.step = Float64FromNode(next)
		}

		if next.GetType() == nodeTypeVariable {
			fd.counter = next.(*Variable)
		}
	}

	// Skip comma
	ts.Next()

	if ts.HasNext() {
		next := ts.Next()
		if next.GetType() == nodeTypeVariable {
			if scopeIsInitalized {
				v, err := GetVar(next.(*Variable).name)
				if err != nil {
					fd.step = v.(*Number).value
				}
			}
		}

		if next.GetType() == nodeTypeNumber {
			fd.step = next.(*Number).value
		}
	}

	return fd
}

// isLeftParen determines the node
// is a left paren
func isLeftParen(node Node) bool {
	if paren, ok := node.(*Paren); ok {
		return paren.typ == parenTypeOpen
	}

	return false
}

// getParenType gets the paren type of a node
func getParenType(node Node) parenType {
	if paren, ok := node.(*Paren); ok {
		return paren.typ
	}

	return parenTypeNil
}

// shouldPopOperator is used in the conversion to RPN.
// It is used when an operator is read and determines
// if it should be popped based on the operator at the
// top of the stack
func shouldPopOperator(topNode Node, opNode Node) bool {
	var ok bool
	var topOp, op *Operator

	if topOp, ok = topNode.(*Operator); !ok {
		log.Fatal("topNode is not an operator")
	}

	if op, ok = opNode.(*Operator); !ok {
		log.Fatal("topNode is not an operator")
	}

	if op.typ == opTypeExponent {
		return opPrecedenceMap[op.typ] < opPrecedenceMap[topOp.typ]
	}

	return opPrecedenceMap[op.typ] <= opPrecedenceMap[topOp.typ]
}
