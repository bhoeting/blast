package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineTypes(t *testing.T) {
	basicLine, _ := newLine("twentySeven = 3 * 9")
	functionLine, _ := newLine("function derp(x, y)")
	ifLine, _ := newLine("if variable == 4 * 7")

	assert.Equal(t, lineTypeBasic, basicLine.t)
	assert.Equal(t, lineTypeFunction, functionLine.t)
	assert.Equal(t, lineTypeIf, ifLine.t)
}
