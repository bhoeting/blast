package blast

import "fmt"

type opType int

const (
	opTypeAddition opType = iota
	opTypeSubtraction
	opTypeMultiplication
	opTypeDivision
	opTypeParen
)

var precedenceMap map[opType]int = map[opType]int{
	opTypeAddition:       1,
	opTypeDivision:       2,
	opTypeSubtraction:    1,
	opTypeMultiplication: 2,
}

const (
	parenTypeOpen  = 1
	parenTypeClose = 2
)

type component interface {
}

type str struct {
	data string
}

type integer struct {
	data int
}

type float struct {
	data float64
}

type paren struct {
	t int
}

type operator struct {
	precedence int
	t          opType
}

type componentStream struct {
	items []component
	size  int
}

func newComponentStream() *componentStream {
	cs := new(componentStream)
	cs.size = 0
	cs.items = make([]component, 0)
	return cs
}

func (cs *componentStream) add(c component) {
	cs.items = append(cs.items, c)
	cs.size++
}

func (cs *componentStream) each(f func(c component, i int)) {
	for i, c := range cs.items {
		f(c, i)
	}
}

func newStr(data string) *str {
	s := new(str)
	s.data = data
	return s
}

func newInteger(data int) *integer {
	i := new(integer)
	i.data = data
	return i
}

func newFloat(data float64) *float {
	f := new(float)
	f.data = data
	return f
}

func newParen(data int) *paren {
	p := new(paren)
	p.t = data
	return p
}

func newOperator(t opType) *operator {
	o := new(operator)
	o.t = t
	o.precedence = precedenceMap[o.t]
	return o
}

func componentFromToken(t *token) component {
	switch t.t {
	case tokenTypeString:
		if data, ok := t.data.(string); ok {
			return newStr(data)
		}
	case tokenTypeInt:
		if data, ok := t.data.(int); ok {
			return newInteger(data)
		}
	case tokenTypeFloat:
		if data, ok := t.data.(float64); ok {
			return newFloat(data)
		}
	case tokenTypeParen:
		if data, ok := t.data.(int); ok {
			return newParen(data)
		}
	case tokenTypeOp:
		if data, ok := t.data.(opType); ok {
			return newOperator(data)
		}
	}

	return newInteger(0)
}

func newComponentStreamFromTokenStream(ts *tokenStream) *componentStream {
	cs := newComponentStream()

	for _, t := range ts.tokens {
		c := componentFromToken(t)
		cs.add(c)
	}

	return cs
}

func (cs *componentStream) string() string {
	final := ""

	for _, c := range cs.items {
		if d, ok := c.(*operator); ok {
			final += fmt.Sprint(d.string())
		}

		if p, ok := c.(*paren); ok {
			final += fmt.Sprint(p.string())
		}

		if i, ok := c.(*integer); ok {
			final += fmt.Sprint(i.string())
		}

		if s, ok := c.(*str); ok {
			final += fmt.Sprint(s.string())
		}

		if f, ok := c.(*float); ok {
			final += fmt.Sprint(f.string())
		}
	}

	return str
}

func (i *integer) string() string {
	return fmt.Sprintf("%v", i.data)
}

func (f *float) string() string {
	return fmt.Sprintf("%v", f.data)
}

func (s *str) string() string {
	return fmt.Sprintf("%v", s.data)
}

func (p *paren) string() string {
	if p.t == parenTypeClose {
		return fmt.Sprint(closeParenIdentifier)
	} else {
		return fmt.Sprint(openParenIdentifier)
	}
}

func (o *operator) string() string {
	switch o.t {
	case opTypeAddition:
		return fmt.Sprint("+")
	case opTypeMultiplication:
		return fmt.Sprint(multiplicationIdentifier)
	case opTypeDivision:
		return fmt.Sprint(divisionIdentifier)
	case opTypeSubtraction:
		return fmt.Sprint(subtractionIdentifier)
	}

	return fmt.Sprint("Invalid op type")
}
