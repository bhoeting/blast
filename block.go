package blast

// Blocks is a slice of `Block`
type Blocks []*Block

// ParseCode turns a string of code
// into a Block
func ParseCode(code string) *Block {
	lr := NewLineReader(code)
	lr.ReadLines()

	bb := NewBlockBuilder(lr)
	return bb.Build()
}

// Add adds a `Block` to the block slice
func (blocks *Blocks) Add(b *Block) {
	*blocks = append(*blocks, b)
}

// Block is a recursive struct to represent sections
// of code and handle scoping. See block types
// for more detail
type Block struct {
	parent *Block
	blocks Blocks
	line   *Line
	typ    blockType
}

// BlockBuilder is a struct that turns
// the lines from a `LineReader` into
// a block
type BlockBuilder struct {
	lr    *LineReader
	depth int
	block *Block
}

// blockType is an int
// representing a block type
type blockType int

const (
	// A single line, independent from other lines.
	// Ex: `x = 3`
	blockTypeBasic = iota
	// An if block, with child line
	// being its own block
	blockTypeIf
	// A function declaration block,
	// with each child line being
	// ins own block
	blockTypeFunction
	// The first block, with each
	// line being its own block
	blockTypeMain
	// A for loop declaration, with
	// each child line being its
	// own block
	blockTypeFor
)

// lineBlockTypeKey is a map to help
// create a `Block` from a `Line`
var lineBlockTypeKey = map[lineType]blockType{
	lineTypeFunction: blockTypeFunction,
	lineTypeBasic:    blockTypeBasic,
	lineTypeIf:       blockTypeIf,
	lineTypeReturn:   blockTypeBasic,
	lineTypeFor:      blockTypeFor,
}

// NewBlock returns a new `Block`
func NewBlock(parent *Block, line *Line) *Block {
	b := new(Block)
	b.parent = parent
	b.line = line
	b.typ = lineBlockTypeKey[line.typ]
	return b
}

// NewBlockBuilder returns a new BlockBuilder
func NewBlockBuilder(lr *LineReader) *BlockBuilder {
	bb := new(BlockBuilder)
	bb.lr = lr
	bb.depth = 0
	bb.block = &Block{
		typ: blockTypeMain,
	}
	return bb
}

// String returns a string representation
// of a block
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

// Build creates the block structure
func (bb *BlockBuilder) Build() *Block {
	for bb.lr.HasNextLine() {
		line := bb.lr.NextLine()

		switch line.typ {
		case lineTypeEnd:
			bb.depth--
			bb.block = bb.block.parent
		case lineTypeBasic, lineTypeReturn:
			bb.block.blocks.Add(NewBlock(bb.block, line))
		case lineTypeIf, lineTypeFor:
			bb.depth++
			newBlock := NewBlock(bb.block, line)
			bb.block.blocks.Add(newBlock)
			bb.block = newBlock
		}
	}

	return bb.block
}

// RunBlocks runs each `Block` in the `Block` slice
func (b *Block) RunBlocks() (Node, bool) {
	var node Node
	var returned bool

	for _, block := range b.blocks {
		node, returned = block.Run()
		if returned {
			return node, true
		}
	}
	return node, false
}

// Run executes the line stored in the block
func (b *Block) Run() (Node, bool) {
	switch b.typ {
	case blockTypeBasic:
		if b.line.typ == lineTypeReturn {
			ns := b.line.NodeStream()
			ns.Next()
			return ns.Chop().Evaluate(), true
		}
		return b.line.Run(), false
	case blockTypeIf:
		return runIfBlock(b)
	case blockTypeFunction:
		return runFuncBlock(b)
	case blockTypeFor:
		return runForBlock(b)
	}

	return &nodeNil{}, false
}

// runIfBlock runs an if `Block`
func runIfBlock(b *Block) (Node, bool) {
	Scopes.New()
	ns := b.line.NodeStream()
	ns.Next()
	condition := ns.Chop()

	if BooleanFromNode(condition.Evaluate()) {
		return b.RunBlocks()
	}

	Scopes.Pop()
	return &nodeNil{}, false
}

// runForBlock runs a for `Block`
func runForBlock(b *Block) (Node, bool) {
	Scopes.New()
	ns := b.line.NodeStream()
	fd := ParseForDeclaration(ns)

	if fd.step == 0 {
		if fd.start < fd.end {
			fd.step = 1
		} else {
			fd.step = -1
		}
	}

	for i := fd.start; i <= fd.end; i += fd.step {
		SetVar(fd.counter.name, NewNumberFromFloat(i))
		if node, returned := b.RunBlocks(); returned {
			return node, true
		}
	}

	Scopes.Pop()
	return &nodeNil{}, false
}

// runFuncBlock is currently not a thing
func runFuncBlock(b *Block) (Node, bool) {
	return &nodeNil{}, false
}
