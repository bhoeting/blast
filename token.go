package blast

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	TokenTypeInteger  = 0
	TokenTypeFloat    = 1
	TokenTypeString   = 2
	TokenTypeOperator = 3
	TokenTypeVariable = 4
)

const (
	OpTypeAddition       = 0
	OpTypeSubtraction    = 1
	OpTypeMultiplication = 2
	OpTypeDivision       = 3
)

type Token struct {
	data interface{}
	t    int
}

// NewToken array returns a new token array
func NewTokenArray(code string) []*Token {
	strTok := ""
	arrToks := make([]*Token, 0)
	lCode := len(code)

	for index, r := range code {
		if index+1 == lCode {
			strTok += string(r)
		}

		if code[index] == TokenDelimiter[0] || index+1 == lCode {
			if index > 0 && strTok[0] == QuoteIdentifier[0] &&
				strTok[len(strTok)-1] != QuoteIdentifier[0] {
				strTok += string(r)
				continue
			}

			arrToks = append(arrToks, NewToken(strTok))
			strTok = ""
		} else {
			strTok += string(r)
		}
	}

	return arrToks
}

// NewToken returns a new token
func NewToken(strTok string) *Token {
	// Handle operators
	if len(strTok) == 1 {
		d := -1
		switch string(strTok[0]) {
		case AdditionOp:
			d = OpTypeAddition
		case MultiplicationOp:
			d = OpTypeMultiplication
		case DivisionOp:
			d = OpTypeDivision
		case SubtractionOp:
			d = OpTypeSubtraction
		}

		if d != -1 {
			return &Token{int(d), TokenTypeOperator}
		}
	}

	// Handle strings
	if strings.Index(strTok, QuoteIdentifier) != -1 {
		return &Token{strTok[1 : len(strTok)-1], TokenTypeString}
	}

	// Handle numbers
	isNum, isDecimal := true, false
	for _, r := range strTok {
		if !unicode.IsNumber(r) && r != '-' {
			if r != '.' {
				isNum = false
				break
			} else {
				isDecimal = true
			}
		}
	}

	if isNum {
		if isDecimal {
			f, _ := strconv.ParseFloat(strTok, 64)
			return &Token{f, TokenTypeFloat}
		} else {
			i, _ := strconv.Atoi(strTok)
			return &Token{int64(i), TokenTypeInteger}
		}
	}

	return &Token{strTok, TokenTypeVariable}
}

func NewTokenFromVariable(v *Variable) *Token {
	switch v.t {
	case VarTypeInteger:
		if d, err := v.v.Integer(); err != nil {
			return &Token{d, TokenTypeInteger}
		}
	case VarTypeFloat:
		if d, err := v.v.Float(); err != nil {
			return &Token{d, TokenTypeFloat}
		}
	case VarTypeString:
		return &Token{v.v.String(), TokenTypeString}
	}

	return new(Token)
}

// SetVariableNameToValue
func (t *Token) SetVariableNameToValue(b *Blast) {
	v, err := b.vMap.Get(t.Str())

	if err == ErrVarNotFound {
		log.Fatalf("Could not get the value of variable %v", t.data)
		os.Exit(1)
	}

	t = NewTokenFromVariable(v)
}

// HandleToken performs an operation defined by oT on
// Tokens t1 and t2.
func HandleTokens(t1 *Token, t2 *Token, oT *Token, b *Blast) *Token {
	// Set the variable Tokens' data to the
	// value of the variable they represent
	if t1.t == TokenTypeVariable {
		t1.SetVariableNameToValue(b)
	}

	if t2.t == TokenTypeVariable {
		t2.SetVariableNameToValue(b)
	}

	// Handle string operations
	if t1.t == TokenTypeString {
		if t2.t == TokenTypeString || t2.t == TokenTypeFloat {
			if oT.Op() != OpTypeAddition {
				log.Fatal("Invalid operator on string")
				os.Exit(0)
			}
			return &Token{t1.String() + t2.String(), TokenTypeString}
		} else {
			switch oT.Op() {
			case OpTypeAddition:
				return &Token{t1.String() + t2.String(), TokenTypeString}
			case OpTypeMultiplication:
				return &Token{strings.Repeat(t1.String(), int(t2.Int())), TokenTypeInteger}
			default:
				log.Fatal("Invalid operator on string")
				os.Exit(0)
			}
		}
	} else {
		if t2.t == TokenTypeString {
			if t1.t == TokenTypeFloat {
				if oT.Op() != OpTypeAddition {
					log.Fatal("Invalid operator on string")
					os.Exit(0)
				}
				return &Token{t1.String() + t2.String(), TokenTypeString}
			} else {
				switch oT.Op() {
				case OpTypeAddition:
					return &Token{t1.String() + t2.String(), TokenTypeString}
				case OpTypeMultiplication:
					return &Token{strings.Repeat(t2.String(), int(t1.Int())), TokenTypeInteger}
				default:
					log.Fatal("Invalid operator on string")
					os.Exit(0)
				}
			}
		}
	}

	// Handle the Tokens as numbers
	value1, value2, tokType := 0.0, 0.0, -1
	if t1.t == TokenTypeInteger {
		if t2.t == TokenTypeInteger {
			value1 = float64(t1.Int())
			value2 = float64(t2.Int())
			tokType = TokenTypeInteger
		} else {
			value1 = float64(t1.Int())
			value2 = t2.Float()
			tokType = TokenTypeFloat
		}
	} else {
		if t2.t == TokenTypeInteger {
			value1 = t1.Float()
			value2 = float64(t2.Int())
			tokType = TokenTypeFloat
		} else {
			value1 = t1.Float()
			value2 = t2.Float()
			tokType = TokenTypeFloat
		}
	}

	if tokType == TokenTypeInteger {
		switch oT.Op() {
		case OpTypeAddition:
			return &Token{int64(value1) + int64(value2), TokenTypeInteger}
		case OpTypeSubtraction:
			return &Token{int64(value1) - int64(value2), TokenTypeInteger}
		case OpTypeMultiplication:
			return &Token{int64(value1) * int64(value2), TokenTypeInteger}
		case OpTypeDivision:
			return &Token{int64(value1) / int64(value2), TokenTypeInteger}
		}
	} else {
		switch oT.Op() {
		case OpTypeAddition:
			return &Token{value1 + value2, TokenTypeInteger}
		case OpTypeSubtraction:
			return &Token{value1 - value2, TokenTypeInteger}
		case OpTypeMultiplication:
			return &Token{value1 * value2, TokenTypeInteger}
		case OpTypeDivision:
			return &Token{value1 / value2, TokenTypeInteger}
		}
	}

	return new(Token)
}

func (t *Token) Int() int64 {
	if integer, ok := t.data.(int64); ok {
		return integer
	}

	return 0
}

func (t *Token) Float() float64 {
	if float, ok := t.data.(float64); ok {
		return float
	}

	return 0.0
}

func (t *Token) Op() int {
	if op, ok := t.data.(int); ok {
		return op
	}

	return -1
}

func (t *Token) Str() string {
	if str, ok := t.data.(string); ok {
		return str
	}

	return ""
}

func (t *Token) Var() string {
	return t.Str()
}

func (t *Token) Type() int {
	return t.t
}

func (t *Token) String() string {
	return fmt.Sprintf("%v", t.data)
}
