package blast

import (
	"fmt"
	"log"
	"strconv"
)

type Token interface {
	GetType() tokenType
	String() string
}

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

type Operator struct {
	typ opType
}

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

func (o *Operator) GetType() tokenType {
	return tokenTypeOperator
}

func (o *Operator) String() string {
	return operatorStrings[o.typ]
}

func NewOperator(strOp string) *Operator {
	operator := new(Operator)

	if ot, ok := operatorKey[strOp]; ok {
		operator.typ = ot
	} else {
		log.Fatalf("Could not parse operator %s", strOp)
	}

	return operator
}

type Number struct {
	value float64
}

func (n *Number) GetType() tokenType {
	return tokenTypeNumber
}

func (n *Number) String() string {
	return fmt.Sprintf("%v", n.value)
}

func NewNumber(strNum string) *Number {
	number := new(Number)
	value, err := strconv.ParseFloat(strNum, 64)

	if err != nil {
		log.Fatal(err.Error())
	}

	number.value = value
	return number
}

func NewNumberFromFloat(value float64) *Number {
	number := new(Number)
	number.value = value
	return number
}

type Boolean struct {
	typ booleanType
}

type booleanType int

const (
	booleanTypeTrue = iota
	booleanTypeFalse
)

func (b *Boolean) GetType() tokenType {
	return tokenTypeBoolean
}

func (b *Boolean) String() string {
	switch b.typ {
	case booleanTypeTrue:
		return "true"
	case booleanTypeFalse:
		return "false"
	}
	return "false"
}

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

type Paren struct {
	typ parenType
}

type parenType int

const (
	parenTypeOpen = iota
	parenTypeClose
	parenTypeNil
)

func (p *Paren) GetType() tokenType {
	return tokenTypeParen
}

func (p *Paren) String() string {
	switch p.typ {
	case parenTypeOpen:
		return "("
	default:
		return ")"
	}
}

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

type String struct {
	value string
}

func (s *String) GetType() tokenType {
	return tokenTypeString
}

func (s *String) String() string {
	return fmt.Sprintf("\"%s\"", s.value)
}

func NewString(s string) *String {
	str := new(String)
	str.value = s
	return str
}

type FunctionCall struct {
	name string
}

func (f *FunctionCall) GetType() tokenType {
	return tokenTypeFuncCall
}

func (f *FunctionCall) String() string {
	return fmt.Sprintf("%s()", f.name)
}

func NewFunctionCall(strName string) *FunctionCall {
	f := new(FunctionCall)
	f.name = strName
	return f
}

type Variable struct {
	name  string
	value Token
}

func (v *Variable) GetType() tokenType {
	return tokenTypeVariable
}

func (v *Variable) String() string {
	return v.name
}

func NewVariable(strName string) *Variable {
	v := new(Variable)
	v.name = strName
	return v
}

type Comma struct{}

func NewComma() *Comma {
	return new(Comma)
}

func (c *Comma) GetType() tokenType {
	return tokenTypeComma
}

func (c *Comma) String() string {
	return ","
}

type tokenNil struct{}

func (n *tokenNil) GetType() tokenType {
	return tokenTypeUnkown
}

func (n *tokenNil) String() string {
	return "<NIL>"
}

type ArgCount int

func (a ArgCount) GetType() tokenType {
	return tokenTypeArgCount
}

func (a ArgCount) String() string {
	return fmt.Sprintf("%d", int(a))
}

func NewArgCount(count int) ArgCount {
	return ArgCount(count)
}

type Reserved struct {
	value string
}

func (r *Reserved) GetType() tokenType {
	return tokenTypeReserved
}

func (r *Reserved) String() string {
	return "<RESERVED>"
}

func NewReserved() *Reserved {
	return &Reserved{}
}

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
