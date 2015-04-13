package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariableMethods(t *testing.T) {
	b := NewBlast()
	b.addVariable(newVariable("x", 200))

	v, err := b.getVariable("x")
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, 200, v.integer())

	b.setVariable("x", 201)
	v, err = b.getVariable("x")
	assert.Equal(t, 201, v.integer())

	b.removeVariable("x")
	_, err = b.getVariable("x")
	assert.IsType(t, new(variableNotDeclaredError), err)

	b.addEmptyVariable("d")
	b.setVariable("d", 5000)
	v, err = b.getVariable("d")
	assert.Equal(t, 5000, v.data)
}

func TestVariableDeclaration(t *testing.T) {
	Init()
	ParseBasicLine("x = 3")
	v, err := B.getVariable("x")

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, 3, v.integer())

	ParseBasicLine("y=3000+3")
	v, err = B.getVariable("y")

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, 3003, v.integer())

	ParseBasicLine("y=300+\"string\"")
	v, err = B.getVariable("y")

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "300string", v.str())
}

func TestVariableUsage(t *testing.T) {
	Init()
	ParseBasicLine("n=3")
	ParseBasicLine("x=4")
	ParseBasicLine("s=\"string\"")

	assert.Equal(t, "3", ParseBasicLine("n"))
	assert.Equal(t, "5", ParseBasicLine("2+n"))
	assert.Equal(t, "100", ParseBasicLine("n=n+97"))
	assert.Equal(t, "\"100string\"", ParseBasicLine("n+s"))
	assert.Equal(t, "\"stringstringstringstring\"", ParseBasicLine("x*s"))
}
