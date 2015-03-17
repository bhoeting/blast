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
	ts := newTokenStream("(23+\"string\")+300.2").combine()

	t.Log(ts.string())

	assert.Equal(t, tokenTypeParen, ts.tokens[0].t)
	assert.Equal(t, tokenTypeInt, ts.tokens[1].t)
	assert.Equal(t, tokenTypeOp, ts.tokens[2].t)
	assert.Equal(t, tokenTypeString, ts.tokens[3].t)
	assert.Equal(t, tokenTypeParen, ts.tokens[4].t)
	assert.Equal(t, tokenTypeOp, ts.tokens[5].t)
	assert.Equal(t, tokenTypeFloat, ts.tokens[6].t)
}
