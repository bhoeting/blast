package blast

import (
	"fmt"
	"strconv"
	"unicode"
)

const (
	tokenTypeUnknown = -1
	tokenTypeOp      = 3
	tokenTypeChar    = 2
	tokenTypeParen   = 1
	tokenTypeQuote   = 4
	tokenTypeSpace   = 5
	tokenTypeString  = 6
	tokenTypeInt     = 7
	tokenTypeFloat   = 8
	tokenTypeVar     = 9
)

type opType int

const (
	opTypeAddition opType = iota
	opTypeSubtraction
	opTypeMultiplication
	opTypeDivision
)

const (
	quoteIdentifier          = "\""
	additionIdentifier       = "+"
	subtractionIdentifier    = "-"
	multiplicationIdentifier = "*"
	divisionIdentifier       = "/"
	ifIdentifier             = "if"
	spaceIdentifier          = " "
	openParenIdentifier      = "("
	closeParenIdentifier     = ")"
	decimalPointIdentifier   = "."
)

const (
	parenTypeOpen  = 1
	parenTypeClose = 2
)

type token struct {
	data interface{}
	t    int
}

type tokenStream struct {
	tokens []*token
	size   int
	index  int
}

func newTokenStream(strTokens string) *tokenStream {
	ts := new(tokenStream)
	ts.tokens = parseTokens(strTokens)
	ts.size = len(ts.tokens)
	ts.index = 0
	return ts
}

func (ts *tokenStream) next() (*token, bool) {
	if ts.index > ts.size-1 {
		return new(token), false
	}

	tok := ts.curr()
	ts.index++

	return tok, true
}

func (ts *tokenStream) combine() *tokenStream {
	replacementTokens := make([]*token, 0)
	tokensToBeCombined := make([]*token, 0)
	prevTokensEndIndex := 0

	ts.each(func(token *token, index int) {
		if token.t != tokenTypeOp || token.t != tokenTypeSpace {
			tokensToBeCombined = append(tokensToBeCombined, token)
		}

		if token.t == tokenTypeOp || index == ts.size-1 {
			replacementTokens = append(replacementTokens,
				combineTokens(tokensToBeCombined[prevTokensEndIndex:]))

			prevTokensEndIndex = len(tokensToBeCombined)
			return
		}
	})

	ts.tokens = replacementTokens
	ts.size = len(replacementTokens) - 1

	return ts
}

func combineTokens(tokens []*token) *token {
	strCombinedTok := ""
	isNum, isDecimal, isStr := true, false, false

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
		} else {
			i, _ := strconv.Atoi(strCombinedTok)
			return newToken(i, tokenTypeInt)
		}
	}

	return newToken(strCombinedTok, tokenTypeVar)
}

// each is a helper method that runs a function
// on each of the tokens in tokenStream's token slice
func (ts *tokenStream) each(f func(tok *token, index int)) {
	for next, ok := ts.next(); ok == true; next, ok = ts.next() {
		f(next, ts.index-1)
	}
}

func (ts *tokenStream) curr() *token {
	return ts.tokens[ts.index]
}

func (t *token) string() string {
	return fmt.Sprintf("%v", t.data)
}

func newToken(data interface{}, t int) *token {
	tok := new(token)
	tok.data = data
	tok.t = t
	return tok
}

func parseTokens(code string) []*token {
	tokens := make([]*token, 0)

	for _, r := range code {
		tokens = append(tokens, parseToken(string(r)))
	}

	return tokens
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
