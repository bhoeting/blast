package blast

import (
	"fmt"
	"strconv"
)

type opType int

const (
	opTypeEqualTo = iota
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
	tokenTypeOperator
)

var tokenTypeStrings = map[tokenType]string{
	tokenTypeIf:       "if",
	tokenTypeInt:      "int",
	tokenTypeVar:      "var",
	tokenTypeEnd:      "end",
	tokenTypeNull:     "null",
	tokenTypeComma:    ",",
	tokenTypeFloat:    "float",
	tokenTypeParen:    "paren",
	tokenTypeString:   "string",
	tokenTypeUnkown:   "unkown",
	tokenTypeOperator: "operator",
}

type parenType int

const (
	parenTypeOpen = iota
	parenTypeClose
)

var operatorIdentifiers = map[string]opType{
	additionIdentifier:           opTypeAddition,
	equalityIdentifier:           opTypeEqualTo,
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
}

var varToToken = map[varType]int{
	varTypeString: tokenTypeString,
	varTypeFloat:  tokenTypeFloat,
	varTypeInt:    tokenTypeInt,
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
	}

	// Get either a float, int, or var token
	for _, lex := range ls.lexemes {
		// If a letter is detected, the token is a var token
		if lex.t == lexemeTypeLetter {
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
	tls := new(lexemeStream)
	tls.start = 0
	isString := false
	var lex *lexeme
	var i int

	addToken := func() {
		if tls.size > 0 {
			tls.end = i
			ts.add(newTokenFromLexemeStream(tls))
			tls.clear()
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

			// Add the actual paren token
			tls.push(lex)
			addToken()

			continue
		}

		if lex.t == lexemeTypeOperator {
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
		str += t.string() + ","
	}

	return str
}

// add adds a token to the tokenStream
func (ts *tokenStream) add(token *token) {
	ts.tokens = append(ts.tokens, token)
	ts.size++
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

// evaluateTokens performs an operation on two tokens
func evaluateTokens(t1 *token, t2 *token, op *token) *token {
	return tokenNull
}

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
	default:
		return fmt.Sprintf("%v", t.data)
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

func (t *token) variable() *variable {
	return t.data.(*variable)
}
