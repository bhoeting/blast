package blast

import "testing"

func TestToken(t *testing.T) {
	iToken := NewToken("888")
	sToken := NewToken("\"derp\"")
	oToken := NewToken("+")
	vToken := NewToken("x")

	if iToken.t != TokenTypeNumber {
		t.Fatal("Failed to get number from token")
		if iToken.Int() != 888 {
			t.Fatalf("Failed to get value, expected %v, got %v", 888, iToken.Int())
		}
	}

	if oToken.t != TokenTypeOperator {
		t.Fatal("Failed to get operator from token")
		if iToken.Op() != OpTypeAddition {
			t.Fatalf("Failed to get value, expected %v, got %v", OpTypeAddition, iToken.Op())
		}
	}

	if sToken.t != TokenTypeString {
		t.Fatal("Failed to get string from token")
		if iToken.Str() != "derp" {
			t.Fatalf("Failed to get value, expected %v, got %v", "derp", iToken.Str())
		}
	}

	if vToken.t != TokenTypeVariable {
		t.Fatal("Failed to get variable from token")
		if iToken.Var() != "x" {
			t.Fatalf("Failed to get value, expected %v, got %v", "x", iToken.Var())
		}
	}
}

func TestArrToken(t *testing.T) {
	arr := NewTokenArray("888 + \"de rp\"")

	if len(arr) != 3 {
		t.Fatalf("Failed to get correct length of TokenArray %v", len(arr))
	}

	if arr[0].t != TokenTypeNumber || arr[1].t != TokenTypeOperator || arr[2].t != TokenTypeString {
		t.Fatalf("Failed to get correct types in token array %v %v %v", arr[0].t, arr[1].t, arr[2].t)
	}

	if arr[0].Int() != int64(888) || arr[1].Op() != OpTypeAddition {
		t.Logf("%v\n%v", arr[0].Int(), arr[1].Op())
		t.Fatalf("%v", "Failed to get correct values in Token Array")
	}
}
