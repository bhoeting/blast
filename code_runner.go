package blast

func (c *code) run() {
	for _, b := range *c {
		b.run()
	}
}

func (b *block) run() string {
	switch b.t {
	case blockTypeBasic:
		return b.line.run()
	case blockTypeIf:
		runIfBlock(b)
	}

	return ""
}

func runIfBlock(b *block) {
	b.line.ts.tokens = b.line.ts.tokens[1:]

	if b.line.ts.parse().boolean() == true {
		runBlocks(b.blocks)
	}
}

func runBlocks(blocks []*block) {
	for _, b := range blocks {
		b.run()
	}
}

func (ln *line) run() string {
	return ln.ts.parse().string()
}
