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
	positions := getOperatorPositions(tokens)
	lenTokens := len(tokens)
	lenPositions := len(positions)

	for lenPositions > 0 {
		pos := positions[0]
		tokens[pos-1] = HandleTokens(tokens[pos-1], tokens[pos+1], tokens[pos], b)

		for tIndex := pos; tIndex < lenTokens-2; tIndex++ {
			tokens[tIndex] = tokens[tIndex+2]
		}

		tokens[lenTokens-1] = new(Token)
		tokens[lenTokens-2] = new(Token)
		tokens = tokens[:lenTokens-2]
		lenTokens -= 2

		copy(positions[0:], positions[1:])
		positions[len(positions)-1] = 0
		positions = positions[:len(positions)-1]
		lenPositions -= 1

		for pIndex := 0; pIndex < lenPositions; pIndex++ {
			if positions[pIndex] > pos {
				positions[pIndex] -= 2
			}
		}
	}

	return tokens[0].String()
}

// getOperatorPositions returns an array of
// operator indicies sorted by the order
// in which their operands are to be evaluated.
func getOperatorPositions(tokens []*Token) []int {
	multAndDiv, addAndSub, strs :=
		make([]int, 0), make([]int, 0), make([]int, 0)

	// Loop through each token, but only act on operators.
	for index, token := range tokens {
		if token.Type() == TokenTypeOperator {
			// If either of the operator's operands is a
			// string, add the index to the strs array.
			if tokens[index-1].Type() == TokenTypeString ||
				tokens[index+1].Type() == TokenTypeString {
				strs = append(strs, index)
				continue
			}

			// If the operator is addition/subtraction,
			// send it to the addAndSub slice.
			if token.Op() == OpTypeAddition ||
				token.Op() == OpTypeSubtraction {
				addAndSub = append(addAndSub, index)
			} else {
				// Otherwise, it must be multiplication/division.
				// send it to the multAndDiv slice.
				multAndDiv = append(multAndDiv, index)
			}
		}
	}

	// Combine the slices in the order of their precedence.
	return append(append(multAndDiv, addAndSub...), strs...)
}
