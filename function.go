package blast

import "log"

type function struct {
	name   string
	blocks []*blocks
	params []*param
}

type param struct {
	name  string
	token *token
}

func parseFunction(ts *tokenStream) *function {
	f := new(function)
	f.name = ts[1]
	var pName = ""
	defaultValue := new(tokenStream)
	isSettingDefaultValue = false

	for _, t := range ts.tokens[3:] {
		if isSettingDefaultValue {
			defaultValue.add(t)
		}

		switch t.t {
		case tokenTypeParen, tokenTypeComma:
			if defaultValue.size == 0 {
				f.params = append(f.params, &param{pName, 0})
			}

			f.params = append(f.params, &param{pName, defaultValue.parse()})
			isSettingDefaultValue = false
			defaultValue.clear()
			pName = ""
		case tokenTypeVar:
			if pName == "" {
				defaultValue.clear()
				pName = t.data.(string)
			}
		case tokenTypeOperator:
			if t.opType() == opTypeAssignment {
				isSettingDefaultValue = true
			}
		}
	}

	return f
}

func callFunction(fName string, params []*token) *token {
	f, err := B.getFunction(fName)

	if err == nil {
		log.Fatal(err)
	}

	f.call(params)
}

func (f *function) call(tokens []*token) *token {
	for i, t := range tokens {
		f.params[i].token = t
	}

	for _, p := range f.params {
		B.setVariable(p.name, p.token.data)
	}

	runBlocks(f.blocks)
}
