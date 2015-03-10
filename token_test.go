package blast

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	iToken := NewToken("888")
	sToken := NewToken("\"derp\"")
	oToken := NewToken("+")
	vToken := NewToken("x")

	// Test int
	if iToken.t != TokenTypeInteger {
		t.Fatal("Failed to get number from token")
	}
	if iToken.Int() != 888 {
		t.Fatalf("Failed to get value, expected %v, got %v", 888, iToken.Int())
	}

	// Test operator
	if oToken.t != TokenTypeOperator {
		t.Fatal("Failed to get operator from token")
	}
	if oToken.Op() != OpTypeAddition {
		t.Fatalf("Failed to get value, expected %v, got %v", OpTypeAddition, oToken.Op())
	}

	// Test string
	if sToken.t != TokenTypeString {
		t.Fatal("Failed to get string from token")
	}
	if sToken.Str() != "derp" {
		t.Fatalf("Failed to get value, expected %v, got %v", "derp", sToken.Str())
	}

	// Test var
	if vToken.t != TokenTypeVariable {
		t.Fatal("Failed to get variable from token")
	}
	if vToken.Var() != "x" {
		t.Fatalf("Failed to get value, expected %v, got %v", "x", vToken.Var())
	}
}

func TestArrToken(t *testing.T) {
	arr := NewTokenArray("\"de rp\" + 300")

	if len(arr) != 3 {
		t.Fatalf("Failed to get correct length of TokenArray %v", len(arr))
	}

	if arr[0].t != TokenTypeString || arr[1].t != TokenTypeOperator || arr[2].t != TokenTypeInteger {
		t.Fatalf("Failed to get correct types in token array %v %v %v", arr[0].t, arr[1].t, arr[2].t)
	}

	if arr[0].Str() != "de rp" || arr[1].Op() != OpTypeAddition || arr[2].Int() != int64(300) {
		t.Log(arr[0].Str(), arr[1].Op(), arr[2].Int())
		t.Fatalf("%v", "Failed to get correct values in Token Array")
	}
}

func TestTokenHandler(t *testing.T) {
	arr := NewTokenArray(fmt.Sprintf("3 * \"%v\"", "derp"))
	tok := HandleTokens(arr[0], arr[2], arr[1], new(Blast))

	if tok.Str() != "derpderpderp" {
		t.Fatalf("Failed to multiply string by interger, expected %v, got %v", "derpderpderp", tok.Str())
	}

	arr = NewTokenArray(fmt.Sprintf("\"%v\" + 3000", "busters "))
	tok = HandleTokens(arr[0], arr[2], arr[1], new(Blast))

	if tok.Str() != "busters 3000" {
		t.Fatalf("Failed to add string and interger, expected %v, got %v", "busters 3000", tok.Str())
	}

	arr = NewTokenArray("2 + 2")
	tok = HandleTokens(arr[0], arr[2], arr[1], new(Blast))
	if tok.Int() != int64(4) {
		t.Fatalf("Failed to add integer and integer, wanted %v, got %v", tok.Int())
	}
}
