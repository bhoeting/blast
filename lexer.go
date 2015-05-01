package blast

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lex returns the lexed form of the input
func Lex(input string) *Lexer {
	lexer := NewLexer(input)
	lexer.Lex()
	return lexer
}

// Lexer is a struct used
// for lexical analysis.
type Lexer struct {
	pos        int
	width      int
	parenDepth int
	itemPos    int
	text       string
	curr       string
	items      []*Item
}

// lexerFn is a recursive func type
// for lexing specific items.
type lexerFn func(l *Lexer) lexerFn

// Item is lexed from a string,
// precursor to a token.
type Item struct {
	typ  itemType
	pos  int
	text string
}

// NewItem returns a new Item
func NewItem(text string, pos int, typ itemType) *Item {
	item := new(Item)
	item.pos = pos
	item.text = text
	item.typ = typ
	return item
}

// String returns a string representation
// of an Item
func (i *Item) String() string {
	return fmt.Sprintf("%s", i.text)
}

// itemType is an int representing
// the type of an item
type itemType int

const (
	itemTypeEOF = iota
	itemTypeNum
	itemTypeBool
	itemTypeString
	itemTypeOperator
	itemTypeIdentifier
	itemTypeOpenParen
	itemTypeCloseParen
	itemTypeComma
	itemTypeIf
	itemTypeFunction
	itemTypeElse
	itemTypeReturn
	itemTypeEnd
	itemTypeFor
)

// String returns a string representation
// of an itemType
func (typ itemType) String() string {
	switch typ {
	case itemTypeNum:
		return "Number"
	case itemTypeBool:
		return "Bool"
	case itemTypeString:
		return "String"
	case itemTypeOperator:
		return "Operator"
	case itemTypeComma:
		return "Comma"
	case itemTypeOpenParen:
		return "Open paren"
	case itemTypeCloseParen:
		return "Close paren"
	case itemTypeIdentifier:
		return "Identifier"
	}

	return "Unknown"
}

var (
	// itemEOF is an item that is
	// returned when undefined
	// behavior occurs.
	itemEOF = &Item{
		typ:  itemTypeEOF,
		pos:  -1,
		text: "<NULL>",
	}

	// reservedKey is used to get
	// a reserved itemType from
	// a reserved word
	reservedKey = map[string]itemType{
		"if":       itemTypeIf,
		"else":     itemTypeElse,
		"true":     itemTypeBool,
		"false":    itemTypeBool,
		"return":   itemTypeReturn,
		"function": itemTypeFunction,
		"end":      itemTypeEnd,
		"for":      itemTypeFor,
	}
)

// eof is returned from Next()
// when there are no more characters
const eof = -1

// NewLexer returns a new Lexer
// to lex the string `text`
func NewLexer(text string) *Lexer {
	lexer := new(Lexer)
	lexer.text = text
	return lexer
}

// String returns a string representation
// of a lexer.
func (l *Lexer) String() string {
	str := ""

	for _, item := range l.items {
		str += item.String()
	}

	return str
}

// Lex is called when the item
// being lexed is currently unkonwn,
// and returns a lexerFn to lex
// a specific item.
func (l *Lexer) Lex() lexerFn {

	if l.HasNext() {
		r := l.Peek()

		// Lex string literal
		if r == '"' {
			return l.LexString()
		}

		// Skip whitespace
		if unicode.IsSpace(r) {
			l.Next()
			return l.Lex()
		}

		switch r {
		// Lex open paren
		case '(':
			l.Consume(l.Next())
			l.PushItem(itemTypeOpenParen)
			return l.Lex()
		// Lex close paren
		case ')':
			l.Consume(l.Next())
			l.PushItem(itemTypeCloseParen)
			return l.Lex()
		// Lex comma
		case ',':
			l.Consume(l.Next())
			l.PushItem(itemTypeComma)
			return l.Lex()
		}

		// Lex identifier
		if isAlphaNumeric(r) {
			return l.LexIdentifier()
		}

		// Lex number (int or float)
		if r == '.' || r == '-' || unicode.IsNumber(r) {
			return l.LexNumber()
		}

		// Lex operator
		if isOperatorPiece(r) {
			return l.LexOperator()
		}
	}

	return l.Stop()
}

// LexIdentifer lexes an identifier, a
// sequence of alpha or numeric characters
// or a _.
func (l *Lexer) LexIdentifier() lexerFn {
	l.ConsumeWhileValid(func(r rune) bool {
		// If we run into a `.` then we might
		// be lexing a float
		return isAlphaNumeric(r) || r == '.'
	})

	// If the identifier is successfully
	// parsed into a float, we will lex
	// a number instead.
	_, err := strconv.ParseFloat(l.curr, 64)
	if err == nil {
		return l.LexNumber()
	}

	l.PushItem(parseItemTypeFromString(l.curr))
	return l.Lex()
}

