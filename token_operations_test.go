package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenOperations(t *testing.T) {
	ts := NewNodeStreamFromLexer(Lex("15.5 -120 true false \"derp\""))
	flt := ts.Next()
	neg := ts.Next()
	//tru := ts.Next()
	//fls := ts.Next()
	//derp := ts.Next()
	assert.Equal(t, -104.5, tokenValue(AddNodes(flt, neg)))
}

func tokenValue(token Node) interface{} {
	switch token.GetType() {
	case nodeTypeNumber, nodeTypeBoolean:
		return Float64FromNode(token)
	case nodeTypeString:
		return StringFromNode(token)
	}
	return nil
}
