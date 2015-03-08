package blast

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// RunFile runs the blast code
// contained in a file
func RunFile(fName string) {
	if !strings.HasSuffix(fName, ".blast") {
		fName += ".blast"
	}

	file, _ := os.Open(fName)
	Read(file)
}

// Read reads from an io.Reader
func Read(r io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(r)

	if err != nil {
		return "", err
	}

	code := string(bytes)

	return code, nil
}
