package blast

// operator represents
// an operator
type operator struct {
	t int
}

// newOperator returns a new operator
func newOperator(t int) *operator {
	o := new(operator)
	o.t = t
	return o
}
