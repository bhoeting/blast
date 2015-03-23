package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testTokenStream tests that the correct
// token type is created for each token
// in the test string
func TestTokenStream(t *testing.T) {
	ts := newTokenStream("2+(2-9) + \"")

	// The desired token types
	tokenTypes := []int{
		tokenTypeChar, tokenTypeOp,
		tokenTypeParen, tokenTypeChar,
		tokenTypeOp, tokenTypeChar,
		tokenTypeParen, tokenTypeSpace,
		tokenTypeOp, tokenTypeSpace,
		tokenTypeQuote,
	}

	// Check that each token matches the
	// desired token types
	ts.each(func(tok *token, index int) {
		if tokenTypes[index] != tok.t {
			t.Fatalf("Could not get correct token types, index: %v want: %v got: %v",
				ts.index, tokenTypes[index], tok.t)
		}
	})
}

func TestCombiningOfTokens(t *testing.T) {
	ts := newTokenStream("(23+\"string\") + 300.2 + \"derp\"").combine()

	assert.Equal(t, tokenTypeParen, ts.tokens[0].t)
	assert.Equal(t, tokenTypeInt, ts.tokens[1].t)
	assert.Equal(t, tokenTypeOp, ts.tokens[2].t)
	assert.Equal(t, tokenTypeString, ts.tokens[3].t)
	assert.Equal(t, tokenTypeParen, ts.tokens[4].t)
	assert.Equal(t, tokenTypeOp, ts.tokens[5].t)
	assert.Equal(t, tokenTypeFloat, ts.tokens[6].t)
	assert.Equal(t, tokenTypeOp, ts.tokens[7].t)
	assert.Equal(t, tokenTypeString, ts.tokens[8].t)
}

func TestTokenOperations(t *testing.T) {
	add := newTokenStream("222+4").combine()
	assert.Equal(t, 226, addTokens(add.tokens[0], add.tokens[2]).integer())

	add = newTokenStream("200+2.2").combine()
	assert.Equal(t, 202.2, addTokens(add.tokens[0], add.tokens[2]).float())

	add = newTokenStream("200+\"string\"").combine()
	assert.Equal(t, "200string", addTokens(add.tokens[0], add.tokens[2]).str())

	mult := newTokenStream("7*8").combine()
	assert.Equal(t, 56, multiplyTokens(mult.tokens[0], mult.tokens[2]).integer())

	mult = newTokenStream("7.7*8.99").combine()
	assert.Equal(t, 69.223, multiplyTokens(mult.tokens[0], mult.tokens[2]).float())

	mult = newTokenStream("3*\"string\"").combine()
	assert.Equal(t, "stringstringstring", multiplyTokens(mult.tokens[0], mult.tokens[2]).str())
}
