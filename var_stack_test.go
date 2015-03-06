package blast

import "testing"

func TestVarStack(t *testing.T) {
	v := NewVarStack()
	v.Push(NewVariable("x", 3))
	v.Push(NewVariable("y", "busta"))
	v.Push(NewVariable("z", 6.6))

	item, err := v.Get("y")

	actual := item.v.String()
	expected := "busta"

	if err != nil || actual != expected {
		t.Fatalf("Failed to get variable '%v' from stack, got %v", expected, actual)
	}

	v.Pop()
	_, err = v.Get("z")

	if err == nil {
		t.Fatal("Failed to pop from the VarStack")
	}

}
