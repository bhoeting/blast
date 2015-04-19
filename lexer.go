package blast

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Lex(input string) *Lexer {
	lexer := NewLexer(input)
	lexer.Lex()
	return lexer
}

type Lexer struct {
	pos        int
	width      int
	parenDepth int
	itemPos    int
	text       string
	curr       string
	items      []*Item
	state      lexerState
}

type lexerState int

type lexerFn func(l *Lexer) lexerFn

const (
	lexerStateInQuotes = iota
	lexerStateDefault
)

type Item struct {
	typ  itemType
	pos  int
	text string
}

func NewItem(text string, pos int, typ itemType) *Item {
	item := new(Item)
	item.pos = pos
	item.text = text
	item.typ = typ
	return item
}

func (i *Item) String() string {
	return fmt.Sprintf("%s", i.text)
}

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
)

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
	itemEOF = &Item{
		typ:  itemTypeEOF,
		pos:  -1,
		text: "<NULL>",
	}

	reservedKey = map[string]itemType{
		"if":       itemTypeIf,
		"else":     itemTypeElse,
		"true":     itemTypeBool,
		"false":    itemTypeBool,
		"return":   itemTypeReturn,
		"function": itemTypeFunction,
		"end":      itemTypeEnd,
	}
)

const eof = -1

func NewLexer(text string) *Lexer {
	lexer := new(Lexer)
	lexer.state = lexerStateDefault
	lexer.text = text
	return lexer
}

func (l *Lexer) String() string {
	str := ""

	for _, item := range l.items {
		str += item.String()
	}

	return str
}

func (l *Lexer) Lex() lexerFn {

	if l.HasNext() {
		r := l.Peek()
		if r == '"' {
			return l.LexString()
		}

		if unicode.IsSpace(r) {
			l.Next()
			return l.Lex()
		}

		switch r {
		case '(':
			l.Consume(l.Next())
			l.PushItem(itemTypeOpenParen)
			return l.Lex()
		case ')':
			l.Consume(l.Next())
			l.PushItem(itemTypeCloseParen)
			return l.Lex()
		case ',':
			l.Consume(l.Next())
			l.PushItem(itemTypeComma)
			return l.Lex()
		}

		if isAlphaNumeric(r) {
			return l.LexIdentifier()
		}

		if r == '.' || r == '-' || unicode.IsNumber(r) {
			return l.LexNumber()
		}

		if isOperatorPiece(r) {
			return l.LexOperator()
		}
	}

	return l.Stop()
}

func (l *Lexer) LexIdentifier() lexerFn {
	var next rune

	l.ConsumeUntilNotValid(func(r rune) bool {
		next = r
		return isAlphaNumeric(r) || r == '.'
	})

	_, err := strconv.ParseFloat(l.curr, 64)

	if err == nil {
		return l.LexNumber()
	}

	l.PushItem(parseItemTypeFromString(l.curr))
	return l.Lex()
}

func (l *Lexer) LexNumber() lexerFn {
	l.ConsumeUntilNotValid(func(r rune) bool {

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

	l.PushItem(itemTypeNum)
	return l.Lex()
}

func (l *Lexer) LexOperator() lexerFn {
	l.ConsumeUntilNotValid(func(r rune) bool {
		return isOperatorPiece(r)
	})

	if !isOperator(l.curr) {
		l.Errorf("Invalid operator %s", l.curr)
	}

	l.PushItem(itemTypeOperator)
	return l.Lex()
}

func (l *Lexer) LexString() lexerFn {
	l.Next()
	l.ConsumeUntilNotValid(func(r rune) bool {
		return r != '"'
	})

	l.PushItem(itemTypeString)
	l.Next()
	return l.Lex()
}

func (l *Lexer) ConsumeUntilNotValid(isValid func(r rune) bool) {
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

func (l *Lexer) Stop() lexerFn {
	l.curr = ""
	l.pos = -1
	return nil
}

func (l *Lexer) PushItem(typ itemType) *Lexer {
	item := NewItem(l.curr, l.pos, typ)
	l.items = append(l.items, item)
	l.curr = ""
	return l
}

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

func (l *Lexer) Next() rune {
	if !l.HasNext() {
		return eof
	}

	r, width := utf8.DecodeRuneInString(l.text[l.pos:])
	l.pos += width
	l.width = width
	return r
}

func (l *Lexer) HasNext() bool {
	return l.pos < len(l.text)
}

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

func (l *Lexer) Backup() *Lexer {
	l.pos -= l.width
	return l
}

func (l *Lexer) Consume(r rune) *Lexer {
	l.curr += string(r)
	return l
}

func (l *Lexer) Flush() *Lexer {
	l.curr = ""
	return l
}

func (l *Lexer) NextItem() *Item {
	if !l.HasNextItem() {
		return itemEOF
	}

	item := l.items[l.itemPos]
	l.itemPos++
	return item
}

func (l *Lexer) PeekItem() *Item {
	i := l.NextItem()
	l.BackupItem()
	return i
}

func (l *Lexer) FirstItem() *Item {
	if !l.HasNextItem() {
		return itemEOF
	}

	i := l.NextItem()
	l.BackupItem()

	return i
}

func (l *Lexer) BackupItem() *Lexer {
	l.itemPos--
	return l
}

func (l *Lexer) HasNextItem() bool {
	return l.itemPos < len(l.items)
}

func (l *Lexer) Errorf(errFmt string, args ...interface{}) {
	log.Fatalf("Blast error at pos "+string(l.pos)+": "+errFmt, args...)
}

func parseItemTypeFromString(text string) itemType {
	if typ, ok := reservedKey[text]; ok {
		return typ
	}

	return itemTypeIdentifier
}

func isOperator(strOp string) bool {
	switch strOp {
	case "+", "-", "*", "/", "=", "==", "&&", "||", "^", "<", "<=", ">", ">=", "!=":
		return true
	}

	return false
}

func isOperatorPiece(r rune) bool {
	switch r {
	case '+', '-', '/', '*', '=', '&', '|', '^', '<', '>', '!':
		return true
	}

	return false
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsNumber(r) || unicode.IsLetter(r)
}

func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}
