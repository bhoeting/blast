package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenOperations(t *testing.T) {
	ts := NewTokenStreamFromLexer(Lex("15.5 -120 true false \"derp\""))
	flt := ts.Next()
	neg := ts.Next()
	//tru := ts.Next()
	//fls := ts.Next()
	//derp := ts.Next()
	assert.Equal(t, -104.5, tokenValue(AddTokens(flt, neg)))
}

func tokenValue(token Token) interface{} {
	switch token.GetType() {
	case tokenTypeNumber, tokenTypeBoolean:
		return NumberFromToken(token)
	case tokenTypeString:
		return StringFromToken(token)
	}
	return nil
}
