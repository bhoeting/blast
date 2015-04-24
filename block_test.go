package blast

import "testing"

// func TestBlockBuilder(t *testing.T) {
// 	block := ParseCode(readTestCode(t, "line_reader_test"))
// 	assert.Equal(t, blockTypeMain, block.typ)
// 	assert.Equal(t, blockTypeFunction, block.blocks[0].typ)
// 	assert.Equal(t, blockTypeIf, block.blocks[0].blocks[0].typ)
// }

func TestCode(t *testing.T) {
	InitScope()
	block := ParseCode(readTestCode(t, "test1.blast"))
	block.RunBlocks()
}
