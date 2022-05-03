package awsngfw

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuthWithKeyInfo(t *testing.T) {
	x := getJwt{
		Expires: 120,
		KeyInfo: &jwtKeyInfo{
			Region: "us-east-1",
			Tenant: "XY",
		},
	}

	b, err := json.Marshal(x)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(b), "KeyInfo") {
		t.Fatalf("KeyInfo not found: %s", b)
	}
}

func TestAuthWithoutKeyInfo(t *testing.T) {
	x := getJwt{
		Expires: 120,
	}

	b, err := json.Marshal(x)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(b), "KeyInfo") {
		t.Fatalf("KeyInfo found: %s", b)
	}
}
