package blast

import "testing"

func TestTokenGetter(t *testing.T) {
	tokens := parseTokens("2+(2-9) + \"")

	tokenTypes := []int{
		tokenTypeChar, tokenTypeOp,
		tokenTypeParen, tokenTypeChar,
		tokenTypeOp, tokenTypeChar,
		tokenTypeParen, tokenTypeSpace,
		tokenTypeOp, tokenTypeSpace,
		tokenTypeQuote,
	}

	for i, _ := range tokenTypes {
		if tokenTypes[i] != tokens[i].t {
			t.Fatalf(
				`Failed to get correct token types, 
				 got %v, expected %v, in index %v`,
				tokens[i].t, tokenTypes[i], i,
			)
		}
	}
}
