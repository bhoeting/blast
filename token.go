package blast

import (
	"strconv"
	"strings"
	"unicode"
)

const (
	TokenTypeNumber   = 0
	TokenTypeString   = 1
	TokenTypeOperator = 2
	TokenTypeVariable = 3
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

		// 3 + "de rp"

		if r == rune(TokenDelimiter[0]) || index+1 == lCode {
			if strTok[0] == QuoteIdentifier[0] {
				strTok += string(r)
				if strTok[len(strTok)-1] != QuoteIdentifier[0] {
					continue
				}
			} else {
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
			return &Token{d, TokenTypeOperator}
		}
	}

	// Handle strings
	if strings.Index(strTok, QuoteIdentifier) != -1 {
		return &Token{strTok, TokenTypeString}
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
			return &Token{f, TokenTypeNumber}
		} else {
			i, _ := strconv.Atoi(strTok)
			return &Token{int64(i), TokenTypeNumber}
		}
	}

	return &Token{strTok, TokenTypeVariable}
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
