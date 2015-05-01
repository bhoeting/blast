package blast

import (
	"log"
	"strings"
)

// AddNodes adds two Nodes into one Node
func AddNodes(n1 Node, n2 Node) Node {
	var result Node

	if n1.GetType() == nodeTypeString || n2.GetType() == nodeTypeString {
		result = NewString(StringFromNode(n1) + StringFromNode(n2))
	} else {
		result = NewNumberFromFloat(Float64FromNode(n1) + Float64FromNode(n2))
	}

	return result
}

// SubtractNodes subtracts two Nodes into one Node
func SubtractNodes(n1 Node, n2 Node) Node {
	if n1.GetType() != nodeTypeString && n2.GetType() != nodeTypeString {
		return NewNumberFromFloat(Float64FromNode(n1) - Float64FromNode(n2))
	}

	log.Fatalf("Cannot subtract %v and %v", n1, n2)
	return &nodeNil{}
}

// MultiplyNodes multiplies two Nodes into one Node
func MultiplyNodes(n1 Node, n2 Node) Node {
	if n1.GetType() == nodeTypeString {
		if n2.GetType() == nodeTypeString {
			log.Fatalf("Cannot multiply tokens %v and %v", n1, n2)
		} else {
			return NewString(strings.Repeat(
				StringFromNode(n1), int(Float64FromNode(n2))))
		}
	} else {
		if n2.GetType() != nodeTypeString {
			return NewNumberFromFloat(Float64FromNode(n1) * Float64FromNode(n2))
		} else {
			return NewString(strings.Repeat(StringFromNode(n2), int(Float64FromNode(n1))))
		}
	}

	log.Fatalf("Could not multiply %v and %v", n1, n2)
	return &nodeNil{}
}

// RaiseNodes raises n1 to the power of n2 into one Node
func RaiseNodes(n1 Node, n2 Node) Node {
	result := 1.0
	num1, num2 := Float64FromNode(n1), Float64FromNode(n2)

	for i := 0; i < int(num2); i++ {
		result *= num1
	}

	return NewNumberFromFloat(result)
}

// AssignNodes assigns the variable represented
// by n1 to the value represented by n2
func AssignNode(n1 Node, n2 Node) Node {
	if n1.GetType() != nodeTypeVariable {
		log.Fatalf("Could not assign %v to %v", n2, n1)
	}

	if v, ok := n1.(*Variable); ok {
		SetVar(v.name, n2)
		return n2
	}

	log.Fatalf("Could not assign %v to %v", n2, n1)
	return &nodeNil{}
}

// DivideNodes divides two Nodes into one
func DivideNodes(n1 Node, n2 Node) Node {
	if n1.GetType() != nodeTypeString && n2.GetType() != nodeTypeString {
		return NewNumberFromFloat(Float64FromNode(n1) / Float64FromNode(n2))
	}

	log.Fatalf("Cannot divide %v and %v", n1, n2)
	return &nodeNil{}
}

// CompareNodes compares two Nodes opNode and returns a
// Boolean Node
func CompareNodes(n1 Node, n2 Node, tokOp Node) Node {
	var result bool
	var num1, num2 float64

	op, ok := tokOp.(*Operator)

	if !ok {
		log.Fatalf("Cannot compare with operator %v", tokOp)
	}

	// If the operator is not == or != then get
	// the numerical values from the Nodes
	if op.typ != opTypeEqualTo && op.typ != opTypeNotEqualTo {
		num1, num2 = Float64FromNode(n1), Float64FromNode(n2)
	}

	switch op.typ {
	case opTypeEqualTo:
		result = nodeIsEqualTo(n1, n2)
	case opTypeNotEqualTo:
		result = !nodeIsEqualTo(n1, n2)
	case opTypeLessThan:
		result = num1 < num2
	case opTypeLessThanOrEqualTo:
		result = num1 <= num2
	case opTypeGreaterThan:
		result = num1 > num2
	case opTypeGreaterThanOrEqualTo:
		result = num1 >= num2
	// TODO: make these their own operation category
	case opTypeAnd:
		result = num1 != 0.0 && num2 != 0.0
	case opTypeOr:
		result = num1 != 0.0 || num2 != 0.0
	}

	return NewBooleanFromBool(result)
}

// tokenIsEqual compares two Nodes and determine if they're equal
func nodeIsEqualTo(n1 Node, n2 Node) bool {
	if n1.GetType() == nodeTypeString || n2.GetType() == nodeTypeString {
		if n1.GetType() != n2.GetType() {
			return false
		}
	}

	if n1.GetType() == nodeTypeString {
		return StringFromNode(n1) == StringFromNode(n2)
	}

	return Float64FromNode(n1) == Float64FromNode(n2)
}
