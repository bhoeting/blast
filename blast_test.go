package blast

import "testing"

func TestGenericParser(t *testing.T) {
	b := new(Blast)
	r := b.ParseGeneric("2 + 2")
	if r != "4" {
		t.Fatalf("Failed to parse generic code, wanted %v, got %v", "4", r)
	}

	r = b.ParseGeneric("3 * \"lol\"")
	if r != "lollollol" {
		t.Fatalf("Failed to parse generic code, wanted %v, got %v", "lollollol", r)
	}
}
