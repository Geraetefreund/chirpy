package auth

import (
	"github.com/alexedwards/argon2id"
	"testing"
)

func TestArgon2id(t *testing.T) {
	password := "pa$$word"
	hash, _ := argon2id.CreateHash(password, argon2id.DefaultParams)

	match, _ := argon2id.ComparePasswordAndHash(password, hash)

	if !match {
		t.Fatal("password did not match hash")
	}
}
