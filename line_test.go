package blast

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineReader(t *testing.T) {
	lineReader := NewLineReader(readTestCode(t, "line_reader_test")).ReadLines()

	assert.Equal(t, lineTypeFunction, lineReader.NextLine().typ)
	assert.Equal(t, lineTypeIf, lineReader.NextLine().typ)
	assert.Equal(t, lineTypeReturn, lineReader.NextLine().typ)
	assert.Equal(t, lineTypeEnd, lineReader.NextLine().typ)
	assert.Equal(t, lineTypeReturn, lineReader.NextLine().typ)
	assert.Equal(t, lineTypeEnd, lineReader.NextLine().typ)
	assert.Equal(t, lineTypeBasic, lineReader.NextLine().typ)
	assert.Equal(t, lineTypeBasic, lineReader.NextLine().typ)
}

func readTestCode(t *testing.T, fName string) string {
	if !strings.HasSuffix(fName, ".blast") {
		fName += ".blast"
	}

	data, err := ioutil.ReadFile("test_code/" + fName)

	if err != nil {
		t.Error(err.Error())
	}

	return string(data)
}
