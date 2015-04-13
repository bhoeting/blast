package blast

import (
	"fmt"
	"log"
	"strconv"
)

type opType int

const (
	opTypeAnd = iota
	opTypeEqualTo
	opTypeAddition
	opTypeDivision
	opTypeLessThan
	opTypeAssignment
	opTypeGreaterThan
	opTypeSubtraction
	opTypeMultiplication
	opTypeLessThanOrEqualTo
	opTypeGreaterThanOrEqualTo
)

type tokenType int

const (
	tokenTypeIf = iota
	tokenTypeInt
	tokenTypeVar
	tokenTypeEnd
	tokenTypeNull
	tokenTypeComma
	tokenTypeFloat
	tokenTypeParen
	tokenTypeUnkown
	tokenTypeString
	tokenTypeReturn
	tokenTypeBoolean
	tokenTypeOperator
	tokenTypeFunction
	tokenTypeFunctionCall
)

var tokenTypeStrings = map[tokenType]string{
	tokenTypeIf:           "if",
	tokenTypeInt:          "int",
	tokenTypeVar:          "var",
	tokenTypeEnd:          "end",
	tokenTypeNull:         "null",
	tokenTypeComma:        ",",
	tokenTypeFloat:        "float",
	tokenTypeParen:        "paren",
	tokenTypeString:       "string",
	tokenTypeUnkown:       "unkown",
	tokenTypeReturn:       "return",
	tokenTypeBoolean:      "boolean",
	tokenTypeOperator:     "operator",
	tokenTypeFunction:     "function",
	tokenTypeFunctionCall: "function_call",
}

type parenType int

const (
	parenTypeOpen = iota
	parenTypeClose
)

var operatorIdentifiers = map[string]opType{
	andIdentifier:                opTypeAnd,
	equalityIdentifier:           opTypeEqualTo,
	additionIdentifier:           opTypeAddition,
	divisionIdentifier:           opTypeDivision,
	lessThanIdentifier:           opTypeLessThan,
	assignmentIdentifier:         opTypeAssignment,
	greaterThanIdentifier:        opTypeGreaterThan,
	subtractionIdentifier:        opTypeSubtraction,
	multiplicationIdentifier:     opTypeMultiplication,
	lessThanOrEqualIdentifier:    opTypeLessThanOrEqualTo,
	greaterThanOrEqualIdentifier: opTypeGreaterThanOrEqualTo,
}

type token struct {
	data  interface{}
	start int
	end   int
	t     tokenType
}

type tokenStream struct {
	tokens []*token
	size   int
	isRPN  bool
}

var varToToken = map[varType]tokenType{
	varTypeString:  tokenTypeString,
	varTypeFloat:   tokenTypeFloat,
	varTypeInt:     tokenTypeInt,
	varTypeBoolean: tokenTypeBoolean,
}

// newTokenFromLexemeStream returns a new token
// from the lexemes in the lexemeStream
func newTokenFromLexemeStream(ls *lexemeStream) *token {
	tType := tokenTypeInt
	tokenData := ls.string()
	start, end := ls.start, ls.end

	// If the stream begins and ends with quotes,
	// return a string token
	if tokenData[0] == quoteIdentifier[0] &&
		tokenData[len(tokenData)-1] == quoteIdentifier[0] {
		return newToken(tokenData[1:len(tokenData)-1], start, end, tokenTypeString)
	}

	// Check for one-lengthed tokens
	if ls.size == 1 {
		switch ls.get(0).t {
		case lexemeTypeParen:
			return newParenToken(tokenData, start, end)
		case lexemeTypeComma:
			return newBasicToken(start, end, tokenTypeComma)
		}
	}

	// Check for reserved words
	switch tokenData {
	case ifIdentifier:
		return newBasicToken(start, end, tokenTypeIf)
	case functionIdentifier:
		return newBasicToken(start, end, tokenTypeFunction)
	case endIdentifier:
		return newBasicToken(start, end, tokenTypeEnd)
	case trueIdentifier, falseIdentifier:
		return newBooleanToken(tokenData, start, end)
	case returnIdentifier:
		return newBasicToken(start, end, tokenTypeReturn)
	}

	// Get either a float, int, or var token
	for _, lex := range ls.lexemes {
		// If a letter is detected, the token is a var token
		if lex.t == lexemeTypeLetter {
			if _, err := B.getVariable(tokenData); err == nil {
				return newToken(tokenData, start, end, tokenTypeVar)
			}

			if _, err := B.getFunction(tokenData); err == nil {
				return newToken(tokenData, start, end, tokenTypeFunctionCall)
			}

			return newToken(tokenData, start, end, tokenTypeVar)
		}

		if lex.t == lexemeTypeOperator {
			tType = tokenTypeOperator
			break
		}

		// If the token is an int and a decimal is detected, the
		// token is a float token
		if lex.t == lexemeTypeDecimalDigit && tType == tokenTypeInt {
			tType = tokenTypeFloat
		}
	}

	switch tType {
	case tokenTypeInt:
		return newIntegerToken(tokenData, start, end)
	case tokenTypeFloat:
		return newFloatToken(tokenData, start, end)
	}

	return newOperatorToken(tokenData, start, end)
}

