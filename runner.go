package blast

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func RunFile(fName string) {

	if !strings.HasSuffix(fName, ".blast") {
		fName += ".blast"
	}

	data, err := ioutil.ReadFile(fName)

	if err != nil {
		fmt.Errorf("%s", err)
	}

	InitScope()
	block := ParseCode(string(data))
	block.RunBlocks()
}
