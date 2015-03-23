package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentFromToken(t *testing.T) {
	pToken := newToken(1, tokenTypeParen)
	iToken := newToken(23, tokenTypeInt)
	fToken := newToken(23.3, tokenTypeFloat)
	sToken := newToken("string", tokenTypeString)
	oToken := newToken(opType(1), tokenTypeOp)

	pComp := componentFromToken(pToken)
	iComp := componentFromToken(iToken)
	fComp := componentFromToken(fToken)
	sComp := componentFromToken(sToken)
	oComp := componentFromToken(oToken)

	assert.IsType(t, new(paren), pComp)
	assert.IsType(t, new(integer), iComp)
	assert.IsType(t, new(float), fComp)
	assert.IsType(t, new(str), sComp)
	assert.IsType(t, new(operator), oComp)
}

func TestComponentStreamFromTokenSream(t *testing.T) {
	ts := newTokenStream("230+(221.3-900) + \"string\"").combine()
	cs := newComponentStreamFromTokenStream(ts)

	assert.IsType(t, new(integer), cs.items[0])
	assert.IsType(t, new(operator), cs.items[1])
	assert.IsType(t, new(paren), cs.items[2])
	assert.IsType(t, new(float), cs.items[3])
	assert.IsType(t, new(operator), cs.items[4])
	assert.IsType(t, new(integer), cs.items[5])
	assert.IsType(t, new(paren), cs.items[6])
	assert.IsType(t, new(operator), cs.items[7])
	assert.IsType(t, new(str), cs.items[8])
}

func BenchmarkComponentStreamFromTokenStream(b *testing.B) {

}
