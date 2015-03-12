package blast

const (
	tokenTypeOp    = 3
	tokenTypeChar  = 2
	tokenTypeParen = 1
	tokenTypeQuote = 4
	tokenTypeSpace = 5
)

const (
	opTypeAddition       = 1
	opTypeSubtraction    = 2
	opTypeMultiplication = 3
	opTypeDivision       = 4
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
)

const (
	parenTypeOpen  = 1
	parenTypeClose = 2
)

type token struct {
	data interface{}
	t    int
}

func newToken(data interface{}, t int) {
	tok := new(token)
	tok.data = data
	tok.t = t
	return t
}

func getTokens(code string) []*Tokens {
	tokens := make([]*token, len(code))

	for i, r := range code {
		tokens = append(tokens, parseToken(strToken))
	}
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

func parseOperator(strToken string) int {
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
