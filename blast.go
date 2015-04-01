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
	vars  map[string]*variable
	funcs map[string]*function
}

// variable is a struct that
// stores a name and value
type variable struct {
	name string
	data interface{}
	t    varType
}

// variableNotDeclaredError is an error
// thrown when a variable is accessed
// but has not been declared
type variableNotDeclaredError struct {
	vName string
}

// functionNotFound is an error thrown
// when a function is accessed but
// has not been instantiated
type functionNotFoundError struct {
	fName string
}

// Error returns a string with an error message
func (err *variableNotDeclaredError) Error() string {
	return fmt.Sprintf("Variable %v not declared.", err.vName)
}

// Error returns a string with an error message
func (err *functionNotFoundError) Error() string {
	return fmt.Sprintf("Variable %v not declared.", err.fName)
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

// Parse parses code
func Parse(strCode string) {
	newCode(strCode).run()
}

// ParseBasicLine turns a line of code into
// language components
func ParseBasicLine(line string) string {
	ls := newLexemeStream(line)
	ts := newTokenStream(ls)
	return ts.parse().string()
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
	if v, ok := b.vars[name]; ok {
		return v, nil
	}

	return new(variable), &variableNotDeclaredError{name}
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

// addFunction adds a function to the funcs map
func (b *Blast) addFunction(name string, f *function) *function {
	b.funcs[name] = f
	return f
}

// getFunction gets a function
func (b *Blast) getFunction(name string) *function {
	if f, ok := b.funcs[name]; ok {
		return f, nil
	}

	return new(function), &functionNotFoundError{name}
}

// string returns a string representation
// of a Blast struct
func (b *Blast) string() string {
	str := "variables: {\n"

	for name, v := range b.vars {
		str += fmt.Sprintf("\t%v: %v\n", name, v.string())
	}

	str += "}"

	return str
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
	return newToken(v.data, 0, 0, varToToken[v.t])
}
