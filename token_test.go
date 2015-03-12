package blast

import "testing"

func TestTokenGetter(t *testing.T) {
	tokens := getTokens("2+(2 - 9)")
	
	if tokens[0] != tokenTypeChar ||
	
}
