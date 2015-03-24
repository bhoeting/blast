package blast

import (
	"errors"
	"fmt"
)

type varType int

const (
	varTypeInt = iota
	varTypeString
	varTypeFloat
	varTypeEmpty
)

var (
	errVarNotFound = errors.New("Variable not declared.")
)

// B is the global
// blast object
var B *Blast

// Blast is a struct
// the stores basic
// language components
type Blast struct {
	vars map[string]*variable
}

// variable is a struct that
// stores a name and value
type variable struct {
	name string
	data interface{}
	t    varType
}

// NewBlast returns a new Blast struct
func NewBlast() *Blast {
	b := new(Blast)
	b.vars = make(map[string]*variable, 0)
	return b
}

// Init prepares the blast
// package for parsing
func Init() {
	B = NewBlast()
}

// Parse turns a string of code into
// language components
func Parse(code string) string {
	ts := newTokenStream(code).combine().toRPN().evaluate()
	return ts.string()
}

// addEmptyVariable adds an empty variable to
// the var map
func (b *Blast) addEmptyVariable(name string) {
	v := new(variable)
	v.t = varTypeEmpty
	b.vars[v.name] = v
}

// addVariable adds a variable to the Blast struct's
// variable map
func (b *Blast) addVariable(v *variable) *variable {
	b.vars[v.name] = v
	return v
}

// getVariable gets a variable
func (b *Blast) getVariable(name string) (*variable, error) {
	if _, ok := b.vars[name]; ok {
		return b.vars[name], nil
	}

	return new(variable), errVarNotFound
}

// removeVariable removes a variable from the
// Blast struct's variable map
func (b *Blast) removeVariable(name string) *Blast {
	delete(b.vars, name)
	return b
}

// setVariable sets a variable's value
func (b *Blast) setVariable(name string, data interface{}) *variable {
	b.vars[name] = newVariable(name, data)
	return b.vars[name]
}

// newVariable returns a new variable
func newVariable(name string, data interface{}) *variable {
	v := new(variable)
	v.data = data
	v.name = name

	switch v.data.(type) {
	case int:
		v.t = varTypeInt
	case float64:
		v.t = varTypeFloat
	case string:
		v.t = varTypeString
	default:
		v.t = varTypeString
	}

	return v
}

// string returns a string representation
// of a variable
func (v *variable) string() string {
	if v.t == varTypeString {
		return "\"" + fmt.Sprintf("%v", v.data) + "\""
	}

	return fmt.Sprintf("%v", v.data)
}

// integer returns the integer
// value of a variable
func (v *variable) integer() int {
	return v.data.(int)
}

// float returns the float
// value of a variable
func (v *variable) float() float64 {
	return v.data.(float64)
}

// str returns the string
// value of a variable
func (v *variable) str() string {
	return v.data.(string)
}

// toToken converts the variable to a token
func (v *variable) toToken() *token {
	fmt.Printf("converting %v to token %v\n", v.name, v.data)
	return newToken(v.data, varToToken[v.t])
}
