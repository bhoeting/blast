package blast

import "log"

type userFunction struct {
	name   string
	lines  *lines
	params []*param
}

type builtinFunction struct {
	name      string
	f         func(tokens []*token) *token
	numParams int
}

type function interface {
	call(tokens []*token) *token
	fName() string
	nParams() int
}

type param struct {
	name  string
	token *token
}

func parseUserFunction(ts *tokenStream) *userFunction {
	f := new(userFunction)
	f.name = ts.tokens[1].data.(string)
	var pName = ""
	defaultValue := new(tokenStream)
	isSettingDefaultValue := false

	for _, t := range ts.tokens[3:] {
		switch t.t {
		case tokenTypeParen, tokenTypeComma:
			if defaultValue.size == 0 {
				f.params = append(f.params, &param{pName, new(token)})
				pName = ""
			} else {
				f.params = append(f.params, &param{pName, defaultValue.parse()})
				isSettingDefaultValue = false
				defaultValue.clear()
				pName = ""
			}
		case tokenTypeUnkown, tokenTypeVar:
			if pName == "" {
				defaultValue.clear()
				pName = t.data.(string)
			}
		case tokenTypeOperator:
			if t.opType() == opTypeAssignment {
				isSettingDefaultValue = true
			} else {
				defaultValue.add(t)
			}
		default:
			if isSettingDefaultValue {
				defaultValue.add(t)
			}
		}
	}

	return f
}

func callFunction(fName string, tokens []*token) (*token, int) {
	f, err := B.getFunction(fName)
	start := len(tokens) - f.nParams()

	if err != nil {
		log.Fatal(err.Error())
	}

	for i := start; i < len(tokens); i++ {
		tokens[i] = evaluateToken(tokens[i])
	}

	return f.call(tokens[start:]), f.nParams()
}

func (f *userFunction) call(tokens []*token) *token {

	for i, t := range tokens {
		f.params[i].token = t
	}

	for _, p := range f.params {
		B.setVariable(p.name, p.token.data)
	}

	token := f.lines.run()

	for _, p := range f.params {
		B.removeVariable(p.name)
	}

	return token
}

func (f *userFunction) fName() string {
	return f.name
}

func (f *userFunction) nParams() int {
	return len(f.params)
}

func (bf *builtinFunction) call(tokens []*token) *token {
	return bf.f(tokens)
}

func (bf *builtinFunction) fName() string {
	return bf.name
}

func (bf *builtinFunction) nParams() int {
	return bf.numParams
}

func newBuiltinFunction(name string, f func(tokens []*token) *token, numParams int) *builtinFunction {
	bf := new(builtinFunction)
	bf.numParams = numParams
	bf.name = name
	bf.f = f
	return bf
}