// LexNumber lexes a number (float or int)
func (l *Lexer) LexNumber() lexerFn {
	l.ConsumeWhileValid(func(r rune) bool {
		// Check that a negative sign
		// only occurs at the beginning
		if r == '-' {
			return len(l.curr) == 0
		}

		// Check that only one decimal point
		// occurs in the string
		if r == '.' {
			if strings.ContainsRune(l.curr, '.') {
				l.Errorf("Too many decimal points in %s", l.curr+string(r))
			} else {
				return true
			}
		}

		// Check the rune is a number
		return unicode.IsNumber(r)
	})

	// If the only character is a `-`,
	// then it is an operator rather
	// than a negative number.
	if l.curr == "-" {
		l.Backup()
		l.curr = ""
		return l.LexOperator()
	}

	l.PushItem(itemTypeNum)
	return l.Lex()
}

// LexOperator lexes an operator
func (l *Lexer) LexOperator() lexerFn {
	l.ConsumeWhileValid(func(r rune) bool {
		return isOperatorPiece(r)
	})

	if !isOperator(l.curr) {
		l.Errorf("Invalid operator %s", l.curr)
	}

	l.PushItem(itemTypeOperator)
	return l.Lex()
}

// LexString lexes a string literal.
func (l *Lexer) LexString() lexerFn {
	l.Next()
	l.ConsumeWhileValid(func(r rune) bool {
		return r != '"'
	})

	l.PushItem(itemTypeString)
	l.Next()
	return l.Lex()
}

// ConsumeWhileValid consumes the next rune until
// `isValid` returns false
func (l *Lexer) ConsumeWhileValid(isValid func(r rune) bool) {
	var r rune

	for {
		r = l.Next()
		if isValid(r) {
			l.Consume(r)
		} else {
			break
		}
	}

	if l.HasNext() || r != eof {
		l.Backup()
	}
}

// Stop stops the lexical analysis
// by returninig nil instead of
// a lexerFn
func (l *Lexer) Stop() lexerFn {
	l.curr = ""
	l.pos = -1
	return nil
}

// PushItem adds an item to the lexer's
// `Item` slice
func (l *Lexer) PushItem(typ itemType) *Lexer {
	item := NewItem(l.curr, l.pos, typ)
	l.items = append(l.items, item)
	l.curr = ""
	return l
}

// AtTerminator determines if the next
// rune is whitespace, a paren, or comma.
func (l *Lexer) AtTerminator() bool {
	r := l.Peek()

	if unicode.IsSpace(r) {
		return true
	}

	switch r {
	case eof, ',', '(', ')':
		return true
	}

	return false
}

// Next returns the next rune
// from the text and increments
// the position by the rune size
func (l *Lexer) Next() rune {
	if !l.HasNext() {
		return eof
	}

	r, width := utf8.DecodeRuneInString(l.text[l.pos:])
	l.pos += width
	l.width = width
	return r
}

// HasNext determines if another
// rune can be read
func (l *Lexer) HasNext() bool {
	return l.pos < len(l.text)
}

// Peek gets the next rune
// without incrementing the
// position
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

// Backup moves the position back
// by the width of the current
// rune
func (l *Lexer) Backup() *Lexer {
	l.pos -= l.width
	return l
}

// Consume adds rune `r` to curr
func (l *Lexer) Consume(r rune) *Lexer {
	l.curr += string(r)
	return l
}

// Flush empties curr
func (l *Lexer) Flush() *Lexer {
	l.curr = ""
	return l
}

// NextItem returns the next item
// from the `Item` slice
func (l *Lexer) NextItem() *Item {
	if !l.HasNextItem() {
		return itemEOF
	}

	item := l.items[l.itemPos]
	l.itemPos++
	return item
}

// PeekItem returns the next item
// without incrementing the position
func (l *Lexer) PeekItem() *Item {
	i := l.NextItem()
	l.BackupItem()
	return i
}

// FirstItem returns the first item
func (l *Lexer) FirstItem() *Item {
	if !l.HasNextItem() {
		return itemEOF
	}

	i := l.NextItem()
	l.BackupItem()

	return i
}

// BackupItem decrements the item position
func (l *Lexer) BackupItem() *Lexer {
	l.itemPos--
	return l
}

// HasNextItem determines there is
// an Item to return
func (l *Lexer) HasNextItem() bool {
	return l.itemPos < len(l.items)
}

// Errorf reports a lexing error
func (l *Lexer) Errorf(errFmt string, args ...interface{}) {
	log.Fatalf("Blast error at pos "+string(l.pos)+": "+errFmt, args...)
}

// parseItemTypeFromString returns the reserved
// item type that matches the string `text` or
// itemTypeIdentifier if there is no match
func parseItemTypeFromString(text string) itemType {
	if typ, ok := reservedKey[text]; ok {
		return typ
	}

	return itemTypeIdentifier
}

// isOperator determines if a sequence
// of operator chars is a valid operator
func isOperator(strOp string) bool {
	switch strOp {
	case "+", "-", "*", "/", "=", "==", "&&", "||", "^", "<", "<=", ">", ">=", "!=", "->":
		return true
	}

	return false
}

// isOperatorPiece determines if `r`
// is a valid piece or an operator
func isOperatorPiece(r rune) bool {
	switch r {
	case '+', '-', '/', '*', '=', '&', '|', '^', '<', '>', '!':
		return true
	}

	return false
}

// isAlphaNumeric determines if `r` is an
// alapha, number or a '_'
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsNumber(r) || unicode.IsLetter(r)
}

// isAlpha determines if r is alpha or a '_'
func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}
