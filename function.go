package blast

import "fmt"

// Function is an interface
// with a call method
type Function interface {
	Call(args *TokenStream) Token
}

// funcNil is returned when there
// is undefined behavior related
// to a Function
var funcNil = &UserFunction{}

// UserFunction is a struct
// that represents a user
// defined function
type UserFunction struct {
	params []*Param
	name   string
	block  *Block
}

// Param is a struct that
// represents a function
// definition parameter
type Param struct {
	name  string
	value Token
}

// goFunc is a func type with a TokenStream
// parameter that returns an interface
type goFunc func(args *TokenStream) interface{}

// BuilinFunction is a struct represeting
// a function built in to Blast that
// is implemented in Go
type BuiltinFunction struct {
	f goFunc
}

// Call runs a UserFunction and returns the
// result as a Token
func (f *UserFunction) Call(args *TokenStream) Token {
	Scopes.New()

	for _, param := range f.params {
		if args.HasNext() {
			arg := args.Next()
			SetVar(param.name, arg)
		} else {
			SetVar(param.name, param.value)
		}
	}

	result, _ := f.block.RunBlocks()
	Scopes.Pop()
	return result
}

// Call runs a BuiltinFunction and returns the result as a Token
func (bf *BuiltinFunction) Call(args *TokenStream) Token {
	result := bf.f(args)

	if result == nil {
		return &tokenNil{}
	}

	switch result.(type) {
	case Token:
		return result.(Token)
	case float64:
		return NewNumberFromFloat(result.(float64))
	case string:
		return NewString(result.(string))
	case bool:
		return NewBooleanFromBool(result.(bool))
	}

	return &tokenNil{}
}

// ParseUserFunction parses a TokenStream into a user
// function definition
func ParseUserFunction(ts *TokenStream) *UserFunction {
	parenDepth := 1
	f := new(UserFunction)
	paramTS := NewTokenStream()

	// Skip the "function"
	ts.Next()

	// Set the name
	f.name = ts.Next().(*FunctionCall).name

	// Skip the first paren
	ts.Next()

	for ts.HasNext() {
		token := ts.Next()
		if token.GetType() == tokenTypeParen {
			switch paren := token.(*Paren); paren.typ {
			case parenTypeOpen:
				parenDepth++
			case parenTypeClose:
				parenDepth--
			}
		}

		if parenDepth == 0 {
			if paramTS.Length() > 0 {
				f.params = append(f.params, ParseParam(paramTS))
			}
			break
		}

		if token.GetType() == tokenTypeComma {
			f.params = append(f.params, ParseParam(paramTS))
			paramTS = NewTokenStream()
		} else {
			paramTS.Push(token)
		}
	}

	return f
}

// ParseParam parses a TokenStream into
// a parameter
func ParseParam(ts *TokenStream) *Param {
	param := new(Param)

	if ts.Length() == 1 {
		return &Param{
			name:  ts.Next().String(),
			value: &tokenNil{},
		}
	}

	param.name = ts.Next().String()

	if ts.Peek().GetType() == tokenTypeOperator &&
		ts.Peek().(*Operator).typ == opTypeAssignment {
		ts.Next()
	}

	param.value = ts.Chop().Evaluate()
	return param
}

// NewBuiltinFunction returns a new BuiltinFunction
func NewBuiltinFunc(f goFunc) *BuiltinFunction {
	return &BuiltinFunction{
		f: f,
	}
}

// LoadBuiltinFunctions adds all the BuiltinFuctions
// to the scope
func LoadBuiltinFunctions() {
	SetFunc("print", NewBuiltinFunc(builtinPrint))
	SetFunc("println", NewBuiltinFunc(builtinPrintln))
	SetFunc("modulus", NewBuiltinFunc(builtinModulus))
}

// builtinPrint prints the Nodes
func builtinPrint(args *TokenStream) interface{} {
	str := ""

	for _, token := range args.tokens {
		str += StringFromToken(token) + " "
	}

	str = str[:len(str)-1]
	fmt.Print(str)

	return nil
}

// builtinPrint prints the Nodes on their own line
func builtinPrintln(args *TokenStream) interface{} {
	str := ""

	for _, token := range args.tokens {
		str += StringFromToken(token) + " "
	}

	if len(str) > 0 {
		str = str[:len(str)-1]
	}

	fmt.Println(str)
	return nil
}

// builtinModulus performs the modulus operation on two Nodes.
// TODO: make this an operator
func builtinModulus(args *TokenStream) interface{} {
	return float64(int(NumberFromToken(args.Next())) % int(NumberFromToken(args.Next())))
}