// newToken returns a new token
func newToken(data interface{}, start int, end int, t tokenType) *token {
	tok := new(token)
	tok.data = data
	tok.start = start
	tok.end = end
	tok.t = t
	return tok
}

// newTokenStream creates a tokenStream from a lexemeStream
func newTokenStream(ls *lexemeStream) *tokenStream {
	ts := new(tokenStream)
	ts.isRPN = false
	tls := new(lexemeStream)
	tls.start = 0
	isString := false
	unkownIndex := -1
	var lex *lexeme
	var i int

	addToken := func() {
		if tls.size > 0 {
			newTok := newTokenFromLexemeStream(tls)
			ts.add(newTok)

			// If a token is unkown, it is either
			// a variable or a function identifier.
			// We will store the index of this token
			// to determine if it is a function when
			// we have more information
			if newTok.t == tokenTypeUnkown {
				unkownIndex = ts.size - 1
			}

			tls.clear()
			tls.end = i
			tls.start = i + 1
		}
	}

	for i, lex = range ls.lexemes {
		if lex.t == lexemeTypeQuote {
			if !isString && tls.size > 0 {
				addToken()
			}

			isString = !isString
		}

		if isString {
			tls.push(lex)
			continue
		}

		if lex.t == lexemeTypeEOL ||
			lex.t == lexemeTypeSpace {
			addToken()
			continue
		}

		if lex.t == lexemeTypeComma ||
			lex.t == lexemeTypeParen {

			// Add token that was
			// being built already
			addToken()

			// If an open paren is detected, then we know
			// the unknown token must be a function call
			if unkownIndex != -1 && lex.text == '(' {
				ts.tokens[unkownIndex].t = tokenTypeFunctionCall
				unkownIndex = -1
			}

			// Add the actual paren token
			tls.push(lex)
			addToken()

			continue
		}

		if lex.t == lexemeTypeOperator {
			unkownIndex = -1
			if tls.top().t == lexemeTypeOperator {
				tls.push(lex)
			} else {
				addToken()
				tls.push(lex)
			}
		} else {
			if tls.top().t != lexemeTypeOperator {
				tls.push(lex)
			} else {
				addToken()
				tls.push(lex)
			}
		}
	}

	addToken()

	return ts
}

// get returns the token at the specified index
func (ts *tokenStream) get(index int) *token {
	return ts.tokens[index]
}

// string returns a string representation
// of a tokenSream
func (ts *tokenStream) string() string {
	str := ""

	for _, t := range ts.tokens {
		str += t.string() + "`"
	}

	return str
}

// add adds a token to the tokenStream
func (ts *tokenStream) add(token *token) {
	ts.tokens = append(ts.tokens, token)
	ts.size++
}

// clear removes the tokens from
// a tokenStream
func (ts *tokenStream) clear() *tokenStream {
	ts.tokens = nil
	ts.size = 0

	return ts
}

// prepend inserts a token at the from of the stream
func (ts *tokenStream) prepend(t *token) {
	ts.tokens = append([]*token{t}, ts.tokens...)
}

// newBasicToken returns a new token
// that doesn't have any special data
func newBasicToken(start int, end int, tType tokenType) *token {
	return newToken(0, start, end, tType)
}

// newParenToken returns a new paren token
func newParenToken(text string, start int, end int) *token {
	if text == openParenIdentifier {
		return newToken(parenTypeOpen, start, end, tokenTypeParen)
	} else {
		return newToken(parenTypeClose, start, end, tokenTypeParen)
	}
}

