package blast

type TokenStream struct {
	pos    int
	size   int
	tokens []Token
}

func (ts *TokenStream) Push(t Token) Token {
	ts.tokens = append(ts.tokens, t)
	ts.size++
	return t
}

func (ts *TokenStream) Pop() Token {
	ts.size--
	t := ts.tokens[ts.size]
	ts.tokens = ts.tokens[:ts.size]
	return t
}

func (ts *TokenStream) Top() Token {
	if ts.size-1 < 0 {
		return &tokenNil{}
	}

	t := ts.tokens[ts.size-1]
	return t
}

func (ts *TokenStream) RemoveLast() Token {
	ts.size--
	t := ts.tokens[0]
	ts.tokens = ts.tokens[1:]
	return t
}

func (ts *TokenStream) Length() int {
	return ts.size
}

func (ts *TokenStream) HasNext() bool {
	return ts.pos < ts.size
}

func (ts *TokenStream) Next() Token {
	if ts.pos >= ts.size {
		return &tokenNil{}
	}

	token := ts.tokens[ts.pos]
	ts.pos++
	return token
}

func (ts *TokenStream) Backup() *TokenStream {
	ts.pos--
	return ts
}

func (ts *TokenStream) Peek() Token {
	t := ts.Next()
	ts.Backup()
	return t
}

func (ts *TokenStream) Evaluate() Token {
	rpn := NewTokenStreamInRPN(ts)
	return EvaluateRPN(rpn)
}

func (ts *TokenStream) String() string {
	str := ""

	for _, token := range ts.tokens {
		str += token.String() + " "
	}

	return str
}

func NewTokenStream() *TokenStream {
	ts := new(TokenStream)
	return ts
}

func (ts *TokenStream) Chop() *TokenStream {
	newTS := NewTokenStream()

	for ts.HasNext() {
		newTS.Push(ts.Next())
	}

	return newTS
}

func (ts *TokenStream) Reverse() {
	reversed := NewTokenStream()

	for len(ts.tokens) != 0 {
		token := ts.Pop()
		reversed.Push(token)
	}

	ts.tokens = reversed.tokens
	ts.size = reversed.size
}

func (ts *TokenStream) Reset() {
	ts.pos = 0
}

func NewTokenStreamFromLexer(l *Lexer) *TokenStream {
	ts := NewTokenStream()

	for l.HasNextItem() {
		switch item := l.NextItem(); item.typ {
		case itemTypeNum:
			ts.Push(NewNumber(item.text))
		case itemTypeBool:
			ts.Push(NewBoolean(item.text))
		case itemTypeString:
			ts.Push(NewString(item.text))
		case itemTypeOperator:
			ts.Push(NewOperator(item.text))
		case itemTypeOpenParen, itemTypeCloseParen:
			ts.Push(NewParen(item.text))
		case itemTypeComma:
			ts.Push(NewComma())
		case itemTypeIdentifier:
			if l.HasNextItem() && l.PeekItem().typ == itemTypeOpenParen {
				ts.Push(NewFunctionCall(item.text))
			} else {
				ts.Push(NewVariable(item.text))
			}
		default:
			ts.Push(NewReserved())
		}
	}

	l.itemPos = 0
	return ts
}
