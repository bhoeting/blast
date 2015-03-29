package blast

import (
	"log"
	"strings"
)

func compareTokens(t1 *token, t2 *token, opType opType) *token {
	var result bool

	switch opType {
	case opTypeEqualTo:
		result = t1.number() == t2.number()
	case opTypeLessThan:
		result = t1.number() < t2.number()
	case opTypeGreaterThan:
		result = t1.number() > t2.number()
	case opTypeLessThanOrEqualTo:
		result = t1.number() <= t2.number()
	case opTypeGreaterThanOrEqualTo:
		result = t1.number() >= t2.number()
	default:
		result = true
	}

	return newToken(result, 0, 0, tokenTypeBoolean)
}

func assignTokens(t1 *token, t2 *token) *token {
	if t1.t == tokenTypeVar && t2.t != tokenTypeVar {
		v := B.setVariable(t1.data.(string), t2.data)
		return v.toToken()
	}

	log.Fatalf("%v", "Left token must be a variable")
	return tokenNull
}

func addTokens(t1 *token, t2 *token) *token {
	if t1.t == tokenTypeString {
		if t2.t == tokenTypeString {
			return newToken(t1.str()+t2.str(), 0, 0, tokenTypeString)
		}
		return newToken(t1.str()+t2.string(), 0, 0, tokenTypeString)
	}

	if t2.t == tokenTypeString {
		return newToken(t1.string()+t2.str(), 0, 0, tokenTypeString)
	}

	if isIntInt(t1, t2) {
		return newToken(t1.integer()+t2.integer(), 0, 0, tokenTypeInt)
	}

	if isIntFloat(t1, t2) {
		return newToken(float64(t1.integer())+t2.float(), 0, 0, tokenTypeFloat)
	}

	if isFloatFloat(t1, t2) {
		return newToken(t1.float()+t2.float(), 0, 0, tokenTypeFloat)
	}

	if isFloatInt(t1, t2) {
		return newToken(t1.float()+float64(t2.integer()), 0, 0, tokenTypeFloat)
	}

	return newToken(t1.string()+t2.string(), 0, 0, tokenTypeString)
}

func subtractTokens(t1 *token, t2 *token) *token {
	if t1.t == tokenTypeString || t2.t == tokenTypeString {
		// todo: throw err
	}

	if isIntInt(t1, t2) {
		return newToken(t1.integer()-t2.integer(), 0, 0, tokenTypeInt)
	}

	if isIntFloat(t1, t2) {
		return newToken(float64(t1.integer())-t2.float(), 0, 0, tokenTypeFloat)
	}

	if isFloatFloat(t1, t2) {
		return newToken(t1.float()-t2.float(), 0, 0, tokenTypeFloat)
	}

	if isFloatInt(t1, t2) {
		return newToken(t1.float()-float64(t2.integer()), 0, 0, tokenTypeFloat)
	}

	return newToken(t1.string()+t2.string(), 0, 0, tokenTypeString)
}

func multiplyTokens(t1 *token, t2 *token) *token {

	if isIntInt(t1, t2) {
		return newToken(t1.integer()*t2.integer(), 0, 0, tokenTypeInt)
	}

	if isIntFloat(t1, t2) {
		return newToken(float64(t1.integer())*t2.float(), 0, 0, tokenTypeFloat)
	}

	if isIntString(t1, t2) {
		return newToken(strings.Repeat(t2.str(), t1.integer()), 0, 0, tokenTypeString)
	}

	if isFloatFloat(t1, t2) {
		return newToken(t1.float()*t2.float(), 0, 0, tokenTypeFloat)
	}

	if isFloatInt(t1, t2) {
		return newToken(t1.float()*float64(t2.integer()), 0, 0, tokenTypeFloat)
	}

	if isStringInt(t1, t2) {
		return newToken(strings.Repeat(t1.str(), t2.integer()), 0, 0, tokenTypeString)
	}

	return tokenNull
}

func divideTokens(t1 *token, t2 *token) *token {
	if isIntInt(t1, t2) {
		return newToken(t1.integer()/t2.integer(), 0, 0, tokenTypeInt)
	}

	if isIntFloat(t1, t2) {
		return newToken(float64(t1.integer())/t2.float(), 0, 0, tokenTypeFloat)
	}

	if isFloatFloat(t1, t2) {
		return newToken(t1.float()/t2.float(), 0, 0, tokenTypeFloat)
	}

	if isFloatInt(t1, t2) {
		return newToken(t1.float()/float64(t2.integer()), 0, 0, tokenTypeFloat)
	}

	return tokenNull
}

func isIntInt(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeInt && t2.t == tokenTypeInt
}

func isIntString(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeInt && t2.t == tokenTypeString
}

func isIntFloat(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeInt && t2.t == tokenTypeFloat
}

func isIntVar(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeString && t2.t == tokenTypeVar
}

func isFloatFloat(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeFloat && t2.t == tokenTypeFloat
}

func isFloatInt(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeFloat && t2.t == tokenTypeInt
}

func isFloatVar(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeFloat && t2.t == tokenTypeVar
}

func isFloatString(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeFloat && t2.t == tokenTypeString
}

func isStringString(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeString && t2.t == tokenTypeString
}

func isStringInt(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeString && t2.t == tokenTypeInt
}

func isStringFloat(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeString && t2.t == tokenTypeInt
}

func isStringVar(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeString && t2.t == tokenTypeVar
}

func isVarVar(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeVar && t2.t == tokenTypeVar
}

func isVarInt(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeVar && t2.t == tokenTypeInt
}

func isVarFloat(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeVar && t2.t == tokenTypeFloat
}

func isVarString(t1 *token, t2 *token) bool {
	return t1.t == tokenTypeVar && t2.t == tokenTypeString
}
