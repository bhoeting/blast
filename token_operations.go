package blast

import (
	"log"
	"strings"
)

// AddNodes adds two Nodes into one Node
func AddTokens(t1 Token, t2 Token) Token {
	var result Token

	if t1.GetType() == tokenTypeString || t2.GetType() == tokenTypeString {
		result = NewString(StringFromToken(t1) + StringFromToken(t2))
	} else {
		result = NewNumberFromFloat(NumberFromToken(t1) + NumberFromToken(t2))
	}

	return result
}

// SubtractNodes subtracts two Nodes into one Node
func SubtractTokens(t1 Token, t2 Token) Token {
	if t1.GetType() != tokenTypeString && t2.GetType() != tokenTypeString {
		return NewNumberFromFloat(NumberFromToken(t1) - NumberFromToken(t2))
	}

	log.Fatalf("Cannot subtract %v and %v", t1, t2)
	return &tokenNil{}
}

// MultiplyNodes multiplies two Nodes into one Node
func MultiplyTokens(t1 Token, t2 Token) Token {
	if t1.GetType() == tokenTypeString {
		if t2.GetType() == tokenTypeString {
			log.Fatalf("Cannot multiply tokens %v and %v", t1, t2)
		} else {
			return NewString(strings.Repeat(
				StringFromToken(t1), int(NumberFromToken(t2))))
		}
	} else {
		if t2.GetType() != tokenTypeString {
			return NewNumberFromFloat(NumberFromToken(t1) * NumberFromToken(t2))
		} else {
			return NewString(strings.Repeat(StringFromToken(t2), int(NumberFromToken(t1))))
		}
	}

	log.Fatalf("Could not multiply %v and %v", t1, t2)
	return &tokenNil{}
}

// RaiseNodes raises n1 to the power of n2 into one Node
func RaiseTokens(t1 Token, t2 Token) Token {
	result := 1.0
	num1, num2 := NumberFromToken(t1), NumberFromToken(t2)

	for i := 0; i < int(num2); i++ {
		result *= num1
	}

	return NewNumberFromFloat(result)
}

// AssignNodes assigns the variable represented
// by n1 to the value represented by n2
func AssignTokens(t1 Token, t2 Token) Token {
	if t1.GetType() != tokenTypeVariable {
		log.Fatalf("Could not assign %v to %v", t2, t1)
	}

	if v, ok := t1.(*Variable); ok {
		SetVar(v.name, t2)
		return t2
	}

	log.Fatalf("Could not assign %v to %v", t2, t1)
	return &tokenNil{}
}

// DivideNodes divides two Nodes into one
func DivideTokens(t1 Token, t2 Token) Token {
	if t1.GetType() != tokenTypeString && t2.GetType() != tokenTypeString {
		return NewNumberFromFloat(NumberFromToken(t1) / NumberFromToken(t2))
	}

	log.Fatalf("Cannot divide %v and %v", t1, t2)
	return &tokenNil{}
}

// CompareNodes compares two Nodes opNode and returns a
// Boolean Node
func CompareTokens(t1 Token, t2 Token, tokOp Token) Token {
	var result bool
	var num1, num2 float64

	op, ok := tokOp.(*Operator)

	if !ok {
		log.Fatalf("Cannot compare with operator %v", tokOp)
	}

	// If the operator is not == or != then get
	// the numerical values from the Nodes
	if op.typ != opTypeEqualTo && op.typ != opTypeNotEqualTo {
		num1, num2 = NumberFromToken(t1), NumberFromToken(t2)
	}

	switch op.typ {
	case opTypeEqualTo:
		result = tokenIsEqualTo(t1, t2)
	case opTypeNotEqualTo:
		result = !tokenIsEqualTo(t1, t2)
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
func tokenIsEqualTo(t1 Token, t2 Token) bool {
	if t1.GetType() == tokenTypeString || t2.GetType() == tokenTypeString {
		if t1.GetType() != t2.GetType() {
			return false
		}
	}

	if t1.GetType() == tokenTypeString {
		return StringFromToken(t1) == StringFromToken(t2)
	}

	return NumberFromToken(t1) == NumberFromToken(t2)
}
