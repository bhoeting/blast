package blast

type Blocks []*Block

func ParseCode(code string) *Block {
	lr := NewLineReader(code)
	lr.ReadLines()

	bb := NewBlockBuilder(lr)
	return bb.Build()
}

func (blocks *Blocks) Add(b *Block) {
	*blocks = append(*blocks, b)
}

type Block struct {
	parent *Block
	blocks Blocks
	line   *Line
	typ    blockType
}

type BlockBuilder struct {
	lr    *LineReader
	depth int
	block *Block
}

type blockType int

const (
	blockTypeBasic = iota
	blockTypeIf
	blockTypeFunction
	blockTypeMain
	blockTypeFor
)

var lineBlockTypeKey = map[lineType]blockType{
	lineTypeFunction: blockTypeFunction,
	lineTypeBasic:    blockTypeBasic,
	lineTypeIf:       blockTypeIf,
	lineTypeReturn:   blockTypeBasic,
	lineTypeFor:      blockTypeFor,
}

func NewBlock(parent *Block, line *Line) *Block {
	b := new(Block)
	b.parent = parent
	b.line = line
	b.typ = lineBlockTypeKey[line.typ]
	return b
}

func NewBlockBuilder(lr *LineReader) *BlockBuilder {
	bb := new(BlockBuilder)
	bb.lr = lr
	bb.depth = 0
	bb.block = &Block{
		typ: blockTypeMain,
	}
	return bb
}

func (b *Block) String() string {
	str := "\n"
	switch b.typ {
	case blockTypeBasic:
		return b.line.String()
	case blockTypeIf, blockTypeFunction:
		str = b.line.String()
		for _, block := range b.blocks {
			str += "\t" + block.String() + "\n"
		}
		str += "end"
	case blockTypeMain:
		for _, block := range b.blocks {
			str += block.String() + "\n"
		}
	}

	return str
}

func (bb *BlockBuilder) Build() *Block {
	for bb.lr.HasNextLine() {
		line := bb.lr.NextLine()

		switch line.typ {
		case lineTypeEnd:
			bb.depth--
			bb.block = bb.block.parent
		case lineTypeBasic, lineTypeReturn:
			bb.block.blocks.Add(NewBlock(bb.block, line))
		case lineTypeIf, lineTypeFunction, lineTypeFor:
			bb.depth++
			newBlock := NewBlock(bb.block, line)
			bb.block.blocks.Add(newBlock)
			bb.block = newBlock
		}
	}

	// if bb.depth > 0 {
	// 	log.Fatalf("%d too many ends", bb.depth)
	// }

	// if bb.depth < 0 {
	// 	log.Fatal("unclosed block")
	// }

	return bb.block
}

func (b *Block) RunBlocks() (Token, bool) {
	var token Token
	var returned bool

	for _, block := range b.blocks {
		token, returned = block.Run()
		if returned {
			return token, true
		}
	}
	return token, false
}

func (b *Block) Run() (Token, bool) {
	switch b.typ {
	case blockTypeBasic:
		if b.line.typ == lineTypeReturn {
			ts := b.line.TokenStream()
			ts.Next()
			return ts.Chop().Evaluate(), true
		}
		return b.line.Run(), false
	case blockTypeIf:
		return runIfBlock(b)
	case blockTypeFunction:
		return runFuncBlock(b)
	case blockTypeFor:
		return runForBlock(b)
	}

	return &tokenNil{}, false
}

func runIfBlock(b *Block) (Token, bool) {
	Scopes.New()
	ts := b.line.TokenStream()
	ts.Next()
	condition := ts.Chop()

	if BooleanFromToken(condition.Evaluate()) {
		return b.RunBlocks()
	}

	Scopes.Pop()
	return &tokenNil{}, false
}

func runForBlock(b *Block) (Token, bool) {
	Scopes.New()
	ts := b.line.TokenStream()
	fd := ParseForDeclaration(ts)

	if fd.step == 0 {
		if fd.start < fd.end {
			fd.step = 1
		} else {
			fd.step = -1
		}
	}

	for i := fd.start; i <= fd.end; i += fd.step {
		SetVar(fd.counter.name, NewNumberFromFloat(i))
		if token, returned := b.RunBlocks(); returned {
			return token, true
		}
	}

	Scopes.Pop()
	return &tokenNil{}, false
}

func runFuncBlock(b *Block) (Token, bool) {
	return &tokenNil{}, false
}
