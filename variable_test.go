package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariableMethods(t *testing.T) {
	b := NewBlast()
	b.addVariable(newVariable("x", 200))

	v, err := b.getVariable("x")
	if err == errVarNotFound {
		t.Fatal(err.Error())
	}
	assert.Equal(t, 200, v.integer())

	b.setVariable("x", 201)
	v, err = b.getVariable("x")
	assert.Equal(t, 201, v.integer())

	b.removeVariable("x")
	_, err = b.getVariable("x")
	assert.Equal(t, errVarNotFound, err)

	b.addEmptyVariable("d")
	b.setVariable("d", 5000)
	v, err = b.getVariable("d")
	assert.Equal(t, 5000, v.data)
}

func TestVariableDeclaration(t *testing.T) {
	Init()
	Parse("x = 3")
	v, err := B.getVariable("x")

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, 3, v.integer())

	Parse("y=3000+3")
	v, err = B.getVariable("y")

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, 3003, v.integer())

	Parse("y=300+\"string\"")
	v, err = B.getVariable("y")

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "300string", v.str())
}

func TestVariableUsage(t *testing.T) {
	Init()
	Parse("n=3")
	Parse("x=4")
	Parse("s=\"string\"")

	assert.Equal(t, "3", Parse("n"))
	assert.Equal(t, "5", Parse("2+n"))
	assert.Equal(t, "100", Parse("n=n+97"))
	assert.Equal(t, "\"100string\"", Parse("n+s"))
	assert.Equal(t, "\"stringstringstringstring\"", Parse("x*s"))
}
