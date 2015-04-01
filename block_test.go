package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockTypes(t *testing.T) {
	b := newBlock([]string{"x = 3"})
	assert.Equal(t, blockTypeBasic, b.t)

	b = newBlock([]string{"if x == 3"})
	assert.Equal(t, blockTypeIf, b.t)
}

func TestBlock(t *testing.T) {
	strCode :=
		`x = 3

if x == 3
x = x + 6
end`

	c := newCode(strCode)
	assert.Equal(t, blockTypeBasic, (*c)[0].t)
	assert.Equal(t, blockTypeIf, (*c)[1].t)
	assert.Equal(t, blockTypeBasic, (*c)[1].blocks[0].t)
}
