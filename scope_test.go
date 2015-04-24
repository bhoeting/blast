package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalScope(t *testing.T) {
	InitScope()

	assert.NotNil(t, GlobalScope())
	assert.NotNil(t, CurrScope())
}

func TestVariables(t *testing.T) {
	InitScope()

	SetVar("x", NewNumber("200"))
	SetVar("y", NewNumber("300"))

	x, errX := GetVar("x")
	y, errY := GetVar("y")

	assert.Nil(t, errX)
	assert.Nil(t, errY)

	assert.Equal(t, "200", x.String())
	assert.Equal(t, "300", y.String())
}
