package blast

import "errors"

type VarStack struct {
	top  *Node
	size uint64
}

type Node struct {
	value *Variable
	next  *Node
}

func NewVarStack() *VarStack {
	v := new(VarStack)
	return v
}

func (v *VarStack) Push(variable *Variable) *Variable {
	v.top = &Node{variable, v.top}
	v.size++
	return variable
}

func (v *VarStack) Pop() (variable *Variable) {
	if v.size <= 0 {
		return new(Variable)
	}

	variable = v.top.value
	v.top = v.top.next
	v.size--

	return
}

func (v *VarStack) Get(vName string) (*Variable, error) {
	n := v.top
	var index uint64

	for index = 0; index < v.size; index++ {
		if n.value.name == vName {
			return n.value, nil
		}

		n = n.next
	}

	return nil, errors.New("Could not find variable")

}

func (v *VarStack) GetByIndex(index uint64) (*Variable, error) {
	n := v.top
	var i uint64

	for i = 0; i < v.size; i++ {
		if index == i {
			return n.value, nil
		}

		n = n.next
	}

	return nil, errors.New("Could not find variable")
}

func (v *VarStack) IsEmpty() bool {
	return v.size == 0
}

func (v *VarStack) Size() int {
	return int(v.size)
}
