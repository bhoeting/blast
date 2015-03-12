package blast

// chunk is an interface
// that represents most
// small pieces of code
// for example, the chunks
// from the string x + 2 * 300
// would be x, +, 2, *, 300
type chunk struct {
}

type value struct {
}

// variable is meant to be
// an anonymous field in
// variable type structs
type variable struct {
	name string
}

// str represents
// a string variable
type str struct {
	chunk
	data string
}

// integer represents
// an integer variable
type integer struct {
	chunk
}

// newStr returns a new string
func newStr(data string) *str {
	s := new(str)
	s.data = data
	return s
}
