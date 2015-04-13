package blast

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (b *Blast) importBuiltinFuncs() {
	b.addFunction("print", newBuiltinFunction("print", builtinPrint, 1))
	b.addFunction("println", newBuiltinFunction("println", builtinPrintln, 1))
	b.addFunction("max", newBuiltinFunction("max", builtinMax, 2))
	b.addFunction("abs", newBuiltinFunction("abs", builtinAbs, 1))
	b.addFunction("get", newBuiltinFunction("get", builtinGet, 1))
	b.addFunction("modulus", newBuiltinFunction("modulus", builtinModulus, 2))
}

func builtinPrint(tokens []*token) *token {
	if tokens[0].t == tokenTypeString {
		fmt.Printf("%v", tokens[0].data)
	} else {
		fmt.Printf("%v", tokens[0].string())
	}

	return tokens[0]
}

func builtinPrintln(tokens []*token) *token {
	if tokens[0].t == tokenTypeString {
		fmt.Printf("%v\n", tokens[0].data)
	} else {
		fmt.Println(tokens[0].string())
	}

	return tokens[0]
}

func builtinAbs(tokens []*token) *token {
	if tokens[0].number() < 0 {
		tokens[0].data = tokens[0].number() * -1
	}
	return tokens[0]
}

func builtinMax(tokens []*token) *token {
	if tokens[0].number() > tokens[1].number() {
		return tokens[0]
	}
	return tokens[1]
}

func builtinModulus(tokens []*token) *token {
	result := tokens[0].integer() % tokens[1].integer()
	return newToken(result, 0, 0, tokenTypeBoolean)
}

func builtinGet(tokens []*token) *token {
	token := newToken("", 0, 0, tokenTypeString)
	resp, err := http.Get(tokens[0].str())

	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		token.data = string(body)
	}

	token.data = "Could not make request"
	return token
}
