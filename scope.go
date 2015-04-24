package blast

import "fmt"

var (
	Scopes            ScopeStack
	scopeIsInitalized = false
)

type Scope struct {
	vars  map[string]Token
	funcs map[string]Function
}

type ScopeStack struct {
	scopes []*Scope
	size   int
}

func InitScope() {
	Scopes.scopes = append(Scopes.scopes, NewScope())
	Scopes.size = 1
	scopeIsInitalized = true
	LoadBuiltinFunctions()
}

func GlobalScope() *Scope {
	return Scopes.scopes[0]
}

func CurrScope() *Scope {
	return Scopes.Top()
}

func SetVar(name string, token Token) {
	Scopes.scopes[Scopes.size-1].SetVar(name, token)
}

func GetVar(name string) (Token, error) {
	return Scopes.scopes[Scopes.size-1].GetVar(name)
}

func GetFunc(name string) (Function, error) {
	return Scopes.scopes[0].GetFunc(name)
}

func SetFunc(name string, f Function) {
	Scopes.scopes[0].SetFunc(name, f)
}

type ErrVarNotFound struct {
	vName string
}

type ErrFuncNotFound struct {
	fName string
}

func (err *ErrVarNotFound) Error() string {
	return fmt.Sprintf("Variable %v not found", err.vName)
}

func (err *ErrFuncNotFound) Error() string {
	return fmt.Sprintf("Function %v not found", err.fName)
}

func NewScope() *Scope {
	s := new(Scope)
	s.vars = make(map[string]Token)
	s.funcs = make(map[string]Function)
	return s
}

func (s *Scope) SetVar(name string, value Token) {
	s.vars[name] = value
}

func (s *Scope) GetVar(name string) (Token, error) {
	if v, ok := s.vars[name]; ok {
		return v, nil
	}

	return &tokenNil{}, &ErrVarNotFound{name}
}

func (s *Scope) GetFunc(name string) (Function, error) {
	if f, ok := s.funcs[name]; ok {
		return f, nil
	}

	return funcNil, &ErrFuncNotFound{name}
}

func (s *Scope) SetFunc(name string, f Function) {
	s.funcs[name] = f
}

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

func (st *ScopeStack) Top() *Scope {
	return st.scopes[st.size-1]
}

func (st *ScopeStack) New() *Scope {
	s := new(Scope)
	s.vars = st.Top().vars
	s.funcs = st.Top().funcs

	st.size++
	st.scopes = append(st.scopes, s)
	return s
}

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
