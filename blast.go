package blast

// Blast stores the main
// components of a
// blast program
type Blast struct {
	vMap *VarMap
}

// NewBlast creates a new Blast
func NewBlast() *Blast {
	b := new(Blast)
	b.vMap = NewVarMap()
	return b
}

// ProcessVariable either declares
// or initalizes a variable
func (b *Blast) ProcessVariable(v *Variable) {
	_, err := b.vMap.Add(v)

	if err == ErrVarExists {
		b.vMap.Set(v)
	}
}

// ParseLine will parse any line of blast code
func (b *Blast) ParseLine(code string) string {
	return b.ParseGeneric(code)
}

// ParseGeneric parses an expression and
// returns its value as a string
func (b *Blast) ParseGeneric(code string) string {
	tokens := NewTokenArray(code)

	if len(tokens) == 1 {
		return tokens[0].String()
	}

	for index, token := range tokens {
		if op := token.Op(); op != -1 {
			return HandleTokens(tokens[index-1], tokens[index+1], token, b).String()
		}
	}

	return ""
}
