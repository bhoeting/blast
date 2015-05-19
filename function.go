package blast

import "fmt"

// Function is an interface
// with a call method
type Function interface {
	Call(args *NodeStream) Node
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
	value Node
}

// goFunc is a func type with a NodeStream
// parameter that returns an interface
type goFunc func(args *NodeStream) interface{}

// BuilinFunction is a struct represeting
// a function built in to Blast that
// is implemented in Go
type BuiltinFunction struct {
	f goFunc
}

// Call runs a UserFunction and returns the
// result as a odSe
func (f *UserFunction) Call(args *NodeStream) Node {
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

// Call runs a BuiltinFunction and returns the result as a Node
func (bf *BuiltinFunction) Call(args *NodeStream) Node {
	result := bf.f(args)

	if result == nil {
		return &nodeNil{}
	}

	switch result.(type) {
	case Node:
		return result.(Node)
	case float64:
		return NewNumberFromFloat(result.(float64))
	case string:
		return NewString(result.(string))
	case bool:
		return NewBooleanFromBool(result.(bool))
	}

	return &nodeNil{}
}

// ParseUserFunction parses a NodeStream into a user
// function definition
func ParseUserFunction(ns *NodeStream) *UserFunction {
	parenDepth := 1
	f := new(UserFunction)
	paramns := NewNodeStream()

	// Skip the "function"
	ns.Next()

	// Set the name
	f.name = ns.Next().(*FunctionCall).name

	// Skip the first paren
	ns.Next()

	for ns.HasNext() {
		node := ns.Next()
		if node.GetType() == nodeTypeParen {
			switch paren := node.(*Paren); paren.typ {
			case parenTypeOpen:
				parenDepth++
			case parenTypeClose:
				parenDepth--
			}
		}

		if parenDepth == 0 {
			if paramns.Length() > 0 {
				f.params = append(f.params, ParseParam(paramns))
			}
			break
		}

		if node.GetType() == nodeTypeComma {
			f.params = append(f.params, ParseParam(paramns))
			paramns = NewNodeStream()
		} else {
			paramns.Push(node)
		}
	}

	return f
}

// ParseParam parses a NodeStream into
// a parameter
func ParseParam(ns *NodeStream) *Param {
	param := new(Param)

	if ns.Length() == 1 {
		return &Param{
			name:  ns.Next().String(),
			value: &nodeNil{},
		}
	}

	param.name = ns.Next().String()

	if ns.Peek().GetType() == nodeTypeOperator &&
		ns.Peek().(*Operator).typ == opTypeAssignment {
		ns.Next()
	}

	param.value = ns.Chop().Evaluate()
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
}

// builtinPrint prinns the Nodes
func builtinPrint(args *NodeStream) interface{} {
	str := ""

	for _, node := range args.nodes {
		str += StringFromNode(node) + " "
	}

	str = str[:len(str)-1]
	fmt.Print(str)

	return nil
}

// builtinPrint prinns the Nodes on their own line
func builtinPrintln(args *NodeStream) interface{} {
	str := ""

	for _, node := range args.nodes {
		str += StringFromNode(node) + " "
	}

	if len(str) > 0 {
		str = str[:len(str)-1]
	}

	fmt.Println(str)
	return nil
}

// builtinModulus performs the modulus operation on two Nodes.
// TODO: make this an operator
func builtinModulus(args *NodeStream) interface{} {
	return float64(int(Float64FromNode(args.Next())) % int(Float64FromNode(args.Next())))
}