// newFloatToken returns a new float token
func newFloatToken(text string, start int, end int) *token {
	data, _ := strconv.ParseFloat(text, 64)
	return newToken(data, start, end, tokenTypeFloat)
}

// newIntegerToken returns a new integer token
func newIntegerToken(text string, start int, end int) *token {
	data, _ := strconv.Atoi(text)
	return newToken(data, start, end, tokenTypeInt)
}

// newOperatorToken returns a new operator token
func newOperatorToken(text string, start int, end int) *token {
	return newToken(operatorIdentifiers[text], start, end, tokenTypeOperator)
}

// newBooleanToken returns a new boolean token
func newBooleanToken(text string, start int, end int) *token {
	data := false

	if text == trueIdentifier {
		data = true
	}

	return newToken(data, start, end, tokenTypeBoolean)
}

// evaluateToken evalutes a single token
func evaluateToken(t *token) *token {
	if t.t == tokenTypeUnkown {
		t = resolveUnknownToken(t)
	}

	if t.t == tokenTypeVar {
		return getTokenFromVariableToken(t)
	}

	return t
}

// evaluateTokens performs an operation on two tokens
func evaluateTokens(t1 *token, t2 *token, op *token) *token {
	if op.opType() != opTypeAssignment {
		t1 = evaluateToken(t1)
	} else {
		if t1.t == tokenTypeUnkown {
			t1 = resolveUnknownToken(t1)
		}
	}

	t2 = evaluateToken(t2)

	switch op.opType() {
	case opTypeAddition:
		return addTokens(t1, t2)
	case opTypeSubtraction:
		return subtractTokens(t1, t2)
	case opTypeMultiplication:
		return multiplyTokens(t1, t2)
	case opTypeDivision:
		return divideTokens(t1, t2)
	case opTypeAssignment:
		t1.t = tokenTypeVar
		return assignTokens(t1, t2)
	default:
		return compareTokens(t1, t2, op.opType())
	}

	return tokenNull
}

// resolveUnkownToken determines if the
// unknown token is a variable or function
// call based on the declared functions
func resolveUnknownToken(t *token) *token {
	_, err := B.getFunction(t.data.(string))

	// If a function was not found
	// the token is a variable
	if err != nil {
		t.t = tokenTypeFunctionCall
	} else {
		t.t = tokenTypeVar
	}

	return t
}

// getTokenFromVariableToken returns a new token
// with the value of the variable specified by
// the token
func getTokenFromVariableToken(t *token) *token {
	v, err := B.getVariable(t.data.(string))

	if err != nil {
		log.Fatalf(err.Error())
	}

	return newToken(v.data, t.start, t.end, varToToken[v.t])
}

// string returns a string representation
// of a token
func (t *token) string() string {
	switch t.t {
	case tokenTypeOperator:
		for opIdentifier, oType := range operatorIdentifiers {
			if oType == t.data.(opType) {
				return fmt.Sprintf("%v", opIdentifier)
			}
		}
		return fmt.Sprintf("%v", "Optype uknown")
	case tokenTypeParen:
		if t.data.(int) == parenTypeClose {
			return fmt.Sprint(")")
		}
		return fmt.Sprint("(")
	case tokenTypeString:
		return fmt.Sprintf("\"%v\"", t.data)
	case tokenTypeIf:
		return fmt.Sprintf("%v", ifIdentifier)
	case tokenTypeEnd:
		return fmt.Sprintf("%v", endIdentifier)
	case tokenTypeFunction:
		return fmt.Sprintf("%v", functionIdentifier)
	case tokenTypeReturn:
		return fmt.Sprintf("%v", returnIdentifier)
	case tokenTypeComma:
		return fmt.Sprintf("%v", ",")
	case tokenTypeFunctionCall:
		return fmt.Sprintf("%v()", t.data)
	}

	return fmt.Sprintf("%v", t.data)
}

func (t *token) str() string {
	return t.data.(string)
}

func (t *token) integer() int {
	return t.data.(int)
}

func (t *token) float() float64 {
	return t.data.(float64)
}

func (t *token) opType() opType {
	return t.data.(opType)
}

func (t *token) parenType() int {
	return t.data.(int)
}

func (t *token) number() float64 {
	if n, ok := t.data.(int); ok {
		return float64(n)
	}

	if n, ok := t.data.(bool); ok {
		if n {
			return 1
		}
		return 0
	}

	return t.data.(float64)
}

func (t *token) boolean() bool {
	return t.data.(bool)
}
