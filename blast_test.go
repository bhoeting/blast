package blast

import (
	"strings"
	"testing"
)

func TestGenericParser(t *testing.T) {
	b := new(Blast)
	r := b.ParseGeneric("2 + 2").String()
	if r != "4" {
		t.Fatalf("Failed to parse generic code, wanted %v, got %v", "4", r)
	}

	r = b.ParseGeneric("3 * \"lol\"").String()
	if r != "lollollol" {
		t.Fatalf("Failed to parse generic code, wanted %v, got %v", "lollollol", r)
	}

	r = b.ParseGeneric("3 * 8 + 3 / 4 * \"derp\"").String()
	if r != strings.Repeat("derp", 24) {
		t.Fatalf("Failed to parse generic code, wanted %v, got %v", strings.Repeat("derp", 24), r)
	}
}

func TestPositions(t *testing.T) {
	positions := getOperatorPositions(NewTokenArray("2 + 2 * 3 - 1 / 3"))

	if len(positions) != 4 {
		t.Fatalf("Failed to get operator positions array length, expected %v, got %v, %v", 4, len(positions), positions)
	}

	if positions[0] != 3 || positions[1] != 7 || positions[2] != 1 || positions[3] != 5 {
		t.Fatal("Failed to get correct operator positions %v", positions)
	}
}

func TestVariableInitalization(t *testing.T) {
	b := NewBlast()
	b.HandleVariable("x = 2")

	t.Logf("%v", b.vMap.items)

	if _, err := b.vMap.Get("x"); err == ErrVarNotFound {
		t.Fatalf("Could not find variable %v", "x")
	}
}

func TestParseLine(t *testing.T) {
	b := NewBlast()
	b.ParseLine("x = 2 + 2")
	v, _ := b.vMap.Get("x")
	got, _ := v.v.Integer()

	if got != int64(4) {
		t.Fatalf("Could not declare variable, expected %v, got %v", 4, got)
	}

	result := b.ParseLine("x + 2")
	if result != "6" {
		t.Fatalf("Expected 6, got %v", result)
	}
}
