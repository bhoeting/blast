package blast

import "testing"

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
	ts := newTokenStream("23+\"string\"+300.2").combine()

	firstToken := ts.tokens[0].data.(int)
	secondToken := ts.tokens[1].data.(string)
	thirdToken := ts.tokens[2].data.(float64)

	if firstToken != 23 {
		t.Fatalf("Could not combine tokens, wanted %v got %v", 21, firstToken)
	}

	if secondToken != "string" {
		t.Fatalf("Could not combine tokens, wanted %v got %v", "string", secondToken)
	}

	if thirdToken != 300.2 {
		t.Fatalf("Could not combine tokens, wanted %v got %v", 300.2, thirdToken)
	}
}
