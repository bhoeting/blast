package blast

import "strings"

const (
	VarTypeBool    = 0
	VarTypeInteger = 1
	VarTypeString  = 2
	VarTypeFloat   = 3
)

// Variable stores a
// Type and a Value
type Variable struct {
	t    int
	v    *Value
	name string
}

// NewVariable returns a new variable
func NewVariable(name string, value interface{}) *Variable {
	variable := new(Variable)
	variable.name = name
	variable.v = NewValue(value)
	return variable
}

// ParseVariableAssignment assigns a variable
// to a value or creates a new variable based
// an assignment statement
func ParseVariableAssignment(input string) {
	assignmentIndex := strings.Index(input, string(AssignmentOp))
	if assignmentIndex == -1 {
		return
	}

	return
}
