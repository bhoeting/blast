package blast

import (
	"log"
	"strings"
)

// LineReader is a struct that
// assists with reading lines
type LineReader struct {
	strLines []string
	lines    []*Line
	size     int
	pos      int
	nLines   int
	lineNum  int
}

var (
	// lineEOF is used when
	// no more lines are
	// available
	lineEOF = &Line{
		typ: lineTypeEOF,
	}

	// itemLineKey is used to determine
	// a line type based on the first
	// Token
	tokenLineKey = map[itemType]lineType{
		itemTypeEnd:       lineTypeEnd,
		tokenTypeIf:       lineTypeIf,
		tokenTypeElse:     lineTypeElse,
		tokenTypeEnd:      lineTypeFor,
		tokenTypeFunction: lineTypeFunction,
		tokenTypeReturn:   lineTypeReturn,
	}
)

// Line is a struct that contains
// a `Lexer` and lineType
type Line struct {
	lexer *Lexer
	typ   lineType
}

// String returns a string representation of a `Line`
func (l *Line) String() string {
	return l.lexer.String()
}

// Run evaluates the `NodeStream`
// produced by the `Lexer`
func (l *Line) Run() Node {
	return l.NodeStream().Evaluate()
}

// NodeStream returns a NodeStream
// from the `Lexer`
func (l *Line) NodeStream() *NodeStream {
	return NewNodeStreamFromLexer(l.lexer)
}

// lineType is a int
// representing a `Line` type
type lineType int

const (
	lineTypeBasic = iota
	lineTypeIf
	lineTypeFor
	lineTypeFunction
	lineTypeEnd
	lineTypeReturn
	lineTypeElse
	lineTypeEOF
	lineTypeBlank
)

// NewLineReader returns a new `LineReader`
func NewLineReader(buffer string) *LineReader {
	lr := new(LineReader)
	lr.strLines = strings.Split(buffer, "\n")
	lr.lineNum = 1
	lr.size = 0
	lr.nLines = len(lr.strLines)
	return lr
}

// ReadLines turns the string slice into `Line` types
func (lr *LineReader) ReadLines() *LineReader {
	index := 0

	for line := lr.next(); line.typ != lineTypeEOF; line = lr.next() {
		if line.typ != lineTypeBlank {
			lr.lines = append(lr.lines, line)
			lr.size++
			index++
		}

		if line.typ == lineTypeFunction {
			lr.getFunction(line)
		}
	}

	lr.pos = 0
	return lr
}

// getFunction read the lines in a function declaration block
// separately from the other blocks and adds a new `UserFunction`
// to the Scope
func (lr *LineReader) getFunction(line *Line) {
	depth := 1
	f := ParseUserFunction(line.NodeStream())
	newReader := new(LineReader)

	for line := lr.next(); line.typ != lineTypeEOF; line = lr.next() {
		if line.typ == lineTypeIf || line.typ == lineTypeFunction || line.typ == lineTypeFor {
			depth++
		}

		if line.typ == lineTypeEnd {
			depth--
			if depth == 0 {
				break
			}
		}

		if line.typ == lineTypeFunction {
			log.Fatal("Cannot have function in function")
		}

		newReader.lines = append(newReader.lines, line)
		newReader.size++
	}

	f.block = NewBlockBuilder(newReader).Build()

	if scopeIsInitalized {
		SetFunc(f.name, f)
	}

	println(len(f.block.blocks))
}

// NextLine returns the next `Line` from
// the `LineReader`
func (lr *LineReader) NextLine() *Line {
	if !lr.HasNextLine() {
		return lineEOF
	}

	line := lr.lines[lr.pos]
	lr.pos++
	lr.lineNum++
	return line
}

// Backup decrements the line position and
// the line number
func (lr *LineReader) Backup() *LineReader {
	lr.pos--
	lr.lineNum--
	return lr
}

// HasNextLine determines if another line
// can be read from the `LineReaderstructs`
func (lr *LineReader) HasNextLine() bool {
	return lr.pos < lr.size
}

// next gets the next string from the `LineReader`'s
// string slice and returns a new `Line` created
// from that string
func (lr *LineReader) next() *Line {
	if lr.pos >= lr.nLines {
		return lineEOF
	}

	line := new(Line)

	if !shouldSkipLine(lr.strLines[lr.pos]) {
		line.lexer = NewLexer(lr.strLines[lr.pos])
		line.lexer.Lex()

		if typ, ok := tokenLineKey[line.lexer.FirstItem().typ]; ok {
			line.typ = typ
		} else {
			line.typ = lineTypeBasic
		}
	} else {
		line.typ = lineTypeBlank
	}

	lr.pos++

	return line
}

// shouldSkipLine determines if a line should be skipped
// (not read by the lexer) if it begines with the
// comment identifier or contains only whitespace
func shouldSkipLine(strLine string) bool {
	noWhiteSpace := strings.TrimSpace(strLine)
	return len(noWhiteSpace) == 0 ||
		strings.HasPrefix(noWhiteSpace, "--")
}
