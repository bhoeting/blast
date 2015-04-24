package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionParsing(t *testing.T) {
	f := ParseUserFunction(NewTokenStreamFromLexer(
		Lex("function max(x = 200, y, z = .9 * (44.4 + 14))")))

	assert.Equal(t, "max", f.name)

	assert.Equal(t, "x", f.params[0].name)
	assert.Equal(t, NewNumberFromFloat(200.0), f.params[0].value.(*Number))

	assert.Equal(t, "y", f.params[1].name)
	assert.IsType(t, &tokenNil{}, f.params[1].value)

	assert.Equal(t, "z", f.params[2].name)
	assert.Equal(t, NewNumberFromFloat(52.56), f.params[2].value.(*Number))
}
