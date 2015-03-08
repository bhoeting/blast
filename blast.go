package blast

// Blast stores the main
// components of a
// blast program
type Blast struct {
	vMap *VarMap
}

// ProcessVariable either declares
// or initalizes a variable
func (b *Blast) ProcessVariable(v *Variable) {
	_, err := b.vMap.Add(v)

	if err == ErrVarExists {
		b.vMap.Set(v)
	}
}

// ParseGeneric parses an expression and
// returns its value as a string
func (b *Blast) ParseGeneric(code string) string {
	tokens := NewTokenArray(code)

	isInt, isFloat, isStr, isVar := false, false, false, true

	for index, token := range tokens {

	}
}
