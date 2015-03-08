package blast

import "testing"

func TestVarMapAdd(t *testing.T) {
	var testValue int64
	testValue = 3000

	vm := NewVarMap()
	vm.Add(NewVariable("test", testValue))

	_, found := vm.items["test"]

	if !found {
		t.Fatal("Failed to Add to VarMap")
	}
}

func TestVarMapGet(t *testing.T) {
	vm := NewVarMap()
	vm.Add(NewVariable("test1", int64(44)))
	vm.Add(NewVariable("test2", int64(55)))
	v, err := vm.Get("test1")

	if err != nil {
		t.Fatal(err.Error())
	}

	if result, _ := v.v.Integer(); result != int64(44) {
		t.Fatal("Failed to get correct value from VarMap")
	}
}

func TestVarMapDelete(t *testing.T) {
	vm := NewVarMap()
	vm.Add(NewVariable("test", 33))
	vm.Remove("test")

	_, err := vm.Get("test")

	if err == nil || vm.Size() != 0 {
		t.Fatal("Failed to remove item from VarMap")
	}
}

func TestVarMapSet(t *testing.T) {
	vm := NewVarMap()
	vm.Add(NewVariable("test", int64(33)))
	vm.Set(NewVariable("test", int64(34)))

	test, err := vm.Get("test")
	integer, err := test.v.Integer()

	if err != nil || integer != int64(34) {
		t.Fatal("Failed to set variable in VarMap")
	}

}
