package blast

// line stores a slice of chunks
type line struct {
	chunks []*chunk
}

// addChunk appends a chunk to the
// line's chunk slice
func (ln *line) addChunk(c *chunk) {
	ln.chunks = append(ln.chunks, c)
}

// parse turns a string of
// code into a line struct
func parse(code string) {
	tokens := parseToken(code)
	isStr := false
	strChunk := ""
	ln := new(line)
	plevel := 0

	for i, t := range tokens {
		if t.t == tokenTypeQuote {
			if isStr {
				isStr = false
				ln.addChunk(newStr(strChunk))
				strChunk = ""
				continue
			}

			isStr = true
		}

		if isStr {
			strChunk += t.str()
			continue
		}

		if t.t == tokenTypeOp {
			ln.addChunk(newOperator(t.data))
		}

	}
}

func getIntFromStr(str string) (int64, error) {

}
