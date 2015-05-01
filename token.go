package blast

import (
	"fmt"
	"log"
	"strconv"
)

// Node is an interface
// for types created
// from nodes
type Node interface {
	GetType() nodeType
	String() string
}

// nodeType is an int
// representation of a
// node type
type nodeType int

const (
	nodeTypeUnkown = iota
	nodeTypeFuncCall
	nodeTypeVariable
	nodeTypeNumber
	nodeTypeString
	nodeTypeParen
	nodeTypeBoolean
	nodeTypeOperator
	nodeTypeComma
	nodeTypeArgCount
	nodeTypeReserved
)

// Operator is a struct
// that stores an opType
type Operator struct {
	typ opType
}

// opType is an int representing
// an operator type
type opType int

const (
	opTypeAddition = iota
	opTypeSubtraction
	opTypeMultiplication
	opTypeDivision
	opTypeExponent
	opTypeAssignment
	opTypeEqualTo
	opTypeNotEqualTo
	opTypeLessThan
	opTypeLessThanOrEqualTo
	opTypeGreaterThan
	opTypeGreaterThanOrEqualTo
	opTypeAnd
	opTypeOr
	opTypeArrow
)

// operatorKey is used to get an
// operator type from a string
var operatorKey = map[string]opType{
	"+":  opTypeAddition,
	"-":  opTypeSubtraction,
	"*":  opTypeMultiplication,
	"/":  opTypeDivision,
	"^":  opTypeExponent,
	"=":  opTypeAssignment,
	"==": opTypeEqualTo,
	"!=": opTypeNotEqualTo,
	"<":  opTypeLessThan,
	"<=": opTypeLessThanOrEqualTo,
	">":  opTypeGreaterThan,
	">=": opTypeGreaterThanOrEqualTo,
	"&&": opTypeAnd,
	"||": opTypeOr,
	"->": opTypeArrow,
}

// operatorStrings is used to get a
// string representation from
// an opType
var operatorStrings = map[opType]string{
	opTypeAddition:             "+",
	opTypeSubtraction:          "-",
	opTypeMultiplication:       "*",
	opTypeDivision:             "/",
	opTypeExponent:             "^",
	opTypeAssignment:           "=",
	opTypeEqualTo:              "==",
	opTypeNotEqualTo:           "!=",
	opTypeLessThan:             "<",
	opTypeLessThanOrEqualTo:    "<=",
	opTypeGreaterThan:          ">",
	opTypeGreaterThanOrEqualTo: ">=",
	opTypeAnd:                  "&&",
	opTypeOr:                   "||",
	opTypeArrow:                "->",
}

// GetType returns nodeTypeOperator
func (o *Operator) GetType() nodeType {
	return nodeTypeOperator
}

// String returns an Operator as a string
func (o *Operator) String() string {
	return operatorStrings[o.typ]
}

// NewOperator returns a new Operator
func NewOperator(strOp string) *Operator {
	operator := new(Operator)

	if ot, ok := operatorKey[strOp]; ok {
		operator.typ = ot
	} else {
		log.Fatalf("Could not parse operator %s", strOp)
	}

	return operator
}

// Number is a struct that stores a float64
type Number struct {
	value float64
}

// GetType returns nodeTypeNumber
func (n *Number) GetType() nodeType {
	return nodeTypeNumber
}

// String returns a Number as a string
func (n *Number) String() string {
	return fmt.Sprintf("%v", n.value)
}

// NewNumber returns a new Number
func NewNumber(strNum string) *Number {
	number := new(Number)
	value, err := strconv.ParseFloat(strNum, 64)

	if err != nil {
		log.Fatal(err.Error())
	}

	number.value = value
	return number
}

// NewNumberFromFloat returns a Number from a float64
func NewNumberFromFloat(value float64) *Number {
	number := new(Number)
	number.value = value
	return number
}

// Boolean is a struct that
// stores a booleanType
type Boolean struct {
	typ booleanType
}

// booleanType is an int that
// represents a booleanType
type booleanType int

const (
	booleanTypeTrue = iota
	booleanTypeFalse
)

// GetType returns nodeTypeBoolean
func (b *Boolean) GetType() nodeType {
	return nodeTypeBoolean
}

// String returns a Boolean as a string
func (b *Boolean) String() string {
	switch b.typ {
	case booleanTypeTrue:
		return "true"
	case booleanTypeFalse:
		return "false"
	}
	return "false"
}

// NewBoolean returns a new Boolean
func NewBoolean(strBool string) *Boolean {
	boolean := new(Boolean)

	switch strBool {
	case "true":
		boolean.typ = booleanTypeTrue
	case "false":
		boolean.typ = booleanTypeFalse
	default:
		log.Fatal("Could not parse boolean")
	}

	return boolean
}

// NewBooleanFromBool returns a new Boolean from a bool
func NewBooleanFromBool(value bool) *Boolean {
	boolean := new(Boolean)

	switch value {
	case true:
		boolean.typ = booleanTypeTrue
	default:
		boolean.typ = booleanTypeFalse
	}

	return boolean
}

