package blast

import (
	"errors"
	"fmt"
)

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

// Value stores a piece of data
type Value struct {
	data interface{}
}

// Errors associated with type mismatching
var (
	ErrBoolTypeMismatch    = errors.New("Could not get type Bool from value.")
	ErrIntegerTypeMismatch = errors.New("Could not get type Integer from value.")
	ErrFloatTypeMismatch   = errors.New("Could not get type Float from value.")
	ErrStringTypeMismatch  = errors.New("Could not get type String from value.")
)

// NewVariable returns a new variable
func NewVariable(name string, value interface{}) *Variable {
	variable := new(Variable)
	variable.name = name
	variable.v = NewValue(value)
	return variable
}

func ParseVariable(code string) *Variable {
	return new(Variable)
}

// NewValue returns a NewValue
func NewValue(d interface{}) *Value {
	v := new(Value)
	v.data = d
	return v
}

// Bool returns the Value's data as a boolean
func (v *Value) Bool() (bool, error) {
	switch v.data.(type) {
	case bool:
		return v.data.(bool), nil
	case int:
		if v.data.(int) == 0 {
			return false, nil
		}
		return true, nil
	case string:
		if v.data.(string) == "" {
			return false, nil
		}
		return true, nil
	default:
		return false, ErrBoolTypeMismatch
	}
}

// Integer returns the Value's data as an Integer
func (v *Value) Integer() (int64, error) {
	switch v.data.(type) {
	case int64:
		return v.data.(int64), nil
	case float64:
		return int64(v.data.(float64)), nil
	case bool:
		if v.data.(bool) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrIntegerTypeMismatch
	}

}

// String returns the Value's data as a String
func (v *Value) String() string {
	switch v.data.(type) {
	case string:
		return v.data.(string)
	default:
		return fmt.Sprintf("%v", v.data)
	}
}

// Float returns the Value's data as a Float
func (v *Value) Float() (float64, error) {
	switch v.data.(type) {
	case float64:
		return v.data.(float64), nil
	case int64:
		return float64(v.data.(int64)), nil
	case bool:
		if v.data.(bool) {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0.0, ErrFloatTypeMismatch
	}
}
