package blast

import "strings"

type lineType int

const (
	lineTypeUnkown = iota
	lineTypeIf
	lineTypeFunction
	lineTypeBasic
	lineTypeEnder
)

type blockType int

const (
	blockTypeBasic = iota
	blockTypeFunction
	blockTypeIf
)

type code []*block

type block struct {
	blocks []*block
	line   *line
	t      blockType
	nLines int
}

var lineToBlock = map[lineType]blockType{
	lineTypeIf:       blockTypeIf,
	lineTypeFunction: blockTypeFunction,
}

type line struct {
	ts *tokenStream
	t  lineType
}

func newCode(strCode string) *code {
	var i int
	c := new(code)
	strLines := strings.Split(strCode, "\n")
	lenStrLines := len(strLines)

	for i < lenStrLines {
		if strings.Trim(strLines[i], " ") == "" {
			i++
			continue
		}

		block := newBlock(strLines[i:])
		*c = append(*c, block)
		i += block.nLines
	}

	return c
}

func (c *code) string() string {
	str := ""

	for _, b := range *c {
		str += b.string() + "\n"
	}

	return str
}

func newBlock(strLines []string) *block {
	if newLine(strLines[0]).t == lineTypeBasic {
		return newBasicBlock(newLine(strLines[0]))
	}

	return newMultiBlock(strLines)
}

func newMultiBlock(strLines []string) *block {
	block := new(block)
	var line *line

	for i, strLine := range strLines {
		block.nLines++

		switch line = newLine(strLine); line.t {
		case lineTypeEnder:
			return block
		case lineTypeBasic:
			block.add(newBasicBlock(line))
		default:
			block.line = line
			block.t = lineToBlock[line.t]
			block.add(newMultiBlock(strLines[i+1:]))
		}
	}

	return block
}

func newBasicBlock(ln *line) *block {
	return &block{
		line:   ln,
		t:      blockTypeBasic,
		nLines: 1,
	}
}

func (b *block) add(block *block) {
	if block.line != nil {
		b.blocks = append(b.blocks, block)
	}
}

func (b *block) string() string {
	str := ""

	if b.line == nil {
		str += ""
	} else {
		if b.t == blockTypeBasic {
			return b.line.ts.string()
		}
	}

	for _, b := range b.blocks {
		str += b.string() + "\n"
	}

	return str
}

func newLine(strCode string) *line {
	ls := newLexemeStream(strCode)
	ts := newTokenStream(ls)

	switch ts.get(0).t {
	case tokenTypeIf:
		return &line{ts, lineTypeIf}
	case tokenFunction:
		return &line{ts, lineTypeFunction}
	case tokenTypeEnd:
		return &line{ts, lineTypeEnder}
	}

	return &line{ts, lineTypeBasic}
}
