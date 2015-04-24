package blast

import (
	"log"
	"strings"
)

type LineReader struct {
	strLines []string
	lines    []*Line
	size     int
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
		itemTypeFor:      lineTypeFor,
		itemTypeFunction: lineTypeFunction,
		itemTypeReturn:   lineTypeReturn,
	}
)

type Line struct {
	lexer *Lexer
	typ   lineType
}

func (l *Line) String() string {
	return l.lexer.String()
}

func (l *Line) Run() Token {
	return l.TokenStream().Evaluate()
}

func (l *Line) TokenStream() *TokenStream {
	return NewTokenStreamFromLexer(l.lexer)
}

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

func NewLineReader(buffer string) *LineReader {
	lr := new(LineReader)
	lr.strLines = strings.Split(buffer, "\n")
	lr.lineNum = 1
	lr.size = 0
	lr.nLines = len(lr.strLines)
	return lr
}

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

func (lr *LineReader) getFunction(line *Line) {
	depth := 1
	f := ParseUserFunction(line.TokenStream())
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
	return lr.pos < lr.size
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
