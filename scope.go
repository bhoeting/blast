package blast

import "fmt"

var (
	// Global scope stack
	Scopes            ScopeStack
	scopeIsInitalized = false
)

// Scope stores a map of Nodes
// and a map of Functions
type Scope struct {
	vars  map[string]Token
	funcs map[string]Function
}

// ScopeStack is a stack
// of Scopes
type ScopeStack struct {
	scopes []*Scope
	size   int
}

// InitScope initalizes the global ScopeStack
// and loads the builtin functions
func InitScope() {
	Scopes.scopes = append(Scopes.scopes, NewScope())
	Scopes.size = 1
	scopeIsInitalized = true
	LoadBuiltinFunctions()
}

// GlobalScope returns the scope
// at the bottom of the scope stack
func GlobalScope() *Scope {
	return Scopes.scopes[0]
}

// GlobalScope returns the scope
// at the top of the scope stack
func CurrScope() *Scope {
	return Scopes.Top()
}

// SetVar sets a variable on the current scope
func SetVar(name string, token Token) {
	Scopes.scopes[Scopes.size-1].SetVar(name, token)
}

// GetVar gets a variable from the current scope
func GetVar(name string) (Token, error) {
	return Scopes.scopes[Scopes.size-1].GetVar(name)
}

// GetFunc returns a function from the current scope
func GetFunc(name string) (Function, error) {
	return Scopes.scopes[0].GetFunc(name)
}

// SetFunc sets a function on the current scope
func SetFunc(name string, f Function) {
	Scopes.scopes[0].SetFunc(name, f)
}

// ErrVarNotFound is thrown when a
// variable that doesn't exist
// is accessed
type ErrVarNotFound struct {
	vName string
}

// ErrFuncNotFound is thrown when a
// function that doens't exist
// is accessed
type ErrFuncNotFound struct {
	fName string
}

// Error returns a string representation of an
// ErrVarNotFound
func (err *ErrVarNotFound) Error() string {
	return fmt.Sprintf("Variable %v not found", err.vName)
}

// Error returns a string representation of an
// ErrFuncNotFound
func (err *ErrFuncNotFound) Error() string {
	return fmt.Sprintf("Function %v not found", err.fName)
}

// NewScope returns a new Scope
func NewScope() *Scope {
	s := new(Scope)
	s.vars = make(map[string]Token)
	s.funcs = make(map[string]Function)
	return s
}

// SetVar sets a variable on the Scope
func (s *Scope) SetVar(name string, value Token) {
	s.vars[name] = value
}

// GetVar gets a variable from the Scope
func (s *Scope) GetVar(name string) (Token, error) {
	if v, ok := s.vars[name]; ok {
		return v, nil
	}

	return &tokenNil{}, &ErrVarNotFound{name}
}

// GetFunc returns the function from the Scope
func (s *Scope) GetFunc(name string) (Function, error) {
	if f, ok := s.funcs[name]; ok {
		return f, nil
	}

	return funcNil, &ErrFuncNotFound{name}
}

// SetFunc sets the a function on the Scope
func (s *Scope) SetFunc(name string, f Function) {
	s.funcs[name] = f
}

// String returns a string representation
// of the Scope
func (s *Scope) String() string {
	str := "vars: {\n"

	for name, v := range s.vars {
		str += "\t" + name + ": " + v.String() + "\n"
	}

	str += "}\n\nfuncs: {\n"

	for name, _ := range s.funcs {
		str += "\t" + name + "\n"
	}

	return str + "}"
}

// Top returns the Scope at the top of the
// ScopeStack
func (st *ScopeStack) Top() *Scope {
	return st.scopes[st.size-1]
}

// New adds a new Scope the to the scope stack
func (st *ScopeStack) New() *Scope {
	s := new(Scope)
	s.vars = st.Top().vars
	s.funcs = st.Top().funcs

	st.size++
	st.scopes = append(st.scopes, s)
	return s
}

// Pop removes the Scope at the top
// of the stack
func (st *ScopeStack) Pop() *Scope {
	top := st.Top()
	st.size--
	curr := st.Top()

	// Move any variables that could
	// have been modified to the
	// new current Scope
	for name, token := range top.vars {
		if _, ok := curr.vars[name]; ok {
			curr.vars[name] = token
		}
	}

	st.scopes = st.scopes[:st.size]
	return top
}
