package blast

import (
	"fmt"
	"strconv"
	"unicode"
)

const (
	tokenTypeNull       = 1
	tokenTypeOp         = 3
	tokenTypeChar       = 2
	tokenTypeParen      = 1
	tokenTypeQuote      = 4
	tokenTypeSpace      = 5
	tokenTypeString     = 6
	tokenTypeInt        = 7
	tokenTypeFloat      = 8
	tokenTypeVar        = 9
	tokenTypeAssignment = 10
)

// Token stores data
// and a token type
type token struct {
	data interface{}
	t    int
}

// tokenStram is a token
// slice wrapper
type tokenStream struct {
	tokens []*token
	size   int
	index  int
}

// newTokenStream returns a new tokenStream
func newTokenStream(strTokens string) *tokenStream {
	ts := new(tokenStream)
	ts.tokens = parseTokens(strTokens)
	ts.size = len(ts.tokens)
	ts.index = 0
	return ts
}

// next returns the next item in the stream
// and a boolean stating if an item was available
func (ts *tokenStream) next() (*token, bool) {
	if ts.index > ts.size-1 {
		return new(token), false
	}

	tok := ts.curr()
	ts.index++

	return tok, true
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

// combine like tokens into a single token,
// example: a stream of [3, 0, 0] will become [300]
func (ts *tokenStream) combine() *tokenStream {
	// The token slice that will replace the current token slice
	var replacementTokens []*token

	// A slice of tokens that are to be combined
	var tokensToBeCombined []*token

	// The index of the last token to be combined
	prevTokensEndIndex := 0

	ts.each(func(token *token, index int) {
		// If the token is a char or quote, prepare them to be combined
		if token.t == tokenTypeChar || token.t == tokenTypeQuote {
			tokensToBeCombined = append(tokensToBeCombined, token)
		}

		// If (the token is an operator or the end of the slice has been reached or
		// the token is a closing paren) and the tokensToBeCombined has been added
		// to since the combining of tokens
		if (token.t == tokenTypeOp || index == ts.size-1 ||
			(token.t == tokenTypeParen && token.data == parenTypeClose)) &&
			len(tokensToBeCombined) != prevTokensEndIndex {

			// Add the combined token to the replacementTokens
			replacementTokens = append(replacementTokens,
				combineTokens(tokensToBeCombined[prevTokensEndIndex:]))

			prevTokensEndIndex = len(tokensToBeCombined)
		}

		// Add operators and parens AFTER any combinations occure
		if token.t == tokenTypeOp || token.t == tokenTypeParen {
			replacementTokens = append(replacementTokens, token)
		}
	})

	ts.tokens = replacementTokens
	ts.size = len(replacementTokens) - 1

	return ts
}

// combineTokens combines the items in a token
// slice into one token with a common token type
func combineTokens(tokens []*token) *token {
	strCombinedTok := ""
	isNum, isDecimal, isStr := true, false, false

	// If the first and last tokens are quotes
	if tokens[0].t == tokenTypeQuote &&
		tokens[len(tokens)-1].t == tokenTypeQuote {
		isStr = true
	}

	for _, t := range tokens {
		if strTok, ok := t.data.(string); ok {
			strCombinedTok += strTok
			if !unicode.IsNumber(rune(strTok[0])) {
				if strTok != subtractionIdentifier {
					if strTok != decimalPointIdentifier {
						isNum = false
					} else {
						isDecimal = true
					}
				}
			}
		}
	}

	if isStr {
		return newToken(strCombinedTok, tokenTypeString)
	}

	if isNum {
		if isDecimal {
			f, _ := strconv.ParseFloat(strCombinedTok, 64)
			return newToken(f, tokenTypeFloat)
		}

		i, _ := strconv.Atoi(strCombinedTok)
		return newToken(i, tokenTypeInt)
	}

	return newToken(strCombinedTok, tokenTypeVar)
}

// each is a helper method that runs a function
// on each of the tokens in tokenStream's token slice
func (ts *tokenStream) each(f func(tok *token, index int)) {
	for next, ok := ts.next(); ok == true; next, ok = ts.next() {
		f(next, ts.index-1)
	}

	ts.index = 0
}

func (ts *tokenStream) curr() *token {
	return ts.tokens[ts.index]
}

func newToken(data interface{}, t int) *token {
	tok := new(token)
	tok.data = data
	tok.t = t
	return tok
}

func parseTokens(code string) []*token {
	var tokens []*token

	for _, r := range code {
		tokens = append(tokens, parseToken(string(r)))
	}

	return tokens
}

func evaluateTokens(t1 *token, t2 *token, op *token) *token {
	switch op.opType() {
	case opTypeAddition:
		return addTokens(t1, t2)
	case opTypeSubtraction:
		return subtractTokens(t1, t2)
	case opTypeMultiplication:
		return multiplyTokens(t1, t2)
	case opTypeDivision:
		return divideTokens(t1, t2)
	}
	return tokenNull
}

func parseToken(strToken string) *token {
	if strToken == spaceIdentifier {
		return newToken(0, tokenTypeSpace)
	}

	if strToken == openParenIdentifier {
		return newToken(parenTypeOpen, tokenTypeParen)
	}

	if strToken == closeParenIdentifier {
		return newToken(parenTypeClose, tokenTypeParen)
	}

	if strToken == quoteIdentifier {
		return newToken(0, tokenTypeQuote)
	}

	if opType := parseOperator(strToken); opType != -1 {
		return newToken(opType, tokenTypeOp)
	}

	return newToken(strToken, tokenTypeChar)
}

func parseOperator(strToken string) opType {
	switch strToken {
	case additionIdentifier:
		return opTypeAddition
	case subtractionIdentifier:
		return opTypeSubtraction
	case multiplicationIdentifier:
		return opTypeMultiplication
	case divisionIdentifier:
		return opTypeDivision
	default:
		return -1
	}
}

func (t *token) string() string {
	switch t.t {
	case tokenTypeQuote:
		return fmt.Sprintf("%v", "\"")
	case tokenTypeSpace:
		return fmt.Sprintf("%v", " ")
	case tokenTypeOp:
		switch t.data.(opType) {
		case opTypeAddition:
			return fmt.Sprint("+")
		case opTypeMultiplication:
			return fmt.Sprint("*")
		case opTypeDivision:
			return fmt.Sprint("/")
		case opTypeSubtraction:
			return fmt.Sprint("-")
		}
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
