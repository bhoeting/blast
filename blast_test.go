package blast

import "testing"

func TestBlast(t *testing.T) {
	Init()
	Parse("x = 300")
	// assert.Equal(t, "300", Parse("x"))
}