// Paren is a struct with a parenTytpe
type Paren struct {
	typ parenType
}

// parenType is an int
// of a paren type
type parenType int

const (
	parenTypeOpen = iota
	parenTypeClose
	parenTypeNil
)

// GetType returns nodeTypeParen
func (p *Paren) GetType() nodeType {
	return nodeTypeParen
}

// String returns a Paren as a string
func (p *Paren) String() string {
	switch p.typ {
	case parenTypeOpen:
		return "("
	default:
		return ")"
	}
}

// NewParen returns a new Paren
func NewParen(strParen string) *Paren {
	paren := new(Paren)

	switch strParen {
	case "(":
		paren.typ = parenTypeOpen
	case ")":
		paren.typ = parenTypeClose
	default:
		log.Fatal("Could not parse paren")
	}

	return paren
}

// String is a struct that
// stores a string
type String struct {
	value string
}

// GetType returns nodeTypeString
func (s *String) GetType() nodeType {
	return nodeTypeString
}

// String returns the string wrapped in quotes
func (s *String) String() string {
	return fmt.Sprintf("\"%s\"", s.value)
}

// NewString returns a new string
func NewString(s string) *String {
	str := new(String)
	str.value = s
	return str
}

// FunctionCall is a struct
// that stores the name of
// a function
type FunctionCall struct {
	name string
}

// GetType returns nodeTypeFuncCall
func (f *FunctionCall) GetType() nodeType {
	return nodeTypeFuncCall
}

// String returns the func name with parens
// after it
func (f *FunctionCall) String() string {
	return fmt.Sprintf("%s()", f.name)
}

// NewFunctionCall returns a new FunctionCall
func NewFunctionCall(strName string) *FunctionCall {
	f := new(FunctionCall)
	f.name = strName
	return f
}

// Variable is a struct that
// stores a name and Node
type Variable struct {
	name  string
	value Node
}

// GetType returns nodeTypeVariable
func (v *Variable) GetType() nodeType {
	return nodeTypeVariable
}

// String returns the variable name
func (v *Variable) String() string {
	return v.name
}

// NewVariable returns a new Variable
func NewVariable(strName string) *Variable {
	v := new(Variable)
	v.name = strName
	return v
}

// Comma is an empty struct
type Comma struct{}

// NewComma returns a new Comma
func NewComma() *Comma {
	return new(Comma)
}

// GetType returns nodeTypeComma
func (c *Comma) GetType() nodeType {
	return nodeTypeComma
}

// String returns a `,`
func (c *Comma) String() string {
	return ","
}

// nodeNil is used for undefined behavior
type nodeNil struct{}

// GetType returns nodeTypeUnknown
func (n *nodeNil) GetType() nodeType {
	return nodeTypeUnkown
}

// String returns <NIL>
func (n *nodeNil) String() string {
	return "<NIL>"
}

// ArgCount is an int
// representing the amount
// of arguments passed in a
// function
type ArgCount int

// GetType returns a nodeTypeArgCount
func (a ArgCount) GetType() nodeType {
	return nodeTypeArgCount
}

// String returns the ArgCount as a string
func (a ArgCount) String() string {
	return fmt.Sprintf("%d", int(a))
}

// NewArgCount returns a new ArgCount
func NewArgCount(count int) ArgCount {
	return ArgCount(count)
}

// Reserved is a struct
// for reserved words and
// stores the reserved word
type Reserved struct {
	value string
}

// GetType returns nodeTypeReserved
func (r *Reserved) GetType() nodeType {
	return nodeTypeReserved
}

// String returns <RESERVED>
func (r *Reserved) String() string {
	return "<RESERVED>"
}

// NewReserved returns a new
// reserved word
func NewReserved() *Reserved {
	return &Reserved{}
}

// NumberFromNode returns a float64
// from a Node
func Float64FromNode(node Node) float64 {
	switch node.GetType() {
	case nodeTypeNumber:
		return node.(*Number).value
	case nodeTypeBoolean:
		switch node.(*Boolean).typ {
		case booleanTypeTrue:
			return 1.0
		case booleanTypeFalse:
			return 0.0
		}
	default:
		log.Fatalf("Could not get numerical value from %v", node)
	}
	return 0.0
}

// StringFromNode returns a string a Node
func StringFromNode(node Node) string {
	switch node.GetType() {
	case nodeTypeNumber, nodeTypeBoolean:
		return node.String()
	case nodeTypeString:
		return node.(*String).value
	}

	log.Fatalf("Could not get string from %v", node)
	return ""
}

// BooleanFromNode returns a bool from a Node
func BooleanFromNode(node Node) bool {
	switch node.GetType() {
	case nodeTypeBoolean:
		switch node.(*Boolean).typ {
		case booleanTypeTrue:
			return true
		case booleanTypeFalse:
			return false
		}
	case nodeTypeNumber:
		if node.(*Number).value == 0.0 {
			return false
		}
		return true
	}

	log.Fatalf("Could not get boolean value from %v", node)
	return false
}
