package crypt

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	rawPwd := "12345"
	encPwd, err := HashPassword(rawPwd)
	if err != nil {
		t.Fatalf("Hashing password was failed, err %v", err)
	}

	if rawPwd == encPwd {
		t.Fatalf("Encoded and raw password should not be equal (pwd: %s)", rawPwd)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	rawPwd := "12345"
	encPwd, err := HashPassword(rawPwd)
	if err != nil {
		t.Fatalf("Hashing password was failed, err %v", err)
	}

	verified := CheckPasswordHash(rawPwd, encPwd)
	if !verified {
		t.Fatalf("Failed to verify the same password(%s) with its hash(%s)", rawPwd, encPwd)
	}
}


func BenchmarkGeneringHashAndChecking(b *testing.B) {
	for i:= 0; i < b.N; i++ {
		rawPwd := fmt.Sprintf("12345-%d", i)
		encPwd, err := HashPassword(rawPwd)
		if err != nil {
			b.Fatalf("Hashing password was failed, err %v", err)
		}
	
		verified := CheckPasswordHash(rawPwd, encPwd)
		if !verified {
			b.Fatalf("Failed to verify the same password(%s) with its hash(%s)", rawPwd, encPwd)
		}
	}
}