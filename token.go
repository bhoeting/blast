package blast

import (
	"fmt"
	"log"
	"strconv"
)

// Node is an interface
// for types created
// from Tokens
type Token interface {
	GetType() tokenType
	String() string
}

// nodeType is an int
// representation of a
// node type
type tokenType int

const (
	tokenTypeUnkown = iota
	tokenTypeFuncCall
	tokenTypeVariable
	tokenTypeNumber
	tokenTypeString
	tokenTypeParen
	tokenTypeBoolean
	tokenTypeOperator
	tokenTypeComma
	tokenTypeArgCount
	tokenTypeReserved
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

// GetType returns tokenTypeOperator
func (o *Operator) GetType() tokenType {
	return tokenTypeOperator
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

// GetType returns tokenTypeNumber
func (n *Number) GetType() tokenType {
	return tokenTypeNumber
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

// GetType returns tokenTypeBoolean
func (b *Boolean) GetType() tokenType {
	return tokenTypeBoolean
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

// GetType returns tokenTypeParen
func (p *Paren) GetType() tokenType {
	return tokenTypeParen
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

// GetType returns tokenTypeString
func (s *String) GetType() tokenType {
	return tokenTypeString
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

// GetType returns tokenTypeFuncCall
func (f *FunctionCall) GetType() tokenType {
	return tokenTypeFuncCall
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
	value Token
}

// GetType returns tokenTypeVariable
func (v *Variable) GetType() tokenType {
	return tokenTypeVariable
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

// GetType returns tokenTypeComma
func (c *Comma) GetType() tokenType {
	return tokenTypeComma
}

// String returns a `,`
func (c *Comma) String() string {
	return ","
}

// tokenNil is used for undefined behavior
type tokenNil struct{}

// GetType returns tokenTypeUnknown
func (n *tokenNil) GetType() tokenType {
	return tokenTypeUnkown
}

// String returns <NIL>
func (n *tokenNil) String() string {
	return "<NIL>"
}

// ArgCount is an int
// representing the amount
// of arguments passed in a
// function
type ArgCount int

// GetType returns a tokenTypeArgCount
func (a ArgCount) GetType() tokenType {
	return tokenTypeArgCount
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

// GetType returns tokenTypeReserved
func (r *Reserved) GetType() tokenType {
	return tokenTypeReserved
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

// NumberFromToken returns a float64
// from a Node
func NumberFromToken(token Token) float64 {
	switch token.GetType() {
	case tokenTypeNumber:
		return token.(*Number).value
	case tokenTypeBoolean:
		switch token.(*Boolean).typ {
		case booleanTypeTrue:
			return 1.0
		case booleanTypeFalse:
			return 0.0
		}
	default:
		log.Fatalf("Could not get numerical value from %v", token)
	}
	return 0.0
}

// StringFromToken returns a string a Node
func StringFromToken(token Token) string {
	switch token.GetType() {
	case tokenTypeNumber, tokenTypeBoolean:
		return token.String()
	case tokenTypeString:
		return token.(*String).value
	}

	log.Fatalf("Could not get string from %v", token)
	return ""
}

// BooleanFromToken returns a bool from a Node
func BooleanFromToken(token Token) bool {
	switch token.GetType() {
	case tokenTypeBoolean:
		switch token.(*Boolean).typ {
		case booleanTypeTrue:
			return true
		case booleanTypeFalse:
			return false
		}
	case tokenTypeNumber:
		if token.(*Number).value == 0.0 {
			return false
		}
		return true
	}

	log.Fatalf("Could not get boolean value from %v", token)
	return false
}
