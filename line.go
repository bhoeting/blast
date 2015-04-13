package blast

import (
	"errors"
	"strings"
)

type lineType int

const (
	lineTypeUnkown = iota
	lineTypeIf
	lineTypeFunction
	lineTypeBasic
	lineTypeEnder
)

var (
	errBlankLine = errors.New("Line is blank")
)

// lines is a pointer
// to a slice of lines
type lines []*line

// line is a struct
// that stores
type line struct {
	ts *tokenStream
	ls *lexemeStream
	t  lineType
}

// newLines returns a new lines type
func newLines(strCode string) *lines {
	lines := new(lines)
	strLines := strings.Split(strCode, "\n")
	lineIndex, nLines := 0, len(strLines)

	for lineIndex < nLines {
		if strings.HasPrefix(strLines[lineIndex], commentIdentifier) {
			lineIndex++
			continue
		}

		line, err := newLine(strLines[lineIndex])

		// Skip blank lines
		if err != nil {
			lineIndex++
			continue
		}

		// If a line is a function declaration,
		// declare the function and skip the
		// lines within the function block
		if line.t == lineTypeFunction {
			skip := declareFunction(strLines[lineIndex:])
			lineIndex += skip
		} else {
			lineIndex++
			lines.push(line)
		}
	}

	return lines
}

// run runs all the lines
func (lines *lines) run() *token {
	t := runLines(lines)
	return t
}

func runLines(lines *lines) *token {
	skip := false
	bLevel := 0
	token := new(token)

	for _, line := range *lines {

		if line.t == lineTypeIf {
			bLevel++
			if !skip {
				if line.ts.parse().boolean() == false {
					skip = true
				}
			}
		}

		if line.t == lineTypeBasic {
			if !skip {
				tempToken, shouldReturn := line.run()
				if shouldReturn {
					token = tempToken
					break
				}
			}
		}

		if line.t == lineTypeEnder {
			bLevel--
			if bLevel == 0 {
				skip = false
			}
		}
	}

	lines.reload()

	return token
}

// declareFunction declares a function and returns the amount of lines
// in the block (includeing the declaration and end lines)
func declareFunction(strLines []string) int {
	bLevel := 0
	lineCount := 0
	lines := new(lines)
	var f *userFunction

loop:
	for _, strLine := range strLines {
		line, err := newLine(strLine)

		// Skip blank lines
		if err != nil {
			lineCount++
			continue
		}

		switch line.t {
		case lineTypeFunction:
			f = parseUserFunction(line.ts)
			bLevel++
		case lineTypeBasic:
			lines.push(line)
		case lineTypeIf:
			bLevel++
			lines.push(line)
		case lineTypeEnder:
			bLevel--
			if bLevel == 0 {
				lineCount++
				break loop
			} else {
				lines.push(line)
			}
		}

		lineCount++
	}

	lines.reload()
	f.lines = lines
	B.addFunction(f.name, f)
	return lineCount
}

// reload re-creates the token
// stream for each line
func (lines *lines) reload() {
	for _, line := range *lines {
		line.ts = newTokenStream(line.ls)
	}
}

// push adds a line to the end of the slice
func (lines *lines) push(line *line) {
	*lines = append(*lines, line)
}

// run executes parses a line of code
// and determines if the value is "returned"
func (ln *line) run() (*token, bool) {
	tempTs := *ln.ts

	if tempTs.get(0).t == tokenTypeReturn {
		return tempTs.parse(), true
	}

	return tempTs.parse(), false
}

func (lines *lines) string() string {
	str := ""

	for _, line := range *lines {
		str += line.ts.string()
	}

	return str
}

// newLine returns a new line
func newLine(strLn string) (*line, error) {
	ls := newLexemeStream(strLn)
	ts := newTokenStream(ls)

	if ts.size == 0 {
		return &line{}, errBlankLine
	}

	switch ts.get(0).t {
	case tokenTypeIf:
		return &line{ts, ls, lineTypeIf}, nil
	case tokenTypeFunction:
		return &line{ts, ls, lineTypeFunction}, nil
	case tokenTypeEnd:
		return &line{ts, ls, lineTypeEnder}, nil
	}

	return &line{ts, ls, lineTypeBasic}, nil
}
