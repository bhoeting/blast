package blast

import "errors"

type VarMap struct {
	items map[string]*Variable
	size  uint
}

// Error to be returned when the requested item could not be found
var ErrVarNotFound = errors.New("Could not find variable in VarMap.")
var ErrVarExists = errors.New("Variable already exists.")

// NewVarMap returns a new VarMap
func NewVarMap() *VarMap {
	v := new(VarMap)
	v.size = 0
	v.items = make(map[string]*Variable, 0)
	return v
}

// Add adds an item to the VarMap
func (vm *VarMap) Add(v *Variable) (*Variable, error) {
	if !vm.Check(v.name) {
		vm.items[v.name] = v
		vm.size++
		return v, nil
	}

	return new(Variable), ErrVarExists
}

// Remove removes an item from the VarMap
func (vm *VarMap) Remove(n string) *Variable {
	if vm.Check(n) {
		item := vm.items[n]
		delete(vm.items, n)
		vm.size--
		return item
	}

	return new(Variable)
}

// Set sets an item to the specified Variable
func (vm *VarMap) Set(v *Variable) *Variable {
	vm.items[v.name] = v
	return v
}

// Get gets an item from the VarMap
func (vm *VarMap) Get(n string) (*Variable, error) {
	if item, found := vm.items[n]; found {
		return item, nil
	}

	return new(Variable), ErrVarNotFound
}

// Check checks if there is an item in the
// VarMap with the specified key
func (vm *VarMap) Check(n string) bool {
	_, found := vm.items[n]
	return found
}

// Size returns the VarMap's size
func (vm *VarMap) Size() uint {
	return vm.size
}
