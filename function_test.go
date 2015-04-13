package blast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionParsing(t *testing.T) {
	ts := newTokenStream(newLexemeStream("function print(x = 3 + 5, y, z = 4)"))
	f := parseUserFunction(ts)

	assert.Equal(t, "print", f.name)
	assert.Equal(t, 8, f.params[0].token.integer())
	assert.Equal(t, "x", f.params[0].name)
	assert.Equal(t, "y", f.params[1].name)
	assert.Nil(t, f.params[1].token.data)
	assert.Equal(t, "z", f.params[2].name)

	ts = newTokenStream(newLexemeStream("function add(n1, n2)"))
	f = parseUserFunction(ts)

	assert.Equal(t, "n1", f.params[0].name)
	assert.Equal(t, "n2", f.params[1].name)
}

func TestFunctionDeclaration(t *testing.T) {
	_ = newLines(`function add(n1, n2)
		if n1 == n2
		return 2 * n1
		end

		return n1 + n2
		end
		x = 3`)

	f, _ := B.getFunction("add")
	uf, _ := f.(*userFunction)

	expected := newTokenStream(newLexemeStream("if n1 == n2 return 2 * n1 end return n1 + n2"))

	assert.Equal(t, expected.string(), uf.lines.string())
}

// func TestFunctionCalling(t *testing.T) {
// 	lines := newLines(`function add(n1, n2)
// 		if n1 == n2
// 		return 2 * n1
// 		end

// 		return n1 + n2
// 		end`)

// 	// lines.run()
// 	println(lines.string())

// 	five := newLines("add(2, 3)").run()
// 	four := newLines("add(2, 2)").run()

// 	assert.Equal(t, 5, five)
// 	assert.Equal(t, 4, four)
// }
