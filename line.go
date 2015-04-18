package blast

import "strings"

type LineReader struct {
	strLines []string
	lines    []*Line
	pos      int
	nLines   int
	lineNum  int
}

var (
	lineEOF = &Line{
		typ: lineTypeEOF,
	}

	itemLineKey = map[itemType]lineType{
		itemTypeEnd:      lineTypeEnd,
		itemTypeIf:       lineTypeIf,
		itemTypeElse:     lineTypeElse,
		itemTypeFunction: lineTypeFunction,
		itemTypeReturn:   lineTypeReturn,
	}
)

type Line struct {
	lexer *Lexer
	typ   lineType
}

type lineType int

const (
	lineTypeBasic = iota
	lineTypeIf
	lineTypeFunction
	lineTypeEnd
	lineTypeReturn
	lineTypeElse
	lineTypeEOF
	lineTypeBlank
)

func NewLineReader(buffer string) *LineReader {
	lr := new(LineReader)
	lr.strLines = strings.Split(buffer, "\n")
	lr.lineNum = 1
	lr.nLines = len(lr.strLines)
	lr.lines = make([]*Line, lr.nLines)
	return lr
}

func (lr *LineReader) ReadLines() *LineReader {
	index := 0

	for line := lr.next(); line.typ != lineTypeEOF; line = lr.next() {
		if line.typ != lineTypeBlank {
			lr.lines[index] = line
			index++
		}
	}

	lr.pos = 0
	return lr
}

func (lr *LineReader) NextLine() *Line {
	if !lr.HasNextLine() {
		return lineEOF
	}

	line := lr.lines[lr.pos]
	lr.pos++
	lr.lineNum++
	return line
}

func (lr *LineReader) Backup() *LineReader {
	lr.pos--
	lr.lineNum--
	return lr
}

func (lr *LineReader) HasNextLine() bool {
	return lr.pos < lr.nLines
}

func (lr *LineReader) next() *Line {
	if lr.pos >= lr.nLines {
		return lineEOF
	}

	line := new(Line)

	if !shouldSkipLine(lr.strLines[lr.pos]) {
		line.lexer = NewLexer(lr.strLines[lr.pos])
		line.lexer.Lex()

		if typ, ok := itemLineKey[line.lexer.FirstItem().typ]; ok {
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

func shouldSkipLine(strLine string) bool {
	noWhiteSpace := strings.TrimSpace(strLine)
	return len(noWhiteSpace) == 0 ||
		strings.HasPrefix(noWhiteSpace, "--")
}
